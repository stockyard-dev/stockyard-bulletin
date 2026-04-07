package main

import (
	"flag"
	"fmt"
	"github.com/stockyard-dev/stockyard-bulletin/internal/server"
	"github.com/stockyard-dev/stockyard-bulletin/internal/store"
	"log"
	"os"
)

func main() {
	portFlag := flag.String("port", "", "")
	dataFlag := flag.String("data", "", "")
	flag.Parse()
	port := os.Getenv("PORT")
	if port == "" {
		port = "9260"
	}
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = ", "
	}
	if *portFlag != "" {
		port = *portFlag
	}
	if *dataFlag != "" {
		dataDir = *dataFlag
	}
	db, err := store.Open(dataDir)
	if err != nil {
		log.Fatalf("bulletin: %v", err)
	}
	defer db.Close()
	srv := server.New(db, server.DefaultLimits(), dataDir)
	fmt.Printf("\n  Stockyard Bulletin\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Questions? hello@stockyard.dev — I read every message\n\n", port, port)
	log.Printf("bulletin: listening on :%s", port)
	log.Fatal(srv.ListenAndServe(":" + port))
}
