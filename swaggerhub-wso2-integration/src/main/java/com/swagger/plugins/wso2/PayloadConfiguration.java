/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package com.swagger.plugins.wso2;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.smartbear.swaggerhub.plugins.PluginExecutionException;
import io.swagger.models.Swagger;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 *
 * Class with the methods to configure the API creation payload
 *
 */
public class PayloadConfiguration {

    private static final Logger log = LoggerFactory.getLogger(PayloadConfiguration.class);

    /**
     * Configures the payload to create an API in the cloud setting relevant values from the swagger definition.
     *
     * @param swagger                   The POJO of the swagger definition
     * @return                          Returns the configured payload
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    public String configurePayload(Swagger swagger, String swaggerJson) throws PluginExecutionException {

        String payload;

        log.debug("Setting default values for the api creation payload");
        String[] schemes = {"http", "https"};
        String[] defaultTier = {"Unlimited"};
        String visibility = "PUBLIC";

        log.debug("Getting info from swagger POJO");
        String name = swagger.getInfo().getTitle();
        String version = swagger.getInfo().getVersion();
        String description = swagger.getInfo().getDescription();

        log.debug("Creating the POJO for the payload to create the API");
        PayloadStructure structure = new PayloadStructure();

        log.debug("Setting the values for the POJO to create an API in the cloud");
        structure.setName(name);
        structure.setVersion(version);
        structure.setDescription(description);
        structure.setContext(swagger.getBasePath());
        structure.setApiDefinition(swaggerJson);
        structure.setIsDefaultVersion(false);
        structure.setTransport(schemes);
        structure.setTiers(defaultTier);
        structure.setVisibility(visibility);
        structure.setEndpointConfig("");
        CorsConfiguration configuration = new CorsConfiguration(false);
        structure.setCorsConfiguration(configuration);

        try {
            log.debug("Converting the POJO to a json string");
            ObjectMapper objectMapper = new ObjectMapper();
            payload = objectMapper.writeValueAsString(structure);
        } catch (JsonProcessingException e) {
            log.error("Error while converting the POJO to json", e);
            throw new com.smartbear.swaggerhub.plugins.
                    PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Error when converting pojo " +
                    "to json");
        }

        log.debug("Returning the payload for creating the API in the cloud");
        return payload;
    }
}
