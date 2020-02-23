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

from datetime import datetime
import os
import requests
import json

abs_path = os.path.abspath(os.path.dirname(__file__))

def log(filename, tag, write_string):
    """
    This function will write the given log output to the given log file
    :param filename: filename to write the output
    :param tag: Log tag
    :param write_string: Message to be written
    :return: None
    """
    
    with open('{}/../../../../../logs/{}'.format(abs_path, filename), 'a+') as log_file:
        log_file.write("[{}] ".format(tag) + str(datetime.now()) + ": " + write_string + "\n")

