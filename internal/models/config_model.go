package models

type Config struct {
	Database struct {
		Host     string `json:"Host" binding:"required"`
		Port     int    `json:"Port" binding:"required"`
		User     string `json:"User" binding:"required"`
		Password string `json:"Password" binding:"required"`
		DBName   string `json:"DBName" binding:"required"`
	} `json:"Database" binding:"required"`
	SecretKey string `json:"SecretKey" binding:"required"`
}
