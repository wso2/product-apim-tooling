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

package com.swagger.pluigns.wso2;

import com.swagger.plugins.wso2.PayloadConfiguration;
import io.swagger.models.Swagger;
import io.swagger.util.Json;
import org.apache.commons.io.IOUtils;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;

public class PayloadConfigurationTest {

    Swagger swagger;
    String expectedPayload;
    String swaggerJson;

    /**
     * This method initializes the variables before tests are executed
     *
     * @throws Exception Thrown if an exception occur
     */
    @Before
    public void init() throws Exception {

        swaggerJson = IOUtils.toString(getClass().getClassLoader()
                .getResourceAsStream("simple-inventory-api.json"));

        swagger = Json.mapper().readValue(swaggerJson, Swagger.class);
    }

    /**
     * This test method is to check whether the payload is returned correctly for creating the API in the cloud
     *
     * @throws Exception Thrown if an exception occur
     */
    @Test
    public void testConfigurePayload() throws Exception {

        expectedPayload = IOUtils.toString(getClass().getClassLoader()
                .getResourceAsStream("api-creation-payload.json"));

        PayloadConfiguration payloadConfiguration = new PayloadConfiguration();
        String payload = payloadConfiguration.configurePayload(swagger, swaggerJson);
        Assert.assertEquals(expectedPayload, payload);
    }
}
