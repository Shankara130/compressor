package queue

import (
	"context"
	"testing"
	"time"

	"github.com/Shankara130/compressor/internal/domain/entity"
)

func TestInMemoryJobQueue_RoundTrip(t *testing.T) {
	q := NewInMemoryJobQueue()

	job := entity.Job{ID: "job-1", MimeType: "image/png"}
	if err := q.Enqueue(job); err != nil {
		t.Fatalf("Enqueue: %v", err)
	}

	got, err := q.Dequeue(context.Background())
	if err != nil {
		t.Fatalf("Dequeue: %v", err)
	}
	if got.ID != job.ID {
		t.Fatalf("expected job %q, got %q", job.ID, got.ID)
	}
}

// TestInMemoryJobQueue_DequeueUnblocksOnCancel verifies the shutdown fix:
// a worker blocked in Dequeue on an empty queue must return when its context
// is cancelled. Otherwise graceful worker shutdown hangs until the timeout.
func TestInMemoryJobQueue_DequeueUnblocksOnCancel(t *testing.T) {
	q := NewInMemoryJobQueue()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		_, err := q.Dequeue(ctx)
		errCh <- err
	}()

	// Give the goroutine time to park on the empty channel.
	select {
	case err := <-errCh:
		t.Fatalf("Dequeue returned before cancel: %v", err)
	case <-time.After(50 * time.Millisecond):
	}

	cancel()

	select {
	case err := <-errCh:
		if err != context.Canceled {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Dequeue did not unblock after context cancel")
	}
}
