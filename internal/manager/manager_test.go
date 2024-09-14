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
	name1 := "existing_actor"
	createdAddress := model.Address{
		Kind: "inventory",
		Name: "created_actor",
	}

	existingInventory := actor.NewInventoryActor(name1)
	test2ID := createdAddress.Hash()

	tests := []struct {
		manager      *Manager
		name         string
		expectedID   string
		expectedName string
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
			name:         name1,
			expectedID:   existingInventory.ID.String(),
			expectedName: existingInventory.Name,
		},
		{
			manager:      NewManager(),
			name:         createdAddress.Name,
			expectedID:   test2ID.String(),
			expectedName: createdAddress.Name,
		},
	}

	for _, tt := range tests {
		tf := func(t *testing.T) {
			address := model.Address{
				Kind: "inventory",
				Name: tt.name,
			}
			inv := tt.manager.actors.Get(address).(*actor.InventoryActor)

			require.Equal(t, tt.expectedID, inv.ID.String())
			require.Equal(t, tt.expectedName, inv.Name)
		}

		t.Run(tt.name, tf)
	}
}

func TestSend(t *testing.T) {
	manager := NewManager()
	address := model.Address{
		Kind: "inventory",
		Name: "test",
	}

	request := &message.BuildRequest{
		Name:     "test",
		Duration: "10s",
	}

	err := manager.Send(context.Background(), address, request, 10*time.Second)

	require.NoError(t, err)
}
