package handler

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
)

type Cut struct {
	cutReady   *atomic.Bool
	cutFlag    *atomic.Bool
	channel    chan *string
	cutContent []string
}

func NewCut() *Cut {
	c := Cut{
		cutReady: &atomic.Bool{},
		cutFlag:  &atomic.Bool{},
		channel:  make(chan *string, bufferSize),
	}
	c.cutFlag.Store(false)
	return &c
}

func (c *Cut) Ready() bool {
	return c.cutReady.Load()
}

func (c *Cut) Channel() chan *string {
	return c.channel
}

func (c *Cut) Run(ctx context.Context) error {
	c.cutReady.Store(true)
	defer c.cutReady.Store(false)
	log.Info("Begin output cut")
	for {
		select {
		case <-ctx.Done():
			log.Info("Stop output cut")
			return nil
		case s := <-c.channel:
			if c.cutFlag.Load() {
				c.cutContent = append(c.cutContent, *s)
			}
		}
	}
}

func (c *Cut) BeginCut() error {
	if !c.cutReady.Load() {
		return fmt.Errorf("cut is not ready, begin cut failed")
	}
	if c.cutFlag.Load() {
		return fmt.Errorf("cutting now, begin cut failed")
	}
	c.cutFlag.Store(true)
	return nil
}

func (c *Cut) StopCut() (string, error) {
	if !c.cutFlag.Load() {
		return "", fmt.Errorf("not cutting, stop cut failed")
	}
	c.cutFlag.Store(false)
	content := strings.Join(c.cutContent, "\n")
	log.Debug("Cut success: ", content)
	return content, nil
}
