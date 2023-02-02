package output

import (
	"bufio"
	"context"
	"go.uber.org/zap"
	"io"
	"terraria-run/internal/agent/handler"
)

var log = zap.S()

type Listener struct {
	outputChannels []chan *string

	// handler
	Record *handler.Record
	Report *handler.Report
	Cut    *handler.Cut
}

func NewListener() *Listener {
	l := Listener{
		Record: handler.NewRecord(),
		Report: handler.NewReport(),
		Cut:    handler.NewCut(),
	}
	l.outputChannels = append(l.outputChannels, l.Record.Channel())
	l.outputChannels = append(l.outputChannels, l.Report.Channel())
	l.outputChannels = append(l.outputChannels, l.Cut.Channel())
	return &l
}

func (l *Listener) Run(ctx context.Context, r io.Reader) error {
	if err := l.startHandlers(ctx); err != nil {
		return err
	}
	log.Info("Begin listen output")
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		out := scanner.Text()
		for _, c := range l.outputChannels {
			c <- &out
		}
		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}
	return nil
}

func (l *Listener) startHandlers(ctx context.Context) error {
	go func() {
		err := l.Record.Run(ctx)
		if err != nil {
			log.Error("Record handler stopped: ", err)
		}
	}()
	go func() {
		err := l.Report.Run(ctx)
		if err != nil {
			log.Error("Report handler stopped: ", err)
		}
	}()
	go func() {
		err := l.Cut.Run(ctx)
		if err != nil {
			log.Error("Cut handler stopped: ", err)
		}
	}()
	return nil
}
