package middleware

import log "github.com/sirupsen/logrus"

type Options struct {
	recoveryFunc RecoveryHandlerFunc
	authFunc     AuthHandlerFunc
	logLevel     log.Level
}

type Option func(*Options)

func WithRecoverFunc(f RecoveryHandlerFunc) Option {
	return func(o *Options) {
		o.recoveryFunc = f
	}
}

func WithAuthFunc(f AuthHandlerFunc) Option {
	return func(o *Options) {
		o.authFunc = f
	}
}

func WithLogLevel(level log.Level) Option {
	return func(o *Options) {
		o.logLevel = level
	}
}
