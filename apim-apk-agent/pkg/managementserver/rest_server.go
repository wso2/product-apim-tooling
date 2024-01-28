package managementserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/config"
)

func init() {
}

// StartInternalServer starts the internal server
func StartInternalServer(port uint) {
	r := gin.Default()

	r.GET("/applications", func(c *gin.Context) {
		applicationList := GetAllApplications()
		c.JSON(http.StatusOK, ApplicationList{List: applicationList})
	})
	r.GET("/subscriptions", func(c *gin.Context) {
		subscriptionList := GetAllSubscriptions()
		c.JSON(http.StatusOK, SubscriptionList{List: subscriptionList})
	})
	r.GET("/applicationmappings", func(c *gin.Context) {
		applicationMappingList := GetAllApplicationMappings()
		c.JSON(http.StatusOK, ApplicationMappingList{List: applicationMappingList})
	})
	r.GET("/applicationkeymappings", func(c *gin.Context) {
		applicationKeyMappingList := GetAllApplicationKeyMappings()
		c.JSON(http.StatusOK, ApplicationKeyMappingList{List: applicationKeyMappingList})
	})
	gin.SetMode(gin.ReleaseMode)
	publicKeyLocation, privateKeyLocation, _ := config.GetKeyLocations()
	r.RunTLS(fmt.Sprintf(":%d", port), publicKeyLocation, privateKeyLocation)
}
