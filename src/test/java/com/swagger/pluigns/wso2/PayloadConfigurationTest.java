package com.swagger.pluigns.wso2;

import com.swagger.plugins.wso2.PayloadConfiguration;
import io.swagger.models.Swagger;
import io.swagger.util.Json;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;

public class PayloadConfigurationTest {

    Swagger swagger;
    String swaggerDefinition;
    String expectedPayload;

    /**
     * This method initializes the variables before tests are executed
     *
     * @throws Exception
     */
    @Before
    public void init() throws Exception {

        swaggerDefinition = "swagger: '2.0'\n" +
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
                "basePath: /simple\n" +
                "schemes:\n" +
                " - https";

        swagger = Json.mapper().readValue(PayloadConfiguration.convertYamlToJson(swaggerDefinition), Swagger.class);
    }

    /**
     * This test method is to check whether the payload is returned correctly for creating the API in the cloud
     *
     * @throws Exception
     */
    @Test
    public void testConfigurePayload() throws Exception {

        expectedPayload = "{\"name\":\"Simple Inventory API\",\"context\":\"/simple\",\"version\":\"1.0.0\"," +
                "\"description\":\"This is a simple API\",\"apiDefinition\":\"{\\\"swagger\\\":\\\"2.0\\\"," +
                "\\\"info\\\":{\\\"description\\\":\\\"This is a simple API\\\",\\\"version\\\":\\\"1.0.0\\\"," +
                "\\\"title\\\":\\\"Simple Inventory API\\\",\\\"contact\\\":{\\\"email\\\":" +
                "\\\"you@your-company.com\\\"},\\\"license\\\":{\\\"name\\\":\\\"Apache 2.0\\\"," +
                "\\\"url\\\":\\\"http://www.apache.org/licenses/LICENSE-2.0.html\\\"}},\\\"host\\\":" +
                "\\\"virtserver.swaggerhub.com\\\",\\\"basePath\\\":\\\"/simple\\\",\\\"tags\\\":[{\\\"name\\\":" +
                "\\\"admins\\\",\\\"description\\\":\\\"Secured Admin-only calls\\\"},{\\\"name\\\":" +
                "\\\"developers\\\",\\\"description\\\":\\\"Operations available to regular developers\\\"}]," +
                "\\\"schemes\\\":[\\\"https\\\"],\\\"paths\\\":{\\\"/inventory\\\":{\\\"get\\\":{\\\"tags\\\":" +
                "[\\\"developers\\\"],\\\"summary\\\":\\\"searches inventory\\\",\\\"description\\\":" +
                "\\\"By passing in the appropriate options, you can search for\\\\navailable inventory in the " +
                "system\\\\n\\\",\\\"operationId\\\":\\\"searchInventory\\\",\\\"produces\\\":" +
                "[\\\"application/json\\\"],\\\"parameters\\\":[{\\\"name\\\":\\\"searchString\\\",\\\"in\\\":" +
                "\\\"query\\\",\\\"description\\\":\\\"pass an optional search string for looking up inventory\\\"," +
                "\\\"required\\\":false,\\\"type\\\":\\\"string\\\"},{\\\"name\\\":\\\"skip\\\",\\\"in\\\":" +
                "\\\"query\\\",\\\"description\\\":\\\"number of records to skip for pagination\\\"," +
                "\\\"required\\\":false,\\\"type\\\":\\\"integer\\\",\\\"minimum\\\":0,\\\"format\\\":" +
                "\\\"int32\\\"},{\\\"name\\\":\\\"limit\\\",\\\"in\\\":\\\"query\\\",\\\"description\\\":" +
                "\\\"maximum number of records to return\\\",\\\"required\\\":false,\\\"type\\\":\\\"integer\\\"," +
                "\\\"maximum\\\":50,\\\"minimum\\\":0,\\\"format\\\":\\\"int32\\\"}],\\\"responses\\\":{\\\"200\\\":" +
                "{\\\"description\\\":\\\"search results matching criteria\\\",\\\"schema\\\":{\\\"type\\\":" +
                "\\\"array\\\",\\\"items\\\":{\\\"$ref\\\":\\\"#/definitions/InventoryItem\\\"}}},\\\"400\\\":" +
                "{\\\"description\\\":\\\"bad input parameter\\\"}}},\\\"post\\\":{\\\"tags\\\":[\\\"admins\\\"]," +
                "\\\"summary\\\":\\\"adds an inventory item\\\",\\\"description\\\":" +
                "\\\"Adds an item to the system\\\",\\\"operationId\\\":\\\"addInventory\\\",\\\"consumes\\\":" +
                "[\\\"application/json\\\"],\\\"produces\\\":[\\\"application/json\\\"],\\\"parameters\\\":" +
                "[{\\\"in\\\":\\\"body\\\",\\\"name\\\":\\\"inventoryItem\\\",\\\"description\\\":" +
                "\\\"Inventory item to add\\\",\\\"required\\\":false,\\\"schema\\\":{\\\"$ref\\\":" +
                "\\\"#/definitions/InventoryItem\\\"}}],\\\"responses\\\":{\\\"201\\\":{\\\"description\\\":" +
                "\\\"item created\\\"},\\\"400\\\":{\\\"description\\\":\\\"invalid input, object invalid\\\"}," +
                "\\\"409\\\":{\\\"description\\\":\\\"an existing item already exists\\\"}}}}},\\\"definitions\\\":" +
                "{\\\"InventoryItem\\\":{\\\"type\\\":\\\"object\\\",\\\"required\\\":[\\\"id\\\"," +
                "\\\"manufacturer\\\",\\\"name\\\",\\\"releaseDate\\\"],\\\"properties\\\":{\\\"id\\\":" +
                "{\\\"type\\\":\\\"string\\\",\\\"format\\\":\\\"uuid\\\",\\\"example\\\":" +
                "\\\"d290f1ee-6c54-4b01-90e6-d701748f0851\\\"},\\\"name\\\":{\\\"type\\\":\\\"string\\\"," +
                "\\\"example\\\":\\\"Widget Adapter\\\"},\\\"releaseDate\\\":{\\\"type\\\":\\\"string\\\"," +
                "\\\"format\\\":\\\"int32\\\",\\\"example\\\":\\\"2016-08-29T09:12:33.001Z\\\"},\\\"manufacturer\\\":" +
                "{\\\"$ref\\\":\\\"#/definitions/Manufacturer\\\"}}},\\\"Manufacturer\\\":" +
                "{\\\"required\\\":[\\\"name\\\"],\\\"properties\\\":{\\\"name\\\":{\\\"type\\\":\\" +
                "\"string\\\",\\\"example\\\":\\\"ACME Corporation\\\"},\\\"homePage\\\":{\\\"type\\\":" +
                "\\\"string\\\",\\\"format\\\":\\\"url\\\",\\\"example\\\":\\\"https://www.acme-corp.com\\\"}," +
                "\\\"phone\\\":{\\\"type\\\":\\\"string\\\",\\\"example\\\":\\\"408-867-5309\\\"}}}}}\"," +
                "\"isDefaultVersion\":false,\"transport\":[\"http\",\"https\"],\"tiers\":[\"Unlimited\"]," +
                "\"visibility\":\"PUBLIC\",\"endpointConfig\":\"\",\"corsConfiguration\":" +
                "{\"corsConfigurationEnabled\":false}}";

        PayloadConfiguration payloadConfiguration = new PayloadConfiguration();
        String payload = payloadConfiguration.configurePayload(swagger);
        Assert.assertEquals(expectedPayload, payload);
    }

    /**
     * This test method is to check whether a swagger yaml is correctly converted into a json
     *
     * @throws Exception
     */
    @Test
    public void testConvertYamlToJson() throws Exception {

        String expectedJson = "{\"swagger\":\"2.0\",\"info\":{\"description\":\"This is a simple API\",\"version\":" +
                "\"1.0.0\",\"title\":\"Simple Inventory API\",\"contact\":{\"email\":\"you@your-company.com\"}," +
                "\"license\":{\"name\":\"Apache 2.0\",\"url\":\"http://www.apache.org/licenses/LICENSE-2.0.html\"}}," +
                "\"host\":\"virtserver.swaggerhub.com\",\"basePath\":\"/simple\",\"tags\":[{\"name\":\"admins\"," +
                "\"description\":\"Secured Admin-only calls\"},{\"name\":\"developers\",\"description\":" +
                "\"Operations available to regular developers\"}],\"schemes\":[\"https\"],\"paths\":" +
                "{\"/inventory\":{\"get\":{\"tags\":[\"developers\"],\"summary\":\"searches inventory\"," +
                "\"description\":\"By passing in the appropriate options, you can search for\\navailable" +
                " inventory in the system\\n\",\"operationId\":\"searchInventory\",\"produces\":" +
                "[\"application/json\"],\"parameters\":[{\"name\":\"searchString\",\"in\":\"query\"," +
                "\"description\":\"pass an optional search string for looking up inventory\",\"required\":" +
                "false,\"type\":\"string\"},{\"name\":\"skip\",\"in\":\"query\",\"description\":\"number of" +
                " records to skip for pagination\",\"required\":false,\"type\":\"integer\",\"minimum\":" +
                "0,\"format\":\"int32\"},{\"name\":\"limit\",\"in\":\"query\",\"description\":" +
                "\"maximum number of records to return\",\"required\":false,\"type\":\"integer\",\"maximum\":5" +
                "0,\"minimum\":0,\"format\":\"int32\"}],\"responses\":{\"200\":{\"description\":\"" +
                "search results matching criteria\",\"schema\":{\"type\":\"array\",\"items\":{\"$ref\":\"" +
                "#/definitions/InventoryItem\"}}},\"400\":{\"description\":\"bad input parameter\"}}},\"post\"" +
                ":{\"tags\":[\"admins\"],\"summary\":\"adds an inventory item\",\"description\":\"Adds an ite" +
                "m to the system\",\"operationId\":\"addInventory\",\"consumes\":[\"application/json\"],\"pro" +
                "duces\":[\"application/json\"],\"parameters\":[{\"in\":\"body\",\"name\":\"inventoryItem\",\"" +
                "description\":\"Inventory item to add\",\"required\":false,\"schema\":{\"$ref\":\"#/definitio" +
                "ns/InventoryItem\"}}],\"responses\":{\"201\":{\"description\":\"item created\"},\"400\":{\"de" +
                "scription\":\"invalid input, object invalid\"},\"409\":{\"description\":\"an existing item al" +
                "ready exists\"}}}}},\"definitions\":{\"InventoryItem\":{\"type\":\"object\",\"required\":[\"i" +
                "d\",\"manufacturer\",\"name\",\"releaseDate\"],\"properties\":{\"id\":{\"type\":\"string\",\"" +
                "format\":\"uuid\",\"example\":\"d290f1ee-6c54-4b01-90e6-d701748f0851\"},\"name\":{\"type\":\"" +
                "string\",\"example\":\"Widget Adapter\"},\"releaseDate\":{\"type\":\"string\",\"format\":\"in" +
                "t32\",\"example\":\"2016-08-29T09:12:33.001Z\"},\"manufacturer\":{\"$ref\":\"#/definitions/Ma" +
                "nufacturer\"}}},\"Manufacturer\":{\"required\":[\"name\"],\"properties\":{\"name\":{\"type\":" +
                "\"string\",\"example\":\"ACME Corporation\"},\"homePage\":{\"type\":\"string\",\"format\":\"u" +
                "rl\",\"example\":\"https://www.acme-corp.com\"},\"phone\":{\"type\":\"string\",\"example\":\"4" +
                "08-867-5309\"}}}}}";

        String json = PayloadConfiguration.convertYamlToJson(swaggerDefinition);
        Assert.assertEquals(expectedJson, json);

    }

}
