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

	indices []uuid.UUID
	items   map[uuid.UUID]T
}

func NewQueue[T Queueable]() *Queue[T] {
	return &Queue[T]{
		mx: &sync.Mutex{},

		indices: make([]uuid.UUID, 0),
		items:   make(map[uuid.UUID]T),
	}
}

func (q *Queue[T]) len() int {
	return len(q.indices)
}

func (q *Queue[T]) Unshift() T {
	q.mx.Lock()
	defer q.mx.Unlock()

	len := q.len()
	if len == 0 {
		return nil
	}

	index := q.indices[0]
	val := q.items[index]

	if len == 1 {
		q.indices = make([]uuid.UUID, 0)
	} else {
		q.indices = q.indices[1:]
	}

	delete(q.items, index)

	return val
}

func (q *Queue[T]) Push(item T) int {
	q.mx.Lock()
	defer q.mx.Unlock()

	index := uuid.New()
	q.indices = append(q.indices, index)
	q.items[index] = item

	return q.len()
}

type InventoryActor struct {
	ID uuid.UUID

	timers map[uuid.UUID]time.Timer

	Buildings *Collection[Building]
	Resources *Collection[Resource]

	BuildQueue *Queue[*message.BuildRequest]
}

func NewInventoryActor(id uuid.UUID) *InventoryActor {
	actor := &InventoryActor{
		ID: id,

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

func (a *InventoryActor) Receive(ctx context.Context, msg proto.Message, res proto.Message) error {
	req, ok := msg.(*message.BuildRequest)
	if !ok {
		return fmt.Errorf("invalid message type")
	}

	slog.Info("actor received message",
		"actor_kind", a.GetKind(),
		"actor_id", a.GetID(),
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
