package grpcapi

import (
	"context"
	"database/sql"
	"log"
	"micro-store/inventory-service/invdb"
	pb "micro-store/inventory-service/proto"
	"net"

	"google.golang.org/grpc"
)

type InventoryServer struct {
	pb.UnimplementedInventoryServiceServer
	db *sql.DB
}

func (s *InventoryServer) AllocateItemQuantity(ctx context.Context, r *pb.AllocateQuantityRequest) (*pb.AllocationResponse, error) {
	if err := invdb.AllocateInventoryItemAmount(s.db, int64(r.ItemId), int64(r.Count)); err != nil {
		msg := err.Error()
		return &pb.AllocationResponse{Ok: false, Message: &msg}, nil
	} else {
		return &pb.AllocationResponse{Ok: true}, nil
	}
}

func Serve(db *sql.DB, bind string) {
	listener, err := net.Listen("tcp", bind)
	if err != nil {
		log.Fatalf("gRPC server error: Failure to bind %v\n", bind)
	}
	grpcServer := grpc.NewServer()
	inventoryServer := InventoryServer{db: db}
	pb.RegisterInventoryServiceServer(grpcServer, &inventoryServer)
	log.Printf("gRPC server listening on: %v...\n", bind)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server error: %v\n", err)
	}
}
