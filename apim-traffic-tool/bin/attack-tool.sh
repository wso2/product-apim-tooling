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
  echo "Attack Tool Options"
  echo "1: DOS attack"
  echo "2: DDOS attack"
  echo "3: Abnormal token usage attack"
  echo "4: Extreme delete attack"
  echo "5: Stolen token attack"
  echo "stop: Stop running attack"
}

# function to start a dos attack
func_DOS() {
  JMPATH=$(cat "$(pwd)"/../config/user-settings.yaml | shyaml get-value path_variables.jmeter)
  PROCESSES=$(cat "$(pwd)"/../config/attack-tool.yaml | shyaml get-value general_config.number_of_processes)
  $JMPATH/jmeter -n -t "$(pwd)"'/../lib/attack-tool/src/jmeter/DOS_Attack.jmx' -JThreads=$PROCESSES -l "$(pwd)"/../logs/jmeter-results-attack_tool.log -j "$(pwd)"/../logs/jmeter-attack_tool.log
  echo "DOS attack finished. See 'logs/jmeter-attack_tool.log' for details"
}

# function to start a ddos attack
func_DDOS() {
  JMPATH=$(cat "$(pwd)"/../config/user-settings.yaml | shyaml get-value path_variables.jmeter)
  PROCESSES=$(cat "$(pwd)"/../config/attack-tool.yaml | shyaml get-value general_config.number_of_processes)
  $JMPATH/jmeter -n -t "$(pwd)"'/../lib/attack-tool/src/jmeter/DDOS_Attack.jmx' -JThreads=$PROCESSES -l "$(pwd)"/../logs/jmeter-results-attack_tool.log -j "$(pwd)"/../logs/jmeter-attack_tool.log
  echo "DDOS attack finished. See 'logs/jmeter-attack_tool.log.log' for details"
}

# function to start an abnormal token usage attack
func_abnormal_token_usage() {
  if command -v python3 &>/dev/null; then
    rm -f "$(pwd)"/../lib/attack-tool/data/runtime_data/attack_processes.pid
    nohup python3 "$(pwd)"/../lib/attack-tool/src/python/abnormal_token_usage.py >> "$(pwd)"/../logs/attack-shell.log 2>&1 &
    echo $! > "$(pwd)"/../data/attack_tool.pid
    echo "Abnormal token usage attack started. See 'logs/attack-tool.log' for details"
  elif command -v python &>/dev/null; then
    rm -f "$(pwd)"/../lib/attack-tool/data/runtime_data/attack_processes.pid
    nohup python "$(pwd)"/../lib/attack-tool/src/python/abnormal_token_usage.py >> "$(pwd)"/../logs/attack-shell.log 2>&1 &
    echo $! > "$(pwd)"/../data/attack_tool.pid
    echo "Abnormal token usage attack started. See 'logs/attack-tool.log' for details"
  else
    echo "Python 3 is required for the command!"
    exit 1
  fi
}

# function to start an extreme delete attack (data deletion attack)
func_extreme_delete() {
  if command -v python3 &>/dev/null; then
    nohup python3 "$(pwd)"/../lib/attack-tool/src/python/extreme_delete.py >> "$(pwd)"/../logs/attack-shell.log 2>&1 &
    echo $! > "$(pwd)"/../data/attack_tool.pid
    echo "Extreme delete attack started. See 'logs/attack-tool.log' for details"
  elif command -v python &>/dev/null; then
    nohup python "$(pwd)"/../lib/attack-tool/src/python/extreme_delete.py >> "$(pwd)"/../logs/attack-shell.log 2>&1 &
    echo $! > "$(pwd)"/../data/attack_tool.pid
    echo "Extreme delete attack started. See 'logs/attack-tool.log' for details"
  else
    echo "Python 3 is required for the command!"
    exit 1
  fi
}

# function to start a stolen token attack
func_stolen_token() {
  if command -v python3 &>/dev/null; then
    rm -f "$(pwd)"/../lib/attack-tool/data/runtime_data/attack_processes.pid
    nohup python3 "$(pwd)"/../lib/attack-tool/src/python/stolen_token.py >> "$(pwd)"/../logs/attack-shell.log 2>&1 &
    echo $! > "$(pwd)"/../data/attack_tool.pid
    echo "Stolen token attack started. See 'logs/attack-tool.log' for details"
  elif command -v python &>/dev/null; then
    rm -f "$(pwd)"/../lib/attack-tool/data/runtime_data/attack_processes.pid
    nohup python "$(pwd)"/../lib/attack-tool/src/python/stolen_token.py >> "$(pwd)"/../logs/attack-shell.log 2>&1 &
    echo $! > "$(pwd)"/../data/attack_tool.pid
    echo "Stolen token attack started. See 'logs/attack-tool.log' for details"
  else
    echo "Python 3 is required for the command!"
    exit 1
  fi
}

# function to stop the attack tool while it is running
func_stop_attack() {
  PID=`cat "$(pwd)"/../data/attack_tool.pid 2>/dev/null`
  if [ -z $PID ];
  then
    echo "Attack Tool is Not Running"
  else
    kill -0 $PID 2>/dev/null
    if [ $? -eq 0 ];
    then
      kill -9 $PID

      if [ -e "$(pwd)"/../lib/attack-tool/data/runtime_data/attack_processes.pid ];
      then
        while IFS= read -r subPID; do
          kill -9 $subPID
        done < "$(pwd)"/../lib/attack-tool/data/runtime_data/attack_processes.pid
      fi
      
      if [ $? -eq 0 ];
      then
          echo "Attack Tool Stopped Successfully"
      fi
    else
      echo "Attack Tool Already Stopped"
    fi
  fi
  > "$(pwd)"/../data/attack_tool.pid
  rm -f "$(pwd)"/../lib/attack-tool/data/runtime_data/attack_processes.pid
}


case "$1" in
  -h)
    func_help
    exit 0
  ;;
  1)
    func_DOS 2>&1 | tee -a "$(pwd)"/../logs/attack-shell.log
    exit 0
  ;;
  2)
    func_DDOS 2>&1 | tee -a "$(pwd)"/../logs/attack-shell.log
    exit 0
  ;;
  3)
    func_abnormal_token_usage 2>&1 | tee -a "$(pwd)"/../logs/attack-shell.log
    exit 0
  ;;
  4)
    func_extreme_delete 2>&1 | tee -a "$(pwd)"/../logs/attack-shell.log
    exit 0
  ;;
  5)
    func_stolen_token 2>&1 | tee -a "$(pwd)"/../logs/attack-shell.log
    exit 0
  ;;
  stop)
    func_stop_attack 2>&1 | tee -a "$(pwd)"/../logs/attack-shell.log
    exit 0
  ;;
  *)
    echo "Invalid argument!"
    func_help
    exit 1
  ;;
esac
