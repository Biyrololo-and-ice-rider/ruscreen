package port

import "github.com/yourorg/meet/internal/domain"

// TokenGenerator is the outgoing port for JWT creation and verification.
// Implementation lives in internal/infra/token/jwt_generator.go.
type TokenGenerator interface {
	// NewToken creates a signed JWT embedding roomID, participantID, and role.
	NewToken(roomID, participantID string, role domain.ParticipantRole) (string, error)

	// Verify parses and validates a JWT, returning its embedded claims.
	Verify(token string) (roomID, participantID string, role domain.ParticipantRole, err error)
}
