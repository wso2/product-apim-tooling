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

import com.smartbear.swaggerhub.plugins.PluginExecutionException;
import org.apache.http.HttpResponse;
import org.apache.http.client.HttpClient;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.client.methods.HttpPut;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.HttpClients;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;

/**
 *
 * Class with the implementations of the methods defined in HttpRequestService interface
 *
 */
public class MakeHttpRequestService implements HttpRequestService {

    private static final Logger log = LoggerFactory.getLogger(MakeHttpRequestService.class);
    private static final String AUTHORIZATION_HEADER = "Authorization";
    private static final String CONTENT_TYPE_HEADER = "Content-Type";

    /**
     * Makes a POST request to the given url and returns the response
     *
     * @param url                       URL that the http call should be made to
     * @param tokenPrefix               Prefix of the token of the request
     * @param token                     Authorization token of the request
     * @param contentType               Content type of the request
     * @param payload                   Body of the request
     * @return                          Returns a Http response
     * @throws PluginExecutionException Thrown when an exception is caught while the plugin executes
     */
    public HttpResponse makePostRequest(String url, String tokenPrefix, String token, String contentType,
                                        StringEntity payload) throws PluginExecutionException {
        HttpResponse response;
        HttpClient httpClient = HttpClients.createDefault();
        HttpPost postRequest = new HttpPost(url);
        postRequest.setHeader(AUTHORIZATION_HEADER, tokenPrefix + " " + token);
        postRequest.setHeader(CONTENT_TYPE_HEADER, contentType);
        postRequest.setEntity(payload);
        try {
            response = httpClient.execute(postRequest);
        } catch (IOException e) {
            log.error("Error making POST request", e);
            throw new com.smartbear.swaggerhub.plugins.
                    PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Error making POST request");
        }
        return response;
    }

    /**
     * Makes a PUT request and returns the response
     *
     * @param url                       URL that the http call should be made to
     * @param tokenPrefix               Prefix of the token of the request
     * @param token                     Authorization token of the request
     * @param contentType               Content type of the request
     * @param payload                   Body of the request
     * @throws PluginExecutionException Thrown when an exception is caught while the plugin executes
     */
    public HttpResponse makePutRequest(String url, String tokenPrefix, String token, String contentType,
                                       StringEntity payload) throws PluginExecutionException {
        HttpResponse response;
        HttpClient httpClient = HttpClients.createDefault();
        HttpPut putRequest = new HttpPut(url);
        putRequest.setHeader(AUTHORIZATION_HEADER, tokenPrefix + " " + token);
        putRequest.setHeader(CONTENT_TYPE_HEADER, contentType);
        putRequest.setEntity(payload);
        try {
            response = httpClient.execute(putRequest);
        } catch (IOException e) {
            log.error("Error making PUT request", e);
            throw new com.smartbear.swaggerhub.plugins.
                    PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Error making PUT request");
        }
        return response;
    }

    /**
     * Makes a GET request and returns the response
     *
     * @param url                       URL that the http call should be made to
     * @param tokenPrefix               Prefix of the token of the request
     * @param token                     Authorization token of the request
     * @param contentType               Content type of the request
     * @return                          Returns the response of the request
     * @throws PluginExecutionException Thrown when an exception is caught while the plugin executes
     */
    public HttpResponse makeGetRequest(String url, String tokenPrefix, String token, String contentType)
            throws com.smartbear.swaggerhub.plugins.PluginExecutionException {
        HttpResponse response;
        HttpClient httpClient = HttpClients.createDefault();
        HttpGet getRequest = new HttpGet(url);
        getRequest.setHeader(AUTHORIZATION_HEADER, tokenPrefix + " " + token);
        getRequest.setHeader(CONTENT_TYPE_HEADER, contentType);
        try {
            response = httpClient.execute(getRequest);
        } catch (IOException e) {
            log.error("Error making GET request", e);
            throw new com.smartbear.swaggerhub.plugins.
                    PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Error making GET request");
        }
        return response;
    }
}
