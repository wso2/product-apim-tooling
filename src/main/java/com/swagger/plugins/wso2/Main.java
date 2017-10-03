package com.swagger.plugins.wso2;

import org.json.simple.parser.ParseException;

import java.io.IOException;


/*****************************************************************
 *  Class name : Main
 * Methods : main
 * Functionality : For testing purposes
 * Visibility : Public
 * ****************************************************************/
public class Main {
    public static void main(String[] args) throws IOException, ParseException, PluginExecutionException {

        String swaggerDef2 = "{\"swagger\":\"2.0\",\"info\":{\"description\":\"This is a test API.\",\"version\":\"1.0.0\",\"title\":\"TestApi\",\"contact\":{\"email\":\"you@your-company.com\"},\"license\":{\"name\":\"Apache 2.0\",\"url\":\"http://www.apache.org/licenses/LICENSE-2.0.html\"}},\"tags\":[{\"name\":\"admins\",\"description\":\"Secured Admin-only calls\"},{\"name\":\"developers\",\"description\":\"Operations available to regular developers\"}],\"paths\":{\"/inventory\":{\"post\":{\"tags\":[\"admins\"],\"summary\":\"adds an inventory item\",\"operationId\":\"addInventory\",\"description\":\"Adds an item to the system\",\"consumes\":[\"application/json\"],\"produces\":[\"application/json\"],\"parameters\":[{\"in\":\"body\",\"name\":\"inventoryItem\",\"description\":\"Inventory item to add\",\"schema\":{\"$ref\":\"#/definitions/InventoryItem\"}}],\"responses\":{\"201\":{\"description\":\"item created\"},\"400\":{\"description\":\"invalid input, object invalid\"},\"409\":{\"description\":\"an existing item already exists\"}}}}},\"definitions\":{\"InventoryItem\":{\"type\":\"object\",\"required\":[\"id\",\"name\",\"manufacturer\",\"releaseDate\"],\"properties\":{\"id\":{\"type\":\"string\",\"format\":\"uuid\",\"example\":\"d290f1ee-6c54-4b01-90e6-d701748f0851\"},\"name\":{\"type\":\"string\",\"example\":\"Widget Adapter\"},\"releaseDate\":{\"type\":\"string\",\"format\":\"int32\",\"example\":\"2016-08-29T09:12:33.001Z\"},\"manufacturer\":{\"$ref\":\"#/definitions/Manufacturer\"}}},\"Manufacturer\":{\"required\":[\"name\"],\"properties\":{\"name\":{\"type\":\"string\",\"example\":\"ACME Corporation\"},\"homePage\":{\"type\":\"string\",\"format\":\"url\",\"example\":\"https://www.acme-corp.com\"},\"phone\":{\"type\":\"string\",\"example\":\"408-867-5309\"}}}}}";

//        PayloadConfiguration configuration = new PayloadConfiguration();
//        configuration.configurePayload("yolom@seekjobs4u.com","yolo4958",  swaggerDef2, "1.3.5", "/newContext");


        Wso2ApiGatewayPlugin plugin = new Wso2ApiGatewayPlugin();
        plugin.afterApiVersionSaved("yolom@seekjobs4u.com","yolo4958", "Yolofernando123", swaggerDef2, "5.2.4", "newContext");

//        Wso2Api api = new Wso2Api();
//        System.out.println(api.getClientIdAndSecret("akilaaroshana@hotmail.com","ms9714", "Akilahotmail123"));

    }
}
