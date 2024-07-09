// Copyright 2021 gorse Project Authors
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

package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/zhenghaoz/gorse/config"
)

var (
	RestAPIRequestSecondsVec = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: config.TracingConfig.GetAppName(),
		Subsystem: "server",
		Name:      "rest_api_request_seconds",
	}, []string{"api"})
)
