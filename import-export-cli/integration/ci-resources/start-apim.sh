#!/bin/bash

APIM_DOWNLOAD='https://wso2.org/jenkins/job/products/job/product-apim/lastStableBuild/org.wso2.am$wso2am/artifact/org.wso2.am/wso2am/4.0.0-SNAPSHOT/wso2am-4.0.0-SNAPSHOT.zip'
APIM_PACK=wso2am-4.0.0-SNAPSHOT

wget $APIM_DOWNLOAD

# Unzip and update deployment.toml in downloaded pack
unzip -q $APIM_PACK.zip -d offset0/distribution
unzip -q $APIM_PACK.zip -d offset1/distribution
cp offset0/deployment.toml offset0/distribution/$APIM_PACK/repository/conf/.
cp offset1/deployment.toml offset1/distribution/$APIM_PACK/repository/conf/.

# Start APIM servers
offset0/distribution/$APIM_PACK/bin/api-manager.sh start 

echo "Starting APIM Instance 1"

offset1/distribution/$APIM_PACK/bin/api-manager.sh start 

echo "Starting APIM Instance 2"

# Check if APIM servers are ready
APIM1_STATUS=404
HALF_MINUTE_COUNTER=0
while [[ $APIM1_STATUS != 200 ]]
do
    if (($HALF_MINUTE_COUNTER > 0)); then        
        sleep 30s               
        if (($HALF_MINUTE_COUNTER % 2 == 0)); then                        
            let MINUTES=($HALF_MINUTE_COUNTER / 2)
            echo "Waited $MINUTES minute(s) for APIM Instance 1 to start"
        fi
    fi  
    ((++HALF_MINUTE_COUNTER))         
    APIM1_STATUS=$(curl --write-out %{http_code} --silent --output /dev/null "http://localhost:9763/services/Version")
done

echo "APIM Instance 1 started"


APIM2_STATUS=404
HALF_MINUTE_COUNTER=0
while [[ $APIM2_STATUS  != 200 ]]
do
    if (($HALF_MINUTE_COUNTER > 0)); then
        sleep 30s        
        if (($HALF_MINUTE_COUNTER % 2 == 0)); then
            let MINUTES=($HALF_MINUTE_COUNTER / 2)
            echo "Waited $MINUTES minute(s) for APIM Instance 2 to start"
        fi
    fi
    ((++HALF_MINUTE_COUNTER))
    APIM2_STATUS=$(curl --write-out %{http_code} --silent --output /dev/null "http://localhost:9764/services/Version")
done

echo "APIM Instance 2 started"


