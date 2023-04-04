package model

import (
	"gorm.io/gorm"
)

func InsertData(DB *gorm.DB, data interface{}) error {
	// Use GORM to insert the data into the database
	err := DB.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func SearchData(db *gorm.DB, data interface{}) (interface{}, error) {
	// Get the User data from the database
	if err := db.Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}
