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

package config

import (
	"sync"
	"time"
)

// Experimenting asynchronous communication between go routines using channels
// This uses singleton pattern where creating a single channel for communication
//
// To get a instance of the channel for a data publisher go routine
//
//	`publisher := NewSender()`
//
// Create a receiver channel in worker go routine
// receiver := NewReceiver()
//
// From publisher go routine, feed string value to the channel
// publisher<- "some value"
//
// In worker go routine, read the value sent by the publisher
// message := <-receiver
var once sync.Once

// C represents the channel to identify modifications added to the configuration file
// TODO: (VirajSalaka) remove this as unused.
var (
	C chan string // better to be interface{} type which could send any type of data.
)

// NewSender initializes the channel if it is not created an returns
func NewSender() chan string {
	once.Do(func() {
		C = make(chan string)
	})
	return C
}

// NewReceiver initializes the channel if it is not created an returns
func NewReceiver() chan string {
	once.Do(func() {
		C = make(chan string)
	})
	return C
}

const (
	//UnassignedAsDeprecated is used by the configurations which are deprecated.
	UnassignedAsDeprecated string = "unassigned-as-deprecated"
)

// Config represents the adapter configuration.
// It is created directly from the configuration toml file.
// Note :
//
//	Don't use toml tag for configuration properties as it may affect environment variable based
//	config resolution.
type Config struct {
	ControlPlane controlPlane `toml:"controlPlane"`
	DataPlane    dataPlane    `toml:"dataPlane"`
	Agent        agent        `toml:"agent"`
	// Metric represents configurations to expose/export go metrics
	Metrics metrics `toml:"metrics"`
}
type agent struct {
	Enabled    bool
	Keystore   keystore
	TrustStore truststore
	Mode       string
}
type keystore struct {
	KeyPath  string
	CertPath string
}

type truststore struct {
	Location string
}

// ControlPlane struct contains configurations related to the API Manager
type controlPlane struct {
	Enabled    bool
	ServiceURL string
	// Deprecated: Use ServiceURL instead.
	ServiceURLDeprecated       string `toml:"serviceUrl"`
	Username                   string
	Password                   string
	SyncApisOnStartUp          bool
	SendRevisionUpdate         bool
	EnvironmentLabels          []string
	RetryInterval              time.Duration
	SkipSSLVerification        bool
	BrokerConnectionParameters brokerConnectionParameters
	HTTPClient                 httpClient
	RequestWorkerPool          requestWorkerPool
	InternalKeyIssuer          string
	ClientID                   string
	ClientSecret               string
}

// Dataplane struct contains the configurations related to the APK
type dataPlane struct {
	Enabled            bool
	K8ResourceEndpoint string
	Namespace          string
}

type requestWorkerPool struct {
	PoolSize int
	// QueueSizePerPool is the size of the queue for each worker pool
	QueueSizePerPool int
	// PauseTimeAfterFailure is the time duration to pause the worker pool after a failure
	PauseTimeAfterFailure time.Duration
}
type brokerConnectionParameters struct {
	EventListeningEndpoints []string
	ReconnectInterval       time.Duration
	ReconnectRetryCount     int
}

type httpClient struct {
	RequestTimeOut time.Duration
}

// Metrics defines the configuration for metrics collection.
type metrics struct {
	Enabled bool
	Type    string
	Port    int32
}
