package actor

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTimer(t *testing.T) {
	tests := []struct {
		label          string
		duration       time.Duration
		timeout        time.Duration
		cancelAfter    time.Duration
		expectedStatus TimerStatus
		expectedError  error
	}{
		{
			label:          "done",
			duration:       100 * time.Millisecond,
			timeout:        200 * time.Millisecond,
			cancelAfter:    0,
			expectedStatus: StatusDone,
			expectedError:  nil,
		},
		{
			label:          "deadline",
			duration:       600 * time.Millisecond,
			timeout:        100 * time.Millisecond,
			cancelAfter:    0,
			expectedStatus: StatusFailed,
			expectedError:  context.DeadlineExceeded,
		},
		{
			label:          "cancelled",
			duration:       500 * time.Millisecond,
			timeout:        200 * time.Millisecond,
			cancelAfter:    100 * time.Millisecond,
			expectedStatus: StatusCancelled,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		tf := func(t *testing.T) {
			queue_id := uuid.New()
			reply := make(chan TimerReply)

			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			a := NewTimerActor(queue_id, tt.duration, reply)

			if tt.cancelAfter > 0 {
				go func() {
					t := time.NewTimer(tt.cancelAfter)
					<-t.C
					a.Stop()
				}()
			}

			actualStatus := StatusFailed
			var actualError error
		Loop:
			for {
				select {
				case r := <-reply:
					actualStatus = r.Status
					slog.Info("reply received", r.Attributes()...)
					break Loop
				case <-ctx.Done():
					if err := ctx.Err(); err != nil {
						actualError = err
					}
					break Loop
				}
			}

			cancel()

			require.Equal(t, tt.expectedStatus, actualStatus)
			require.Equal(t, tt.expectedError, actualError)
		}

		t.Run(tt.label, tf)
	}
}
