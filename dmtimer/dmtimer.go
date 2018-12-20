package dmtimer

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type DmTimers struct {
	list map[string]*time.Timer
	mu   sync.Mutex
}

func (d *DmTimers) Get(key string) *time.Timer {
	return d.list[key]
}

func (d *DmTimers) Add(key string, value *time.Timer) {
	d.mu.Lock()
	if d.list == nil {
		d.list = make(map[string]*time.Timer)
	}
	d.list[key] = value
	d.mu.Unlock()
}

func (d *DmTimers) Del(key string, value *time.Timer) {
	d.mu.Lock()
	delete(d.list, key)
	d.mu.Unlock()
}

func (d *DmTimers) Keys() []string {
	keys := []string{}
	for k := range d.list {
		keys = append(keys, k)
		fmt.Println(k)
	}

	return keys
}

func (d *DmTimers) Len() int {
	return len(d.list)
}
func ParseTimerID(url string) (string, error) {
	// expeted format foo.com/ping/123456asdf we want the part afer /ping/
	urlParts := strings.Split(url, "/")
	if len(urlParts) <= 2 {
		return "", fmt.Errorf("Url did not include a Timer ID")
	}

	// get rid of any stray spaces
	timerID := strings.TrimSpace(urlParts[2])

	return timerID, nil
}
