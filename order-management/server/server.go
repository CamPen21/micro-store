package main

import (
	"database/sql"
	"log"
	"micro-store/order-management/grpcapi"
	"micro-store/order-management/orderdb"
	pb "micro-store/order-management/proto"
	"sync"

	"github.com/alexflint/go-arg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var args struct {
	DbPath            string `arg:"env:ORDERS_DB_PATH"`
	BindGrcp          string `arg:"env:ORDERS_GRCP_BIND_ADDR"`
	InventoryBindGrcp string `arg:"env:INVENTORY_GRCP_BIND_ADDR"`
}

func main() {
	arg.MustParse(&args)
	if args.DbPath == "" {
		args.DbPath = "./orders.db"
	}
	if args.BindGrcp == "" {
		args.BindGrcp = ":8081"
	}
	if args.InventoryBindGrcp == "" {
		args.InventoryBindGrcp = ":3000"
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
		conn, err := grpc.Dial(args.InventoryBindGrcp, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			clientErr := err.Error()
			log.Fatalf("Couldn't connect to inventory: %v", clientErr)
		}
		defer conn.Close()
		client := pb.NewInventoryServiceClient(conn)
		log.Println("Starting gRPC server...")
		grpcapi.Serve(db, args.BindGrcp, client)
		defer sWg.Done()
	}()

	sWg.Wait()
}
