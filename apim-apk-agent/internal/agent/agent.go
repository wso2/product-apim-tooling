/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

// Package agent contains the implementation to start the agent
package agent

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/wso2/apk/adapter/pkg/health"
	"github.com/wso2/apk/common-go-libs/loggers"
	"github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/config"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/internal/eventhub"
	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/internal/loggers"
	logging "github.com/wso2/product-apim-tooling/apim-apk-agent/internal/logging"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/internal/messaging"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/internal/synchronizer"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/managementserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var (
	debug       bool
	onlyLogging bool

	port    uint
	alsPort uint

	mode string
)

const (
	ads                      = "ads"
	amqpProtocol             = "amqp"
	grpcMaxConcurrentStreams = 1000000
)

func init() {
	flag.BoolVar(&debug, "debug", true, "Use debug logging")
	flag.BoolVar(&onlyLogging, "onlyLogging", false, "Only demo AccessLogging Service")
	flag.UintVar(&port, "port", 18000, "Management server port")
	flag.UintVar(&alsPort, "als", 18090, "Accesslog server port")
	flag.StringVar(&mode, "ads", ads, "Management server type (ads, xds, rest)")
}

// Run starts the XDS server and Rest API server.
func Run(conf *config.Config) {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt)
	// TODO: (VirajSalaka) Support the REST API Configuration via flags only if it is a valid requirement
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.LoggerInternalMsg.Debugf("Run method started with context : %v", ctx)

	// log config watcher
	watcherLogConf, _ := fsnotify.NewWatcher()
	logConfigPath, errC := config.GetLogConfigPath()
	if errC == nil {
		errC = watcherLogConf.Add(logConfigPath)
	}

	if errC != nil {
		logger.LoggerInternalMsg.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Error reading the log configs, error: %v", errC.Error()))
	}

	logger.LoggerInternalMsg.Info("Starting apim-apk-agent ....")
	eventHubEnabled := conf.ControlPlane.Enabled

	// Load initial data from control plane
	eventhub.LoadInitialData(conf)

	if eventHubEnabled {
		var connectionURLList = conf.ControlPlane.BrokerConnectionParameters.EventListeningEndpoints
		if strings.Contains(connectionURLList[0], amqpProtocol) {
			go messaging.ProcessEvents(conf)
		}
	}

	// Load initial KM data from control plane
	synchronizer.FetchKeyManagersOnStartUp(conf)

	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions, grpc.KeepaliveParams(
		keepalive.ServerParameters{
			Time:    time.Duration(5 * time.Minute),
			Timeout: time.Duration(20 * time.Second),
		}),
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
	)
	publicKeyLocation, privateKeyLocation, truststoreLocation := config.GetKeyLocations()
	cert, err := config.GetServerCertificate(publicKeyLocation, privateKeyLocation)

	caCertPool := config.GetTrustedCertPool(truststoreLocation)

	if err == nil {
		grpcOptions = append(grpcOptions, grpc.Creds(
			credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{cert},
				ClientAuth:   tls.RequireAndVerifyClientCert,
				ClientCAs:    caCertPool,
			}),
		))
	} else {
		loggers.LoggerAPKOperator.Warn("failed to initiate the ssl context: ", err)
		panic(err)
	}

	grpcOptions = append(grpcOptions, grpc.KeepaliveParams(
		keepalive.ServerParameters{
			Time:    time.Duration(5 * time.Minute),
			Timeout: time.Duration(20 * time.Second),
		}),
	)
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1100, logging.BLOCKER, "Failed to listen on port: %v, error: %v", port, err.Error()))
	}
	apkmgt.RegisterEventStreamServiceServer(grpcServer, &managementserver.EventServer{})
	loggers.LoggerAPKOperator.Info("port: ", port, " APK agent Listening for gRPC connections")
	go func() {
		loggers.LoggerAPKOperator.Info("Starting GRPC server.")
		health.CommonEnforcerGrpcService.SetStatus(true)
		if err = grpcServer.Serve(lis); err != nil {
			health.CommonEnforcerGrpcService.SetStatus(false)
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error1101, logging.BLOCKER, "Failed to start XDS GRPS server, error: %v", err.Error()))
		}
	}()
OUTER:
	for {
		select {
		case l := <-watcherLogConf.Events:
			switch l.Op.String() {
			case "WRITE":
				logger.LoggerInternalMsg.Info("Loading updated log config file...")
				config.ClearLogConfigInstance()
				logger.UpdateLoggers()
			}
		case s := <-sig:
			switch s {
			case os.Interrupt:
				logger.LoggerInternalMsg.Info("Shutting down...")
				break OUTER
			}
		}
	}
	logger.LoggerInternalMsg.Info("Bye!")
}
