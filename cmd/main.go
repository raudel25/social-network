package main

import (
	"flag"
	"fmt"
	"social-network-api/internal/routes"
)

func main() {
	portPtr := flag.Int("port", 8080, "Default port is 8080")
	flag.Parse()

	router := routes.SetupRoutes()

	port := *portPtr
	router.Run(fmt.Sprintf(":%d", port))
}
