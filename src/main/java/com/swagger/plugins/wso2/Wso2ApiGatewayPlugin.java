package com.swagger.plugins.wso2;

import io.swagger.models.Swagger;
import io.swagger.util.*;
import org.apache.commons.lang3.StringUtils;
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
    private String context;


    /**
     * This method is triggered before saving the API to ensure whether a valid identifier exists.
     * @param triggeredByUUID
     * @param objectPath
     * @param swaggerYaml The swagger drfinition of the API to be exported to the cloud
     * @param forceUpdate
     * @param isPrivate
     * @return Returns the swaggerYaml after ensuring a valid api identifier exists
     * @throws PluginExecutionException
     */

    //@Override
    public String beforeApiVersionSaved(String triggeredByUUID, String objectPath, String swaggerYaml,
                                        Boolean forceUpdate,/*Collection<SpecEntry> links,*/
                                        Boolean isPrivate) throws PluginExecutionException {
        Swagger swagger;

        try {
            swagger = Json.mapper().readValue(PayloadConfiguration.convertYamlToJson(swaggerYaml), Swagger.class);
            if (swagger == null) {
                throw new Exception();
            }
        } catch (Exception e) {
            throw  new PluginExecutionException("Swagger definition is invalid or not readable");
        }

        if (swagger.getVendorExtensions() == null || swagger.getVendorExtensions().get(WSO2_API_ID_EXTENSION) == null) {
            if (StringUtils.isBlank(apiId)) {
                swagger.vendorExtension(WSO2_API_ID_EXTENSION, System.currentTimeMillis());
            } else {
              swagger.vendorExtension(WSO2_API_ID_EXTENSION,apiId);
            }
        } else {
            return swaggerYaml;
        }

        try {
            return Yaml.mapper().writeValueAsString(swagger);
        } catch (Exception e) {
            throw new PluginExecutionException("Swagger definition is invalid or not readable ");
        }
    }

    /**
     * Perform the deployment to WSO2 API cloud during the SAVE operation.
     * @param swaggerYaml The swagger drfinition of the API to be exported to the cloud
     * @throws PluginExecutionException
     * @throws IOException
     * @throws ParseException
     */

    //@Override
    public void afterApiVersionSaved(/*String triggeredByUUID, String objectPath, */String swaggerYaml)
            throws PluginExecutionException, IOException, ParseException {

        PayloadConfiguration configuration = new PayloadConfiguration();
        configure();
        String creationPayload = configuration.configurePayload(userEmail, userOrganizationKey, swaggerYaml, context);

        Wso2Api api = new Wso2Api();
        String accessToken;

        accessToken = api.getAccessToken(userEmail, userOrganizationKey, userPassword);
        api.saveAPI(creationPayload, accessToken);
    }

    /**
     * Called when triggering the Plugin.  The `config` value will be supplied and populated with all
     * `required` values from the configurationSchema, and any optional values specified by the user
     * in the integrations interface.
     */

    //@Override
    public void configure() {
        this.userEmail = "yolom@seekjobs4u.com";
        this.userOrganizationKey = "yolo4958";
        this.userPassword = "Yolofernando123";
        // The "/" is compulsory for the context otherwise the server will return code 400 instead of 409.
        this.context = "/"+"test";
    }



}
