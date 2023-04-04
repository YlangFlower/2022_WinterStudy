package model

import (
	"fmt"
	"log"
	"study/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	var cfg config.Config
	config.GetSettings(&cfg)
	cfg.Database.Type = "mysql"

	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Protocol, cfg.Database.Host, cfg.Database.Port, cfg.Database.Db)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return nil, err
	}
	// Migrate all tables
	if err := AutoMigrateModels(db); err != nil {
		log.Println("AutoMigrateModels error :", err)
	}
	return db, nil
}

func AutoMigrateModels(db *gorm.DB) error {
	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).AutoMigrate(
		&City{}, &District{}, &Dong{}, &User{}, &ProductCategory{},
		&ProductStatus{}, &Product{}, &ChRoom{}, &Chat{}, &Board{},
		&Comment{}, &CoC{},
	); err != nil {
		return err
	}
	// create the 'boards' table
	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).AutoMigrate(&Board{}); err != nil {
		return err
	}
	// create the 'posts' table with foreign key constraint
	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).AutoMigrate(&Post{}); err != nil {
		return err
	}
	return nil
}

// called from request file
func SetupAssociations(db *gorm.DB, user *User) error {
	for _, dongID := range user.DongIDs {
		var dong Dong
		if err := db.Where("id = ?", dongID).First(&dong).Error; err != nil {
			return err
		}
		if err := db.Model(user).Association("Dongs").Append(&dong); err != nil {
			return err
		}
	}
	return nil
}
