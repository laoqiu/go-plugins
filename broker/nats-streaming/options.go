package stan

import (
	"context"
	"time"

	"github.com/micro/go-micro/broker"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
)

type nOptionsKey struct{}
type snOptionsKey struct{}
type clusterIdKey struct{}
type clientIdKey struct{}
type manualAckModeKey struct{}
type durableNameKey struct{}
type deliverAllAvailableKey struct{}
type startAtTimeKey struct{}
type startAtSequenceKey struct{}
type startWithLastReceivedKey struct{}
type ackWaitKey struct{}
type maxInflightKey struct{}

// NatsOptions accepts nats.Options
func NatsOptions(opts nats.Options) broker.Option {
	return func(o *broker.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, nOptionsKey{}, opts)
	}
}

func StanOptions(opts stan.Options) broker.Option {
	return func(o *broker.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, snOptionsKey{}, opts)
	}
}

func ClusterID(e string) broker.Option {
	return func(o *broker.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, clusterIdKey{}, e)
	}
}

func ClientID(e string) broker.Option {
	return func(o *broker.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, clientIdKey{}, e)
	}
}

func SetManualAckMode() broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, manualAckModeKey{}, true)
	}
}

func DurableName(name string) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, durableNameKey{}, name)
	}
}

func DeliverAllAvailable() broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, deliverAllAvailableKey{}, true)
	}
}

func StartAtTime(start time.Time) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, startAtTimeKey{}, start)
	}
}

func StartAtSequence(seq uint64) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, startAtSequenceKey{}, seq)
	}
}

func StartWithLastReceived() broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, startWithLastReceivedKey{}, true)
	}
}

func AckWait(t time.Duration) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, ackWaitKey{}, t)
	}
}

func MaxInflight(n int) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, maxInflightKey{}, n)
	}
}
