package actor

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueueUnshift(t *testing.T) {
	tests := []struct {
		label         string
		queue         []*BuildRequest
		expectedValue *BuildRequest
		expectedQueue []*BuildRequest
	}{
		{
			label:         "empty",
			queue:         make([]*BuildRequest, 0),
			expectedValue: nil,
			expectedQueue: make([]*BuildRequest, 0),
		},
		{
			label: "one",
			queue: []*BuildRequest{
				{name: "test_1"},
			},
			expectedValue: &BuildRequest{name: "test_1"},
			expectedQueue: make([]*BuildRequest, 0),
		},
		{
			label: "more",
			queue: []*BuildRequest{
				{name: "test_1"},
				{name: "test_2"},
			},
			expectedValue: &BuildRequest{name: "test_1"},
			expectedQueue: []*BuildRequest{
				{name: "test_2"},
			},
		},
	}

	for _, tt := range tests {
		tf := func(t *testing.T) {
			queue := Queue[*BuildRequest]{
				mx:    &sync.Mutex{},
				items: tt.queue,
			}

			val := queue.Unshift()

			if tt.expectedValue == nil {
				require.Nil(t, val)
			} else {
				require.Equal(t, tt.expectedValue.name, val.name)
			}

			require.ElementsMatch(t, tt.expectedQueue, queue.items)
		}

		t.Run(tt.label, tf)
	}
}
