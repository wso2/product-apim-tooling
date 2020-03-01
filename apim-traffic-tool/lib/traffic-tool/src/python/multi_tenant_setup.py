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
import base64
from utils import request_methods, log

# variables
logger = log.setLogger('multi_tenant_setup')
abs_path = ""
tenant_list = []
app_list = {}
app_api_subs = {}
app_tenant = {}


def loadConfig():
    """
    This function will load and set the configuration data
    :return: None
    """
    global abs_path, app_list, app_api_subs

    # load application details
    with open(abs_path + '/../../../../config/apim.yaml', 'r') as file:
        apim_config = yaml.load(file, Loader=yaml.FullLoader)
    apps = apim_config['apps']

    for app in apps:
        app_list[app['name']] = []
        app_api_subs[app['name']] = app['api_subscriptions'].split(',')


def genMultiTenantPattern():
    """
    This function will divide each app and relavant APIs among tenants, and write to tenant_details.yaml file
    :return: None
    """
    global tenant_list, app_tenant
    tenant_app = {}
    tenant_api = {}

    # divide apps among tenants
    i = 0
    for app in app_list:
        i += 1
        tenant = tenant_list[i%len(tenant_list)]
        apis = app_api_subs.get(app)

        if tenant['name'] not in tenant_app:
            tenant_app[tenant['name']] = []
        tenant_app.get(tenant['name']).append(app)

        if tenant['name'] not in tenant_api:
            tenant_api[tenant['name']] = set()
        for api in apis:
            tenant_api.get(tenant['name']).add(api)

        app_tenant[app] = tenant['name']
    
    for tenant in tenant_api:
        tenant_api[tenant] = list(tenant_api.get(tenant))
    
    # generate tenant_details.yaml file
    with open(abs_path + '/../../data/scenario/tenant_details.yaml', 'w') as f:
        yaml.dump({'tenant_apps': tenant_app}, f, sort_keys=False)
        yaml.dump({'tenant_apis': tenant_api}, f, sort_keys=False)


def genUserTenant():
    """
    This function will generate tenant list for each user and, write to tenant_details.yaml file
    :return: None
    """
    global app_tenant
    user_tenant = {}

    with open(abs_path + '/../../data/scenario/user_details.yaml', 'r') as user_file:
        user_data = yaml.load(user_file, Loader=yaml.FullLoader)

    for user_record in user_data['users']:
        username = user_record['username']
        apps = user_record['applications'].split(',')
        tenant_list = set()

        for app in apps:
            tenant_list.add(app_tenant.get(app.strip()))
        
        user_tenant[username] = list(tenant_list)

    # write to tenant_details.yaml file
    with open(abs_path + '/../../data/scenario/tenant_details.yaml', 'a+') as f:
        yaml.dump({'user_tenants': user_tenant}, f, sort_keys=False)


if __name__ == "__main__":
    abs_path = os.path.abspath(os.path.dirname(__file__))
    loadConfig()

    with open(abs_path + '/../../../../config/apim.yaml', 'r') as config_file:
        apim_config = yaml.load(config_file, Loader=yaml.FullLoader)
    tenant_list = apim_config['multi_tenancy']['tenant_list']

    if str(apim_config['multi_tenancy']['enabled']).lower() == 'true':
        genMultiTenantPattern()
        genUserTenant()
    else:
        logger.info("Multi tenancy is not enabled!")
