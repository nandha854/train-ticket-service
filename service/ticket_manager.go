package service

import (
	"context"
	"log"
	"sync"

	pb "github.com/nandha854/train-ticket-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TicketManager handles ticket purchases, retrievals, and modifications.
// It interacts with SeatManager to manage seat assignments for tickets.
type TicketManager struct {
	pb.UnimplementedTicketServiceServer
	SeatManager *SeatManager
	Receipts    map[string]*pb.TicketReceipt
	mu          sync.Mutex
}

// NewTicketManager initializes a new TicketManager with a SeatManager and an empty receipts map.
func NewTicketManager() *TicketManager {
	return &TicketManager{
		SeatManager: NewSeatManager(),
		Receipts:    make(map[string]*pb.TicketReceipt),
	}
}

// PurchaseTicket processes a ticket purchase request, assigns a seat, and returns a ticket receipt.
func (t *TicketManager) PurchaseTicket(ctx context.Context, req *pb.PurchaseTicketRequest) (*pb.TicketReceipt, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	log.Printf("PurchaseTicket request received: %+v", req)

	if req.User == nil || req.User.Email == "" || req.From == "" || req.To == "" {
		log.Printf("PurchaseTicket request missing required fields: %+v", req)
		return nil, status.Error(codes.InvalidArgument, "missing required fields")
	}

	if req.From != "London" || req.To != "France" {
		log.Printf("PurchaseTicket request with invalid station: From=%s, To=%s", req.From, req.To)
		return nil, status.Error(codes.InvalidArgument, "invalid station")
	}

	seat, section, err := t.SeatManager.AssignSeat()
	if err != nil {
		log.Printf("PurchaseTicket seat assignment failed: %v", err)
		return nil, err
	}

	receipt := &pb.TicketReceipt{
		User:  req.User,
		From:  req.From,
		To:    req.To,
		Price: 20.00,
		Seat:  &pb.Seat{SeatNumber: int32(seat), Section: section},
	}

	t.Receipts[req.User.Email] = receipt

	log.Printf("PurchaseTicket successful: %+v", receipt)
	return receipt, nil
}

// GetReceipt retrieves the ticket receipt for a given email.
func (t *TicketManager) GetReceipt(ctx context.Context, req *pb.GetReceiptRequest) (*pb.TicketReceipt, error) {
	log.Printf("GetReceipt request received: %+v", req)

	if req.Email == "" {
		log.Printf("GetReceipt request missing email: %+v", req)
		return nil, status.Error(codes.InvalidArgument, "missing required fields")
	}

	receipt, ok := t.Receipts[req.Email]
	if !ok {
		log.Printf("GetReceipt not found for email: %s", req.Email)
		return nil, status.Error(codes.NotFound, "ticket receipt not found")
	}

	log.Printf("GetReceipt successful: %+v", receipt)
	return receipt, nil
}

// GetUsersBySection retrieves a list of users seated in a specific section.
func (t *TicketManager) GetUsersBySection(ctx context.Context, req *pb.GetUsersBySectionRequest) (*pb.UsersBySectionResponse, error) {
	log.Printf("GetUsersBySection request received: %+v", req)

	if req.Section == "" {
		log.Printf("GetUsersBySection request missing section: %+v", req)
		return nil, status.Error(codes.InvalidArgument, "missing required fields")
	}

	users := []*pb.UserTicket{}
	for _, receipt := range t.Receipts {
		if receipt.Seat.Section == req.Section {
			users = append(users, &pb.UserTicket{User: receipt.User, Seat: receipt.Seat})
		}
	}

	log.Printf("GetUsersBySection successful: section=%s, users=%d", req.Section, len(users))
	return &pb.UsersBySectionResponse{Users: users}, nil
}

// RemoveUser cancels a ticket and releases the assigned seat.
func (t *TicketManager) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	log.Printf("RemoveUser request received: %+v", req)

	if req.Email == "" {
		log.Printf("RemoveUser request missing email: %+v", req)
		return nil, status.Error(codes.InvalidArgument, "missing required fields")
	}

	receipt, ok := t.Receipts[req.Email]
	if !ok {
		log.Printf("RemoveUser not found for email: %s", req.Email)
		return nil, status.Error(codes.NotFound, "ticket receipt not found")
	}

	if err := t.SeatManager.ReleaseSeat(int(receipt.Seat.SeatNumber), receipt.Seat.Section); err != nil {
		log.Printf("RemoveUser seat release failed: %v", err)
		return nil, err
	}

	delete(t.Receipts, req.Email)

	log.Printf("RemoveUser successful: email=%s", req.Email)
	return &pb.RemoveUserResponse{Message: "Ticket cancelled successfully"}, nil
}

// ModifyUserSeat changes the seat assignment for a user.
func (t *TicketManager) ModifyUserSeat(ctx context.Context, req *pb.ModifyUserSeatRequest) (*pb.TicketReceipt, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	log.Printf("ModifyUserSeat request received: %+v", req)

	if req.Email == "" || req.NewSeat == nil || req.NewSeat.Section == "" || req.NewSeat.SeatNumber == 0 {
		log.Printf("ModifyUserSeat request missing required fields: %+v", req)
		return nil, status.Error(codes.InvalidArgument, "missing required fields")
	}

	receipt, ok := t.Receipts[req.Email]
	if !ok {
		log.Printf("ModifyUserSeat not found for email: %s", req.Email)
		return nil, status.Error(codes.NotFound, "ticket receipt not found")
	}

	if err := t.SeatManager.ModifySeat(int(receipt.Seat.SeatNumber), receipt.Seat.Section, int(req.NewSeat.SeatNumber), req.NewSeat.Section); err != nil {
		log.Printf("ModifyUserSeat seat modification failed: %v", err)
		return nil, err
	}

	receipt.Seat = req.NewSeat

	log.Printf("ModifyUserSeat successful: %+v", receipt)
	return receipt, nil
}
