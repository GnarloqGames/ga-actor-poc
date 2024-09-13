package actor

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type TimerStatus string

const (
	AttributeTimerID  string = "id"
	AttributeQueueID  string = "queue_id"
	AttributeDuration string = "duration"
	AttributeStatus   string = "status"

	StatusFailed    TimerStatus = "failed"
	StatusDone      TimerStatus = "done"
	StatusCancelled TimerStatus = "cancelled"
)

type TimerActor struct {
	ID       uuid.UUID
	QueueID  uuid.UUID
	Duration time.Duration

	stop chan struct{}
}

func (t TimerActor) Attributes() []any {
	return []any{
		AttributeTimerID, t.ID.String(),
		AttributeQueueID, t.QueueID.String(),
		AttributeDuration, t.Duration.String(),
	}
}

type TimerReply struct {
	ID      uuid.UUID
	QueueID uuid.UUID
	Status  TimerStatus
}

func (t TimerReply) Attributes() []any {
	return []any{
		AttributeTimerID, t.ID.String(),
		AttributeQueueID, t.QueueID.String(),
		AttributeStatus, t.Status,
	}
}

func NewTimerActor(queueID uuid.UUID, duration time.Duration, replyChan chan TimerReply) *TimerActor {
	actor := &TimerActor{
		ID:       uuid.New(),
		QueueID:  queueID,
		Duration: duration,
		stop:     make(chan struct{}),
	}
	actorAttributes := actor.Attributes()

	go func() {
		timer := time.NewTimer(duration)

		reply := TimerReply{
			ID:      actor.ID,
			QueueID: actor.QueueID,
			Status:  StatusFailed,
		}

	Loop:
		for {
			select {
			case <-timer.C:
				slog.Info("timer expired",
					actorAttributes...,
				)
				reply.Status = StatusDone
				break Loop
			case <-actor.stop:
				slog.Info("timer received stop signal",
					actorAttributes...,
				)
				timer.Stop()
				reply.Status = StatusCancelled
				break Loop
			}
		}

		replyChan <- reply
	}()

	return actor
}

func (t *TimerActor) Stop() {
	t.stop <- struct{}{}
}
