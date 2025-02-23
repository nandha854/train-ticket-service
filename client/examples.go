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

	user1 := &proto.User{Email: "test23@example.com", FirstName: "Nandha", LastName: "Kumar"}
    purchaseResp1, err := client.PurchaseTicket(context.Background(), &proto.PurchaseTicketRequest{
        From: "London",
        To:   "France",
        User: user1,
    })


    if err != nil {
        log.Fatalf("PurchaseTicket failed: %v", err)
    }
    log.Printf("Purchased Ticket: %v", purchaseResp1)

	user2 := &proto.User{Email: "test3@example.com", FirstName: "Nandha", LastName: "Kumar"}
    purchaseResp2, err := client.PurchaseTicket(context.Background(), &proto.PurchaseTicketRequest{
        From: "London",
        To:   "France",
        User: user2,
    })


    if err != nil {			// Simple round-robin to assign seats

        log.Fatalf("PurchaseTicket failed: %v", err)
    }
    log.Printf("Purchased Ticket: %v", purchaseResp2)

	// Get Ticket
	getResp, err := client.GetReceipt(context.Background(), &proto.GetReceiptRequest{ Email: user.Email })
	if err != nil {
		log.Fatalf("GetTicket failed: %v", err)
	}
	log.Printf("Get Ticket: %v", getResp)

	// Get User by Section
	userResp, err := client.GetUsersBySection(context.Background(), &proto.GetUsersBySectionRequest{ Section: "A" })
	if err != nil {
		log.Fatalf("GetUsersBySection failed: %v", err)
	}
	log.Printf("Get Users by Section: %v", userResp)

	// Modify User Seat	
	modifyResp, err := client.ModifyUserSeat(context.Background(), &proto.ModifyUserSeatRequest{
		Email: user.Email,
		NewSeat: &proto.Seat{SeatNumber: getResp.Seat.GetSeatNumber() + 1, Section: getResp.Seat.GetSection()},
	})
	if err != nil {
		log.Fatalf("ModifyUserSeat failed: %v", err)
	}
	log.Printf("Modify User Seat: %v", modifyResp)

	// Remove User
	removeResp, err := client.RemoveUser(context.Background(), &proto.RemoveUserRequest{ Email: user.Email })
	if err != nil {
		log.Fatalf("RemoveUser failed: %v", err)
	}
	log.Printf("Remove User: %v", removeResp)
}