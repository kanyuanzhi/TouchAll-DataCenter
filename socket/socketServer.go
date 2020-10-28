package socket

import (
	"dataCenter/dbDrivers"
	"dataCenter/models"
	"dataCenter/protocal"
	"dataCenter/utils"
	"dataCenter/websocket"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net"
)

// 接口服务器
type SocketServer struct {
	wsClients     *websocket.WsClients
	mongoConn     *mongo.Database
	useMongodb    bool
	mysqlConn     *sqlx.DB
	useMysql      bool
	socketClients *SocketClients
}

func NewSocketServer(wsClients *websocket.WsClients, socketClients *SocketClients) *SocketServer {
	config := utils.NewConfig()
	useMongodb := config.GetValue("mongodb.use").(bool)
	useMysql := config.GetValue("mysql.use").(bool)
	return &SocketServer{
		wsClients:     wsClients,
		mongoConn:     nil,
		useMongodb:    useMongodb,
		mysqlConn:     nil,
		useMysql:      useMysql,
		socketClients: socketClients,
	}
}

func (socketServer *SocketServer) Start() {
	config := utils.NewConfig()
	port := config.GetSocketConfig().(string)
	l, err := net.Listen("tcp", ":"+port)
	log.Printf("Start SocketServer on port %s", port)
	if err != nil {
		log.Println(err)
		return
	}

	if socketServer.useMongodb {
		socketServer.mongoConn = dbDrivers.GetMongodbConn()
		if err != nil {
			log.Println(err)
			return
		}
	}
	if socketServer.useMysql {
		socketServer.mysqlConn = dbDrivers.GetMysqlConn()
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		socketServer.socketClients.members[&conn] = true
		go socketServer.handleConn(conn)
	}
}

// 处理socket连接，完成数据解包
func (socketServer *SocketServer) handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
		delete(socketServer.socketClients.members, &conn)
	}()

	tempBuffer := make([]byte, 0)
	readerChannel := make(chan []byte, 16)
	go socketServer.reader(readerChannel)
	for {
		var buffer = make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		tempBuffer = protocal.Unpack(append(tempBuffer, buffer[:n]...), readerChannel)
	}
}

func (socketServer *SocketServer) reader(readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			// 由于此时不知道数据类型，不能转为相应的结构体，因此需要先把json转map，以获取data_type
			var m map[string]interface{}
			json.Unmarshal(data, &m)
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
						utils.InsertManyRecords(documents, socketServer.mongoConn.Collection("person_awareness"))
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
				// EquipmentBasicInformationAwareness，不做推送，只更新数据库
				if socketServer.useMysql == true {
					var equipmentBasicInformationAwareness models.EquipmentBasicInformationAwareness
					json.Unmarshal(data, &equipmentBasicInformationAwareness)
					go func(equipmentBasicInformationAwareness *models.EquipmentBasicInformationAwareness) {
						utils.InsertEquipmentBasicInformation(equipmentBasicInformationAwareness, socketServer.mysqlConn)
					}(&equipmentBasicInformationAwareness)
				}
			case 31:
				// EquipmentStatusAwareness
				var equipmentStatusAwareness models.EquipmentsStatusAwareness
				json.Unmarshal(data, &equipmentStatusAwareness)
				//log.Println(equipmentStatusAwareness)

				go func(equipmentStatusAwareness *models.EquipmentsStatusAwareness) {
					if socketServer.useMongodb == true {
						utils.InsertOneRecord(equipmentStatusAwareness, socketServer.mongoConn.Collection("equipment_awareness"))
					}
				}(&equipmentStatusAwareness)

				go func(equipmentStatusAwareness *models.EquipmentsStatusAwareness) {
					socketServer.wsClients.EquipmentStatusBroadcast <- equipmentStatusAwareness
				}(&equipmentStatusAwareness)
			}
		}
	}
}
