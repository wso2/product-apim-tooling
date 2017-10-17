package com.swagger.plugins.wso2;

import io.swagger.models.Model;
import io.swagger.models.ModelImpl;
import io.swagger.models.Swagger;
import io.swagger.models.properties.StringProperty;
import io.swagger.util.*;
import org.apache.commons.lang3.StringUtils;
import org.json.simple.parser.ParseException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;


/*****************************************************************
 *  Class name : Wso2ApiGatewayPlugin
 * Methods : beforeApiVersionSaved, afterApiVersionSaved, configure, getConfigurationSchema
 * Functionality : Contains the methods to check api identifier, obtain user inputs and push the API to the cloud
 * Visibility : Public
 * ****************************************************************/
public class Wso2ApiGatewayPlugin {

    private static final Logger LOGGER = LoggerFactory.getLogger(Wso2ApiGatewayPlugin.class);
    private static final String WSO2_API_ID_EXTENSION = "x-wso2-api-id";
    private String userEmail;
    private String userPassword;
    private String userOrganizationKey;
    private String apiId;               //apiId will be used when calling the beforeApiVersionSaved method
    private String context;


    /**
     * This method is triggered before saving the API to ensure whether a valid identifier exists.
     *
     * @param triggeredByUUID
     * @param objectPath
     * @param swaggerYaml The swagger drfinition of the API to be exported to the cloud
     * @param forceUpdate
     * @param isPrivate
     * @return Returns the swaggerYaml after ensuring a valid api identifier exists
     * @throws PluginExecutionException
     */
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
            LOGGER.error("Swagger definition is invalid or not readable");
            throw  new PluginExecutionException("Swagger definition is invalid or not readable");
        }

        if (swagger.getVendorExtensions() == null || swagger.getVendorExtensions().get(WSO2_API_ID_EXTENSION) == null) {
            if (StringUtils.isBlank(apiId)) {
                swagger.vendorExtension(WSO2_API_ID_EXTENSION, System.currentTimeMillis());
            } else {
              swagger.vendorExtension(WSO2_API_ID_EXTENSION, apiId);
            }
        } else {
            return swaggerYaml;
        }

        try {
            return Yaml.mapper().writeValueAsString(swagger);
        } catch (Exception e) {
            LOGGER.error("Swagger definition is invalid or not readable");
            throw new PluginExecutionException("Swagger definition is invalid or not readable ");
        }
    }

    /**
     * Perform the deployment to WSO2 API cloud during the SAVE operation.
     *
     * @param swaggerYaml The swagger definition of the API to be exported to the cloud
     * @throws PluginExecutionException
     * @throws IOException
     * @throws ParseException
     */
    public void afterApiVersionSaved(/*String triggeredByUUID, String objectPath, */String swaggerYaml)
            throws PluginExecutionException, IOException, ParseException {

        PayloadConfiguration configuration = new PayloadConfiguration();

        LOGGER.info("Obtaining user inputs");
        configure();

        LOGGER.info("Creating the payload from user inputs");
        String creationPayload = configuration.configurePayload(userEmail, userOrganizationKey, swaggerYaml, context);

        Wso2Api api = new Wso2Api();
        api.getClientIdAndSecret(userEmail, userOrganizationKey, userPassword);
        api.getAccessToken(userEmail, userOrganizationKey, userPassword);
        api.saveAPI(creationPayload);
    }

    /**
     * Called when triggering the Plugin.  The `config` value will be supplied and populated with all
     * `required` values from the configurationSchema, and any optional values specified by the user
     * in the integrations interface.
     */
    public void configure() {
        this.userEmail = "fpk11@geronra.com";
        this.userOrganizationKey = "testusercompany";
        this.userPassword = "Testuser1";
        // The "/" is compulsory for the context otherwise the server will return code 400 instead of 409.
        this.context = "/"+"simple";
    }

    /**
     * Returns a model with user to obtain user's email address, password, organization key and context
     *
     * @return Returns a model
     */
    public Model getConfigurationSchema() {
        ModelImpl model = new ModelImpl()
                .description("Exports your API to WSO2 API cloud.  When you save, your API will be exported to" +
                        " your WSO2 API cloud account. The Vendor extension 'x-wso2-api-id' will be added to this" +
                        " Swagger definition file in order to link it with the API saved into your WSO2 API cloud" +
                        " account.")
                .property("userEmail", new StringProperty()
                        .required(true)
                        .title("Email")
                        .description("The email address of the WSO2 API cloud account."))
                .property("password", new StringProperty()
                        .title("userPassword")
                        .description("The password of the WSO2 API cloud account."))
                .property("organizationKey", new StringProperty()
                        .required(true)
                        .title("organizationKey")
                        .description("The organization key of the WSO2 API cloud account."))
                .property("context", new StringProperty()
                        .required(true)
                        .title("Context")
                        .description("The context of the API exported to WSO2 API cloud."));
        return model;
    }

}
