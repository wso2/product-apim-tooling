package com.swagger.plugins.wso2;

import io.swagger.models.Swagger;
import io.swagger.util.Json;
import org.apache.commons.codec.binary.Base64;
import org.apache.commons.io.Charsets;
import org.apache.commons.io.IOUtils;
import org.apache.http.HttpResponse;
import org.apache.http.entity.StringEntity;
import org.json.simple.JSONArray;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.io.UnsupportedEncodingException;


/*****************************************************************
 * Class name : WSO2Api
 * Attributes : httpClient, content, response, parser, encodedIdAndSecret, accessToken
 * Constants : API_CLOUD_URL, DYNAMIC_CLIENT_REGISTRATION_URL, TOKEN_API_URL
 * Methods : getAuthorizationPayloadForClientIdAndSecret, getAuthorizationPayloadForTokenApi, setClientIdAndSecret,
 * getAccessToken, getApiIdentifier, updateApi, saveAPI
 * Functionality : Contains the methods to obtain the access token and push the API to the cloud and update the API
 * Visibility : Public
 * ****************************************************************/
public class WSO2Api {

    private Logger log = LoggerFactory.getLogger(WSO2Api.class);

    private static final String API_CLOUD_URL = "https://api.cloud.wso2.com/api/am/publisher/v0.11/apis/";
    private static final String DYNAMIC_CLIENT_REGISTRATION_URL = "https://api.cloud.wso2.com/client-registration/" +
                                                                  "v0.11/register";
    private static final String TOKEN_API_URL = "https://gateway.api.cloud.wso2.com/token";

    private HttpRequestService httpRequestService;
    private String encodedIdAndSecret;


    /**
     * Instantiates an instance of HttpRequestService interface with the instantiation of WSO2Api class
     *
     * @param httpRequestService Interface with method signatures to be implemented
     */
    public WSO2Api(HttpRequestService httpRequestService) {
        this.httpRequestService = httpRequestService;
    }

    /**
     * Returns the payload for request of getting client id and secret.
     *
     * @param email             Email of the cloud account to export the API
     * @param organizationKey   The key generated in the API cloud for the given credentials, unique for the WSO2 API
     *                          cloud account
     * @return                  Returns the authorization payload for obtaining client id and client secret
     */
    private String getAuthorizationPayloadForClientIdAndSecret(String email, String organizationKey) {
        return "{\n" +
                "    \"clientName\": \"rest_api_publisher\",\n" +
                "    \"owner\":\"" + email + "@" + organizationKey + "\",\n" +
                "    \"grantType\": \"password refresh_token\",\n" +
                "    \"saasApp\": true\n" +
                "}";
    }

    /**
     *
     * @param email Email of the cloud account to export the API
     * @param organizationKey The key generated in the API cloud for the given credentials, unique for the
     *                                  cloud
     * @param password Password of the cloud account to export the API
     * @param scope Scope of the token to be obtained
     * @return Returns the payload for requesting the access token
     */
    private String getAuthorizationPayloadForTokenApi(String email, String organizationKey, String password,
                                                       String scope) {
        return "scope=apim:api_" + scope + "&grant_type=password&username=" +
                email + "@" + organizationKey + "&password=" + password;

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
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    private void setClientIdAndSecret(String email, String organizationKey, String password) throws
            PluginExecutionException {

        HttpResponse response;
        StringEntity authorizationPayload;
        String stringToEncode = email + "@" + organizationKey + ":" + password;
        String encodedString = Base64.encodeBase64String(stringToEncode.getBytes(Charsets.UTF_8));

        try {
            authorizationPayload = new StringEntity(getAuthorizationPayloadForClientIdAndSecret(email,
                    organizationKey));
        } catch (UnsupportedEncodingException unsupportedEncodingException) {
            log.error("The character encoding is not supported for the payload", unsupportedEncodingException);
            throw new PluginExecutionException("The character encoding is not supported");
        }

        String clientId;
        String clientSecret;

        try {
            log.debug("Calling dynamic client registration endpoint");

            response = httpRequestService.makePostRequest(DYNAMIC_CLIENT_REGISTRATION_URL,
                    "Basic", encodedString, "application/json", authorizationPayload);

            if (response.getStatusLine().getStatusCode() == 401) {
                log.error("Error making the request to dynamic client registration endpoint, the request" +
                        " is unauthorized");
                throw new PluginExecutionException("Unauthorized request, check email, organizationKey, password");
            } else if (response.getStatusLine().getStatusCode() == 400) {
                log.error("Error making the request to dynamic client registration endpoint, malformed request body");
                throw new PluginExecutionException("Bad Request, check content");
            } else {
                if (response.getStatusLine().getStatusCode() == 200) {
                    log.debug("Client id and client secret request successful");
                } else {
                    log.error("Client id and client secret request unsuccessful");
                    throw new PluginExecutionException("Requesting client id and client secret was unsuccessful");
                }
            }

            String content = new String(IOUtils.toByteArray(response.getEntity().getContent()), Charsets.UTF_8);
            log.debug("Received the response");

            JSONParser parser = new JSONParser();
            JSONObject clientIdAndSecretJson = (JSONObject) parser.parse(content);
            clientId = clientIdAndSecretJson.get("clientId").toString();
            clientSecret = clientIdAndSecretJson.get("clientSecret").toString();

        } catch (IOException ioException) {
            log.error("Error while making the http call to dynamic client registration endpoint", ioException);
            throw new PluginExecutionException("Error while making http call to dynamic client registration endpoint");
        } catch (ParseException parseException) {
            log.error("Error while parsing the http response to json", parseException);
            throw new PluginExecutionException("Error parsing the response to json");
        }

        String toEncode = clientId + ":" + clientSecret;

        log.debug("Obtaining the encoded clientId and clientSecret");
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
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    private String getAccessToken(String email, String organizationKey, String password, String scope) throws
            PluginExecutionException {

        HttpResponse response;
        String content;
        StringEntity authorizationPayload;
        JSONObject accessTokenJson;
        setClientIdAndSecret(email, organizationKey, password);

        try {
            authorizationPayload = new StringEntity(getAuthorizationPayloadForTokenApi(email, organizationKey,
                    password, scope));
        } catch (UnsupportedEncodingException unsupportedEncodingException) {
            log.error("The character encoding is not supported for the payload", unsupportedEncodingException);
            throw new PluginExecutionException("The character encoding is not supported");
        }

        try {
            log.debug("Issuing REST call to Token API");

            response = httpRequestService.makePostRequest(TOKEN_API_URL, "Basic", encodedIdAndSecret,
                    "application/x-www-form-urlencoded", authorizationPayload);

            if (response.getStatusLine().getStatusCode() == 401) {
                log.error("Error making the request to token API, the request is unauthorized");
                throw new PluginExecutionException("Unauthorized request, check email, organizationKey, password");
            } else if (response.getStatusLine().getStatusCode() == 400) {
                log.error("Error making the request to token API, malformed request body");
                throw new PluginExecutionException("Bad Request, check content");
            } else {
                if (response.getStatusLine().getStatusCode() == 200) {
                    log.debug("Access token request successful");
                } else {
                    log.error("Access token request unsuccessful");
                    throw new PluginExecutionException("Error obtaining the access token");
                }
            }

            content = new String(IOUtils.toByteArray(response.getEntity().getContent()), Charsets.UTF_8);
            log.debug("Received the response");

            JSONParser parser = new JSONParser();
            accessTokenJson = (JSONObject) parser.parse(content);

        } catch (IOException ioException) {
            log.error("Error while making the http call to token API", ioException);
            throw new PluginExecutionException("Error while making http call to token API");
        } catch (ParseException parseException) {
            log.error("Error while parsing the http response to json", parseException);
            throw new PluginExecutionException("Error while parsing the http response to json");
        }

        log.debug("Obtaining the access token");
        return accessTokenJson.get("access_token").toString();
    }

    /**
     * Checks whether the API already exists in the cloud according to the name and the version of the API
     *
     * @param email                     Email of the cloud account to export the API
     * @param organizationKey           The key generated in the API cloud for the given credentials, unique for the
     *                                  cloud
     * @param password                  Password of the cloud account to export the API
     * @param swagger                   The POJO of the swagger definition
     * @return                          Returns true if the API already exists, else returns false
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    private String getApiIdentifier(String email, String organizationKey, String password, Swagger swagger) throws
            PluginExecutionException {

        HttpResponse response;
        String content;
        JSONObject responseJson;

        String accessToken = getAccessToken(email, organizationKey, password, "view");
        response = httpRequestService.makeGetRequest(API_CLOUD_URL, "Bearer", accessToken,
                "application/json");

        if (response.getStatusLine().getStatusCode() == 200) {
            log.debug("The API list is returned");
        } else if (response.getStatusLine().getStatusCode() == 406) {
            log.error("Not acceptable, the requested media is not supported");
            throw new PluginExecutionException("Not acceptable, the requested media is not supported");
        } else {
            log.error("The API list is not returned");
            throw new PluginExecutionException("The API list is not returned");
        }

        try {
            content = new String(IOUtils.toByteArray(response.getEntity().getContent()), Charsets.UTF_8);
        } catch (IOException ioException) {
            log.error("Error while searching APIs.", ioException);
            throw new PluginExecutionException("Error while searching APIs");
        }

        JSONParser parser = new JSONParser();
        try {
            responseJson = (JSONObject) parser.parse(content);
        } catch (ParseException parseException) {
            log.error("Error while parsing search API response", parseException);
            throw  new PluginExecutionException("Error parsing search API response to json");
        }

        String version = swagger.getInfo().getVersion();
        String name = swagger.getInfo().getTitle();

        JSONArray apiList = (JSONArray) responseJson.get("list");

        JSONObject api;
        int length = apiList.size();
        int apiIndex;

        for (apiIndex = 0; apiIndex < length; apiIndex++) {
            api = (JSONObject) apiList.get(apiIndex);
            if (name.equals(api.get("name")) && version.equals(api.get("version"))) {
                return api.get("id").toString();
            } else {
                if (apiIndex == length - 1) {
                    return null;
                }
            }
        }
        return null;
    }

    /**
     *
     * @param email                     Email of the cloud account to export the API
     * @param organizationKey           The key generated in the API cloud for the given credentials, unique for the
     *                                  cloud
     * @param password                  Password of the cloud account to export the API
     * @param apiIdentifier             Identifier of the API
     * @param payload                   Body of the request to be made
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    private void updateApi(String email, String organizationKey, String password, String apiIdentifier,
                          StringEntity payload) throws PluginExecutionException {

        String accessToken = getAccessToken(email, organizationKey, password, "create");
        HttpResponse response = httpRequestService.makePutRequest(API_CLOUD_URL + apiIdentifier, "Bearer",
                accessToken, "application/json", payload);
        System.out.println(response.getStatusLine().getStatusCode());
        if (response.getStatusLine().getStatusCode() == 200) {
            log.debug("The API is updated");
        } else {
            if (response.getStatusLine().getStatusCode() == 400) {
                log.error("Bad Request. Invalid request or validation error");
                throw new PluginExecutionException("Bad Request. Invalid request or validation error");
            } else if (response.getStatusLine().getStatusCode() == 403) {
                log.error("Forbidden. The request must be conditional but no condition has been specified.");
                throw new PluginExecutionException("Forbidden. The request must be conditional but no condition" +
                        " has been specified.");
            } else if (response.getStatusLine().getStatusCode() == 404) {
                log.error("Nothing to update");
                throw new PluginExecutionException("Nothing to update");
            } else if (response.getStatusLine().getStatusCode() == 412) {
                log.error("Precondition Failed. The request has not been performed because one of the preconditions" +
                        " is not met.");
                throw new PluginExecutionException("Precondition Failed. The request has not been performed because" +
                        " one of the preconditions is not met.");
            } else {
                log.error("The API is not updated");
            }
        }
    }

    /**
     * Creates an API in the api cloud and prints the response of the details of the API made
     *
     * @param email                     Email of the cloud account to export the API
     * @param organizationKey           The key generated in the API cloud for the given credentials, unique for the
     *                                  cloud
     * @param password                  Password of the cloud account to export the API
     * @param payload                   Body of the API http request to create an API
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    public void saveAPI(String email, String organizationKey, String password, String payload) throws
            PluginExecutionException {

        StringEntity creationPayload;
        HttpResponse postResponse;
        Swagger swagger;
        JSONObject swaggerDefinitionJson;

        JSONParser parser = new JSONParser();

        try {
            swaggerDefinitionJson = (JSONObject) parser.parse(payload);
        } catch (ParseException parseException) {
            log.error("Error parsing payload", parseException);
            throw new PluginExecutionException("Error parsing payload");
        }

        String definition = swaggerDefinitionJson.get("apiDefinition").toString();

        try {
            swagger = Json.mapper().readValue(definition, Swagger.class);
        } catch (IOException ioException) {
            log.error("Error parsing swagger definition", ioException);
            throw new PluginExecutionException("Error parsing swagger definition");
        }

        String apiIdentifier = getApiIdentifier(email, organizationKey, password, swagger);

        try {
            creationPayload = new StringEntity(payload);
        } catch (UnsupportedEncodingException unsupportedEncodingException) {
            log.error("The character encoding is not supported for the payload", unsupportedEncodingException);
            throw new PluginExecutionException("The character encoding is not supported");
        }

        if (apiIdentifier == null) {

            String accessToken = getAccessToken(email, organizationKey, password, "create");

            log.debug("Creating the API in the cloud");
            postResponse = httpRequestService.makePostRequest(API_CLOUD_URL, "Bearer", accessToken,
                    "application/json", creationPayload);

            if (postResponse.getStatusLine().getStatusCode() == 201) {
                log.debug("The API is created in the cloud");
            } else {
                if (postResponse.getStatusLine().getStatusCode() == 400) {
                log.error("Error creating the API, already exists with a different context");
                throw new PluginExecutionException("Bad content");
                } else if (postResponse.getStatusLine().getStatusCode() == 415) {
                    log.error("Unsupported media type");
                    throw new PluginExecutionException("Error creating the API, unsupported media type");
                } else {
                    log.debug("The API is not created in the cloud");
                }
            }

        } else {
            updateApi(email, organizationKey, password, apiIdentifier, creationPayload);
        }
    }
}

