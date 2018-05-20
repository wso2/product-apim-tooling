## WSO2 API-M Update Test Tool

* This tool able to verify the customer data with newly updated WUM distribution of WSO2 API Manager (WSO2 API-M). It makes sure that the customer specific data will not cause any problem when deploying and executing on the newly updated WUM distributions.

* It extracts some selected subset of data from currently running instance of API-M and push it to the newly WUM updated version and check it.

* This tool built using Apache JMeter 4.0.

Prerequisites:
==================================

    1. Install Oracle Java SE Development Kit (JDK) version 1.8 and set the JAVA_HOME environment variable.

    2. Last stable deployment details of the API Manager (vanilla or WUM updated)

    3. Newly downloaded WUM updated distribution

Installation & Running
==================================

1. Extract the wso2am-update-test-tool.zip and go to that (wso2am-update-test-tool/) directory.

2. Configure the "config.json" file provided under "conf/" folder.
	Note: Following configurations need to be done before executing this tool. (It's already set up to defaults values).

 	* Please provide user-name and password of your currently running distribution and newly updated WUM distribution.
		"username":  "john",
		"password":  "john123",
	* Please provide requested URLs.
		1. Endpoint to obtain the consumer key/secret key pair by calling the dynamic client registration .
			"clientRegisterUrl": "https://localhost:9443/client-registration/v0.12/register",\
		2. Endpoint to get an access token with the preferred grant type (with specified scope)
			"accessTokenUrl":    "https://localhost:8243/token",
		3. Endpoint to get details or create new APIs from API publisher
			"APIsPublisherUrl":  "https://localhost:9443/api/am/publisher/v0.12/apis",
		4. Endpoint to get APIs details from API Store
			"APIsStoreUrl":      "https://localhost:9443/api/am/store/v0.12/apis",
		5. Endpoint to get details or create applications from API Store
			"applicationsStoreUrl":  "https://localhost:9443/api/am/store/v0.12/applications",
		6. Endpoint to add a new subscription providing the id of the API and the application on API Store.
			"subscriptionsStoreUrl": "https://localhost:9443/api/am/store/v0.12/subscriptions",
	* Please provide number of APIs and Applications you want to consider
		"numberOfAPIsToConsider":  "10",
		"numberOfApplicationsToConsider":  "10"

3. Run the run.sh or run.bat script based on your operating system.
	on Linux-
		sh run.sh or ./run.sh
	on Windows-
		run.bat
	* Please find generated inputs files to recreate the APIs and Applications on new distributions under "resources" directory.
	* Also find log.jtl file to view the summary results.

4. If you already have specified input details to create APIs/Applications on new server (without grabbing APIs/Applications details from latest stable deployment)
		  1. Please keep those files under corresponding folders(apis/applications) under resources directory.
  		2. Make sure you specified newly downloaded WUM updated distribution details under new instance in conf/config.json file.
   Now you can run the run.sh script with specified command line argument. (um - user defined mode )
		sh run.sh um or ./run.sh um
