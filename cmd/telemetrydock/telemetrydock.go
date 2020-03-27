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

	"github.com/sodafoundation/telemetry/pkg/db"
	"github.com/sodafoundation/telemetry/pkg/dock"
	"github.com/sodafoundation/telemetry/pkg/model"
	. "github.com/sodafoundation/telemetry/pkg/utils/config"
	"github.com/sodafoundation/telemetry/pkg/utils/constants"
	"github.com/sodafoundation/telemetry/pkg/utils/daemon"
	"github.com/sodafoundation/telemetry/pkg/utils/logs"
)

func init() {
	// Load global configuration from specified config file.
	CONF.Load()

	// Parse some configuration fields from command line. and it will override the value which is got from config file.
	flag.StringVar(&CONF.TelemetryDock.ApiEndpoint, "api-endpoint", CONF.TelemetryDock.ApiEndpoint, "Listen endpoint of dock service")
	flag.StringVar(&CONF.TelemetryDock.DockType, "dock-type", CONF.TelemetryDock.DockType, "Type of dock service")
	flag.BoolVar(&CONF.TelemetryDock.Daemon, "daemon", CONF.TelemetryDock.Daemon, "Run app as a daemon with -daemon=true")
	flag.DurationVar(&CONF.TelemetryDock.LogFlushFrequency, "log-flush-frequency", CONF.TelemetryDock.LogFlushFrequency, "Maximum number of seconds between log flushes")
	flag.Parse()

	daemon.CheckAndRunDaemon(CONF.TelemetryDock.Daemon)
}

func main() {
	// Open OpenSDS dock service log file.
	logs.InitLogs(CONF.TelemetryDock.LogFlushFrequency)
	defer logs.FlushLogs()

	// Set up database session.
	db.Init(&CONF.Database)

	// FixMe: TelemetryDock attacher service needs to specify the endpoint via configuration file,
	//  so add this temporarily.
	listenEndpoint := constants.OpensdsDockBindEndpoint
	if CONF.TelemetryDock.DockType == model.DockTypeAttacher {
		listenEndpoint = CONF.TelemetryDock.ApiEndpoint
	}
	// Construct dock module grpc server struct and run dock server process.
	ds := dock.NewDockServer(CONF.TelemetryDock.DockType, listenEndpoint)
	if err := ds.Run(); err != nil {
		panic(err)
	}
}
