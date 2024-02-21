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
package managementserver

import (
	apkmgt "github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/service/apkmgt"
	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/loggers"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/utils"
	"google.golang.org/grpc/metadata"
)

// EventServer struct use to hold event server
type EventServer struct {
	apkmgt.UnimplementedEventStreamServiceServer
}

// StreamEvents streams events to the enforcer
func (s EventServer) StreamEvents(req *apkmgt.Request, srv apkmgt.EventStreamService_StreamEventsServer) error {
	// Read metadata from the request context
	md, ok := metadata.FromIncomingContext(srv.Context())
	if !ok {
		logger.LoggerMgtServer.Errorf("error : %v", "Failed to get metadata from the request context")
		return nil
		// Handle the case where metadata is not present
	}
	commonControllerID := md.Get("common-controller-uuid")
	logger.LoggerMgtServer.Debugf("Enforcer ID : %v", commonControllerID[0])
	utils.AddClientConnection(commonControllerID[0], srv)
	utils.SendInitialEvent(srv)
	<-srv.Context().Done()
	logger.LoggerMgtServer.Infof("Connection closed by the client : %v", commonControllerID[0])
	utils.DeleteClientConnection(commonControllerID[0])
	return nil // Client closed the connection
}
