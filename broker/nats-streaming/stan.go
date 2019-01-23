// Package nats provides a NATS broker
package stan

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/codec/json"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
)

var (
	DefaultClusterID = "test-cluster"
	DefaultClientID  = "test-client"
)

type nbroker struct {
	sync.RWMutex
	addrs  []string
	conn   stan.Conn
	opts   broker.Options
	nopts  nats.Options
	snopts stan.Options
}

type subscriber struct {
	topic string
	s     stan.Subscription
	opts  broker.SubscribeOptions
}

type publication struct {
	m   *broker.Message
	msg *stan.Msg
}

func init() {
	cmd.DefaultBrokers["stan"] = NewBroker
}

func (n *publication) Topic() string {
	return n.msg.Subject
}

func (n *publication) Message() *broker.Message {
	return n.m
}

func (n *publication) Ack() error {
	return n.msg.Ack()
}

func (n *subscriber) Options() broker.SubscribeOptions {
	return n.opts
}

func (n *subscriber) Topic() string {
	return n.topic
}

func (n *subscriber) Unsubscribe() error {
	return n.s.Unsubscribe()
}

func (n *nbroker) Address() string {
	if n.conn != nil && n.conn.NatsConn().IsConnected() {
		return n.conn.NatsConn().ConnectedUrl()
	}
	if len(n.addrs) > 0 {
		return n.addrs[0]
	}

	return ""
}

func setAddrs(addrs []string) []string {
	var cAddrs []string
	for _, addr := range addrs {
		if len(addr) == 0 {
			continue
		}
		if !strings.HasPrefix(addr, "nats://") {
			addr = "nats://" + addr
		}
		cAddrs = append(cAddrs, addr)
	}
	if len(cAddrs) == 0 {
		cAddrs = []string{stan.DefaultNatsURL}
	}
	return cAddrs
}

func (n *nbroker) Connect() error {
	n.RLock()
	if n.conn != nil && n.conn.NatsConn().IsConnected() {
		n.RUnlock()
		return nil
	}
	n.RUnlock()

	opts := n.nopts
	opts.Servers = n.addrs
	opts.Secure = n.opts.Secure
	opts.TLSConfig = n.opts.TLSConfig

	// secure might not be set
	if n.opts.TLSConfig != nil {
		opts.Secure = true
	}

	nc, err := opts.Connect()
	if err != nil {
		return err
	}
	snc, err := stan.Connect(n.getClusterID(), n.getClientID(), stan.NatsConn(nc))
	if err != nil {
		return err
	}
	n.Lock()
	n.conn = snc
	n.Unlock()
	return nil
}

func (n *nbroker) Disconnect() error {
	n.RLock()
	n.conn.Close()
	n.RUnlock()
	return nil
}

func (n *nbroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	n.addrs = setAddrs(n.opts.Addrs)
	return nil
}

func (n *nbroker) Options() broker.Options {
	return n.opts
}

func (n *nbroker) Publish(topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	b, err := n.opts.Codec.Marshal(msg)
	if err != nil {
		return err
	}
	n.RLock()
	defer n.RUnlock()
	return n.conn.Publish(topic, b)
}

func (n *nbroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	opt := broker.SubscribeOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	var (
		err   error
		sub   stan.Subscription
		sopts []stan.SubscriptionOption
	)

	manualAcks := false
	if opt.Context != nil {
		manualAcks, _ = opt.Context.Value(manualAckModeKey{}).(bool)
		if manualAcks {
			sopts = append(sopts, stan.SetManualAckMode())
		}
	}

	if opt.Context != nil {
		if durableName, ok := opt.Context.Value(durableNameKey{}).(string); ok && len(durableName) > 0 {
			sopts = append(sopts, stan.DurableName(durableName))
		}
		if ackWait, ok := opt.Context.Value(ackWaitKey{}).(time.Duration); ok {
			sopts = append(sopts, stan.AckWait(ackWait))
		}
		if maxInflight, ok := opt.Context.Value(maxInflightKey{}).(int); ok {
			sopts = append(sopts, stan.MaxInflight(maxInflight))
		}
		if v, ok := opt.Context.Value(deliverAllAvailableKey{}).(bool); ok && v {
			sopts = append(sopts, stan.DeliverAllAvailable())
		} else if v, ok := opt.Context.Value(startWithLastReceivedKey{}).(bool); ok && v {
			sopts = append(sopts, stan.StartWithLastReceived())
		} else if v, ok := opt.Context.Value(startAtTimeKey{}).(time.Time); ok {
			sopts = append(sopts, stan.StartAtTime(v))
		} else if v, ok := opt.Context.Value(startAtSequenceKey{}).(uint64); ok {
			sopts = append(sopts, stan.StartAtSequence(v))
		}
	}

	fn := func(msg *stan.Msg) {
		var m broker.Message
		if opt.AutoAck && manualAcks {
			msg.Ack()
		}
		if err := n.opts.Codec.Unmarshal(msg.Data, &m); err != nil {
			return
		}
		if err := handler(&publication{m: &m, msg: msg}); err != nil {
			return
		}
	}

	n.RLock()
	if len(opt.Queue) > 0 {
		sub, err = n.conn.QueueSubscribe(topic, opt.Queue, fn, sopts...)
	} else {
		sub, err = n.conn.Subscribe(topic, fn, sopts...)
	}
	n.RUnlock()
	if err != nil {
		return nil, err
	}
	return &subscriber{topic: topic, s: sub, opts: opt}, nil
}

func (n *nbroker) String() string {
	return "stan"
}

func NewBroker(opts ...broker.Option) broker.Broker {
	options := broker.Options{
		// Default codec
		Codec:   json.Marshaler{},
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	natsOpts := nats.GetDefaultOptions()
	if n, ok := options.Context.Value(nOptionsKey{}).(nats.Options); ok {
		natsOpts = n
	}

	// broker.Options have higher priority than stan.Options
	// only if Addrs, Secure or TLSConfig were not set through a broker.Option
	// we read them from stan.Option
	if len(options.Addrs) == 0 {
		options.Addrs = natsOpts.Servers
	}

	if !options.Secure {
		options.Secure = natsOpts.Secure
	}

	if options.TLSConfig == nil {
		options.TLSConfig = natsOpts.TLSConfig
	}

	stanOpts := stan.DefaultOptions
	if n, ok := options.Context.Value(snOptionsKey{}).(stan.Options); ok {
		stanOpts = n
	}

	nb := &nbroker{
		opts:   options,
		nopts:  natsOpts,
		snopts: stanOpts,
		addrs:  setAddrs(options.Addrs),
	}

	return nb
}

func (n *nbroker) getClusterID() string {
	if e, ok := n.opts.Context.Value(clusterIdKey{}).(string); ok {
		return e
	}
	return DefaultClusterID
}

func (n *nbroker) getClientID() string {
	if e, ok := n.opts.Context.Value(clientIdKey{}).(string); ok {
		return e
	}
	return DefaultClientID
}
