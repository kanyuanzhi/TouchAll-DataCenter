package socket

import (
	"TouchAll-DataCenter/config"
	"TouchAll-DataCenter/models"
	"TouchAll-DataCenter/protocal"
	"TouchAll-DataCenter/utils"
	"TouchAll-DataCenter/websocket"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

// 接口服务器
type SocketServer struct {
	wsClients     *websocket.WsClients
	useMongodb    bool
	useMysql      bool
	socketClients *SocketClients
}

func NewSocketServer(wsClients *websocket.WsClients, socketClients *SocketClients) *SocketServer {
	config := config.NewConfig()
	useMongodb := config.GetValue("mongodb.use").(bool)
	useMysql := config.GetValue("mysql.use").(bool)
	return &SocketServer{
		wsClients:     wsClients,
		useMongodb:    useMongodb,
		useMysql:      useMysql,
		socketClients: socketClients,
	}
}

func (socketServer *SocketServer) Start() {
	config := config.NewConfig()
	port := config.GetSocketConfig().(string)
	l, err := net.Listen("tcp", ":"+port)
	log.Printf("Start the SocketServer of the data center on port %s", port)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		socketServer.socketClients.members[conn] = true
		go socketServer.handleConn(conn)
	}
}

// 处理socket连接，完成数据解包
func (socketServer *SocketServer) handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	tempBuffer := make([]byte, 0)
	readerChannel := make(chan []byte, 512)
	go socketServer.reader(readerChannel, conn)
	for {
		var buffer = make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			// 设备与数据中心连接断开后，依据连接conn找到设备ID，并删除设备状态流中该设备的状态信息
			socketServer.wsClients.EquipmentStatusStream.StatusStreamSyncMap.Delete(socketServer.socketClients.equipmentMembers[conn])
			delete(socketServer.socketClients.members, conn)
			delete(socketServer.socketClients.equipmentMembers, conn)
			fmt.Println(err.Error())
			return
		}
		tempBuffer = protocal.Unpack(append(tempBuffer, buffer[:n]...), readerChannel)
	}
}

func (socketServer *SocketServer) reader(readerChannel chan []byte, conn net.Conn) {
	for {
		select {
		case data := <-readerChannel:
			// 由于此时不知道数据类型，不能转为相应的结构体，因此需要先把json转map，以获取data_type
			var m map[string]interface{}
			json.Unmarshal(data, &m)
			if m["data_type"] == nil {
				return
			}
			switch int(m["data_type"].(float64)) {
			case 10:
				// PeopleAwareness
				var peopleAwareness models.PeopleAwareness
				json.Unmarshal(data, &peopleAwareness)
				personAwarenessData := peopleAwareness.PersonAwarenessData

				go func(personAwarenessData []*models.PersonAwareness) {
					// 存入mongodb
					// TODO: 先存personAwareness，再存peopleAwareness
					if len(personAwarenessData) != 0 && socketServer.useMongodb == true {
						documents := make([]interface{}, 0)
						for i := 0; i < len(personAwarenessData); i++ {
							documents = append(documents, personAwarenessData[i])
						}
						utils.InsertManyRecords(documents, "person_awareness")
					}
				}(personAwarenessData)

				go func(peopleAwareness *models.PeopleAwareness) {
					// 通过websocket推送至前端网页
					socketServer.wsClients.PeopleBroadcast <- peopleAwareness
					socketServer.wsClients.PersonBroadcast <- peopleAwareness.PersonAwarenessData
				}(&peopleAwareness)
			case 11:
				// PersonAwareness，在PeopleAwareness中已处理
				continue
			case 20:
				// EnvironmentAwareness
				continue
			case 30:
				// EquipmentBasicInformationAwareness，设备基本信息，不做推送，只更新数据库

				// 响应给设备的信息
				responseForEquipmentBasicInformation := models.NewResponseForEquipmentBasicInformation()
				responseForEquipmentBasicInformation.DataType = 32
				responseForEquipmentBasicInformation.UseMysql = socketServer.useMysql
				responseForEquipmentBasicInformation.UseMongodb = socketServer.useMongodb

				go func(data []byte, conn net.Conn, responseForEquipmentBasicInformation *models.ResponseForEquipmentBasicInformation) {
					if socketServer.useMysql == true {
						// 根据network_mac查询是否在数据库mysql中注册，
						// 若注册则更新并返回数据库中注册ID、验证状态，若未注册则随机注册一个ID并与"未验证"状态一起返回至设备处理,
						// 设备将该ID添加至发送到数据中心的状态信息中
						var equipmentBasicInformationAwareness models.EquipmentBasicInformationAwareness
						_ = json.Unmarshal(data, &equipmentBasicInformationAwareness)
						isEquipmentVisited, equipmentID, authenticated := utils.IsEquipmentNetworkMacExisted(equipmentBasicInformationAwareness.Network.NetworkMac)
						if isEquipmentVisited {
							responseForEquipmentBasicInformation.EquipmentID = equipmentID
							responseForEquipmentBasicInformation.Authenticated = authenticated
							equipmentBasicInformationAwareness.EquipmentID = equipmentID
							utils.UpdateEquipmentBasicInformation(equipmentBasicInformationAwareness)
						} else {
							// 此时设备未在前端页面注册，交由数据中心自动注册
							rand.Seed(time.Now().Unix())
							temporalEquipmentID := rand.Intn(65535)
							for {
								if !utils.IsEquipmentIDExisted(temporalEquipmentID) {
									break
								} else {
									temporalEquipmentID = rand.Intn(65535)
								}
							}
							responseForEquipmentBasicInformation.EquipmentID = temporalEquipmentID
							responseForEquipmentBasicInformation.Authenticated = 0 // 由数据中心自动注册的设备为未验证设备
							equipmentBasicInformationAwareness.EquipmentID = temporalEquipmentID
							utils.InsertEquipmentBasicInformation(equipmentBasicInformationAwareness)
						}
					} else {
						// 未连接mysql，仅做测试用
						responseForEquipmentBasicInformation := models.NewResponseForEquipmentBasicInformation()
						responseForEquipmentBasicInformation.EquipmentID = 0
						responseForEquipmentBasicInformation.Authenticated = 0 // 由数据中心自动注册的设备为未验证设备
					}
					responseForEquipment, _ := json.Marshal(responseForEquipmentBasicInformation)
					_, err := conn.Write(responseForEquipment)
					if err != nil {
						log.Println(err.Error())
					}
				}(data, conn, responseForEquipmentBasicInformation)
			case 31:
				// EquipmentStatusAwareness，设备状态信息，存入Mongodb并汇总至设备状态流EquipmentStatusStream
				var equipmentStatusAwareness models.EquipmentStatusAwareness
				_ = json.Unmarshal(data, &equipmentStatusAwareness)

				go func(equipmentStatusAwareness *models.EquipmentStatusAwareness) {
					if socketServer.useMongodb == true {
						_, _ = utils.InsertOneRecord(equipmentStatusAwareness, "equipment_awareness")
					}
				}(&equipmentStatusAwareness)

				// 跟踪每个设备的连接情况
				socketServer.socketClients.equipmentMembers[conn] = equipmentStatusAwareness.EquipmentID

				// 将每个设备的状态信息合并到状态流EquipmentStatusStream中，以设备ID作区分
				//socketServer.wsClients.EquipmentStatusStream.StatusStream[equipmentStatusAwareness.EquipmentID] = equipmentStatusAwareness
				// 使用并发安全的字典类型sync.Map
				socketServer.wsClients.EquipmentStatusStream.StatusStreamSyncMap.Store(equipmentStatusAwareness.EquipmentID, equipmentStatusAwareness)

			//go func(equipmentStatusAwareness *models.EquipmentStatusAwareness) {
			//	socketServer.wsClients.EquipmentStatusBroadcast <- equipmentStatusAwareness
			//}(&equipmentStatusAwareness)
			case 50:
				responseForCamera := models.NewResponseForCamera()
				responseForCamera.DataType = 51
				responseForCamera.UseMysql = socketServer.useMysql
				responseForCamera.UseMongodb = socketServer.useMongodb
				go func(data []byte, conn net.Conn, responseForCamera *models.ResponseForCamera) {
					if socketServer.useMysql {
						var camera models.Camera
						_ = json.Unmarshal(data, &camera)
						isCameraVisited, cameraID, authenticated := utils.IsCameraHostExisted(camera.CameraHost)
						if isCameraVisited {
							responseForCamera.CameraID = cameraID
							responseForCamera.Authenticated = authenticated
						} else {
							// 此时监控摄像机未在前端页面注册，交由数据中心自动注册
							rand.Seed(time.Now().Unix())
							temporalCameraID := rand.Intn(65535)
							for {
								if !utils.IsEquipmentIDExisted(temporalCameraID) {
									break
								} else {
									temporalCameraID = rand.Intn(65535)
								}
							}
							responseForCamera.CameraID = temporalCameraID
							responseForCamera.Authenticated = 0
							camera.CameraID = temporalCameraID
							utils.InsertCamera(camera)
						}
						responseForCamera, _ := json.Marshal(responseForCamera)
						_, err := conn.Write(responseForCamera)
						if err != nil {
							log.Println(err.Error())
						}
					}
				}(data, conn, responseForCamera)
			default:
				return
			}
		}
	}
}
