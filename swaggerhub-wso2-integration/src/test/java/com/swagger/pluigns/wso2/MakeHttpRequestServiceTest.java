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

import com.swagger.plugins.wso2.MakeHttpRequestService;
import org.apache.commons.io.IOUtils;
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

        String payload = IOUtils.toString(getClass().getClassLoader()
                .getResourceAsStream("api-creation-payload.json"));

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
