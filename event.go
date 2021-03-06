package timewheel

import (
	"fmt"
	"sync"
	"time"
)

var globalEventPool = newEventPool()

// ExpireFunc represents a function will be executed when a event is trigged.
type ExpireFunc func()

// An Event represents an elemenet of the events in the timer.
type Event struct {
	slotPos int // mark timeWheel slot index
	index   int // index in the min heap structure

	ttl    time.Duration // wait delay time
	expire time.Time     // due timestamp
	fn     ExpireFunc    // callback function

	next    *Event
	cron    bool // repeat task
	cronNum int  // cron circle num
	alone   bool // indicates event is alone or in the free linked-list of timer
}

// clear field
func (e *Event) clear() {
	e.index = 0
	e.slotPos = 0
	e.cron = false
	e.fn = nil
	e.alone = false
}

// Less is used to compare expiration with other events.
func (e *Event) Less(o *Event) bool {
	return e.expire.Before(o.expire)
}

// Delay is used to give the duration that event will expire.
func (e *Event) Delay() time.Duration {
	return e.expire.Sub(time.Now())
}

func (e *Event) String() string {
	return fmt.Sprintf("index %d ttl %v, expire at %v", e.index, e.ttl, e.expire)
}

func newEventPool() *eventPool {
	return &eventPool{}
}

type eventPool struct {
	p sync.Pool
}

func (ep *eventPool) get() *Event {
	if t, _ := ep.p.Get().(*Event); t != nil {
		return t
	}

	return new(Event)
}

func (ep *eventPool) put(ev *Event) {
	ep.p.Put(ev)
}
