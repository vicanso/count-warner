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
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWarner(t *testing.T) {
	assert := assert.New(t)
	max := 10
	warner := NewWarner(time.Millisecond, max)
	warner.ResetOnWarn = true
	key := "abcd"
	done := false
	warner.On(func(k string, count int) {
		done = true
		assert.Equal(max+1, count)
	})
	for index := 0; index < max+1; index++ {
		warner.Inc(key, 1)
	}
	assert.True(done)
}

func TestWarnerClearExpired(t *testing.T) {
	assert := assert.New(t)
	warner := NewWarner(time.Millisecond, 10)
	warner.Inc("abcd", 1)
	assert.NotNil(warner.get("abcd"))

	time.Sleep(2 * time.Millisecond)
	warner.ClearExpired()
	assert.Nil(warner.get("abcd"))
}

func TestParallel(t *testing.T) {
	warner := NewWarner(time.Second, 10)
	go func() {
		for i := 0; i < 100; i++ {
			warner.Inc("1", 1)
		}
	}()
	go func() {
		for i := 0; i < 100; i++ {
			warner.Inc(strconv.Itoa(i), 1)
		}
	}()
	go func() {
		for i := 0; i < 100; i++ {
			warner.ClearExpired()
		}
	}()
	time.Sleep(100 * time.Millisecond)
}
