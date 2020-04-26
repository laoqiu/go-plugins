// Package nats provides a NATS broker
// From https://github.com/micro/go-micro/blob/master/broker/nats/options.go
package nats

import (
	"context"

	"github.com/laoqiu/go-plugins/broker"
	nats "github.com/nats-io/nats.go"
)

type optionsKey struct{}
type drainConnectionKey struct{}

// setBrokerOption returns a function to setup a context with given value
func setBrokerOption(k, v interface{}) broker.Option {
	return func(o *broker.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

// Options accepts nats.Options
func Options(opts nats.Options) broker.Option {
	return setBrokerOption(optionsKey{}, opts)
}

// DrainConnection will drain subscription on close
func DrainConnection() broker.Option {
	return setBrokerOption(drainConnectionKey{}, struct{}{})
}
