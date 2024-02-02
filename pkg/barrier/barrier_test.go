package barrier

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBarrier(t *testing.T) {
	t.Parallel()
	t.Run("successfully record bad event", func(t *testing.T) {
		b := New(10, 1*time.Second)
		b.RecordBadEvent()
		assert.Equal(t, 9, b.Value)
	})
	t.Run("errors if the instance capacity is set to a negative number", func(t *testing.T) {
		ctx := context.Background()
		b := New(-1, 1*time.Second)
		assert.Error(t, b.Start(ctx))
	})
	t.Run("record bad event doesn't go to a negative number", func(t *testing.T) {
		b := New(0, 1*time.Second)
		b.RecordBadEvent()
		assert.Equal(t, 0, b.Value)
	})
	t.Run("isBadState returns true when instance reaches 0", func(t *testing.T) {
		b := New(0, 1*time.Second)
		assert.True(t, true, b.IsBadState())
	})
	t.Run("successfully refill barrier at the refill rate", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		b := New(10, 1*time.Second)
		go func() {
			require.NoError(t, b.Start(ctx))
		}()
		b.RecordBadEvent()
		b.RecordBadEvent()
		assert.Equal(t, 8, b.Value)
		time.Sleep(3 * time.Second)
		cancel()
		assert.Equal(t, 10, b.Value)
	})
}
