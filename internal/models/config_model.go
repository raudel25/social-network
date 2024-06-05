package models

type Database struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
	DBName   string `json:"dbName" binding:"required"`
}
type Config struct {
	Database  `json:"Database" binding:"required"`
	SecretKey string `json:"SecretKey" binding:"required"`
}
