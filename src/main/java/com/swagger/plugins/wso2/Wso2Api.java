package com.swagger.plugins.wso2;

import org.apache.commons.codec.binary.Base64;
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

    private HttpsURLConnection connection;
    private DataOutputStream dataOutputStream;
    private BufferedReader inDataStream;
    private StringBuffer response;
    private URL url;
    private JSONParser parser;

    /*
    * Method name : getClientIdAndSecret
    * Functionality : Obtains the client ID and client secret for the given email, organization key and password
    * @param : String, String, String
    * @return : String
    * */
    public String getClientIdAndSecret(String email, String organizationKey, String password) throws IOException {
        String encodeString = email + "@" + organizationKey + ":" + password;
        String authorizationPayload = "{\n" +
                "    \"callbackUrl\": \"www.google.lk\",\n" +
                "    \"clientName\": \"rest_api_publisher\",\n" +
                "    \"owner\":\""+email+"@"+organizationKey+"\",\n" +
                "    \"grantType\": \"password refresh_token\",\n" +
                "    \"saasApp\": true\n" +
                "}";
        try {
            url = new URL(DYNAMIC_CLIENT_REGISTRATION_URL);
            try {
                connection = (HttpsURLConnection) url.openConnection();
                connection.setRequestMethod("POST");
                connection.setRequestProperty("Content-Type", "application/json");
                connection.setRequestProperty("Authorization", "Basic " + Base64.encodeBase64String(encodeString.getBytes()));
                connection.setDoOutput(true);

                dataOutputStream = new DataOutputStream(connection.getOutputStream());
                dataOutputStream.writeBytes(authorizationPayload);

                inDataStream = new BufferedReader(new InputStreamReader(connection.getInputStream()));
                String line;
                response = new StringBuffer();
                while ((line = inDataStream.readLine()) != null)
                    response.append(line);
            } catch (IOException e) {
                e.printStackTrace();
            }
        } catch (MalformedURLException e) {
            e.printStackTrace();
        } catch (IOException e) {
            e.printStackTrace();
        } finally {
            dataOutputStream.flush();
            dataOutputStream.close();
            inDataStream.close();
        }
        return response.toString();
    }


    /*
    * Method name : getAccessToken
    * Functionality : Obtains the access token with the use of given client ID and client secret
    * @param : String, String, String
    * @return : String
    * */
    public String getAccessToken(String email, String organizationKey, String password) throws IOException {

        String clientIdAndSecret = getClientIdAndSecret(email, organizationKey, password);

        parser = new JSONParser();
        String clientId = null;
        String clientSecret = null;
        try {
            JSONObject clientIdAndSecretJson = (JSONObject) parser.parse(clientIdAndSecret);
            clientId = clientIdAndSecretJson.get("clientId").toString();
            clientSecret = clientIdAndSecretJson.get("clientSecret").toString();
        } catch (ParseException e) {
            e.printStackTrace();
        }

        String toEncode = clientId+":"+clientSecret;

        try {
            url = new URL(TOKEN_API_URL);
            connection = (HttpsURLConnection) url.openConnection();
            connection.setRequestMethod("POST");
            connection.setRequestProperty("Content-Type","application/x-www-form-urlencoded");
            connection.setRequestProperty("Authorization", "Basic "+ Base64.encodeBase64String(toEncode.getBytes()));
            connection.setDoOutput(true);

            dataOutputStream = new DataOutputStream(connection.getOutputStream());
            dataOutputStream.writeBytes("scope=apim:api_create&grant_type=password&username="+email+"@ms9714&password="+password);

            response = new StringBuffer();
            inDataStream = new BufferedReader(new InputStreamReader(connection.getInputStream()));
            String line;
            response = new StringBuffer();
            while ((line = inDataStream.readLine()) != null)
                response.append(line);
        } catch (IOException e) {
            e.printStackTrace();
        } finally {
            dataOutputStream.flush();
            dataOutputStream.close();
            inDataStream.close();
        }
        parser = new JSONParser();
        String accessToken = null;
        try {
            JSONObject accessTokenJson = (JSONObject) parser.parse(response.toString());
            accessToken = accessTokenJson.get("access_token").toString();
        } catch (ParseException e) {
            e.printStackTrace();
        }
        return accessToken;
    }





    /*
    * Method name : saveAPI
    * Functionality : Creates an API in the api cloud
    * @param : String, String
    * @return : void
    * */
    public void saveAPI(String swagger, String accessToken) {

        try {
            url = new URL(API_CREATE_CLOUD_URL);
            try {
                connection = (HttpsURLConnection) url.openConnection();
                connection.setRequestMethod("POST");
                connection.setRequestProperty("Authorization", "Bearer "+ accessToken);
                connection.setRequestProperty("Content-Type", "application/json");
                connection.setDoOutput(true);

                dataOutputStream = new DataOutputStream(connection.getOutputStream());
                dataOutputStream.writeBytes(swagger);


                inDataStream = new BufferedReader(new InputStreamReader(connection.getInputStream()));
                String line;
                StringBuffer response = new StringBuffer();
                while ((line = inDataStream.readLine()) != null)
                    response.append(line);



                System.out.println(response.toString());
                LOGGER.info("The api is created in the cloud successfully.");
                LOGGER.warn("The access token expires in an hour");
            } catch (IOException e) {
                e.printStackTrace();
            } finally {
                dataOutputStream.flush();
                dataOutputStream.close();
                inDataStream.close();
            }
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}

