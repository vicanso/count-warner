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
	warner.On(func(k string, c Count) {
		done = true
		assert.Equal(max+1, c.Value)
	})
	for index := 0; index < max+1; index++ {
		warner.Inc(key, 1)
	}
	assert.True(done)
	assert.Nil(warner.m[key])
}

func TestWarnerReset(t *testing.T) {
	key := "abcd"
	assert := assert.New(t)
	warner := NewWarner(time.Second, 10)
	warner.Inc(key, 1)
	assert.Equal(1, len(warner.m))
	warner.Reset(key)
	assert.Equal(0, len(warner.m))
}

func TestWarnerClearExpired(t *testing.T) {
	assert := assert.New(t)
	warner := NewWarner(time.Millisecond, 10)
	warner.Inc("abcd", 1)
	assert.Equal(1, len(warner.m))
	time.Sleep(2 * time.Millisecond)
	warner.ClearExpired()
	assert.Equal(0, len(warner.m))
}
