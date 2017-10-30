package com.swagger.plugins.wso2;

import org.json.simple.parser.ParseException;

import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;


/*****************************************************************
 * Class name : Main
 * Methods : main
 * Functionality : Executes the application
 * Visibility : Public
 * ****************************************************************/
public class Main {
    public static void main(String[] args) throws IOException, ParseException, PluginExecutionException {

        String path1 = "/home/akila/Akila/swagger-wso2-integration-local/src/test/java/com/swagger/pluigns/" +
                "wso2/SwaggerSimpleApi.yaml";
        String path2 = "/home/akila/Akila/swagger-wso2-integration-local/src/test/java/com/swagger/pluigns/" +
                "wso2/SwaggerPetStore.yaml";
        String path3 = "/home/akila/Akila/swagger-wso2-integration-local/src/test/java/com/swagger/pluigns/" +
                "wso2/SwaggerOAuth2Implicit.yaml";
        String path4 = "/home/akila/Akila/swagger-wso2-integration-local/src/test/java/com/swagger/pluigns/" +
                "wso2/SwaggerHomeIOT.yaml";
        String path5 = "/home/akila/Akila/swagger-wso2-integration-local/src/test/java/com/swagger/pluigns/" +
                "wso2/OAuth2Password.yaml";
        String path6 = "/home/akila/Akila/swagger-wso2-integration-local/src/test/java/com/swagger/pluigns/" +
                "wso2/OAuth2Application.yaml";
        String path7 = "/home/akila/Akila/swagger-wso2-integration-local/src/test/java/com/swagger/pluigns/" +
                "wso2/OAuth2AccessCode.yaml";
        String path8 = "/home/akila/Akila/swagger-wso2-integration-local/src/test/java/com/swagger/pluigns/" +
                "wso2/PhoneVerifyApi.yaml";

        WSO2ApiGatewayPlugin plugin = new WSO2ApiGatewayPlugin();
        String swaggerYaml = plugin.beforeApiVersionSaved(fileReader(path1));
        plugin.afterApiVersionSaved(swaggerYaml);
    }

    /**
     * This method not mandatory for the application, just implemented to get the content of different yaml files.
     * Therefore, findBug errors can be omitted.
     * @param path Path of the yaml file
     * @return     Returns the content of the Yaml file of the provided path.
     */
    public static String fileReader(String path) throws IOException {

        StringBuffer buffer = new StringBuffer();
        BufferedReader br = null;
        FileReader fr = null;

        try {
            fr = new FileReader(path);
            br = new BufferedReader(fr);

            String sCurrentLine;

            while ((sCurrentLine = br.readLine()) != null) {
                buffer.append(sCurrentLine);
                buffer.append("\n");
            }
        } catch (IOException e) {
            throw e;
        } finally {
            try {
                if (br != null) {
                    br.close();
                }
                if (fr != null) {
                    fr.close();
                }
            } catch (IOException ex) {
                throw ex;
            }
        }
        return buffer.toString();
    }

}
