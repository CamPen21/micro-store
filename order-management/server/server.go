package main

import (
	"database/sql"
	"log"
	"micro-store/order-management/grpcapi"
	"micro-store/order-management/orderdb"
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
		args.DbPath = "./orders.db"
	}
	if args.BindGrcp == "" {
		args.BindGrcp = ":8081"
	}
	db, err := sql.Open("sqlite3", args.DbPath)
	if err != nil {
		log.Fatal("Couldn't connect to the database")
	}
	defer db.Close()

	orderdb.Initialize(db)

	var sWg sync.WaitGroup

	sWg.Add(1)
	go func() {
		log.Println("Starting gRPC server...")
		grpcapi.Serve(db, args.BindGrcp)
		defer sWg.Done()
	}()

	sWg.Wait()
}
