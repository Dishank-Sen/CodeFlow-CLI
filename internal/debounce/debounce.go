package debounce

import (
	"sync"
	"time"
)

type Debouncer struct {
	mu     sync.Mutex
	timers map[string]*time.Timer
}

func NewDebouncer() *Debouncer {
	return &Debouncer{
		timers: make(map[string]*time.Timer),
	}
}

func (d *Debouncer) Debounce(key string, delay time.Duration, fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if timer, ok := d.timers[key]; ok {
		timer.Stop()
	}

	d.timers[key] = time.AfterFunc(delay, func() {
		d.mu.Lock()
		delete(d.timers, key)
		d.mu.Unlock()

		fn()
	})
}

func GetDebounceTime() int64{
	// reads the config file for debounce time
	
}