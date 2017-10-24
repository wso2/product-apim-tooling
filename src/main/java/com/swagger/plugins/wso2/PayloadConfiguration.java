package com.swagger.plugins.wso2;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import io.swagger.models.Swagger;
import io.swagger.parser.SwaggerParser;
import io.swagger.util.Json;
import io.swagger.util.Yaml;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;


/*****************************************************************
 * Class name : PayloadConfiguration
 * Methods : configurePayload, payloadToPojo
 * Functionality : Contains the method to convert string payload to a java pojo
 * Visibility : Public
 * ****************************************************************/
public class PayloadConfiguration {

    private static final Logger log = LoggerFactory.getLogger(PayloadConfiguration.class);

    /**
     * Configures the payload to create an API in the cloud setting relevant values from the swagger definition.
     *
     * @param email                     Email of the cloud account to export the API
     * @param organizationKey           The key generated using email and password
     * @param swagger                   The POJO of the swagger definition
     * @return                          Returns the configured payload
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    public String configurePayload(String email, String organizationKey, Swagger swagger) throws
            PluginExecutionException, IOException {

        String payload;
        String swaggerYaml;

        log.debug("Setting default values for the api creation payload");
        String[] schemes = {"http", "https"};
        String[] defaultTier = {"Unlimited"};
        String visibility = "PUBLIC";

        log.debug("Converting swagger POJO to a yaml to assign to the payload");
        try {
            swaggerYaml = Yaml.mapper().writeValueAsString(swagger);
        } catch (Exception exception) {
            log.error("Swagger definition is invalid or not readable", exception);
            throw new PluginExecutionException("Swagger definition is invalid or not readable");
        }

        log.debug("Getting info from swagger POJO");
        String name = swagger.getInfo().getTitle();
        String version = swagger.getInfo().getVersion();
        String description = swagger.getInfo().getDescription();

        log.debug("Creating the POJO for the payload to create the API");
        PayloadStructure structure = new PayloadStructure();

        log.debug("Setting the values for the POJO to create an API in the cloud");
        structure.setName(name);
        structure.setVersion(version);
        structure.setDescription(description);
        structure.setContext(swagger.getBasePath());
        structure.setProvider(email + "@" + organizationKey);
        structure.setApiDefinition(convertYamlToJson(swaggerYaml));
        structure.setIsDefaultVersion(false);
        structure.setTransport(schemes);
        structure.setTiers(defaultTier);
        structure.setVisibility(visibility);
        structure.setEndpointConfig("");
        CorsConfiguration configuration = new CorsConfiguration(false);
        structure.setCorsConfiguration(configuration);

        log.debug("Converting the POJO to a json string");

        try {
            ObjectMapper objectMapper = new ObjectMapper();
            payload = objectMapper.writeValueAsString(structure);
        } catch (JsonProcessingException jsonProcessingException) {
            log.error("Error while converting the POJO to json", jsonProcessingException);
            throw new PluginExecutionException("Error when converting pojo to json");
        }

        log.debug("Returning the payload for creating the API in the cloud");
        return payload;
    }

    /**
     * Converts a Yaml to Json
     *
     * @param swaggerYaml   The swagger definition of the API to be exported to the cloud
     * @return              Returns the json string of the Yaml
     */
    public static String convertYamlToJson(String swaggerYaml) throws PluginExecutionException {
        Swagger swagger;
        String json;
        JsonNode jsonNode;
        log.debug("Parsing Yaml to JSON");
        try {
            JsonNode node = Yaml.mapper().readValue(swaggerYaml, JsonNode.class);
            swagger = new SwaggerParser().read(node);

            json =  Json.pretty(swagger);

            log.debug("Minifying the JSON");
            ObjectMapper objectMapper = new ObjectMapper();
            jsonNode = objectMapper.readValue(json, JsonNode.class);

        } catch (IOException ioException) {
            log.error("Error while converting the Yaml to Json", ioException);
            throw  new PluginExecutionException("Error converting the Yaml to Json");
        }

        log.debug("Returning the minified JSON");
        return jsonNode.toString();
    }

}
