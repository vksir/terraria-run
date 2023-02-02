package handler

import (
	"context"
	"os"
	"sync/atomic"
	"terraria-run/internal/common/constant"
)

type Record struct {
	ready   *atomic.Bool
	channel chan *string
}

func NewRecord() *Record {
	r := Record{
		ready:   &atomic.Bool{},
		channel: make(chan *string, bufferSize),
	}
	return &r
}

func (r *Record) Ready() bool {
	return r.ready.Load()
}

func (r *Record) Channel() chan *string {
	return r.channel
}

func (r *Record) Run(ctx context.Context) error {
	w, err := os.OpenFile(constant.TModLoaderLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		return err
	}
	defer func(w *os.File) {
		err := w.Close()
		if err != nil {
			log.Error("Close file failed", err)
		}
	}(w)
	log.Info("Begin output record")
	for {
		select {
		case <-ctx.Done():
			log.Info("Stop output record")
			return nil
		case s := <-r.channel:
			log.Debugf("[tModLoader] %s", *s)
			_, err = w.Write([]byte(*s + "\n"))
			if err != nil {
				log.Error("Write record failed: ", err)
			}
		}
	}
}
