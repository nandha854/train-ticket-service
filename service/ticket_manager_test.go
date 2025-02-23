package service

import (
	"context"
	"testing"

	pb "github.com/nandha854/train-ticket-service/proto"
	"github.com/stretchr/testify/assert"
)

func TestNewTicketManager(t *testing.T) {
	tm := NewTicketManager()
	assert.NotNil(t, tm.SeatManager, "SeatManager should be initialized")
}

func TestPurchaseTicket(t *testing.T) {
	tm := NewTicketManager()

	tests := []struct {
		name        string
		request     *pb.PurchaseTicketRequest
		expectError bool
	}{
		{
			name: "Valid Ticket Purchase",
			request: &pb.PurchaseTicketRequest{
				User: &pb.User{FirstName: "Nandha", LastName: "Kumar", Email: "test@example.com"},
				From: "London",
				To:   "France",
			},
			expectError: false,
		},
		{
			name: "Missing User Info",
			request: &pb.PurchaseTicketRequest{
				From: "London",
				To:   "France",
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := tm.PurchaseTicket(context.Background(), tc.request)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp.Seat)
			}
		})
	}
}

func TestGetReceipt(t *testing.T) {
	tm := NewTicketManager()

	userEmail := "receipt@example.com"
	tm.Receipts[userEmail] = &pb.TicketReceipt{
		User: &pb.User{FirstName: "Test", LastName: "User", Email: userEmail},
		Seat: &pb.Seat{SeatNumber: 5, Section: "A"},
		From: "Chennai",
		To:   "Bangalore",
	}

	tests := []struct {
		name        string
		request     *pb.GetReceiptRequest
		expectError bool
	}{
		{
			name:        "Valid Receipt",
			request:     &pb.GetReceiptRequest{Email: userEmail},
			expectError: false,
		},
		{
			name:        "Invalid Email",
			request:     &pb.GetReceiptRequest{Email: "invalid@example.com"},
			expectError: true,
		},
		{
			name:        "Missing Email",
			request:     &pb.GetReceiptRequest{},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := tm.GetReceipt(context.Background(), tc.request)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestGetUsersBySection(t *testing.T) {
	tm := NewTicketManager()

	// Adding users to Section "A"
	tm.Receipts["user1@example.com"] = &pb.TicketReceipt{
		User: &pb.User{FirstName: "Nandha", LastName: "Kumar", Email: "user1@example.com"},
		Seat: &pb.Seat{SeatNumber: 1, Section: "A"},
	}
	tm.Receipts["user2@example.com"] = &pb.TicketReceipt{
		User: &pb.User{FirstName: "Test", LastName: "User", Email: "user2@example.com"},
		Seat: &pb.Seat{SeatNumber: 2, Section: "A"},
	}

	tests := []struct {
		name        string
		request     *pb.GetUsersBySectionRequest
		expectCount int
		expectError bool
	}{
		{
			name:        "Valid Section with Users",
			request:     &pb.GetUsersBySectionRequest{Section: "A"},
			expectCount: 2,
			expectError: false,
		},
		{
			name:        "Valid Section with No Users",
			request:     &pb.GetUsersBySectionRequest{Section: "B"},
			expectCount: 0,
			expectError: false,
		},
		{
			name:        "Invalid Section",
			request:     &pb.GetUsersBySectionRequest{Section: ""},
			expectCount: 0,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := tm.GetUsersBySection(context.Background(), tc.request)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, resp.Users, tc.expectCount)
			}
		})
	}
}

func TestModifyUserSeat(t *testing.T) {
	tm := NewTicketManager()

	userEmail := "modify@example.com"
	seatNumber, section := 10, "A"
	tm.SeatManager.Sections[section].AvailableSeats[seatNumber] = "Assigned"

	tm.Receipts[userEmail] = &pb.TicketReceipt{
		User: &pb.User{FirstName: "Kumar", LastName: "Test", Email: userEmail},
		Seat: &pb.Seat{SeatNumber: int32(seatNumber), Section: section},	
	}

	tests := []struct {
		name        string
		request     *pb.ModifyUserSeatRequest
		expectError bool
	}{
		{
			name: "Valid Seat Modification",
			request: &pb.ModifyUserSeatRequest{
				Email:   userEmail,
				NewSeat: &pb.Seat{SeatNumber: 20, Section: "A"},
			},
			expectError: false,
		},
		{
			name: "Invalid Email",
			request: &pb.ModifyUserSeatRequest{
				Email:   "invalid@example.com",
				NewSeat: &pb.Seat{SeatNumber: 30, Section: "A"},
			},
			expectError: true,
		},
		{
			name: "Missing Seat Information",
			request: &pb.ModifyUserSeatRequest{
				Email:   userEmail,
				NewSeat: nil,
			},
			expectError: true,
		},
		{
			name: "Assigning Already Taken Seat",
			request: &pb.ModifyUserSeatRequest{
				Email:   userEmail,
				NewSeat: &pb.Seat{SeatNumber: 20, Section: "A"},
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			resp, err := tm.ModifyUserSeat(context.Background(), tc.request)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tc.request.NewSeat.SeatNumber, resp.Seat.SeatNumber)
			}
		})
	}
}

func TestRemoveUser(t *testing.T) {
	tm := NewTicketManager()

	userEmail := "remove@example.com"
	seatNumber, section, err := tm.SeatManager.AssignSeat()

	tm.Receipts[userEmail] = &pb.TicketReceipt{
		User: &pb.User{FirstName: "Kumar", LastName: "Test", Email: userEmail},
		Seat: &pb.Seat{SeatNumber: int32(seatNumber), Section: section},	
	}
	assert.NoError(t, err, "Seat should be assigned successfully before removing user")

	resp, err := tm.RemoveUser(context.Background(), &pb.RemoveUserRequest{Email: userEmail})
	assert.NoError(t, err, "User removal should be successful")
	assert.Equal(t, "Ticket cancelled successfully", resp.Message)
	assert.NotContains(t, tm.Receipts, userEmail, "User should be removed from receipts")
}

