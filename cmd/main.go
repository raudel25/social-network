package main

import (
	"flag"
	"fmt"
	"social-network-api/internal/config"
)

func main() {
	config.ConnectDatabase()

	portPtr := flag.Int("port", 8080, "Default port is 8080")
	flag.Parse()

	router := config.SetupRoutes() // Asegúrate de pasar db a SetupRoutes si es necesario

	port := *portPtr
	router.Run(fmt.Sprintf(":%d", port))
}
