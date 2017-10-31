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
     * @param url                       URL that the http call shoud be made to
     * @param tokenPrefix               Prefix of the token of the request
     * @param token                     Authorization token of the request
     * @param contentType               Content type of the request
     * @param payload                   Body of the request
     * @return                          Returns a Http response
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    HttpResponse makePostRequest(String url, String tokenPrefix, String token, String contentType,
                                 StringEntity payload) throws PluginExecutionException;

    /**
     *
     * @param url                       URL that the http call shoud be made to
     * @param tokenPrefix               Prefix of the token of the request
     * @param token                     Authorization token of the request
     * @param contentType               Content type of the request
     * @param payload                   Body of the request
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    HttpResponse makePutRequest(String url, String tokenPrefix, String token, String contentType, StringEntity payload)
            throws com.smartbear.swaggerhub.plugins.PluginExecutionException;

    /**
     * Makes a GET request and returns the response
     *
     * @param url                       URL that the http call shoud be made to
     * @param tokenPrefix               Prefix of the token of the request
     * @param token                     Authorization token of the request
     * @param contentType               Content type of the request
     * @return                          Returns the response of the request
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    HttpResponse makeGetRequest(String url, String tokenPrefix, String token, String contentType)
            throws com.smartbear.swaggerhub.plugins.PluginExecutionException;
}
