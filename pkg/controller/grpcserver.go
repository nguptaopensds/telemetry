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

package controller

import (

	"net"


	log "github.com/golang/glog"
	pb "github.com/telemetry/pkg/model/proto"
	"google.golang.org/grpc"
)

func NewGrpcServer(port string) *GrpcServer {
	ctrl := NewController()
	return &GrpcServer{
		Controller: ctrl,
		Port:       port,
	}
}

type GrpcServer struct {
	*Controller
	Port string
}

// Run method would start the listen mechanism of controller module.
func (g *GrpcServer) Run() error {
	// New Grpc Server
	s := grpc.NewServer()
	// Register controller service.
	pb.RegisterControllerServer(s, g)
	//pb.RegisterFileShareControllerServer(s, g)

	// Listen the controller server port.
	lis, err := net.Listen("tcp", g.Port)
	if err != nil {
		log.Fatalf("failed to listen: %+v", err)
		return err
	}

	log.Info("osdslet server initialized! Start listening on port:", lis.Addr())

	// Start controller server watching loop.
	defer s.Stop()
	return s.Serve(lis)
}

