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
package ranking

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zhenghaoz/gorse/base"
	"github.com/zhenghaoz/gorse/base/task"
	"github.com/zhenghaoz/gorse/model"
	"io"
	"testing"
)

type mockMatrixFactorizationForSearch struct {
	model.BaseModel
}

func newMockMatrixFactorizationForSearch(numEpoch int) *mockMatrixFactorizationForSearch {
	return &mockMatrixFactorizationForSearch{model.BaseModel{Params: model.Params{model.NEpochs: numEpoch}}}
}

func (m *mockMatrixFactorizationForSearch) Complexity() int {
	panic("implement me")
}

func (m *mockMatrixFactorizationForSearch) Bytes() int {
	panic("implement me")
}

func (m *mockMatrixFactorizationForSearch) GetUserFactor(_ int32) []float32 {
	panic("implement me")
}

func (m *mockMatrixFactorizationForSearch) GetItemFactor(_ int32) []float32 {
	panic("implement me")
}

func (m *mockMatrixFactorizationForSearch) IsUserPredictable(_ int32) bool {
	panic("implement me")
}

func (m *mockMatrixFactorizationForSearch) IsItemPredictable(_ int32) bool {
	panic("implement me")
}

func (m *mockMatrixFactorizationForSearch) Marshal(_ io.Writer) error {
	panic("implement me")
}

func (m *mockMatrixFactorizationForSearch) Unmarshal(_ io.Reader) error {
	panic("implement me")
}

func (m *mockMatrixFactorizationForSearch) Invalid() bool {
	panic("implement me")
}

func (m *mockMatrixFactorizationForSearch) GetUserIndex() base.Index {
	panic("don't call me")
}

func (m *mockMatrixFactorizationForSearch) GetItemIndex() base.Index {
	panic("don't call me")
}

func (m *mockMatrixFactorizationForSearch) Fit(_, _ *DataSet, cfg *FitConfig) Score {
	score := float32(0)
	score += m.Params.GetFloat32(model.NFactors, 0.0)
	score += m.Params.GetFloat32(model.InitMean, 0.0)
	score += m.Params.GetFloat32(model.InitStdDev, 0.0)
	cfg.Task.Add(m.Params.GetInt(model.NEpochs, 0))
	return Score{NDCG: score}
}

func (m *mockMatrixFactorizationForSearch) Predict(_, _ string) float32 {
	panic("don't call me")
}

func (m *mockMatrixFactorizationForSearch) InternalPredict(_, _ int32) float32 {
	panic("don't call me")
}

func (m *mockMatrixFactorizationForSearch) Clear() {
	// do nothing
}

func (m *mockMatrixFactorizationForSearch) GetParamsGrid(_ bool) model.ParamsGrid {
	return model.ParamsGrid{
		model.NFactors:   []interface{}{1, 2, 3, 4},
		model.InitMean:   []interface{}{4, 3, 2, 1},
		model.InitStdDev: []interface{}{4, 4, 4, 4},
	}
}

type mockRunner struct {
	mock.Mock
}

func (r *mockRunner) Lock() {
	r.Called()
}

func (r *mockRunner) UnLock() {
	r.Called()
}

func newFitConfigForSearch() *FitConfig {
	t := task.NewTask("test", 100)
	return &FitConfig{
		Verbose: 1,
		Task:    t,
	}
}

func TestGridSearchCV(t *testing.T) {
	m := &mockMatrixFactorizationForSearch{}
	fitConfig := newFitConfigForSearch()
	runner := new(mockRunner)
	runner.On("Lock")
	runner.On("UnLock")
	r := GridSearchCV(m, nil, nil, m.GetParamsGrid(false), 0, fitConfig, runner)
	runner.AssertCalled(t, "Lock")
	runner.AssertCalled(t, "UnLock")
	assert.Equal(t, float32(12), r.BestScore.NDCG)
	assert.Equal(t, model.Params{
		model.NFactors:   4,
		model.InitMean:   4,
		model.InitStdDev: 4,
	}, r.BestParams)
}

func TestRandomSearchCV(t *testing.T) {
	m := &mockMatrixFactorizationForSearch{}
	fitConfig := newFitConfigForSearch()
	runner := new(mockRunner)
	runner.On("Lock")
	runner.On("UnLock")
	r := RandomSearchCV(m, nil, nil, m.GetParamsGrid(false), 63, 0, fitConfig, runner)
	runner.AssertCalled(t, "Lock")
	runner.AssertCalled(t, "UnLock")
	assert.Equal(t, float32(12), r.BestScore.NDCG)
	assert.Equal(t, model.Params{
		model.NFactors:   4,
		model.InitMean:   4,
		model.InitStdDev: 4,
	}, r.BestParams)
}

func TestModelSearcher(t *testing.T) {
	runner := new(mockRunner)
	runner.On("Lock")
	runner.On("UnLock")
	searcher := NewModelSearcher(2, 63, task.NewConstantJobsAllocator(1), false)
	searcher.models = []MatrixFactorization{newMockMatrixFactorizationForSearch(2)}
	tk := task.NewTask("test", searcher.Complexity())
	err := searcher.Fit(NewMapIndexDataset(), NewMapIndexDataset(), tk, runner)
	assert.NoError(t, err)
	_, m, score := searcher.GetBestModel()
	assert.Equal(t, float32(12), score.NDCG)
	assert.Equal(t, model.Params{
		model.NEpochs:    2,
		model.NFactors:   4,
		model.InitMean:   4,
		model.InitStdDev: 4,
	}, m.GetParams())
	assert.Equal(t, searcher.Complexity(), tk.Done)
}
