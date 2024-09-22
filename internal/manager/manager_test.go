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

	ctx := context.WithValue(context.Background(), model.KeyID, uuid.New())
	existingInventory := actor.InventoryActorFactory(ctx)
	test2ID := createdAddress.Hash()

	tests := []struct {
		manager    *Manager
		id         uuid.UUID
		expectedID string
	}{
		{
			manager: &Manager{
				actors: map[string]*ActorCollection{
					"inventory": {
						mx: &sync.Mutex{},
						actors: map[uuid.UUID]model.Actor{
							existingInventory.GetID(): existingInventory,
						},
					},
				},
			},
			id:         existingInventory.GetID(),
			expectedID: existingInventory.GetID().String(),
		},
		{
			manager:    NewManager(),
			id:         test2ID,
			expectedID: test2ID.String(),
		},
	}

	for _, tt := range tests {
		tf := func(t *testing.T) {
			if len(tt.manager.actors) == 0 {
				err := tt.manager.NewKind("inventory", actor.InventoryActorFactory)
				require.NoError(t, err)
			}
			address := model.Address{
				Kind: "inventory",
				ID:   tt.id,
			}
			actors, ok := tt.manager.actors[address.Kind]
			require.True(t, ok)
			inv := actors.Get(address).(*actor.InventoryActor)

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
