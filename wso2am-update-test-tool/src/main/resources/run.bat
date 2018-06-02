@echo off
echo "welcome to WSO2 API-M update test tool"
cd .\apache-jmeter-4.0\
.\bin\jmeter -n -t ..\wso2am-update_test.jmx -l ..\log.jtl
