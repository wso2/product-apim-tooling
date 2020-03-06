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

import logging
import logging.handlers
import yaml
import os

# variables
log_level = logging.INFO
abs_path = os.path.abspath(os.path.dirname(__file__))
log_file = abs_path + '/../../../../../logs/attack-tool.log'


def set_log_level():
    global log_level

    with open(abs_path + '/../../../../../config/user-settings.yaml', 'r') as f:
        user_settings = yaml.load(f, Loader=yaml.FullLoader)
    user_log_level = user_settings['log_level']

    if user_log_level.strip().lower() == "debug":
        log_level = logging.DEBUG
    elif user_log_level.strip().lower() == "info":
        log_level = logging.INFO
    elif user_log_level.strip().lower() == "warning":
        log_level = logging.WARNING
    elif user_log_level.strip().lower() == "error":
        log_level = logging.ERROR
    else:
        log_level = logging.CRITICAL


def set_logger(logger_name):
    global log_level, log_file

    logger = logging.getLogger(logger_name)
    logger.setLevel(log_level)
    formatter = logging.Formatter('[%(levelname)s] %(asctime)s: %(filename)s: %(message)s')
    formatter1 = logging.Formatter('[%(levelname)s] %(asctime)s: %(message)s')

    file_handler = logging.handlers.RotatingFileHandler(
        log_file, maxBytes=10485760, backupCount=10, encoding='utf-8'
    )
    file_handler.setFormatter(formatter)
    logger.addHandler(file_handler)

    stream_handler = logging.StreamHandler()
    stream_handler.setFormatter(formatter1)
    logger.addHandler(stream_handler)

    return logger


set_log_level()
