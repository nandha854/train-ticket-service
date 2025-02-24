package service

import (
	"fmt"
	"sync"
)

// SeatManager handles the assignment, release, and modification of seats.
// It manages seats across different sections in a round-robin manner.
type SeatManager struct {
	Sections    map[string]*Section
	mu          sync.Mutex
	nextSections []string
	nextSection int
}

type Section struct {
	Name           string
	MaxSeats       int
	AvailableSeats map[int]string
}

type SectionConfigs struct {
	SectionName string
	MaxSeats    int
}


// NewSeatManager initializes a new SeatManager with predefined sections and seats.
func NewSeatManager(sectionConfigs []SectionConfigs) *SeatManager {

	sections := make(map[string]*Section)
	nextSections := []string{}
	for _, sectionConfig := range sectionConfigs {
		sections[sectionConfig.SectionName] = &Section{
			Name: sectionConfig.SectionName,
			MaxSeats: sectionConfig.MaxSeats,
			AvailableSeats: initializeSeats(sectionConfig.MaxSeats),
		}
		nextSections = append(nextSections, sectionConfig.SectionName)
	}

	return &SeatManager{
		Sections: sections,
		nextSections: nextSections,
		nextSection: 0,
	}
}

// initializeSeats creates a map of seats marked as "Available".
func initializeSeats(count int) map[int]string {
	seats := make(map[int]string)
	for i := 1; i <= count; i++ {
		seats[i] = "Available"
	}
	return seats
}

// AssignSeat assigns the next available seat in a round-robin manner.
func (s *SeatManager) AssignSeat() (int, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	section := s.Sections[s.nextSections[s.nextSection]]

	for seat, available := range section.AvailableSeats {
		if available == "Available" {
			section.AvailableSeats[seat] = "Assigned"

			// Simple round-robin to assign seats
			s.nextSection = (s.nextSection + 1) % len(s.nextSections)
			return seat, section.Name, nil
		}
	}

	return 0, "", fmt.Errorf("no seats available")
}

// ReleaseSeat releases an assigned seat, making it available again.
func (s *SeatManager) ReleaseSeat(seat int, seatSection string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

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

// ModifySeat changes the seat assignment from one seat to another.
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
