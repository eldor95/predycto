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
	"github.com/chewxy/math32"
	"github.com/scylladb/go-set"
	"github.com/scylladb/go-set/strset"
	"github.com/zhenghaoz/gorse/base"
	"github.com/zhenghaoz/gorse/model"
	"github.com/zhenghaoz/gorse/model/click"
	"github.com/zhenghaoz/gorse/model/ranking"
	"github.com/zhenghaoz/gorse/storage/cache"
	"github.com/zhenghaoz/gorse/storage/data"
	"go.uber.org/zap"
	"sort"
	"time"
)

const (
	RankingTop10NDCG      = "NDCG@10"
	RankingTop10Precision = "Precision@10"
	RankingTop10Recall    = "Recall@10"
	ClickPrecision        = "Precision"
	ClickThroughRate      = "ClickThroughRate"
	ActiveUsersYesterday  = "ActiveUsersYesterday"
	ActiveUsersMonthly    = "ActiveUsersMonthly"
)

// popItem updates popular items for the database.
func (m *Master) popItem(items []data.Item, feedback []data.Feedback) {
	base.Logger().Info("collect popular items", zap.Int("n_cache", m.GorseConfig.Database.CacheSize))
	// create item mapping
	itemMap := make(map[string]data.Item)
	for _, item := range items {
		itemMap[item.ItemId] = item
	}
	// count feedback
	timeWindowLimit := time.Now().AddDate(0, 0, -m.GorseConfig.Recommend.PopularWindow)
	count := make(map[string]int)
	for _, fb := range feedback {
		if fb.Timestamp.After(timeWindowLimit) {
			count[fb.ItemId]++
		}
	}
	// collect pop items
	popItems := make(map[string]*base.TopKStringFilter)
	popItems[""] = base.NewTopKStringFilter(m.GorseConfig.Database.CacheSize)
	for itemId, f := range count {
		popItems[""].Push(itemId, float32(f))
		item := itemMap[itemId]
		for _, label := range item.Labels {
			if _, exists := popItems[label]; !exists {
				popItems[label] = base.NewTopKStringFilter(m.GorseConfig.Database.CacheSize)
			}
			popItems[label].Push(itemId, float32(f))
		}
	}
	// write back
	for label, topItems := range popItems {
		result, scores := topItems.PopAll()
		if err := m.CacheClient.SetScores(cache.PopularItems, label, cache.CreateScoredItems(result, scores)); err != nil {
			base.Logger().Error("failed to cache popular items", zap.Error(err))
		}
	}
	if err := m.CacheClient.SetString(cache.GlobalMeta, cache.LastUpdatePopularTime, base.Now()); err != nil {
		base.Logger().Error("failed to cache popular items", zap.Error(err))
	}
}

// latest updates latest items.
func (m *Master) latest(items []data.Item) {
	base.Logger().Info("collect latest items", zap.Int("n_cache", m.GorseConfig.Database.CacheSize))
	var err error
	latestItems := make(map[string]*base.TopKStringFilter)
	latestItems[""] = base.NewTopKStringFilter(m.GorseConfig.Database.CacheSize)
	// find latest items
	for _, item := range items {
		if !item.Timestamp.IsZero() {
			latestItems[""].Push(item.ItemId, float32(item.Timestamp.Unix()))
			for _, label := range item.Labels {
				if _, exist := latestItems[label]; !exist {
					latestItems[label] = base.NewTopKStringFilter(m.GorseConfig.Database.CacheSize)
				}
				latestItems[label].Push(item.ItemId, float32(item.Timestamp.Unix()))
			}
		}
	}
	for label, topItems := range latestItems {
		result, scores := topItems.PopAll()
		if err = m.CacheClient.SetScores(cache.LatestItems, label, cache.CreateScoredItems(result, scores)); err != nil {
			base.Logger().Error("failed to cache latest items", zap.Error(err))
		}
	}
	if err = m.CacheClient.SetString(cache.GlobalMeta, cache.LastUpdateLatestTime, base.Now()); err != nil {
		base.Logger().Error("failed to cache latest items time", zap.Error(err))
	}
}

// similar updates neighbors for the database.
func (m *Master) similar(items []data.Item, dataset *ranking.DataSet, similarity string) {
	base.Logger().Info("collect similar items", zap.Int("n_cache", m.GorseConfig.Database.CacheSize))
	// create progress tracker
	completed := make(chan []interface{}, 1000)
	go func() {
		completedCount := 0
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case _, ok := <-completed:
				if !ok {
					return
				}
				completedCount++
			case <-ticker.C:
				base.Logger().Debug("collect similar items",
					zap.Int("n_complete_items", completedCount),
					zap.Int("n_items", dataset.ItemCount()))
			}
		}
	}()

	// pre-ranking
	for _, feedbacks := range dataset.ItemFeedback {
		sort.Ints(feedbacks)
	}

	if err := base.Parallel(dataset.ItemCount(), m.GorseConfig.Master.FitJobs, func(workerId, jobId int) error {
		users := dataset.ItemFeedback[jobId]
		// Collect candidates
		itemSet := set.NewIntSet()
		for _, u := range users {
			itemSet.Add(dataset.UserFeedback[u]...)
		}
		// Ranking
		nearItems := base.NewTopKFilter(m.GorseConfig.Database.CacheSize)
		for j := range itemSet.List() {
			if j != jobId {
				var score float32
				score = dotInt(dataset.ItemFeedback[jobId], dataset.ItemFeedback[j])
				if similarity == model.SimilarityCosine {
					score /= math32.Sqrt(float32(len(dataset.ItemFeedback[jobId])))
					score /= math32.Sqrt(float32(len(dataset.ItemFeedback[j])))
				}
				nearItems.Push(j, score)
			}
		}
		elem, scores := nearItems.PopAll()
		recommends := make([]string, len(elem))
		for i := range recommends {
			recommends[i] = dataset.ItemIndex.ToName(elem[i])
		}
		if err := m.CacheClient.SetScores(cache.SimilarItems, dataset.ItemIndex.ToName(jobId), cache.CreateScoredItems(recommends, scores)); err != nil {
			return err
		}
		completed <- nil
		return nil
	}); err != nil {
		base.Logger().Error("failed to cache similar items", zap.Error(err))
	}
	close(completed)
	if err := m.CacheClient.SetString(cache.GlobalMeta, cache.LastUpdateNeighborTime, base.Now()); err != nil {
		base.Logger().Error("failed to cache similar items", zap.Error(err))
	}
}

func dotString(a, b []string) float32 {
	i, j, sum := 0, 0, float32(0)
	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			i++
			j++
			sum++
		} else if a[i] < b[j] {
			i++
		} else if a[i] > b[j] {
			j++
		}
	}
	return sum
}

func dotInt(a, b []int) float32 {
	i, j, sum := 0, 0, float32(0)
	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			i++
			j++
			sum++
		} else if a[i] < b[j] {
			i++
		} else if a[i] > b[j] {
			j++
		}
	}
	return sum
}

// fitRankingModel fits ranking model using passed dataset. After model fitted, following states are changed:
// 1. Ranking model version are increased.
// 2. Ranking model score are updated.
// 3. Ranking model, version and score are persisted to local cache.
func (m *Master) fitRankingModel(dataSet *ranking.DataSet, rankingModel ranking.Model) {
	base.Logger().Info("fit ranking model", zap.Int("n_jobs", m.GorseConfig.Master.FitJobs))
	// training model
	trainSet, testSet := dataSet.Split(0, 0)
	score := rankingModel.Fit(trainSet, testSet, ranking.NewFitConfig().SetJobs(m.GorseConfig.Master.FitJobs))
	// update match model
	m.rankingModelMutex.Lock()
	m.rankingModel = rankingModel
	m.rankingModelVersion++
	m.rankingScore = score
	m.rankingModelMutex.Unlock()
	base.Logger().Info("fit ranking model complete",
		zap.String("version", fmt.Sprintf("%x", m.rankingModelVersion)))
	if err := m.DataClient.InsertMeasurement(data.Measurement{Name: RankingTop10NDCG, Value: score.NDCG, Timestamp: time.Now()}); err != nil {
		base.Logger().Error("failed to insert measurement", zap.Error(err))
	}
	if err := m.DataClient.InsertMeasurement(data.Measurement{Name: RankingTop10Recall, Value: score.Recall, Timestamp: time.Now()}); err != nil {
		base.Logger().Error("failed to insert measurement", zap.Error(err))
	}
	if err := m.DataClient.InsertMeasurement(data.Measurement{Name: RankingTop10Precision, Value: score.Precision, Timestamp: time.Now()}); err != nil {
		base.Logger().Error("failed to insert measurement", zap.Error(err))
	}
	if err := m.CacheClient.SetString(cache.GlobalMeta, cache.LastFitRankingModelTime, base.Now()); err != nil {
		base.Logger().Error("failed to write meta", zap.Error(err))
	}
	if err := m.CacheClient.SetString(cache.GlobalMeta, cache.LastRankingModelVersion, fmt.Sprintf("%x", m.rankingModelVersion)); err != nil {
		base.Logger().Error("failed to write meta", zap.Error(err))
	}
	// caching model
	m.localCache.RankingModelName = m.rankingModelName
	m.localCache.RankingModelVersion = m.rankingModelVersion
	m.localCache.RankingModel = rankingModel
	m.localCache.RankingModelScore = score
	m.localCache.UserIndex = m.userIndex
	if m.localCache.ClickModel == nil || m.localCache.ClickModel.Invalid() {
		base.Logger().Info("wait click model")
	} else if err := m.localCache.WriteLocalCache(); err != nil {
		base.Logger().Error("failed to write local cache", zap.Error(err))
	} else {
		base.Logger().Info("write model to local cache",
			zap.String("ranking_model_name", m.localCache.RankingModelName),
			zap.String("ranking_model_version", base.Hex(m.localCache.RankingModelVersion)),
			zap.Float32("ranking_model_score", m.localCache.RankingModelScore.NDCG),
			zap.Any("ranking_model_params", m.localCache.RankingModel.GetParams()))
	}
}

func (m *Master) analyze() error {
	// pull existed click through rates
	clickThroughRates, err := m.DataClient.GetMeasurements(ClickThroughRate, 30)
	if err != nil {
		return err
	}
	existed := strset.New()
	for _, clickThroughRate := range clickThroughRates {
		existed.Add(clickThroughRate.Timestamp.String())
	}
	// update click through rate
	for i := 30; i > 0; i-- {
		dateTime := time.Now().AddDate(0, 0, -i)
		date := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 0, 0, 0, 0, time.UTC)
		if !existed.Has(date.String()) {
			// click through clickThroughRate
			startTime := time.Now()
			clickThroughRate, err := m.DataClient.GetClickThroughRate(date, m.GorseConfig.Database.ClickFeedbackTypes, m.GorseConfig.Database.ReadFeedbackType)
			if err != nil {
				return err
			}
			err = m.DataClient.InsertMeasurement(data.Measurement{
				Name:      ClickThroughRate,
				Timestamp: date,
				Value:     float32(clickThroughRate),
			})
			if err != nil {
				return err
			}
			base.Logger().Info("update click through rate",
				zap.String("date", date.String()),
				zap.Duration("time_used", time.Since(startTime)),
				zap.Float64("click_through_rate", clickThroughRate))
		}
	}

	// pull existed active users
	yesterdayActiveUsers, err := m.DataClient.GetMeasurements(ActiveUsersYesterday, 1)
	if err != nil {
		return err
	}
	yesterdayDatetime := time.Now().AddDate(0, 0, -1)
	yesterdayDate := time.Date(yesterdayDatetime.Year(), yesterdayDatetime.Month(), yesterdayDatetime.Day(), 0, 0, 0, 0, time.UTC)
	if len(yesterdayActiveUsers) == 0 || !yesterdayActiveUsers[0].Timestamp.Equal(yesterdayDate) {
		startTime := time.Now()
		activeUsers, err := m.DataClient.CountActiveUsers(time.Now().AddDate(0, 0, -1))
		if err != nil {
			return err
		}
		err = m.DataClient.InsertMeasurement(data.Measurement{
			Name:      ActiveUsersYesterday,
			Timestamp: yesterdayDate,
			Value:     float32(activeUsers),
		})
		if err != nil {
			return err
		}
		base.Logger().Info("update active users of yesterday",
			zap.String("date", yesterdayDate.String()),
			zap.Duration("time_used", time.Since(startTime)),
			zap.Int("active_users", activeUsers))
	}

	monthActiveUsers, err := m.DataClient.GetMeasurements(ActiveUsersMonthly, 1)
	if err != nil {
		return err
	}
	if len(monthActiveUsers) == 0 || !monthActiveUsers[0].Timestamp.Equal(yesterdayDate) {
		startTime := time.Now()
		activeUsers, err := m.DataClient.CountActiveUsers(time.Now().AddDate(0, 0, -30))
		if err != nil {
			return err
		}
		err = m.DataClient.InsertMeasurement(data.Measurement{
			Name:      ActiveUsersMonthly,
			Timestamp: yesterdayDate,
			Value:     float32(activeUsers),
		})
		if err != nil {
			return err
		}
		base.Logger().Info("update active users of this month",
			zap.String("date", yesterdayDate.String()),
			zap.Duration("time_used", time.Since(startTime)),
			zap.Int("active_users", activeUsers))
	}

	return nil
}

// fitClickModel fits click model using latest data. After model fitted, following states are changed:
// 1. Click model version are increased.
// 2. Click model score are updated.
// 3. Click model, version and score are persisted to local cache.
func (m *Master) fitClickModel(fm click.FactorizationMachine) {
	base.Logger().Info("fit click model", zap.Int("n_jobs", m.GorseConfig.Master.FitJobs))

	// pull dataset
	dataset, err := click.LoadDataFromDatabase(m.DataClient,
		m.GorseConfig.Database.PositiveFeedbackType,
		m.GorseConfig.Database.ReadFeedbackType)
	if err != nil {
		base.Logger().Error("failed to pull dataset", zap.Error(err))
		return
	}

	// training model
	trainSet, testSet := dataset.Split(0.2, 0)
	score := fm.Fit(trainSet, testSet, click.NewFitConfig().SetJobs(m.GorseConfig.Master.FitJobs))
	// update match model
	m.clickModelMutex.Lock()
	m.clickModel = fm
	m.clickScore = score
	m.clickModelVersion++
	m.clickModelMutex.Unlock()
	base.Logger().Info("fit click model complete",
		zap.String("version", fmt.Sprintf("%x", m.clickModelVersion)))
	if err := m.DataClient.InsertMeasurement(data.Measurement{Name: ClickPrecision, Value: score.Precision, Timestamp: time.Now()}); err != nil {
		base.Logger().Error("failed to insert measurement", zap.Error(err))
	}

	// caching model
	m.localCache.ClickModelScore = m.clickScore
	m.localCache.ClickModelVersion = m.clickModelVersion
	m.localCache.ClickModel = fm
	if m.localCache.RankingModel == nil || m.localCache.RankingModel.Invalid() {
		base.Logger().Info("wait ranking model")
	} else if err = m.localCache.WriteLocalCache(); err != nil {
		base.Logger().Error("failed to write local cache", zap.Error(err))
	} else {
		base.Logger().Info("write model to local cache",
			zap.String("click_model_version", base.Hex(m.localCache.ClickModelVersion)),
			zap.Float32("click_model_score", score.Precision),
			zap.Any("click_model_params", m.localCache.ClickModel.GetParams()))
	}
}

// searchClickModel searches best hyper-parameters for factorization machines.
func (m *Master) searchClickModel() error {
	dataset, err := click.LoadDataFromDatabase(m.DataClient,
		m.GorseConfig.Database.ClickFeedbackTypes,
		m.GorseConfig.Database.ReadFeedbackType)
	if err != nil {
		base.Logger().Error("failed to pull data from click model search", zap.Error(err))
		return err
	}
	trainSet, testSet := dataset.Split(0.2, 0)
	return m.clickModelSearcher.Fit(trainSet, testSet)
}
