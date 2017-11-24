/*
*  Copyright (c) 2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
     * @throws Exception Thrown if an exception occur
     */
    @Before
    public void init() throws Exception {

        httpRequestService = Mockito.mock(HttpRequestService.class);
        api = new WSO2Api(httpRequestService);

        accessToken = "f96f3947-c5a5-3d79-9ff9-2729de4f9d42";

        payload = IOUtils.toString(getClass().getClassLoader().getResourceAsStream("api-creation-payload.json"));

        payloadEntity = new StringEntity(payload);
    }

    /**
     * This test method is for checking whether the API is updated if the API already exists in the WSO2 API management
     * cloud
     *
     * @throws Exception Thrown if an exception occur
     */
    @Test
    public void testSaveApiForApiAlreadyExists() throws Exception {
        String stream = IOUtils.toString(getClass().getClassLoader()
                .getResourceAsStream("inventory-api-search-response.json"));
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
     * @throws Exception Thrown if an exception occur
     */
    @Test
    public void testSaveApiForNoMatchingApiFound() throws Exception {
        String stream = IOUtils.toString(getClass().getClassLoader()
                .getResourceAsStream("unmatched-apis-search-response.json"));
        mockGetRequest(stream, 200);
        mockPostRequest(201);

        api.saveAPI(accessToken, payload);
        Mockito.verify(httpRequestService, Mockito.times(1)).makePostRequest(Mockito.anyString(),
                Mockito.anyString(), Mockito.anyString(), Mockito.anyString(), Mockito.any(StringEntity.class));
    }

    /**
     * This test method is to check whether the API is created if there are no APIs in the WSO2 API management cloud
     *
     * @throws Exception Thrown if an exception occur
     */
    @Test
    public void testSaveApiForNoApiInTheCloud() throws Exception {
        String stream = IOUtils.toString(getClass().getClassLoader()
                .getResourceAsStream("empty-api-list-search-response.json"));
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
     * @throws Exception Thrown if an exception occur
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
     * @throws Exception Thrown if an exception occur
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
     * @throws Exception Thrown if an exception occur
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
     * @throws Exception Thrown if an exception occur
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForBadApiUpdateRequest() throws Exception {
        String stream = IOUtils.toString(getClass().getClassLoader()
                .getResourceAsStream("inventory-api-search-response.json"));
        mockGetRequest(stream, 200);
        mockPutRequest(400);
        api.saveAPI(accessToken, payload);
    }

    /**
     * This test method is to check whether an exception is thrown if the request is not conditional
     *
     * @throws Exception Thrown if an exception occur
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForConditionalUpdateRequest() throws Exception {
        String stream = IOUtils.toString(getClass().getClassLoader()
                .getResourceAsStream("inventory-api-search-response.json"));
        mockGetRequest(stream, 200);
        mockPutRequest(403);
        api.saveAPI(accessToken, payload);
    }

    /**
     * This test method is to check whether an exception is thrown if the resource to update is not found
     *
     * @throws Exception Thrown if an exception occur
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForResourceNotFoundToUpdate() throws Exception {
        String stream = IOUtils.toString(getClass().getClassLoader()
                .getResourceAsStream("inventory-api-search-response.json"));
        mockGetRequest(stream, 200);
        mockPutRequest(404);
        api.saveAPI(accessToken, payload);
    }

    /**
     * This test method is to check whether an exception is thrown if one of the preconditions are not met
     *
     * @throws Exception Thrown if an exception occur
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForUpdateRequestPrecondition() throws Exception {
        String stream = IOUtils.toString(getClass().getClassLoader()
                .getResourceAsStream("inventory-api-search-response.json"));
        mockGetRequest(stream, 200);
        mockPutRequest(412);
        api.saveAPI(accessToken, payload);
    }

    /**
     * This test method is to check whether an exception is thrown if the API update request is unauthorized
     *
     * @throws Exception Thrown if an exception occur
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForUnauthorizedUpdateRequest() throws Exception {
        String stream = IOUtils.toString(getClass().getClassLoader()
                .getResourceAsStream("inventory-api-search-response.json"));
        mockGetRequest(stream, 200);
        mockPutRequest(401);
        api.saveAPI(accessToken, payload);
    }

    /**
     *
     *
     * @throws Exception Thrown if an exception occur
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForApiIsNotUpdated() throws Exception {
        String stream = IOUtils.toString(getClass().getClassLoader()
                .getResourceAsStream("inventory-api-search-response.json"));
        mockGetRequest(stream, 200);
        mockPutRequest(999);
        api.saveAPI(accessToken, payload);
        Mockito.verify(httpRequestService, Mockito.times(1)).makePutRequest(Mockito.anyString(),
                Mockito.anyString(), Mockito.anyString(), Mockito.anyString(), Mockito.any(StringEntity.class));
    }

    /**
     * This test method is to check whether an exception is thrown if the API create request has bad content
     *
     * @throws Exception Thrown if an exception occur
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForBadContentInApiCreationRequest() throws Exception {
        String stream = IOUtils.toString(getClass().getClassLoader()
                .getResourceAsStream("empty-api-list-search-response.json"));
        mockGetRequest(stream, 200);
        mockPostRequest(400);
        api.saveAPI(accessToken, payload);
    }

    /**
     * This test method is to check whether an exception is thrown if API create request is unauthorized
     *
     * @throws Exception Thrown if an exception occur
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForUnauthorizedApiCreationRequest() throws Exception {
        String stream = IOUtils.toString(getClass().getClassLoader()
                .getResourceAsStream("empty-api-list-search-response.json"));
        mockGetRequest(stream, 200);
        mockPostRequest(401);
        api.saveAPI(accessToken, payload);
    }

    /**
     * This test method is to check whether an exception is thrown if media type of the API create request is not
     * supported
     *
     * @throws Exception Thrown if an exception occur
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForUnsupportedMediaTypeInApiCreateRequest() throws Exception {
        String stream = IOUtils.toString(getClass().getClassLoader()
                .getResourceAsStream("empty-api-list-search-response.json"));
        mockGetRequest(stream, 200);
        mockPostRequest(415);
        api.saveAPI(accessToken, payload);
    }

    /**
     *
     *
     * @throws Exception Thrown if an exception occur
     */
    @Test (expected = PluginExecutionException.class)
    public void testSaveApiForApiIsNotCreatedInCloud() throws Exception {
        String stream = IOUtils.toString(getClass().getClassLoader()
                .getResourceAsStream("empty-api-list-search-response.json"));
        mockGetRequest(stream, 200);
        mockPostRequest(999);
        api.saveAPI(accessToken, payload);
    }

    /**
     * Contains code for mocking the GET request
     *
     * @param stream Hard coded data stream for mocking
     * @param code Response code of the request
     * @throws Exception Thrown if an exception occur
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
     * @throws Exception Thrown if an exception occur
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
     * @throws Exception Thrown if an exception occur
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
