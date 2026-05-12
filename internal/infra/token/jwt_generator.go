package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourorg/meet/internal/domain"
)

type claims struct {
	jwt.RegisteredClaims
	RoomID        string `json:"room_id"`
	ParticipantID string `json:"participant_id"`
	Role          string `json:"role"`
}

// JWTGenerator implements port.TokenGenerator using HMAC-SHA256 JWT.
type JWTGenerator struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTGenerator(secret string) *JWTGenerator {
	return &JWTGenerator{secret: []byte(secret), ttl: 24 * time.Hour}
}

func (g *JWTGenerator) NewToken(roomID, participantID string, role domain.ParticipantRole) (string, error) {
	c := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		RoomID:        roomID,
		ParticipantID: participantID,
		Role:          string(role),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signed, err := token.SignedString(g.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

func (g *JWTGenerator) Verify(tokenStr string) (roomID, participantID string, role domain.ParticipantRole, err error) {
	t, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return g.secret, nil
	})
	if err != nil || !t.Valid {
		return "", "", "", domain.ErrInvalidToken
	}
	c, ok := t.Claims.(*claims)
	if !ok {
		return "", "", "", domain.ErrInvalidToken
	}
	return c.RoomID, c.ParticipantID, domain.ParticipantRole(c.Role), nil
}
