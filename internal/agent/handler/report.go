package handler

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"sync/atomic"
	"time"
)

const (
	TypeServerActive = "SERVER_ACTIVE"
)

var rawEvents = getEvents()

type Report struct {
	ready   *atomic.Bool
	lock    *sync.Mutex
	channel chan *string
	events  []*Event
}

type Event struct {
	PatternString string
	Pattern       *regexp.Regexp
	Format        string
	Level         string
	Time          int64
	Msg           string
	Type          string
}

func NewReport() *Report {
	r := Report{
		ready:   &atomic.Bool{},
		lock:    &sync.Mutex{},
		channel: make(chan *string, bufferSize),
	}
	return &r
}

func (r *Report) Ready() bool {
	return r.ready.Load()
}

func (r *Report) Channel() chan *string {
	return r.channel
}

func (r *Report) Start(ctx context.Context) error {
	log.Info("Begin output report")
	go func() {
		defer r.ready.Store(false)
		for {
			select {
			case <-ctx.Done():
				log.Warn("Output report stopped")
				return
			case s := <-r.channel:
				if e := parseEvent(s); e != nil {
					r.lock.Lock()
					r.events = append(r.events, e)
					if len(r.events) > reportSize {
						r.events = r.events[len(r.events)-reportSize:]
					}
					r.lock.Unlock()
					log.Infof("Report event: %+v", *e)
				}
			}
		}
	}()
	r.ready.Store(true)
	return nil
}

func (r *Report) GetEvents() ([]*Event, error) {
	r.lock.Lock()
	events := make([]*Event, len(r.events))
	copy(events, r.events)
	r.lock.Unlock()
	return events, nil
}

func parseEvent(s *string) *Event {
	for i := range rawEvents {
		if res := rawEvents[i].Pattern.FindStringSubmatch(*s); res != nil {
			var args []any
			res = res[1:]
			for i := range res {
				args = append(args, res[i])
			}
			e := Event{
				Level: rawEvents[i].Level,
				Time:  time.Now().Unix(),
				Msg:   fmt.Sprintf(rawEvents[i].Format, args...),
				Type:  rawEvents[i].Type,
			}
			return &e
		}
	}
	return nil
}

func getEvents() []*Event {
	events := []*Event{
		{
			// Finding Mods...
			PatternString: `Finding Mods`,
			Format:        "正在加载 Mod",
			Level:         "info",
		},
		{
			// Creating world
			PatternString: `Creating world`,
			Format:        "正在生成世界",
			Level:         "info",
		},
		{
			// Loading world data: 1%
			PatternString: `Loading world data: 1%`,
			Format:        "正在加载世界",
			Level:         "info",
		},
		{
			// Server started
			PatternString: `Server started`,
			Format:        "服务器启动成功",
			Level:         "warning",
			Type:          TypeServerActive,
		},
	}
	for i := range events {
		events[i].Pattern = regexp.MustCompile(events[i].PatternString)
	}
	return events
}
