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
logger = log.setLogger('multi_tenant_create_scenario')

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
tenant_config_details = None
tenant_details = {}


def loadConfig():
    """
    This function will load and set the configuration data
    :return: None
    """
    global abs_path, apim_version, token_registration_endpoint, token_endpoint, publisher_api_endpoint, store_application_endpoint, store_subs_endpoint
    global user_signup_endpoint, gateway_protocol, gateway_host, gateway_servelet_port_https, gateway_servelet_port_http, nio_pt_transport_port
    global production_endpoint, sandbox_endpoint, api_throttling_tier, api_visibility, app_throttling_tier, token_validity_period, subscription_tier
    global tenant_config_details, tenant_details

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

    tenant_config_details = apim_config['multi_tenancy']['tenant_list']

    with open(abs_path + '/../../../../config/traffic-tool.yaml', 'r') as traffic_config_file:
        traffic_config = yaml.load(traffic_config_file, Loader=yaml.FullLoader)

    production_endpoint = str(traffic_config['api']['production_endpoint'])
    sandbox_endpoint = str(traffic_config['api']['sandbox_endpoint'])
    api_throttling_tier = str(traffic_config['api']['throttling_tier'])
    api_visibility = str(traffic_config['api']['visibility'])
    app_throttling_tier = str(traffic_config['application']['throttling_tier'])
    token_validity_period = str(traffic_config['application']['token_validity_period'])
    subscription_tier = str(traffic_config['subscription_tier'])

    with open(abs_path + '/../../data/scenario/tenant_details.yaml', 'r') as tenant_file:
        tenant_details = yaml.load(tenant_file, Loader=yaml.FullLoader)

    # set tenant admin details for each tenant
    tenant_config_temp = {}
    for tenant in tenant_config_details:
        concat_value = tenant['admin_username'].strip() + ':' + tenant['admin_password'].strip()
        auth_b64 = base64.b64encode(concat_value.encode('utf-8')).decode('utf-8')
        
        tenant_config_temp[tenant['name']] = {
            "admin_username": tenant['admin_username'].strip(),
            "admin_password": tenant['admin_password'].strip(),
            "admin_b64": auth_b64
        }

    tenant_config_details = tenant_config_temp


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
    global swagger_template, template_components, apis
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


def setupTenantKeys():
    """
    This function will get and set API creation, publish and subscribe tokens for each tenant
    :return: True if successful. False otherwise
    """
    global tenant_config_details

    # set id, secret, tokens for each tenant
    for tenant in tenant_config_details:
        tenant_data = tenant_config_details.get(tenant)

        # get client id and client secret
        client_id, client_secret = request_methods.getIDSecret(gateway_protocol, gateway_host, gateway_servelet_port_https, token_registration_endpoint, tenant_data['admin_username'], tenant_data['admin_b64'])
        
        if client_id == None or client_secret == None:
            logger.error("Fetching client id, client secret unsuccessful for tenant: {}. Aborting task...".format(tenant))
            return False
        logger.info("Successfully fetched client id, client secret for tenant: {}".format(tenant))

        concat_value = client_id + ":" + client_secret
        b64_encoded = base64.b64encode(concat_value.encode('utf-8')).decode('utf-8')
        tenant_config_details[tenant]['b64_keySecret'] = b64_encoded

        # get access token to create APIs
        access_token_create = request_methods.getAccessToken(gateway_protocol, gateway_host, nio_pt_transport_port, token_endpoint, b64_encoded, 'apim:api_create apim:api_view', tenant_data['admin_username'], tenant_data['admin_password'])[0]

        if access_token_create == None:
            logger.error("Getting API creation access token failed for tenant: {}!. Aborting task...".format(tenant))
            return False
        logger.info("Successfully received API creation access token for tenant: {}".format(tenant))
        tenant_config_details[tenant]['create_token'] = access_token_create

        # get access token to publish
        access_token_publish = request_methods.getAccessToken(gateway_protocol, gateway_host, nio_pt_transport_port, token_endpoint, b64_encoded, 'apim:api_publish', tenant_data['admin_username'], tenant_data['admin_password'])[0]
        if access_token_publish == None:
            logger.error("Getting API publishing access token failed for tenant: {}!. Aborting task...".format(tenant))
            return False
        logger.info("Successfully received API publishing access token for tenant: {}".format(tenant))
        tenant_config_details[tenant]['publish_token'] = access_token_publish

        # get access token to subscribe
        access_token_subs = request_methods.getAccessToken(gateway_protocol, gateway_host, nio_pt_transport_port, token_endpoint, b64_encoded, 'apim:subscribe apim:api_view', tenant_data['admin_username'], tenant_data['admin_password'])[0]

        if access_token_subs == None:
            logger.error("Getting subscription access token failed for tenant: {}!. Aborting task...".format(tenant))
            return False
        logger.info("Successfully received subscribe access token for tenant: {}".format(tenant))
        tenant_config_details[tenant]['subscribe_token'] = access_token_subs

    return True


def createAndPublishAPIs():
    """
    This function will create and publish relevant APIs in each tenant
    :return: None
    """
    global tenant_config_details, tenant_details, api_ids
    created_count = 0
    published_count = 0

    # read all API data from api_details.yaml file
    with open(abs_path + '/../../../../config/api_details.yaml', 'r') as config_file:
        api_config = yaml.load(config_file, Loader=yaml.FullLoader)

    apis = api_config['apis']
    api_data = {}

    for api in apis:
        api_data[api['name']] = {
            'description': api['description'],
            'tags': api['tags'],
            'context': api['context'],
            'version': api['version'],
            'resources': api['resources']
        }

    logger.info("API data read from api_details.yaml")

    # clear previous stored api id's
    with open(abs_path + '/../../data/runtime_data/api_ids_multi_tenant.csv', 'w') as f:
        f.write('')

    # iterate over each tenant
    for tenant in tenant_details['tenant_apis']:
        api_list = tenant_details['tenant_apis'][tenant]
        tenant_conf = tenant_config_details.get(tenant)
        
        # iterate the procedure for each API
        for api in api_list:
            # create new API
            api_id = request_methods.createAPI(
                gateway_protocol, gateway_host, gateway_servelet_port_https, publisher_api_endpoint, tenant_conf.get('create_token'), 
                api, api_data.get(api)['description'], api_data.get(api)['context'], api_data.get(api)['version'], 
                swagger_definitions.get(api.lower()), api_data.get(api)['tags'], api_throttling_tier, api_visibility, 
                production_endpoint, sandbox_endpoint, tenant_conf.get('admin_username'))

            if not api_id:
                logger.error("API creation Failed!. API name: {}, Tenant: {}. Retrying...".format(api, tenant))
                api_id = request_methods.createAPI(
                    gateway_protocol, gateway_host, gateway_servelet_port_https, publisher_api_endpoint, tenant_conf.get('create_token'), 
                    api, api_data.get(api)['description'], api_data.get(api)['context'], api_data.get(api)['version'], 
                    swagger_definitions.get(api.lower()), api_data.get(api)['tags'], api_throttling_tier, api_visibility, 
                    production_endpoint, sandbox_endpoint, tenant_conf.get('admin_username'))

                if not api_id:
                    logger.error("API creation Failed!. API name: {}, Tenant: {}".format(api, tenant))
                else:
                    logger.info("API created Successfully. API name: {}, API ID: {}, Tenant: {}".format(api, api_id, tenant))
                    if tenant not in api_ids.keys():
                        api_ids[tenant] = {api: api_id}
                    else:
                        api_ids[tenant][api] = api_id

                    created_count += 1

                    with open(abs_path + '/../../data/runtime_data/api_ids_multi_tenant.csv', 'a+') as f:
                        f.write(tenant + ',' + api_id + '\n')
            
            else:
                logger.info("API created Successfully. API name: {}, API ID: {}, Tenant: {}".format(api, api_id, tenant))
                if tenant not in api_ids.keys():
                    api_ids[tenant] = {api: api_id}
                else:
                    api_ids[tenant][api] = api_id
                
                created_count += 1

                with open(abs_path + '/../../data/runtime_data/api_ids_multi_tenant.csv', 'a+') as f:
                    f.write(tenant + ',' + api_id + '\n')

            # publish the API
            ret_val = request_methods.publishAPI(gateway_protocol, 
                gateway_host, gateway_servelet_port_https, publisher_api_endpoint, tenant_conf.get('publish_token'), api_id)

            if not ret_val:
                logger.error("API publishing Failed!. API name: {}, Tenant: {}. Retrying...".format(api, tenant))
                ret_val = request_methods.publishAPI(gateway_protocol, gateway_host, 
                    gateway_servelet_port_https, publisher_api_endpoint, tenant_conf.get('publish_token'), api_id)
                if not ret_val:
                    logger.error("API publishing Failed!. API name: {}, Tenant: {}.".format(api, tenant))
                else:
                    logger.info("API published Successfully. API name: {}, API ID: {}, Tenant: {}".format(api, api_id, tenant))
                    published_count += 1
            else:
                logger.info("API published Successfully. API name: {}, API ID: {}, Tenant: {}".format(api, api_id, tenant))
                published_count += 1

    out_txt = "API creation process completed. Total {} APIs created. Total {} APIs published".format(str(created_count), str(published_count))
    logger.info(out_txt)


def createApplicationsAndSubscribe():
    """
    This function will create relevant applications in each tenant and subscribe to APIs
    :return: None
    """
    global api_ids
    created_count = 0
    app_api_sub = ""

    # read all app data from apim.yaml file
    with open(abs_path + '/../../../../config/apim.yaml', 'r') as config_file:
        apim_config = yaml.load(config_file, Loader=yaml.FullLoader)

    apps = apim_config['apps']
    app_data = {}

    for app in apps:
        app_data[app['name']] = {
            'description': app['description'],
            'api_subscriptions': app['api_subscriptions']
        }

    logger.info("Application data read from apim.yaml")

    # clear application key, secret file
    with open(abs_path + '/../../data/scenario/api_invoke_key_secret.csv', 'w') as file:
        file.write('')

    # clear previous stored app id's
    with open(abs_path + '/../../data/runtime_data/app_ids_multi_tenant.csv', 'w') as f:
        f.write('')

    # iterate over each tenant
    for tenant in tenant_details['tenant_apps']:
        app_list = tenant_details['tenant_apps'][tenant]
        tenant_conf = tenant_config_details.get(tenant)
        
        # iterate the procedure for each application
        for app in app_list:
            api_subs = app_data.get(app)['api_subscriptions'].split(',')
            app_api_sub += app + ": "

            # create new Application
            app_id = request_methods.createApplication(
                gateway_protocol, gateway_host, gateway_servelet_port_https, store_application_endpoint, tenant_conf.get('subscribe_token'), 
                app, app_data.get(app)['description'], app_throttling_tier)

            if not app_id:
                logger.error("App creation Failed!. App name: {}, Tenant: {}. Retrying...".format(app, tenant))
                app_id = request_methods.createApplication(gateway_protocol, 
                    gateway_host, gateway_servelet_port_https, store_application_endpoint, tenant_conf.get('subscribe_token'), 
                    app, app_data.get(app)['description'], app_throttling_tier)
                
                if not app_id:
                    logger.error("App creation Failed!. App name: {}, Tenant: {}".format(app, tenant))
                else:
                    logger.info("Application created Successfully. App name: {}, App ID: {}, Tenant: {}".format(app, app_id, tenant))
                    created_count += 1

                    with open(abs_path + '/../../data/runtime_data/app_ids_multi_tenant.csv', 'a+') as f:
                        f.write(tenant + ',' + app_id + '\n')
            
            else:
                logger.info("Application created Successfully. App name: {}, App ID: {}, Tenant: {}".format(app, app_id, tenant))
                created_count += 1

                with open(abs_path + '/../../data/runtime_data/app_ids_multi_tenant.csv', 'a+') as f:
                    f.write(tenant + ',' + app_id + '\n')

            #subscribe for each API
            sub_count = 0
            for api in api_subs:
                ret_val = request_methods.subscribe(
                    gateway_protocol, gateway_host, gateway_servelet_port_https, store_subs_endpoint, tenant_conf.get('subscribe_token'), 
                    subscription_tier, api_ids.get(tenant)[api], app_id)

                if not ret_val:
                    logger.error("API subscription Failed!. App name: {}, API name: {}, Tenant: {}. Retrying...".format(app, api, tenant))
                    ret_val = request_methods.subscribe(
                    gateway_protocol, gateway_host, gateway_servelet_port_https, store_subs_endpoint, tenant_conf.get('subscribe_token'), 
                    subscription_tier, api_ids.get(tenant)[api], app_id)

                    if not ret_val:
                        logger.error("API subscription Failed!. App name: {}, API name: {}, Tenant: {}".format(app, api, tenant))
                    else:
                        logger.info("API subscription successful. App name: {}, API name: {}, Tenant: {}".format(app, api, tenant))
                        sub_count += 1
                else:
                    logger.info("API subscription successful. App name: {}, API name: {}, Tenant: {}".format(app, api, tenant))
                    sub_count += 1
            app_api_sub += str(sub_count) + " "

            # generate production key
            keygen_endpoint = str(store_application_endpoint) + '/generate-keys?applicationId=' + str(app_id)

            key, secret = request_methods.genProductionKey(gateway_protocol, gateway_host, 
                gateway_servelet_port_https, keygen_endpoint, tenant_conf.get('subscribe_token'), token_validity_period)

            if key == None:
                logger.error("App key generation Failed!. App name: {}, Tenant: {}. Retrying...".format(app, tenant))
                key, secret = request_methods.genProductionKey(gateway_protocol, gateway_host, 
                    gateway_servelet_port_https, keygen_endpoint, tenant_conf.get('subscribe_token'), token_validity_period)

                if key == None:
                    logger.error("App key generation Failed!. App name: {}, Tenant: {}".format(app, tenant))
                    app_api_sub += "(key: false), "
                else:
                    logger.info("App key generation successful. App name: {}, Tenant: {}".format(app, tenant))
                    app_api_sub += "(key: true), "
            else:
                logger.info("App key generation successful. App name: {}, Tenant: {}".format(app, tenant))
                app_api_sub += "(key: true), "

            # write key secret to api_invoke_key_secret.csv file
            concat_value = key + ":" + secret
            b64_encoded = base64.b64encode(concat_value.encode('utf-8')).decode('utf-8')

            with open(abs_path + '/../../data/scenario/api_invoke_key_secret.csv', 'a+') as file:
                file.write(app + ',' + b64_encoded + '\n')

    out_txt = "Application creation process completed. Total {} Apps created. Following subscriptions happen:\n{}".format(str(created_count), app_api_sub)
    logger.info(out_txt)


def createUsers():
    """
    This function will self-signup all the users in carbon
    :return: None
    """
    global tenant_details, tenant_config_details
    signup_count = 0

    with open(abs_path + '/../../data/scenario/user_details.yaml', 'r') as user_file:
        user_data = yaml.load(user_file, Loader=yaml.FullLoader)

    if '2.' in apim_version:
        for user in user_data['users']:
            for tenant in tenant_details['user_tenants'].get(user['username']):
                username = user['username'] + '@' + tenant
                all_fields_values = user['firstname'] +'|' + user['lastname'] + '|' + user['organization'] + '|' + user['country'] + '|' + user['email'] + '|' + user['no_land'] + '|' + user['no_mobile'] + '|' + user['IM'] + '|' + user['url']

                ret_val = request_methods.selfSignupStoreAPI('http', gateway_host, 
                    gateway_servelet_port_http, user_signup_endpoint, username, user['password'], 
                    all_fields_values)

                if not ret_val:
                    logger.error("User creation Failed!. username: {}, Tenant: {}. Retrying...".format(user['username'], tenant))

                    ret_val = request_methods.selfSignupStoreAPI('http', gateway_host, 
                        gateway_servelet_port_http, user_signup_endpoint, username, user['password'], 
                        all_fields_values)
                    
                    if not ret_val:
                        logger.error("User creation Failed!. username: {}, Tenant: {}".format(user['username'], tenant))
                    else:
                        logger.info("User creation Successful!. username: {}, Tenant: {}".format(user['username'], tenant))
                        signup_count += 1
                else:
                    signup_count += 1

    else:
        for user in user_data['users']:
            for tenant in tenant_details['user_tenants'].get(user['username']):
                username = user['username'] + '@' + tenant
                admin_b64 = tenant_config_details[tenant]['admin_b64']

                ret_val = request_methods.selfSignupIS(gateway_protocol, gateway_host, 
                    gateway_servelet_port_https, user_signup_endpoint, admin_b64, user['username'], 
                    user['password'], user['firstname'], user['lastname'], user['email'], user['country'], 
                    user['organization'], user['no_land'], user['no_mobile'], user['IM'], user['url'])

                if not ret_val:
                    logger.error("User creation Failed!. username: {}, Tenant: {}. Retrying...".format(user['username'], tenant))

                    ret_val = request_methods.selfSignupIS(gateway_protocol, gateway_host, 
                        gateway_servelet_port_https, user_signup_endpoint, admin_b64, user['username'], 
                        user['password'], user['firstname'], user['lastname'], user['email'], user['country'], 
                        user['organization'], user['no_land'], user['no_mobile'], user['IM'], user['url'])
                    
                    if not ret_val:
                        logger.error("User creation Failed!. username: {}, Tenant: {}".format(user['username'], tenant))
                    else:
                        logger.info("User creation Successful!. username: {}, Tenant: {}".format(user['username'], tenant))
                        signup_count += 1
                else:
                    signup_count += 1

    out_txt = "User self signup completed. Total {} user accounts signed up".format(signup_count)
    logger.info(out_txt)


if __name__ == "__main__":
    loadConfig()
    loadSwagger()
    if swaggerCheck() and setupTenantKeys():
        createAndPublishAPIs()
        createApplicationsAndSubscribe()
        createUsers()
