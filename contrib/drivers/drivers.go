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
This module defines an standard table of storage driver. The default storage
driver is sample driver used for testing. If you want to use other storage
plugin, just modify Init() and Clean() method.

*/

package drivers

import (
	"github.com/sodafoundation/telemetry/contrib/drivers/ceph"
	"github.com/sodafoundation/telemetry/contrib/drivers/huawei/oceanstor"
	"github.com/sodafoundation/telemetry/contrib/drivers/lvm"
	
	"github.com/sodafoundation/telemetry/contrib/drivers/utils/config"
	"github.com/sodafoundation/telemetry/pkg/model"

)

// VolumeDriver is an interface for exposing some operations of different volume
// drivers, currently support sample, lvm, ceph, cinder and so forth.



func CleanMetricDriver(d MetricDriver) MetricDriver {
	// Execute different clean operations according to the MetricDriver type.
	switch d.(type) {
	case *lvm.MetricDriver:
		break
	default:
		break
	}
	_ = d.Teardown()
	d = nil

	return d
}

type MetricDriver interface {
	//Any initialization the metric driver does while starting.
	Setup() error
	//Any operation the metric driver does while stopping.
	Teardown() error
	// Collect metrics for all supported resources
	CollectMetrics() ([]*model.MetricSpec, error)
}

// Init
func InitMetricDriver(resourceType string) MetricDriver {
	var d MetricDriver
	switch resourceType {
	case config.LVMDriverType:
		d = &lvm.MetricDriver{}
		break
	case config.CephDriverType:
		d = &ceph.MetricDriver{}
		break
	case config.HuaweiOceanStorBlockDriverType:
		d = &oceanstor.MetricDriver{}
		break
	default:
		//d = &sample.Driver{}
		break
	}
	d.Setup()
	return d
}
