package com.swagger.plugins.wso2;

import org.apache.commons.codec.binary.Base64;
import org.apache.commons.io.Charsets;
import org.apache.commons.io.IOUtils;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;


/*****************************************************************
 * Class name : Wso2Api
 * Attributes : httpClient, content, response, parser, encodedIdAndSecret, accessToken
 * Constants : API_CREATE_CLOUD_URL, DYNAMIC_CLIENT_REGISTRATION_URL, TOKEN_API_URL
 * Methods : getAuthorizationPayload, getClientIdAndSecret, getAccessToken, saveAPI, getHttpResponse
 * Functionality : Contains the methods to obtain the access token and push the API to the cloud
 * Visibility : Public
 * ****************************************************************/
public class Wso2Api {

    private Logger log = LoggerFactory.getLogger(Wso2Api.class);

    private static final String API_CREATE_CLOUD_URL = "https://api.cloud.wso2.com/api/am/publisher/v0.11/apis";
    private static final String DYNAMIC_CLIENT_REGISTRATION_URL = "https://api.cloud.wso2.com/client-registration/" +
                                                                  "v0.11/register";
    private static final String TOKEN_API_URL = "https://gateway.api.cloud.wso2.com/token";

    private CloseableHttpClient httpClient;
    private String content;
    private CloseableHttpResponse response;
    private JSONParser parser;

    private String encodedIdAndSecret;
    private String accessToken;

    /**
     * Returns the payload for request of getting client id and secret.
     *
     * @param email             Email of the cloud account to export the API
     * @param organizationKey   The key generated in the API cloud for the given credentials, unique for the WSO2 API
     *                          cloud account
     * @return                  Returns the authorization payload for obtaining client id and client secret
     */
    private String getAuthorizationPayload(String email, String organizationKey) {
        return "{\n" +
                "    \"clientName\": \"rest_api_publisher\",\n" +
                "    \"owner\":\"" + email + "@" + organizationKey + "\",\n" +
                "    \"grantType\": \"password refresh_token\",\n" +
                "    \"saasApp\": true\n" +
                "}";
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
     * @throws ParseException           Thrown if an error is occurred when parsing the content of the response to json
     *                                  object
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    private void setClientIdAndSecret(String email, String organizationKey, String password) throws IOException,
            ParseException, PluginExecutionException {

        String stringToEncode = email + "@" + organizationKey + ":" + password;
        String encodedString = Base64.encodeBase64String(stringToEncode.getBytes(Charsets.UTF_8));
        StringEntity authorizationPayload = new StringEntity(getAuthorizationPayload(email, organizationKey));

        String clientId;
        String clientSecret;

        try {
            log.info("Calling dynamic client registration endpoint");
            response = getHttpResponse(DYNAMIC_CLIENT_REGISTRATION_URL, encodedString, authorizationPayload);

            if (response.getStatusLine().getStatusCode() == 401) {
                log.error("Error making the request to dynamic client registration endpoint, the request" +
                        " is unauthorized");
                throw new PluginExecutionException("Unauthorized request, check email, organizationKey, password");
            } else if (response.getStatusLine().getStatusCode() == 400) {
                log.error("Error making the request to dynamic client registration endpoint, malformed request body");
                throw new PluginExecutionException("Bad Request, check content");
            }

            content = new String(IOUtils.toByteArray(response.getEntity().getContent()), Charsets.UTF_8);
            log.info("Received the response");

            parser = new JSONParser();
            JSONObject clientIdAndSecretJson = (JSONObject) parser.parse(content);
            clientId = clientIdAndSecretJson.get("clientId").toString();
            clientSecret = clientIdAndSecretJson.get("clientSecret").toString();

        } catch (IOException ioException) {
            log.error("Error while making the http call to dynamic client registration endpoint", ioException);
            throw ioException;
        } catch (ParseException parseException) {
            log.error("Error while parsing the http response to json", parseException);
            throw parseException;
        } finally {
            httpClient.close();
        }

        String toEncode = clientId + ":" + clientSecret;

        log.info("Obtaining the encoded clientId and clientSecret");
        encodedIdAndSecret = Base64.encodeBase64String(toEncode.getBytes(Charsets.UTF_8));
    }

    /**
     * Returns the access token obtained using the client Id and client secret
     *
     * @param email                     Email of the cloud account to export the API
     * @param organizationKey           The key generated in the API cloud for the given credentials, unique for the
     *                                  cloud
     *                                  account
     * @param password                  Password of the cloud account to export the API
     * @throws IOException              Thrown if an error is occurred in input/output process
     * @throws ParseException           Thrown if an error is occured when parsing the content of the response to json
     *                                  object
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    private void setAccessToken(String email, String organizationKey, String password) throws IOException,
            ParseException, PluginExecutionException {

        StringEntity authorizationPayload = new StringEntity("scope=apim:api_create&grant_type=password&username=" +
                email + "@" + organizationKey + "&password=" + password);

        JSONObject accessTokenJson;

        try {
            log.info("Issuing REST call to Token API");
            response = getHttpResponse(TOKEN_API_URL, encodedIdAndSecret, authorizationPayload);

            if (response.getStatusLine().getStatusCode() == 401) {
                log.error("Error making the request to token API, the request is unauthorized");
                throw new PluginExecutionException("Unauthorized request, check email, organizationKey, password");
            } else if (response.getStatusLine().getStatusCode() == 400) {
                log.error("Error making the request to token API, malformed request body");
                throw new PluginExecutionException("Bad Request, check content");
            }

            content = new String(IOUtils.toByteArray(response.getEntity().getContent()), Charsets.UTF_8);
            log.info("Received the response");

            parser = new JSONParser();
            accessTokenJson = (JSONObject) parser.parse(content);

        } catch (IOException ioException) {
            log.error("Error while making the http call to token api", ioException);
            throw ioException;
        } catch (ParseException parseException) {
            log.error("Error while parsing the http response to json", parseException);
            throw parseException;
        } finally {
            httpClient.close();
        }

        log.info("Obtaining the access token");
        accessToken = accessTokenJson.get("access_token").toString();

    }

    /**
     * Creates an API in the api cloud and prints the response of the details of the API made
     *
     * @param email                     Email of the cloud account to export the API
     * @param organizationKey           The key generated in the API cloud for the given credentials, unique for the
     *                                  cloud
     * @param password                  Password of the cloud account to export the API
     * @param payload                   Payload for the api creation http request
     * @throws IOException              Thrown if an error is occurred in input/output process
     * @throws PluginExecutionException Custom exception to make the exception more readable
     * @throws ParseException           Thrown if an error is occurred when parsing the content of the response to
     *                                  json object
     */
    public void saveAPI(String email, String organizationKey, String password, String payload) throws IOException,
            PluginExecutionException, ParseException {

        setClientIdAndSecret(email, organizationKey, password);
        setAccessToken(email, organizationKey, password);

        StringEntity creationPayload = new StringEntity(payload);

        try {
            log.info("Creating the API in the cloud");
            response = getHttpResponse(API_CREATE_CLOUD_URL, accessToken, creationPayload);

            if (response.getStatusLine().getStatusCode() == 401) {
                log.error("Error while creating the API, the request is unauthorized");
                throw new PluginExecutionException("Unauthorized request");
            } else if (response.getStatusLine().getStatusCode() == 409) {
                log.error("Error while creating the API, the API already exists");
                throw new PluginExecutionException("An API with the same name and the context already exists");
            } else if (response.getStatusLine().getStatusCode() == 400) {
                log.error("Error creating the API, already exists with a different context");
                throw new PluginExecutionException("Bad content");
            }

            content = new String(IOUtils.toByteArray(response.getEntity().getContent()), Charsets.UTF_8);

            if (response.getStatusLine().getStatusCode() == 201) {
                log.info("The API is created in the cloud");
            } else {
                log.info("The API is not created in the cloud");
            }

        } catch (IOException ioException) {
            log.error("Error while making the request to create the API", ioException);
            throw ioException;
        } finally {
            httpClient.close();
        }
    }

    /**
     * Makes http requests to the specified URLs with suitable token prefix and content-type and returns the response
     *
     * @param url           URL which the request is made
     * @param token         Token for authoring the request
     * @param payload       Body of the request
     * @return              Returns the response of the http request made
     * @throws IOException  Thrown if an error is occurred in input/output process
     */
    private CloseableHttpResponse getHttpResponse(String url, String token, StringEntity payload) throws IOException {

        String tokenPrefix = "Basic ";
        String contentType = "application/json";
        if (url.equals(API_CREATE_CLOUD_URL)) {
            tokenPrefix = "Bearer ";
        }
        if (url.equals(TOKEN_API_URL)) {
            contentType = "application/x-www-form-urlencoded";
        }

        httpClient = HttpClients.createDefault();
        HttpPost apiRequest = new HttpPost(url);
        apiRequest.setHeader("Authorization", tokenPrefix + token);
        apiRequest.setHeader("Content-Type", contentType);
        apiRequest.setEntity(payload);

        response = httpClient.execute(apiRequest);
        return response;
    }
}

