// Copyright 2017 The OpenSDS Authors.
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
This module implements a entry into the OpenSDS northbound service.
*/

package controller

import (
	"context"
	log "github.com/golang/glog"
	osdsCtx "github.com/sodafoundation/telemetry/pkg/context"

	"github.com/sodafoundation/telemetry/pkg/controller/metrics"

	"github.com/sodafoundation/telemetry/pkg/db"
	"github.com/sodafoundation/telemetry/pkg/model"
	pb "github.com/sodafoundation/telemetry/pkg/model/proto"
)

const (
	CREATE_LIFECIRCLE_FLAG = iota + 1
	GET_LIFECIRCLE_FLAG
	LIST_LIFECIRCLE_FLAG
	DELETE_LIFECIRCLE_FLAG
	EXTEND_LIFECIRCLE_FLAG
)

func NewController() *Controller {

	metricsCtrl := metrics.NewController()
	return &Controller{

		metricsController:   metricsCtrl,

	}
}

type Controller struct {

	metricsController   metrics.Controller

}

// CreateVolume implements pb.ControllerServer.CreateVolume

func (c *Controller) GetMetrics(context context.Context, opt *pb.GetMetricsOpts) (*pb.GenericResponse, error) {
	log.Info("in controller get metrics methods")

	var result []*model.MetricSpec
	var err error

	if opt.StartTime == "" && opt.EndTime == "" {
		// no start and end time specified, get the latest value of this metric
		result, err = c.metricsController.GetLatestMetrics(opt)
	} else if opt.StartTime == opt.EndTime {
		// same start and end time specified, get the value of this metric at that timestamp
		result, err = c.metricsController.GetInstantMetrics(opt)
	} else {
		// range of start and end time is specified
		result, err = c.metricsController.GetRangeMetrics(opt)
	}

	if err != nil {
		log.Errorf("get metrics failed: %s\n", err.Error())

		return pb.GenericResponseError(err), err
	}

	return pb.GenericResponseResult(result), err
}

func (c *Controller) CollectMetrics(context context.Context, opt *pb.CollectMetricsOpts) (*pb.GenericResponse, error) {
	log.V(5).Info("in controller collect metrics methods")

	ctx := osdsCtx.NewContextFromJson(opt.GetContext())
	dockSpec, err := db.C.ListDocks(ctx)
	if err != nil {
		log.Errorf("list dock failed in CollectMetrics method: %s", err.Error())
		return pb.GenericResponseError(err), err
	}
	for i, d := range dockSpec {
		if d.DriverName == opt.DriverName {
			log.Infof("driver found driver: %s", d.DriverName)
			dockInfo, err := db.C.GetDock(ctx, dockSpec[i].BaseModel.Id)
			if err != nil {
				log.Errorf("error %s when search dock in db by dock id: %s", err.Error(), dockSpec[i].BaseModel.Id)
				return pb.GenericResponseError(err), err

			}
			c.metricsController.SetDock(dockInfo)
			result, err := c.metricsController.CollectMetrics(opt)
			if err != nil {
				log.Errorf("collectMetrics failed: %s", err.Error())

				return pb.GenericResponseError(err), err
			}

			return pb.GenericResponseResult(result), nil
		}
	}
	return nil, nil
}

func (c *Controller) GetUrls(context.Context, *pb.NoParams) (*pb.GenericResponse, error) {
	log.V(5).Info("in controller get urls method")

	var result *map[string]model.UrlDesc
	var err error

	result, err = c.metricsController.GetUrls()

	// make return array
	arrUrls := make([]model.UrlSpec, 0)

	for k, v := range *result {
		// make each url spec
		urlSpec := model.UrlSpec{}
		urlSpec.Name = k
		urlSpec.Url = v.Url
		urlSpec.Desc = v.Desc
		// add to the array
		arrUrls = append(arrUrls, urlSpec)
	}

	if err != nil {
		log.Errorf("get urls failed: %s", err.Error())
		return pb.GenericResponseError(err), err
	}

	return pb.GenericResponseResult(arrUrls), err
}
