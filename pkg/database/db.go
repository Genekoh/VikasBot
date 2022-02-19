package database

import (
	"errors"
	"fmt"
	"github.com/Genekoh/VikasBot/pkg/bot"
	"github.com/Genekoh/VikasBot/pkg/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDatabase() error {
	environ, err := env.GetEnviron("")
	if err != nil {
		return err
	}

	host := environ["DB_HOST"]
	port := environ["DB_PORT"]
	user := environ["DB_USER"]
	name := environ["DB_NAME"]
	password := environ["DB_PASSWORD"]

	dsn := fmt.Sprintf(
		"host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, name, password, port,
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	fmt.Println("Connected to Database")

	err = db.AutoMigrate(&ServerConfig{})
	if err != nil {
		return err
	}

	return nil
}

func GetServerConfig(guildId string) (ServerConfig, error) {
	var conf ServerConfig
	x := db.Model(&ServerConfig{}).First(&conf, guildId)

	err := x.Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return conf, err
		}

		serverConfig := ServerConfig{
			GuildId: guildId,
			Prefix:  bot.DefaultPrefix,
		}
		db.Create(&serverConfig)

		err = db.First(&conf, guildId).Error
		if err != nil {
			return conf, err
		}
	}

	return conf, nil
}

func UpdateServerConfig(conf ServerConfig) {
	db.Save(&conf)
}
