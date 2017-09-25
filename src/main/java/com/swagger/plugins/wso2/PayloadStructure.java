package com.swagger.plugins.wso2;


/*****************************************************************
 *  Class name : PayloadStructure
 * Functionality : Contains getters and setters of the elements of the payload
 * Visibility : Public
 * ****************************************************************/
public class PayloadStructure {

    private String name;
    private String context;
    private String version;
    private String provider;
    private String apiDefinition;
    private boolean isDefaultVersion;
    private String[] transport;
    private String[] tiers;
    private String visibility;
    private String endpointConfig;
    private CorsConfiguration corsConfiguration;

    public String getApiDefinition() {
        return apiDefinition;
    }

    public void setApiDefinition(String apiDefinition) {
        this.apiDefinition = apiDefinition;
    }

    public String getEndpointConfig() {
        return endpointConfig;
    }

    public void setEndpointConfig(String endpointConfig) {
        this.endpointConfig = endpointConfig;
    }

    public CorsConfiguration getCorsConfiguration() {
        return corsConfiguration;
    }

    public void setCorsConfiguration(CorsConfiguration corsConfiguration) {
        this.corsConfiguration = corsConfiguration;
    }

    public String[] getTiers() {
        return tiers;
    }

    public void setTiers(String[] tiers) {
        this.tiers = tiers;
    }

    public boolean getIsDefaultVersion() {
        return isDefaultVersion;
    }

    public void setIsDefaultVersion(boolean defaultVersion) {
        isDefaultVersion = defaultVersion;
    }

    public String[] getTransport() {
        return transport;
    }

    public void setTransport(String[] transport) {
        this.transport = transport;
    }

    public String getVisibility() {
        return visibility;
    }

    public void setVisibility(String visibility) {
        this.visibility = visibility;
    }

    public String getProvider() {
        return provider;
    }

    public void setProvider(String provider) {
        this.provider = provider;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getContext() {
        return context;
    }

    public void setContext(String context) {
        this.context = context;
    }

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
    }
}
