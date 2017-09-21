package com.swagger.plugins.wso2;

import com.fasterxml.jackson.databind.JsonNode;
import io.swagger.models.Info;
import io.swagger.models.Swagger;
import io.swagger.parser.SwaggerParser;
import io.swagger.util.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;

public class Wso2ApiGatewayPlugin {

    private static final Logger LOGGER = LoggerFactory.getLogger(Wso2ApiGatewayPlugin.class);


    /*
    * Trigger before an API is saved.
    * */
    public String beforeApiVersionSaved(String swaggerYaml) throws Exception {
        Swagger swagger = null;

        try {
            LOGGER.info("Convert swaggerYaml to a POJO");
            swagger = Json.mapper().readValue(convertYamlToJson(swaggerYaml), Swagger.class);
            if (swagger == null) {
                throw new Exception();
            }
        } catch (IOException e) {
            e.printStackTrace();
            LOGGER.warn("swagger object is null");
        }

        return null;

    }


    /*
    * Triggered after an API is saved.
    * */
    public void afterApiVersionSaved(String swaggerYaml) {
        Swagger swagger = null;
        try {
            swagger = Json.mapper().readValue(convertYamlToJson(swaggerYaml), Swagger.class);
        } catch (IOException e) {
            e.printStackTrace();
            LOGGER.error("Swagger definition is invalid or not readable");
        }

        Info swaggerInfo = swagger.getInfo();
        String apiName = swaggerInfo.getTitle();
        String apiVersion = swaggerInfo.getVersion();



    }

    public String convertYamlToJson(String swaggerYaml) {
        Swagger swagger = null;
        try {
            JsonNode node = Yaml.mapper().readValue(swaggerYaml, JsonNode.class);
            swagger = new SwaggerParser().read(node);
        } catch (IOException e) {
            e.printStackTrace();
        }
        return Json.pretty(swagger);
    }

}
