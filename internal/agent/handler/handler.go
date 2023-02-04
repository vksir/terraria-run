package handler

import (
	"context"
	"go.uber.org/zap"
)

const (
	bufferSize = 128
	reportSize = 32
)

var log = zap.S()

type Handler interface {
	Ready() bool
	Channel() chan *string
	Start(context.Context) error
}
