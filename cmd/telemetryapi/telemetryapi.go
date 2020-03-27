// Copyright 2018 The OpenSDS Authors.
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

/*
This module implements a entry into the OpenSDS REST service.

*/

package main

import (
	"flag"

	"github.com/sodafoundation/telemetry/pkg/api"
	"github.com/sodafoundation/telemetry/pkg/db"
	. "github.com/sodafoundation/telemetry/pkg/utils/config"
	"github.com/sodafoundation/telemetry/pkg/utils/daemon"
	"github.com/sodafoundation/telemetry/pkg/utils/logs"
)

func init() {
	// Load global configuration from specified config file.
	CONF.Load()

	// Parse some configuration fields from command line. and it will override the value which is got from config file.
	flag.StringVar(&CONF.TelemetryApi.ApiEndpoint, "api-endpoint", CONF.TelemetryApi.ApiEndpoint, "Listen endpoint of api-server service")
	flag.DurationVar(&CONF.TelemetryApi.LogFlushFrequency, "log-flush-frequency", CONF.TelemetryApi.LogFlushFrequency, "Maximum number of seconds between log flushes")
	flag.BoolVar(&CONF.TelemetryApi.Daemon, "daemon", CONF.TelemetryApi.Daemon, "Run app as a daemon with -daemon=true")
	// prometheus related
	flag.StringVar(&CONF.TelemetryApi.PrometheusConfHome, "prometheus-conf-home", CONF.TelemetryApi.PrometheusConfHome, "Prometheus conf. path")
	flag.StringVar(&CONF.TelemetryApi.PrometheusUrl, "prometheus-url", CONF.TelemetryApi.PrometheusUrl, "Prometheus URL")
	flag.StringVar(&CONF.TelemetryApi.PrometheusConfFile, "prometheus-conf-file", CONF.TelemetryApi.PrometheusConfFile, "Prometheus conf. file")
	// alert manager related
	flag.StringVar(&CONF.TelemetryApi.AlertmgrConfHome, "alertmgr-conf-home", CONF.TelemetryApi.AlertmgrConfHome, "Alert manager conf. home")
	flag.StringVar(&CONF.TelemetryApi.AlertMgrUrl, "alertmgr-url", CONF.TelemetryApi.AlertMgrUrl, "Alert manager listen endpoint")
	flag.StringVar(&CONF.TelemetryApi.AlertmgrConfFile, "alertmgr-conf-file", CONF.TelemetryApi.AlertmgrConfFile, "Alert manager conf. file")
	// grafana related
	flag.StringVar(&CONF.TelemetryApi.GrafanaConfHome, "grafana-conf-home", CONF.TelemetryApi.GrafanaConfHome, "Grafana conf. home")
	flag.StringVar(&CONF.TelemetryApi.GrafanaRestartCmd, "grafana-restart-cmd", CONF.TelemetryApi.GrafanaRestartCmd, "Grafana restart command")
	flag.StringVar(&CONF.TelemetryApi.GrafanaConfFile, "grafana-conf-file", CONF.TelemetryApi.GrafanaConfFile, "Grafana conf file")
	flag.StringVar(&CONF.TelemetryApi.GrafanaUrl, "grafana-url", CONF.TelemetryApi.GrafanaUrl, "Grafana listen endpoint")
	// prometheus and alert manager configuration reload url
	flag.StringVar(&CONF.TelemetryApi.ConfReloadUrl, "conf-reload-url", CONF.TelemetryApi.ConfReloadUrl, "Prometheus and Alert manager conf. reload URL")
	flag.Parse()

	daemon.CheckAndRunDaemon(CONF.TelemetryApi.Daemon)
}

func main() {
	// Open OpenSDS orchestrator service log file.
	logs.InitLogs(CONF.TelemetryApi.LogFlushFrequency)
	defer logs.FlushLogs()

	// Set up database session.
	db.Init(&CONF.Database)

	// Start OpenSDS northbound REST service.
	api.Run(CONF.TelemetryApi)
}
