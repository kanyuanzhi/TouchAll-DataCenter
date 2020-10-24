package models

// DataType:
// 10: PeopleAwareness
// 11: PersonAwareness
// 20: EnvironmentAwareness
// 30: EquipmentBasicInformationAwareness
// 31: EquipmentStatusAwareness

type PeopleAwareness struct {
	DataType            int                `json:"data_type" bson:"data_type"`
	Camera              int                `json:"camera" bson:"camera"`
	UpdateTime          int64              `json:"update_time" bson:"update_time"`
	PersonAwarenessData []*PersonAwareness `json:"person_awareness_data" bson:"person_awareness_data"`
}

func NewPeopleAwareness(camera int, updateTime int64, personAwarenessData []*PersonAwareness) *PeopleAwareness {
	return &PeopleAwareness{
		DataType:            1,
		Camera:              camera,
		UpdateTime:          updateTime,
		PersonAwarenessData: personAwarenessData,
	}
}

type PersonAwareness struct {
	DataType       int    `json:"data_type" bson:"data_type"`
	Name           string `json:"name" bson:"name"`
	Camera         int    `json:"camera" bson:"camera"`
	UpdateTime     int64  `json:"update_time" bson:"update_time"`
	BoarderInFrame []int  `json:"boarder_in_frame" bson:"boarder_in_frame"`
}

func NewPersonAwareness(name string, camera int, boarderInFrame []int, updateTime int64) *PersonAwareness {
	return &PersonAwareness{
		DataType:       2,
		Name:           name,
		Camera:         camera,
		BoarderInFrame: boarderInFrame,
		UpdateTime:     updateTime,
	}
}
