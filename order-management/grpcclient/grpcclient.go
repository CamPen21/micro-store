package grpcclient

import (
	"context"
	"log"
	pb "micro-store/order-management/proto"
	"time"
)

type AllocationResponse struct {
	Ok      bool
	Message *string
}

func allocateResponseFromPb(r *pb.AllocationResponse) *AllocationResponse {
	return &AllocationResponse{r.Ok, r.Message}
}

func AllocateItemQuantity(client pb.InventoryServiceClient, itemId, quantity int64) (*AllocationResponse, error) {
	log.Println("Allocating item amount")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := client.AllocateItemQuantity(ctx, &pb.AllocateQuantityRequest{ItemId: int32(itemId), Count: int32(quantity)})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return allocateResponseFromPb(response), nil
}
