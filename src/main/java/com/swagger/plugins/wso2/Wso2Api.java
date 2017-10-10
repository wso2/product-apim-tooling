package com.swagger.plugins.wso2;

import org.apache.commons.codec.binary.Base64;
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
 * Methods : getAuthorizationPayload, getClientIdAndSecret, getAccessToken, SaveAPI
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
    private BufferedReader inDataStream;
    private InputStream inputStream;
    private InputStreamReader inputStreamReader;
    private StringBuffer responseBody;
    private CloseableHttpResponse response;
    private JSONParser parser;

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
     * @param email Email of the cloud account to export the API
     * @param organizationKey the key generated in the API cloud for the given credentials, unique for the cloud account
     * @param password Password of the cloud account to export the API
     * @return Returns the Base64 encoded string of the format [clientId:clientSeret]
     * @throws ParseException
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    public String getClientIdAndSecret(String email, String organizationKey, String password) throws IOException,
            ParseException, PluginExecutionException {

        String encodeString = email + "@" + organizationKey + ":" + password;
        String encodedString = Base64.encodeBase64String(encodeString.getBytes("UTF-8"));
        StringEntity authorizationPayload = new StringEntity(getAuthorizationPayload(email, organizationKey));

        String clientId;
        String clientSecret;

        try {

            response = getHttpResponse(DYNAMIC_CLIENT_REGISTRATION_URL, encodedString, authorizationPayload);

            LOGGER.info("Calling dynamic client registration endpoint");

            if (response.getStatusLine().getStatusCode() == 401) {
                throw new PluginExecutionException("Unauthorized request, check email, organizationKey, password");
            } else if (response.getStatusLine().getStatusCode() == 400) {
                throw new PluginExecutionException("Bad Request, check content");
            }

            inputStream = response.getEntity().getContent();
            inputStreamReader = new InputStreamReader(inputStream, "UTF-8");
            inDataStream = new BufferedReader(inputStreamReader);
            responseBody = new StringBuffer();
            String line;
            while ((line = inDataStream.readLine()) != null) {
                responseBody.append(line);
            }

            LOGGER.info("Received the payload");

            parser = new JSONParser();
            JSONObject clientIdAndSecretJson = (JSONObject) parser.parse(responseBody.toString());
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
        String encodedIdAndSecret = Base64.encodeBase64String(toEncode.getBytes("UTF-8"));

        LOGGER.info("Returning the encoded clientId and clientSecret");
        return encodedIdAndSecret;
    }


    /**
     * Returns the access token obtained using the client Id and client secret
     * @param email Email of the cloud account to export the API
     * @param organizationKey the key generated in the API cloud for the given credentials, unique for the cloud account
     * @param password Password of the cloud account to export the API
     * @return Returns the access token for creating an API in the cloud
     * @throws IOException
     * @throws ParseException
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    public String getAccessToken(String email, String organizationKey, String password) throws IOException,
            ParseException, PluginExecutionException {

        String clientIdAndSecret = getClientIdAndSecret(email, organizationKey, password);
        StringEntity authorizationPayload = new StringEntity("scope=apim:api_create&grant_type=password&username="+
                email+"@"+organizationKey+"&password="+password);

        JSONObject accessTokenJson;
        String accessToken;

        try {

            response = getHttpResponse(TOKEN_API_URL, clientIdAndSecret, authorizationPayload);

            LOGGER.info("Issuing REST call to Token API");

            if (response.getStatusLine().getStatusCode() == 401) {
                throw new PluginExecutionException("Unauthorized request, check email, organizationKey, password");
            } else if (response.getStatusLine().getStatusCode() == 400) {
                throw new PluginExecutionException("Bad Request, check content");
            }

            inDataStream = new BufferedReader(new InputStreamReader(response.getEntity().getContent(), "UTF-8"));
            responseBody = new StringBuffer();
            String line;
            while ((line = inDataStream.readLine()) != null) {
                responseBody.append(line);
            }

            LOGGER.info("Received response");

            parser = new JSONParser();
            accessTokenJson = (JSONObject) parser.parse(responseBody.toString());

        } catch (IOException e) {
            throw e;
        } catch (ParseException e) {
            throw e;
        }
        finally {
            httpClient.close();
        }

        accessToken = accessTokenJson.get("access_token").toString();
        LOGGER.info("Returning the access token");
        return accessToken;
    }

    /**
     * Creates an API in the api cloud and prints the response of the details of the API made
     * @param payload Payload for the api creation http request
     * @param accessToken Access token for creating the API in the cloud
     * @throws IOException
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    public void saveAPI(String payload, String accessToken) throws IOException, PluginExecutionException {

        StringEntity creationPayload = new StringEntity(payload);

        try {

            response = getHttpResponse(API_CREATE_CLOUD_URL, accessToken, creationPayload);

            LOGGER.info("Creating the API in the cloud");

            if (response.getStatusLine().getStatusCode() == 401) {
                throw new PluginExecutionException("Unauthorized request");
            } else if(response.getStatusLine().getStatusCode() == 409) {
                throw new PluginExecutionException("An API with the same name and the context already exists");
            } else if (response.getStatusLine().getStatusCode() == 400) {
                throw new PluginExecutionException("Bad content");
            }

            inDataStream = new BufferedReader(new InputStreamReader(response.getEntity().getContent(), "UTF-8"));
            responseBody = new StringBuffer();
            String line;
            while ((line = inDataStream.readLine()) != null) {
                responseBody.append(line);
            }

            if (response.getStatusLine().getStatusCode() == 201) {
                LOGGER.info("The API is created in the cloud");
            }

        } catch (IOException e) {
            throw e;
        } finally {
            httpClient.close();
        }

        System.out.println(responseBody.toString());
    }

    /**
     * Makes http requests to the specified URLs with suitable token prefix and content-type and returns the response
     * @param url URL which the request is made
     * @param token Token for authoring the request
     * @param payload Body of the request
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

