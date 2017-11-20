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

import com.swagger.plugins.wso2.MakeHttpRequestService;
import org.apache.http.HttpResponse;
import org.apache.http.client.HttpClient;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.client.methods.HttpPut;
import org.apache.http.entity.StringEntity;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;
import org.mockito.Mockito;

public class MakeHttpRequestServiceTest {

    HttpResponse response;
    HttpClient httpClient;

    String url;
    String contentType;
    String accessToken;
    String tokenPrefix;
    StringEntity payloadEntity;
    MakeHttpRequestService makeHttpRequestService;

    @Before
    public void init() throws Exception {
        response = Mockito.mock(HttpResponse.class);
        httpClient = Mockito.mock(HttpClient.class);

        url = "https://api.cloud.wso2.com/api/am/publisher/v0.11/apis/";
        contentType = "application/json";
        accessToken = "10b093ee-e557-38ab-87d0-3e607ef0a684";
        tokenPrefix = "Bearer";
        makeHttpRequestService = new MakeHttpRequestService();

        String payload = "{\"name\":\"Simple Inventory API\",\"context\":\"/simple\",\"version\":\"1.0.0\"," +
                "\"description\":\"This is a simple API\",\"apiDefinition\":\"{\\\"swagger\\\":\\\"2.0\\\"," +
                "\\\"info\\\":{\\\"description\\\":\\\"This is a simple API\\\",\\\"version\\\":\\\"1.0.0\\\"," +
                "\\\"title\\\":\\\"Simple Inventory API\\\",\\\"contact\\\":{\\\"email\\\":" +
                "\\\"you@your-company.com\\\"},\\\"license\\\":{\\\"name\\\":\\\"Apache 2.0\\\"," +
                "\\\"url\\\":\\\"http://www.apache.org/licenses/LICENSE-2.0.html\\\"}},\\\"host\\\":" +
                "\\\"virtserver.swaggerhub.com\\\",\\\"basePath\\\":\\\"/simple\\\",\\\"tags\\\":[{\\\"name\\\":" +
                "\\\"admins\\\",\\\"description\\\":\\\"Secured Admin-only calls\\\"},{\\\"name\\\":" +
                "\\\"developers\\\",\\\"description\\\":\\\"Operations available to regular developers\\\"}]," +
                "\\\"schemes\\\":[\\\"https\\\"],\\\"paths\\\":{\\\"/inventory\\\":{\\\"get\\\":{\\\"tags\\\":" +
                "[\\\"developers\\\"],\\\"summary\\\":\\\"searches inventory\\\",\\\"description\\\":" +
                "\\\"By passing in the appropriate options, you can search for\\\\navailable inventory in the " +
                "system\\\\n\\\",\\\"operationId\\\":\\\"searchInventory\\\",\\\"produces\\\":" +
                "[\\\"application/json\\\"],\\\"parameters\\\":[{\\\"name\\\":\\\"searchString\\\",\\\"in\\\":" +
                "\\\"query\\\",\\\"description\\\":\\\"pass an optional search string for looking up inventory\\\"," +
                "\\\"required\\\":false,\\\"type\\\":\\\"string\\\"},{\\\"name\\\":\\\"skip\\\",\\\"in\\\":" +
                "\\\"query\\\",\\\"description\\\":\\\"number of records to skip for pagination\\\"," +
                "\\\"required\\\":false,\\\"type\\\":\\\"integer\\\",\\\"minimum\\\":0,\\\"format\\\":" +
                "\\\"int32\\\"},{\\\"name\\\":\\\"limit\\\",\\\"in\\\":\\\"query\\\",\\\"description\\\":" +
                "\\\"maximum number of records to return\\\",\\\"required\\\":false,\\\"type\\\":\\\"integer\\\"," +
                "\\\"maximum\\\":50,\\\"minimum\\\":0,\\\"format\\\":\\\"int32\\\"}],\\\"responses\\\":{\\\"200\\\":" +
                "{\\\"description\\\":\\\"search results matching criteria\\\",\\\"schema\\\":{\\\"type\\\":" +
                "\\\"array\\\",\\\"items\\\":{\\\"$ref\\\":\\\"#/definitions/InventoryItem\\\"}}},\\\"400\\\":" +
                "{\\\"description\\\":\\\"bad input parameter\\\"}}},\\\"post\\\":{\\\"tags\\\":[\\\"admins\\\"]," +
                "\\\"summary\\\":\\\"adds an inventory item\\\",\\\"description\\\":" +
                "\\\"Adds an item to the system\\\",\\\"operationId\\\":\\\"addInventory\\\",\\\"consumes\\\":" +
                "[\\\"application/json\\\"],\\\"produces\\\":[\\\"application/json\\\"],\\\"parameters\\\":" +
                "[{\\\"in\\\":\\\"body\\\",\\\"name\\\":\\\"inventoryItem\\\",\\\"description\\\":" +
                "\\\"Inventory item to add\\\",\\\"required\\\":false,\\\"schema\\\":{\\\"$ref\\\":" +
                "\\\"#/definitions/InventoryItem\\\"}}],\\\"responses\\\":{\\\"201\\\":{\\\"description\\\":" +
                "\\\"item created\\\"},\\\"400\\\":{\\\"description\\\":\\\"invalid input, object invalid\\\"}," +
                "\\\"409\\\":{\\\"description\\\":\\\"an existing item already exists\\\"}}}}},\\\"definitions\\\":" +
                "{\\\"InventoryItem\\\":{\\\"type\\\":\\\"object\\\",\\\"required\\\":[\\\"id\\\"," +
                "\\\"manufacturer\\\",\\\"name\\\",\\\"releaseDate\\\"],\\\"properties\\\":{\\\"id\\\":" +
                "{\\\"type\\\":\\\"string\\\",\\\"format\\\":\\\"uuid\\\",\\\"example\\\":" +
                "\\\"d290f1ee-6c54-4b01-90e6-d701748f0851\\\"},\\\"name\\\":{\\\"type\\\":\\\"string\\\"," +
                "\\\"example\\\":\\\"Widget Adapter\\\"},\\\"releaseDate\\\":{\\\"type\\\":\\\"string\\\"," +
                "\\\"format\\\":\\\"int32\\\",\\\"example\\\":\\\"2016-08-29T09:12:33.001Z\\\"},\\\"manufacturer\\\":" +
                "{\\\"$ref\\\":\\\"#/definitions/Manufacturer\\\"}}},\\\"Manufacturer\\\":" +
                "{\\\"required\\\":[\\\"name\\\"],\\\"properties\\\":{\\\"name\\\":{\\\"type\\\":\\" +
                "\"string\\\",\\\"example\\\":\\\"ACME Corporation\\\"},\\\"homePage\\\":{\\\"type\\\":" +
                "\\\"string\\\",\\\"format\\\":\\\"url\\\",\\\"example\\\":\\\"https://www.acme-corp.com\\\"}," +
                "\\\"phone\\\":{\\\"type\\\":\\\"string\\\",\\\"example\\\":\\\"408-867-5309\\\"}}}}}\"," +
                "\"isDefaultVersion\":false,\"transport\":[\"http\",\"https\"],\"tiers\":[\"Unlimited\"]," +
                "\"visibility\":\"PUBLIC\",\"endpointConfig\":\"\",\"corsConfiguration\":" +
                "{\"corsConfigurationEnabled\":false}}";

        payloadEntity = new StringEntity(payload);
    }

    @Test
    public void testMakeGetRequest() throws Exception {
        HttpGet get = Mockito.mock(HttpGet.class);

        Mockito.when(httpClient.execute(get)).thenReturn(response);

        makeHttpRequestService.makeGetRequest(url, tokenPrefix, accessToken, contentType);
        Assert.assertNotNull(response);
    }

    @Test
    public void testMakePutRequest() throws Exception {
        HttpPut put = Mockito.mock(HttpPut.class);

        Mockito.when(httpClient.execute(put)).thenReturn(response);

        makeHttpRequestService.makePutRequest(url, tokenPrefix, accessToken, contentType, payloadEntity);
        Assert.assertNotNull(response);
    }

    @Test
    public void testMakePostRequest() throws Exception {
        HttpPost post = Mockito.mock(HttpPost.class);

        Mockito.when(httpClient.execute(post)).thenReturn(response);

        makeHttpRequestService.makePostRequest(url, tokenPrefix, accessToken, contentType, payloadEntity);
        Assert.assertNotNull(response);
    }

}
