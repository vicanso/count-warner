// Copyright 2019 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package warner

import (
	"sync"
	"time"
)

type (
	// Count count
	Count struct {
		CreatedAt int64
		Value     int
	}
	// warner warner
	warner struct {
		Duration    time.Duration
		Max         int
		ResetOnWarn bool
		mu          sync.Mutex
		m           map[string]*Count
		listerns    []Listener
	}
	// Listener warner listener
	Listener func(key string, c Count)
)

// NewWarner create a new warner
func NewWarner(duration time.Duration, max int) *warner {
	return &warner{
		m:        make(map[string]*Count),
		listerns: make([]Listener, 0),
		Duration: duration,
		Max:      max,
	}
}

// Inc increase or decrease count value
func (w *warner) Inc(key string, value int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	c, ok := w.m[key]
	now := time.Now().UnixNano()
	if !ok || c.CreatedAt+w.Duration.Nanoseconds() < now {
		c = &Count{
			CreatedAt: now,
		}
		w.m[key] = c
	}
	c.Value += value
	if c.Value > w.Max {
		for _, fn := range w.listerns {
			fn(key, *c)
		}
		if w.ResetOnWarn {
			delete(w.m, key)
		}
	}
}

// Reset reset the count value
func (w *warner) Reset(key string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.m, key)
}

// ClearExpired clear expired count
func (w *warner) ClearExpired() {
	w.mu.Lock()
	defer w.mu.Unlock()
	m := make(map[string]*Count)
	now := time.Now().UnixNano()
	d := w.Duration.Nanoseconds()
	for k, c := range w.m {
		// 只保留未过期的
		if c.CreatedAt+d > now {
			m[k] = c
		}
	}
	// 更换map
	w.m = m
}

// On on warn event
func (w *warner) On(ln Listener) {
	w.listerns = append(w.listerns, ln)
}
