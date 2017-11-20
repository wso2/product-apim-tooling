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

import com.smartbear.swaggerhub.plugins.PluginExecutionException;
import com.swagger.plugins.wso2.HttpRequestService;
import com.swagger.plugins.wso2.WSO2Api;
import org.apache.commons.io.IOUtils;
import org.apache.http.HttpEntity;
import org.apache.http.HttpResponse;
import org.apache.http.StatusLine;
import org.apache.http.entity.StringEntity;
import org.junit.Before;
import org.junit.Test;
import org.mockito.Mockito;

public class WSO2ApiTest {

    WSO2Api api;
    HttpRequestService httpRequestService;
    String accessToken;
    String payload;
    StringEntity payloadEntity;

    /**
     * This method initializes the variables before tests are executed
     *
     * @throws Exception
     */
    @Before
    public void init() throws Exception {

        httpRequestService = Mockito.mock(HttpRequestService.class);
        api = new WSO2Api(httpRequestService);

        accessToken = "f96f3947-c5a5-3d79-9ff9-2729de4f9d42";

        payload = "{\"name\":\"Simple Inventory API\",\"context\":\"/simple\",\"version\":\"1.0.0\",\"description\":" +
                "\"This is a simple API\",\"apiDefinition\":\"{\\\"swagger\\\":\\\"2.0\\\",\\\"info\\\":" +
                "{\\\"description\\\":\\\"This is a simple API\\\",\\\"version\\\":\\\"1.0.0\\\"," +
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

    /**
     * This test method is for checking whether the API is updated if the API already exists in the WSO2 API management
     * cloud
     *
     * @throws Exception
     */
    @Test
    public void testSaveApiForApiAlreadyExists() throws Exception {
        String stream = "{\n" +
                "    \"count\": 1,\n" +
                "    \"next\": \"\",\n" +
                "    \"previous\": \"\",\n" +
                "    \"list\": [\n" +
                "        {\n" +
                "            \"id\": \"5e26110f-9991-4778-ba36-fd117ddd1fa9\",\n" +
                "            \"name\": \"Simple Inventory API\",\n" +
                "            \"description\": \"This is a simple API\",\n" +
                "            \"context\": \"/t/testusercom8168/simple\",\n" +
                "            \"version\": \"1.0.0\",\n" +
                "            \"provider\": \"fpk11@geronra.com@testusercom8168\",\n" +
                "            \"status\": \"CREATED\",\n" +
                "            \"thumbnailUri\": null\n" +
                "        }\n" +
                "    ]\n" +
                "}";
        mockGetRequest(stream, 200);
        mockPutRequest(200);

        api.saveAPI(accessToken, payload);
        Mockito.verify(httpRequestService, Mockito.times(1)).makePutRequest(Mockito.anyString(),
                Mockito.anyString(), Mockito.anyString(), Mockito.anyString(), Mockito.any(StringEntity.class));
    }

    /**
     * This test method is for checking whether the API is created if the API is not found among the existing APIs in
     * the WSO2 API management cloud
     *
     * @throws Exception
     */
    @Test
    public void testSaveApiForNoMatchingApiFound() throws Exception {
        String stream = "{\n" +
                "    \"count\": 2,\n" +
                "    \"next\": \"\",\n" +
                "    \"previous\": \"\",\n" +
                "    \"list\": [\n" +
                "        {\n" +
                "            \"id\": \"d425dd83-1851-4371-8c83-996aba6a45ac\",\n" +
                "            \"name\": \"Sample Implicit Flow OAuth2 Project\",\n" +
                "            \"description\": \"This is an example of using OAuth2 Implicit Flow in a " +
                "specification to describe security to your API.\",\n" +
                "            \"context\": \"/t/testusercom8168/api/v1\",\n" +
                "            \"version\": \"1.0.0\",\n" +
                "            \"provider\": \"fpk11@geronra.com@testusercom8168\",\n" +
                "            \"status\": \"CREATED\"\n" +
                "        },\n" +
                "        {\n" +
                "            \"id\": \"69d85f37-dd1a-4e08-a2d9-9392bdd5530e\",\n" +
                "            \"name\": \"Swagger Petstore\",\n" +
                "            \"description\": \"This is a sample server Petstore yahoo server.  " +
                "You can find\\nout more about Swagger at\\n[http://swagger.io](http://swagger.io) or " +
                "on\\n[irc.freenode.net, #swagger](http://swagger.io/irc/).\\n\",\n" +
                "            \"context\": \"/t/testusercom8168/v2\",\n" +
                "            \"version\": \"1.0.0\",\n" +
                "            \"provider\": \"fpk11@geronra.com@testusercom8168\",\n" +
                "            \"status\": \"CREATED\"\n" +
                "        }\n" +
                "    ]\n" +
                "}";
        mockGetRequest(stream, 200);
        mockPostRequest(201);

        api.saveAPI(accessToken, payload);
        Mockito.verify(httpRequestService, Mockito.times(1)).makePostRequest(Mockito.anyString(),
                Mockito.anyString(), Mockito.anyString(), Mockito.anyString(), Mockito.any(StringEntity.class));
    }

    /**
     * This test method is to check whether the API is created if there are no APIs in the WSO2 API management cloud
     *
     * @throws Exception
     */
    @Test
    public void testSaveApiForNoApiInTheCloud() throws Exception {
        String stream = "{\n" +
                "    \"count\": 0,\n" +
                "    \"next\": \"\",\n" +
                "    \"previous\": \"\",\n" +
                "    \"list\": []\n" +
                "}";
        mockGetRequest(stream, 200);
        mockPostRequest(201);

        api.saveAPI(accessToken, payload);

        Mockito.verify(httpRequestService, Mockito.times(1)).makePostRequest(Mockito.anyString(),
                Mockito.anyString(), Mockito.anyString(), Mockito.anyString(), Mockito.any(StringEntity.class));
    }

    /**
     * This test method is to check whether an exception is thrown if no API list is returned from the WSO2 API
     * management cloud
     *
     * @throws Exception
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForNoApiListIsReturenedError() throws Exception {
        HttpResponse responseGet = Mockito.mock(HttpResponse.class);
        StatusLine statusLineGet = Mockito.mock(StatusLine.class);
        Mockito.when(responseGet.getStatusLine()).thenReturn(statusLineGet);
        Mockito.when(statusLineGet.getStatusCode()).thenReturn(999);
        Mockito.when(httpRequestService.makeGetRequest(Mockito.anyString(), Mockito.anyString(), Mockito.anyString(),
                Mockito.anyString())).thenReturn(responseGet);

        api.saveAPI(accessToken, payload);
    }

    /**
     * This test method is to check whether an exception is thrown if the request to search APIs is unauthorized
     *
     * @throws Exception
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForUnauthorizedRequest() throws Exception {
        HttpResponse responseGet = Mockito.mock(HttpResponse.class);
        StatusLine statusLineGet = Mockito.mock(StatusLine.class);
        Mockito.when(responseGet.getStatusLine()).thenReturn(statusLineGet);
        Mockito.when(statusLineGet.getStatusCode()).thenReturn(401);
        Mockito.when(httpRequestService.makeGetRequest(Mockito.anyString(), Mockito.anyString(), Mockito.anyString(),
                Mockito.anyString())).thenReturn(responseGet);

        api.saveAPI(accessToken, payload);
    }

    /**
     * This test method is to check whether an exception is thrown if the content type is not supported
     *
     * @throws Exception
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForUnsupportedMediaType() throws Exception {
        HttpResponse responseGet = Mockito.mock(HttpResponse.class);
        StatusLine statusLineGet = Mockito.mock(StatusLine.class);
        Mockito.when(responseGet.getStatusLine()).thenReturn(statusLineGet);
        Mockito.when(statusLineGet.getStatusCode()).thenReturn(406);
        Mockito.when(httpRequestService.makeGetRequest(Mockito.anyString(), Mockito.anyString(), Mockito.anyString(),
                Mockito.anyString())).thenReturn(responseGet);

        api.saveAPI(accessToken, payload);

    }

    /**
     * This test method is to check whether an exception is thrown the API update request contains bad content
     *
     * @throws Exception
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForBadApiUpdateRequest() throws Exception {
        String stream = "{\n" +
                "    \"count\": 1,\n" +
                "    \"next\": \"\",\n" +
                "    \"previous\": \"\",\n" +
                "    \"list\": [\n" +
                "        {\n" +
                "            \"id\": \"5e26110f-9991-4778-ba36-fd117ddd1fa9\",\n" +
                "            \"name\": \"Simple Inventory API\",\n" +
                "            \"description\": \"This is a simple API\",\n" +
                "            \"context\": \"/t/testusercom8168/simple\",\n" +
                "            \"version\": \"1.0.0\",\n" +
                "            \"provider\": \"fpk11@geronra.com@testusercom8168\",\n" +
                "            \"status\": \"CREATED\",\n" +
                "            \"thumbnailUri\": null\n" +
                "        }\n" +
                "    ]\n" +
                "}";
        mockGetRequest(stream, 200);
        mockPutRequest(400);
        api.saveAPI(accessToken, payload);
    }

    /**
     * This test method is to check whether an exception is thrown if the request is not conditional
     *
     * @throws Exception
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForConditionalUpdateRequest() throws Exception {
        String stream = "{\n" +
                "    \"count\": 1,\n" +
                "    \"next\": \"\",\n" +
                "    \"previous\": \"\",\n" +
                "    \"list\": [\n" +
                "        {\n" +
                "            \"id\": \"5e26110f-9991-4778-ba36-fd117ddd1fa9\",\n" +
                "            \"name\": \"Simple Inventory API\",\n" +
                "            \"description\": \"This is a simple API\",\n" +
                "            \"context\": \"/t/testusercom8168/simple\",\n" +
                "            \"version\": \"1.0.0\",\n" +
                "            \"provider\": \"fpk11@geronra.com@testusercom8168\",\n" +
                "            \"status\": \"CREATED\",\n" +
                "            \"thumbnailUri\": null\n" +
                "        }\n" +
                "    ]\n" +
                "}";
        mockGetRequest(stream, 200);
        mockPutRequest(403);
        api.saveAPI(accessToken, payload);
    }

    /**
     * This test method is to check whether an exception is thrown if the resource to update is not found
     *
     * @throws Exception
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForResourceNotFoundToUpdate() throws Exception {
        String stream = "{\n" +
                "    \"count\": 1,\n" +
                "    \"next\": \"\",\n" +
                "    \"previous\": \"\",\n" +
                "    \"list\": [\n" +
                "        {\n" +
                "            \"id\": \"5e26110f-9991-4778-ba36-fd117ddd1fa9\",\n" +
                "            \"name\": \"Simple Inventory API\",\n" +
                "            \"description\": \"This is a simple API\",\n" +
                "            \"context\": \"/t/testusercom8168/simple\",\n" +
                "            \"version\": \"1.0.0\",\n" +
                "            \"provider\": \"fpk11@geronra.com@testusercom8168\",\n" +
                "            \"status\": \"CREATED\",\n" +
                "            \"thumbnailUri\": null\n" +
                "        }\n" +
                "    ]\n" +
                "}";
        mockGetRequest(stream, 200);
        mockPutRequest(404);
        api.saveAPI(accessToken, payload);
    }

    /**
     * This test method is to check whether an exception is thrown if one of the preconditions are not met
     *
     * @throws Exception
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForUpdateRequestPrecondition() throws Exception {
        String stream = "{\n" +
                "    \"count\": 1,\n" +
                "    \"next\": \"\",\n" +
                "    \"previous\": \"\",\n" +
                "    \"list\": [\n" +
                "        {\n" +
                "            \"id\": \"5e26110f-9991-4778-ba36-fd117ddd1fa9\",\n" +
                "            \"name\": \"Simple Inventory API\",\n" +
                "            \"description\": \"This is a simple API\",\n" +
                "            \"context\": \"/t/testusercom8168/simple\",\n" +
                "            \"version\": \"1.0.0\",\n" +
                "            \"provider\": \"fpk11@geronra.com@testusercom8168\",\n" +
                "            \"status\": \"CREATED\",\n" +
                "            \"thumbnailUri\": null\n" +
                "        }\n" +
                "    ]\n" +
                "}";
        mockGetRequest(stream, 200);
        mockPutRequest(412);
        api.saveAPI(accessToken, payload);
    }

    /**
     * This test method is to check whether an exception is thrown if the API update request is unauthorized
     *
     * @throws Exception
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForUnauthorizedUpdateRequest() throws Exception {
        String stream = "{\n" +
                "    \"count\": 1,\n" +
                "    \"next\": \"\",\n" +
                "    \"previous\": \"\",\n" +
                "    \"list\": [\n" +
                "        {\n" +
                "            \"id\": \"5e26110f-9991-4778-ba36-fd117ddd1fa9\",\n" +
                "            \"name\": \"Simple Inventory API\",\n" +
                "            \"description\": \"This is a simple API\",\n" +
                "            \"context\": \"/t/testusercom8168/simple\",\n" +
                "            \"version\": \"1.0.0\",\n" +
                "            \"provider\": \"fpk11@geronra.com@testusercom8168\",\n" +
                "            \"status\": \"CREATED\",\n" +
                "            \"thumbnailUri\": null\n" +
                "        }\n" +
                "    ]\n" +
                "}";
        mockGetRequest(stream, 200);
        mockPutRequest(401);
        api.saveAPI(accessToken, payload);
    }

    /**
     *
     *
     * @throws Exception
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForApiIsNotUpdated() throws Exception {
        String stream = "{\n" +
                "    \"count\": 1,\n" +
                "    \"next\": \"\",\n" +
                "    \"previous\": \"\",\n" +
                "    \"list\": [\n" +
                "        {\n" +
                "            \"id\": \"5e26110f-9991-4778-ba36-fd117ddd1fa9\",\n" +
                "            \"name\": \"Simple Inventory API\",\n" +
                "            \"description\": \"This is a simple API\",\n" +
                "            \"context\": \"/t/testusercom8168/simple\",\n" +
                "            \"version\": \"1.0.0\",\n" +
                "            \"provider\": \"fpk11@geronra.com@testusercom8168\",\n" +
                "            \"status\": \"CREATED\",\n" +
                "            \"thumbnailUri\": null\n" +
                "        }\n" +
                "    ]\n" +
                "}";
        mockGetRequest(stream, 200);
        mockPutRequest(999);
        api.saveAPI(accessToken, payload);
        Mockito.verify(httpRequestService, Mockito.times(1)).makePutRequest(Mockito.anyString(),
                Mockito.anyString(), Mockito.anyString(), Mockito.anyString(), Mockito.any(StringEntity.class));
    }

    /**
     * This test method is to check whether an exception is thrown if the API create request has bad content
     *
     * @throws Exception
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForBadContentInApiCreationRequest() throws Exception {
        String stream = "{\n" +
                "    \"count\": 0,\n" +
                "    \"next\": \"\",\n" +
                "    \"previous\": \"\",\n" +
                "    \"list\": []\n" +
                "}";
        mockGetRequest(stream, 200);
        mockPostRequest(400);
        api.saveAPI(accessToken, payload);
    }

    /**
     * This test method is to check whether an exception is thrown if API create request is unauthorized
     *
     * @throws Exception
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForUnauthorizedApiCreationRequest() throws Exception {
        String stream = "{\n" +
                "    \"count\": 0,\n" +
                "    \"next\": \"\",\n" +
                "    \"previous\": \"\",\n" +
                "    \"list\": []\n" +
                "}";
        mockGetRequest(stream, 200);
        mockPostRequest(401);
        api.saveAPI(accessToken, payload);
    }

    /**
     * This test method is to check whether an exception is thrown if media type of the API create request is not
     * supported
     *
     * @throws Exception
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForUnsupportedMediaTypeInApiCreateRequest() throws Exception {
        String stream = "{\n" +
                "    \"count\": 0,\n" +
                "    \"next\": \"\",\n" +
                "    \"previous\": \"\",\n" +
                "    \"list\": []\n" +
                "}";
        mockGetRequest(stream, 200);
        mockPostRequest(415);
        api.saveAPI(accessToken, payload);
    }

    /**
     *
     *
     * @throws Exception
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForApiIsNotCreatedInCloud() throws Exception {
        String stream = "{\n" +
                "    \"count\": 0,\n" +
                "    \"next\": \"\",\n" +
                "    \"previous\": \"\",\n" +
                "    \"list\": []\n" +
                "}";
        mockGetRequest(stream, 200);
        mockPostRequest(999);
        api.saveAPI(accessToken, payload);
    }

    /**
     * Contains code for mocking the GET request
     *
     * @param stream Hard coded data stream for mocking
     * @param code Response code of the request
     * @throws Exception
     */
    private void mockGetRequest(String stream, int code) throws Exception {
        HttpResponse responseGet = Mockito.mock(HttpResponse.class);
        StatusLine statusLineGet = Mockito.mock(StatusLine.class);
        HttpEntity httpEntity = Mockito.mock(HttpEntity.class);
        Mockito.when(responseGet.getEntity()).thenReturn(httpEntity);
        Mockito.when(httpEntity.getContent()).thenReturn(IOUtils.toInputStream(stream));
        Mockito.when(responseGet.getStatusLine()).thenReturn(statusLineGet);
        Mockito.when(statusLineGet.getStatusCode()).thenReturn(code);
        Mockito.when(httpRequestService.makeGetRequest(Mockito.anyString(), Mockito.anyString(), Mockito.anyString(),
                Mockito.anyString())).thenReturn(responseGet);
    }

    /**
     * Contains code for mocking the POST request
     *
     * @param code Response code of the request
     * @throws Exception
     */
    private void mockPostRequest(int code) throws Exception {
        HttpResponse responsePost = Mockito.mock(HttpResponse.class);
        StatusLine statusLinePost = Mockito.mock(StatusLine.class);
        Mockito.when(responsePost.getStatusLine()).thenReturn(statusLinePost);
        Mockito.when(statusLinePost.getStatusCode()).thenReturn(code);
        Mockito.when(httpRequestService.makePostRequest(Mockito.anyString(), Mockito.anyString(), Mockito.anyString(),
                Mockito.anyString(), Mockito.any(StringEntity.class))).thenReturn(responsePost);
    }

    /**
     * Contains code for mocking the PUT request
     *
     * @param code Response code of the request
     * @throws Exception
     */
    private void mockPutRequest(int code) throws Exception {
        HttpResponse responsePut = Mockito.mock(HttpResponse.class);
        StatusLine statusLinePut = Mockito.mock(StatusLine.class);
        Mockito.when(responsePut.getStatusLine()).thenReturn(statusLinePut);
        Mockito.when(statusLinePut.getStatusCode()).thenReturn(code);
        Mockito.when(httpRequestService.makePutRequest(Mockito.anyString(), Mockito.anyString(), Mockito.anyString(),
                Mockito.anyString(), Mockito.any(StringEntity.class))).thenReturn(responsePut);
    }

}
