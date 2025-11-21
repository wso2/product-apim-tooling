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

// Package messaging holds the implementation for event listeners functions
package messaging

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wso2/product-apim-tooling/apim-apk-agent/internal/eventhub"
	k8sclient "github.com/wso2/product-apim-tooling/apim-apk-agent/internal/k8sClient"
	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/internal/loggers"
	eventhubTypes "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/eventhub/types"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/logging"
	msg "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/messaging"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// constants related to key manager events
const (
	keyManagerConfigEvent = "key_manager_configuration"
	actionAdd             = "add"
	actionUpdate          = "update"
	actionDelete          = "delete"
	superTenantDomain     = "carbon.super"
)

// handleKMConfiguration
func handleKMConfiguration(c client.Client) {
	for d := range msg.KeyManagerChannel {
		var notification msg.EventKeyManagerNotification
		var keyManager eventhubTypes.KeyManager
		unmarshalErr := json.Unmarshal([]byte(string(d.Body)), &notification)
		if unmarshalErr != nil {
			logger.LoggerMessaging.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Error occurred while unmarshalling key manager event data %v", unmarshalErr.Error()),
				Severity:  logging.CRITICAL,
				ErrorCode: 2000,
			})
			return
		}
		logger.LoggerMessaging.Infof("Event %s is received", notification.Event.PayloadData.EventType)

		var decodedByte, err = base64.StdEncoding.DecodeString(notification.Event.PayloadData.Value)

		if err != nil {
			if _, ok := err.(base64.CorruptInputError); ok {
				logger.LoggerMessaging.ErrorC(logging.ErrorDetails{
					Message:   "\nbase64 input is corrupt, check the provided key",
					Severity:  logging.MINOR,
					ErrorCode: 2001,
				})
			}
			logger.LoggerMessaging.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Error occurred while decoding the notification event %v", err.Error()),
				Severity:  logging.CRITICAL,
				ErrorCode: 2002,
			})
			return
		}

		if strings.EqualFold(keyManagerConfigEvent, notification.Event.PayloadData.EventType) {
			if strings.EqualFold(actionDelete, notification.Event.PayloadData.Action) {
				k8sclient.DeleteTokenIssuersCR(c, notification.Event.PayloadData.Name, notification.Event.PayloadData.TenantDomain)
			} else if decodedByte != nil {
				logger.LoggerMessaging.Debugf("decoded stream %s", string(decodedByte))
				kmConfigMapErr := json.Unmarshal([]byte(string(decodedByte)), &keyManager)
				if kmConfigMapErr != nil {
					logger.LoggerMessaging.ErrorC(logging.ErrorDetails{
						Message:   fmt.Sprintf("Error occurred while unmarshalling key manager config map %v", kmConfigMapErr),
						Severity:  logging.CRITICAL,
						ErrorCode: 2003,
					})
					return
				}
				if strings.EqualFold(actionAdd, notification.Event.PayloadData.Action) ||
					strings.EqualFold(actionUpdate, notification.Event.PayloadData.Action) {
					logger.LoggerMessaging.Debugf("Key Managers received: %v", keyManager)
					resolvedKeyManager := eventhub.MarshalKeyManager(&keyManager)
					logger.LoggerMessaging.Debugf("Resolved Key Managers received: %v", resolvedKeyManager)
					if strings.EqualFold(actionAdd, notification.Event.PayloadData.Action) {
						k8sclient.CreateAndUpdateTokenIssuersCR(resolvedKeyManager, c)
					} else {
						err = k8sclient.UpdateTokenIssuersCR(resolvedKeyManager, c)
						if err != nil {
							k8sclient.CreateAndUpdateTokenIssuersCR(resolvedKeyManager, c)
						}
					}
				}
			}
		}
		logger.LoggerMessaging.Info("handle: deliveries channel closed")
		d.Ack(false)
	}
}
