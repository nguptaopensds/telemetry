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
This module implements the entry into operations of storageDock module.
*/

package dock

import (
	"context"

	"net"

	log "github.com/golang/glog"

	"github.com/telemetry/contrib/drivers"
	"github.com/telemetry/pkg/dock/discovery"


	"github.com/telemetry/pkg/model"
	pb "github.com/telemetry/pkg/model/proto"
	"google.golang.org/grpc"

)

// dockServer is used to implement pb.DockServer
type dockServer struct {
	Port string
	// Discoverer represents the mechanism of DockHub discovering the storage
	// capabilities from different backends.
	Discoverer discovery.DockDiscoverer
	// Driver represents the specified backend resource. This field is used
	// for initializing the specified volume driver.
	//Driver drivers.VolumeDriver
	// Metrics driver to collect metrics
	MetricDriver drivers.MetricDriver

	// FileShareDriver represents the specified backend resource. This field is used
	// for initializing the specified file share driver.
	//FileShareDriver filesharedrivers.FileShareDriver
}

// NewDockServer returns a dockServer instance.
func NewDockServer(dockType, port string) *dockServer {
	return &dockServer{
		Port:       port,
		Discoverer: discovery.NewDockDiscoverer(dockType),
	}
}

// Run method would automatically discover dock and pool resources from
// backends, and then start the listen mechanism of dock module.
func (ds *dockServer) Run() error {
	// New Grpc Server
	s := grpc.NewServer()
	// Register dock service.
	pb.RegisterProvisionDockServer(s, ds)
	//pb.RegisterAttachDockServer(s, ds)
	//pb.RegisterFileShareDockServer(s, ds)

	// Trigger the discovery and report loop so that the dock service would
	// update the capabilities from backends automatically.
	if err := func() error {
		var err error
		if err = ds.Discoverer.Init(); err != nil {
			return err
		}
		ctx := &discovery.Context{
			StopChan: make(chan bool),
			ErrChan:  make(chan error),
			MetaChan: make(chan string),
		}
		go discovery.DiscoveryAndReport(ds.Discoverer, ctx)
		go func(ctx *discovery.Context) {
			if err = <-ctx.ErrChan; err != nil {
				log.Error("when calling capabilty report method:", err)
				ctx.StopChan <- true
			}
		}(ctx)
		return err
	}(); err != nil {
		return err
	}

	// Listen the dock server port.
	lis, err := net.Listen("tcp", ds.Port)
	if err != nil {
		log.Fatalf("failed to listen: %+v", err)
		return err
	}

	log.Info("Dock server initialized! Start listening on port:", lis.Addr())

	// Start dock server watching loop.
	defer s.Stop()
	return s.Serve(lis)
}




// Collect the specified metrics from the metric driver
func (ds *dockServer) CollectMetrics(ctx context.Context, opt *pb.CollectMetricsOpts) (*pb.GenericResponse, error) {
	log.V(5).Info("in dock CollectMetrics methods")
	ds.MetricDriver = drivers.InitMetricDriver(opt.GetDriverName())

	defer drivers.CleanMetricDriver(ds.MetricDriver)

	log.Infof("dock server receive CollectMetrics request, vr =%s", opt)

	result, err := ds.MetricDriver.CollectMetrics()
	if err != nil {
		log.Errorf("error occurred in dock module for collect metrics: %s", err.Error())
		return pb.GenericResponseError(err), err
	}

	return pb.GenericResponseResult(result), nil
}



// GetMetrics method is only defined to make ProvisioinDock service consistent with
// ProvisionController service, so this method is not allowed to be called.
func (ds *dockServer) GetMetrics(context.Context, *pb.GetMetricsOpts) (*pb.GenericResponse, error) {
	return nil, &model.NotImplementError{"method GetMetrics has not been implemented yet"}
}

// GetUrls method is only defined to make ProvisioinDock service consistent with
// ProvisionController service, so this method is not allowed to be called.
func (ds *dockServer) GetUrls(context.Context, *pb.NoParams) (*pb.GenericResponse, error) {
	return nil, &model.NotImplementError{"method GetUrls has not been implemented yet"}
}
