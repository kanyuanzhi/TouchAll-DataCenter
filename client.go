package main

//
//func main() {
//	conn, err := net.Dial("tcp", ":9090")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer conn.Close()
//	count := 0
//	for {
//		//time.Sleep(1*time.Second)
//		boarderInFrame := []int{1, 2, 3, 4}
//		personAwarenessData := make([]*models.PersonAwareness, 0)
//		for i := 0; i < 3; i++ {
//			personAwareness := models.NewPersonAwareness("kan", 1, boarderInFrame, 1111)
//			personAwarenessData = append(personAwarenessData, personAwareness)
//		}
//		peopleAwareness := models.NewPeopleAwareness(1, 111, personAwarenessData)
//
//		data, _ := json.Marshal(peopleAwareness)
//		conn.Write(protocal.Pack(data))
//		count++
//		if count == 10 {
//			break
//		}
//	}
//}
