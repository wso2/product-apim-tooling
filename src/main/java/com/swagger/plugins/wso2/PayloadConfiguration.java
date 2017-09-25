package com.swagger.plugins.wso2;

import com.fasterxml.jackson.databind.ObjectMapper;

import java.io.IOException;


/*****************************************************************
 *  Class name : PayloadConfiguration
 * Methods : payloadToPojo
 * Functionality : Contains the method to convert string payload to a java pojo
 * Visibility : Public
 * ****************************************************************/
public class PayloadConfiguration {


    /*
    * Method name : payloadToPojo
    * Functionality : Converts the string payload to a java pojo
    * @param : String
    * @return : void
    * */
    public void payloadToPojo(String payload) throws IOException {
        ObjectMapper mapper = new ObjectMapper();
        PayloadStructure structure = mapper.readValue(payload, PayloadStructure.class);
        System.out.println(structure.getApiDefinition());
    }

}
