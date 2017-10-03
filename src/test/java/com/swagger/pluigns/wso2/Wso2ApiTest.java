package com.swagger.pluigns.wso2;

import com.swagger.plugins.wso2.Wso2Api;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.CloseableHttpClient;
import org.junit.Test;
import org.junit.Assert;
import org.mockito.Mockito;

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
        HttpPost clientIdAndSecretRequest = new HttpPost("https://api.cloud.wso2.com/client-registration/v0.11/register");
        clientIdAndSecretRequest.setHeader("Authorization","Basic YWtpbGFhcm9zaGFuYUBob3RtYWlsLmNvbUBtczk3MTQ6QWtpbGFob3RtYWlsMTIz");
        clientIdAndSecretRequest.setHeader("Content-Type","application/json");
        clientIdAndSecretRequest.setEntity(authorizationPayload);
        Mockito.when(httpClient.execute(clientIdAndSecretRequest)).thenReturn(response);

        String token = api.getClientIdAndSecret("akilaaroshana@hotmail.com", "ms9714", "Akilahotmail123");
        Assert.assertEquals("UE03azdsUWZaR05LTXlxdHp0VHJmczRFbEtvYTpxcGRjQlF1X3pVTkhmMmZwSGVjTjBydThocU1h", token);
    }
}
