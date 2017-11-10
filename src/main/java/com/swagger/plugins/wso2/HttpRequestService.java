package com.swagger.plugins.wso2;

import com.smartbear.swaggerhub.plugins.PluginExecutionException;
import org.apache.http.HttpResponse;
import org.apache.http.entity.StringEntity;

/*****************************************************************
 * Interface name :   HttpRequestService
 * Abstract Methods : makePostRequest, makePutRequest, makeGetRequest
 * Functionality :    Contains the signatures of the methods to be implemented to make the http requests
 * Visibility :       Public
 * ****************************************************************/
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
            throws com.smartbear.swaggerhub.plugins.PluginExecutionException;

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
            throws com.smartbear.swaggerhub.plugins.PluginExecutionException;
}
