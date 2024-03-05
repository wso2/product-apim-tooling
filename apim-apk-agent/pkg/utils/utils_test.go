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
package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	apkmgt "github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/service/apkmgt"
	subscription "github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/subscription"
)

// Mock for apkmgt.EventStreamService_StreamEventsServer
type MockEventStreamServer struct {
	mock.Mock
	grpc.ServerStream // Embedding grpc.ServerStream
}

func (m *MockEventStreamServer) Send(event *subscription.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

// Override Recv method
func (m *MockEventStreamServer) Recv() (*subscription.Event, error) {
	args := m.Called()
	return args.Get(0).(*subscription.Event), args.Error(1)
}

func (m *MockEventStreamServer) Context() context.Context {
	return context.Background()
}

// TestAddDeleteAndGetAllClientConnections tests AddClientConnection, DeleteClientConnection,
// and GetAllClientConnections functions
func TestAddDeleteAndGetAllClientConnections(t *testing.T) {
	// Create a new mock server
	mockServer := new(MockEventStreamServer)

	// Test AddClientConnection
	clientID := "1ed4120e15fab0833626a36d08ffa3ad7bb9d9a6"
	AddClientConnection(clientID, mockServer)
	assert.Equal(t, 1, len(clientConnections), "Client connection should be added")

	// Test GetAllClientConnections
	allConnections := GetAllClientConnections()
	assert.Equal(t, 1, len(allConnections), "Should return all client connections")

	// Test DeleteClientConnection
	DeleteClientConnection(clientID)
	assert.Equal(t, 0, len(clientConnections), "Client connection should be deleted")
}

func TestSendInitialEventToAllConnectedClients(t *testing.T) {
	// Prepare mock connections
	mockConnection1 := new(MockEventStreamServer)
	mockConnection2 := new(MockEventStreamServer)

	// Set up expectations for Send method
	mockConnection1.On("Send", mock.Anything).Return(nil).Once()
	mockConnection2.On("Send", mock.Anything).Return(nil).Once()

	// Add mock connections to clientConnections
	AddClientConnection("client1", mockConnection1)
	AddClientConnection("client2", mockConnection2)

	// Test positive case: event sent to all clients
	SendInitialEventToAllConnectedClients()

	// Assert that the expectations were met
	mockConnection1.AssertExpectations(t)
	mockConnection2.AssertExpectations(t)

	// Test negative case: no clients connected
	clientConnections = make(map[string]apkmgt.EventStreamService_StreamEventsServer)
	SendInitialEventToAllConnectedClients()

	// Assert that no event is sent when there are no connections
	mockConnection1.AssertNotCalled(t, "Send")
	mockConnection2.AssertNotCalled(t, "Send")
}

func TestSendInitialEvent(t *testing.T) {
	mockConnection := new(MockEventStreamServer)
	mockConnection.On("Send", mock.Anything).Return(nil).Once()
	AddClientConnection("client1", mockConnection)
	SendInitialEventToAllConnectedClients()
	mockConnection.AssertExpectations(t)
	clientConnections = make(map[string]apkmgt.EventStreamService_StreamEventsServer)
	SendInitialEventToAllConnectedClients()
	mockConnection.AssertNotCalled(t, "Send")
}

func TestSendEvent(t *testing.T) {
	// Prepare mock connections
	mockConnection1 := new(MockEventStreamServer)
	mockConnection2 := new(MockEventStreamServer)

	// Set up expectations for Send method
	event := &subscription.Event{ /* Initialize event data */ }
	mockConnection1.On("Send", event).Return(nil).Once()
	mockConnection2.On("Send", event).Return(nil).Once()

	// Add mock connections to clientConnections
	AddClientConnection("client1", mockConnection1)
	AddClientConnection("client2", mockConnection2)

	// Test the function
	SendEvent(event)

	// Assert that the expectations were met
	mockConnection1.AssertExpectations(t)
	mockConnection2.AssertExpectations(t)
}

func TestGetUniqueIDOfApplicationMapping(t *testing.T) {
	td := []struct {
		applicationUUID  string
		subscriptionUUID string
		expectedUniqueID string
	}{
		{"app1", "sub1", "4c6ab46fa7f03acc96e35d9c69fbb2de113845f3"},
		{"app2", "sub2", "211068632fd5179d6f4a1e035a2eddb72dfc334b"},
		{"app3", "sub3", "54a386348dea9b0da068307859827120da33acb8"},
	}

	for _, test := range td {
		t.Run(fmt.Sprintf("Test with appUUID: %s, subUUID: %s", test.applicationUUID, test.subscriptionUUID), func(t *testing.T) {
			actualUniqueID := GetUniqueIDOfApplicationMapping(test.applicationUUID, test.subscriptionUUID)
			if actualUniqueID != test.expectedUniqueID {
				t.Errorf("Expected unique ID: %s, but got: %s", test.expectedUniqueID, actualUniqueID)
			}
		})
	}
}

func TestGetUniqueIDOfApplicationKeyMapping(t *testing.T) {
	td := []struct {
		applicationUUID  string
		keyType          string
		securityScheme   string
		envID            string
		organization     string
		expectedUniqueID string
	}{
		{"app1", "type1", "scheme1", "env1", "org1", "c5b532560059d12143b2e8f12d43268b58f3c659"},
		{"app2", "type2", "scheme2", "env2", "org2", "a236da939779ce8f21f69000fe8b4bf2f66e2792"},
		{"app3", "type3", "scheme3", "env3", "org3", "e068ce4a4f64601073d7f73f09ef868fd2e9afec"},
	}

	for _, test := range td {
		t.Run(fmt.Sprintf("Test with appUUID: %s, keyType: %s, securityScheme: %s, envID: %s, organization: %s",
			test.applicationUUID, test.keyType, test.securityScheme, test.envID, test.organization), func(t *testing.T) {
			actualUniqueID := GetUniqueIDOfApplicationKeyMapping(test.applicationUUID, test.keyType, test.securityScheme, test.envID, test.organization)
			if actualUniqueID != test.expectedUniqueID {
				t.Errorf("Expected unique ID: %s, but got: %s", test.expectedUniqueID, actualUniqueID)
			}
		})
	}
}
