package com.swagger.plugins.wso2;

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

/*****************************************************************
 * Class name : MakeHttpRequestService
 * Implemented interfaces : HttpRequestService
 * Methods : makePostRequest, makePutRequest, makeGetRequest
 * Functionality : Contains the implementations of the methods in the HttpRequestService interface
 * Visibility : Public
 * ****************************************************************/
public class MakeHttpRequestService implements HttpRequestService {

    private static final Logger log = LoggerFactory.getLogger(MakeHttpRequestService.class);

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
    public HttpResponse makePostRequest(String url, String tokenPrefix, String token, String contentType,
                                        StringEntity payload) throws PluginExecutionException {
        HttpResponse response;
        HttpClient httpClient = HttpClients.createDefault();
        HttpPost postRequest = new HttpPost(url);
        postRequest.setHeader("Authorization", tokenPrefix + " " + token);
        postRequest.setHeader("Content-Type", contentType);
        postRequest.setEntity(payload);
        try {
            response = httpClient.execute(postRequest);
        } catch (IOException ioException) {
            log.error("Error making POST request", ioException);
            throw new PluginExecutionException("Error making POST request");
        }
        return response;
    }

    /**
     * Makes a POST request and returns the response
     *
     * @param url                       URL that the http call shoud be made to
     * @param tokenPrefix               Prefix of the token of the request
     * @param token                     Authorization token of the request
     * @param contentType               Content type of the request
     * @param payload                   Body of the request
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    public HttpResponse makePutRequest(String url, String tokenPrefix, String token, String contentType,
                                       StringEntity payload) throws PluginExecutionException {
        HttpResponse response;
        HttpClient httpClient = HttpClients.createDefault();
        HttpPut putRequest = new HttpPut(url);
        putRequest.setHeader("Authorization", tokenPrefix + " " + token);
        putRequest.setHeader("Content-Type", contentType);
        putRequest.setEntity(payload);
        try {
            response = httpClient.execute(putRequest);
        } catch (IOException ioException) {
            log.error("Error making PUT request", ioException);
            throw new PluginExecutionException("Error making PUT request");
        }
        return response;
    }

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
    public HttpResponse makeGetRequest(String url, String tokenPrefix, String token, String contentType)
            throws PluginExecutionException {
        HttpResponse response;
        HttpClient httpClient = HttpClients.createDefault();
        HttpGet getRequest = new HttpGet(url);
        getRequest.setHeader("Authorization", tokenPrefix + " " + token);
        getRequest.setHeader("Content-Type", contentType);
        try {
            response = httpClient.execute(getRequest);
        } catch (IOException ioException) {
            log.error("Error making GET request", ioException);
            throw new PluginExecutionException("Error making GET request");
        }
        return response;
    }
}
