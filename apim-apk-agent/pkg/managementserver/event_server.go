package managementserver

import (
	"github.com/wso2/apk/common-go-libs/loggers"
	apkmgt "github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/service/apkmgt"
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
		loggers.LoggerAPKOperator.Errorf("error : %v", "Failed to get metadata from the request context")
		return nil
		// Handle the case where metadata is not present
	}
	commonControllerID := md.Get("common-controller-uuid")
	loggers.LoggerAPKOperator.Debugf("Enforcer ID : %v", commonControllerID[0])
	utils.AddClientConnection(commonControllerID[0], srv)
	utils.SendInitialEvent(srv)
	<-srv.Context().Done()
	loggers.LoggerAPKOperator.Infof("Connection closed by the client : %v", commonControllerID[0])
	utils.DeleteClientConnection(commonControllerID[0])
	return nil // Client closed the connection
}
