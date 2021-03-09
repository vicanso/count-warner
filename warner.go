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

	"go.uber.org/atomic"
)

type (
	// Count count
	Count struct {
		createdAt time.Time
		value     *atomic.Int32
		emitted   *atomic.Bool
	}
	// warner warner
	warner struct {
		Duration    time.Duration
		Max         int
		ResetOnWarn bool
		m           *sync.Map
		listerns    []Listener
	}
	// Listener warner listener
	Listener func(key string, count int)
)

// NewWarner returns a new warner
func NewWarner(duration time.Duration, max int) *warner {
	return &warner{
		m:        &sync.Map{},
		listerns: make([]Listener, 0),
		Duration: duration,
		Max:      max,
	}
}

// Inc value for the key, value can be < 0
func (w *warner) Inc(key string, value int) {
	// warner并不需要精准的判断，
	// 因此如果并发时可能会覆盖数据，
	// 对整体的告警影响不大
	c := w.get(key)
	now := time.Now()
	if c != nil && c.createdAt.Add(w.Duration).Before(now) {
		c = nil
	}
	if c == nil {
		c = &Count{
			createdAt: now,
			value:     atomic.NewInt32(0),
			emitted:   atomic.NewBool(false),
		}
		w.m.Store(key, c)
	}
	count := c.value.Add(int32(value))
	// 仅触发一次
	if count > int32(w.Max) && !c.emitted.Swap(true) {
		for _, fn := range w.listerns {
			fn(key, int(count))
		}
		// 如果设置了告警后重置
		if w.ResetOnWarn {
			c.value.Store(0)
			c.emitted.Store(false)
		}
	}
}

// On adds a listener for warner
func (w *warner) On(ln Listener) {
	w.listerns = append(w.listerns, ln)
}

// ClearExpired clears the expired key
func (w *warner) ClearExpired() {
	now := time.Now()
	w.m.Range(func(key, value interface{}) bool {
		c, ok := value.(*Count)
		if ok && c.createdAt.Add(w.Duration).Before(now) {
			w.m.Delete(key)
		}
		return true
	})
}

// get get count value
func (w *warner) get(key string) *Count {
	v, ok := w.m.Load(key)
	if !ok {
		return nil
	}
	c, _ := v.(*Count)
	return c
}
