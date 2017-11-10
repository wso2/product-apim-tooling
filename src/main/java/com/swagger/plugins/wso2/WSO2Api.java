package com.swagger.plugins.wso2;

import com.smartbear.swaggerhub.plugins.PluginExecutionException;
import io.swagger.models.Swagger;
import io.swagger.util.Json;
import org.apache.commons.io.Charsets;
import org.apache.commons.io.IOUtils;
import org.apache.http.HttpResponse;
import org.apache.http.entity.StringEntity;
import org.json.simple.JSONArray;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.io.UnsupportedEncodingException;


/*****************************************************************
 * Class name :    WSO2Api
 * Attributes :    httpRequestService
 * Constants :     PUBLISHER_API_ENDPOINT
 * Methods :       getApiIdentifier, updateApi, saveApi
 * Functionality : Contains the methods to export/update the API to/in the WSO2 API Cloud
 * Visibility :    Public
 * ****************************************************************/
public class WSO2Api {

    private Logger log = LoggerFactory.getLogger(WSO2Api.class);

    private static final String PUBLISHER_API_ENDPOINT = "https://api.cloud.wso2.com/api/am/publisher/v0.11/apis/";

    private HttpRequestService httpRequestService;

    /**
     * Instantiates an instance of HttpRequestService interface with the instantiation of WSO2Api class
     *
     * @param httpRequestService Interface with method signatures to be implemented
     */
    public WSO2Api(HttpRequestService httpRequestService) {
        this.httpRequestService = httpRequestService;
    }

    /**
     * Checks whether the API already exists in the cloud according to the name and the version of the API
     *
     * @param accessToken               Access Token to authorize SwaggerHub to access WSO2 API Cloud
     * @param swagger                   The POJO of the swagger definition
     * @return                          Returns true if the API already exists, else returns false
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    private String getApiIdentifier(String accessToken, Swagger swagger) throws PluginExecutionException {

        HttpResponse response;
        String content;
        JSONObject responseJson;

        response = httpRequestService.makeGetRequest(PUBLISHER_API_ENDPOINT, "Bearer", accessToken,
                "application/json");

        if (response.getStatusLine().getStatusCode() == 200) {
            log.debug("The API list is returned");
        } else {
            if (response.getStatusLine().getStatusCode() == 406) {
                log.error("Not acceptable, the requested media is not supported");
                throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Not acceptable, the " +
                        "requested media is not supported");
            } else if (response.getStatusLine().getStatusCode() == 401) {
                log.error("Unauthorized request");
                throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Unauthorized request");
            } else {
                log.error("The API list is not returned");
                throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "The API list is not" +
                        " returned");
            }
        }

        try {
            content = new String(IOUtils.toByteArray(response.getEntity().getContent()), Charsets.UTF_8);
        } catch (IOException e) {
            log.error("Error while searching APIs.", e);
            throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Error while searching APIs");
        }

        JSONParser parser = new JSONParser();
        try {
            responseJson = (JSONObject) parser.parse(content);
        } catch (ParseException e) {
            log.error("Error while parsing search API response", e);
            throw  new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Error parsing search API" +
                    " response to json");
        }

        String version = swagger.getInfo().getVersion();
        String name = swagger.getInfo().getTitle();

        JSONArray apiList = (JSONArray) responseJson.get("list");

        JSONObject api;
        int length;
        int apiIndex;

        if (!apiList.isEmpty()) {
            length = apiList.size();
            for (apiIndex = 0; apiIndex < length; apiIndex++) {
                api = (JSONObject) apiList.get(apiIndex);
                if (name.equals(api.get("name")) && version.equals(api.get("version"))) {
                    return api.get("id").toString();
                } else {
                    if (apiIndex == length - 1) {
                        return null;
                    }
                }
            }
        } else {
            return null;
        }
        return null;
    }

    /**
     * MakES a HTTP request to the publisher API to update an API
     *
     * @param accessToken               Access Token to authorize SwaggerHub to access WSO2 API Cloud
     * @param apiIdentifier             Identifier of the API
     * @param payload                   Body of the request to be made
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    private void updateApi(String accessToken, String apiIdentifier,
                          StringEntity payload) throws PluginExecutionException {

        HttpResponse response = httpRequestService.makePutRequest(PUBLISHER_API_ENDPOINT + apiIdentifier,
                "Bearer",
                accessToken, "application/json", payload);

        if (response.getStatusLine().getStatusCode() == 200) {
            log.debug("The API is updated");
        } else {
            if (response.getStatusLine().getStatusCode() == 400) {
                log.error("Bad Request. Invalid request or validation error");
                throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Bad Request. Invalid " +
                        "request or validation error");
            } else if (response.getStatusLine().getStatusCode() == 403) {
                log.error("Forbidden. The request must be conditional but no condition has been specified.");
                throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Forbidden. The request " +
                        "must be conditional but no condition " +
                        "has been specified.");
            } else if (response.getStatusLine().getStatusCode() == 404) {
                log.error("Nothing to update");
                throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Nothing to update");
            } else if (response.getStatusLine().getStatusCode() == 412) {
                log.error("Precondition Failed. The request has not been performed because one of the preconditions" +
                        " is not met.");
                throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Precondition Failed. The " +
                        "request has not been performed because one of the preconditions is not met.");
            } else if (response.getStatusLine().getStatusCode() == 401) {
                log.error("Unauthorized request");
                throw new com.smartbear.swaggerhub.plugins.
                        PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Unauthorized request");
            } else {
                log.error("The API is not updated");
                throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "The API is not updated");
            }
        }
    }

    /**
     * Makes a HTTP request to create an API
     *
     * @param accessToken Access Token to authorize SwaggerHub to access WSO2 API Cloud
     * @param payload     Body of the request to be made
     * @throws PluginExecutionException
     */
    private void createApi(String accessToken, StringEntity payload) throws PluginExecutionException {

        log.debug("Creating the API in the cloud");
        HttpResponse postResponse = httpRequestService.makePostRequest(PUBLISHER_API_ENDPOINT, "Bearer",
                accessToken, "application/json", payload);

        if (postResponse.getStatusLine().getStatusCode() == 201) {
            log.debug("The API is created in the cloud");
        } else {
            if (postResponse.getStatusLine().getStatusCode() == 400) {
                log.error("Error creating the API, already exists with a different context");
                throw new com.smartbear.swaggerhub.plugins.
                        PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Bad content");
            } else if (postResponse.getStatusLine().getStatusCode() == 415) {
                log.error("Unsupported media type");
                throw new com.smartbear.swaggerhub.plugins.
                        PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Error creating the API," +
                        " unsupported media type");
            } else if (postResponse.getStatusLine().getStatusCode() == 401) {
                log.error("Unauthorized request");
                throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Unauthorized request");
            } else {
                log.debug("The API is not created in the cloud");
                throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "The API is not" +
                        " created in the WSO2 API Cloud");
            }
        }

    }

    /**
     * Creates an API in the api cloud and prints the response of the details of the API made
     *
     * @param accessToken               Access Token to authorize SwaggerHub to access WSO2 API Cloud
     * @param payload                   Body of the API http request to create an API
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    public void saveAPI(String accessToken, String payload) throws PluginExecutionException {

        StringEntity creationPayload;
        Swagger swagger;
        JSONObject swaggerDefinitionJson;

        JSONParser parser = new JSONParser();

        try {
            swaggerDefinitionJson = (JSONObject) parser.parse(payload);
        } catch (ParseException e) {
            log.error("Error parsing payload", e);
            throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Error parsing payload");
        }

        String definition = swaggerDefinitionJson.get("apiDefinition").toString();

        try {
            swagger = Json.mapper().readValue(definition, Swagger.class);
        } catch (IOException e) {
            log.error("Error parsing swagger definition", e);
            throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Error parsing swagger " +
                    "definition");
        }

        String apiIdentifier = getApiIdentifier(accessToken, swagger);

        try {
            creationPayload = new StringEntity(payload);
        } catch (UnsupportedEncodingException e) {
            log.error("The character encoding is not supported for the payload", e);
            throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "The character encoding is" +
                    " not supported");
        }

        if (apiIdentifier == null) {
            createApi(accessToken, creationPayload);
        } else {
            updateApi(accessToken, apiIdentifier, creationPayload);
        }
    }
}

