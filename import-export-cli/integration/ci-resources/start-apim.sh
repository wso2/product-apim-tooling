#!/bin/bash

APIM_DOWNLOAD=https://github.com/wso2/product-apim/releases/download/v4.0.0-alpha/wso2am-4.0.0-alpha.zip
APIM_PACK=wso2am-4.0.0-alpha

wget $APIM_DOWNLOAD

# Unzip and update deployment.toml in downloaded pack
unzip -q $APIM_PACK.zip -d offset0/distribution
unzip -q $APIM_PACK.zip -d offset1/distribution
cp offset0/deployment.toml offset0/distribution/$APIM_PACK/repository/conf/.
cp offset1/deployment.toml offset1/distribution/$APIM_PACK/repository/conf/.

# Start APIM servers
offset0/distribution/$APIM_PACK/bin/wso2server.sh start 

echo "Starting APIM Instance 1"

offset1/distribution/$APIM_PACK/bin/wso2server.sh start 

echo "Starting APIM Instance 2"

# Check if APIM servers are ready
APIM1_STATUS=404
while [[ $APIM1_STATUS != 200 ]]
do
    APIM1_STATUS=$(curl --write-out %{http_code} --silent --output /dev/null "http://localhost:9763/services/Version")
done

echo "APIM Instance 1 started"


APIM2_STATUS=404
while [[ $APIM2_STATUS  != 200 ]]
do
    APIM2_STATUS=$(curl --write-out %{http_code} --silent --output /dev/null "http://localhost:9764/services/Version")
done

echo "APIM Instance 2 started"


