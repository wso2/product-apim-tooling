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
import json
import copy
from utils import request_methods, log

# variables
logger = log.setLogger('create_api_scenario')

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
swagger_template = None
template_components = None
api_ids = {}

tenant_name = "super"
admin_username = "admin"
admin_password = "admin"
admin_b64 = "YWRtaW46YWRtaW4="

def loadConfig():
    """
    This function will load and set the configuration data
    :return: None
    """
    global abs_path, apim_version, token_registration_endpoint, token_endpoint, publisher_api_endpoint, store_application_endpoint, store_subs_endpoint
    global user_signup_endpoint, gateway_protocol, gateway_host, gateway_servelet_port_https, gateway_servelet_port_http, nio_pt_transport_port
    global production_endpoint, sandbox_endpoint, api_throttling_tier, api_visibility, app_throttling_tier, token_validity_period, subscription_tier
    global tenant_name, admin_username, admin_password, admin_b64

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

    tenant_name = apim_config['main_tenant']['tenant_name']
    admin_username = apim_config['main_tenant']['admin_username']
    admin_password = apim_config['main_tenant']['admin_password']
    admin_b64 = apim_config['main_tenant']['admin_b64']

    with open(abs_path + '/../../../../config/traffic-tool.yaml', 'r') as traffic_config_file:
        traffic_config = yaml.load(traffic_config_file, Loader=yaml.FullLoader)

    production_endpoint = str(traffic_config['api']['production_endpoint'])
    sandbox_endpoint = str(traffic_config['api']['sandbox_endpoint'])
    api_throttling_tier = str(traffic_config['api']['throttling_tier'])
    api_visibility = str(traffic_config['api']['visibility'])
    app_throttling_tier = str(traffic_config['application']['throttling_tier'])
    token_validity_period = str(traffic_config['application']['token_validity_period'])
    subscription_tier = str(traffic_config['subscription_tier'])


def genSwagger(filename, api_object):
    """
    This function will generate the swagger definition for a given api object
    :param filename: name of the swagger file to be created
    :param api_object: dictionary containing api details
    :return: True if successful. False otherwise
    """
    global swagger_template, template_components
    
    try:
        template = copy.deepcopy(swagger_template)

        template['info']['title'] = api_object['name']
        template['info']['version'] = str(api_object['version'])
        template['info']['description'] = api_object['description']

        for resource in api_object['resources']:
            resource_template = template_components.get(resource['method'].upper())
            
            path = ""
            if resource['path'][0] != "/":
                path = '/' + resource['path']
            else:
                path = resource['path']
            
            if resource['method'].upper() == "DELETE":
                path += r"/{itemId}"
            elif resource['method'].upper() == "PUT":
                path += r"/{itemId}"
            
            template['paths'][path] = resource_template
            
        with open(abs_path + '/../../../../data/swagger/{}.json'.format(filename), 'w') as f:
            json.dump(template, f)

        return True
    
    except Exception as err:
        logger.exception("Swagger creation failed for API: {}. Error: {}".format(filename, err))
        return False


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
    global swagger_template, template_components
    swagger_generated = []
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
            if swagger_template == None:
                with open(abs_path + '/../../data/tool_data/swagger_template.json', 'r') as f:
                    swagger_template = json.load(f)
                with open(abs_path + '/../../data/tool_data/template_components.json', 'r') as f:
                    template_components = json.load(f)

            if genSwagger(api_name, api):
                swagger_generated.append(api_name)
            else:
                swagger_not_found.append(api_name)

    # generate the final result
    if len(swagger_not_found) >= 1:
        res_txt = "Swagger files not found for following APIs: {}".format(swagger_not_found)
        logger.error(res_txt)
        return False
    elif len(swagger_generated) >= 1:
        res_txt = "Swagger files not found and generated for following APIs. {}".format(swagger_generated)
        logger.info(res_txt)

    logger.info("Swagger check successful. Total API definitions found: {}".format(str(len(apis)-len(swagger_not_found))))
    return True


def createAndPublishAPIs():
    """
    This function will create and publish all APIs
    :return: None
    """
    global api_ids
    created_count = 0
    published_count = 0

    # get id and secret
    client_id, client_secret = request_methods.getIDSecret(gateway_protocol, gateway_host, gateway_servelet_port_https, token_registration_endpoint, admin_username, admin_b64)

    if client_id == None or client_secret == None:
        logger.error("Fetching client id, client secret unsuccessful!. Aborting task...")
        return
    logger.info("Successfully fetched client id, client secret")

    concat_value = client_id + ":" + client_secret
    b64_encoded = base64.b64encode(concat_value.encode('utf-8')).decode('utf-8')

    # get access token to create APIs
    access_token_create = request_methods.getAccessToken(gateway_protocol, gateway_host, nio_pt_transport_port, token_endpoint, b64_encoded, 'apim:api_create apim:api_view', admin_username, admin_password)[0]

    if access_token_create == None:
        logger.error("Getting API creation access token failed!. Aborting task...")
        return
    logger.info("Successfully received API creation access token")

    # get access token to publish
    access_token_publish = request_methods.getAccessToken(gateway_protocol, gateway_host, nio_pt_transport_port, token_endpoint, b64_encoded, 'apim:api_publish', admin_username, admin_password)[0]
    if access_token_publish == None:
        logger.error("Getting API publishing access token failed!. Aborting task...")
        return
    logger.info("Successfully received API publishing access token")

    # read all API data from api_details.yaml file
    with open(abs_path + '/../../../../config/api_details.yaml', 'r') as config_file:
        api_config = yaml.load(config_file, Loader=yaml.FullLoader)

    apis = api_config['apis']
    logger.info("API data read from api_details.yaml")

    # clear previous stored api id's
    with open(abs_path + '/../../data/runtime_data/api_ids.csv', 'w') as f:
        f.write('')

    # iterate the procedure for each API in the config file
    for api in apis:
        api_name = api['name']

        # create new API
        api_id = request_methods.createAPI(gateway_protocol, gateway_host, gateway_servelet_port_https, publisher_api_endpoint, access_token_create, api_name, api['description'], api['context'], api['version'], swagger_definitions.get(api_name.lower()), api['tags'], api_throttling_tier, api_visibility, production_endpoint, sandbox_endpoint, admin_username)

        if not api_id:
            logger.error("API creation Failed!. API name: {}. Retrying...".format(api_name))
            api_id = request_methods.createAPI(gateway_protocol, gateway_host, gateway_servelet_port_https, publisher_api_endpoint, access_token_create, api_name, api['description'], api['context'], api['version'], swagger_definitions.get(api_name.lower()), api['tags'], api_throttling_tier, api_visibility, production_endpoint, sandbox_endpoint, admin_username)
            if not api_id:
                logger.error("API creation Failed!. API name: {}".format(api_name))
            else:
                logger.info("API created Successfully. API name: {}, API ID: {}".format(api_name, api_id))
                api_ids[api_name] = api_id
                created_count += 1
        else:
            logger.info("API created Successfully. API name: {}, API ID: {}".format(api_name, api_id))
            api_ids[api_name] = api_id
            created_count += 1

        # publish the API
        ret_val = request_methods.publishAPI(gateway_protocol, gateway_host, gateway_servelet_port_https, publisher_api_endpoint, access_token_publish, api_id)

        if not ret_val:
            logger.error("API publishing Failed!. API name: {}. Retrying...".format(api_name))
            ret_val = request_methods.publishAPI(gateway_protocol, gateway_host, gateway_servelet_port_https, publisher_api_endpoint, access_token_publish, api_id)
            if not ret_val:
                logger.error("API publishing Failed!. API name: {}.".format(api_name))
            else:
                logger.info("API published Successfully. API name: {}, API ID: {}".format(api_name, api_id))
                published_count += 1
        else:
            logger.info("API published Successfully. API name: {}, API ID: {}".format(api_name, api_id))
            published_count += 1

    out_txt = "API creation process completed. Total {} APIs created. Total {} APIs published".format(str(created_count), str(published_count))
    logger.info(out_txt)


def createApplicationsAndSubscribe():
    """
    This function will create all applications and subscribe to APIs
    :return: None
    """
    created_count = 0
    app_api_sub = ""

    # get id and secret
    client_id, client_secret = request_methods.getIDSecret(gateway_protocol, gateway_host, gateway_servelet_port_https, token_registration_endpoint, admin_username, admin_b64)

    if client_id == None or client_secret == None:
        logger.error("Fetching client id, client secret unsuccessful!. Aborting task...")
        return
    logger.info("Successfully fetched client id, client secret")

    concat_value = client_id + ":" + client_secret
    b64_encoded = base64.b64encode(concat_value.encode('utf-8')).decode('utf-8')

    # get subscriber access token
    access_token_subs = request_methods.getAccessToken(gateway_protocol, gateway_host, nio_pt_transport_port, token_endpoint, b64_encoded, 'apim:subscribe apim:api_view', admin_username, admin_password)[0]

    if access_token_subs == None:
        logger.error("Getting subscription access token failed!. Aborting task...")
        return
    logger.info("Successfully received subscribe access token")

    # read all app data from apim.yaml file
    with open(abs_path + '/../../../../config/apim.yaml', 'r') as config_file:
        apim_config = yaml.load(config_file, Loader=yaml.FullLoader)

    apps = apim_config['apps']
    logger.info("Application data read from apim.yaml")

    # clear application key, secret file
    with open(abs_path + '/../../data/scenario/api_invoke_key_secret.csv', 'w') as file:
        file.write('')

    # clear previous stored app id's
    with open(abs_path + '/../../data/runtime_data/app_ids.csv', 'w') as f:
        f.write('')

    # iterate the procedure for each application in the config file
    for app in apps:
        app_name = app['name']
        api_subs = app['api_subscriptions'].split(',')
        app_api_sub += app_name + ": "

        # create new Application
        app_id = request_methods.createApplication(gateway_protocol, gateway_host, gateway_servelet_port_https, store_application_endpoint, access_token_subs, app_name, app['description'], app_throttling_tier)

        if not app_id:
            logger.error("App creation Failed!. App name: {}. Retrying...".format(app_name))
            app_id = request_methods.createApplication(gateway_protocol, gateway_host, gateway_servelet_port_https, store_application_endpoint, access_token_subs, app_name, app['description'], app_throttling_tier)
            if not app_id:
                logger.error("App creation Failed!. App name: {}".format(app_name))
            else:
                logger.info("Application created Successfully. App name: {}, App ID: {}".format(app_name, app_id))
                created_count += 1
        else:
            logger.info("Application created Successfully. App name: {}, App ID: {}".format(app_name, app_id))
            created_count += 1

        #subscribe for each API
        sub_count = 0
        for api in api_subs:
            ret_val = request_methods.subscribe(gateway_protocol, gateway_host, gateway_servelet_port_https, store_subs_endpoint, access_token_subs, subscription_tier, api_ids.get(api), app_id)

            if not ret_val:
                logger.error("API subscription Failed!. App name: {}, API name: {}. Retrying...".format(app_name, api))
                ret_val = request_methods.subscribe(gateway_protocol, gateway_host, gateway_servelet_port_https, store_subs_endpoint, access_token_subs, subscription_tier, api_ids.get(api), app_id)
                if not ret_val:
                    logger.error("API subscription Failed!. App name: {}, API name: {}".format(app_name, api))
                else:
                    logger.info("API subscription successful. App name: {}, API name: {}".format(app_name, api))
                    sub_count += 1
            else:
                logger.info("API subscription successful. App name: {}, API name: {}".format(app_name, api))
                sub_count += 1
        app_api_sub += str(sub_count) + " "

        # generate production key
        keygen_endpoint = str(store_application_endpoint) + '/generate-keys?applicationId=' + str(app_id)
        key, secret = request_methods.genProductionKey(gateway_protocol, gateway_host, gateway_servelet_port_https, keygen_endpoint, access_token_subs, token_validity_period)

        if key == None:
            logger.error("App key generation Failed!. App name: {}. Retrying...".format(app_name))
            key, secret = request_methods.genProductionKey(gateway_protocol, gateway_host, gateway_servelet_port_https, keygen_endpoint, access_token_subs, token_validity_period)
            if key == None:
                logger.error("App key generation Failed!. App name: {}".format(app_name))
                app_api_sub += "(key: false), "
            else:
                logger.info("App key generation successful. App name: {}".format(app_name))
                app_api_sub += "(key: true), "
        else:
            logger.info("App key generation successful. App name: {}".format(app_name))
            app_api_sub += "(key: true), "

        # write key secret to api_invoke_key_secret.csv file
        concat_value = key + ":" + secret
        b64_encoded = base64.b64encode(concat_value.encode('utf-8')).decode('utf-8')

        with open(abs_path + '/../../data/scenario/api_invoke_key_secret.csv', 'a+') as file:
            file.write(app_name + ',' + b64_encoded + '\n')

    out_txt = "Application creation process completed. Total {} Apps created. Following subscriptions happen:\n{}".format(str(created_count), app_api_sub)
    logger.info(out_txt)


def createUsers():
    """
    This function will self-signup all the users in carbon
    :return: None
    """
    signup_count = 0

    with open(abs_path + '/../../data/scenario/user_details.yaml', 'r') as user_file:
        user_data = yaml.load(user_file, Loader=yaml.FullLoader)

    uname_suffix = ''
    if tenant_name.lower() == "super" or tenant_name.lower() == "carbon.super":
        uname_suffix = ''
    else:
        uname_suffix = '@' + tenant_name

    if '2.' in apim_version:
        for user in user_data['users']:
            username = user['username'] + uname_suffix
            all_fields_values = user['firstname'] +'|' + user['lastname'] + '|' + user['organization'] + '|' + user['country'] + '|' + user['email'] + '|' + user['no_land'] + '|' + user['no_mobile'] + '|' + user['IM'] + '|' + user['url']
            ret_val = request_methods.selfSignupStoreAPI('http', gateway_host, gateway_servelet_port_http, user_signup_endpoint, username, user['password'], all_fields_values)

            if not ret_val:
                logger.error("User creation Failed!. username: {}. Retrying...".format(user['username']))
                ret_val = request_methods.selfSignupStoreAPI('http', gateway_host, gateway_servelet_port_http, user_signup_endpoint, username, user['password'], all_fields_values)
                if not ret_val:
                    logger.error("User creation Failed!. username: {}".format(user['username']))
                else:
                    logger.info("User creation Successful!. username: {}".format(user['username']))
                    signup_count += 1
            else:
                signup_count += 1

    else:
        for user in user_data['users']:
            ret_val = request_methods.selfSignupIS(gateway_protocol, gateway_host, gateway_servelet_port_https, user_signup_endpoint, admin_b64, user['username'], user['password'], user['firstname'], user['lastname'], user['email'], user['country'], user['organization'], user['no_land'], user['no_mobile'], user['IM'], user['url'])

            if not ret_val:
                logger.error("User creation Failed!. username: {}. Retrying...".format(user['username']))
                ret_val = request_methods.selfSignupIS(gateway_protocol, gateway_host, gateway_servelet_port_https, user_signup_endpoint, admin_b64, user['username'], user['password'], user['firstname'], user['lastname'], user['email'], user['country'], user['organization'], user['no_land'], user['no_mobile'], user['IM'], user['url'])
                if not ret_val:
                    logger.error("User creation Failed!. username: {}".format(user['username']))
                else:
                    logger.info("User creation Successful!. username: {}".format(user['username']))
                    signup_count += 1
            else:
                signup_count += 1

    out_txt = "User self signup completed. Total {} users signed up".format(signup_count)
    logger.info(out_txt)


if __name__ == "__main__":
    loadConfig()
    loadSwagger()
    if swaggerCheck():
        createAndPublishAPIs()
        createApplicationsAndSubscribe()
        createUsers()
