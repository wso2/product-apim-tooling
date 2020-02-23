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
import glob
from base64 import b64encode
from utils import request_methods, util_methods


# variables
abs_path = ""
apim_version = ""
token_registration_endpoint = ""
token_endpoint = ""
publisher_api_endpoint = ""
store_application_endpoint = ""
store_subs_endpoint = ""
user_signup_endpoint = ""

gateway_protocol = ""
gateway_host = ""
gateway_servelet_port_https = ""
gateway_servelet_port_http = ""
nio_pt_transport_port = ""

production_endpoint = ""
sandbox_endpoint = ""
api_throttling_tier = ""
api_visibility = ""
app_throttling_tier = ""
token_validity_period = ""
subscription_tier = ""

swagger_definitions = {}
api_ids = {}


def loadConfig():
    """
    This function will load and set the configuration data
    :return: None
    """
    global abs_path, apim_version, token_registration_endpoint, token_endpoint, publisher_api_endpoint, store_application_endpoint, store_subs_endpoint
    global user_signup_endpoint, gateway_protocol, gateway_host, gateway_servelet_port_https, gateway_servelet_port_http, nio_pt_transport_port
    global production_endpoint, sandbox_endpoint, api_throttling_tier, api_visibility, app_throttling_tier, token_validity_period, subscription_tier

    abs_path = os.path.abspath(os.path.dirname(__file__))

    with open(abs_path + '/../../../../config/apim.yaml', 'r') as config_file:
        apim_config = yaml.load(config_file, Loader=yaml.FullLoader)

    apim_version = str(apim_config['apim_version'])
    token_registration_endpoint = str(apim_config['apim_endpoints']['token_registration_endpoint'])
    token_endpoint = str(apim_config['apim_endpoints']['token_endpoint'])
    publisher_api_endpoint = str(apim_config['apim_endpoints']['publisher_api'])
    store_application_endpoint = str(apim_config['apim_endpoints']['store_application'])
    store_subs_endpoint = str(apim_config['apim_endpoints']['store_subscription'])
    user_signup_endpoint = str(apim_config['apim_endpoints']['user_signup'])
    
    gateway_protocol = str(apim_config['management_console']['protocol'])
    gateway_host = str(apim_config['management_console']['host'])
    gateway_servelet_port_https = str(apim_config['management_console']['servlet_transport_port_https'])
    gateway_servelet_port_http = str(apim_config['management_console']['servlet_transport_port_http'])
    nio_pt_transport_port = str(apim_config['api_manager']['nio_pt_transport_port'])

    with open(abs_path + '/../../../../config/traffic-tool.yaml', 'r') as traffic_config_file:
        traffic_config = yaml.load(traffic_config_file, Loader=yaml.FullLoader)

    production_endpoint = str(traffic_config['api']['production_endpoint'])
    sandbox_endpoint = str(traffic_config['api']['sandbox_endpoint'])
    api_throttling_tier = str(traffic_config['api']['throttling_tier'])
    api_visibility = str(traffic_config['api']['visibility'])
    app_throttling_tier = str(traffic_config['application']['throttling_tier'])
    token_validity_period = str(traffic_config['application']['token_validity_period'])
    subscription_tier = str(traffic_config['subscription_tier'])


def loadSwagger():
    """
    This function will load and set the swagger definitions for APIs
    :return: None
    """

    global swagger_definitions

    # load all the json files in the swagger directory
    path = abs_path + '/../../../../data/swagger/*.json'
    files=glob.glob(path)

    # read each swagger file and convert to swagger string
    for file in files:
        filename = os.path.basename(file).split('.')[0]

        f = open(file, 'r')
        content = f.read()
        f.close()

        content = content.replace('\b', '\\b').replace('\f', '\\f').replace('\n', '\\n').replace('\r', '\\r').replace('\t', '\\t').replace('"', '\\"')
        content = content.replace('    ', '')

        swagger_definitions[filename] = content


def swaggerCheck():
    """
    This function will check if the swagger definitions exist for all APIs
    :return: True if exists. False otherwise
    """
    swagger_not_found = []

    # check swagger exists for all APIs
    swagger_names = swagger_definitions.keys()

    # read API data from api_details.yaml file
    with open(abs_path + '/../../../../config/api_details.yaml', 'r') as config_file:
        api_config = yaml.load(config_file, Loader=yaml.FullLoader)

    # check swagger files exists for all APIs
    apis = api_config['apis']
    
    for api in apis:
        api_name = api['name'].lower()

        if api_name not in swagger_names:
            swagger_not_found.append(api_name)
    
    # generate the final result
    if len(swagger_not_found) >= 1:
        res_txt = "Swagger files not found for following APIs: {}".format(swagger_not_found)
        util_methods.log("traffic-tool.log", "ERROR", res_txt)
        print("[ERROR] {}".format(res_txt))
        return False

    util_methods.log("traffic-tool.log", "INFO", "Swagger files found for all APIs. Total API definitions found: {}".format(str(len(apis)-len(swagger_not_found))))
    return True


def createAndPublishAPIs():
    """
    This function will create and publish all APIs
    :return: None
    """
    created_count = 0
    published_count = 0

    # get id and secret
    client_id, client_secret = request_methods.getIDSecret(gateway_protocol, gateway_host, gateway_servelet_port_https, token_registration_endpoint)
    
    if client_id == None or client_secret == None:
        util_methods.log("traffic-tool.log", "ERROR", "Fetching client id, client secret unsuccessful!. Aborting task...")
        print("[ERROR] Fetching client id, client secret unsuccessful!. Aborting task...")
        return
    util_methods.log("traffic-tool.log", "INFO", "Successfully fetched client id, client secret")

    concat_value = client_id + ":" + client_secret
    b64_encoded = base64.b64encode(concat_value.encode('utf-8')).decode('utf-8')

    # get access token to create APIs
    access_token_create = request_methods.getAccessToken(gateway_protocol, gateway_host, nio_pt_transport_port, token_endpoint, b64_encoded, 'apim:api_create apim:api_view')[0]

    if access_token_create == None:
        util_methods.log("traffic-tool.log", "ERROR", "Getting API creation access token failed!. Aborting task...")
        print("[ERROR] Getting API creation access token failed!. Aborting task...")
        return
    util_methods.log("traffic-tool.log", "INFO", "Successfully received API creation access token")

    # get access token to publish
    access_token_publish = request_methods.getAccessToken(gateway_protocol, gateway_host, nio_pt_transport_port, token_endpoint, b64_encoded, 'apim:api_publish')[0]
    if access_token_publish == None:
        util_methods.log("traffic-tool.log", "ERROR", "Getting API publishing access token failed!. Aborting task...")
        print("[ERROR] Getting API publishing access token failed!. Aborting task...")
        return
    util_methods.log("traffic-tool.log", "INFO", "Successfully received API publishing access token")

    # read all API data from api_details.yaml file
    with open(abs_path + '/../../../../config/api_details.yaml', 'r') as config_file:
        api_config = yaml.load(config_file, Loader=yaml.FullLoader)

    apis = api_config['apis']
    util_methods.log("traffic-tool.log", "INFO", "API data read from api_details.yaml")

    # iterate the procedure for each API in the config file
    for api in apis:
        api_name = api['name']
        
        # create new API
        api_id = request_methods.createAPI(gateway_protocol, gateway_host, gateway_servelet_port_https, publisher_api_endpoint, access_token_create, api_name, api['description'], api['context'], api['version'], swagger_definitions.get(api_name.lower()), api['tags'], api_throttling_tier, api_visibility, production_endpoint, sandbox_endpoint)

        if not api_id:
            util_methods.log("traffic-tool.log", "ERROR", "API creation Failed!. API name: {}. Retrying...".format(api_name))
            api_id = request_methods.createAPI(gateway_protocol, gateway_host, gateway_servelet_port_https, publisher_api_endpoint, access_token_create, api_name, api['description'], api['context'], api['version'], swagger_definitions.get(api_name.lower()), api['tags'], api_throttling_tier, api_visibility, production_endpoint, sandbox_endpoint)
            if not api_id:
                util_methods.log("traffic-tool.log", "ERROR", "API creation Failed!. API name: {}".format(api_name))
            else:
                util_methods.log("traffic-tool.log", "INFO", "API created Successfully. API name: {}, API ID: {}".format(api_name, api_id))
                api_ids[api_name] = api_id
                created_count += 1
        else:
            util_methods.log("traffic-tool.log", "INFO", "API created Successfully. API name: {}, API ID: {}".format(api_name, api_id))
            api_ids[api_name] = api_id
            created_count += 1

        # publish the API
        ret_val = request_methods.publishAPI(gateway_protocol, gateway_host, gateway_servelet_port_https, publisher_api_endpoint, access_token_publish, api_id)

        if not ret_val:
            util_methods.log("traffic-tool.log", "ERROR", "API publishing Failed!. API name: {}. Retrying...".format(api_name))
            ret_val = request_methods.publishAPI(gateway_protocol, gateway_host, gateway_servelet_port_https, publisher_api_endpoint, access_token_publish, api_id)
            if not ret_val:
                util_methods.log("traffic-tool.log", "ERROR", "API publishing Failed!. API name: {}.".format(api_name))
            else:
                util_methods.log("traffic-tool.log", "INFO", "API published Successfully. API name: {}, API ID: {}".format(api_name, api_id))
                published_count += 1
        else:
            util_methods.log("traffic-tool.log", "INFO", "API published Successfully. API name: {}, API ID: {}".format(api_name, api_id))
            published_count += 1

    out_txt = "API creation process completed. Total {} APIs created. Total {} APIs published".format(str(created_count), str(published_count))
    util_methods.log("traffic-tool.log", "INFO", out_txt)
    print("[INFO] {}".format(out_txt))


def createApplicationsAndSubscribe():
    """
    This function will create all applications and subscribe to APIs
    :return: None
    """
    created_count = 0
    app_api_sub = ""
    
    # get id and secret
    client_id, client_secret = request_methods.getIDSecret(gateway_protocol, gateway_host, gateway_servelet_port_https, token_registration_endpoint)
    
    if client_id == None or client_secret == None:
        util_methods.log("traffic-tool.log", "ERROR", "Fetching client id, client secret unsuccessful!. Aborting task...")
        print("[ERROR] Fetching client id, client secret unsuccessful!. Aborting task...")
        return
    util_methods.log("traffic-tool.log", "INFO", "Successfully fetched client id, client secret")

    concat_value = client_id + ":" + client_secret
    b64_encoded = base64.b64encode(concat_value.encode('utf-8')).decode('utf-8')

    # get subscriber access token
    access_token_subs = request_methods.getAccessToken(gateway_protocol, gateway_host, nio_pt_transport_port, token_endpoint, b64_encoded, 'apim:subscribe apim:api_view')[0]

    if access_token_subs == None:
        util_methods.log("traffic-tool.log", "ERROR", "Getting subscription access token failed!. Aborting task...")
        print("[ERROR] Getting subscription access token failed!. Aborting task...")
        return
    util_methods.log("traffic-tool.log", "INFO", "Successfully received subscribe access token")

    # read all app data from apim.yaml file
    with open(abs_path + '/../../../../config/apim.yaml', 'r') as config_file:
        apim_config = yaml.load(config_file, Loader=yaml.FullLoader)

    apps = apim_config['apps']
    util_methods.log("traffic-tool.log", "INFO", "Application data read from apim.yaml")

    # clear application key, secret file
    with open(abs_path + '/../../data/scenario/api_invoke_key_secret.csv', 'w') as file:
        file.write('')

    # iterate the procedure for each application in the config file
    for app in apps:
        app_name = app['name']
        api_subs = app['api_subscriptions'].split(',')
        app_api_sub += app_name + ": "

        # create new Application
        app_id = request_methods.createApplication(gateway_protocol, gateway_host, gateway_servelet_port_https, store_application_endpoint, access_token_subs, app_name, app['description'], app_throttling_tier)

        if not app_id:
            util_methods.log("traffic-tool.log", "ERROR", "App creation Failed!. App name: {}. Retrying...".format(app_name))
            app_id = request_methods.createApplication(gateway_protocol, gateway_host, gateway_servelet_port_https, store_application_endpoint, access_token_subs, app_name, app['description'], app_throttling_tier)
            if not app_id:
                util_methods.log("traffic-tool.log", "ERROR", "App creation Failed!. App name: {}".format(app_name))
            else:
                util_methods.log("traffic-tool.log", "INFO", "Application created Successfully. App name: {}, App ID: {}".format(app_name, app_id))
                created_count += 1
        else:
            util_methods.log("traffic-tool.log", "INFO", "Application created Successfully. App name: {}, App ID: {}".format(app_name, app_id))
            created_count += 1
        
        #subscribe for each API
        sub_count = 0
        for api in api_subs:
            ret_val = request_methods.subscribe(gateway_protocol, gateway_host, gateway_servelet_port_https, store_subs_endpoint, access_token_subs, subscription_tier, api_ids.get(api), app_id)

            if not ret_val:
                util_methods.log("traffic-tool.log", "ERROR", "API subscription Failed!. App name: {}, API name: {}. Retrying...".format(app_name, api))
                ret_val = request_methods.subscribe(gateway_protocol, gateway_host, gateway_servelet_port_https, store_subs_endpoint, access_token_subs, subscription_tier, api_ids.get(api), app_id)
                if not ret_val:
                    util_methods.log("traffic-tool.log", "ERROR", "API subscription Failed!. App name: {}, API name: {}".format(app_name, api))
                else:
                    util_methods.log("traffic-tool.log", "INFO", "API subscription successful. App name: {}, API name: {}".format(app_name, api))
                    sub_count += 1
            else:
                util_methods.log("traffic-tool.log", "INFO", "API subscription successful. App name: {}, API name: {}".format(app_name, api))
                sub_count += 1
        app_api_sub += str(sub_count) + " "

        # generate production key
        keygen_endpoint = str(store_application_endpoint) + '/generate-keys?applicationId=' + str(app_id)
        key, secret = request_methods.genProductionKey(gateway_protocol, gateway_host, gateway_servelet_port_https, keygen_endpoint, access_token_subs, token_validity_period)

        if key == None:
            util_methods.log("traffic-tool.log", "ERROR", "App key generation Failed!. App name: {}. Retrying...".format(app_name))
            key, secret = request_methods.genProductionKey(gateway_protocol, gateway_host, gateway_servelet_port_https, keygen_endpoint, access_token_subs, token_validity_period)
            if key == None:
                util_methods.log("traffic-tool.log", "ERROR", "App key generation Failed!. App name: {}".format(app_name))
                app_api_sub += "false, "
            else:
                util_methods.log("traffic-tool.log", "INFO", "App key generation successful. App name: {}".format(app_name))
                app_api_sub += "true, "
        else:
            util_methods.log("traffic-tool.log", "INFO", "App key generation successful. App name: {}".format(app_name))
            app_api_sub += "true, "

        # write key secret to api_invoke_key_secret.csv file
        concat_value = key + ":" + secret
        b64_encoded = base64.b64encode(concat_value.encode('utf-8')).decode('utf-8')

        with open(abs_path + '/../../data/scenario/api_invoke_key_secret.csv', 'a+') as file:
            file.write(app_name + ',' + b64_encoded + '\n')

    out_txt = "Application creation process completed. Total {} Apps created. Following subscriptions happen\n{}".format(str(created_count), app_api_sub)
    util_methods.log("traffic-tool.log", "INFO", out_txt)
    print("[INFO] {}".format(out_txt))


def createUsers():
    """
    This function will self-signup all the users in carbon
    :return: None
    """
    signup_count = 0
    
    with open(abs_path + '/../../data/scenario/user_details.yaml', 'r') as user_file:
        user_data = yaml.load(user_file, Loader=yaml.FullLoader)

    if '2.' in apim_version:
        for user in user_data['users']:
            all_fields_values = user['firstname'] +'|' + user['lastname'] + '|' + user['organization'] + '|' + user['country'] + '|' + user['email'] + '|' + user['no_land'] + '|' + user['no_mobile'] + '|' + user['IM'] + '|' + user['url']
            ret_val = request_methods.selfSignupStoreAPI('http', gateway_host, gateway_servelet_port_http, user_signup_endpoint, user['username'], user['password'], all_fields_values)

            if not ret_val:
                util_methods.log("traffic-tool.log", "ERROR", "User creation Failed!. username: {}. Retrying...".format(user['username']))
                ret_val = request_methods.selfSignupStoreAPI('http', gateway_host, gateway_servelet_port_http, user_signup_endpoint, user['username'], user['password'], all_fields_values)
                if not ret_val:
                    util_methods.log("traffic-tool.log", "ERROR", "User creation Failed!. username: {}".format(user['username']))
                else:
                    util_methods.log("traffic-tool.log", "INFO", "User creation Successful!. username: {}".format(user['username']))
                    signup_count += 1
            else:
                signup_count += 1

    else:
        for user in user_data['users']:
            ret_val = request_methods.selfSignupIS(gateway_protocol, gateway_host, gateway_servelet_port_https, user_signup_endpoint, user['username'], user['password'], user['firstname'], user['lastname'], user['email'], user['country'], user['organization'], user['no_land'], user['no_mobile'], user['IM'], user['url'])

            if not ret_val:
                util_methods.log("traffic-tool.log", "ERROR", "User creation Failed!. username: {}. Retrying...".format(user['username']))
                ret_val = request_methods.selfSignupIS(gateway_protocol, gateway_host, gateway_servelet_port_https, user_signup_endpoint, user['username'], user['password'], user['firstname'], user['lastname'], user['email'], user['country'], user['organization'], user['no_land'], user['no_mobile'], user['IM'], user['url'])
                if not ret_val:
                    util_methods.log("traffic-tool.log", "ERROR", "User creation Failed!. username: {}".format(user['username']))
                else:
                    util_methods.log("traffic-tool.log", "INFO", "User creation Successful!. username: {}".format(user['username']))
                    signup_count += 1
            else:
                signup_count += 1

    out_txt = "User self signup completed. Total {} users signed up".format(signup_count)
    util_methods.log("traffic-tool.log", "INFO", out_txt)
    print("[INFO] {}".format(out_txt))


if __name__ == "__main__":
    loadConfig()
    loadSwagger()
    if swaggerCheck():
        createAndPublishAPIs()
        createApplicationsAndSubscribe()
        createUsers()
