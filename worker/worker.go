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

package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/juju/errors"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/scylladb/go-set"
	"github.com/scylladb/go-set/strset"
	"github.com/zhenghaoz/gorse/base"
	"github.com/zhenghaoz/gorse/config"
	"github.com/zhenghaoz/gorse/model/click"
	"github.com/zhenghaoz/gorse/model/ranking"
	"github.com/zhenghaoz/gorse/protocol"
	"github.com/zhenghaoz/gorse/storage/cache"
	"github.com/zhenghaoz/gorse/storage/data"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"math"
	"math/rand"
	"net/http"
	"time"
)

const batchSize = 10000

// Worker manages states of a worker node.
type Worker struct {
	// worker config
	cfg        *config.Config
	jobs       int
	workerName string
	httpHost   string
	httpPort   int
	masterHost string
	masterPort int
	testMode   bool
	cacheFile  string

	// database connection
	cachePath   string
	cacheClient cache.Database
	dataPath    string
	dataClient  data.Database

	// master connection
	masterClient protocol.MasterClient

	// user index
	latestUserIndexVersion  int64
	currentUserIndexVersion int64
	userIndex               base.Index

	// ranking model
	latestRankingModelVersion  int64
	currentRankingModelVersion int64
	rankingModel               ranking.MatrixFactorization

	// click model
	latestClickModelVersion  int64
	currentClickModelVersion int64
	clickModel               click.FactorizationMachine

	// peers
	peers []string
	me    string

	// events
	ticker     *time.Ticker
	syncedChan chan bool // meta synced events
	pulledChan chan bool // model pulled events
}

// NewWorker creates a new worker node.
func NewWorker(masterHost string, masterPort int, httpHost string, httpPort, jobs int, cacheFile string) *Worker {
	return &Worker{
		// database
		dataClient:  data.NoDatabase{},
		cacheClient: cache.NoDatabase{},
		// config
		cacheFile:  cacheFile,
		masterHost: masterHost,
		masterPort: masterPort,
		httpHost:   httpHost,
		httpPort:   httpPort,
		jobs:       jobs,
		cfg:        (*config.Config)(nil).LoadDefaultIfNil(),
		// events
		ticker:     time.NewTicker(time.Minute),
		syncedChan: make(chan bool, 1024),
		pulledChan: make(chan bool, 1024),
	}
}

// Sync this worker to the master.
func (w *Worker) Sync() {
	defer base.CheckPanic()
	base.Logger().Info("start meta sync", zap.Int("meta_timeout", w.cfg.Master.MetaTimeout))
	for {
		var meta *protocol.Meta
		var err error
		if meta, err = w.masterClient.GetMeta(context.Background(),
			&protocol.NodeInfo{
				NodeType: protocol.NodeType_WorkerNode,
				NodeName: w.workerName,
				HttpPort: int64(w.httpPort),
			}); err != nil {
			base.Logger().Error("failed to get meta", zap.Error(err))
			goto sleep
		}

		// load master config
		err = json.Unmarshal([]byte(meta.Config), &w.cfg)
		if err != nil {
			base.Logger().Error("failed to parse master config", zap.Error(err))
			goto sleep
		}

		// connect to data store
		if w.dataPath != w.cfg.Database.DataStore {
			base.Logger().Info("connect data store", zap.String("database", w.cfg.Database.DataStore))
			if w.dataClient, err = data.Open(w.cfg.Database.DataStore); err != nil {
				base.Logger().Error("failed to connect data store", zap.Error(err))
				goto sleep
			}
			w.dataPath = w.cfg.Database.DataStore
		}

		// connect to cache store
		if w.cachePath != w.cfg.Database.CacheStore {
			base.Logger().Info("connect cache store", zap.String("database", w.cfg.Database.CacheStore))
			if w.cacheClient, err = cache.Open(w.cfg.Database.CacheStore); err != nil {
				base.Logger().Error("failed to connect cache store", zap.Error(err))
				goto sleep
			}
			w.cachePath = w.cfg.Database.CacheStore
		}

		// check ranking model version
		w.latestRankingModelVersion = meta.RankingModelVersion
		if w.latestRankingModelVersion != w.currentRankingModelVersion {
			base.Logger().Info("new ranking model found",
				zap.String("old_version", base.Hex(w.currentRankingModelVersion)),
				zap.String("new_version", base.Hex(w.latestRankingModelVersion)))
			w.syncedChan <- true
		}

		// check click model version
		w.latestClickModelVersion = meta.ClickModelVersion
		if w.latestClickModelVersion != w.currentClickModelVersion {
			base.Logger().Info("new click model found",
				zap.String("old_version", base.Hex(w.currentClickModelVersion)),
				zap.String("new_version", base.Hex(w.latestClickModelVersion)))
			w.syncedChan <- true
		}

		// check user index version
		w.latestUserIndexVersion = meta.UserIndexVersion
		if w.latestUserIndexVersion != w.currentUserIndexVersion {
			base.Logger().Info("new user index found",
				zap.String("old_version", base.Hex(w.currentUserIndexVersion)),
				zap.String("new_version", base.Hex(w.latestUserIndexVersion)))
			w.syncedChan <- true
		}

		w.peers = meta.Workers
		w.me = meta.Me
	sleep:
		if w.testMode {
			return
		}
		time.Sleep(time.Duration(w.cfg.Master.MetaTimeout) * time.Second)
	}
}

// Pull user index and ranking model from master.
func (w *Worker) Pull() {
	defer base.CheckPanic()
	for range w.syncedChan {
		pulled := false

		// pull user index
		if w.latestUserIndexVersion != w.currentUserIndexVersion {
			base.Logger().Info("start pull user index")
			if userIndexReceiver, err := w.masterClient.GetUserIndex(context.Background(),
				&protocol.VersionInfo{Version: w.latestUserIndexVersion},
				grpc.MaxCallRecvMsgSize(math.MaxInt)); err != nil {
				base.Logger().Error("failed to pull user index", zap.Error(err))
			} else {
				// encode user index
				var userIndex base.Index
				userIndex, err = protocol.UnmarshalIndex(userIndexReceiver)
				if err != nil {
					base.Logger().Error("fail to unmarshal user index", zap.Error(err))
				} else {
					w.userIndex = userIndex
					w.currentUserIndexVersion = w.latestUserIndexVersion
					base.Logger().Info("synced user index",
						zap.String("version", base.Hex(w.currentUserIndexVersion)))
					pulled = true
				}
			}
		}

		// pull ranking model
		if w.latestRankingModelVersion != w.currentRankingModelVersion {
			base.Logger().Info("start pull ranking model")
			if rankingModelReceiver, err := w.masterClient.GetRankingModel(context.Background(),
				&protocol.VersionInfo{Version: w.latestRankingModelVersion},
				grpc.MaxCallRecvMsgSize(math.MaxInt)); err != nil {
				base.Logger().Error("failed to pull ranking model", zap.Error(err))
			} else {
				var rankingModel ranking.MatrixFactorization
				rankingModel, err = protocol.UnmarshalRankingModel(rankingModelReceiver)
				if err != nil {
					base.Logger().Error("failed to unmarshal ranking model", zap.Error(err))
				} else {
					w.rankingModel = rankingModel
					w.currentRankingModelVersion = w.latestRankingModelVersion
					base.Logger().Info("synced ranking model",
						zap.String("version", base.Hex(w.currentRankingModelVersion)))
					pulled = true
				}
			}
		}

		// pull click model
		if w.latestClickModelVersion != w.currentClickModelVersion {
			base.Logger().Info("start pull click model")
			if clickModelReceiver, err := w.masterClient.GetClickModel(context.Background(),
				&protocol.VersionInfo{Version: w.latestClickModelVersion},
				grpc.MaxCallRecvMsgSize(math.MaxInt)); err != nil {
				base.Logger().Error("failed to pull click model", zap.Error(err))
			} else {
				var clickModel click.FactorizationMachine
				clickModel, err = protocol.UnmarshalClickModel(clickModelReceiver)
				if err != nil {
					base.Logger().Error("failed to unmarshal click model", zap.Error(err))
				} else {
					w.clickModel = clickModel
					w.currentClickModelVersion = w.latestClickModelVersion
					base.Logger().Info("synced click model",
						zap.String("version", base.Hex(w.currentClickModelVersion)))
					pulled = true
				}
			}
		}

		if w.testMode {
			return
		}
		if pulled {
			w.pulledChan <- true
		}
	}
}

// ServeMetrics serves Prometheus metrics.
func (w *Worker) ServeMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", w.httpHost, w.httpPort), nil)
	if err != nil {
		base.Logger().Fatal("failed to start http server", zap.Error(err))
	}
}

// Serve as a worker node.
func (w *Worker) Serve() {
	rand.Seed(time.Now().UTC().UnixNano())
	// open local store
	state, err := LoadLocalCache(w.cacheFile)
	if err != nil {
		base.Logger().Error("failed to load persist state", zap.Error(err),
			zap.String("path", state.path))
	}
	if state.WorkerName == "" {
		state.WorkerName = base.GetRandomName(0)
		err = state.WriteLocalCache()
		if err != nil {
			base.Logger().Fatal("failed to write meta", zap.Error(err))
		}
	}
	w.workerName = state.WorkerName
	base.Logger().Info("start worker",
		zap.Int("n_jobs", w.jobs),
		zap.String("worker_name", w.workerName))

	// connect to master
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", w.masterHost, w.masterPort), grpc.WithInsecure())
	if err != nil {
		base.Logger().Fatal("failed to connect master", zap.Error(err))
	}
	w.masterClient = protocol.NewMasterClient(conn)

	go w.Sync()
	go w.Pull()
	go w.ServeMetrics()

	loop := func() {
		if w.userIndex == nil {
			base.Logger().Debug("user index doesn't exist")
		} else {
			// split users
			workingUsers, err := split(w.userIndex, w.peers, w.me)
			if err != nil {
				base.Logger().Error("failed to split users", zap.Error(err),
					zap.String("me", w.me),
					zap.Strings("workers", w.peers))
				return
			}

			// recommendation
			w.Recommend(workingUsers)
		}
	}

	for {
		select {
		case <-w.ticker.C:
			loop()
		case <-w.pulledChan:
			loop()
		}
	}
}

// Recommend items to users. The workflow of recommendation is:
// 1. Skip inactive users.
// 2. Load historical items.
// 3. Load positive items if KNN used.
// 4. Generate recommendation.
// 5. Save result.
// 6. Insert cold-start items into results.
// 7. Rank items in results by click-through-rate.
// 8. Refresh cache.
func (w *Worker) Recommend(users []string) {
	// load user index
	userIndexer := w.userIndex
	base.Logger().Info("ranking recommendation",
		zap.Int("n_working_users", len(users)),
		zap.Int("n_jobs", w.jobs),
		zap.Int("cache_size", w.cfg.Database.CacheSize))
	// progress tracker
	completed := make(chan interface{}, 1000)
	taskName := fmt.Sprintf("Generate offline recommendation [%s]", w.workerName)
	if w.masterClient != nil {
		if _, err := w.masterClient.StartTask(context.Background(),
			&protocol.StartTaskRequest{Name: taskName, Total: int64(len(users))}); err != nil {
			base.Logger().Error("failed to report start task", zap.Error(err))
		}
	}
	// pull items from database
	itemCache, itemCategories, err := w.pullItems()
	if err != nil {
		base.Logger().Error("failed to pull items", zap.Error(err))
		return
	}
	go func() {
		defer base.CheckPanic()
		completedCount := 0
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case _, ok := <-completed:
				if !ok {
					return
				}
				completedCount++
			case <-ticker.C:
				if w.masterClient != nil {
					if _, err := w.masterClient.UpdateTask(context.Background(),
						&protocol.UpdateTaskRequest{Name: taskName, Done: int64(completedCount)}); err != nil {
						base.Logger().Error("failed to report update task", zap.Error(err))
					}
				}
				base.Logger().Info("ranking recommendation",
					zap.Int("n_complete_users", completedCount),
					zap.Int("n_working_users", len(users)))
			}
		}
	}()
	// recommendation
	startTime := time.Now()
	userFeedbackCache := NewFeedbackCache(w.dataClient, w.cfg.Database.PositiveFeedbackType...)
	err = base.Parallel(len(users), w.jobs, func(workerId, jobId int) error {
		userStartTime := time.Now()
		userId := users[jobId]
		// convert to user index
		userIndex := userIndexer.ToNumber(userId)
		// skip inactive users before max recommend period
		if !w.checkRecommendCacheTimeout(userId, itemCategories) {
			return nil
		}

		// load historical items
		historyItems, err := loadUserHistoricalItems(w.dataClient, userId)
		excludeSet := set.NewStringSet(historyItems...)
		if err != nil {
			base.Logger().Error("failed to pull user feedback",
				zap.String("user_id", userId), zap.Error(err))
			return errors.Trace(err)
		}

		// load positive items
		var positiveItems []string
		if w.cfg.Recommend.EnableItemBasedRecommend {
			positiveItems, err = userFeedbackCache.GetUserFeedback(userId)
			if err != nil {
				base.Logger().Error("failed to pull user feedback",
					zap.String("user_id", userId), zap.Error(err))
				return errors.Trace(err)
			}
		}

		// create candidates container
		candidates := make(map[string][][]string)
		candidates[""] = make([][]string, 0)
		for _, category := range itemCategories {
			candidates[category] = make([][]string, 0)
		}

		// Recommender #1: collaborative filtering.
		if w.cfg.Recommend.EnableColRecommend && w.rankingModel != nil && w.rankingModel.IsUserPredictable(userIndex) {
			itemIds := w.rankingModel.GetItemIndex().GetNames()
			localStartTime := time.Now()
			recItemsFilters := make(map[string]*base.TopKStringFilter)
			recItemsFilters[""] = base.NewTopKStringFilter(w.cfg.Database.CacheSize)
			for _, category := range itemCategories {
				recItemsFilters[category] = base.NewTopKStringFilter(w.cfg.Database.CacheSize)
			}
			for itemIndex, itemId := range itemIds {
				if !excludeSet.Has(itemId) && itemCache.IsAvailable(itemId) && w.rankingModel.IsItemPredictable(int32(itemIndex)) {
					prediction := w.rankingModel.InternalPredict(userIndex, int32(itemIndex))
					recItemsFilters[""].Push(itemId, prediction)
					for _, category := range itemCache[itemId].Categories {
						recItemsFilters[category].Push(itemId, prediction)
					}
				}
			}
			// save result
			for category, recItemsFilter := range recItemsFilters {
				recommendItems, recommendScores := recItemsFilter.PopAll()
				candidates[category] = append(candidates[category], recommendItems)
				if err = w.cacheClient.SetCategoryScores(cache.CollaborativeRecommend, userId, category, cache.CreateScoredItems(recommendItems, recommendScores)); err != nil {
					base.Logger().Error("failed to cache collaborative filtering recommendation result", zap.Error(err))
					return errors.Trace(err)
				}
			}
			CollaborativeRecommendSeconds.Observe(time.Since(localStartTime).Seconds())
		} else if w.rankingModel == nil {
			base.Logger().Warn("no collaborative filtering model")
		} else if !w.rankingModel.IsUserPredictable(userIndex) {
			base.Logger().Warn("user is unpredictable", zap.String("user_id", userId))
		}

		// Recommender #2: item-based.
		if w.cfg.Recommend.EnableItemBasedRecommend {
			localStartTime := time.Now()
			for _, category := range append([]string{""}, itemCategories...) {
				// collect candidates
				scores := make(map[string]float32)
				for _, itemId := range positiveItems {
					// load similar items
					similarItems, err := w.cacheClient.GetCategoryScores(cache.ItemNeighbors, itemId, category, 0, w.cfg.Database.CacheSize)
					if err != nil {
						base.Logger().Error("failed to load similar items", zap.Error(err))
						return errors.Trace(err)
					}
					// add unseen items
					for _, item := range similarItems {
						if !excludeSet.Has(item.Id) && itemCache.IsAvailable(item.Id) {
							scores[item.Id] += item.Score
						}
					}
				}
				// collect top k
				filter := base.NewTopKStringFilter(w.cfg.Database.CacheSize)
				for id, score := range scores {
					filter.Push(id, score)
				}
				ids, _ := filter.PopAll()
				candidates[category] = append(candidates[category], ids)
			}
			ItemBasedRecommendSeconds.Observe(time.Since(localStartTime).Seconds())
		}

		// Recommender #3: insert user-based items
		if w.cfg.Recommend.EnableUserBasedRecommend {
			localStartTime := time.Now()
			scores := make(map[string]float32)
			// load similar users
			similarUsers, err := w.cacheClient.GetScores(cache.UserNeighbors, userId, 0, w.cfg.Database.CacheSize)
			if err != nil {
				base.Logger().Error("failed to load similar users", zap.Error(err))
				return errors.Trace(err)
			}
			for _, user := range similarUsers {
				// load historical feedback
				similarUserPositiveItems, err := userFeedbackCache.GetUserFeedback(user.Id)
				if err != nil {
					base.Logger().Error("failed to pull user feedback",
						zap.String("user_id", userId), zap.Error(err))
					return errors.Trace(err)
				}
				// add unseen items
				for _, itemId := range similarUserPositiveItems {
					if !excludeSet.Has(itemId) && itemCache.IsAvailable(itemId) {
						scores[itemId] += user.Score
					}
				}
			}
			// collect top k
			filters := make(map[string]*base.TopKStringFilter)
			filters[""] = base.NewTopKStringFilter(w.cfg.Database.CacheSize)
			for _, category := range itemCategories {
				filters[category] = base.NewTopKStringFilter(w.cfg.Database.CacheSize)
			}
			for id, score := range scores {
				filters[""].Push(id, score)
				for _, category := range itemCache[id].Categories {
					filters[category].Push(id, score)
				}
			}
			for category, filter := range filters {
				ids, _ := filter.PopAll()
				candidates[category] = append(candidates[category], ids)
			}
			UserBasedRecommendSeconds.Observe(time.Since(localStartTime).Seconds())
		}

		// Recommender #4: latest items.
		if w.cfg.Recommend.EnableLatestRecommend {
			localStartTime := time.Now()
			for _, category := range append([]string{""}, itemCategories...) {
				latestItems, err := w.cacheClient.GetSort(cache.Key(cache.LatestItems, category), 0, w.cfg.Database.CacheSize)
				if err != nil {
					base.Logger().Error("failed to load latest items", zap.Error(err))
					return errors.Trace(err)
				}
				var recommend []string
				for _, latestItem := range latestItems {
					if !excludeSet.Has(latestItem.Id) && itemCache.IsAvailable(latestItem.Id) {
						recommend = append(recommend, latestItem.Id)
					}
				}
				candidates[category] = append(candidates[category], recommend)
			}
			LoadLatestRecommendCacheSeconds.Observe(time.Since(localStartTime).Seconds())
		}

		// Recommender #5: popular items.
		if w.cfg.Recommend.EnablePopularRecommend {
			localStartTime := time.Now()
			for _, category := range append([]string{""}, itemCategories...) {
				popularItems, err := w.cacheClient.GetSort(cache.Key(cache.PopularItems, category), 0, w.cfg.Database.CacheSize)
				if err != nil {
					base.Logger().Error("failed to load popular items", zap.Error(err))
					return errors.Trace(err)
				}
				var recommend []string
				for _, popularItem := range popularItems {
					if !excludeSet.Has(popularItem.Id) && itemCache.IsAvailable(popularItem.Id) {
						recommend = append(recommend, popularItem.Id)
					}
				}
				candidates[category] = append(candidates[category], recommend)
			}
			LoadPopularRecommendCacheSeconds.Observe(time.Since(localStartTime).Seconds())
		}

		// rank items
		results := make(map[string][]cache.Scored)
		for category, catCandidates := range candidates {
			if w.cfg.Recommend.EnableClickThroughPrediction && w.clickModel != nil {
				results[category], err = w.rankByClickTroughRate(userId, catCandidates, itemCache)
				if err != nil {
					base.Logger().Error("failed to rank items", zap.Error(err))
					return errors.Trace(err)
				}
			} else {
				results[category] = mergeAndShuffle(catCandidates)
			}
		}

		// explore latest and popular
		for category, result := range results {
			results[category], err = w.exploreRecommend(result, excludeSet, category)
			if err != nil {
				base.Logger().Error("failed to explore latest and popular items", zap.Error(err))
				return errors.Trace(err)
			}

			if err = w.cacheClient.SetCategoryScores(cache.OfflineRecommend, userId, category, results[category]); err != nil {
				base.Logger().Error("failed to cache recommendation", zap.Error(err))
				return errors.Trace(err)
			}
		}
		if err = w.cacheClient.SetTime(cache.LastUpdateUserRecommendTime, userId, time.Now()); err != nil {
			base.Logger().Error("failed to cache recommendation time", zap.Error(err))
		}

		// refresh cache
		err = w.refreshCache(userId)
		if err != nil {
			base.Logger().Error("failed to refresh cache", zap.Error(err))
			return errors.Trace(err)
		}
		completed <- nil
		GenerateRecommendSeconds.Observe(time.Since(userStartTime).Seconds())
		return nil
	})
	close(completed)
	if err != nil {
		base.Logger().Error("failed to continue offline recommendation", zap.Error(err))
		return
	}
	if w.masterClient != nil {
		if _, err := w.masterClient.FinishTask(context.Background(),
			&protocol.FinishTaskRequest{Name: taskName}); err != nil {
			base.Logger().Error("failed to report finish task", zap.Error(err))
		}
	}
	base.Logger().Info("complete ranking recommendation",
		zap.String("used_time", time.Since(startTime).String()))
}

// rankByClickTroughRate ranks items by predicted click-through-rate.
func (w *Worker) rankByClickTroughRate(userId string, candidates [][]string, itemCache map[string]data.Item) ([]cache.Scored, error) {
	startTime := time.Now()
	var err error
	// concat candidates
	memo := strset.New()
	var itemIds []string
	for _, v := range candidates {
		for _, itemId := range v {
			if !memo.Has(itemId) {
				memo.Add(itemId)
				itemIds = append(itemIds, itemId)
			}
		}
	}
	// download user
	var user data.User
	user, err = w.dataClient.GetUser(userId)
	if err != nil {
		return nil, errors.Trace(err)
	}
	// download items
	items := make([]data.Item, 0, len(itemIds))
	for _, itemId := range itemIds {
		if item, exist := itemCache[itemId]; exist {
			items = append(items, item)
		} else {
			base.Logger().Warn("item doesn't exists in database", zap.String("item_id", itemId))
		}
	}
	// rank by CTR
	topItems := base.NewTopKStringFilter(w.cfg.Database.CacheSize)
	for _, item := range items {
		topItems.Push(item.ItemId, w.clickModel.Predict(userId, item.ItemId, user.Labels, item.Labels))
	}
	elems, scores := topItems.PopAll()
	CTRRecommendSeconds.Observe(time.Since(startTime).Seconds())
	return cache.CreateScoredItems(elems, scores), nil
}

func mergeAndShuffle(candidates [][]string) []cache.Scored {
	memo := strset.New()
	pos := make([]int, len(candidates))
	var recommend []cache.Scored
	for {
		// filter out ended slice
		var src []int
		for i := range candidates {
			if pos[i] < len(candidates[i]) {
				src = append(src, i)
			}
		}
		if len(src) == 0 {
			break
		}
		// select a slice randomly
		j := src[rand.Intn(len(src))]
		candidateId := candidates[j][pos[j]]
		pos[j]++
		if !memo.Has(candidateId) {
			memo.Add(candidateId)
			recommend = append(recommend, cache.Scored{Score: 0, Id: candidateId})
		}
	}
	return recommend
}

func (w *Worker) exploreRecommend(exploitRecommend []cache.Scored, excludeSet *strset.Set, category string) ([]cache.Scored, error) {
	localExcludeSet := excludeSet.Copy()
	// create thresholds
	explorePopularThreshold := 0.0
	if threshold, exist := w.cfg.Recommend.GetExploreRecommend("popular"); exist {
		explorePopularThreshold = threshold
	}
	exploreLatestThreshold := explorePopularThreshold
	if threshold, exist := w.cfg.Recommend.GetExploreRecommend("latest"); exist {
		exploreLatestThreshold += threshold
	}
	// load popular items
	popularItems, err := w.cacheClient.GetSort(cache.Key(cache.PopularItems, category), 0, w.cfg.Database.CacheSize)
	if err != nil {
		return nil, errors.Trace(err)
	}
	// load the latest items
	latestItems, err := w.cacheClient.GetSort(cache.Key(cache.LatestItems, category), 0, w.cfg.Database.CacheSize)
	if err != nil {
		return nil, errors.Trace(err)
	}
	// explore recommendation
	var exploreRecommend []cache.Scored
	for range exploitRecommend {
		dice := rand.Float64()
		var recommendItem cache.Scored
		if dice < explorePopularThreshold && len(popularItems) > 0 {
			recommendItem = popularItems[0]
			popularItems = popularItems[1:]
		} else if dice < exploreLatestThreshold && len(latestItems) > 0 {
			recommendItem = latestItems[0]
			latestItems = latestItems[1:]
		} else if len(exploitRecommend) > 0 {
			recommendItem = exploitRecommend[0]
			exploitRecommend = exploitRecommend[1:]
		} else {
			break
		}
		if !localExcludeSet.Has(recommendItem.Id) {
			localExcludeSet.Add(recommendItem.Id)
			exploreRecommend = append(exploreRecommend, recommendItem)
		}
	}
	return exploreRecommend, nil
}

// checkRecommendCacheTimeout checks if recommend cache stale.
// 1. if cache is empty, stale.
// 2. if active time > recommend time, stale.
// 3. if recommend time + timeout < now, stale.
func (w *Worker) checkRecommendCacheTimeout(userId string, categories []string) bool {
	var activeTime, recommendTime time.Time
	// check cache
	for _, category := range append([]string{""}, categories...) {
		items, err := w.cacheClient.GetCategoryScores(cache.OfflineRecommend, userId, category, 0, -1)
		if err != nil {
			base.Logger().Error("failed to read meta", zap.String("user_id", userId), zap.Error(err))
			return true
		} else if len(items) == 0 {
			return true
		}
	}
	// read active time
	var err error
	activeTime, err = w.cacheClient.GetTime(cache.LastModifyUserTime, userId)
	if err != nil {
		base.Logger().Error("failed to read meta", zap.Error(err))
		return true
	}
	// read recommend time
	recommendTime, err = w.cacheClient.GetTime(cache.LastUpdateUserRecommendTime, userId)
	if err != nil {
		base.Logger().Error("failed to read meta", zap.Error(err))
		return true
	}
	// check time
	if activeTime.Unix() < recommendTime.Unix() {
		timeoutTime := recommendTime.Add(time.Hour * 24 * time.Duration(w.cfg.Recommend.RefreshRecommendPeriod))
		return timeoutTime.Unix() < time.Now().Unix()
	}
	return true
}

func loadUserHistoricalItems(database data.Database, userId string) ([]string, error) {
	items := make([]string, 0)
	feedbacks, err := database.GetUserFeedback(userId, false)
	if err != nil {
		return nil, err
	}
	for _, feedback := range feedbacks {
		items = append(items, feedback.ItemId)
	}
	return items, nil
}

func (w *Worker) refreshCache(userId string) error {
	var timeLimit *time.Time
	// read recommend time
	recommendTime, err := w.cacheClient.GetTime(cache.LastUpdateUserRecommendTime, userId)
	if err == nil {
		timeLimit = &recommendTime
	} else {
		return errors.Trace(err)
	}
	// clear cache
	err = w.cacheClient.ClearScores(cache.IgnoreItems, userId)
	if err != nil {
		return errors.Trace(err)
	}
	// load cache
	feedback, err := w.dataClient.GetUserFeedback(userId, true)
	if err != nil {
		return errors.Trace(err)
	}
	var items []cache.Scored
	for _, v := range feedback {
		if v.Timestamp.Unix() > timeLimit.Unix() {
			items = append(items, cache.Scored{Id: v.ItemId, Score: float32(v.Timestamp.Unix())})
		}
	}
	err = w.cacheClient.AppendScores(cache.IgnoreItems, userId, items...)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

// split users between worker nodes.
func split(userIndex base.Index, nodes []string, me string) ([]string, error) {
	// locate me
	pos := -1
	for i, node := range nodes {
		if node == me {
			pos = i
		}
	}
	if pos == -1 {
		return nil, fmt.Errorf("current node isn't in worker nodes")
	}
	// split users
	users := userIndex.GetNames()
	workingUsers := make([]string, 0)
	for ; pos < len(users); pos += len(nodes) {
		workingUsers = append(workingUsers, users[pos])
	}
	base.Logger().Info("allocate working users",
		zap.Int("n_working_users", len(workingUsers)),
		zap.Int("n_users", len(users)))
	return workingUsers, nil
}

func (w *Worker) pullItems() (ItemCache, []string, error) {
	// pull items from database
	itemCache := make(ItemCache)
	itemCategories := strset.New()
	itemChan, errChan := w.dataClient.GetItemStream(batchSize, nil)
	for batchItems := range itemChan {
		for _, item := range batchItems {
			itemCache[item.ItemId] = item
			itemCategories.Add(item.Categories...)
		}
	}
	if err := <-errChan; err != nil {
		return nil, nil, errors.Trace(err)
	}
	return itemCache, itemCategories.List(), nil
}

// ItemCache is alias of map[string]data.Item.
type ItemCache map[string]data.Item

// IsAvailable means the item exists in database and is not hidden.
func (c ItemCache) IsAvailable(itemId string) bool {
	if item, exist := c[itemId]; exist {
		return !item.IsHidden
	} else {
		return false
	}
}

// FeedbackCache is the cache for user feedbacks.
type FeedbackCache struct {
	Types  []string
	Cache  cmap.ConcurrentMap
	Client data.Database
}

// NewFeedbackCache creates a new FeedbackCache.
func NewFeedbackCache(client data.Database, feedbackTypes ...string) *FeedbackCache {
	return &FeedbackCache{
		Types:  feedbackTypes,
		Client: client,
		Cache:  cmap.New(),
	}
}

// GetUserFeedback gets user feedback from cache or database.
func (c *FeedbackCache) GetUserFeedback(userId string) ([]string, error) {
	if tmp, ok := c.Cache.Get(userId); ok {
		return tmp.([]string), nil
	} else {
		items := make([]string, 0)
		feedbacks, err := c.Client.GetUserFeedback(userId, false, c.Types...)
		if err != nil {
			return nil, err
		}
		for _, feedback := range feedbacks {
			items = append(items, feedback.ItemId)
		}
		c.Cache.Set(userId, items)
		return items, nil
	}
}
