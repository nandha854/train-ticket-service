package main

import (
	"log"
	"net"

	pb "github.com/nandha854/train-ticket-service/proto"
	"github.com/nandha854/train-ticket-service/service"
	"google.golang.org/grpc"
)

func main(){
	// Create a new gRPC server 
	server := grpc.NewServer() 
	sectionConfigs := []service.SectionConfigs{
		{SectionName: "A", MaxSeats: 50},
		{SectionName: "B", MaxSeats: 50},
	}

	// Initialize a new SeatManager
	seatManager := service.NewSeatManager(sectionConfigs)

	// Initialize a stationConnection
	connectionStations := map[string]float64{
		"London-France": 20.00,
	}

	// Register the service with the server 
	pb.RegisterTicketServiceServer(server, service.NewTicketManager(seatManager, connectionStations)) 

	// Start listening on a port (e.g., 50051) 
	listen, err := net.Listen("tcp", ":50051") 

	if err != nil { 
		log.Fatalf("failed to listen: %v", err)
	} 

	// Start the gRPC server 
	log.Println("Server listening on port 50051...") 
	
	if err := server.Serve(listen); err != nil {
		 log.Fatalf("failed to serve: %v", err)
	}
}