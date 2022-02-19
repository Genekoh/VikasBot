package database

type ServerConfig struct {
	GuildId string `gorm:"primaryKey"`
	Prefix  string
}
