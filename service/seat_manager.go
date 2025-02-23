package service

import (
	"fmt"
	"sync"

)

type SeatManager struct {
	Sections map[string]*Section
	mu       sync.Mutex
	nextSection string
}

type Section struct {
	Name           string
	MaxSeats       int
	AvailableSeats map[int]string
}

const MaxSeat = 100

func NewSeatManager() *SeatManager {
	return &SeatManager{
		Sections: map[string]*Section{
			"A": {Name: "A", MaxSeats: MaxSeat, AvailableSeats: initializeSeats(MaxSeat)},
			"B": {Name: "B", MaxSeats: MaxSeat, AvailableSeats: initializeSeats(MaxSeat)},
		},
		nextSection: "A",
	}
}

func initializeSeats(count int) map[int]string {
	seats := make(map[int]string)
	for i := 1; i <= count; i++ {
		seats[i] = "Available"
	}
	return seats
}

func (s *SeatManager) AssignSeat() (int, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Assign seat

	section := s.Sections[s.nextSection]

	for seat, available := range section.AvailableSeats { 
		if available == "Available" {
			section.AvailableSeats[seat] = "Assigned"

			// Simple round-robin to assign seats
			if s.nextSection == "A" {
				s.nextSection = "B"
			} else {
				s.nextSection = "A"
			}
			return seat, section.Name, nil
		}
	}

	return 0, "", fmt.Errorf("no seats available")
}

func (s *SeatManager) ReleaseSeat(seat int, seatSection string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Release seat

	section, ok := s.Sections[seatSection]

	if !ok {
		return fmt.Errorf("Section not found")
	}

	if section.AvailableSeats[seat] == "Assigned" {
		section.AvailableSeats[seat] = "Available"
		return nil
	}

	return fmt.Errorf("seat is not assigned yet")
}

func (s *SeatManager) ModifySeat(seat int, seatSection string, newSeat int, newSection string) error {
	// Validate inputs before locking
	oldSection, ok := s.Sections[seatSection]
	if !ok {
		return fmt.Errorf("old section not found")
	}

	nwSection, ok := s.Sections[newSection]
	if !ok {
		return fmt.Errorf("new section not found")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check seat assignment
	if oldSection.AvailableSeats[seat] != "Assigned" {
		return fmt.Errorf("old seat is not assigned")
	}

	if nwSection.AvailableSeats[newSeat] != "Available" {
		return fmt.Errorf("new seat is not available")
	}

	// Swap seat assignments
	oldSection.AvailableSeats[seat] = "Available"
	nwSection.AvailableSeats[newSeat] = "Assigned"

	return nil
}

