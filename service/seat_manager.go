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

const MaxSeats = 100

func NewSeatManager() *SeatManager {
	return &SeatManager{
		Sections: map[string]*Section{
			"A": {Name: "A", MaxSeats: 100, AvailableSeats: initializeSeats(100)},
			"B": {Name: "B", MaxSeats: 100, AvailableSeats: initializeSeats(100)},
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
	s.mu.Lock()
	defer s.mu.Unlock()
	// Modify seat

	oldSection, ok := s.Sections[seatSection]

	if !ok {
		return fmt.Errorf("section not found")
	}

	nwSection, ok := s.Sections[newSection]

	if !ok {
		return fmt.Errorf("section not found")
	}

	if nwSection.AvailableSeats[newSeat] != "Available" {
		return fmt.Errorf("seat is not available")
	}

	if oldSection.AvailableSeats[seat] == "Assigned" {
		s.ReleaseSeat(seat, seatSection)
	}

	nwSection.AvailableSeats[newSeat] = "Assigned"

	return nil
}