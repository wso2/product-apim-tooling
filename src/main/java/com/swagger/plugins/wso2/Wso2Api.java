package com.swagger.plugins.wso2;

import org.apache.commons.codec.binary.Base64;
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

import java.io.*;


/*****************************************************************
 *  Class name : Wso2Api
 * Attributes : httpClient, inDataStream, inputStream, inputStreamReader, responseBody, response, parser
 * Constants : API_CREATE_CLOUD_URL, DYNAMIC_CLIENT_REGISTRATION_URL, TOKEN_API_URL
 * Methods : getAuthorizationPayload, getClientIdAndSecret, getAccessToken, saveAPI, getHttpResponse
 * Functionality : Contains the methods to obtain the access token and push the API to the cloud
 * Visibility : Public
 * ****************************************************************/
public class Wso2Api implements HttpConnectionService{

    private Logger LOGGER = LoggerFactory.getLogger(Wso2Api.class);

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
     * @param organizationKey   The key generated in the API cloud for the given credentials, unique for the cloud
     *                          account
     * @return
     */
    public String getAuthorizationPayload(String email, String organizationKey) {
        return "{\n" +
                "    \"callbackUrl\": \"www.google.lk\",\n" +
                "    \"clientName\": \"rest_api_publisher\",\n" +
                "    \"owner\":\""+email+"@"+organizationKey+"\",\n" +
                "    \"grantType\": \"password refresh_token\",\n" +
                "    \"saasApp\": true\n" +
                "}";
    }

    /**
     * Returns the Base64 encoded string of the format (clientId:clientSecret) for a given email organization key
     * and password.
     *
     * @param email             Email of the cloud account to export the API
     * @param organizationKey   The key generated in the API cloud for the given credentials, unique for the cloud
     *                          account
     * @param password          Password of the cloud account to export the API
     * @return Returns the Base64 encoded string of the format [clientId:clientSeret]
     * @throws ParseException
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    public void getClientIdAndSecret(String email, String organizationKey, String password) throws IOException,
            ParseException, PluginExecutionException {

        String stringToEncode = email + "@" + organizationKey + ":" + password;
        String encodedString = Base64.encodeBase64String(stringToEncode.getBytes("UTF-8"));
        StringEntity authorizationPayload = new StringEntity(getAuthorizationPayload(email, organizationKey));

        String clientId;
        String clientSecret;

        try {

            LOGGER.info("Calling dynamic client registration endpoint");

            response = getHttpResponse(DYNAMIC_CLIENT_REGISTRATION_URL, encodedString, authorizationPayload);

            if (response.getStatusLine().getStatusCode() == 401) {
                throw new PluginExecutionException("Unauthorized request, check email, organizationKey, password");
            } else if (response.getStatusLine().getStatusCode() == 400) {
                throw new PluginExecutionException("Bad Request, check content");
            }

            content = new String(IOUtils.toByteArray(response.getEntity().getContent()));
            LOGGER.info("Received the response");

            parser = new JSONParser();
            JSONObject clientIdAndSecretJson = (JSONObject) parser.parse(content);
            clientId = clientIdAndSecretJson.get("clientId").toString();
            clientSecret = clientIdAndSecretJson.get("clientSecret").toString();

        } catch (IOException e) {
            throw e;
        } catch (ParseException e) {
            throw e;
        }finally {
            httpClient.close();
        }

        String toEncode = clientId+":"+clientSecret;

        LOGGER.info("Obtaining the encoded clientId and clientSecret");
        encodedIdAndSecret = Base64.encodeBase64String(toEncode.getBytes("UTF-8"));
    }


    /**
     * Returns the access token obtained using the client Id and client secret
     *
     * @param email             Email of the cloud account to export the API
     * @param organizationKey   the key generated in the API cloud for the given credentials, unique for the cloud
     *                          account
     * @param password          Password of the cloud account to export the API
     * @return Returns the access token for creating an API in the cloud
     * @throws IOException
     * @throws ParseException
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    public void getAccessToken(String email, String organizationKey, String password) throws IOException,
            ParseException, PluginExecutionException {

        StringEntity authorizationPayload = new StringEntity("scope=apim:api_create&grant_type=password&username="+
                email+"@"+organizationKey+"&password="+password);

        JSONObject accessTokenJson;

        try {

            LOGGER.info("Issuing REST call to Token API");
            response = getHttpResponse(TOKEN_API_URL, encodedIdAndSecret, authorizationPayload);

            if (response.getStatusLine().getStatusCode() == 401) {
                throw new PluginExecutionException("Unauthorized request, check email, organizationKey, password");
            } else if (response.getStatusLine().getStatusCode() == 400) {
                throw new PluginExecutionException("Bad Request, check content");
            }

            content = new String(IOUtils.toByteArray(response.getEntity().getContent()));
            LOGGER.info("Received the response");

            parser = new JSONParser();
            accessTokenJson = (JSONObject) parser.parse(content);

        } catch (IOException e) {
            throw e;
        } catch (ParseException e) {
            throw e;
        }
        finally {
            httpClient.close();
        }

        LOGGER.info("Obtaining the access token");
        accessToken = accessTokenJson.get("access_token").toString();

    }

    /**
     * Creates an API in the api cloud and prints the response of the details of the API made
     *
     * @param payload       Payload for the api creation http request
     * @throws IOException
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    public void saveAPI(String payload) throws IOException, PluginExecutionException {

        StringEntity creationPayload = new StringEntity(payload);

        try {

            LOGGER.info("Creating the API in the cloud");
            response = getHttpResponse(API_CREATE_CLOUD_URL, accessToken, creationPayload);

            if (response.getStatusLine().getStatusCode() == 401) {
                throw new PluginExecutionException("Unauthorized request");
            } else if(response.getStatusLine().getStatusCode() == 409) {
                throw new PluginExecutionException("An API with the same name and the context already exists");
            } else if (response.getStatusLine().getStatusCode() == 400) {
                throw new PluginExecutionException("Bad content");
            }

            content = new String(IOUtils.toByteArray(response.getEntity().getContent()));

            if (response.getStatusLine().getStatusCode() == 201) {
                LOGGER.info("The API is created in the cloud");
            } else {
                LOGGER.info("The API is not created in the cloud");
            }

        } catch (IOException e) {
            throw e;
        } finally {
            httpClient.close();
        }

        System.out.println(content);
    }

    /**
     * Makes http requests to the specified URLs with suitable token prefix and content-type and returns the response
     *
     * @param url       URL which the request is made
     * @param token     Token for authoring the request
     * @param payload   Body of the request
     * @return Returns the response of the http request made
     * @throws IOException
     */
    public CloseableHttpResponse getHttpResponse(String url, String token, StringEntity payload) throws IOException {

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

