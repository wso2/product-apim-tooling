echo "welcome to WSO2 API-M update test tool"
if [ -z "$1" ]
then
	sh apache-jmeter-4.0/bin/jmeter -n -t wso2am-update_test.jmx -l log.jtl
elif [ "$1" = "um" ]
then
	echo "checking requirements to run in user defined mode..."
	count1=`ls resources/apis/*.json 2>/dev/null | wc -l`
	count2=`ls resources/applications/*.json 2>/dev/null | wc -l`
	if [ $count1 = 0 ] && [ $count2 = 0 ]
	then
		echo "you haven't specified any input files under directories \"resources/apis\" and \"resources/applications\"
Please specify it before to start"
		exit
	else
		echo "you have provided" $count1 "API(s) and" $count2 "application(s) input files."
		sh apache-jmeter-4.0/bin/jmeter -n -t resume_wso2am-update_test.jmx -l log.jtl
	fi
else
	echo invalid arugument specified. please try again.
fi
