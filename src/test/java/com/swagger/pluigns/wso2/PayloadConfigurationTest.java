package com.swagger.pluigns.wso2;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.swagger.plugins.wso2.CorsConfiguration;
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

    String inputYaml = "swagger: '2.0'\n" +
            "info:\n" +
            "  description: This is a simple API\n" +
            "  version: 1.0.0\n" +
            "  title: Simple Inventory API\n" +
            "  # put the contact info for your development or API team\n" +
            "  contact:\n" +
            "    email: you@your-company.com\n" +
            "\n" +
            "  license:\n" +
            "    name: Apache 2.0\n" +
            "    url: http://www.apache.org/licenses/LICENSE-2.0.html\n" +
            "\n" +
            "# tags are used for organizing operations\n" +
            "tags:\n" +
            "- name: admins\n" +
            "  description: Secured Admin-only calls\n" +
            "- name: developers\n" +
            "  description: Operations available to regular developers\n" +
            "\n" +
            "paths:\n" +
            "  /inventory:\n" +
            "    get:\n" +
            "      tags:\n" +
            "      - developers\n" +
            "      summary: searches inventory\n" +
            "      operationId: searchInventory\n" +
            "      description: |\n" +
            "        By passing in the appropriate options, you can search for\n" +
            "        available inventory in the system\n" +
            "      produces:\n" +
            "      - application/json\n" +
            "      parameters:\n" +
            "      - in: query\n" +
            "        name: searchString\n" +
            "        description: pass an optional search string for looking up inventory\n" +
            "        required: false\n" +
            "        type: string\n" +
            "      - in: query\n" +
            "        name: skip\n" +
            "        description: number of records to skip for pagination\n" +
            "        type: integer\n" +
            "        format: int32\n" +
            "        minimum: 0\n" +
            "      - in: query\n" +
            "        name: limit\n" +
            "        description: maximum number of records to return\n" +
            "        type: integer\n" +
            "        format: int32\n" +
            "        minimum: 0\n" +
            "        maximum: 50\n" +
            "      responses:\n" +
            "        200:\n" +
            "          description: search results matching criteria\n" +
            "          schema:\n" +
            "            type: array\n" +
            "            items:\n" +
            "              $ref: '#/definitions/InventoryItem'\n" +
            "        400:\n" +
            "          description: bad input parameter\n" +
            "    post:\n" +
            "      tags:\n" +
            "      - admins\n" +
            "      summary: adds an inventory item\n" +
            "      operationId: addInventory\n" +
            "      description: Adds an item to the system\n" +
            "      consumes:\n" +
            "      - application/json\n" +
            "      produces:\n" +
            "      - application/json\n" +
            "      parameters:\n" +
            "      - in: body\n" +
            "        name: inventoryItem\n" +
            "        description: Inventory item to add\n" +
            "        schema:\n" +
            "          $ref: '#/definitions/InventoryItem'\n" +
            "      responses:\n" +
            "        201:\n" +
            "          description: item created\n" +
            "        400:\n" +
            "          description: invalid input, object invalid\n" +
            "        409:\n" +
            "          description: an existing item already exists\n" +
            "definitions:\n" +
            "  InventoryItem:\n" +
            "    type: object\n" +
            "    required:\n" +
            "    - id\n" +
            "    - name\n" +
            "    - manufacturer\n" +
            "    - releaseDate\n" +
            "    properties:\n" +
            "      id:\n" +
            "        type: string\n" +
            "        format: uuid\n" +
            "        example: d290f1ee-6c54-4b01-90e6-d701748f0851\n" +
            "      name:\n" +
            "        type: string\n" +
            "        example: Widget Adapter\n" +
            "      releaseDate:\n" +
            "        type: string\n" +
            "        format: int32\n" +
            "        example: 2016-08-29T09:12:33.001Z\n" +
            "      manufacturer:\n" +
            "        $ref: '#/definitions/Manufacturer'\n" +
            "  Manufacturer:\n" +
            "    required:\n" +
            "    - name\n" +
            "    properties:\n" +
            "      name:\n" +
            "        type: string\n" +
            "        example: ACME Corporation\n" +
            "      homePage:\n" +
            "        type: string\n" +
            "        format: url\n" +
            "        example:  https://www.acme-corp.com\n" +
            "      phone:\n" +
            "        type: string\n" +
            "        example: 408-867-5309\n" +
            "# Added by API Auto Mocking Plugin\n" +
            "host: virtserver.swaggerhub.com\n" +
            "basePath: /yolom4/SimpleApi/1.0.0\n" +
            "schemes:\n" +
            " - https";

    String extectedJson = "{\n" +
            "  \"swagger\" : \"2.0\",\n" +
            "  \"info\" : {\n" +
            "    \"description\" : \"This is a simple API\",\n" +
            "    \"version\" : \"1.0.0\",\n" +
            "    \"title\" : \"Simple Inventory API\",\n" +
            "    \"contact\" : {\n" +
            "      \"email\" : \"you@your-company.com\"\n" +
            "    },\n" +
            "    \"license\" : {\n" +
            "      \"name\" : \"Apache 2.0\",\n" +
            "      \"url\" : \"http://www.apache.org/licenses/LICENSE-2.0.html\"\n" +
            "    }\n" +
            "  },\n" +
            "  \"host\" : \"virtserver.swaggerhub.com\",\n" +
            "  \"basePath\" : \"/yolom4/SimpleApi/1.0.0\",\n" +
            "  \"tags\" : [ {\n" +
            "    \"name\" : \"admins\",\n" +
            "    \"description\" : \"Secured Admin-only calls\"\n" +
            "  }, {\n" +
            "    \"name\" : \"developers\",\n" +
            "    \"description\" : \"Operations available to regular developers\"\n" +
            "  } ],\n" +
            "  \"schemes\" : [ \"https\" ],\n" +
            "  \"paths\" : {\n" +
            "    \"/inventory\" : {\n" +
            "      \"get\" : {\n" +
            "        \"tags\" : [ \"developers\" ],\n" +
            "        \"summary\" : \"searches inventory\",\n" +
            "        \"description\" : \"By passing in the appropriate options, you can search for\\navailable inventory in the system\\n\",\n" +
            "        \"operationId\" : \"searchInventory\",\n" +
            "        \"produces\" : [ \"application/json\" ],\n" +
            "        \"parameters\" : [ {\n" +
            "          \"name\" : \"searchString\",\n" +
            "          \"in\" : \"query\",\n" +
            "          \"description\" : \"pass an optional search string for looking up inventory\",\n" +
            "          \"required\" : false,\n" +
            "          \"type\" : \"string\"\n" +
            "        }, {\n" +
            "          \"name\" : \"skip\",\n" +
            "          \"in\" : \"query\",\n" +
            "          \"description\" : \"number of records to skip for pagination\",\n" +
            "          \"required\" : false,\n" +
            "          \"type\" : \"integer\",\n" +
            "          \"minimum\" : 0,\n" +
            "          \"format\" : \"int32\"\n" +
            "        }, {\n" +
            "          \"name\" : \"limit\",\n" +
            "          \"in\" : \"query\",\n" +
            "          \"description\" : \"maximum number of records to return\",\n" +
            "          \"required\" : false,\n" +
            "          \"type\" : \"integer\",\n" +
            "          \"maximum\" : 50,\n" +
            "          \"minimum\" : 0,\n" +
            "          \"format\" : \"int32\"\n" +
            "        } ],\n" +
            "        \"responses\" : {\n" +
            "          \"200\" : {\n" +
            "            \"description\" : \"search results matching criteria\",\n" +
            "            \"schema\" : {\n" +
            "              \"type\" : \"array\",\n" +
            "              \"items\" : {\n" +
            "                \"$ref\" : \"#/definitions/InventoryItem\"\n" +
            "              }\n" +
            "            }\n" +
            "          },\n" +
            "          \"400\" : {\n" +
            "            \"description\" : \"bad input parameter\"\n" +
            "          }\n" +
            "        }\n" +
            "      },\n" +
            "      \"post\" : {\n" +
            "        \"tags\" : [ \"admins\" ],\n" +
            "        \"summary\" : \"adds an inventory item\",\n" +
            "        \"description\" : \"Adds an item to the system\",\n" +
            "        \"operationId\" : \"addInventory\",\n" +
            "        \"consumes\" : [ \"application/json\" ],\n" +
            "        \"produces\" : [ \"application/json\" ],\n" +
            "        \"parameters\" : [ {\n" +
            "          \"in\" : \"body\",\n" +
            "          \"name\" : \"inventoryItem\",\n" +
            "          \"description\" : \"Inventory item to add\",\n" +
            "          \"required\" : false,\n" +
            "          \"schema\" : {\n" +
            "            \"$ref\" : \"#/definitions/InventoryItem\"\n" +
            "          }\n" +
            "        } ],\n" +
            "        \"responses\" : {\n" +
            "          \"201\" : {\n" +
            "            \"description\" : \"item created\"\n" +
            "          },\n" +
            "          \"400\" : {\n" +
            "            \"description\" : \"invalid input, object invalid\"\n" +
            "          },\n" +
            "          \"409\" : {\n" +
            "            \"description\" : \"an existing item already exists\"\n" +
            "          }\n" +
            "        }\n" +
            "      }\n" +
            "    }\n" +
            "  },\n" +
            "  \"definitions\" : {\n" +
            "    \"InventoryItem\" : {\n" +
            "      \"type\" : \"object\",\n" +
            "      \"required\" : [ \"id\", \"manufacturer\", \"name\", \"releaseDate\" ],\n" +
            "      \"properties\" : {\n" +
            "        \"id\" : {\n" +
            "          \"type\" : \"string\",\n" +
            "          \"format\" : \"uuid\",\n" +
            "          \"example\" : \"d290f1ee-6c54-4b01-90e6-d701748f0851\"\n" +
            "        },\n" +
            "        \"name\" : {\n" +
            "          \"type\" : \"string\",\n" +
            "          \"example\" : \"Widget Adapter\"\n" +
            "        },\n" +
            "        \"releaseDate\" : {\n" +
            "          \"type\" : \"string\",\n" +
            "          \"format\" : \"int32\",\n" +
            "          \"example\" : \"2016-08-29T09:12:33.001Z\"\n" +
            "        },\n" +
            "        \"manufacturer\" : {\n" +
            "          \"$ref\" : \"#/definitions/Manufacturer\"\n" +
            "        }\n" +
            "      }\n" +
            "    },\n" +
            "    \"Manufacturer\" : {\n" +
            "      \"required\" : [ \"name\" ],\n" +
            "      \"properties\" : {\n" +
            "        \"name\" : {\n" +
            "          \"type\" : \"string\",\n" +
            "          \"example\" : \"ACME Corporation\"\n" +
            "        },\n" +
            "        \"homePage\" : {\n" +
            "          \"type\" : \"string\",\n" +
            "          \"format\" : \"url\",\n" +
            "          \"example\" : \"https://www.acme-corp.com\"\n" +
            "        },\n" +
            "        \"phone\" : {\n" +
            "          \"type\" : \"string\",\n" +
            "          \"example\" : \"408-867-5309\"\n" +
            "        }\n" +
            "      }\n" +
            "    }\n" +
            "  }\n" +
            "}";

//    public String fileReader(String filePath) throws IOException {
//
//        String out = null;
//        StringBuffer stringBuffer = null;
//        BufferedReader bufferedReader = null;
//        FileReader fileReader = null;
//
//        try {
//            fileReader = new FileReader(filePath);
//            bufferedReader = new BufferedReader(fileReader);
//
//            String sCurrentLine;
//            stringBuffer = new StringBuffer();
//
//            while ((sCurrentLine = bufferedReader.readLine()) != null) {
////                stringBuffer.append(sCurrentLine);
//                out += sCurrentLine+"\n";
//            }
//        } catch (IOException e) {
//            throw e;
//        } finally {
//            try {
//                if (bufferedReader != null) {
//                    bufferedReader.close();
//                }
//                if (fileReader != null) {
//                    fileReader.close();
//                }
//            } catch (IOException ex) {
//                ex.printStackTrace();
//            }
//        }
//        return out;
//    }

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

    String expectedPayload = "{\"name\":\"Simple Inventory API\",\"context\":\"sample\"," +
            "\"version\":\"1.0.0\",\"provider\":\"akilaaroshana@hotmail.com@ms9714\"," +
            "\"apiDefinition\":\"{\\n  \\\"swagger\\\" : \\\"2.0\\\",\\n  " +
            "\\\"info\\\" : {\\n    \\\"description\\\" : \\\"This is a simple API\\\",\\n   " +
            " \\\"version\\\" : \\\"1.0.0\\\",\\n    \\\"title\\\" : \\\"Simple Inventory API\\\",\\n  " +
            "  \\\"contact\\\" : {\\n      \\\"email\\\" : \\\"you@your-company.com\\\"\\n    },\\n  " +
            "  \\\"license\\\" : {\\n      \\\"name\\\" : \\\"Apache 2.0\\\",\\n      \\\"url\\\" : " +
            "\\\"http://www.apache.org/licenses/LICENSE-2.0.html\\\"\\n    }\\n  },\\n  \\\"host\\\" :" +
            " \\\"virtserver.swaggerhub.com\\\",\\n  \\\"basePath\\\" : \\\"/yolom4/SimpleApi/1.0.0\\\",\\n " +
            " \\\"tags\\\" : [ {\\n    \\\"name\\\" : \\\"admins\\\",\\n    \\\"description\\\" : " +
            "\\\"Secured Admin-only calls\\\"\\n  }, {\\n    \\\"name\\\" : \\\"developers\\\",\\n " +
            "   \\\"description\\\" : \\\"Operations available to regular developers\\\"\\n  " +
            "} ],\\n  \\\"schemes\\\" : [ \\\"https\\\" ],\\n  \\\"paths\\\" : {\\n    " +
            "\\\"/inventory\\\" : {\\n      \\\"get\\\" : {\\n        \\\"tags\\\" : [" +
            " \\\"developers\\\" ],\\n        \\\"summary\\\" : \\\"searches inventory\\\",\\n   " +
            "     \\\"description\\\" : \\\"By passing in the appropriate options, you can search " +
            "for\\\\navailable inventory in the system\\\\n\\\",\\n        \\\"operationId\\\" :" +
            " \\\"searchInventory\\\",\\n        \\\"produces\\\" : [ \\\"application/json\\\" ],\\n  " +
            "      \\\"parameters\\\" : [ {\\n          \\\"name\\\" : \\\"searchString\\\",\\n       " +
            "   \\\"in\\\" : \\\"query\\\",\\n          \\\"description\\\" : \\\"pass an optional search " +
            "string for looking up inventory\\\",\\n          \\\"required\\\" : false,\\n        " +
            "  \\\"type\\\" : \\\"string\\\"\\n        }, {\\n          \\\"name\\\" : \\\"skip\\\",\\n  " +
            "        \\\"in\\\" : \\\"query\\\",\\n          \\\"description\\\" : \\\"number of records " +
            "to skip for pagination\\\",\\n          \\\"required\\\" : false,\\n          \\\"type\\\"" +
            " : \\\"integer\\\",\\n          \\\"minimum\\\" : 0,\\n          \\\"format\\\" : " +
            "\\\"int32\\\"\\n        }, {\\n          \\\"name\\\" : \\\"limit\\\",\\n         " +
            " \\\"in\\\" : \\\"query\\\",\\n          \\\"description\\\" : \\\"maximum number " +
            "of records to return\\\",\\n          \\\"required\\\" : false,\\n          \\\"type\\\" " +
            ": \\\"integer\\\",\\n          \\\"maximum\\\" : 50,\\n          \\\"minimum\\\" : 0,\\n  " +
            "        \\\"format\\\" : \\\"int32\\\"\\n        } ],\\n        \\\"responses\\\" : {\\n  " +
            "        \\\"200\\\" : {\\n            \\\"description\\\" : \\\"search results matching " +
            "criteria\\\",\\n            \\\"schema\\\" : {\\n              \\\"type\\\" : \\\"array\\\",\\n  " +
            "            \\\"items\\\" : {\\n                \\\"$ref\\\" : " +
            "\\\"#/definitions/InventoryItem\\\"\\n              }\\n            }\\n        " +
            "  },\\n          \\\"400\\\" : {\\n            \\\"description\\\" :" +
            " \\\"bad input parameter\\\"\\n          }\\n        }\\n      },\\n     " +
            " \\\"post\\\" : {\\n        \\\"tags\\\" : [ \\\"admins\\\" ],\\n       " +
            " \\\"summary\\\" : \\\"adds an inventory item\\\",\\n        \\\"description\\\" : " +
            "\\\"Adds an item to the system\\\",\\n        \\\"operationId\\\" : \\\"addInventory\\\",\\n " +
            "       \\\"consumes\\\" : [ \\\"application/json\\\" ],\\n        \\\"produces\\\" :" +
            " [ \\\"application/json\\\" ],\\n        \\\"parameters\\\" : [ {\\n          \\\"in\\\" : " +
            "\\\"body\\\",\\n          \\\"name\\\" : \\\"inventoryItem\\\",\\n         " +
            " \\\"description\\\" : \\\"Inventory item to add\\\",\\n          \\\"required\\\" : " +
            "false,\\n          \\\"schema\\\" : {\\n            \\\"$ref\\\" :" +
            " \\\"#/definitions/InventoryItem\\\"\\n          }\\n        } ],\\n       " +
            " \\\"responses\\\" : {\\n          \\\"201\\\" : {\\n            \\\"description\\\"" +
            " : \\\"item created\\\"\\n          },\\n          \\\"400\\\" : {\\n            " +
            "\\\"description\\\" : \\\"invalid input, object invalid\\\"\\n          },\\n   " +
            "       \\\"409\\\" : {\\n            \\\"description\\\" : \\\"an existing item" +
            " already exists\\\"\\n          }\\n        }\\n      }\\n    }\\n  },\\n " +
            " \\\"definitions\\\" : {\\n    \\\"InventoryItem\\\" : {\\n      \\\"type\\\" :" +
            " \\\"object\\\",\\n      \\\"required\\\" : [ \\\"id\\\", \\\"manufacturer\\\"," +
            " \\\"name\\\", \\\"releaseDate\\\" ],\\n      \\\"properties\\\" : {\\n        " +
            "\\\"id\\\" : {\\n          \\\"type\\\" : \\\"string\\\",\\n          \\\"format\\\" : " +
            "\\\"uuid\\\",\\n          \\\"example\\\" : \\\"d290f1ee-6c54-4b01-90e6-d701748f0851\\\"\\n   " +
            "     },\\n        \\\"name\\\" : {\\n          \\\"type\\\" : \\\"string\\\",\\n      " +
            "    \\\"example\\\" : \\\"Widget Adapter\\\"\\n        },\\n        \\\"releaseDate\\\" " +
            ": {\\n          \\\"type\\\" : \\\"string\\\",\\n          \\\"format\\\" :" +
            " \\\"int32\\\",\\n          \\\"example\\\" : \\\"2016-08-29T09:12:33.001Z\\\"\\n " +
            "       },\\n        \\\"manufacturer\\\" : {\\n          \\\"$ref\\\" : " +
            "\\\"#/definitions/Manufacturer\\\"\\n        }\\n      }\\n    },\\n    \\\"Manufacturer\\\" " +
            ": {\\n      \\\"required\\\" : [ \\\"name\\\" ],\\n      \\\"properties\\\" : {\\n      " +
            "  \\\"name\\\" : {\\n          \\\"type\\\" : \\\"string\\\",\\n          \\\"example\\\" " +
            ": \\\"ACME Corporation\\\"\\n        },\\n        \\\"homePage\\\" : {\\n          " +
            "\\\"type\\\" : \\\"string\\\",\\n          \\\"format\\\" : \\\"url\\\",\\n      " +
            "    \\\"example\\\" : \\\"https://www.acme-corp.com\\\"\\n        },\\n        " +
            "\\\"phone\\\" : {\\n          \\\"type\\\" : \\\"string\\\",\\n          \\\"example\\\" " +
            ": \\\"408-867-5309\\\"\\n        }\\n      }\\n    }\\n  }\\n}\",\"" +
            "isDefaultVersion\":false,\"transport\":[\"http\",\"https\"],\"tiers\":[\"Unlimited\"]" +
            ",\"visibility\":\"PUBLIC\",\"endpointConfig\":\"\"" +
            ",\"corsConfiguration\":{\"corsConfigurationEnabled\":false}}";

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
