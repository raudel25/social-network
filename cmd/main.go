package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"social-network-api/internal/config"
	"social-network-api/internal/models"
)

func main() {
	portPtr := flag.Int("port", 5000, "Default port is 8080")
	flag.Parse()

	jsonFile, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	var configModel models.Config
	err = json.Unmarshal(jsonFile, &configModel)
	if err != nil {
		log.Fatalf("Error parsing JSON: %s", err)
	}

	router := config.SetupApi(configModel)

	port := *portPtr
	router.Run(fmt.Sprintf(":%d", port))
}
