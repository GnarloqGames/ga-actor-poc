package actor

import (
	"testing"

	"github.com/gnarloqgames/ga-actor-poc/message"
	"github.com/stretchr/testify/require"
)

func TestQueueUnshift(t *testing.T) {
	tests := []struct {
		label         string
		queue         []*message.BuildRequest
		expectedValue *message.BuildRequest
		expectedQueue []*message.BuildRequest
	}{
		{
			label:         "empty",
			queue:         make([]*message.BuildRequest, 0),
			expectedValue: nil,
			expectedQueue: make([]*message.BuildRequest, 0),
		},
		{
			label: "one",
			queue: []*message.BuildRequest{
				{Name: "test_1"},
			},
			expectedValue: &message.BuildRequest{Name: "test_1"},
			expectedQueue: make([]*message.BuildRequest, 0),
		},
		{
			label: "more",
			queue: []*message.BuildRequest{
				{Name: "test_1"},
				{Name: "test_2"},
			},
			expectedValue: &message.BuildRequest{Name: "test_1"},
			expectedQueue: []*message.BuildRequest{
				{Name: "test_2"},
			},
		},
	}

	for _, tt := range tests {
		queue := NewQueue[*message.BuildRequest]()
		for _, item := range tt.queue {
			queue.Push(item) //nolint
		}

		tf := func(t *testing.T) {
			val := queue.Unshift()

			if tt.expectedValue == nil {
				require.Nil(t, val)
			} else {
				require.Equal(t, tt.expectedValue.Name, val.Name)
			}

			actualItems := make([]*message.BuildRequest, 0)
			for _, item := range queue.items {
				actualItems = append(actualItems, item)
			}
			require.ElementsMatch(t, tt.expectedQueue, actualItems)
		}

		t.Run(tt.label, tf)
	}
}
