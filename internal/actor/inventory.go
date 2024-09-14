package actor

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gnarloqgames/ga-actor-poc/internal/model"
	"github.com/gnarloqgames/ga-actor-poc/message"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

var _ model.Actor = (*InventoryActor)(nil)

type Building struct {
	id   uuid.UUID
	name string
}

type Resource struct {
	id     uuid.UUID
	name   string
	amount uint
}

type StoredResource struct {
	id        uuid.UUID
	name      string
	amount    uint
	temporary bool
}

type Storeable interface {
	Building | Resource | StoredResource
}

type Collection[T Storeable] struct {
	mx *sync.Mutex

	items map[uuid.UUID]T
}

func NewCollection[T Storeable]() *Collection[T] {
	return &Collection[T]{
		mx: &sync.Mutex{},

		items: make(map[uuid.UUID]T),
	}
}

type Queueable interface {
	*message.BuildRequest
}

type Queue[T Queueable] struct {
	mx *sync.Mutex

	items []T
}

func NewQueue[T Queueable]() *Queue[T] {
	return &Queue[T]{
		mx:    &sync.Mutex{},
		items: make([]T, 0),
	}
}

func (q *Queue[T]) len() int {
	return len(q.items)
}

func (q *Queue[T]) Unshift() T {
	q.mx.Lock()
	defer q.mx.Unlock()

	len := q.len()
	if len == 0 {
		return nil
	}

	val := q.items[0]

	if len == 1 {
		q.items = make([]T, 0)
	} else {
		q.items = q.items[1:]
	}

	return val
}

func (q *Queue[T]) Push(item T) int {
	q.mx.Lock()
	defer q.mx.Unlock()

	q.items = append(q.items, item)

	return q.len()
}

type InventoryActor struct {
	ID   uuid.UUID
	Name string

	timers map[uuid.UUID]time.Timer

	Buildings *Collection[Building]
	Resources *Collection[Resource]

	BuildQueue *Queue[*message.BuildRequest]
}

func NewInventoryActor(name string) *InventoryActor {
	address := model.Address{
		Kind: "inventory",
		Name: name,
	}

	actor := &InventoryActor{
		ID:   address.Hash(),
		Name: name,

		timers: make(map[uuid.UUID]time.Timer),

		Buildings: NewCollection[Building](),
		Resources: NewCollection[Resource](),

		BuildQueue: NewQueue[*message.BuildRequest](),
	}

	go func() {
		for {
			if actor.BuildQueue.len() == 0 {
				continue
			}

			task := actor.BuildQueue.Unshift()

			slog.Warn("task done", "name", task.Name)
		}
	}()

	return actor
}

func (a *InventoryActor) GetID() string {
	return a.ID.String()
}

func (a *InventoryActor) GetKind() string {
	return "inventory"
}

func (a *InventoryActor) GetName() string {
	return a.Name
}

func (a *InventoryActor) Receive(ctx context.Context, msg proto.Message, res proto.Message) error {
	req, ok := msg.(*message.BuildRequest)
	if !ok {
		return fmt.Errorf("invalid message type")
	}

	slog.Info("actor received message",
		"actor_kind", a.GetKind(),
		"actor_id", a.GetID(),
		"actor_name", a.GetName(),
		"building_name", req.Name,
	)

	newLen := a.BuildQueue.Push(req)

	slog.Info("added request to build queue",
		"name", req.Name,
		"duration", req.Duration,
		"len", newLen,
	)

	return nil
}

func (a *InventoryActor) Start(ctx context.Context) {
	slog.Info("starting actor", "kind", "inventory", "id", a.ID.String())
}

func (a *InventoryActor) Destroy(ctx context.Context) {
	slog.Info("stopping actor", "kind", "inventory", "id", a.ID.String())
}

type BuildResponse struct {
}
