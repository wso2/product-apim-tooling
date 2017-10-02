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
public class Wso2Api {

    private Logger LOGGER = LoggerFactory.getLogger(Wso2Api.class);

    private static final String API_CREATE_CLOUD_URL = "https://api.cloud.wso2.com/api/am/publisher/v0.11/apis";
    private static final String DYNAMIC_CLIENT_REGISTRATION_URL = "https://api.cloud.wso2.com/client-registration/v0.11/register";
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

    /*
    * Method name : getClientIdAndSecret
    * Functionality : Returns the Base64 encoded string of the format (clientId:clientSecret) for a given email, organization key and password.
    * @param : String, String, String
    * @return : String
    * */
    public String getClientIdAndSecret(String email, String organizationKey, String password) throws IOException, ParseException, PluginExecutionException {

        String encodeString = email + "@" + organizationKey + ":" + password;
        String encodedString = Base64.encodeBase64String(encodeString.getBytes("UTF-8"));
        StringEntity authorizationPayload = new StringEntity(getAuthorizationPayload(email, organizationKey));

        String clientId;
        String clientSecret;

        try {
            httpClient = HttpClients.createDefault();
            HttpPost clientIdAndSecretRequest = new HttpPost(DYNAMIC_CLIENT_REGISTRATION_URL);
            clientIdAndSecretRequest.setHeader("Authorization","Basic "+encodedString);
            clientIdAndSecretRequest.setHeader("Content-Type","application/json");
            clientIdAndSecretRequest.setEntity(authorizationPayload);

            LOGGER.info("Calling dynamic client registration endpoint");
            response = httpClient.execute(clientIdAndSecretRequest);


            if (response.getStatusLine().getStatusCode() == 401) {
                throw new PluginExecutionException("Unauthorized request, check email, organizationKey, password");
            } else if (response.getStatusLine().getStatusCode() == 400) {
                throw new PluginExecutionException("Bad Request, check content");
            }

            inputStream = response.getEntity().getContent();
            inputStreamReader = new InputStreamReader(inputStream);
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

        System.out.println("Response code of :");
        System.out.println(response.getStatusLine().getStatusCode());
        LOGGER.info("Returning the encoded clientId and clientSecret");
        return encodedIdAndSecret;
    }


    /*
    * Method name : getAccessToken
    * Functionality : Returns the access token obtained using the client Id and client secret.
    * @param : String, String, String
    * @return : String
    * */
    public String getAccessToken(String email, String organizationKey, String password) throws IOException, ParseException, PluginExecutionException {

        String clientIdAndSecret = getClientIdAndSecret(email, organizationKey, password);
        StringEntity authorizationPayload = new StringEntity("scope=apim:api_create&grant_type=password&username="+email+"@ms9714&password="+password);

        JSONObject accessTokenJson;
        String accessToken;

        try {
            httpClient = HttpClients.createDefault();
            HttpPost accessTokenRequest = new HttpPost(TOKEN_API_URL);
            accessTokenRequest.setHeader("Authorization","Basic "+clientIdAndSecret);
            accessTokenRequest.setHeader("Content-Type", "application/x-www-form-urlencoded");
            accessTokenRequest.setEntity(authorizationPayload);

            LOGGER.info("Issuing REST call to Token API");
            response = httpClient.execute(accessTokenRequest);

            if (response.getStatusLine().getStatusCode() == 401) {
                throw new PluginExecutionException("Unauthorized request, check email, organizationKey, password");
            } else if (response.getStatusLine().getStatusCode() == 400) {
                throw new PluginExecutionException("Bad Request, check content");
            }

            inDataStream = new BufferedReader(new InputStreamReader(response.getEntity().getContent()));
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
        System.out.println("Response code :");
        System.out.println(response.getStatusLine().getStatusCode());
        LOGGER.info("Returning the access token");
        return accessToken;
    }

    /*
    * Method name : saveAPI
    * Functionality : Creates an API in the api cloud and prints the response of the details of the API made.
    * @param : String, String
    * @return : void
    * */
    public void saveAPI(String swagger, String accessToken) throws IOException, PluginExecutionException {

        StringEntity creationPayload = new StringEntity(swagger);

        try {
            httpClient = HttpClients.createDefault();
            HttpPost createApiRequest = new HttpPost(API_CREATE_CLOUD_URL);
            createApiRequest.setHeader("Authorization","Bearer "+accessToken);
            createApiRequest.setHeader("Content-Type","application/json");
            createApiRequest.setEntity(creationPayload);

            LOGGER.info("Creating the API in the cloud");
            response = httpClient.execute(createApiRequest);

            if (response.getStatusLine().getStatusCode() == 401) {
                throw new PluginExecutionException("Unauthorized request");
            } else if(response.getStatusLine().getStatusCode() == 409) {
                throw new PluginExecutionException("An API with the same name and the context already exists");
            }

            inDataStream = new BufferedReader(new InputStreamReader(response.getEntity().getContent()));
            responseBody = new StringBuffer();
            String line;
            while ((line = inDataStream.readLine()) != null) {
                responseBody.append(line);
            }

            LOGGER.info("The API is created in the cloud");

        } catch (IOException e) {
            throw e;
        } finally {
            httpClient.close();
        }

        System.out.println("Response code :");
        System.out.println(response.getStatusLine().getStatusCode());
        System.out.println(responseBody.toString());
    }
}

