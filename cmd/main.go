package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"social-network-api/internal/config"
	"social-network-api/internal/models"

	"github.com/joho/godotenv"
)

func main() {
	portPtr := flag.Int("port", 5000, "Default port is 5000")
	envPtr := flag.String("env", "dev", "Default env is dev")

	flag.Parse()

	env := *envPtr

	if env == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %s", err)
		}
	}

	configModel := getConfig()

	router := config.SetupApi(configModel)

	port := *portPtr
	router.Run(fmt.Sprintf(":%d", port))
}

func getConfig() models.Config {
	hostConfig := os.Getenv("DB_HOST")
	portConfig, _ := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 32)
	dbNameConfig := os.Getenv("DB_NAME")
	userConfig := os.Getenv("DB_USER")
	passwordConfig := os.Getenv("DB_PASSWORD")
	secretKeyConfig := os.Getenv("SECRET_KEY")
	connectionStringConfig := os.Getenv("CONNECTION_STRING")

	return models.Config{ConnectionString: connectionStringConfig, Database: models.Database{Host: hostConfig, Port: int(portConfig), User: userConfig, DBName: dbNameConfig, Password: passwordConfig}, SecretKey: secretKeyConfig}
}
