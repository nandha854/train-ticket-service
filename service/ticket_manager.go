package service

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/nandha854/train-ticket-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TicketManager struct {
	pb.UnimplementedTicketServiceServer
	SeatManager *SeatManager
	Receipts    map[string]*pb.TicketReceipt
	mu          sync.Mutex
}

func NewTicketManager() *TicketManager {
	return &TicketManager{
		SeatManager: NewSeatManager(),
		Receipts:    make(map[string]*pb.TicketReceipt),
	}
}

func (t *TicketManager) PurchaseTicket(ctx context.Context, req *pb.PurchaseTicketRequest) (*pb.TicketReceipt, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if req.User == nil || req.User.Email == "" || req.From == "" || req.To == "" {
		return nil, status.Error(codes.InvalidArgument, "missing required fields")
	}

	seat, section, err := t.SeatManager.AssignSeat()
	if err != nil {
		return nil, err
	}

	receipt := &pb.TicketReceipt{
		User:    req.User,
		From:   req.From,
		To:     req.To,
		Price: 20,
		Seat:  &pb.Seat{SeatNumber: int32(seat), Section: section},
	}

	t.Receipts[req.User.Email] = receipt

	return receipt, nil
}

func (t *TicketManager) GetReceipt(ctx context.Context, req *pb.GetReceiptRequest) (*pb.TicketReceipt, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "missing required fields")
	}

	receipt, ok := t.Receipts[req.Email]
	if !ok {
		return nil, status.Error(codes.NotFound, "ticket receipt not found")
	}

	return receipt, nil
}

func (t *TicketManager) GetUsersBySection(ctx context.Context, req *pb.GetUsersBySectionRequest) (*pb.UsersBySectionResponse, error) {

	if req.Section == "" {
		return nil, status.Error(codes.InvalidArgument, "missing required fields")
	}

	users := []*pb.UserTicket{}
	fmt.Println(len(t.Receipts))
	for _, receipt := range t.Receipts {
		if receipt.Seat.Section == req.Section {
			users = append(users, &pb.UserTicket{User: receipt.User, Seat: receipt.Seat})
		}
	}

	return &pb.UsersBySectionResponse{Users: users}, nil
}

func (t *TicketManager) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "missing required fields")
	}

	receipt, ok := t.Receipts[req.Email]
	if !ok {
		return nil, status.Error(codes.NotFound, "ticket receipt not found")
	}

	if err := t.SeatManager.ReleaseSeat(int(receipt.Seat.SeatNumber), receipt.Seat.Section); err != nil {
		return nil, err
	}

	delete(t.Receipts, req.Email)

	return &pb.RemoveUserResponse{Message: "Ticket cancelled successfully"}, nil
}

func (t *TicketManager) ModifyUserSeat(ctx context.Context, req *pb.ModifyUserSeatRequest) (*pb.TicketReceipt, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if req.Email == "" || req.NewSeat == nil || req.NewSeat.Section == "" || req.NewSeat.SeatNumber == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing required fields")
	}

	receipt, ok := t.Receipts[req.Email]
	if !ok {
		return nil, status.Error(codes.NotFound, "ticket receipt not found")
	}

	if err := t.SeatManager.ModifySeat(int(receipt.Seat.SeatNumber), receipt.Seat.Section, int(req.NewSeat.SeatNumber), req.NewSeat.Section); err != nil {
		return nil, err
	}

	receipt.Seat = req.NewSeat

	return receipt, nil
}

