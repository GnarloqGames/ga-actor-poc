package manager

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gnarloqgames/ga-actor-poc/internal/actor"
	"github.com/gnarloqgames/ga-actor-poc/internal/model"
	"github.com/gnarloqgames/ga-actor-poc/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetInventory(t *testing.T) {
	createdAddress := model.Address{
		Kind: "inventory",
		ID:   uuid.New(),
	}

	existingInventory := actor.NewInventoryActor(uuid.New())
	test2ID := createdAddress.Hash()

	tests := []struct {
		manager    *Manager
		id         uuid.UUID
		expectedID string
	}{
		{
			manager: &Manager{
				actors: &ActorCollection{
					mx: &sync.Mutex{},
					actors: map[uuid.UUID]model.Actor{
						existingInventory.ID: existingInventory,
					},
				},
			},
			id:         existingInventory.ID,
			expectedID: existingInventory.ID.String(),
		},
		{
			manager:    NewManager(),
			id:         test2ID,
			expectedID: test2ID.String(),
		},
	}

	for _, tt := range tests {
		tf := func(t *testing.T) {
			address := model.Address{
				Kind: "inventory",
				ID:   tt.id,
			}
			inv := tt.manager.actors.Get(address).(*actor.InventoryActor)

			require.Equal(t, tt.expectedID, inv.ID.String())
		}

		t.Run(tt.id.String(), tf)
	}
}

func TestSend(t *testing.T) {
	manager := NewManager()
	address := model.Address{
		Kind: "inventory",
		ID:   uuid.New(),
	}

	request := &message.BuildRequest{
		Name:     "test",
		Duration: "10s",
	}

	err := manager.Send(context.Background(), address, request, 10*time.Second)

	require.NoError(t, err)
}
