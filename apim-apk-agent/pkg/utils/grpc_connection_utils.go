package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wso2/apk/common-go-libs/constants"
	"github.com/wso2/apk/common-go-libs/loggers"
	apkmgt "github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/subscription"
)

var clientConnections = make(map[string]apkmgt.EventStreamService_StreamEventsServer)

// AddClientConnection adds a client connection to the map
func AddClientConnection(clientID string, stream apkmgt.EventStreamService_StreamEventsServer) {
	clientConnections[clientID] = stream
}

// DeleteClientConnection deletes a client connection from the map
func DeleteClientConnection(clientID string) {
	delete(clientConnections, clientID)
}

// GetAllClientConnections returns all client connections
func GetAllClientConnections() map[string]apkmgt.EventStreamService_StreamEventsServer {
	return clientConnections
}

// SendInitialEventToAllConnectedClients sends initial event to all connected clients
func SendInitialEventToAllConnectedClients() {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)

	event := subscription.Event{
		Uuid:      uuid.New().String(),
		Type:      constants.AllEvents,
		TimeStamp: milliseconds,
	}
	loggers.LoggerAPKOperator.Debugf("Sending initial event to all clients: %v", &event)
	for _, stream := range GetAllClientConnections() {
		stream.Send(&event)
	}
}

// SendInitialEvent sends initial event to the enforcer
func SendInitialEvent(srv apkmgt.EventStreamService_StreamEventsServer) {
	currentTime := time.Now()
	milliseconds := currentTime.UnixNano() / int64(time.Millisecond)

	event := subscription.Event{
		Uuid:      uuid.New().String(),
		Type:      constants.AllEvents,
		TimeStamp: milliseconds,
	}
	loggers.LoggerAPKOperator.Debugf("Sending initial event to client: %v", &event)
	srv.Send(&event)
}

// SendEvent sends event to the common-controllers
func SendEvent(event *subscription.Event) {
	loggers.LoggerAPKOperator.Infof("Sending event to all clients: %v", event)
	for clientID, stream := range GetAllClientConnections() {
		err := stream.Send(event)
		if err != nil {
			loggers.LoggerAPKOperator.Errorf("Error sending event to client %s: %v", clientID, err)
		} else {
			loggers.LoggerAPKOperator.Debugf("Event sent to client %s", clientID)
		}
	}
}

// GetUniqueIDOfApplicationMapping returns the unique id of an application mapping
func GetUniqueIDOfApplicationMapping(applicationUUID string, subscriptionUUID string) string {
	hasher := sha1.New()
	hasher.Write([]byte(fmt.Sprintf("%s:%s", applicationUUID, subscriptionUUID)))

	// Get the final hash sum
	hashSum := hasher.Sum(nil)

	// Convert the hash sum to a hexadecimal string
	hashString := hex.EncodeToString(hashSum)
	return hashString
}

// GetUniqueIDOfApplicationKeyMapping returns the unique id of an application key mapping
func GetUniqueIDOfApplicationKeyMapping(applicationUUID string, keyType string, securityScheme string, envID string, organization string) string {
	hasher := sha1.New()
	loggers.LoggerAPKOperator.Infof("Hashing: %s:%s:%s:%s:%s", applicationUUID, keyType, securityScheme, envID, organization)
	hasher.Write([]byte(fmt.Sprintf("%s:%s:%s:%s:%s", applicationUUID, keyType, securityScheme, envID, organization)))
	// Get the final hash sum
	hashSum := hasher.Sum(nil)
	// Convert the hash sum to a hexadecimal string
	hashString := hex.EncodeToString(hashSum)
	return hashString
}
