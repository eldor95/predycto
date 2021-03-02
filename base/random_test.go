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
package base

import (
	"github.com/chewxy/math32"
	"github.com/stretchr/testify/assert"
	"github.com/zhenghaoz/gorse/floats"
	"testing"
)

const randomEpsilon = 0.1

func TestRandomGenerator_MakeNormalMatrix(t *testing.T) {
	rng := NewRandomGenerator(0)
	vec := rng.NormalMatrix(1, 1000, 1, 2)[0]
	assert.False(t, math32.Abs(floats.Mean(vec)-1) > randomEpsilon)
	assert.False(t, math32.Abs(floats.StdDev(vec)-2) > randomEpsilon)
}

func TestRandomGenerator_MakeUniformMatrix(t *testing.T) {
	rng := NewRandomGenerator(0)
	vec := rng.UniformMatrix(1, 1000, 1, 2)[0]
	assert.False(t, floats.Min(vec) < 1)
	assert.False(t, floats.Max(vec) > 2)
}
