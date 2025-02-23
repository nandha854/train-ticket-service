package service

import (
	"context"
	"sync"

	pb "github.com/nandha854/train-ticket-service/proto"
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




