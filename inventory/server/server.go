package main

import (
	"database/sql"
	"log"
	"micro-store/inventory-service/grpcapi"
	"micro-store/inventory-service/invdb"
	"sync"

	"github.com/alexflint/go-arg"
)

var args struct {
	DbPath   string `arg:"env:INVENTORY_DB_PATH"`
	BindGrcp string `arg:"env:INVENTORY_GRCP_BIND_ADDR"`
}

func main() {
	arg.MustParse(&args)
	if args.DbPath == "" {
		args.DbPath = "./inventory.db"
	}
	if args.BindGrcp == "" {
		args.BindGrcp = ":3000"
	}
	db, err := sql.Open("sqlite3", args.DbPath)
	if err != nil {
		log.Fatal("Couldn't connect to the database")
	}
	defer db.Close()

	invdb.Initialize(db)

	var sWg sync.WaitGroup

	sWg.Add(1)
	go func() {
		log.Println("Starting gRPC server...")
		grpcapi.Serve(db, args.BindGrcp)
		defer sWg.Done()
	}()

	sWg.Wait()
}
