/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package com.swagger.plugins.wso2;

import com.fasterxml.jackson.databind.JsonNode;
import com.smartbear.config.Configuration;
import com.smartbear.swaggerhub.plugins.AnnotatedPlugin;
import com.smartbear.swaggerhub.plugins.Plugin;
import com.smartbear.swaggerhub.plugins.PluginConfigurationException;
import com.smartbear.swaggerhub.plugins.PluginExecution;
import com.smartbear.swaggerhub.plugins.PluginExecutionException;
import com.smartbear.swaggerhub.plugins.PluginInfo;
import com.smartbear.swaggerhub.plugins.lifecycle.AfterApiVersionSavedPlugin;
import com.smartbear.swaggerhub.plugins.lifecycle.BeforeApiVersionSavedPlugin;
import com.smartbear.swaggerhub.plugins.lifecycle.SwaggerhubLifecycleEvent;
import com.smartbear.swaggerhub.plugins.model.SpecEntry;

import io.swagger.models.Info;
import io.swagger.models.Model;
import io.swagger.models.ModelImpl;
import io.swagger.models.Swagger;
import io.swagger.models.auth.AuthorizationValue;
import io.swagger.models.properties.StringProperty;

import io.swagger.parser.SwaggerParser;
import io.swagger.util.Json;
import io.swagger.util.Yaml;
import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.util.Arrays;
import java.util.Collection;
import java.util.List;

/**
 *
 * Info of the plugin
 *
 */
@PluginInfo(
        name = "WSO2 API Cloud",
        version = "1.0.0",
        id = "WSO2ApiGatewayPlugin",
        createdBy = "WSO2",
        failureBehaviour = Plugin.FailureBehavior.CONTINUE
)

/**
 *
 * Class with the methods to execute the plugin correctly
 *
 */
public class WSO2ApiGatewayPlugin extends AnnotatedPlugin implements BeforeApiVersionSavedPlugin,
        AfterApiVersionSavedPlugin {

    private static final Logger log = LoggerFactory.getLogger(WSO2ApiGatewayPlugin.class);

    private String accessToken;

    public WSO2ApiGatewayPlugin() {
        super(Arrays.asList(SwaggerhubLifecycleEvent.BEFORE_API_VERSION_SAVED.name(),
                SwaggerhubLifecycleEvent.AFTER_API_VERSION_SAVED.name()));
    }

    /**
     * Triggers before saving the API to make sure the swagger definition contains the info section and the 'basepath'
     *
     * @param triggeredByUUID
     * @param objectPath
     * @param swaggerYaml
     * @param forceUpdate
     * @param links
     * @param isPrivate
     * @return Returns the swagger definition
     * @throws PluginExecutionException Thrown when an exception is caught while the plugin executes
     */
    @Override
    public String beforeApiVersionSaved(String triggeredByUUID, String objectPath, String swaggerYaml,
                                        Boolean forceUpdate, Collection<SpecEntry> links, Boolean isPrivate) throws
            PluginExecutionException {

        Swagger swagger;
        try {
            swagger = Json.mapper().readValue((swaggerYaml), Swagger.class);
        } catch (Exception exception) {
            log.error("Swagger definition is invalid or not readable", exception);
            throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT,
                    "Swagger definition is invalid or not readable");
        }
        Info info = swagger.getInfo();
        if (info == null || (StringUtils.isBlank(info.getTitle()))) {
            log.error("Missing info section in definition");
            throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Definition must have an " +
                    "info section with a unique title.");
        }
        if (StringUtils.isBlank(swagger.getBasePath())) {
            log.error("Missing base path in definition");
            throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Definition must have a " +
                    "unique 'basepath' section");
        }
        return swaggerYaml;
    }

    /**
     * Perform the deployment to WSO2 API cloud during the SAVE operation
     *
     * @param triggeredByUUID
     * @param objectPath
     * @param swaggerYaml
     * @throws PluginExecutionException Thrown when an exception is caught while the plugin executes
     */
    @Override
    public void afterApiVersionSaved(String triggeredByUUID, String objectPath, String swaggerYaml) throws
            PluginExecutionException {

        Swagger swagger;
        String swaggerJson = convertSwaggerToJson(swaggerYaml);

        PayloadConfiguration configuration = new PayloadConfiguration();

        log.debug("Creating the payload from user inputs");
        String creationPayload;

        try {
            swagger = Json.mapper().readValue(swaggerJson, Swagger.class);
            creationPayload = configuration.configurePayload(swagger, swaggerJson);
        } catch (IOException e) {
            log.error("Error while input/output operation");
            throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Error while input/output" +
                    " operation");
        }

        WSO2Api api = new WSO2Api(new MakeHttpRequestService());
        api.saveAPI(accessToken, creationPayload);
    }

    /**
     * Called when triggering the Plugin. The `config` value will be supplied and populated with all
     * `required` values from the configurationSchema, and any optional values specified by the user
     * in the integrations interface.
     *
     * @param config
     * @throws PluginConfigurationException
     */
    @Override
    public void configure(Configuration config) throws PluginConfigurationException {
        this.accessToken = config.getString("token").trim();
    }

    /**
     * Use a `Model` to define the user-specified configuration values
     *
     * @param configuration
     * @return Returns a model to obtain the user values
     */
    @Override
    public Model getConfigurationSchema(Configuration configuration) {
        StringProperty property = new StringProperty();
        property._enum(SwaggerhubLifecycleEvent.BEFORE_API_VERSION_SAVED.name());
        property._enum(SwaggerhubLifecycleEvent.AFTER_API_VERSION_SAVED.name());
        return (new ModelImpl().description("Exports the API to your WSO2 API Cloud instance. When you save, your" +
                " API will be exported to your WSO2 API Cloud account. If the API already exists in your WSO2 API " +
                "Cloud instance, the API will be updated.")).property("token",
                (new StringProperty()).required(true).title("WSO2 API Access Token"));
    }

    @Override
    protected boolean process(PluginExecution pluginExecution) {
        return false;
    }

    /**
     * Helper to convert a Swagger definition, in YAML, into JSON and also fully resolve
     * any remotely defined objects, which may be private
     *
     * @param swaggerYaml               Swagger definition in yaml format
     * @return                          Returns the json of the yaml provided
     * @throws PluginExecutionException Thrown when an exception is caught while the plugin executes
     **/
    private String convertSwaggerToJson(String swaggerYaml) throws PluginExecutionException {
        try {
            List<AuthorizationValue> auth = null;
            if (accessToken != null) {
                AuthorizationValue authValue = new AuthorizationValue()
                        .type("header")
                        .keyName("Authorization")
                        .value(accessToken);
                auth = Arrays.asList(authValue);
            }
            JsonNode node = Yaml.mapper().readValue(swaggerYaml, JsonNode.class);
            Swagger swagger = new SwaggerParser().read(node, auth, true);
            return Json.pretty(swagger);
        } catch (IOException e) {
            throw new PluginExecutionException(PluginExecutionException.INVALID_INPUT, "Error while converting " +
                    "SwaggerYaml to Json");
        }
    }
}
