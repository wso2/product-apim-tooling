package com.swagger.pluigns.wso2;

import com.swagger.plugins.wso2.Main;
import com.swagger.plugins.wso2.PayloadConfiguration;
import com.swagger.plugins.wso2.PluginExecutionException;
import com.swagger.plugins.wso2.Wso2Api;
import org.json.simple.parser.ParseException;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.CloseableHttpClient;
import org.junit.Test;
import org.junit.Assert;
import org.mockito.Mockito;

import java.io.IOException;

public class Wso2ApiTest {

//    @Test
//    public void getClientIdAndSecretTestHappyPath() throws Exception {
//        Wso2Api api = new Wso2Api();
//        StringEntity authorizationPayload = new StringEntity("{\n" +
//                "    \"callbackUrl\": \"www.google.lk\",\n" +
//                "    \"clientName\": \"rest_api_publisher\",\n" +
//                "    \"owner\":\"akilaaroshana@hotmail.com@ms9714\",\n" +
//                "    \"grantType\": \"password refresh_token\",\n" +
//                "    \"saasApp\": true\n" +
//                "}");
//
//        CloseableHttpResponse response = Mockito.mock(CloseableHttpResponse.class);
//        CloseableHttpClient httpClient = Mockito.mock(CloseableHttpClient.class);
//        HttpPost clientIdAndSecretRequest = new HttpPost("https://api.cloud.wso2.com/client-registration/v0.11/" +
//                "register");
//        clientIdAndSecretRequest.setHeader("Authorization","Basic " +
//                "YWtpbGFhcm9zaGFuYUBob3RtYWlsLmNvbUBtczk3MTQ6QWtpbGFob3RtYWlsMTIz");
//        clientIdAndSecretRequest.setHeader("Content-Type","application/json");
//        clientIdAndSecretRequest.setEntity(authorizationPayload);
//        Mockito.when(httpClient.execute(clientIdAndSecretRequest)).thenReturn(response);
//
//        String token = api.getClientIdAndSecret("yolom@seekjobs4u.com", "yolo4958",
//                "Yolofernando123");
//        System.out.println(token);
//        Assert.assertEquals("WEtCZWVwd213VF81SHZlaXNDd3BEVkRKa3dFYTpPRjBUaE5BUFF1R0s1VnVJS2ZQRGhUYWZ5a2Nh",
//                token);
//    }
//
//    @Test
//    public void getAccessToken() throws ParseException, PluginExecutionException, IOException {
//        Wso2Api api = Mockito.mock(Wso2Api.class);
//        Mockito.when(api.getClientIdAndSecret("yolom@seekjobs4u.com", "yolo4958",
//                "Yolofernando123")).thenReturn("WEtCZWVwd213VF81SHZlaXNDd3BEVkRKa3dFYT" +
//                "pPRjBUaE5BUFF1R0s1VnVJS2ZQRGhUYWZ5a2Nh");
//        StringEntity authorizationPayload = new StringEntity("scope=apim:api_create&grant_type=" +
//                "password&username=yolom@seekjobs4u.com@yolo4958&password=Yolofernando123");
//
//        CloseableHttpClient httpClient = Mockito.mock(CloseableHttpClient.class);
//        CloseableHttpResponse response = Mockito.mock(CloseableHttpResponse.class);
//        HttpPost accessTokenRequest = new HttpPost("https://gateway.api.cloud.wso2.com/token");
//        accessTokenRequest.setHeader("Authorization","Basic WEtCZWVwd213VF81SHZlaXNDd3BEVkRKa3dFYT" +
//                "pPRjBUaE5BUFF1R0s1VnVJS2ZQRGhUYWZ5a2Nh");
//        accessTokenRequest.setHeader("Content-Type", "application/x-www-form-urlencoded");
//        accessTokenRequest.setEntity(authorizationPayload);
//
//        Mockito.when(httpClient.execute(accessTokenRequest)).thenReturn(response);
//
//        Wso2Api wso2Api = new Wso2Api();
//        String returnedAccessToken = wso2Api.getAccessToken("yolom@seekjobs4u.com", "yolo4958",
//                "Yolofernando123");
//
//        //Since there is no way of knowing the access token before it is issued, just tested for obtaining the access
//        //token
//        if (returnedAccessToken != null) {
//            Assert.assertTrue(true);
//        }
//    }
//
//    @Test
//    public void saveApiTest() throws ParseException, PluginExecutionException, IOException {
//
//        String yaml = Main.fileReader("/home/akila/Akila/swagger-wso2-integration-local/src/test/" +
//                "java/com/swagger/pluigns/wso2/SwaggerSimpleApi.yaml");
//
//        Wso2Api wso2Api = Mockito.mock(Wso2Api.class);
//        String accessToken = wso2Api.getAccessToken("yolom@seekjobs4u.com", "yolo4958",
//                "Yolofernando123");
//        System.out.println(accessToken);
//        PayloadConfiguration configuration = new PayloadConfiguration();
//
//        String payload = configuration.configurePayload("yolom@seekjobs4u.com", "yolo4958", yaml, "sample");
//
//        StringEntity stringEntity = new StringEntity(payload);
//
//
//        CloseableHttpClient httpClient = Mockito.mock(CloseableHttpClient.class);
//        HttpPost createApiRequest = new HttpPost("https://api.cloud.wso2.com/api/am/publisher/v0.11/apis");
//        createApiRequest.setHeader("Authorization","Bearer "+accessToken);
//        createApiRequest.setHeader("Content-Type","application/json");
//        createApiRequest.setEntity(stringEntity);
//
//        CloseableHttpResponse response = Mockito.mock(CloseableHttpResponse.class);
//        Mockito.when(httpClient.execute(createApiRequest)).thenReturn(response);
//
//
//        Wso2Api api = new Wso2Api();
//        api.saveAPI(payload, accessToken);
//
//        if (response.getStatusLine().getStatusCode() == 400) {
//            Assert.assertTrue(true);
//        }
//
//
//    }
    //Test for apis with same name but with different versions.
}
