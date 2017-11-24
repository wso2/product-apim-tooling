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

package com.swagger.plugins.wso2;

import com.smartbear.swaggerhub.plugins.PluginExecutionException;
import org.apache.http.HttpResponse;
import org.apache.http.entity.StringEntity;

/**
 *
 * Interface with the method signatures to make HTTP requests
 *
 */
public interface HttpRequestService {

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
    HttpResponse makePostRequest(String url, String tokenPrefix, String token, String contentType,
                                 StringEntity payload) throws PluginExecutionException;

    /**
     *
     * @param url                       URL that the http call should be made to
     * @param tokenPrefix               Prefix of the token of the request
     * @param token                     Authorization token of the request
     * @param contentType               Content type of the request
     * @param payload                   Body of the request
     * @throws PluginExecutionException Thrown when an exception is caught while the plugin executes
     */
    HttpResponse makePutRequest(String url, String tokenPrefix, String token, String contentType, StringEntity payload)
            throws PluginExecutionException;

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
    HttpResponse makeGetRequest(String url, String tokenPrefix, String token, String contentType)
            throws PluginExecutionException;
}
