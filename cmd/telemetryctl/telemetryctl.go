// Copyright 2019 The OpenSDS Authors.
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

	c "github.com/sodafoundation/telemetry/pkg/controller"
	"github.com/sodafoundation/telemetry/pkg/db"
	. "github.com/sodafoundation/telemetry/pkg/utils/config"
	"github.com/sodafoundation/telemetry/pkg/utils/constants"
	"github.com/sodafoundation/telemetry/pkg/utils/daemon"
	"github.com/sodafoundation/telemetry/pkg/utils/logs"
)

func init() {
	// Load global configuration from specified config file.
	CONF.Load()

	// Parse some configuration fields from command line. and it will override the value which is got from config file.
	flag.StringVar(&CONF.TelemetryCtl.ApiEndpoint, "api-endpoint", CONF.TelemetryCtl.ApiEndpoint, "Listen endpoint of controller service")
	flag.BoolVar(&CONF.TelemetryCtl.Daemon, "daemon", CONF.TelemetryCtl.Daemon, "Run app as a daemon with -daemon=true")
	flag.DurationVar(&CONF.TelemetryCtl.LogFlushFrequency, "log-flush-frequency", CONF.TelemetryCtl.LogFlushFrequency, "Maximum number of seconds between log flushes")

	flag.StringVar(&CONF.TelemetryCtl.PrometheusPushMechanism, "prometheus-push-mechanism", CONF.TelemetryCtl.PrometheusPushMechanism, "Prometheus push mechanism")
	flag.StringVar(&CONF.TelemetryCtl.PushGatewayUrl, "prometheus-push-gateway-url", CONF.TelemetryCtl.PushGatewayUrl, "Prometheus push gateway URL")
	flag.StringVar(&CONF.TelemetryCtl.NodeExporterWatchFolder, "node-exporter-watch-folder", CONF.TelemetryCtl.NodeExporterWatchFolder, "Node exporter watch folder")
	flag.StringVar(&CONF.TelemetryCtl.KafkaEndpoint, "kafka-endpoint", CONF.TelemetryCtl.KafkaEndpoint, "Kafka endpoint")
	flag.StringVar(&CONF.TelemetryCtl.KafkaTopic, "kafka-topic", CONF.TelemetryCtl.KafkaTopic, "Kafka topic")
	flag.StringVar(&CONF.TelemetryCtl.GrafanaUrl, "grafana-url", CONF.TelemetryCtl.GrafanaUrl, "Grafana listen endpoint")
	flag.StringVar(&CONF.TelemetryCtl.AlertMgrUrl, "alertmgr-url", CONF.TelemetryCtl.AlertMgrUrl, "Alert manager listen endpoint")
	flag.Parse()

	daemon.CheckAndRunDaemon(CONF.TelemetryCtl.Daemon)
}

func main() {
	// Open OpenSDS orchestrator service log file.
	logs.InitLogs(CONF.TelemetryCtl.LogFlushFrequency)
	defer logs.FlushLogs()

	// Set up database session.
	db.Init(&CONF.Database)

	// Construct controller module grpc server struct and run controller server process.
	if err := c.NewGrpcServer(constants.OpensdsCtrBindEndpoint).Run(); err != nil {
		panic(err)
	}
}
