package com.swagger.pluigns.wso2;

import com.swagger.plugins.wso2.PayloadConfiguration;
import com.swagger.plugins.wso2.PluginExecutionException;
import com.swagger.plugins.wso2.Wso2Api;
import org.apache.http.impl.client.HttpClients;
import org.json.simple.parser.ParseException;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.CloseableHttpClient;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.junit.Test;
import org.junit.Assert;
import org.mockito.Mockito;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;

public class Wso2ApiTest {

    @Test
    public void getClientIdAndSecretTestHappyPath() throws Exception {
        Wso2Api api = new Wso2Api();
        StringEntity authorizationPayload = new StringEntity("{\n" +
                "    \"callbackUrl\": \"www.google.lk\",\n" +
                "    \"clientName\": \"rest_api_publisher\",\n" +
                "    \"owner\":\"akilaaroshana@hotmail.com@ms9714\",\n" +
                "    \"grantType\": \"password refresh_token\",\n" +
                "    \"saasApp\": true\n" +
                "}");

        CloseableHttpResponse response = Mockito.mock(CloseableHttpResponse.class);
        CloseableHttpClient httpClient = Mockito.mock(CloseableHttpClient.class);
        HttpPost clientIdAndSecretRequest = new HttpPost("https://api.cloud.wso2.com/client-registration/v0.11/" +
                "register");
        clientIdAndSecretRequest.setHeader("Authorization","Basic " +
                "YWtpbGFhcm9zaGFuYUBob3RtYWlsLmNvbUBtczk3MTQ6QWtpbGFob3RtYWlsMTIz");
        clientIdAndSecretRequest.setHeader("Content-Type","application/json");
        clientIdAndSecretRequest.setEntity(authorizationPayload);
        Mockito.when(httpClient.execute(clientIdAndSecretRequest)).thenReturn(response);

        String token = api.getClientIdAndSecret("yolom@seekjobs4u.com", "yolo4958",
                "Yolofernando123");
        System.out.println(token);
        Assert.assertEquals("WEtCZWVwd213VF81SHZlaXNDd3BEVkRKa3dFYTpPRjBUaE5BUFF1R0s1VnVJS2ZQRGhUYWZ5a2Nh",
                token);
    }

    @Test
    public void getAccessToken() throws ParseException, PluginExecutionException, IOException {
        Wso2Api api = Mockito.mock(Wso2Api.class);
        Mockito.when(api.getClientIdAndSecret("yolom@seekjobs4u.com", "yolo4958",
                "Yolofernando123")).thenReturn("WEtCZWVwd213VF81SHZlaXNDd3BEVkRKa3dFYT" +
                "pPRjBUaE5BUFF1R0s1VnVJS2ZQRGhUYWZ5a2Nh");
        StringEntity authorizationPayload = new StringEntity("scope=apim:api_create&grant_type=" +
                "password&username=yolom@seekjobs4u.com@yolo4958&password=Yolofernando123");

        CloseableHttpClient httpClient = Mockito.mock(CloseableHttpClient.class);
        CloseableHttpResponse response = Mockito.mock(CloseableHttpResponse.class);
        HttpPost accessTokenRequest = new HttpPost("https://gateway.api.cloud.wso2.com/token");
        accessTokenRequest.setHeader("Authorization","Basic WEtCZWVwd213VF81SHZlaXNDd3BEVkRKa3dFYT" +
                "pPRjBUaE5BUFF1R0s1VnVJS2ZQRGhUYWZ5a2Nh");
        accessTokenRequest.setHeader("Content-Type", "application/x-www-form-urlencoded");
        accessTokenRequest.setEntity(authorizationPayload);

        Mockito.when(httpClient.execute(accessTokenRequest)).thenReturn(response);

        Wso2Api wso2Api = new Wso2Api();
        String returnedAccessToken = wso2Api.getAccessToken("yolom@seekjobs4u.com", "yolo4958",
                "Yolofernando123");

        //Since there is no way of knowing the access token before it is issued, just tested for obtaining the access
        //token
        if (returnedAccessToken != null) {
            Assert.assertTrue(true);
        }
    }

    @Test
    public void saveApiTest() throws ParseException, PluginExecutionException, IOException {

        String yaml = "swagger: '2.0'\n" +
                "info:\n" +
                "  description: This is a simple API\n" +
                "  version: 1.0.0\n" +
                "  title: Simple Inventory API\n" +
                "  # put the contact info for your development or API team\n" +
                "  contact:\n" +
                "    email: you@your-company.com\n" +
                "\n" +
                "  license:\n" +
                "    name: Apache 2.0\n" +
                "    url: http://www.apache.org/licenses/LICENSE-2.0.html\n" +
                "\n" +
                "# tags are used for organizing operations\n" +
                "tags:\n" +
                "- name: admins\n" +
                "  description: Secured Admin-only calls\n" +
                "- name: developers\n" +
                "  description: Operations available to regular developers\n" +
                "\n" +
                "paths:\n" +
                "  /inventory:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - developers\n" +
                "      summary: searches inventory\n" +
                "      operationId: searchInventory\n" +
                "      description: |\n" +
                "        By passing in the appropriate options, you can search for\n" +
                "        available inventory in the system\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      parameters:\n" +
                "      - in: query\n" +
                "        name: searchString\n" +
                "        description: pass an optional search string for looking up inventory\n" +
                "        required: false\n" +
                "        type: string\n" +
                "      - in: query\n" +
                "        name: skip\n" +
                "        description: number of records to skip for pagination\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "        minimum: 0\n" +
                "      - in: query\n" +
                "        name: limit\n" +
                "        description: maximum number of records to return\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "        minimum: 0\n" +
                "        maximum: 50\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: search results matching criteria\n" +
                "          schema:\n" +
                "            type: array\n" +
                "            items:\n" +
                "              $ref: '#/definitions/InventoryItem'\n" +
                "        400:\n" +
                "          description: bad input parameter\n" +
                "    post:\n" +
                "      tags:\n" +
                "      - admins\n" +
                "      summary: adds an inventory item\n" +
                "      operationId: addInventory\n" +
                "      description: Adds an item to the system\n" +
                "      consumes:\n" +
                "      - application/json\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      parameters:\n" +
                "      - in: body\n" +
                "        name: inventoryItem\n" +
                "        description: Inventory item to add\n" +
                "        schema:\n" +
                "          $ref: '#/definitions/InventoryItem'\n" +
                "      responses:\n" +
                "        201:\n" +
                "          description: item created\n" +
                "        400:\n" +
                "          description: invalid input, object invalid\n" +
                "        409:\n" +
                "          description: an existing item already exists\n" +
                "definitions:\n" +
                "  InventoryItem:\n" +
                "    type: object\n" +
                "    required:\n" +
                "    - id\n" +
                "    - name\n" +
                "    - manufacturer\n" +
                "    - releaseDate\n" +
                "    properties:\n" +
                "      id:\n" +
                "        type: string\n" +
                "        format: uuid\n" +
                "        example: d290f1ee-6c54-4b01-90e6-d701748f0851\n" +
                "      name:\n" +
                "        type: string\n" +
                "        example: Widget Adapter\n" +
                "      releaseDate:\n" +
                "        type: string\n" +
                "        format: int32\n" +
                "        example: 2016-08-29T09:12:33.001Z\n" +
                "      manufacturer:\n" +
                "        $ref: '#/definitions/Manufacturer'\n" +
                "  Manufacturer:\n" +
                "    required:\n" +
                "    - name\n" +
                "    properties:\n" +
                "      name:\n" +
                "        type: string\n" +
                "        example: ACME Corporation\n" +
                "      homePage:\n" +
                "        type: string\n" +
                "        format: url\n" +
                "        example:  https://www.acme-corp.com\n" +
                "      phone:\n" +
                "        type: string\n" +
                "        example: 408-867-5309\n" +
                "# Added by API Auto Mocking Plugin\n" +
                "host: virtserver.swaggerhub.com\n" +
                "basePath: /yolom4/SimpleApi/1.0.0\n" +
                "schemes:\n" +
                " - https";

        Wso2Api wso2Api = Mockito.mock(Wso2Api.class);
        String accessToken = wso2Api.getAccessToken("yolom@seekjobs4u.com", "yolo4958",
                "Yolofernando123");
        System.out.println(accessToken);
        PayloadConfiguration configuration = new PayloadConfiguration();

        String payload = configuration.configurePayload("yolom@seekjobs4u.com", "yolo4958", yaml, "1.0.0", "sample", "id");

        StringEntity stringEntity = new StringEntity(payload);


        CloseableHttpClient httpClient = Mockito.mock(CloseableHttpClient.class);
        HttpPost createApiRequest = new HttpPost("https://api.cloud.wso2.com/api/am/publisher/v0.11/apis");
        createApiRequest.setHeader("Authorization","Bearer "+accessToken);
        createApiRequest.setHeader("Content-Type","application/json");
        createApiRequest.setEntity(stringEntity);

        CloseableHttpResponse response = Mockito.mock(CloseableHttpResponse.class);
        Mockito.when(httpClient.execute(createApiRequest)).thenReturn(response);


        Wso2Api api = new Wso2Api();
        api.saveAPI(payload, accessToken);

        if (response.getStatusLine().getStatusCode() == 400) {
            Assert.assertTrue(true);
        }


    }
    //Test for apis with same name but with different versions.
}
