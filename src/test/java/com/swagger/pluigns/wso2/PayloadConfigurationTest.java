package com.swagger.pluigns.wso2;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.swagger.plugins.wso2.CorsConfiguration;
import com.swagger.plugins.wso2.Main;
import com.swagger.plugins.wso2.PayloadConfiguration;
import com.swagger.plugins.wso2.PayloadStructure;
import com.swagger.plugins.wso2.PluginExecutionException;
import io.swagger.models.Info;
import io.swagger.models.Swagger;
import io.swagger.parser.SwaggerParser;
import io.swagger.util.Json;
import io.swagger.util.Yaml;
import org.junit.Assert;
import org.junit.Test;

import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;

public class PayloadConfigurationTest {

    String inputYaml = Main.fileReder("/home/akila/Akila/swagger-wso2-integration-local/src/test/" +
            "java/com/swagger/pluigns/wso2/SwaggerSimpleApi.yaml");

    String extectedJson = "";

    @Test
    public void convertYamlToJsonTestHappyPath() throws IOException {

        Swagger swagger;
        try {
            JsonNode node = Yaml.mapper().readValue(inputYaml, JsonNode.class);
            swagger = new SwaggerParser().read(node);
        } catch (IOException e) {
            throw e;
        }

        String json = Json.pretty(swagger);
        String got = PayloadConfiguration.convertYamlToJson(inputYaml);
        Assert.assertEquals(extectedJson, got);
    }

    String expectedPayload = "";

    @Test
    public void configurePayloadTest() throws PluginExecutionException, IOException {

        PayloadStructure structure = new PayloadStructure();
        ObjectMapper objectMapper;
        Swagger swagger;
        String payload;
        String[] schemes = {"http","https"};
        String[] defaultTier = {"Unlimited"};

        try {
            objectMapper = new ObjectMapper();
            swagger = Json.mapper().readValue(PayloadConfiguration.convertYamlToJson(inputYaml), Swagger.class);
        } catch (Exception e) {
            throw new PluginExecutionException("Swagger definition is invalid or not readable");
        }

        Info info = new Info();
        String name = swagger.getInfo().getTitle();
        swagger.setInfo(info.title(name));
        swagger.setInfo(info.version("1.0.0"));
        swagger.setBasePath("sample");

        structure.setName(swagger.getInfo().getTitle());
        structure.setVersion(swagger.getInfo().getVersion());
        structure.setContext(swagger.getBasePath());
        structure.setProvider("akilaaroshana@hotmail.com@ms9714");
        structure.setApiDefinition(PayloadConfiguration.convertYamlToJson(inputYaml));
        structure.setIsDefaultVersion(false);
        structure.setTransport(schemes);
        structure.setTiers(defaultTier);
        structure.setVisibility("PUBLIC");
        structure.setEndpointConfig("");
        CorsConfiguration configuration = new CorsConfiguration(false);
        structure.setCorsConfiguration(configuration);

        try {
            payload = objectMapper.writeValueAsString(structure);
        } catch (JsonProcessingException e) {
            throw new PluginExecutionException("Error when converting pojo to json");
        }

        Assert.assertEquals(expectedPayload, payload);

    }

}
