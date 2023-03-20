// Copyright 2020 gorse Project Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package master

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhenghaoz/gorse/base/task"
	"github.com/zhenghaoz/gorse/config"
	"github.com/zhenghaoz/gorse/storage/cache"
	"github.com/zhenghaoz/gorse/storage/data"
)

type mockMaster struct {
	Master
}

func (m *mockMaster) Close() {
}

func newMockMaster(t *testing.T) *mockMaster {
	s := new(mockMaster)
	s.taskMonitor = task.NewTaskMonitor()
	// open database
	var err error
	s.Settings = config.NewSettings()
	s.DataClient, err = data.Open(fmt.Sprintf("sqlite://%s/data.db", t.TempDir()), "")
	assert.NoError(t, err)
	s.CacheClient, err = cache.Open(fmt.Sprintf("sqlite://%s/cache.db", t.TempDir()), "")
	assert.NoError(t, err)
	// init database
	err = s.DataClient.Init()
	assert.NoError(t, err)
	err = s.CacheClient.Init()
	assert.NoError(t, err)
	return s
}
