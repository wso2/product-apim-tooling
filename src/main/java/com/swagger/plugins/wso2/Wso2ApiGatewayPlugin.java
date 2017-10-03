package com.swagger.plugins.wso2;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import io.swagger.models.Scheme;
import io.swagger.models.Swagger;
import io.swagger.parser.SwaggerParser;
import io.swagger.util.*;
import org.json.simple.parser.ParseException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;


/*****************************************************************
 *  Class name : Wso2ApiGatewayPlugin
 * Methods : beforeApiVersionSaved, afterApiVersionSaved
 * Functionality : Contains the method to push the API to the cloud
 * Visibility : Public
 * ****************************************************************/
public class Wso2ApiGatewayPlugin {

    private static final Logger LOGGER = LoggerFactory.getLogger(Wso2ApiGatewayPlugin.class);
    private static final String WSO2_API_ID_EXTENSION = "x-wso2-api-id";
    private String userEmail;
    private String userPassword;
    private String userOrganizationKey;
    private String apiId;
    private String apiVersion;
    private String accessToken;

    /*
    * Trigger before an API is saved.
    * */
//    public String beforeApiVersionSaved(String swaggerYaml) throws Exception {
//        Swagger swagger = null;
//
//        try {
//            LOGGER.info("Convert swaggerYaml to a POJO");
//            swagger = Json.mapper().readValue(convertYamlToJson(swaggerYaml), Swagger.class);
//            if (swagger == null) {
//                throw new Exception();
//            }
//        } catch (IOException e) {
//            e.printStackTrace();
//            LOGGER.warn("swagger object is null");
//        }
//
//        return null;
//
//    }


    /*
    * Method name : afterApiVersionSaved
    * Functionality : Sets the essential values of the swagger definition elements to payload elements
    * @param : String
    * @return : void
    * */
    public void afterApiVersionSaved(String email, String organizationKey, String password, String swaggerDefinition, String version, String context) throws IOException, ParseException, PluginExecutionException {

        PayloadConfiguration configuration = new PayloadConfiguration();

        String pa = configuration.configurePayload(email, organizationKey, swaggerDefinition, version, context);

        Wso2Api api = new Wso2Api();
        accessToken = api.getAccessToken(email,organizationKey,password);
        api.saveAPI(pa,accessToken);
    }
}
