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

        String swaggerDef2 = "{\"swagger\":\"2.0\",\"info\":{\"description\":\"This is a test API.\"," +
                "\"version\":\"1.0.0\",\"title\":\"TestApi\",\"contact\":{\"email\":\"you@your-company.com\"}," +
                "\"license\":{\"name\":\"Apache 2.0\",\"url\":\"http://www.apache.org/licenses/LICENSE-2.0.html\"}}," +
                "\"tags\":[{\"name\":\"admins\",\"description\":\"Secured Admin-only calls\"}" +
                ",{\"name\":\"developers\",\"description\":\"Operations available to regular developers\"}]" +
                ",\"paths\":{\"/inventory\":{\"post\":{\"tags\":[\"admins\"],\"summary\":\"adds an inventory item\"" +
                ",\"operationId\":\"addInventory\",\"description\":\"Adds an item to the system\"" +
                ",\"consumes\":[\"application/json\"],\"produces\":[\"application/json\"]," +
                "\"parameters\":[{\"in\":\"body\",\"name\":\"inventoryItem\"," +
                "\"description\":\"Inventory item to add\",\"schema\":{\"$ref\":\"#/definitions/InventoryItem\"}}]," +
                "\"responses\":{\"201\":{\"description\":\"item created\"},\"400\":{\"description\":\"invalid input, " +
                "object invalid\"},\"409\":{\"description\":\"an existing item already exists\"}}}}}," +
                "\"definitions\":{\"InventoryItem\":{\"type\":\"object\",\"required\":[\"id\",\"name\"," +
                "\"manufacturer\",\"releaseDate\"],\"properties\":{\"id\":{\"type\":\"string\"," +
                "\"format\":\"uuid\",\"example\":\"d290f1ee-6c54-4b01-90e6-d701748f0851\"}," +
                "\"name\":{\"type\":\"string\",\"example\":\"Widget Adapter\"}," +
                "\"releaseDate\":{\"type\":\"string\",\"format\":\"int32\"," +
                "\"example\":\"2016-08-29T09:12:33.001Z\"}," +
                "\"manufacturer\":{\"$ref\":\"#/definitions/Manufacturer\"}}}," +
                "\"Manufacturer\":{\"required\":[\"name\"],\"properties\":{\"name\":{\"type\":\"string\"," +
                "\"example\":\"ACME Corporation\"},\"homePage\":{\"type\":\"string\",\"format\":\"url\"," +
                "\"example\":\"https://www.acme-corp.com\"},\"phone\":{\"type\":\"string\"," +
                "\"example\":\"408-867-5309\"}}}}}";

        String swaggerYaml =
                "swagger: '2.0'\n" +
                "info:\n" +
                "  description: This is a test API.\n" +
                "  version: 1.0.0\n" +
                "  title: TestApi\n" +
                "  contact:\n" +
                "    email: you@your-company.com\n" +
                "  license:\n" +
                "    name: Apache 2.0\n" +
                "    url: http://www.apache.org/licenses/LICENSE-2.0.html\n" +
                "tags:\n" +
                "- name: admins\n" +
                "  description: Secured Admin-only calls\n" +
                "- name: developers\n" +
                "  description: Operations available to regular developers\n" +
                "paths:\n" +
                "  \"/inventory\":\n" +
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
                "          \"$ref\": \"#/definitions/InventoryItem\"\n" +
                "      responses:\n" +
                "        '201':\n" +
                "          description: item created\n" +
                "        '400':\n" +
                "          description: invalid input, object invalid\n" +
                "        '409':\n" +
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
                "        example: '2016-08-29T09:12:33.001Z'\n" +
                "      manufacturer:\n" +
                "        \"$ref\": \"#/definitions/Manufacturer\"\n" +
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
                "        example: https://www.acme-corp.com\n" +
                "      phone:\n" +
                "        type: string\n" +
                "        example: 408-867-5309\n";

        String randomSwagger = "swagger: '2.0'\n" +
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
                "        example: 408-867-5309";

        String randomSwagger2 = "swagger: '2.0'\n" +
                "info:\n" +
                "  version: 1.0.0\n" +
                "  title: home-iot-api\n" +
                "  description: The API for the EatBacon IOT project\n" +
                "\n" +
                "consumes:\n" +
                "  - application/json\n" +
                "produces:\n" +
                "  - application/json\n" +
                "paths:\n" +
                "  /devices:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - Device\n" +
                "      description: returns all registered devices\n" +
                "      operationId: getDevices\n" +
                "      parameters:\n" +
                "      - in: query\n" +
                "        name: skip\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "        description: number of records to skip\n" +
                "      - in: query\n" +
                "        name: limit\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "        description: max number of records to return\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: All the devices\n" +
                "          schema:\n" +
                "            type: array\n" +
                "            items:\n" +
                "              type: string\n" +
                "              format: uri\n" +
                "              example: http://10.0.0.225:8080\n" +
                "    post:\n" +
                "      tags:\n" +
                "      - Device\n" +
                "      operationId: register\n" +
                "      parameters:\n" +
                "      - in: body\n" +
                "        name: device\n" +
                "        required: false\n" +
                "        schema:\n" +
                "          $ref: '#/definitions/DeviceRegistrationInfo'\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: successfully registered device\n" +
                "\n" +
                "  /lighting/dimmers/{deviceId}/{value}:\n" +
                "    post:\n" +
                "      tags:\n" +
                "      - \"Z-Wave\"\n" +
                "      operationId: \"setDimmer\"\n" +
                "      parameters:\n" +
                "      - name: \"deviceId\"\n" +
                "        in: \"path\"\n" +
                "        required: true\n" +
                "        type: \"string\"\n" +
                "      - name: \"value\"\n" +
                "        in: \"path\"\n" +
                "        required: true\n" +
                "        type: \"integer\"\n" +
                "        maximum: 100.0\n" +
                "        minimum: 0.0\n" +
                "        format: \"int32\"\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: \"response\"\n" +
                "          schema:\n" +
                "            $ref: \"#/definitions/ApiResponse\"\n" +
                "      x-swagger-router-controller: \"ZWave\"\n" +
                "  /lighting/dimmers/{deviceId}/{value}/timer/{timeunit}:\n" +
                "    post:\n" +
                "      tags:\n" +
                "      - Z-Wave\n" +
                "      description: sets a dimmer to a specific value on a timer\n" +
                "      operationId: setDimmerTimer\n" +
                "      parameters:\n" +
                "      - name: deviceId\n" +
                "        in: path\n" +
                "        required: true\n" +
                "        type: string\n" +
                "      - name: value\n" +
                "        in: path\n" +
                "        required: true\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "      - name: timeunit\n" +
                "        in: path\n" +
                "        required: true\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "      - name: units\n" +
                "        in: query\n" +
                "        required: false\n" +
                "        type: string\n" +
                "        default: milliseconds\n" +
                "        enum:\n" +
                "        - seconds\n" +
                "        - minutes\n" +
                "        - milliseconds\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: response\n" +
                "          schema:\n" +
                "            $ref: '#/definitions/ApiResponse'\n" +
                "      x-swagger-router-controller: ZWave\n" +
                "  /lighting/switches/{deviceId}:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - Z-Wave\n" +
                "      operationId: getSwitchState\n" +
                "      parameters:\n" +
                "      - name: deviceId\n" +
                "        in: path\n" +
                "        required: true\n" +
                "        type: string\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: response\n" +
                "          schema:\n" +
                "            $ref: '#/definitions/DeviceState'\n" +
                "      x-swagger-router-controller: ZWave\n" +
                "  /lighting/switches/{deviceId}/{value}:\n" +
                "    post:\n" +
                "      tags:\n" +
                "      - Z-Wave\n" +
                "      operationId: setSwitch\n" +
                "      parameters:\n" +
                "      - name: deviceId\n" +
                "        in: path\n" +
                "        required: true\n" +
                "        type: string\n" +
                "      - name: value\n" +
                "        in: path\n" +
                "        required: true\n" +
                "        type: string\n" +
                "        enum:\n" +
                "        - true\n" +
                "        - false\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: response\n" +
                "          schema:\n" +
                "            $ref: '#/definitions/ApiResponse'\n" +
                "      x-swagger-router-controller: ZWave\n" +
                "  /lighting/switches/{deviceId}/{value}/timer/{minutes}:\n" +
                "    post:\n" +
                "      tags:\n" +
                "      - Z-Wave\n" +
                "      description: sets a switch to a specific value on a timer\n" +
                "      operationId: setSwitchTimer\n" +
                "      parameters:\n" +
                "      - name: deviceId\n" +
                "        in: path\n" +
                "        required: true\n" +
                "        type: string\n" +
                "      - name: value\n" +
                "        in: path\n" +
                "        required: true\n" +
                "        type: string\n" +
                "        enum:\n" +
                "        - true\n" +
                "        - false\n" +
                "      - name: minutes\n" +
                "        in: path\n" +
                "        required: true\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: response\n" +
                "          schema:\n" +
                "            $ref: '#/definitions/ApiResponse'\n" +
                "      x-swagger-router-controller: ZWave\n" +
                "  /lightingSummary:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - Z-Wave\n" +
                "      operationId: getLightingSummary\n" +
                "      parameters: []\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: ok\n" +
                "          schema:\n" +
                "            $ref: '#/definitions/LightingSummary'\n" +
                "      x-swagger-router-controller: ZWave\n" +
                "  /temperature:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - Environment\n" +
                "      operationId: temperatureSummary\n" +
                "      parameters: []\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: ok\n" +
                "          schema:\n" +
                "            $ref: '#/definitions/TemperatureSummary'\n" +
                "      x-swagger-router-controller: Environment\n" +
                "  /temperature/forecast/{days}:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - Environment\n" +
                "      operationId: getForecast\n" +
                "      parameters:\n" +
                "      - name: days\n" +
                "        in: path\n" +
                "        required: true\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: the forecast\n" +
                "          schema:\n" +
                "            $ref: '#/definitions/ForecastResponse'\n" +
                "      x-swagger-router-controller: Environment\n" +
                "  /temperature/{zoneId}:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - Environment\n" +
                "      operationId: getZoneTemperature\n" +
                "      parameters:\n" +
                "      - name: zoneId\n" +
                "        in: path\n" +
                "        required: true\n" +
                "        type: string\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: Zone temperature\n" +
                "          schema:\n" +
                "            $ref: '#/definitions/TemperatueZoneStatus'\n" +
                "      x-swagger-router-controller: Environment\n" +
                "  /temperature/{zoneId}/heater:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - Environment\n" +
                "      description: gets the state of the heater\n" +
                "      operationId: getHeaterState\n" +
                "      parameters:\n" +
                "      - name: zoneId\n" +
                "        in: path\n" +
                "        required: true\n" +
                "        type: string\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: heater state\n" +
                "          schema:\n" +
                "            $ref: '#/definitions/HeaterState'\n" +
                "      x-swagger-router-controller: Environment\n" +
                "  /temperature/{zoneId}/heater/{state}:\n" +
                "    post:\n" +
                "      tags:\n" +
                "      - Environment\n" +
                "      description: turns the heater on or off\n" +
                "      operationId: setHeaterState\n" +
                "      parameters:\n" +
                "      - name: zoneId\n" +
                "        in: path\n" +
                "        required: true\n" +
                "        type: string\n" +
                "      - name: state\n" +
                "        in: path\n" +
                "        required: true\n" +
                "        type: string\n" +
                "        enum:\n" +
                "        - false\n" +
                "        - true\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: Status of the operation\n" +
                "          schema:\n" +
                "            $ref: '#/definitions/ApiResponse'\n" +
                "      x-swagger-router-controller: Environment\n" +
                "  /zones:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - Zones\n" +
                "      operationId: getZones\n" +
                "      parameters: []\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: ok\n" +
                "          schema:\n" +
                "            type: array\n" +
                "            items:\n" +
                "              type: string\n" +
                "      x-swagger-router-controller: Zones\n" +
                "  /zones/{zoneId}/quiet:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - Zones\n" +
                "      operationId: quietZone\n" +
                "      parameters:\n" +
                "      - name: zoneId\n" +
                "        in: path\n" +
                "        required: true\n" +
                "        type: string\n" +
                "        enum:\n" +
                "        - basement\n" +
                "        - first-floor\n" +
                "        - second-floor\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: ok\n" +
                "      x-swagger-router-controller: Zones\n" +
                "definitions:\n" +
                "  LightingSummary:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      zones:\n" +
                "        type: array\n" +
                "        items:\n" +
                "          $ref: '#/definitions/LightingZone'\n" +
                "      zoneStatus:\n" +
                "        type: array\n" +
                "        items:\n" +
                "          $ref: '#/definitions/LightingZoneStatus'\n" +
                "    description: ok\n" +
                "  LightingZone:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      id:\n" +
                "        type: string\n" +
                "      name:\n" +
                "        type: string\n" +
                "      deviceId:\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "      deviceType:\n" +
                "        type: string\n" +
                "        enum:\n" +
                "        - dimmer\n" +
                "        - switch\n" +
                "      zone:\n" +
                "        type: string\n" +
                "  LightingZoneStatus:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      id:\n" +
                "        type: string\n" +
                "      name:\n" +
                "        type: string\n" +
                "      lastUpdate:\n" +
                "        type: string\n" +
                "        format: date-time\n" +
                "      level:\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "    description: the status of the lighting zone.\n" +
                "  TemperatureSummary:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      zones:\n" +
                "        type: array\n" +
                "        items:\n" +
                "          $ref: '#/definitions/TemperatureZone'\n" +
                "      zoneStatus:\n" +
                "        type: array\n" +
                "        items:\n" +
                "          $ref: '#/definitions/TemperatueZoneStatus'\n" +
                "    description: ok\n" +
                "  TemperatureZone:\n" +
                "    type: object\n" +
                "    required:\n" +
                "    - id\n" +
                "    - name\n" +
                "    properties:\n" +
                "      id:\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "        description: the unique identifier for the zone\n" +
                "      name:\n" +
                "        type: string\n" +
                "      inputPosition:\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "      outputPosition:\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "      zone:\n" +
                "        type: string\n" +
                "    description: a single temperature zone\n" +
                "  TemperatueZoneStatus:\n" +
                "    type: object\n" +
                "    required:\n" +
                "    - id\n" +
                "    - timestamp\n" +
                "    - value\n" +
                "    properties:\n" +
                "      id:\n" +
                "        type: string\n" +
                "        description: the unique identifier for the zone\n" +
                "      name:\n" +
                "        type: string\n" +
                "        description: the name of the zone\n" +
                "      value:\n" +
                "        type: number\n" +
                "        format: double\n" +
                "        description: the temperature in the zone\n" +
                "      units:\n" +
                "        type: string\n" +
                "        description: the temperature units\n" +
                "        enum:\n" +
                "        - celcius\n" +
                "        - fahrenheit\n" +
                "        default: fahrenheit\n" +
                "      timestamp:\n" +
                "        type: string\n" +
                "        format: date-time\n" +
                "        description: the timestamp when the temperature was measured\n" +
                "    description: status of a single zone\n" +
                "  ApiResponse:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      code:\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "      message:\n" +
                "        type: string\n" +
                "        example: 'everything is ok'\n" +
                "  HeaterState:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      id:\n" +
                "        type: string\n" +
                "      state:\n" +
                "        type: string\n" +
                "  DeviceState:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      id:\n" +
                "        type: string\n" +
                "      name:\n" +
                "        type: string\n" +
                "      lastUpdate:\n" +
                "        type: string\n" +
                "        format: date-time\n" +
                "      level:\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "  ForecastResponse:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      city:\n" +
                "        $ref: '#/definitions/City'\n" +
                "      values:\n" +
                "        type: array\n" +
                "        items:\n" +
                "          $ref: '#/definitions/Forecast'\n" +
                "  Forecast:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      date:\n" +
                "        type: string\n" +
                "        format: date-time\n" +
                "      pressure:\n" +
                "        type: number\n" +
                "        format: double\n" +
                "      humidity:\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "      windSpeed:\n" +
                "        type: number\n" +
                "        format: double\n" +
                "      clouds:\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "      temperature:\n" +
                "        $ref: '#/definitions/ForecastTemperature'\n" +
                "      weather:\n" +
                "        $ref: '#/definitions/WeatherForecast'\n" +
                "  City:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      id:\n" +
                "        type: integer\n" +
                "        format: int64\n" +
                "      name:\n" +
                "        type: string\n" +
                "      lat:\n" +
                "        type: number\n" +
                "        format: double\n" +
                "      lon:\n" +
                "        type: number\n" +
                "        format: double\n" +
                "      country:\n" +
                "        type: string\n" +
                "  ForecastTemperature:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      low:\n" +
                "        type: number\n" +
                "        format: double\n" +
                "      high:\n" +
                "        type: number\n" +
                "        format: double\n" +
                "      morning:\n" +
                "        type: number\n" +
                "        format: double\n" +
                "      day:\n" +
                "        type: number\n" +
                "        format: double\n" +
                "      evening:\n" +
                "        type: number\n" +
                "        format: double\n" +
                "      night:\n" +
                "        type: number\n" +
                "        format: double\n" +
                "  WeatherForecast:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      summary:\n" +
                "        type: string\n" +
                "      description:\n" +
                "        type: string\n" +
                "      icon:\n" +
                "        type: string\n" +
                "  DeviceRegistrationInfo:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      uri:\n" +
                "        type: string\n" +
                "        format: uri\n" +
                "        example: http://10.0.0.220:8080\n" +
                "      id:\n" +
                "        type: string\n" +
                "        format: uuid\n" +
                "        example: '0729a580-2240-11e6-9eb5-0002a5d5c51b'";

        String petStore = "swagger: '2.0'\n" +
                "info:\n" +
                "  description: |\n" +
                "    This is a sample server Petstore server.  You can find \n" +
                "    out more about Swagger at \n" +
                "    [http://swagger.io](http://swagger.io) or on \n" +
                "    [irc.freenode.net, #swagger](http://swagger.io/irc/).\n" +
                "  version: 1.0.0\n" +
                "  title: Swagger Petstore\n" +
                "  termsOfService: http://swagger.io/terms/\n" +
                "  contact:\n" +
                "    email: apiteam@swagger.io\n" +
                "  license:\n" +
                "    name: Apache 2.0\n" +
                "    url: http://www.apache.org/licenses/LICENSE-2.0.html\n" +
                "host: petstore.swagger.io\n" +
                "basePath: /v2\n" +
                "tags:\n" +
                "- name: pet\n" +
                "  description: Everything about your Pets\n" +
                "  externalDocs:\n" +
                "    description: Find out more\n" +
                "    url: http://swagger.io\n" +
                "- name: store\n" +
                "  description: Access to Petstore orders\n" +
                "- name: user\n" +
                "  description: Operations about user\n" +
                "  externalDocs:\n" +
                "    description: Find out more about our store\n" +
                "    url: http://swagger.io\n" +
                "schemes:\n" +
                "- http\n" +
                "paths:\n" +
                "  /pet:\n" +
                "    post:\n" +
                "      tags:\n" +
                "      - pet\n" +
                "      summary: Add a new pet to the store\n" +
                "      operationId: addPet\n" +
                "      consumes:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters:\n" +
                "      - in: body\n" +
                "        name: body\n" +
                "        description: Pet object that needs to be added to the store\n" +
                "        required: true\n" +
                "        schema:\n" +
                "          $ref: '#/definitions/Pet'\n" +
                "      responses:\n" +
                "        405:\n" +
                "          description: Invalid input\n" +
                "      security:\n" +
                "      - petstore_auth:\n" +
                "        - write:pets\n" +
                "        - read:pets\n" +
                "    put:\n" +
                "      tags:\n" +
                "      - pet\n" +
                "      summary: Update an existing pet\n" +
                "      operationId: updatePet\n" +
                "      consumes:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters:\n" +
                "      - in: body\n" +
                "        name: body\n" +
                "        description: Pet object that needs to be added to the store\n" +
                "        required: true\n" +
                "        schema:\n" +
                "          $ref: '#/definitions/Pet'\n" +
                "      responses:\n" +
                "        400:\n" +
                "          description: Invalid ID supplied\n" +
                "        404:\n" +
                "          description: Pet not found\n" +
                "        405:\n" +
                "          description: Validation exception\n" +
                "      security:\n" +
                "      - petstore_auth:\n" +
                "        - write:pets\n" +
                "        - read:pets\n" +
                "  /pet/findByStatus:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - pet\n" +
                "      summary: Finds Pets by status\n" +
                "      description: Multiple status values can be provided with comma separated strings\n" +
                "      operationId: findPetsByStatus\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters:\n" +
                "      - name: status\n" +
                "        in: query\n" +
                "        description: Status values that need to be considered for filter\n" +
                "        required: true\n" +
                "        type: array\n" +
                "        items:\n" +
                "          type: string\n" +
                "          enum:\n" +
                "          - available\n" +
                "          - pending\n" +
                "          - sold\n" +
                "          default: available\n" +
                "        collectionFormat: multi\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: successful operation\n" +
                "          schema:\n" +
                "            type: array\n" +
                "            items:\n" +
                "              $ref: '#/definitions/Pet'\n" +
                "        400:\n" +
                "          description: Invalid status value\n" +
                "      security:\n" +
                "      - petstore_auth:\n" +
                "        - write:pets\n" +
                "        - read:pets\n" +
                "  /pet/findByTags:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - pet\n" +
                "      summary: Finds Pets by tags\n" +
                "      description: Muliple tags can be provided with comma separated strings. Use\\\n" +
                "        \\ tag1, tag2, tag3 for testing.\n" +
                "      operationId: findPetsByTags\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters:\n" +
                "      - name: tags\n" +
                "        in: query\n" +
                "        description: Tags to filter by\n" +
                "        required: true\n" +
                "        type: array\n" +
                "        items:\n" +
                "          type: string\n" +
                "        collectionFormat: multi\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: successful operation\n" +
                "          schema:\n" +
                "            type: array\n" +
                "            items:\n" +
                "              $ref: '#/definitions/Pet'\n" +
                "        400:\n" +
                "          description: Invalid tag value\n" +
                "      security:\n" +
                "      - petstore_auth:\n" +
                "        - write:pets\n" +
                "        - read:pets\n" +
                "      deprecated: true\n" +
                "  /pet/{petId}:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - pet\n" +
                "      summary: Find pet by ID\n" +
                "      description: Returns a single pet\n" +
                "      operationId: getPetById\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters:\n" +
                "      - name: petId\n" +
                "        in: path\n" +
                "        description: ID of pet to return\n" +
                "        required: true\n" +
                "        type: integer\n" +
                "        format: int64\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: successful operation\n" +
                "          schema:\n" +
                "            $ref: '#/definitions/Pet'\n" +
                "        400:\n" +
                "          description: Invalid ID supplied\n" +
                "        404:\n" +
                "          description: Pet not found\n" +
                "      security:\n" +
                "      - api_key: []\n" +
                "    post:\n" +
                "      tags:\n" +
                "      - pet\n" +
                "      summary: Updates a pet in the store with form data\n" +
                "      operationId: updatePetWithForm\n" +
                "      consumes:\n" +
                "      - application/x-www-form-urlencoded\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters:\n" +
                "      - name: petId\n" +
                "        in: path\n" +
                "        description: ID of pet that needs to be updated\n" +
                "        required: true\n" +
                "        type: integer\n" +
                "        format: int64\n" +
                "      - name: name\n" +
                "        in: formData\n" +
                "        description: Updated name of the pet\n" +
                "        required: false\n" +
                "        type: string\n" +
                "      - name: status\n" +
                "        in: formData\n" +
                "        description: Updated status of the pet\n" +
                "        required: false\n" +
                "        type: string\n" +
                "      responses:\n" +
                "        405:\n" +
                "          description: Invalid input\n" +
                "      security:\n" +
                "      - petstore_auth:\n" +
                "        - write:pets\n" +
                "        - read:pets\n" +
                "    delete:\n" +
                "      tags:\n" +
                "      - pet\n" +
                "      summary: Deletes a pet\n" +
                "      operationId: deletePet\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters:\n" +
                "      - name: api_key\n" +
                "        in: header\n" +
                "        required: false\n" +
                "        type: string\n" +
                "      - name: petId\n" +
                "        in: path\n" +
                "        description: Pet id to delete\n" +
                "        required: true\n" +
                "        type: integer\n" +
                "        format: int64\n" +
                "      responses:\n" +
                "        400:\n" +
                "          description: Invalid ID supplied\n" +
                "        404:\n" +
                "          description: Pet not found\n" +
                "      security:\n" +
                "      - petstore_auth:\n" +
                "        - write:pets\n" +
                "        - read:pets\n" +
                "  /pet/{petId}/uploadImage:\n" +
                "    post:\n" +
                "      tags:\n" +
                "      - pet\n" +
                "      summary: uploads an image\n" +
                "      operationId: uploadFile\n" +
                "      consumes:\n" +
                "      - multipart/form-data\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      parameters:\n" +
                "      - name: petId\n" +
                "        in: path\n" +
                "        description: ID of pet to update\n" +
                "        required: true\n" +
                "        type: integer\n" +
                "        format: int64\n" +
                "      - name: additionalMetadata\n" +
                "        in: formData\n" +
                "        description: Additional data to pass to server\n" +
                "        required: false\n" +
                "        type: string\n" +
                "      - name: file\n" +
                "        in: formData\n" +
                "        description: file to upload\n" +
                "        required: false\n" +
                "        type: file\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: successful operation\n" +
                "          schema:\n" +
                "            $ref: '#/definitions/ApiResponse'\n" +
                "      security:\n" +
                "      - petstore_auth:\n" +
                "        - write:pets\n" +
                "        - read:pets\n" +
                "  /store/inventory:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - store\n" +
                "      summary: Returns pet inventories by status\n" +
                "      description: Returns a map of status codes to quantities\n" +
                "      operationId: getInventory\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      parameters: []\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: successful operation\n" +
                "          schema:\n" +
                "            type: object\n" +
                "            additionalProperties:\n" +
                "              type: integer\n" +
                "              format: int32\n" +
                "      security:\n" +
                "      - api_key: []\n" +
                "  /store/order:\n" +
                "    post:\n" +
                "      tags:\n" +
                "      - store\n" +
                "      summary: Place an order for a pet\n" +
                "      operationId: placeOrder\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters:\n" +
                "      - in: body\n" +
                "        name: body\n" +
                "        description: order placed for purchasing the pet\n" +
                "        required: true\n" +
                "        schema:\n" +
                "          $ref: '#/definitions/Order'\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: successful operation\n" +
                "          schema:\n" +
                "            $ref: '#/definitions/Order'\n" +
                "        400:\n" +
                "          description: Invalid Order\n" +
                "  /store/order/{orderId}:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - store\n" +
                "      summary: Find purchase order by ID\n" +
                "      description: For valid response try integer IDs with value >= 1 and <= 10.\\\n" +
                "        \\ Other values will generated exceptions\n" +
                "      operationId: getOrderById\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters:\n" +
                "      - name: orderId\n" +
                "        in: path\n" +
                "        description: ID of pet that needs to be fetched\n" +
                "        required: true\n" +
                "        type: integer\n" +
                "        maximum: 10.0\n" +
                "        minimum: 1.0\n" +
                "        format: int64\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: successful operation\n" +
                "          schema:\n" +
                "            $ref: '#/definitions/Order'\n" +
                "        400:\n" +
                "          description: Invalid ID supplied\n" +
                "        404:\n" +
                "          description: Order not found\n" +
                "    delete:\n" +
                "      tags:\n" +
                "      - store\n" +
                "      summary: Delete purchase order by ID\n" +
                "      description: For valid response try integer IDs with positive integer value.\\\n" +
                "        \\ Negative or non-integer values will generate API errors\n" +
                "      operationId: deleteOrder\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters:\n" +
                "      - name: orderId\n" +
                "        in: path\n" +
                "        description: ID of the order that needs to be deleted\n" +
                "        required: true\n" +
                "        type: integer\n" +
                "        minimum: 1.0\n" +
                "        format: int64\n" +
                "      responses:\n" +
                "        400:\n" +
                "          description: Invalid ID supplied\n" +
                "        404:\n" +
                "          description: Order not found\n" +
                "  /user:\n" +
                "    post:\n" +
                "      tags:\n" +
                "      - user\n" +
                "      summary: Create user\n" +
                "      description: This can only be done by the logged in user.\n" +
                "      operationId: createUser\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters:\n" +
                "      - in: body\n" +
                "        name: body\n" +
                "        description: Created user object\n" +
                "        required: true\n" +
                "        schema:\n" +
                "          $ref: '#/definitions/User'\n" +
                "      responses:\n" +
                "        default:\n" +
                "          description: successful operation\n" +
                "  /user/createWithArray:\n" +
                "    post:\n" +
                "      tags:\n" +
                "      - user\n" +
                "      summary: Creates list of users with given input array\n" +
                "      operationId: createUsersWithArrayInput\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters:\n" +
                "      - in: body\n" +
                "        name: body\n" +
                "        description: List of user object\n" +
                "        required: true\n" +
                "        schema:\n" +
                "          type: array\n" +
                "          items:\n" +
                "            $ref: '#/definitions/User'\n" +
                "      responses:\n" +
                "        default:\n" +
                "          description: successful operation\n" +
                "  /user/createWithList:\n" +
                "    post:\n" +
                "      tags:\n" +
                "      - user\n" +
                "      summary: Creates list of users with given input array\n" +
                "      operationId: createUsersWithListInput\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters:\n" +
                "      - in: body\n" +
                "        name: body\n" +
                "        description: List of user object\n" +
                "        required: true\n" +
                "        schema:\n" +
                "          type: array\n" +
                "          items:\n" +
                "            $ref: '#/definitions/User'\n" +
                "      responses:\n" +
                "        default:\n" +
                "          description: successful operation\n" +
                "  /user/login:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - user\n" +
                "      summary: Logs user into the system\n" +
                "      operationId: loginUser\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters:\n" +
                "      - name: username\n" +
                "        in: query\n" +
                "        description: The user name for login\n" +
                "        required: true\n" +
                "        type: string\n" +
                "      - name: password\n" +
                "        in: query\n" +
                "        description: The password for login in clear text\n" +
                "        required: true\n" +
                "        type: string\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: successful operation\n" +
                "          schema:\n" +
                "            type: string\n" +
                "          headers:\n" +
                "            X-Rate-Limit:\n" +
                "              type: integer\n" +
                "              format: int32\n" +
                "              description: calls per hour allowed by the user\n" +
                "            X-Expires-After:\n" +
                "              type: string\n" +
                "              format: date-time\n" +
                "              description: date in UTC when token expires\n" +
                "        400:\n" +
                "          description: Invalid username/password supplied\n" +
                "  /user/logout:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - user\n" +
                "      summary: Logs out current logged in user session\n" +
                "      operationId: logoutUser\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters: []\n" +
                "      responses:\n" +
                "        default:\n" +
                "          description: successful operation\n" +
                "  /user/{username}:\n" +
                "    get:\n" +
                "      tags:\n" +
                "      - user\n" +
                "      summary: Get user by user name\n" +
                "      operationId: getUserByName\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters:\n" +
                "      - name: username\n" +
                "        in: path\n" +
                "        description: The name that needs to be fetched. Use user1 for testing. \n" +
                "        required: true\n" +
                "        type: string\n" +
                "      responses:\n" +
                "        200:\n" +
                "          description: successful operation\n" +
                "          schema:\n" +
                "            $ref: '#/definitions/User'\n" +
                "        400:\n" +
                "          description: Invalid username supplied\n" +
                "        404:\n" +
                "          description: User not found\n" +
                "    put:\n" +
                "      tags:\n" +
                "      - user\n" +
                "      summary: Updated user\n" +
                "      description: This can only be done by the logged in user.\n" +
                "      operationId: updateUser\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters:\n" +
                "      - name: username\n" +
                "        in: path\n" +
                "        description: name that need to be updated\n" +
                "        required: true\n" +
                "        type: string\n" +
                "      - in: body\n" +
                "        name: body\n" +
                "        description: Updated user object\n" +
                "        required: true\n" +
                "        schema:\n" +
                "          $ref: '#/definitions/User'\n" +
                "      responses:\n" +
                "        400:\n" +
                "          description: Invalid user supplied\n" +
                "        404:\n" +
                "          description: User not found\n" +
                "    delete:\n" +
                "      tags:\n" +
                "      - user\n" +
                "      summary: Delete user\n" +
                "      description: This can only be done by the logged in user.\n" +
                "      operationId: deleteUser\n" +
                "      produces:\n" +
                "      - application/json\n" +
                "      - application/xml\n" +
                "      parameters:\n" +
                "      - name: username\n" +
                "        in: path\n" +
                "        description: The name that needs to be deleted\n" +
                "        required: true\n" +
                "        type: string\n" +
                "      responses:\n" +
                "        400:\n" +
                "          description: Invalid username supplied\n" +
                "        404:\n" +
                "          description: User not found\n" +
                "securityDefinitions:\n" +
                "  petstore_auth:\n" +
                "    type: oauth2\n" +
                "    authorizationUrl: http://petstore.swagger.io/oauth/dialog\n" +
                "    flow: implicit\n" +
                "    scopes:\n" +
                "      write:pets: modify pets in your account\n" +
                "      read:pets: read your pets\n" +
                "  api_key:\n" +
                "    type: apiKey\n" +
                "    name: api_key\n" +
                "    in: header\n" +
                "definitions:\n" +
                "  Order:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      id:\n" +
                "        type: integer\n" +
                "        format: int64\n" +
                "      petId:\n" +
                "        type: integer\n" +
                "        format: int64\n" +
                "      quantity:\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "      shipDate:\n" +
                "        type: string\n" +
                "        format: date-time\n" +
                "      status:\n" +
                "        type: string\n" +
                "        description: Order Status\n" +
                "        enum:\n" +
                "        - placed\n" +
                "        - approved\n" +
                "        - delivered\n" +
                "      complete:\n" +
                "        type: boolean\n" +
                "        default: false\n" +
                "    xml:\n" +
                "      name: Order\n" +
                "  Category:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      id:\n" +
                "        type: integer\n" +
                "        format: int64\n" +
                "      name:\n" +
                "        type: string\n" +
                "    xml:\n" +
                "      name: Category\n" +
                "  User:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      id:\n" +
                "        type: integer\n" +
                "        format: int64\n" +
                "      username:\n" +
                "        type: string\n" +
                "      firstName:\n" +
                "        type: string\n" +
                "      lastName:\n" +
                "        type: string\n" +
                "      email:\n" +
                "        type: string\n" +
                "      password:\n" +
                "        type: string\n" +
                "      phone:\n" +
                "        type: string\n" +
                "      userStatus:\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "        description: User Status\n" +
                "    xml:\n" +
                "      name: User\n" +
                "  Tag:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      id:\n" +
                "        type: integer\n" +
                "        format: int64\n" +
                "      name:\n" +
                "        type: string\n" +
                "    xml:\n" +
                "      name: Tag\n" +
                "  Pet:\n" +
                "    type: object\n" +
                "    required:\n" +
                "    - name\n" +
                "    - photoUrls\n" +
                "    properties:\n" +
                "      id:\n" +
                "        type: integer\n" +
                "        format: int64\n" +
                "      category:\n" +
                "        $ref: '#/definitions/Category'\n" +
                "      name:\n" +
                "        type: string\n" +
                "        example: doggie\n" +
                "      photoUrls:\n" +
                "        type: array\n" +
                "        xml:\n" +
                "          name: photoUrl\n" +
                "          wrapped: true\n" +
                "        items:\n" +
                "          type: string\n" +
                "      tags:\n" +
                "        type: array\n" +
                "        xml:\n" +
                "          name: tag\n" +
                "          wrapped: true\n" +
                "        items:\n" +
                "          $ref: '#/definitions/Tag'\n" +
                "      status:\n" +
                "        type: string\n" +
                "        description: pet status in the store\n" +
                "        enum:\n" +
                "        - available\n" +
                "        - pending\n" +
                "        - sold\n" +
                "    xml:\n" +
                "      name: Pet\n" +
                "  ApiResponse:\n" +
                "    type: object\n" +
                "    properties:\n" +
                "      code:\n" +
                "        type: integer\n" +
                "        format: int32\n" +
                "      type:\n" +
                "        type: string\n" +
                "      message:\n" +
                "        type: string\n" +
                "externalDocs:\n" +
                "  description: Find out more about Swagger\n" +
                "  url: http://swagger.io";

//        PayloadConfiguration configuration = new PayloadConfiguration();
//        configuration.configurePayload("yolom@seekjobs4u.com","yolo4958",  swaggerDef2, "1.3.5", "/newContext");


        Wso2ApiGatewayPlugin plugin = new Wso2ApiGatewayPlugin();
        plugin.afterApiVersionSaved(petStore);

//        Wso2Api api = new Wso2Api();
//        System.out.println(api.getClientIdAndSecret("akilaaroshana@hotmail.com","ms9714", "Akilahotmail123"));

//        Wso2ApiGatewayPlugin plugin = new Wso2ApiGatewayPlugin();
//        plugin.afterApiVersionSaved(swaggerYaml);

    }
}
