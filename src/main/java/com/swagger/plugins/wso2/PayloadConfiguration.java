package com.swagger.plugins.wso2;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import io.swagger.models.Swagger;
import io.swagger.util.Json;

import java.io.IOException;


/*****************************************************************
 *  Class name : PayloadConfiguration
 * Methods : payloadToPojo
 * Functionality : Contains the method to convert string payload to a java pojo
 * Visibility : Public
 * ****************************************************************/
public class PayloadConfiguration {

    ObjectMapper mapper;

    /*
    * Method name : payloadToPojo
    * Functionality : Converts the string payload to a java pojo
    * @param : String
    * @return : void
    * */
    public void payloadToPojo(String payload) throws IOException {
        mapper = new ObjectMapper();
        PayloadStructure structure = mapper.readValue(payload, PayloadStructure.class);
    }

    /*
    * Method name : pojoToPayload
    * Functionality : Converts java pojo to a json
    * @param : PayloadStructure
    * @return : void
    * */
    public void pojoToPayload(PayloadStructure payloadStructure) throws JsonProcessingException {
        mapper = new ObjectMapper();
        String jsonPayload = mapper.writerWithDefaultPrettyPrinter().writeValueAsString(payloadStructure);
    }

    public String configurePayload(String email, String organizationKey, String swaggerDefinition, String context) throws IOException {
        Swagger swagger = null;
        String[] schemes = {"http","https"};
        String[] defaultTier = {"Unlimited"};
        ObjectMapper objectMapper = new ObjectMapper();

        //Getting pojo into the method
        PayloadStructure structure = new PayloadStructure();
        swagger = Json.mapper().readValue(swaggerDefinition, Swagger.class);

        structure.setName(swagger.getInfo().getTitle());
        structure.setVersion(swagger.getInfo().getVersion());
        structure.setContext(context);
        structure.setProvider(email+"@"+organizationKey);
        structure.setApiDefinition(swaggerDefinition);
        structure.setIsDefaultVersion(false);
        structure.setTransport(schemes);
        structure.setTiers(defaultTier);
        structure.setVisibility("PUBLIC");
        structure.setEndpointConfig("");
        CorsConfiguration configuration = new CorsConfiguration(false);
        structure.setCorsConfiguration(configuration);

        //Converting the pojo to json
        String payload = objectMapper.writeValueAsString(structure);
        return payload;
    }

}
