package main

import (
	"flag"
	"fmt"
	"social-network-api/internal/config"
	"social-network-api/internal/db"
)

func main() {
	portPtr := flag.Int("port", 8080, "Default port is 8080")
	flag.Parse()

	db := db.ConnectDatabase()
	router := config.SetupApi(db)

	port := *portPtr
	router.Run(fmt.Sprintf(":%d", port))
}
