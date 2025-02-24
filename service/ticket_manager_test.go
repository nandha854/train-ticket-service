package service

import (
	"context"
	"testing"

	pb "github.com/nandha854/train-ticket-service/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createTestTicketManager() *TicketManager {
    sectionConfigs := []SectionConfigs{
        {SectionName: "A", MaxSeats: 50},
        {SectionName: "B", MaxSeats: 50},
    }
	stationConnection := map[string]float64{
		"London-France": 20.00,
	}
    seatManager := NewSeatManager(sectionConfigs)
    return NewTicketManager(seatManager, stationConnection)
}

func TestNewTicketManager(t *testing.T) {
    tm := createTestTicketManager()
    assert.NotNil(t, tm.SeatManager, "SeatManager should be initialized")
}

func TestPurchaseTicket(t *testing.T) {
    tm := createTestTicketManager()

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
        {
            name: "Invalid Station Info",
            request: &pb.PurchaseTicketRequest{
                From: "Chennai",
                To:   "Coimbatore",
            },
            expectError: true,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            resp, err := tm.PurchaseTicket(context.Background(), tc.request)

            if tc.expectError {
                assert.Error(t, err)
                st, ok := status.FromError(err)
                assert.True(t, ok)
                assert.Equal(t, codes.InvalidArgument, st.Code())
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp.Seat)
            }
        })
    }
}

func TestGetReceipt(t *testing.T) {
    tm := createTestTicketManager()

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
        expectCode  codes.Code
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
            expectCode:  codes.NotFound,
        },
        {
            name:        "Missing Email",
            request:     &pb.GetReceiptRequest{},
            expectError: true,
            expectCode:  codes.InvalidArgument,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            resp, err := tm.GetReceipt(context.Background(), tc.request)

            if tc.expectError {
                assert.Error(t, err)
                st, ok := status.FromError(err)
                assert.True(t, ok)
                assert.Equal(t, tc.expectCode, st.Code())
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
            }
        })
    }
}

func TestGetUsersBySection(t *testing.T) {
    tm := createTestTicketManager()

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
        expectCode  codes.Code
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
            expectCode:  codes.InvalidArgument,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            resp, err := tm.GetUsersBySection(context.Background(), tc.request)

            if tc.expectError {
                assert.Error(t, err)
                st, ok := status.FromError(err)
                assert.True(t, ok)
                assert.Equal(t, tc.expectCode, st.Code())
            } else {
                assert.NoError(t, err)
                assert.Len(t, resp.Users, tc.expectCount)
            }
        })
    }
}

func TestModifyUserSeat(t *testing.T) {
    tm := createTestTicketManager()

    userEmail := "modify@example.com"
    seatNumber, section := 10, "A"

    // Assign a seat using SeatManager
    tm.SeatManager.Sections[section].AvailableSeats[seatNumber] = "Assigned"

    tm.Receipts[userEmail] = &pb.TicketReceipt{
        User: &pb.User{FirstName: "Kumar", LastName: "Test", Email: userEmail},
        Seat: &pb.Seat{SeatNumber: int32(seatNumber), Section: section},
    }

    tests := []struct {
        name        string
        request     *pb.ModifyUserSeatRequest
        expectError bool
        expectCode  codes.Code
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
            expectCode:  codes.NotFound,
        },
        {
            name: "Missing Seat Information",
            request: &pb.ModifyUserSeatRequest{
                Email:   userEmail,
                NewSeat: nil,
            },
            expectError: true,
            expectCode:  codes.InvalidArgument,
        },
        {
            name: "Assigning Already Taken Seat",
            request: &pb.ModifyUserSeatRequest{
                Email:   userEmail,
                NewSeat: &pb.Seat{SeatNumber: 20, Section: "A"},
            },
            expectError: true,
			expectCode:  codes.InvalidArgument,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            resp, err := tm.ModifyUserSeat(context.Background(), tc.request)

            if tc.expectError {
                assert.Error(t, err)
                st, ok := status.FromError(err)
                assert.True(t, ok)
                assert.Equal(t, tc.expectCode, st.Code())
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.Equal(t, tc.request.NewSeat.SeatNumber, resp.Seat.SeatNumber)
            }
        })
    }
}

func TestRemoveUser(t *testing.T) {
    tm := createTestTicketManager()

    userEmail := "remove@example.com"
    // Assign a seat using SeatManager
    seatNumber := 1
    section := "A"
    tm.SeatManager.Sections[section].AvailableSeats[seatNumber] = "Assigned"

    tm.Receipts[userEmail] = &pb.TicketReceipt{
        User: &pb.User{FirstName: "Kumar", LastName: "Test", Email: userEmail},
        Seat: &pb.Seat{SeatNumber: int32(seatNumber), Section: section},
    }

    resp, err := tm.RemoveUser(context.Background(), &pb.RemoveUserRequest{Email: userEmail})
    assert.NoError(t, err, "User removal should be successful")
    assert.Equal(t, "Ticket cancelled successfully", resp.Message)
    assert.NotContains(t, tm.Receipts, userEmail, "User should be removed from receipts")
    assert.Equal(t, "Available", tm.SeatManager.Sections[section].AvailableSeats[seatNumber], "Seat should be available after removal")
}