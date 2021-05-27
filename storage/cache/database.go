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

package cache

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"strings"
)

const (
	// IgnoreItems is these items that a user has read.
	IgnoreItems        = "ignore_items"
	PopularItems       = "popular_items"
	LatestItems        = "latest_items"
	SimilarItems       = "similar_items"
	CollaborativeItems = "collaborative_items"
	SubscribeItems     = "subscribe_items"

	GlobalMeta                  = "global_meta"
	CollectPopularTime          = "last_update_popular_time"
	CollectLatestTime           = "last_update_latest_time"
	CollectSimilarTime          = "last_update_similar_time"
	FitMatrixFactorizationTime  = "last_fit_match_model_time"
	FitFactorizationMachineTime = "last_fit_rank_model_time"
	MatrixFactorizationVersion  = "latest_match_model_version"
	FactorizationMachineVersion = "latest_rank_model_version"

	LastActiveTime          = "last_active_time"
	LastUpdateRecommendTime = "last_update_recommend_time"
	NumUsers                = "num_users"
	NumItems                = "num_items"
	NumPositiveFeedback     = "num_pos_feedback"
)

var ErrObjectNotExist = fmt.Errorf("object not exists")
var ErrNoDatabase = fmt.Errorf("no database specified")

type ScoredItem struct {
	ItemId string
	Score  float32
}

func CreateScoredItems(itemIds []string, scores []float32) []ScoredItem {
	if len(itemIds) != len(scores) {
		panic("the length of itemIds and scores should be equal")
	}
	items := make([]ScoredItem, len(itemIds))
	for i := range items {
		items[i].ItemId = itemIds[i]
		items[i].Score = scores[i]
	}
	return items
}

func RemoveScores(items []ScoredItem) []string {
	itemIds := make([]string, len(items))
	for i := range itemIds {
		itemIds[i] = items[i].ItemId
	}
	return itemIds
}

type Database interface {
	Close() error
	SetScores(prefix, name string, items []ScoredItem) error
	GetScores(prefix, name string, begin int, end int) ([]ScoredItem, error)
	ClearList(prefix, name string) error
	AppendList(prefix, name string, items ...string) error
	GetList(prefix, name string) ([]string, error)
	GetString(prefix, name string) (string, error)
	SetString(prefix, name string, val string) error
	GetInt(prefix, name string) (int, error)
	SetInt(prefix, name string, val int) error
}

const redisPrefix = "redis://"

// Open a connection to a database.
func Open(path string) (Database, error) {
	if strings.HasPrefix(path, redisPrefix) {
		addr := path[len(redisPrefix):]
		database := new(Redis)
		database.client = redis.NewClient(&redis.Options{Addr: addr})
		return database, nil
	}
	return nil, errors.Errorf("Unknown database: %s", path)
}
