package main

import (
	"context"
	"flag"
	"log"

	"github.com/nandha854/train-ticket-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
    addr = flag.String("addr", "localhost:50051", "server address")
)

func main() {
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    client := proto.NewTicketServiceClient(conn)

    // Purchase Ticket
    user := &proto.User{Email: "test@example.com", FirstName: "Nandha", LastName: "Kumar"}
    purchaseResp, err := client.PurchaseTicket(context.Background(), &proto.PurchaseTicketRequest{
        From: "London",
        To:   "France",
        User: user,
    })
    if err != nil {
        log.Fatalf("PurchaseTicket failed: %v", err)
    }
    log.Printf("Purchased Ticket: %v", purchaseResp)
}
