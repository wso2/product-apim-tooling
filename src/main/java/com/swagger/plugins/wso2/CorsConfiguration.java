package com.swagger.plugins.wso2;

/*****************************************************************
 *  Class name : CorsConfiguration
 * Attributes : corsConfigurationEnabled
 * Methods : isCorsConfigurationEnabled, setCorsConfigurationEnabled
 * Functionality : Represents the nested element(corsConfiguration) of the payload
 * Visibility : Public
 * ****************************************************************/
public class CorsConfiguration {

    private boolean corsConfigurationEnabled;

    public CorsConfiguration(boolean value) {
        this.corsConfigurationEnabled = value;
    }

    public boolean isCorsConfigurationEnabled() {
        return corsConfigurationEnabled;
    }

    public void setCorsConfigurationEnabled(boolean corsConfigurationEnabled) {
        this.corsConfigurationEnabled = corsConfigurationEnabled;
    }
}
