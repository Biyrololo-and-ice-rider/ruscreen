package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourorg/meet/internal/domain"
)

var now = time.Now().UTC()

// ---- Room ----

func TestNewRoom_OK(t *testing.T) {
	r, err := domain.NewRoom("abc12345", now)
	require.NoError(t, err)
	assert.Equal(t, "abc12345", r.ID)
	assert.Equal(t, domain.RoomStatusActive, r.Status)
	assert.Empty(t, r.Participants)
	assert.Nil(t, r.ScreenSharerID)
}

func TestRoom_CanJoin_TrueWhenActiveAndNotFull(t *testing.T) {
	r, _ := domain.NewRoom("x", now)
	assert.True(t, r.CanJoin())
}

func TestRoom_CanJoin_FalseWhenFull(t *testing.T) {
	r, _ := domain.NewRoom("x", now)
	for i := 0; i < 4; i++ {
		p, _ := domain.NewParticipant("id", "name", domain.RoleGuest, now)
		r.Participants = append(r.Participants, p)
	}
	assert.False(t, r.CanJoin())
}

func TestRoom_CanJoin_FalseWhenClosed(t *testing.T) {
	r, _ := domain.NewRoom("x", now)
	_ = r.Close(now)
	assert.False(t, r.CanJoin())
}

func TestRoom_Close_OK(t *testing.T) {
	r, _ := domain.NewRoom("x", now)
	err := r.Close(now)
	require.NoError(t, err)
	assert.Equal(t, domain.RoomStatusClosed, r.Status)
	assert.NotNil(t, r.ClosedAt)
}

func TestRoom_Close_AlreadyClosed(t *testing.T) {
	r, _ := domain.NewRoom("x", now)
	_ = r.Close(now)
	err := r.Close(now)
	assert.ErrorIs(t, err, domain.ErrRoomAlreadyClosed)
}

func TestRoom_HasScreenSharer(t *testing.T) {
	r, _ := domain.NewRoom("x", now)
	assert.False(t, r.HasScreenSharer())
	id := "p1"
	r.ScreenSharerID = &id
	assert.True(t, r.HasScreenSharer())
}

// ---- Participant ----

func TestNewParticipant_OK(t *testing.T) {
	p, err := domain.NewParticipant("id1", "Alice", domain.RoleOrganizer, now)
	require.NoError(t, err)
	assert.Equal(t, "Alice", p.Name)
	assert.True(t, p.CameraEnabled)
	assert.True(t, p.MicEnabled)
	assert.True(t, p.IsOrganizer())
}

func TestNewParticipant_EmptyName(t *testing.T) {
	_, err := domain.NewParticipant("id1", "", domain.RoleGuest, now)
	assert.ErrorIs(t, err, domain.ErrEmptyName)
}

func TestNewParticipant_NameTooLong(t *testing.T) {
	long := make([]rune, 51)
	for i := range long {
		long[i] = 'a'
	}
	_, err := domain.NewParticipant("id1", string(long), domain.RoleGuest, now)
	assert.ErrorIs(t, err, domain.ErrNameTooLong)
}

func TestParticipant_IsOrganizer(t *testing.T) {
	org, _ := domain.NewParticipant("id1", "Bob", domain.RoleOrganizer, now)
	guest, _ := domain.NewParticipant("id2", "Eve", domain.RoleGuest, now)
	assert.True(t, org.IsOrganizer())
	assert.False(t, guest.IsOrganizer())
}

// ---- ChatMessage ----

func TestNewChatMessage_OK(t *testing.T) {
	msg, err := domain.NewChatMessage("m1", "r1", "p1", "Alice", "hello", now)
	require.NoError(t, err)
	assert.Equal(t, "hello", msg.Text)
}

func TestNewChatMessage_EmptyText(t *testing.T) {
	_, err := domain.NewChatMessage("m1", "r1", "p1", "Alice", "", now)
	assert.ErrorIs(t, err, domain.ErrEmptyMessage)
}

func TestNewChatMessage_TooLong(t *testing.T) {
	long := make([]rune, 2001)
	for i := range long {
		long[i] = 'x'
	}
	_, err := domain.NewChatMessage("m1", "r1", "p1", "Alice", string(long), now)
	assert.ErrorIs(t, err, domain.ErrMessageTooLong)
}
