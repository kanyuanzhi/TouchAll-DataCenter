package models

type PeopleAwareness struct {
	DataType            int                `json:"data_type" bson:"data_type"`
	CameraID            int                `json:"camera_id" bson:"camera_id"`
	UpdatedAt           int64              `json:"updated_at" bson:"updated_at"`
	PersonAwarenessData []*PersonAwareness `json:"person_awareness_data" bson:"person_awareness_data"`
}

type PersonAwareness struct {
	DataType       int    `json:"data_type" bson:"data_type"`
	Name           string `json:"name" bson:"name"`
	CameraID       int    `json:"camera_id" bson:"camera_id"`
	UpdatedAt      int64  `json:"updated_at" bson:"updated_at"`
	BoarderInFrame []int  `json:"boarder_in_frame" bson:"boarder_in_frame"`
}
