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

// Package config contains the implementation and data structures related to configurations and
// configuration (log and adapter config) parsing. If a new configuration is introduced to the adapter
// configuration file, the corresponding change needs to be added to the relevant data stucture as well.
package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"sync"

	toml "github.com/pelletier/go-toml"
	logger "github.com/sirupsen/logrus"
	pkgconf "github.com/wso2/apk/adapter/pkg/config"
	"github.com/wso2/apk/adapter/pkg/logging"
)

var (
	onceConfigRead      sync.Once
	onceGetDefaultVhost sync.Once
	adapterConfig       *Config
	defaultVhost        map[string]string
	e                   error
)

// DefaultGatewayName represents the name of the default gateway
const DefaultGatewayName = "Default"

// DefaultGatewayVHost represents the default vhost of default gateway environment if it is not configured
// For /testtoken and /health check, if user not configured default env, we have no vhost
const DefaultGatewayVHost = "localhost"

const (
	// The environtmental variable which represents the path of the distribution in host machine.
	mgwHomeEnvVariable = "MGW_HOME"
	// RelativeConfigPath is the relative file path where the configuration file is.
	relativeConfigPath = "/conf/config.toml"
)

// Constants related to utility functions
const (
	tenantDomainSeparator = "@"
	superTenantDomain     = "carbon.super"
)

// ReadConfigs implements adapter configuration read operation. The read operation will happen only once, hence
// the consistancy is ensured.
//
// If the "MGW_HOME" variable is set, the configuration file location would be picked relative to the
// variable's value ("/conf/config.toml"). otherwise, the "MGW_HOME" variable would be set to the directory
// from where the executable is called from.
//
// Returns the configuration object that is initialized with default values. Changes to the default
// configuration object is achieved through the configuration file.
func ReadConfigs() (*Config, error) {
	onceConfigRead.Do(func() {
		adapterConfig = defaultConfig
		_, err := os.Stat(pkgconf.GetApkHome() + relativeConfigPath)
		if err != nil {
			loggerConfig.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Configuration file not found : %s", err.Error()),
				Severity:  logging.BLOCKER,
				ErrorCode: 1000,
			})
		}
		content, readErr := ioutil.ReadFile(pkgconf.GetApkHome() + relativeConfigPath)
		if readErr != nil {
			loggerConfig.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Error reading configurations : %s", readErr.Error()),
				Severity:  logging.BLOCKER,
				ErrorCode: 1001,
			})
			return
		}
		parseErr := toml.Unmarshal(content, adapterConfig)
		if parseErr != nil {
			loggerConfig.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Error parsing the configurations : %s", parseErr.Error()),
				Severity:  logging.BLOCKER,
				ErrorCode: 1002,
			})
			return
		}

		adapterConfig.resolveDeprecatedProperties()

		pkgconf.ResolveConfigEnvValues(reflect.ValueOf(&(adapterConfig.ControlPlane)).Elem(), "ControlPlane", true)
	})
	return adapterConfig, e
}

// SetConfig sets the given configuration to the adapter configuration
func SetConfig(conf *Config) {
	adapterConfig = conf
}

// SetDefaultConfig sets the default configuration to the adapter configuration
func SetDefaultConfig() {
	adapterConfig = defaultConfig
}

// ReadLogConfigs implements adapter/proxy log-configuration read operation.The read operation will happen only once, hence
// the consistancy is ensured.
//
// If the "MGW_HOME" variable is set, the log configuration file location would be picked relative to the
// variable's value ("/conf/log_config.toml"). otherwise, the "MGW_HOME" variable would be set to the directory
// from where the executable is called from.
//
// Returns the log configuration object mapped from the configuration file during the startup.
func ReadLogConfigs() *pkgconf.LogConfig {
	return pkgconf.ReadLogConfigs()
}

// ClearLogConfigInstance removes the existing configuration.
// Then the log configuration can be re-initialized.
func ClearLogConfigInstance() {
	pkgconf.ClearLogConfigInstance()
}

// GetLogConfigPath returns the file location of the log-config path
func GetLogConfigPath() (string, error) {
	return pkgconf.GetLogConfigPath()
}

// GetMgwHome reads the MGW_HOME environmental variable and returns the value.
// This represent the directory where the distribution is located.
// If the env variable is not present, the directory from which the executable is triggered will be assigned.
func GetMgwHome() string {
	return pkgconf.GetApkHome()
}

// GetControlPlaneConnectedTenantDomain returns the tenant domain of the user used to authenticate with event hub.
func GetControlPlaneConnectedTenantDomain() string {
	// Read configurations to get the control plane authenticated user
	conf, _ := ReadConfigs()

	// Populate data from the config
	cpTenantAdminUser := conf.ControlPlane.Username
	tenantDomain := strings.Split(cpTenantAdminUser, tenantDomainSeparator)
	if len(tenantDomain) > 1 {
		return tenantDomain[len(tenantDomain)-1]
	}
	return superTenantDomain
}

func (config *Config) resolveDeprecatedProperties() {
	if config.ControlPlane.ServiceURLDeprecated != UnassignedAsDeprecated {
		printDeprecatedWarningLog("controlPlane.serviceUrl", "controlPlane.serviceURL")
		config.ControlPlane.ServiceURL = config.ControlPlane.ServiceURLDeprecated
	}
}

func printDeprecatedWarningLog(deprecatedTerm, currentTerm string) {
	logger.Warnf("%s is deprecated. Use %s instead", deprecatedTerm, currentTerm)
}
