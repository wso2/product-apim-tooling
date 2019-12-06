#!/bin/sh

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
  echo "DOS attack finished. See 'logs/attack-shell.log' for details"
}

# function to start a ddos attack
func_DDOS() {
  JMPATH=$(cat "$(pwd)"/../config/user-settings.yaml | shyaml get-value path_variables.jmeter)
  PROCESSES=$(cat "$(pwd)"/../config/attack-tool.yaml | shyaml get-value general_config.number_of_processes)
  $JMPATH/jmeter -n -t "$(pwd)"'/../lib/attack-tool/src/jmeter/DDOS_Attack.jmx' -JThreads=$PROCESSES -l "$(pwd)"/../logs/jmeter-results-attack_tool.log -j "$(pwd)"/../logs/jmeter-attack_tool.log
  echo "DDOS attack finished. See 'logs/attack-shell.log' for details"
}

# function to start an abnormal token usage attack
func_abnormal_token_usage() {
  if command -v python3 &>/dev/null; then
    nohup python3 "$(pwd)"/../lib/attack-tool/src/python/abnormal_token_usage.py >> "$(pwd)"/../logs/attack-shell.log 2>&1 &
    echo $! > "$(pwd)"/../data/attack_tool.pid
    echo "Abnormal token usage attack started. See 'logs/attack-shell.log' for details"
  elif command -v python &>/dev/null; then
    nohup python "$(pwd)"/../lib/attack-tool/src/python/abnormal_token_usage.py >> "$(pwd)"/../logs/attack-shell.log 2>&1 &
    echo $! > "$(pwd)"/../data/attack_tool.pid
    echo "Abnormal token usage attack started. See 'logs/attack-shell.log' for details"
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
    echo "Extreme delete attack started. See 'logs/attack-shell.log' for details"
  elif command -v python &>/dev/null; then
    nohup python "$(pwd)"/../lib/attack-tool/src/python/extreme_delete.py >> "$(pwd)"/../logs/attack-shell.log 2>&1 &
    echo $! > "$(pwd)"/../data/attack_tool.pid
    echo "Extreme delete attack started. See 'logs/attack-shell.log' for details"
  else
    echo "Python 3 is required for the command!"
    exit 1
  fi
}

# function to start a stolen token attack
func_stolen_token() {
  if command -v python3 &>/dev/null; then
    nohup python3 "$(pwd)"/../lib/attack-tool/src/python/stolen_token.py >> "$(pwd)"/../logs/attack-shell.log 2>&1 &
    echo $! > "$(pwd)"/../data/attack_tool.pid
    echo "Stolen token attack started. See 'logs/attack-shell.log' for details"
  elif command -v python &>/dev/null; then
    nohup python "$(pwd)"/../lib/attack-tool/src/python/stolen_token.py >> "$(pwd)"/../logs/attack-shell.log 2>&1 &
    echo $! > "$(pwd)"/../data/attack_tool.pid
    echo "Stolen token attack started. See 'logs/attack-shell.log' for details"
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
      if [ $? -eq 0 ];
      then
          echo "Attack Tool Stopped Successfully"
      fi
    else
      echo "Attack Tool Already Stopped"
    fi
  fi
  > "$(pwd)"/../data/attack_tool.pid
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
