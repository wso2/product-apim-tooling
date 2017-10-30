package com.swagger.plugins.wso2;

import io.swagger.models.Info;
import io.swagger.models.Model;
import io.swagger.models.ModelImpl;
import io.swagger.models.Swagger;
import io.swagger.models.properties.StringProperty;
import io.swagger.util.Json;
import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;


/*****************************************************************
 *  Class name :   WSO2ApiGatewayPlugin
 * Methods :       beforeApiVersionSaved, afterApiVersionSaved, configure, getConfigurationSchema
 * Attributes :    log, userEmail, userPassword, userOrganizationKey, apiId, context
 * Functionality : Contains the methods to check api identifier, obtain user inputs and push the API to the cloud
 * Visibility :    Public
 * ****************************************************************/
public class WSO2ApiGatewayPlugin {

    private static final Logger log = LoggerFactory.getLogger(WSO2ApiGatewayPlugin.class);

    private String accessToken;

    public String beforeApiVersionSaved(/*String triggeredByUUID, String objectPath, */String swaggerYaml/*,
     Boolean forceUpdate, Collection<SpecEntry> links, Boolean isPrivate*/) throws PluginExecutionException {

        Swagger swagger;
        try {
            swagger = Json.mapper().readValue(PayloadConfiguration.convertYamlToJson(swaggerYaml), Swagger.class);
        } catch (Exception exception) {
            log.error("Swagger definition is invalid or not readable", exception);
            throw new PluginExecutionException("Swagger definition is invalid or not readable");
        }
        Info info = swagger.getInfo();
        if (info == null || (StringUtils.isBlank(info.getTitle()))) {
            log.error("Missing info section in definition");
            throw new PluginExecutionException("Definition must have an info section with a unique title.");
        }
        if (StringUtils.isBlank(swagger.getBasePath())) {
            log.error("Missing base path in definition");
            throw new PluginExecutionException("Definition must have a unique 'basepath' section");
        }
        return swaggerYaml;
    }

    /**
     * Perform the deployment to WSO2 API cloud during the SAVE operation.
     *
     * @param swaggerYaml               The swagger definition of the API to be exported to the cloud
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    public void afterApiVersionSaved(/*String triggeredByUUID, String objectPath, */String swaggerYaml)
            throws PluginExecutionException {

        Swagger swagger;

        PayloadConfiguration configuration = new PayloadConfiguration();

        log.debug("Obtaining user inputs");
        configure();

        log.debug("Creating the payload from user inputs");
        String creationPayload;

        try {
            swagger = Json.mapper().readValue(PayloadConfiguration.convertYamlToJson(swaggerYaml), Swagger.class);
            creationPayload = configuration.configurePayload(swagger);
        } catch (IOException e) {
            log.error("Error while input/output operation");
            throw new PluginExecutionException("Error while input/output operation");
        }

        WSO2Api api = new WSO2Api(new MakeHttpRequestService());
        api.saveAPI(accessToken, creationPayload);
    }

    /**
     * Called when triggering the Plugin. The `config` value will be supplied and populated with all
     * `required` values from the configurationSchema, and any optional values specified by the user
     * in the integrations interface.
     */
    private void configure() {
        //For now, user inputs are set to get as command line arguments
        this.accessToken = System.getProperty("token");
    }

    /**
     * Returns a model with user to obtain user's email address, password, organization key and context
     *
     * @return Returns a model
     */
    public Model getConfigurationSchema() {
        ModelImpl model = new ModelImpl()
                .description("Exports the API to your WSO2 API Cloud instance. When you save, your API will" +
                        " be exported to your WSO2 API Cloud account. If the API already exists in your WSO2 API " +
                        "Cloud instance, the API will be updated.")
                .property("token", new StringProperty()
                        .required(true)
                        .title("WSO2 API Cloud API Access Token")
                        .description("Access token to authorize SwaggerHub to export/update APIs to/in the WSO2 API " +
                                "Cloud"));
        return model;
    }

}
