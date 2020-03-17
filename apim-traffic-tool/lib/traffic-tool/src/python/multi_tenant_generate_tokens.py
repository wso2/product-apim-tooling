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

import os
import yaml
import csv
from utils import request_methods, log


# variables
logger = log.setLogger('multi_tenant_generate_tokens')
abs_path = ""
gateway_protocol = ""
gateway_host = ""
nio_pt_transport_port = ""
token_endpoint = ""
user_app = []
app_key_secret = {}
app_tenant = {}


def loadConfig():
    """
    This function will load and set the configuration data
    :return: None
    """
    global abs_path, gateway_protocol, gateway_host, nio_pt_transport_port, token_endpoint, app_tenant

    abs_path = os.path.abspath(os.path.dirname(__file__))

    with open(abs_path + '/../../../../config/apim.yaml', 'r') as config_file:
        apim_config = yaml.load(config_file, Loader=yaml.FullLoader)

    gateway_protocol = str(apim_config['management_console']['protocol'])
    gateway_host = str(apim_config['management_console']['host'])
    nio_pt_transport_port = str(apim_config['api_manager']['nio_pt_transport_port'])
    token_endpoint = str(apim_config['apim_endpoints']['token_endpoint'])

    with open(abs_path + '/../../data/scenario/tenant_details.yaml', 'r') as f:
        tenant_apps = yaml.load(f, Loader=yaml.FullLoader)['tenant_apps']

    for tenant in tenant_apps:
        for app in tenant_apps.get(tenant):
            app_tenant[app] = tenant


def loadUserAppPattern():
    """
    This function will read and load all user-application combinations
    :return: None
    """
    global user_app

    with open(abs_path + '/../../data/scenario/user_details.yaml', 'r') as user_file:
        user_data = yaml.load(user_file, Loader=yaml.FullLoader)

    for user in user_data['users']:
        username = user['username']
        apps = user['applications'].split(',')
        for app in apps:
            user_app.append([username, app.strip()])

    logger.info("User app pattern loaded successfully")


def loadKeySecrets():
    """
    This function will load b64 encrypted key:secret values for each application
    :return: None
    """
    global app_key_secret

    with open(abs_path + '/../../data/scenario/api_invoke_key_secret.csv', 'r') as f:
        reader = csv.reader(f)
        app_key_secret = {rows[0]:rows[1] for rows in reader}

    logger.info("Consumer key, secret loaded successfully")


def generateTokenList():
    """
    This function will generate access tokens for all user-app combinations and write to api_invoke_tokens.csv file
    :return: None
    """
    global user_app, app_key_secret, app_tenant

    # clear api_invoke_token.csv file
    with open(abs_path + '/../../data/scenario/api_invoke_tokens.csv', 'w') as f:
        f.write('username,app_name,access_token\n')

    # iterate for each user-app combination
    for uapp in user_app:
        tenant = app_tenant.get(uapp[1])
        username = uapp[0] + '@' + tenant
        
        access_token = request_methods.generateInvokeToken(
            gateway_protocol, gateway_host, nio_pt_transport_port, token_endpoint, 
            app_key_secret.get(uapp[1]), username, uapp[0], 'apim:api_view')[0]

        if access_token == None:
            logger.error("API Invoke token generation Failed!. Username: {}. Application: {}. Retrying...".format(username, uapp[1]))
            
            access_token = request_methods.generateInvokeToken(
                gateway_protocol, gateway_host, nio_pt_transport_port, token_endpoint, 
                app_key_secret.get(uapp[1]), username, uapp[0], 'apim:api_view')[0]

            if access_token == None:
                logger.error("API Invoke token generation Failed!. Username: {}. Application: {}".format(username, uapp[1]))
            else:
                logger.info("API Invoke token generation successful!. Username: {}. Application: {}. Retrying...".format(username, uapp[1]))

        with open(abs_path + '/../../data/scenario/api_invoke_tokens.csv', 'a+') as f:
            f.write(uapp[0] + ',' + uapp[1] + ',' + access_token + '\n')
        
    logger.info("API Invoke token generation completed")


# execute the functions
if __name__ == "__main__":
    loadConfig()
    loadUserAppPattern()
    loadKeySecrets()
    generateTokenList()
