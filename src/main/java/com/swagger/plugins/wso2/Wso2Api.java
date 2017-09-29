package com.swagger.plugins.wso2;

import org.apache.commons.codec.binary.Base64;
import org.apache.http.HttpResponse;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.net.ssl.HttpsURLConnection;
import java.io.*;
import java.net.MalformedURLException;
import java.net.URL;


/*****************************************************************
 *  Class name : Wso2Api
 * Attributes : HttpsUrlConnection, DataOutPutStream, BufferReader, StringBuffer, URL, JSONParser
 * Methods : getClientIdAndSecret, getAccessToken, SaveAPI
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
    private StringBuffer responseBody;
    private HttpResponse response;
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
    * Functionality : Obtains the client ID and client secret for the given email, organization key and password
    * @param : String, String, String
    * @return : String
    * */
    public String getClientIdAndSecret(String email, String organizationKey, String password) throws IOException, ParseException {

        String encodeString = email + "@" + organizationKey + ":" + password;
        String encodedString = Base64.encodeBase64String(encodeString.getBytes());
        StringEntity authorizationPayload = new StringEntity(getAuthorizationPayload(email, organizationKey));

        httpClient = HttpClients.createDefault();
        HttpPost clientIdAndSecretRequest = new HttpPost(DYNAMIC_CLIENT_REGISTRATION_URL);
        clientIdAndSecretRequest.setHeader("Authorization","Basic "+encodedString);
        clientIdAndSecretRequest.setHeader("Content-Type","application/json");
        clientIdAndSecretRequest.setEntity(authorizationPayload);

        response = httpClient.execute(clientIdAndSecretRequest);

        //Use this space to handle exceptions

        inDataStream = new BufferedReader(new InputStreamReader(response.getEntity().getContent()));
        responseBody = new StringBuffer();
        String line;
        while ((line = inDataStream.readLine()) != null) {
            responseBody.append(line);
        }
        httpClient.close();

        parser = new JSONParser();
        String clientId = null;
        String clientSecret = null;

        JSONObject clientIdAndSecretJson = (JSONObject) parser.parse(responseBody.toString());
        clientId = clientIdAndSecretJson.get("clientId").toString();
        clientSecret = clientIdAndSecretJson.get("clientSecret").toString();

        String toEncode = clientId+":"+clientSecret;
        String encodedIdAndSecret = Base64.encodeBase64String(toEncode.getBytes());

        return encodedIdAndSecret;

    }


    /*
    * Method name : getAccessToken
    * Functionality : Obtains the access token with the use of given client ID and client secret
    * @param : String, String, String
    * @return : String
    * */
    public String getAccessToken(String email, String organizationKey, String password) throws IOException, ParseException {

        String clientIdAndSecret = getClientIdAndSecret(email, organizationKey, password);
        StringEntity authorizationPayload = new StringEntity("scope=apim:api_create&grant_type=password&username="+email+"@ms9714&password="+password);

        httpClient = HttpClients.createDefault();
        HttpPost accessTokenRequest = new HttpPost(TOKEN_API_URL);
        accessTokenRequest.setHeader("Authorization","Basic "+clientIdAndSecret);
        accessTokenRequest.setHeader("Content-Type", "application/x-www-form-urlencoded");
        accessTokenRequest.setEntity(authorizationPayload);

        response = httpClient.execute(accessTokenRequest);

        //Use this space to handle exceptions

        inDataStream = new BufferedReader(new InputStreamReader(response.getEntity().getContent()));
        responseBody = new StringBuffer();
        String line;
        while ((line = inDataStream.readLine()) != null) {
            responseBody.append(line);
        }
        httpClient.close();

        parser = new JSONParser();
        String accessToken;

        JSONObject accessTokenJson = (JSONObject) parser.parse(responseBody.toString());
        accessToken = accessTokenJson.get("access_token").toString();

        return accessToken;
    }





    /*
    * Method name : saveAPI
    * Functionality : Creates an API in the api cloud
    * @param : String, String
    * @return : void
    * */
    public void saveAPI(String swagger, String accessToken) throws IOException {

        StringEntity creationPayload = new StringEntity(swagger);

        httpClient = HttpClients.createDefault();
        HttpPost createApiRequest = new HttpPost(API_CREATE_CLOUD_URL);
        createApiRequest.setHeader("Authorization","Bearer "+accessToken);
        createApiRequest.setHeader("Content-Type","application/json");
        createApiRequest.setEntity(creationPayload);

        response = httpClient.execute(createApiRequest);

        //Use this space to handle exceptions

        inDataStream = new BufferedReader(new InputStreamReader(response.getEntity().getContent()));
        responseBody = new StringBuffer();
        String line;
        while ((line = inDataStream.readLine()) != null) {
            responseBody.append(line);
        }
        httpClient.close();

        System.out.println(responseBody.toString());
    }
}

