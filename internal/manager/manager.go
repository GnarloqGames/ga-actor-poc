package manager

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gnarloqgames/ga-actor-poc/internal/model"
	"google.golang.org/protobuf/proto"
)

type actorFactory func(ctx context.Context) model.Actor

type Manager struct {
	actors map[string]*ActorCollection
}

func NewManager() *Manager {
	return &Manager{
		actors: make(map[string]*ActorCollection),
	}
}

func (m *Manager) NewKind(kind string, factory actorFactory) error {
	if _, ok := m.actors[kind]; ok {
		return fmt.Errorf("kind is already registered")
	}

	collection := NewActorCollection(factory)
	m.actors[kind] = collection

	return nil
}

func (m *Manager) Send(ctx context.Context, address model.Address, msg proto.Message, timeout time.Duration) error {
	actorCollection, ok := m.actors[address.Kind]
	if !ok {
		return fmt.Errorf("kind %s is not registered", address.Kind)
	}

	actor := actorCollection.Get(address)

	slog.Info("sending message to actor",
		"recipient_kind", actor.GetKind(),
		"recipient_id", actor.GetID(),
	)

	return actor.Receive(ctx, msg, nil)
}
