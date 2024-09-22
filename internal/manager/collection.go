package manager

import (
	"sync"

	"github.com/gnarloqgames/ga-actor-poc/internal/actor"
	"github.com/gnarloqgames/ga-actor-poc/internal/model"
	"github.com/google/uuid"
)

type ActorCollection struct {
	mx *sync.Mutex

	actors map[uuid.UUID]model.Actor
}

func NewActorCollection() *ActorCollection {
	return &ActorCollection{
		mx: &sync.Mutex{},

		actors: make(map[uuid.UUID]model.Actor),
	}
}

func (i *ActorCollection) Get(address model.Address) model.Actor {
	i.mx.Lock()
	defer i.mx.Unlock()

	inv, ok := i.actors[address.Hash()]
	if !ok {
		inv = actor.NewInventoryActor(address.ID)
	}

	return inv
}
