#!/bin/sh

# Copyright (c) 2020, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
#
# WSO2 Inc. licenses this file to you under the Apache License,
# Version 2.0 (the "License"); you may not use this file except
# in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied. See the License for the
# specific language governing permissions and limitations
# under the License.

# help function to print help commands when invalid arguments are given
func_help() {
  echo "Traffic Tool Options"
  echo "setup: Setup scenario"
  echo "start: Simulate traffic"
  echo "stop : Stop traffic tool"
  echo "gen  : Generate traffic data (without invoking)"
  echo "clean: Cleanup scenario data in API Manager"
  echo "adv  : View advanced options"
}

# help function to print advanced commands
func_advance_help() {
  echo "Advanced Options"
  echo "1: Generate random user details"
  echo "2: Generate data for example scenario"
  echo "3: Generate invoke scenario"
  echo "4: Create scenario in APIM"
  echo "5: Generate access tokens"
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

# function to generate the scenario distribution
func_gen_scenario_distribution() {
  if command -v python3 &>/dev/null; then
    python3 "$(pwd)"/../lib/traffic-tool/src/python/gen_scenario_yaml.py
  elif command -v python &>/dev/null; then
    python "$(pwd)"/../lib/traffic-tool/src/python/gen_scenario_yaml.py
  else
    echo "Python 3 is required for the command!"
    exit 1
  fi
}

# function to create APIs, applications and users according to the scenario in APIM
func_create_scenario() {
  if [ -e "$(pwd)"/../lib/traffic-tool/data/scenario/user_details.yaml -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/invoke_scenario.yaml ];
  then
    rm -f "$(pwd)"/../lib/traffic-tool/data/scenario/api_invoke_key_secret.csv

    if command -v python3 &>/dev/null; then
      python3 "$(pwd)"/../lib/traffic-tool/src/python/create_api_scenario.py
    elif command -v python &>/dev/null; then
      python "$(pwd)"/../lib/traffic-tool/src/python/create_api_scenario.py
    else
      echo "Python 3 is required for the command!"
      exit 1
    fi
  else
    echo "Missing one or more required files in the 'scenario/' directory"
    exit 1
  fi
}

# function to generate invoke tokens and the user scenario pool
func_gen_tokens() {
  if [ -e "$(pwd)"/../lib/traffic-tool/data/scenario/user_details.yaml -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/invoke_scenario.yaml -a -e "$(pwd)"/../lib/traffic-tool/data/scenario/api_invoke_key_secret.csv ];
  then
    rm -f "$(pwd)"/../lib/traffic-tool/data/scenario/api_invoke_tokens.csv

    if command -v python3 &>/dev/null; then
      python3 "$(pwd)"/../lib/traffic-tool/src/python/generate_tokens.py
      python3 "$(pwd)"/../lib/traffic-tool/src/python/gen_invoke_scenario.py
    elif command -v python &>/dev/null; then
      python "$(pwd)"/../lib/traffic-tool/src/python/generate_tokens.py
      python "$(pwd)"/../lib/traffic-tool/src/python/gen_invoke_scenario.py
    else
      echo "Python 3 is required for the command!"
      exit 1
    fi
  else
    echo "Missing one or more required files in the 'scenario/' directory"
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

# function to remove created APIs, applications and users from the API Manager
func_cleanup() {
  if [ -e "$(pwd)"/../lib/traffic-tool/data/runtime_data/api_ids.csv -a -e "$(pwd)"/../lib/traffic-tool/data/runtime_data/app_ids.csv ];
  then
    if command -v python3 &>/dev/null; then
      python3 "$(pwd)"/../lib/traffic-tool/src/python/cleanup_scenario.py
    elif command -v python &>/dev/null; then
      python "$(pwd)"/../lib/traffic-tool/src/python/cleanup_scenario.py
    else
      echo "Python 3 is required for the command!"
      exit 1
    fi
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
  setup)
    echo "Enter '0' to setup for custom scenario. Any other key to generate a random scenario"
    read KEY
    if [ $KEY == "0" ]
    then
      func_gen_example_scenario 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
      func_create_scenario 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
      func_gen_tokens 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    else
      func_gen_user_details 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
      func_gen_scenario_distribution 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
      func_create_scenario 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
      func_gen_tokens 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    fi
    exit 0
  ;;
  start)
    func_traffic 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  stop)
    func_stop_traffic 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  gen)
    func_gen_invoke_data 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  clean)
    func_cleanup 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  adv)
    func_advance_help
    exit 0
  ;;
  1)
    func_gen_user_details 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  2)
    func_gen_example_scenario 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  3)
    func_gen_scenario_distribution 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  4)
    func_create_scenario 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  5)
    func_gen_tokens 2>&1 | tee -a "$(pwd)"/../logs/traffic-shell.log
    exit 0
  ;;
  *)
    echo "Invalid argument!"
    func_help
    exit 1
  ;;
esac
