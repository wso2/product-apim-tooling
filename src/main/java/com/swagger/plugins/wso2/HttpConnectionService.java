package com.swagger.plugins.wso2;

import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.entity.StringEntity;

import java.io.IOException;

public interface HttpConnectionService {

    CloseableHttpResponse getHttpResponse(String url, String token, StringEntity payload) throws IOException;

}
