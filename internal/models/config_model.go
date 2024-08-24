package models

type Database struct {
	Host             string `json:"host" binding:"required"`
	Port             int    `json:"port" binding:"required"`
	User             string `json:"user" binding:"required"`
	Password         string `json:"password" binding:"required"`
	DBName           string `json:"dbName" binding:"required"`
	ConnectionString string `json:"connectionString" binding:"required"`
}
type Config struct {
	ConnectionString string   `json:"connectionString" binding:"required"`
	Database         Database `json:"database" binding:"required"`
	SecretKey        string   `json:"secretKey" binding:"required"`
}
