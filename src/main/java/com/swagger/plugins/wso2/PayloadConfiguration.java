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
     * @param swaggerYaml               The swagger definition of the API to be exported to the cloud
     * @param context                   The context of the API to be exported to the cloud
     * @return                          Returns the configured payload
     * @throws PluginExecutionException Custom exception to make the exception more readable
     */
    public String configurePayload(String email, String organizationKey, String swaggerYaml,
                                   String context) throws PluginExecutionException, IOException {

        ObjectMapper objectMapper;
        Swagger swagger;
        String payload;

        log.info("Setting default values for the api creation payload");
        String[] schemes = {"http", "https"};
        String[] defaultTier = {"Unlimited"};
        String visibility = "PUBLIC";

        log.info("Creating the swagger POJO for extracting info from the swagger definition");
        try {
            objectMapper = new ObjectMapper();
            swagger = Json.mapper().readValue(convertYamlToJson(swaggerYaml), Swagger.class);
        } catch (Exception ioException) {
            log.error("Error while creating the swagger POJO", ioException);
            throw new PluginExecutionException("Swagger definition is invalid or not readable");
        }

        log.info("Getting info from swagger POJO");
        String name = swagger.getInfo().getTitle();
        String version = swagger.getInfo().getVersion();
        String description = swagger.getInfo().getDescription();

        log.info("Setting the context for the API");
        swagger.setBasePath(context);

        log.info("Creating the POJO for the payload to create the API");
        PayloadStructure structure = new PayloadStructure();

        log.info("Setting the values for the POJO to create an API in the cloud");
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

        log.info("Converting the POJO to a json string");
        try {
            payload = objectMapper.writeValueAsString(structure);
        } catch (JsonProcessingException jsonProcessingException) {
            log.error("Error while converting the POJO to json", jsonProcessingException);
            throw new PluginExecutionException("Error when converting pojo to json");
        }
        log.info("Returning the payload for creating the API in the cloud");
        return payload;
    }

    /**
     * Converts a Yaml to Json
     *
     * @param swaggerYaml   The swagger definition of the API to be exported to the cloud
     * @return              Returns the json string of the Yaml
     */
    public static String convertYamlToJson(String swaggerYaml) throws IOException {
        Swagger swagger;

        log.info("Parsing Yaml to JSON");
        try {
            JsonNode node = Yaml.mapper().readValue(swaggerYaml, JsonNode.class);
            swagger = new SwaggerParser().read(node);
        } catch (IOException ioException) {
            log.error("Error while converting the Yaml to Json", ioException);
            throw  ioException;
        }
        String json =  Json.pretty(swagger);

        log.info("Minifying the JSON");
        ObjectMapper objectMapper = new ObjectMapper();
        JsonNode jsonNode = objectMapper.readValue(json, JsonNode.class);

        log.info("Returning the minified JSON");
        return jsonNode.toString();
    }

}
