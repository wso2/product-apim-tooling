#!/bin/sh

if [ $# -eq 0 ]
  then
    echo "Pass the WSO2 Micro Integrator zip archive path as an argument"
    echo "Example:"
    echo "sh setup_MI.sh /home/micro-integrator/distribution/target/wso2mi-4.2.0-SNAPSHOT.zip"
    exit 1
fi

WSO2_MI_ZIP_PATH=$1
WSO2_MI_TAG=mi-test-apictl
WSO2_MI_TEST_DATA_DIR=./testdata/mi
WSO2_MI_MYSQL_DB_TAG=mi-mysql-db
MYSQL_TEST_DATA_DIR=./testdata/mysql

WSO2_MI_ZIP_FILE=${WSO2_MI_ZIP_PATH##*/}
WSO2_MI_EXTRACTED=${WSO2_MI_ZIP_FILE%.*}

if unzip $WSO2_MI_ZIP_PATH -d $WSO2_MI_TEST_DATA_DIR > /dev/null; then
    echo "Extracting MI Archive completed"
else
    echo "Error extracting MI Archive"
    exit 1
fi

mv $WSO2_MI_TEST_DATA_DIR/$WSO2_MI_EXTRACTED $WSO2_MI_TEST_DATA_DIR/wso2mi

if docker build -t $WSO2_MI_TAG $WSO2_MI_TEST_DATA_DIR > /dev/null; then
    echo "Building WSO2 MI docker image completed"
    rm rm -rf $WSO2_MI_TEST_DATA_DIR/wso2mi
else
    echo "Error building WSO2 MI docker image"
    rm rm -rf $WSO2_MI_TEST_DATA_DIR/wso2mi
    exit 1
fi

if docker build -t $WSO2_MI_MYSQL_DB_TAG $MYSQL_TEST_DATA_DIR > /dev/null; then
    echo "Building MySQL docker image completed"
else
    echo "Error building MySQL docker image"
    exit 1
fi

if docker network create mi-test-net > /dev/null; then
    echo "Creating docker network completed"
else
    echo "Error creating docker network"
    exit 1
fi

if docker run -d --network mi-test-net --name sqlcontainer $WSO2_MI_MYSQL_DB_TAG > /dev/null; then
    echo "MySQL Ccntainer started successfully"
else
    echo "Error starting MySQL container"
    exit 1
fi

if docker run -d --network mi-test-net --name micontainer $WSO2_MI_TAG > /dev/null; then
    echo "WSO2 MI Ccntainer started successfully"
else
    echo "Error starting WSO2 MI container"
    exit 1
fi

WSO2_MI_CONTAINER_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' micontainer)

echo "\nUse ${WSO2_MI_CONTAINER_IP} as the host for MI instance in config.yaml before you run tests"
echo "Use cleanup_MI.sh to remove the containers, images and the network created to run apictl integration tests for MI"
