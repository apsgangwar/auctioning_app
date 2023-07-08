package main

import (
	"log"
	"net/http"
	"os"

	"auction/db"
	"auction/router"
	"auction/svc"
)

func main() {
	// Connect to the database
	svc.DB = db.GetConnection()
	defer svc.DB.Close()

	// Initialize the routes
	router.Initialize()

	// Start the server
	port := os.Getenv("AUCTION_SVC_PORT")
	log.Println("Server started running on http://127.0.0.1:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
