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
import csv
from utils import request_methods, log

# variables
logger = log.setLogger('multi_tenant_cleanup_scenario')
abs_path = ""
token_registration_endpoint = ""
token_endpoint = ""
publisher_api_endpoint = ""
store_application_endpoint = ""
delete_user_soap_endpoint = ""

gateway_protocol = ""
gateway_host = ""
gateway_servelet_port_https = ""
nio_pt_transport_port = ""

user_tenant = {}
tenant_config_details = None

def loadConfig():
    """
    This function will load and set the configuration data
    :return: None
    """
    global abs_path, token_registration_endpoint, token_endpoint, publisher_api_endpoint, store_application_endpoint, delete_user_soap_endpoint
    global gateway_protocol, gateway_host, gateway_servelet_port_https, nio_pt_transport_port
    global user_tenant, tenant_config_details
    
    abs_path = os.path.abspath(os.path.dirname(__file__))

    with open(abs_path + '/../../../../config/apim.yaml', 'r') as config_file:
        apim_config = yaml.load(config_file, Loader=yaml.FullLoader)

    token_registration_endpoint = str(apim_config['apim_endpoints']['token_registration_endpoint'])
    token_endpoint = str(apim_config['apim_endpoints']['token_endpoint'])
    publisher_api_endpoint = str(apim_config['apim_endpoints']['publisher_api'])
    store_application_endpoint = str(apim_config['apim_endpoints']['store_application'])
    delete_user_soap_endpoint = str(apim_config['apim_endpoints']['delete_user'])
    
    gateway_protocol = str(apim_config['management_console']['protocol'])
    gateway_host = str(apim_config['management_console']['host'])
    gateway_servelet_port_https = str(apim_config['management_console']['servlet_transport_port_https'])
    nio_pt_transport_port = str(apim_config['api_manager']['nio_pt_transport_port'])

    tenant_config_details = apim_config['multi_tenancy']['tenant_list']

    with open(abs_path + '/../../data/scenario/tenant_details.yaml', 'r') as tenant_file:
        user_tenant = yaml.load(tenant_file, Loader=yaml.FullLoader)['user_tenants']

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


def setupTenantKeys():
    """
    This function will get and set API creation and subscribe tokens for each tenant
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

        # get API creation access token
        access_token_create = request_methods.getAccessToken(gateway_protocol, gateway_host, nio_pt_transport_port, token_endpoint, b64_encoded, 'apim:api_create apim:api_view', tenant_data['admin_username'], tenant_data['admin_password'])[0]

        if access_token_create == None:
            logger.error("Getting API creation access token failed for tenant: {}!. Aborting task...".format(tenant))
            return False
        logger.info("Successfully received API creation access token for tenant: {}".format(tenant))
        tenant_config_details[tenant]['create_token'] = access_token_create

        # get API subscribe access token
        access_token_subs = request_methods.getAccessToken(gateway_protocol, gateway_host, nio_pt_transport_port, token_endpoint, b64_encoded, 'apim:subscribe apim:api_view', tenant_data['admin_username'], tenant_data['admin_password'])[0]

        if access_token_subs == None:
            logger.error("Getting subscription access token failed for tenant: {}!. Aborting task...".format(tenant))
            return False
        logger.info("Successfully received subscribe access token for tenant: {}".format(tenant))
        tenant_config_details[tenant]['subscribe_token'] = access_token_subs

    return True


def removeApplications():
    """
    This function will remove all created applications from API Manager (apps in multiple tenants)
    :return: None
    """
    global tenant_config_details
    remove_count = 0

    # iterate for each application
    with open(abs_path + '/../../data/runtime_data/app_ids_multi_tenant.csv', 'r') as f:
        reader = csv.reader(f)

        for app_id in reader:
            tenant_conf = tenant_config_details.get(app_id[0])

            deleted = request_methods.deleteAppAPI(
                gateway_protocol, gateway_host, gateway_servelet_port_https, store_application_endpoint, 
                tenant_conf.get('subscribe_token'), app_id[1])

            if not deleted:
                logger.error("Application removing Failed!. App id: {}, Tenant: {}. Retrying...".format(app_id[1], app_id[0]))

                deleted = request_methods.deleteAppAPI(gateway_protocol, gateway_host, 
                    gateway_servelet_port_https, store_application_endpoint, 
                    tenant_conf.get('subscribe_token'), app_id[1])

                if not deleted:
                    logger.error("Application removing Failed!. App id: {}, Tenant: {}".format(app_id[1], app_id[0]))
                else:
                    logger.info("Application removed successfully!. App id: {}, Tenant: {}".format(app_id[1], app_id[0]))
                    remove_count += 1
            else:
                logger.info("Application removed successfully!. App id: {}, Tenant: {}".format(app_id[1], app_id[0]))
                remove_count += 1
            
        logger.info("Application deletion process completed. Total {} applications removed".format(str(remove_count)))


def removeAPIs():
    """
    This function will remove all created APIs from API Manager (apis in multiple tenants)
    :return: None
    """
    global tenant_config_details
    remove_count = 0

    # iterate for each API
    with open(abs_path + '/../../data/runtime_data/api_ids_multi_tenant.csv', 'r') as f:
        reader = csv.reader(f)

        for api_id in reader:
            tenant_conf = tenant_config_details.get(api_id[0])

            deleted = request_methods.deleteAppAPI(
                gateway_protocol, gateway_host, gateway_servelet_port_https, publisher_api_endpoint, 
                tenant_conf.get('create_token'), api_id[1])

            if not deleted:
                logger.error("API removing Failed!. API id: {}, Tenant: {}. Retrying...".format(api_id[1], api_id[0]))
                
                deleted = request_methods.deleteAppAPI(gateway_protocol, gateway_host, 
                    gateway_servelet_port_https, publisher_api_endpoint, 
                    tenant_conf.get('create_token'), api_id[1])

                if not deleted:
                    logger.error("API removing Failed!. API id: {}, Tenant: {}".format(api_id[1], api_id[0]))
                else:
                    logger.info("API removed successfully!. API id: {}, Tenant: {}".format(api_id[1], api_id[0]))
                    remove_count += 1
            else:
                logger.info("API removed successfully!. API id: {}, Tenant: {}".format(api_id[1], api_id[0]))
                remove_count += 1
            
        logger.info("API deletion process completed. Total {} APIs removed".format(str(remove_count)))


def removeUsers():
    """
    This function will remove all created user accounts from carbon (users in multiple tenants)
    :return: None
    """
    global user_tenant
    remove_count = 0

    for user in user_tenant:
        for tenant in user_tenant.get(user):
            admin_b64 = tenant_config_details.get(tenant)['admin_b64']
            
            removed = request_methods.removeUserSOAP(
                gateway_protocol, gateway_host, gateway_servelet_port_https, delete_user_soap_endpoint, 
                user, admin_b64)

            if not removed:
                logger.error("User deletion Failed!. username: {}, Tenant: {}. Retrying...".format(user, tenant))

                removed = request_methods.removeUserSOAP(gateway_protocol, gateway_host, 
                    gateway_servelet_port_https, delete_user_soap_endpoint, user, admin_b64)

                if not removed:
                    logger.error("User deletion Failed!. username: {}, Tenant: {}".format(user, tenant))
                else:
                    logger.info("User removed successfully!. username: {}, Tenant: {}".format(user, tenant))
                    remove_count += 1
            else:
                remove_count += 1

    logger.info("User deletion process completed. Total {} user accounts removed".format(str(remove_count)))


if __name__ == "__main__":
    loadConfig()
    if setupTenantKeys():
        removeApplications()
        removeAPIs()
        removeUsers()
