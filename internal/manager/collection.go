package manager

import (
	"context"
	"sync"

	"github.com/gnarloqgames/ga-actor-poc/internal/model"
	"github.com/google/uuid"
)

type ActorCollection struct {
	mx *sync.Mutex

	actors  map[uuid.UUID]model.Actor
	factory actorFactory
}

func NewActorCollection(factoryFn actorFactory) *ActorCollection {
	return &ActorCollection{
		mx: &sync.Mutex{},

		actors:  make(map[uuid.UUID]model.Actor),
		factory: factoryFn,
	}
}

func (i *ActorCollection) Get(address model.Address) model.Actor {
	i.mx.Lock()
	defer i.mx.Unlock()

	inv, ok := i.actors[address.ID]
	if !ok {
		ctx := context.WithValue(context.Background(), model.KeyID, address.ID)
		inv = i.factory(ctx)
		i.actors[address.ID] = inv
	}

	return inv
}
