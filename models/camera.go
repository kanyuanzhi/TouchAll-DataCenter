package models

import "gorm.io/gorm"

type Camera struct {
	gorm.Model
	DataType      int    `json:"data_type"`
	CameraID      int    `json:"camera_id" gorm:"index"`
	CameraName    string `json:"camera_name"`
	CameraHost    string `json:"camera_host"`
	Description   string `json:"description"`
	Authenticated int    `json:"authenticated"`
}

type ResponseForCamera struct {
	DataType      int  `json:"data_type"`
	CameraID      int  `json:"camera_id"`
	Authenticated int  `json:"authenticated"`
	UseMysql      bool `json:"use_mysql"`
	UseMongodb    bool `json:"use_mongodb"`
}

func NewResponseForCamera() *ResponseForCamera {
	return &ResponseForCamera{
		DataType: 51,
	}
}
