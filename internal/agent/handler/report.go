package handler

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"sync/atomic"
	"terraria-run/internal/common/model"
	"time"
)

type Report struct {
	ready   *atomic.Bool
	lock    *sync.Mutex
	channel chan *string
	events  []*model.Event
}

type Event struct {
	PatternString string
	Pattern       *regexp.Regexp
	Format        string
	Level         string
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

func (r *Report) Run(ctx context.Context) error {
	r.ready.Store(true)
	defer r.ready.Store(false)
	events := getEvents()
	for i := range events {
		events[i].Pattern = regexp.MustCompile(events[i].PatternString)
	}
	log.Info("Begin output report")
	for {
		select {
		case <-ctx.Done():
			log.Info("Stop output report")
			return nil
		case s := <-r.channel:
			for i := range events {
				if res := events[i].Pattern.FindStringSubmatch(*s); res != nil {
					var args []any
					res = res[1:]
					for i := range res {
						args = append(args, res[i])
					}
					e := model.Event{
						Time:  time.Now().Unix(),
						Msg:   fmt.Sprintf(events[i].Format, args...),
						Level: events[i].Level,
					}
					r.lock.Lock()
					r.events = append(r.events, &e)
					if len(r.events) > reportSize {
						r.events = r.events[len(r.events)-reportSize:]
					}
					r.lock.Unlock()
					log.Infof("Create event: %+v", e)
				}
			}
		}
	}
}

func (r *Report) GetEvents() ([]*model.Event, error) {
	if !r.Ready() {
		return nil, fmt.Errorf("report is not ready, get events failed")
	}
	var events []*model.Event
	r.lock.Lock()
	copy(r.events, events)
	r.lock.Unlock()
	return events, nil
}

func getEvents() []*Event {
	return []*Event{
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
			// Server started
			PatternString: `Server started`,
			Format:        "服务器启动成功",
			Level:         "info",
		},
	}
}
