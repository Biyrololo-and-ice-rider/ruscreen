package domain

import "time"

// ParticipantRole defines the role of a participant in a room.
type ParticipantRole string

const (
	RoleOrganizer ParticipantRole = "ORGANIZER"
	RoleGuest     ParticipantRole = "GUEST"
)

// Participant represents a person in a conference room.
type Participant struct {
	ID            string
	Name          string
	Role          ParticipantRole
	CameraEnabled bool
	MicEnabled    bool
	JoinedAt      time.Time
}

// NewParticipant creates a Participant, validating name invariants.
func NewParticipant(id, name string, role ParticipantRole, joinedAt time.Time) (Participant, error) {
	p := Participant{
		ID:            id,
		Name:          name,
		Role:          role,
		CameraEnabled: true,
		MicEnabled:    true,
		JoinedAt:      joinedAt,
	}
	if err := p.validateName(); err != nil {
		return Participant{}, err
	}
	return p, nil
}

// IsOrganizer returns true if the participant holds the ORGANIZER role.
func (p *Participant) IsOrganizer() bool {
	return p.Role == RoleOrganizer
}

func (p *Participant) validateName() error {
	if p.Name == "" {
		return ErrEmptyName
	}
	if len([]rune(p.Name)) > 50 {
		return ErrNameTooLong
	}
	return nil
}
