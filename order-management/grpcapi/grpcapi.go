package grpcapi

import (
	"context"
	"database/sql"
	"log"
	"micro-store/order-management/orderdb"
	pb "micro-store/order-management/proto"
	"net"

	"google.golang.org/grpc"
)

func orderedItemFromPbOrderItem(item *pb.OrderItem) *orderdb.OrderedItem {
	return &orderdb.OrderedItem{ItemId: int64(item.ItemId), Quantity: int64(item.Quantity)}
}

type OrderManagementService struct {
	pb.UnimplementedOrderManagementServiceServer
	db *sql.DB
}

func (s *OrderManagementService) PlaceOrder(ctx context.Context, r *pb.PlaceOrderRequest) (*pb.OrderPlacementResponse, error) {
	tx, err := s.db.Begin()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer tx.Rollback()
	orderedItems := make([]*orderdb.OrderedItem, 0)
	for _, pbItem := range r.OrderItems {
		orderedItems = append(orderedItems, orderedItemFromPbOrderItem(pbItem))
	}
	if orderId, err := orderdb.CreateOrder(tx, orderedItems); err != nil {
		return &pb.OrderPlacementResponse{Success: false, OrderId: -1}, err
	} else {
		return &pb.OrderPlacementResponse{Success: true, OrderId: int32(orderId)}, nil
	}
}

func (s *OrderManagementService) CancelOrder(ctx context.Context, r *pb.CancelOrderRequest) (*pb.OrderPlacementResponse, error) {
	tx, err := s.db.Begin()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer tx.Rollback()
	if err := orderdb.CancelOrder(tx, int64(r.OrderId)); err != nil {
		return &pb.OrderPlacementResponse{Success: false, OrderId: r.OrderId}, err
	} else {
		return &pb.OrderPlacementResponse{Success: true, OrderId: r.OrderId}, nil
	}
}

func Serve(db *sql.DB, bind string) {
	listener, err := net.Listen("tcp", bind)
	if err != nil {
		log.Fatalf("gRPC server error: Failure to bind %v\n", bind)
	}
	grpcServer := grpc.NewServer()
	inventoryServer := OrderManagementService{db: db}
	pb.RegisterOrderManagementServiceServer(grpcServer, &inventoryServer)
	log.Printf("gRPC server listening on: %v...\n", bind)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server error: %v\n", err)
	}
}
