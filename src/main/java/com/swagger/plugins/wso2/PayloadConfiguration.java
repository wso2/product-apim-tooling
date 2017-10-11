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
 *  Class name : PayloadConfiguration
 * Methods : payloadToPojo
 * Functionality : Contains the method to convert string payload to a java pojo
 * Visibility : Public
 * ****************************************************************/
public class PayloadConfiguration {

    private static final Logger LOGGER = LoggerFactory.getLogger(PayloadConfiguration.class);

    /**
     * Configures the paylod to create an API in the cloud setting relevant values from the swagger definition.
     *
     * @param email             Email of the cloud account to export the API
     * @param organizationKey   The key generated using email and password
     * @param swaggerYaml       The swagger definition of the API to be exported to the cloud
     * @param context           The context for the API in the cloud
     * @return Returns the configured payload
     * @throws PluginExecutionException
     */
    public String configurePayload(String email, String organizationKey, String swaggerYaml,
                                   String context) throws PluginExecutionException, IOException {

        ObjectMapper objectMapper;
        Swagger swagger;
        String payload;

        LOGGER.info("Setting default values for the payload");
        String[] schemes = {"http","https"};
        String[] defaultTier = {"Unlimited"};
        String visibility = "PUBLIC";

        LOGGER.info("Creating the swagger POJO for extracting info from the swagger definition");
        try {
            objectMapper = new ObjectMapper();
            swagger = Json.mapper().readValue(convertYamlToJson(swaggerYaml), Swagger.class);
        } catch (Exception e) {
            throw new PluginExecutionException("Swagger definition is invalid or not readable");
        }

        LOGGER.info("Getting info from swagger POJO");
        String name = swagger.getInfo().getTitle();
        String version = swagger.getInfo().getVersion();
        String description = swagger.getInfo().getDescription();

        LOGGER.info("Setting the context for the API");
        swagger.setBasePath(context);

        LOGGER.info("Creating the POJO for the payload to create the API");
        PayloadStructure structure = new PayloadStructure();

        LOGGER.info("Setting the essential values for the POJO to create an API in the cloud");
        structure.setName(name);
        structure.setVersion(version);
        structure.setDescription(description);
        structure.setContext(swagger.getBasePath());
        structure.setProvider(email+"@"+organizationKey);
        structure.setApiDefinition(convertYamlToJson(swaggerYaml));
        structure.setIsDefaultVersion(false);
        structure.setTransport(schemes);
        structure.setTiers(defaultTier);
        structure.setVisibility(visibility);
        structure.setEndpointConfig("");
        CorsConfiguration configuration = new CorsConfiguration(false);
        structure.setCorsConfiguration(configuration);

        LOGGER.info("Converting the POJO to a json string");
        try {
            payload = objectMapper.writeValueAsString(structure);
        } catch (JsonProcessingException e) {
            throw new PluginExecutionException("Error when converting pojo to json");
        }
        LOGGER.info("Returning the payload for creating the API in the cloud");
        return payload;
    }

    /**
     * Converts a Yaml to Json
     *
     * @param swaggerYaml The swagger drfinition of the API to be exported to the cloud
     * @return Returns the json string of the yaml
     */
    public static String convertYamlToJson(String swaggerYaml) throws IOException {
        Swagger swagger = null;

        LOGGER.info("Parsing Yaml to JSON");
        try {
            JsonNode node = Yaml.mapper().readValue(swaggerYaml, JsonNode.class);
            swagger = new SwaggerParser().read(node);
        } catch (IOException e) {
            e.printStackTrace();
        }
        String json =  Json.pretty(swagger);

        LOGGER.info("Minifying the JSON");
        ObjectMapper objectMapper = new ObjectMapper();
        JsonNode jsonNode = objectMapper.readValue(json, JsonNode.class);

        LOGGER.info("Returning the minified JSON");
        return jsonNode.toString();
    }

}
