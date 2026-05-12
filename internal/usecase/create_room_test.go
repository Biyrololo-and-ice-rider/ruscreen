package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourorg/meet/internal/domain"
	"github.com/yourorg/meet/internal/usecase"
)

// --- hand-written mocks (replace with gomock-generated via `make generate`) ---

type mockRoomRepo struct {
	saved      *domain.Room
	findResult domain.Room
	findErr    error
	saveErr    error
}

func (m *mockRoomRepo) Save(_ context.Context, r domain.Room) error {
	m.saved = &r
	return m.saveErr
}
func (m *mockRoomRepo) FindByID(_ context.Context, _ string) (domain.Room, error) {
	return m.findResult, m.findErr
}
func (m *mockRoomRepo) UpdateStatus(_ context.Context, _ string, _ domain.RoomStatus) error {
	return nil
}
func (m *mockRoomRepo) AddParticipant(_ context.Context, _ string, _ domain.Participant) error {
	return nil
}
func (m *mockRoomRepo) RemoveParticipant(_ context.Context, _ string, _ string) error { return nil }
func (m *mockRoomRepo) SetScreenSharer(_ context.Context, _ string, _ *string) error  { return nil }

type mockTokenGen struct {
	token  string
	genErr error
}

func (m *mockTokenGen) NewToken(_, _ string, _ domain.ParticipantRole) (string, error) {
	return m.token, m.genErr
}
func (m *mockTokenGen) Verify(t string) (string, string, domain.ParticipantRole, error) {
	switch t {
	case "valid-organizer":
		return "room1", "p1", domain.RoleOrganizer, nil
	case "valid-guest":
		return "room1", "p2", domain.RoleGuest, nil
	default:
		return "", "", "", domain.ErrInvalidToken
	}
}

// ---- helpers ----

func activeRoom(t *testing.T) domain.Room {
	t.Helper()
	r, err := domain.NewRoom("room1", time.Now().UTC())
	require.NoError(t, err)
	return r
}

// ---- CreateRoom ----

func TestCreateRoom_OK(t *testing.T) {
	repo := &mockRoomRepo{}
	tok := &mockTokenGen{token: "signed-jwt"}
	uc := usecase.NewCreateRoom(repo, tok)

	resp, err := uc.Execute(context.Background(), usecase.CreateRoomRequest{OrganizerName: "Alice"})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.RoomID)
	assert.Equal(t, "/room/"+resp.RoomID, resp.InviteURL)
	assert.Equal(t, "signed-jwt", resp.OrganizerToken)
	require.NotNil(t, repo.saved)
	assert.Equal(t, domain.RoomStatusActive, repo.saved.Status)
	assert.Len(t, repo.saved.Participants, 1)
	assert.True(t, repo.saved.Participants[0].IsOrganizer())
}

func TestCreateRoom_EmptyName(t *testing.T) {
	uc := usecase.NewCreateRoom(&mockRoomRepo{}, &mockTokenGen{})
	_, err := uc.Execute(context.Background(), usecase.CreateRoomRequest{OrganizerName: ""})
	assert.ErrorIs(t, err, domain.ErrEmptyName)
}

func TestCreateRoom_RepoFailure(t *testing.T) {
	repo := &mockRoomRepo{saveErr: errors.New("redis down")}
	uc := usecase.NewCreateRoom(repo, &mockTokenGen{token: "t"})
	_, err := uc.Execute(context.Background(), usecase.CreateRoomRequest{OrganizerName: "Alice"})
	assert.Error(t, err)
}

// ---- JoinRoom ----

func TestJoinRoom_OK(t *testing.T) {
	repo := &mockRoomRepo{findResult: activeRoom(t)}
	tok := &mockTokenGen{token: "guest-jwt"}
	uc := usecase.NewJoinRoom(repo, tok)

	resp, err := uc.Execute(context.Background(), usecase.JoinRoomRequest{
		RoomID:          "room1",
		ParticipantName: "Bob",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.ParticipantID)
	assert.Equal(t, "guest-jwt", resp.ParticipantToken)
}

func TestJoinRoom_EmptyName(t *testing.T) {
	uc := usecase.NewJoinRoom(&mockRoomRepo{}, &mockTokenGen{})
	_, err := uc.Execute(context.Background(), usecase.JoinRoomRequest{
		RoomID:          "r1",
		ParticipantName: "",
	})
	assert.ErrorIs(t, err, domain.ErrEmptyName)
}

func TestJoinRoom_RoomNotFound(t *testing.T) {
	repo := &mockRoomRepo{findErr: domain.ErrRoomNotFound}
	uc := usecase.NewJoinRoom(repo, &mockTokenGen{})
	_, err := uc.Execute(context.Background(), usecase.JoinRoomRequest{
		RoomID:          "x",
		ParticipantName: "Bob",
	})
	assert.ErrorIs(t, err, domain.ErrRoomNotFound)
}

func TestJoinRoom_RoomFull(t *testing.T) {
	room := activeRoom(t)
	for i := 0; i < 4; i++ {
		p, _ := domain.NewParticipant("p", "name", domain.RoleGuest, time.Now())
		room.Participants = append(room.Participants, p)
	}
	repo := &mockRoomRepo{findResult: room}
	uc := usecase.NewJoinRoom(repo, &mockTokenGen{token: "t"})
	_, err := uc.Execute(context.Background(), usecase.JoinRoomRequest{
		RoomID:          "room1",
		ParticipantName: "Latecomer",
	})
	assert.ErrorIs(t, err, domain.ErrRoomFull)
}

func TestJoinRoom_RoomClosed(t *testing.T) {
	room := activeRoom(t)
	_ = room.Close(time.Now())
	repo := &mockRoomRepo{findResult: room}
	uc := usecase.NewJoinRoom(repo, &mockTokenGen{token: "t"})
	_, err := uc.Execute(context.Background(), usecase.JoinRoomRequest{
		RoomID:          "room1",
		ParticipantName: "Bob",
	})
	assert.ErrorIs(t, err, domain.ErrRoomClosed)
}
