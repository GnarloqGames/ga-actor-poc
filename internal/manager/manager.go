package manager

import (
	"context"
	"log/slog"
	"time"

	"github.com/gnarloqgames/ga-actor-poc/internal/model"
	"google.golang.org/protobuf/proto"
)

type Manager struct {
	actors *ActorCollection
}

func NewManager() *Manager {
	return &Manager{
		actors: NewActorCollection(),
	}
}

func (m *Manager) Send(ctx context.Context, address model.Address, msg proto.Message, timeout time.Duration) error {
	actor := m.actors.Get(address)

	slog.Info("sending message to actor",
		"recipient_kind", actor.GetKind(),
		"recipient_id", actor.GetID(),
		"recipient_name", actor.GetName(),
	)

	return actor.Receive(ctx, msg, nil)
}
