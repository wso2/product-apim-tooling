package com.swagger.plugins.wso2;

import java.io.IOException;

public class Main {
    public static void main(String[] args) throws IOException {
        String testYaml = "{\n" +
                "   \"name\": \"TestAPI\",\n" +
                "   \"description\": \"This is a test API.\\r\\n\",\n" +
                "   \"context\": \"/test\",\n" +
                "   \"version\": \"1.0.0\",\n" +
                "   \"provider\": \"akilaaroshana@hotmail.com@ms9714\",\n" +
                "   \"apiDefinition\": \"{\\\"swagger\\\":\\\"2.0\\\",\\\"info\\\":{\\\"description\\\":\\\"This is a test API.\\\",\\\"version\\\":\\\"1.0.0\\\",\\\"title\\\":\\\"TestAPI\\\",\\\"contact\\\":{\\\"email\\\":\\\"you@your-company.com\\\"},\\\"license\\\":{\\\"name\\\":\\\"Apache 2.0\\\",\\\"url\\\":\\\"http://www.apache.org/licenses/LICENSE-2.0.html\\\"}},\\\"tags\\\":[{\\\"name\\\":\\\"admins\\\",\\\"description\\\":\\\"Secured Admin-only calls\\\"},{\\\"name\\\":\\\"developers\\\",\\\"description\\\":\\\"Operations available to regular developers\\\"}],\\\"paths\\\":{\\\"/inventory\\\":{\\\"post\\\":{\\\"tags\\\":[\\\"admins\\\"],\\\"summary\\\":\\\"adds an inventory item\\\",\\\"operationId\\\":\\\"addInventory\\\",\\\"description\\\":\\\"Adds an item to the system\\\",\\\"consumes\\\":[\\\"application/json\\\"],\\\"produces\\\":[\\\"application/json\\\"],\\\"parameters\\\":[{\\\"in\\\":\\\"body\\\",\\\"name\\\":\\\"inventoryItem\\\",\\\"description\\\":\\\"Inventory item to add\\\",\\\"schema\\\":{\\\"$ref\\\":\\\"#/definitions/InventoryItem\\\"}}],\\\"responses\\\":{\\\"201\\\":{\\\"description\\\":\\\"item created\\\"},\\\"400\\\":{\\\"description\\\":\\\"invalid input, object invalid\\\"},\\\"409\\\":{\\\"description\\\":\\\"an existing item already exists\\\"}}}}},\\\"definitions\\\":{\\\"InventoryItem\\\":{\\\"type\\\":\\\"object\\\",\\\"required\\\":[\\\"id\\\",\\\"name\\\",\\\"manufacturer\\\",\\\"releaseDate\\\"],\\\"properties\\\":{\\\"id\\\":{\\\"type\\\":\\\"string\\\",\\\"format\\\":\\\"uuid\\\",\\\"example\\\":\\\"d290f1ee-6c54-4b01-90e6-d701748f0851\\\"},\\\"name\\\":{\\\"type\\\":\\\"string\\\",\\\"example\\\":\\\"Widget Adapter\\\"},\\\"releaseDate\\\":{\\\"type\\\":\\\"string\\\",\\\"format\\\":\\\"int32\\\",\\\"example\\\":\\\"2016-08-29T09:12:33.001Z\\\"},\\\"manufacturer\\\":{\\\"$ref\\\":\\\"#/definitions/Manufacturer\\\"}}},\\\"Manufacturer\\\":{\\\"required\\\":[\\\"name\\\"],\\\"properties\\\":{\\\"name\\\":{\\\"type\\\":\\\"string\\\",\\\"example\\\":\\\"ACME Corporation\\\"},\\\"homePage\\\":{\\\"type\\\":\\\"string\\\",\\\"format\\\":\\\"url\\\",\\\"example\\\":\\\"https://www.acme-corp.com\\\"},\\\"phone\\\":{\\\"type\\\":\\\"string\\\",\\\"example\\\":\\\"408-867-5309\\\"}}}}}\",\n" +
                "   \"wsdlUri\": null,\n" +
                "   \"status\": \"CREATED\",\n" +
                "   \"responseCaching\": \"Disabled\",\n" +
                "   \"cacheTimeout\": 300,\n" +
                "   \"destinationStatsEnabled\": false,\n" +
                "   \"isDefaultVersion\": false,\n" +
                "   \"type\": \"HTTP\",\n" +
                "   \"transport\":    [\n" +
                "      \"http\",\n" +
                "      \"https\"\n" +
                "   ],\n" +
                "   \"tiers\": [\"Unlimited\"],\n" +
                "   \"maxTps\":    {\n" +
                "      \"sandbox\": 5000,\n" +
                "      \"production\": 1000\n" +
                "   },\n" +
                "   \"visibility\": \"PUBLIC\",\n" +
                "   \"visibleRoles\": [],\n" +
                "   \"visibleTenants\": [],\n" +
                "   \"endpointConfig\": \"\",\n" +
                "   \"endpointSecurity\":    {\n" +
                "      \"username\": \"user\",\n" +
                "      \"type\": \"basic\",\n" +
                "      \"password\": \"pass\"\n" +
                "   },\n" +
                "   \"gatewayEnvironments\": \"Production and Sandbox\",\n" +
                "   \"sequences\": [],\n" +
                "   \"subscriptionAvailability\": null,\n" +
                "   \"subscriptionAvailableTenants\": [],\n" +
                "   \"businessInformation\":    {\n" +
                "      \"businessOwnerEmail\": \"marketing@TEST.com\",\n" +
                "      \"technicalOwnerEmail\": \"architecture@TEST.com\",\n" +
                "      \"technicalOwner\": \"John Doe\",\n" +
                "      \"businessOwner\": \"Jane Roe\"\n" +
                "   },\n" +
                "   \"corsConfiguration\":    {\n" +
                "      \"accessControlAllowOrigins\": [\"*\"],\n" +
                "      \"accessControlAllowHeaders\":       [\n" +
                "         \"authorization\",\n" +
                "         \"Access-Control-Allow-Origin\",\n" +
                "         \"Content-Type\",\n" +
                "         \"SOAPAction\"\n" +
                "      ],\n" +
                "      \"accessControlAllowMethods\":       [\n" +
                "         \"GET\",\n" +
                "         \"PUT\",\n" +
                "         \"POST\",\n" +
                "         \"DELETE\",\n" +
                "         \"PATCH\",\n" +
                "         \"OPTIONS\"\n" +
                "      ],\n" +
                "      \"accessControlAllowCredentials\": false,\n" +
                "      \"corsConfigurationEnabled\": false\n" +
                "   }\n" +
                "}";;

        Wso2Api api = new Wso2Api();
        String accessToken = api.getAccessToken("akilaaroshana@hotmail.com","ms9714","Akilahotmail123");
        api.saveApi(testYaml,accessToken);
    }
}
