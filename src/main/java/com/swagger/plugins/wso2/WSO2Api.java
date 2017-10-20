package com.swagger.plugins.wso2;

import org.apache.commons.codec.binary.Base64;
import org.apache.commons.io.Charsets;
import org.apache.commons.io.IOUtils;
import org.apache.http.HttpResponse;
import org.apache.http.entity.StringEntity;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.io.UnsupportedEncodingException;


/*****************************************************************
 * Class name : WSO2Api
 * Attributes : httpClient, content, response, parser, encodedIdAndSecret, accessToken
 * Constants : API_CREATE_CLOUD_URL, DYNAMIC_CLIENT_REGISTRATION_URL, TOKEN_API_URL
 * Methods : getAuthorizationPayload, getClientIdAndSecret, getAccessToken, saveAPI, makeHttpRequest
 * Functionality : Contains the methods to obtain the access token and push the API to the cloud
 * Visibility : Public
 * ****************************************************************/
public class WSO2Api {

    private Logger log = LoggerFactory.getLogger(WSO2Api.class);

    private static final String API_CREATE_CLOUD_URL = "https://api.cloud.wso2.com/api/am/publisher/v0.11/apis";
    private static final String DYNAMIC_CLIENT_REGISTRATION_URL = "https://api.cloud.wso2.com/client-registration/" +
                                                                  "v0.11/register";
    private static final String TOKEN_API_URL = "https://gateway.api.cloud.wso2.com/token";

    private HttpRequestService httpRequestService;

    public WSO2Api(HttpRequestService httpRequestService) {
        this.httpRequestService = httpRequestService;
    }

    /**
     * Returns the payload for request of getting client id and secret.
     *
     * @param email             Email of the cloud account to export the API
     * @param organizationKey   The key generated in the API cloud for the given credentials, unique for the WSO2 API
     *                          cloud account
     * @return                  Returns the authorization payload for obtaining client id and client secret
     */
    private JSONObject getAuthorizationPayload(String email, String organizationKey) throws ParseException {
        String payload =  "{\n" +
                "    \"clientName\": \"rest_api_publisher\",\n" +
                "    \"owner\":\"" + email + "@" + organizationKey + "\",\n" +
                "    \"grantType\": \"password refresh_token\",\n" +
                "    \"saasApp\": true\n" +
                "}";

        JSONParser parser = new JSONParser();
        return (JSONObject) parser.parse(payload);
    }

    /**
     * Returns the Base64 encoded string of the format (clientId:clientSecret) for a given email organization key
     * and password.
     *
     * @param email                     Email of the cloud account to export the API
     * @param organizationKey           The key generated in the API cloud for the given credentials, unique for the
     *                                  cloud
     *                                  account
     * @param password                  Password of the cloud account to export the API
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    private String getClientIdAndSecret(String email, String organizationKey, String password) throws
            PluginExecutionException {

        HttpResponse response;
        StringEntity authorizationPayload;
        String stringToEncode = email + "@" + organizationKey + ":" + password;
        String encodedString = Base64.encodeBase64String(stringToEncode.getBytes(Charsets.UTF_8));

        try {
            authorizationPayload = new StringEntity(getAuthorizationPayload(email, organizationKey)
                    .toString());
        } catch (UnsupportedEncodingException unsupportedEncodingException) {
            log.error("The character encoding is not supported for the payload", unsupportedEncodingException);
            throw new PluginExecutionException("The character encoding is not supported");
        } catch (ParseException e) {
            log.error("Erro while parsing");
            throw new PluginExecutionException("Error while parsing");
        }

        String clientId;
        String clientSecret;

        try {
            log.debug("Calling dynamic client registration endpoint");

            response = httpRequestService.makePostRequest(DYNAMIC_CLIENT_REGISTRATION_URL,
                    "Basic", encodedString, "application/json", authorizationPayload);

            if (response.getStatusLine().getStatusCode() == 401) {
                log.error("Error making the request to dynamic client registration endpoint, the request" +
                        " is unauthorized");
                throw new PluginExecutionException("Unauthorized request, check email, organizationKey, password");
            } else if (response.getStatusLine().getStatusCode() == 400) {
                log.error("Error making the request to dynamic client registration endpoint, malformed request body");
                throw new PluginExecutionException("Bad Request, check content");
            } else {
                if (response.getStatusLine().getStatusCode() == 200) {
                    log.debug("Client id and client secret request successful");
                } else {
                    log.error("Client id and client secret request unsuccessful");
                    throw new PluginExecutionException("Requesting client id and client secret was unsuccessful");
                }
            }

            String content = new String(IOUtils.toByteArray(response.getEntity().getContent()), Charsets.UTF_8);
            log.debug("Received the response");

            JSONParser parser = new JSONParser();
            JSONObject clientIdAndSecretJson = (JSONObject) parser.parse(content);
            clientId = clientIdAndSecretJson.get("clientId").toString();
            clientSecret = clientIdAndSecretJson.get("clientSecret").toString();

        } catch (IOException ioException) {
            log.error("Error while making the http call to dynamic client registration endpoint", ioException);
            throw new PluginExecutionException("Error while making http call to dynamic client registration endpoint");
        } catch (ParseException parseException) {
            log.error("Error while parsing the http response to json", parseException);
            throw new PluginExecutionException("Error parsing the response to json");
        }

        String toEncode = clientId + ":" + clientSecret;

        log.debug("Obtaining the encoded clientId and clientSecret");
        return Base64.encodeBase64String(toEncode.getBytes(Charsets.UTF_8));
    }

    /**
     * Returns the access token obtained using the client Id and client secret
     *
     * @param email                     Email of the cloud account to export the API
     * @param organizationKey           The key generated in the API cloud for the given credentials, unique for the
     *                                  cloud
     *                                  account
     * @param password                  Password of the cloud account to export the API
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    private String getAccessToken(String email, String organizationKey, String password) throws
            PluginExecutionException {

        HttpResponse response;
        String content;
        StringEntity authorizationPayload;
        JSONObject accessTokenJson;
        String encodedIdAndSecret = getClientIdAndSecret(email, organizationKey, password);

        try {
            authorizationPayload = new StringEntity("scope=apim:api_create&grant_type=password&username=" +
                    email + "@" + organizationKey + "&password=" + password);
        } catch (UnsupportedEncodingException unsupportedEncodingException) {
            log.error("The character encoding is not supported for the payload", unsupportedEncodingException);
            throw new PluginExecutionException("The character encoding is not supported");
        }

        try {
            log.debug("Issuing REST call to Token API");

            response = httpRequestService.makePostRequest(TOKEN_API_URL, "Basic", encodedIdAndSecret,
                    "application/x-www-form-urlencoded", authorizationPayload);

            if (response.getStatusLine().getStatusCode() == 401) {
                log.error("Error making the request to token API, the request is unauthorized");
                throw new PluginExecutionException("Unauthorized request, check email, organizationKey, password");
            } else if (response.getStatusLine().getStatusCode() == 400) {
                log.error("Error making the request to token API, malformed request body");
                throw new PluginExecutionException("Bad Request, check content");
            } else {
                if (response.getStatusLine().getStatusCode() == 200) {
                    log.debug("Access token request successful");
                } else {
                    log.error("Access token request unsuccessful");
                    throw new PluginExecutionException("Error obtaining the access token");
                }
            }

            content = new String(IOUtils.toByteArray(response.getEntity().getContent()), Charsets.UTF_8);
            log.debug("Received the response");

            JSONParser parser = new JSONParser();
            accessTokenJson = (JSONObject) parser.parse(content);

        } catch (IOException ioException) {
            log.error("Error while making the http call to token API", ioException);
            throw new PluginExecutionException("Error while making http call to token API");
        } catch (ParseException parseException) {
            log.error("Error while parsing the http response to json", parseException);
            throw new PluginExecutionException("Error while parsing the http response to json");
        }

        log.debug("Obtaining the access token");
        return accessTokenJson.get("access_token").toString();
    }

    /**
     * Creates an API in the api cloud and prints the response of the details of the API made
     *
     * @param email                     Email of the cloud account to export the API
     * @param organizationKey           The key generated in the API cloud for the given credentials, unique for the
     *                                  cloud
     * @param password                  Password of the cloud account to export the API
     * @param payload                   Payload for the api creation http request
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    public void saveAPI(String email, String organizationKey, String password, String payload) throws
            PluginExecutionException {

        HttpResponse response;
        StringEntity creationPayload;
        String accessToken = getAccessToken(email, organizationKey, password);

        try {
            creationPayload = new StringEntity(payload);
        } catch (UnsupportedEncodingException unsupportedEncodingException) {
            log.error("The character encoding is not supported for the payload", unsupportedEncodingException);
            throw new PluginExecutionException("The character encoding is not supported");
        }

            log.debug("Creating the API in the cloud");
            response = httpRequestService.makePostRequest(API_CREATE_CLOUD_URL, "Bearer", accessToken,
                "application/json", creationPayload);

        if (response.getStatusLine().getStatusCode() == 401) {
            log.error("Error while creating the API, the request is unauthorized");
            throw new PluginExecutionException("Unauthorized request");
        } else if (response.getStatusLine().getStatusCode() == 409) {
            log.error("Error while creating the API, the API already exists");
            throw new PluginExecutionException("An API with the same name and the context already exists");
        } else if (response.getStatusLine().getStatusCode() == 400) {
            log.error("Error creating the API, already exists with a different context");
            throw new PluginExecutionException("Bad content");
        } else if (response.getStatusLine().getStatusCode() == 415) {
            log.error("Unsupported media type");
            throw new PluginExecutionException("Error creating the API, unsupported media type");
        } else {
            if (response.getStatusLine().getStatusCode() == 201) {
                log.debug("The API is created in the cloud");
            } else {
                log.debug("The API is not created in the cloud");
            }
        }
    }
}

