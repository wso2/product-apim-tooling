package com.swagger.plugins.wso2;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import io.swagger.models.Swagger;
import io.swagger.parser.SwaggerParser;
import io.swagger.util.Json;
import io.swagger.util.Yaml;

import java.io.IOException;


/*****************************************************************
 *  Class name : PayloadConfiguration
 * Methods : payloadToPojo
 * Functionality : Contains the method to convert string payload to a java pojo
 * Visibility : Public
 * ****************************************************************/
public class PayloadConfiguration {

    /**
     * Configures the paylod to create an API in the cloud setting relevant values from the swagger definition.
     * @param email
     * @param organizationKey
     * @param swaggerYaml
     * @param context
     * @return Returns the configured payload
     * @throws PluginExecutionException
     */
    public String configurePayload(String email, String organizationKey, String swaggerYaml,
                                   String context) throws PluginExecutionException, IOException {

        PayloadStructure structure = new PayloadStructure();
        ObjectMapper objectMapper;
        Swagger swagger;
        String payload;

        //Setting default values
        String[] schemes = {"http","https"};
        String[] defaultTier = {"Unlimited"};
        String visibility = "PUBLIC";

        try {
            objectMapper = new ObjectMapper();
            swagger = Json.mapper().readValue(convertYamlToJson(swaggerYaml), Swagger.class);
        } catch (Exception e) {
            throw new PluginExecutionException("Swagger definition is invalid or not readable");
        }

        String name = swagger.getInfo().getTitle();
        String version = swagger.getInfo().getVersion();
        String description = swagger.getInfo().getDescription();

        swagger.setBasePath(context);

        //Setting values to the pojo
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

        //Converting the pojo to json
        try {
            payload = objectMapper.writeValueAsString(structure);
        } catch (JsonProcessingException e) {
            throw new PluginExecutionException("Error when converting pojo to json");
        }
        return payload;
    }

    /**
     * Converts a Yaml to Json
     * @param swaggerYaml The swagger drfinition of the API to be exported to the cloud
     * @return Returns the json string of the yaml
     */
    public static String convertYamlToJson(String swaggerYaml) throws IOException {
        Swagger swagger = null;
        try {
            JsonNode node = Yaml.mapper().readValue(swaggerYaml, JsonNode.class);
            swagger = new SwaggerParser().read(node);
        } catch (IOException e) {
            e.printStackTrace();
        }
        String json =  Json.pretty(swagger);

        //Minifying json and return
        ObjectMapper objectMapper = new ObjectMapper();
        JsonNode jsonNode = objectMapper.readValue(json, JsonNode.class);
        return jsonNode.toString();
    }

}
