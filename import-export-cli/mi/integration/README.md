# APICTL Integration tests for Micro Integrator

## Pre-requisites for running integration tests
1. In order to run the integration tests you need to have the archive(zip) file of the Micro Integrator you need to test.

2. To setup the Micro Integrator and the MySQL database needed for the tests you can use the `setup_MI.sh` script. You need to pass the zip file path of the Micro Integrator to it as shown in the below example.

   `sh setup_MI.sh /home/micro-integrator/distribution/target/wso2mi-4.0.0-SNAPSHOT.zip`

   This script will do the following,

      - Update the `deployment.toml` with the database configuration.
      - Copy the database connector jars to lib folder.
      - Build a docker image of the updated Micro Integrator and run it.
      - Build a docker image of the preconfigured MySQL database and run it.

   Once all of the above steps are completed, it will show the IP address of the Micro Integrator container. You need to use it in the `mi/integration/config.yaml` file as explained in the next step.

   To remove remove the containers, images and the network created to run apictl integration tests for MI, the `cleanup_MI.sh` script can be used.

3. The `mi/integration/config.yaml` contains the global configurations of the integration tests for MI related commands. It currently governs the following,

- *MI Environment configs* 

  MI environment instances mentioned previously. Default configuration is local instance with port offset **0** as shown below. These must be changed to suite your MI instances if they vary.
  
  **NOTE** You need to replace localhost with the MI container IP address you obtained from the above step.

 ```
   name: testing
   host: localhost
   offset: 0
```

3. Build the apictl source to create the archive distribution of your choice.


## Executing command

| All commands must be run from *mi/integration* directory

### Flags ###

- Required:

   `-archive` :  apictl archive file that is to be tested

- Optional:

   `-run` : Run specific test fucntion only

   `-v` : Print verbose test output, useful for debugging

   `-logtransport` : Print http transport level request/responses



### Command ###

- Basic command

```
go test -archive <apictl archive name>

example: go test -archive apictl-4.0.2-linux-x64.tar.gz

```

- Run a specific test function only

```
go test -archive <apictl archive name> -test.run <Test function name or partial name regex>

example: go test -archive apictl-4.0.2-linux-x64.tar.gz -test.run TestGetConnectors
```

- Print verbose output

```
go test  -archive <apictl archive name> -test.v

example: go test -archive apictl-4.0.2-linux-x64.tar.gz -test.v
```

- Print http transport request/responses

```
go test -archive <apictl archive name> -logtransport

example: go test -archive apictl-4.0.2-linux-x64.tar.gz -logtransport
```
