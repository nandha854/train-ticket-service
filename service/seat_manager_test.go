package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test NewSeatManager initializes correctly
func TestNewSeatManager(t *testing.T) {
	seatManager := NewSeatManager()

	assert.NotNil(t, seatManager, "SeatManager should be initialized")
	assert.Equal(t, 2, len(seatManager.Sections), "Should have 2 sections")
	assert.Contains(t, seatManager.Sections, "A", "Section A should exist")
	assert.Contains(t, seatManager.Sections, "B", "Section B should exist")
	assert.Equal(t, MaxSeat, len(seatManager.Sections["A"].AvailableSeats), "Section A should have 100 seats")
	assert.Equal(t, MaxSeat, len(seatManager.Sections["B"].AvailableSeats), "Section B should have 100 seats")
}

// Table-driven test for AssignSeat
func TestAssignSeat(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*SeatManager) // Setup function to pre-fill data
		expectErr   bool
		expectedSec string
	}{
		{
			name:        "Assign first available seat",
			setup:       func(sm *SeatManager) {}, // No setup, should work normally
			expectErr:   false,
			expectedSec: "A",
		},
		{
			name: "All seats occupied",
			setup: func(sm *SeatManager) {
				for i := 0; i < MaxSeat*2; i++ {
					_, _, _ = sm.AssignSeat() // Assign all seats
				}
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seatManager := NewSeatManager()
			tt.setup(seatManager)

			seat, section, err := seatManager.AssignSeat()

			if tt.expectErr {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Did not expect an error but got one")
				assert.Greater(t, seat, 0, "Seat number should be greater than 0")
				assert.NotEmpty(t, section, "Section should not be empty")
				assert.Equal(t, tt.expectedSec, section, fmt.Sprintf("Expected section %s but got %s", tt.expectedSec, section))
			}
		})
	}
}

// Test Seat Release
func TestReleaseSeat(t *testing.T) {
	seatManager := NewSeatManager()
	seat, section, _ := seatManager.AssignSeat()

	t.Run("Successfully release a seat", func(t *testing.T) {
		err := seatManager.ReleaseSeat(seat, section)
		assert.NoError(t, err, "Releasing an assigned seat should not return an error")
		assert.Equal(t, "Available", seatManager.Sections[section].AvailableSeats[seat], "Seat should be available after release")
	})

	t.Run("Releasing unassigned seat should fail", func(t *testing.T) {
		err := seatManager.ReleaseSeat(seat, section)
		assert.Error(t, err, "Expected an error when releasing an already available seat")
	})
}

// Test Modify Seat
func TestModifySeat(t *testing.T) {
	seatManager := NewSeatManager()
	seat, section, _ := seatManager.AssignSeat()
	newSeat := 2 // We assume seat 2 is available

	t.Run("Modify seat successfully", func(t *testing.T) {
		err := seatManager.ModifySeat(seat, section, newSeat, section)
		assert.NoError(t, err, "Modifying seat should not return an error")
		assert.Equal(t, "Available", seatManager.Sections[section].AvailableSeats[seat], "Old seat should be available after modification")
		assert.Equal(t, "Assigned", seatManager.Sections[section].AvailableSeats[newSeat], "New seat should be assigned")
	})

	t.Run("Modify to an occupied seat should fail", func(t *testing.T) {
		err := seatManager.ModifySeat(seat, section, newSeat, section)
		assert.Error(t, err, "Expected error when modifying to an already assigned seat")
	})

	t.Run("Modify seat in non-existent section should fail", func(t *testing.T) {
		err := seatManager.ModifySeat(seat, section, 10, "C") // Section C does not exist
		assert.Error(t, err, "Expected error when modifying to a non-existent section")
	})
}
