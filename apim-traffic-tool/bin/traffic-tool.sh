#!/bin/sh

# help function to print help commands when invalid arguments are given
func_help() {
  echo "Traffic Tool Options"
  echo "1: Generate data for example scenario"
  echo "2: Create scenario in APIM"
  echo "3: Generate access tokens"
  echo "4: Generate traffic data (without invoking)"
  echo "5: Simulate traffic"
  echo "all: Setup scenario and simulate traffic"
  echo "stop: Stop traffic tool"
  echo "clean: Cleanup scenario data in API Manager"
  # echo "user_details: Generate random user details"
}

# function to generate a set of random user details
func_gen_user_details() {
  if command -v python3 &>/dev/null; then
    python3 "$(pwd)"/../lib/traffic-tool/src/python/gen_user_details.py 0
  elif command -v python &>/dev/null; then
    python "$(pwd)"/../lib/traffic-tool/src/python/gen_user_details.py 0
  else
    echo "Python 3 is required for the command!"
    exit 1
  fi
}

# function to generate random user details and distribute them according to the example scenario
func_gen_example_scenario() {
  if command -v python3 &>/dev/null; then
    python3 "$(pwd)"/../lib/traffic-tool/src/python/gen_user_details.py 1
  elif command -v python &>/dev/null; then
    python "$(pwd)"/../lib/traffic-tool/src/python/gen_user_details.py 1
  else
    echo "Python 3 is required for the command!"
    exit 1
  fi
}

# function to create APIs, applications and users according to the scenario in APIM
func_create_scenario() {
  SCENARIONAME=$(cat "$(pwd)"/../config/traffic-tool.yaml | shyaml get-value scenario_name)
  JMPATH=$(cat "$(pwd)"/../config/user-settings.yaml | shyaml get-value path_variables.jmeter)

  if [ -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/api_creation.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/api_creation_swagger.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/app_creation.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/app_api_subscription_admin.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/user_generation.csv ];
  then
    rm -f "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/api_invoke_key_secret.csv
    $JMPATH/jmeter -n -t "$(pwd)"'/../lib/traffic-tool/src/jmeter/create_api_scenario.jmx' -l "$(pwd)"/../logs/jmeter-results-traffic_tool.log -j "$(pwd)"/../logs/jmeter-traffic_tool.log
    echo "Script execution completed"
  else
    echo "Missing one or more required files in the 'scenario/$SCENARIONAME/' directory"
    exit 1
  fi
}

# function to generate invoke tokens and the user scenario pool
func_gen_tokens() {
  SCENARIONAME=$(cat "$(pwd)"/../config/traffic-tool.yaml | shyaml get-value scenario_name)
  JMPATH=$(cat "$(pwd)"/../config/user-settings.yaml | shyaml get-value path_variables.jmeter)

  if [ -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/app_creation.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/user_app_pattern.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/api_invoke_key_secret.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/invoke_scenario.yaml ];
  then
    rm -f "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/api_invoke_tokens.csv
    $JMPATH/jmeter -n -t "$(pwd)"'/../lib/traffic-tool/src/jmeter/generate_token_list.jmx' -l "$(pwd)"/../logs/jmeter-results-traffic_tool.log -j "$(pwd)"/../logs/jmeter-traffic_tool.log
    echo "Token generation completed"
  else
    echo "Missing one or more required files in the 'scenario/$SCENARIONAME/' directory"
    exit 1
  fi
  if command -v python3 &>/dev/null; then
    python3 "$(pwd)"/../lib/traffic-tool/src/python/gen_invoke_scenario.py
    echo "Script execution completed"
  elif command -v python &>/dev/null; then
    python "$(pwd)"/../lib/traffic-tool/src/python/gen_invoke_scenario.py
    echo "Script execution completed"
  else
    echo "Python 3 is required for the command!"
    exit 1
  fi
}

# function to generate traffic data without invoking APIs
func_gen_invoke_data() {
  if [ -e "$(pwd)"/../lib/traffic-tool/data/runtime_data/scenario_pool.sav ];
  then
    echo "Enter filename (without file extension):"
    read FILENAME
    chmod +x "$(pwd)"/../lib/traffic-tool/src/python/gen_invoke_data.py

    if command -v python3 &>/dev/null; then
      rm -f "$(pwd)"/../lib/traffic-tool/data/runtime_data/traffic_processes.pid
      nohup python3 "$(pwd)"/../lib/traffic-tool/src/python/gen_invoke_data.py $FILENAME >> "$(pwd)"/../logs/traffic-shell.log 2>&1 &
      echo $! > "$(pwd)"/../data/traffic_tool.pid
      echo "Data generation script started"
    elif command -v python &>/dev/null; then
      rm -f "$(pwd)"/../lib/traffic-tool/data/runtime_data/traffic_processes.pid
      nohup python "$(pwd)"/../lib/traffic-tool/src/python/gen_invoke_data.py $FILENAME >> "$(pwd)"/../logs/traffic-shell.log 2>&1 &
      echo $! > "$(pwd)"/../data/traffic_tool.pid
      echo "Data generation script started"
    else
      echo "Python 3 is required for the command!"
      exit 1
    fi
  else
    echo "Missing 'scenario_pool.sav' file"
    exit 1
  fi
}

# function to simulate a traffic on APIM
func_traffic() {
  if [ -e "$(pwd)"/../lib/traffic-tool/data/runtime_data/scenario_pool.sav ];
  then
    echo "Enter filename (without file extension): "
    read FILENAME
    echo "Enter script execution time in minutes: "
    read EXECTIME
    chmod +x "$(pwd)"/../lib/traffic-tool/src/python/invoke_API.py

    if command -v python3 &>/dev/null; then
      rm -f "$(pwd)"/../lib/traffic-tool/data/runtime_data/traffic_processes.pid
      nohup python3 "$(pwd)"/../lib/traffic-tool/src/python/invoke_API.py $FILENAME $EXECTIME >> "$(pwd)"/../logs/traffic-shell.log 2>&1 &
      echo $! > "$(pwd)"/../data/traffic_tool.pid
      echo "Traffic tool started. Wait $EXECTIME minutes to complete the script"
    elif command -v python &>/dev/null; then
      rm -f "$(pwd)"/../lib/traffic-tool/data/runtime_data/traffic_processes.pid
      nohup python "$(pwd)"/../lib/traffic-tool/src/python/invoke_API.py $FILENAME $EXECTIME >> "$(pwd)"/../logs/traffic-shell.log 2>&1 &
      echo $! > "$(pwd)"/../data/traffic_tool.pid
      echo "Traffic tool started. Wait $EXECTIME minutes to complete the script"
    else
      echo "Python 3 is required for the command!"
      exit 1
    fi
  else
    echo "Missing 'scenario_pool.sav' file"
    exit 1
  fi
}

# function to stop the traffic tool while it is running
func_stop_traffic() {
  PID=`cat "$(pwd)"/../data/traffic_tool.pid 2>/dev/null`
  if [ -z $PID ];
  then
    echo "Traffic Tool is Not Running"
  else
    kill -0 $PID 2>/dev/null
    if [ $? -eq 0 ];
    then
      kill -9 $PID

      while IFS= read -r subPID; do
        kill -9 $subPID
      done < "$(pwd)"/../lib/traffic-tool/data/runtime_data/traffic_processes.pid

      if [ $? -eq 0 ];
      then
          echo "Traffic Tool Stopped Successfully"
      fi
    else
      echo "Traffic Tool Already Stopped"
    fi
  fi
  > "$(pwd)"/../data/traffic_tool.pid
  > "$(pwd)"/../lib/traffic-tool/data/runtime_data/traffic_processes.pid
}

# function to generate random user details, generate user distribution, create example scenario, generate access tokens
# and simulate the traffic from a single command
func_all() {
  echo "Enter filename (without file extension): "
  read FILENAME
  echo "Enter script execution time in minutes: "
  read EXECTIME

  SCENARIONAME=$(cat "$(pwd)"/../config/traffic-tool.yaml | shyaml get-value scenario_name)
  JMPATH=$(cat "$(pwd)"/../config/user-settings.yaml | shyaml get-value path_variables.jmeter)

  if command -v python3 &>/dev/null; then
    python3 "$(pwd)"/../lib/traffic-tool/src/python/gen_user_details.py 1

    if [ -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/api_creation.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/api_creation_swagger.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/app_creation.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/app_api_subscription_admin.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/user_generation.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/user_app_pattern.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/invoke_scenario.yaml ];
    then
      rm -f "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/api_invoke_key_secret.csv
      $JMPATH/jmeter -n -t "$(pwd)"'/../lib/traffic-tool/src/jmeter/create_api_scenario.jmx' -l "$(pwd)"/../logs/jmeter-results-traffic_tool.log -j "$(pwd)"/../logs/jmeter-traffic_tool.log
      echo "Scenario created in APIM"

      rm -f "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/api_invoke_tokens.csv
      $JMPATH/jmeter -n -t "$(pwd)"'/../lib/traffic-tool/src/jmeter/generate_token_list.jmx' -l "$(pwd)"/../logs/jmeter-results-traffic_tool.log -j "$(pwd)"/../logs/jmeter-traffic_tool.log
      echo "Token generation completed"

      python3 "$(pwd)"/../lib/traffic-tool/src/python/gen_invoke_scenario.py

      if [ -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/api_invoke_tokens.csv ];
      then
        rm -f "$(pwd)"/../lib/traffic-tool/data/runtime_data/traffic_processes.pid
        chmod +x "$(pwd)"/../lib/traffic-tool/src/python/invoke_API.py
        nohup python3 "$(pwd)"/../lib/traffic-tool/src/python/invoke_API.py $FILENAME $EXECTIME >> "$(pwd)"/../logs/traffic-shell.log 2>&1 &
        echo $! > "$(pwd)"/../data/traffic_tool.pid
        echo "Traffic tool started. Wait $EXECTIME minutes to complete the script"
      else
        echo "Missing token file in the 'data/scenario/$SCENARIONAME/' directory"
        exit 1
      fi
    else
      echo "Missing one or more required files in the 'data/scenario/$SCENARIONAME/data/' directory"
      exit 1
    fi
  elif command -v python &>/dev/null; then
    python "$(pwd)"/../lib/traffic-tool/src/python/gen_user_details.py 1

    if [ -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/api_creation.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/api_creation_swagger.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/app_creation.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/app_api_subscription_admin.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/user_generation.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/user_app_pattern.csv -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/data/invoke_scenario.yaml ];
    then
      rm -f "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/api_invoke_key_secret.csv
      $JMPATH/jmeter -n -t "$(pwd)"'/../lib/traffic-tool/src/jmeter/create_api_scenario.jmx' -l "$(pwd)"/../logs/jmeter-results-traffic_tool.log -j "$(pwd)"/../logs/jmeter-traffic_tool.log
      echo "Scenario created in APIM"

      rm -f "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/api_invoke_tokens.csv
      $JMPATH/jmeter -n -t "$(pwd)"'/../lib/traffic-tool/src/jmeter/generate_token_list.jmx' -l "$(pwd)"/../logs/jmeter-results-traffic_tool.log -j "$(pwd)"/../logs/jmeter-traffic_tool.log
      echo "Token generation completed"

      python "$(pwd)"/../lib/traffic-tool/src/python/gen_invoke_scenario.py

      if [ -e "$(pwd)"/../lib/traffic-tool/data/scenario/$SCENARIONAME/api_invoke_tokens.csv ];
      then
        rm -f "$(pwd)"/../lib/traffic-tool/data/runtime_data/traffic_processes.pid
        chmod +x "$(pwd)"/../lib/traffic-tool/src/python/invoke_API.py
        nohup python "$(pwd)"/../lib/traffic-tool/src/python/invoke_API.py $FILENAME $EXECTIME >> "$(pwd)"/../logs/traffic-shell.log 2>&1 &
        echo $! > "$(pwd)"/../data/traffic_tool.pid
        echo "Traffic tool started. Wait $EXECTIME minutes to complete the script"
      else
        echo "Missing token file in the 'data/scenario/$SCENARIONAME/' directory"
        exit 1
      fi
    else
      echo "Missing one or more required files in the 'data/scenario/$SCENARIONAME/data/' directory"
      exit 1
    fi
  else
    echo "Python 3 is required for the command!"
    exit 1
  fi
}

# function to remove created APIs, applications and users from the API Manager
func_cleanup() {
  if [ -e "$(pwd)"/../lib/traffic-tool/data/runtime_data/api_ids.csv -a -e "$(pwd)"/../lib/traffic-tool/data/runtime_data/app_ids.csv ];
  then
    JMPATH=$(cat "$(pwd)"/../config/user-settings.yaml | shyaml get-value path_variables.jmeter)

    $JMPATH/jmeter -n -t "$(pwd)"'/../lib/traffic-tool/src/jmeter/cleanup_api_manager.jmx' -l "$(pwd)"/../logs/jmeter-results-traffic_tool.log -j "$(pwd)"/../logs/jmeter-traffic_tool.log
    > "$(pwd)"/../lib/traffic-tool/data/runtime_data/api_ids.csv
    > "$(pwd)"/../lib/traffic-tool/data/runtime_data/app_ids.csv
    echo "Script execution completed"
  else
    echo "Missing required data files"
    exit 1
  fi
}


case "$1" in
  -h)
    func_help
    exit 0
  ;;
  1)
    func_gen_example_scenario 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  2)
    func_create_scenario 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  3)
    func_gen_tokens 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  4)
    func_gen_invoke_data 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  5)
    func_traffic 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  all)
    func_all 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  stop)
    func_stop_traffic 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  clean)
    func_cleanup 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  user_details)
    func_gen_user_details 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  *)
    echo "Invalid argument!"
    func_help
    exit 1
  ;;
esac
