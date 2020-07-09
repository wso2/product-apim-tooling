#compdef _apictl apictl


function _apictl {
  local -a commands

  _arguments -C \
    '(-h --help)'{-h,--help}'[help for apictl]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "add:Add an API to the kubernetes cluster"
      "add-env:Add Environment to Config file"
      "change:Change a configuration in K8s cluster resource"
      "change-status:Change Status of an API"
      "delete:Delete an API/APIProduct/Application in an environment"
      "export:Export an API Product in an environment"
      "export-api:Export API"
      "export-apis:Export APIs for migration"
      "export-app:Export App"
      "get-keys:Generate access token to invoke the API or API Product"
      "help:Help about any command"
      "import:Import an API Product to an environment"
      "import-api:Import API"
      "import-app:Import App"
      "init:Initialize a new project in given path"
      "install:Install an operator in the configured K8s cluster"
      "list:List APIs/APIProducts/Applications in an environment or List the environments"
      "login:Login to an API Manager"
      "logout:Logout to from an API Manager"
      "remove:Remove an environment"
      "set:Set configuration parameters"
      "uninstall:Uninstall an operator in the configured K8s cluster"
      "update:Update an API to the kubernetes cluster"
      "vcs:Checks status and deploys projects"
      "version:Display Version on current apictl"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  add)
    _apictl_add
    ;;
  add-env)
    _apictl_add-env
    ;;
  change)
    _apictl_change
    ;;
  change-status)
    _apictl_change-status
    ;;
  delete)
    _apictl_delete
    ;;
  export)
    _apictl_export
    ;;
  export-api)
    _apictl_export-api
    ;;
  export-apis)
    _apictl_export-apis
    ;;
  export-app)
    _apictl_export-app
    ;;
  get-keys)
    _apictl_get-keys
    ;;
  help)
    _apictl_help
    ;;
  import)
    _apictl_import
    ;;
  import-api)
    _apictl_import-api
    ;;
  import-app)
    _apictl_import-app
    ;;
  init)
    _apictl_init
    ;;
  install)
    _apictl_install
    ;;
  list)
    _apictl_list
    ;;
  login)
    _apictl_login
    ;;
  logout)
    _apictl_logout
    ;;
  remove)
    _apictl_remove
    ;;
  set)
    _apictl_set
    ;;
  uninstall)
    _apictl_uninstall
    ;;
  update)
    _apictl_update
    ;;
  vcs)
    _apictl_vcs
    ;;
  version)
    _apictl_version
    ;;
  esac
}


function _apictl_add {
  local -a commands

  _arguments -C \
    '(-h --help)'{-h,--help}'[help for add]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "api:handle APIs in kubernetes cluster "
      "help:Help about any command"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  api)
    _apictl_add_api
    ;;
  help)
    _apictl_add_help
    ;;
  esac
}

function _apictl_add_api {
  _arguments \
    '(-a --apiEndPoint)'{-a,--apiEndPoint}'[]:' \
    '(*-e *--env)'{\*-e,\*--env}'[Environment variables to be passed to deployment]:' \
    '(*-f *--from-file)'{\*-f,\*--from-file}'[Path to swagger file]:' \
    '(-h --help)'{-h,--help}'[help for api]' \
    '(-i --image)'{-i,--image}'[Image of the API. If specified, ignores the value of --override]:' \
    '(-m --mode)'{-m,--mode}'[Property to override the deploying mode. Available modes: privateJet, sidecar]:' \
    '(-n --name)'{-n,--name}'[Name of the API]:' \
    '--namespace[namespace of API]:' \
    '--override[Property to override the existing docker image with the given name and version]' \
    '--replicas[replica set]:' \
    '(-v --version)'{-v,--version}'[Property to override the API version]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_add_help {
  _arguments \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_add-env {
  _arguments \
    '--admin[Admin endpoint for the environment]:' \
    '--apim[API Manager endpoint for the environment]:' \
    '--devportal[DevPortal endpoint for the environment]:' \
    '(-e --environment)'{-e,--environment}'[Name of the environment to be added]:' \
    '(-h --help)'{-h,--help}'[help for add-env]' \
    '--publisher[Publisher endpoint for the environment]:' \
    '--registration[Registration endpoint for the environment]:' \
    '--token[Token endpoint for the environment]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}


function _apictl_change {
  local -a commands

  _arguments -C \
    '(-h --help)'{-h,--help}'[help for change]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "help:Help about any command"
      "registry:Change the registry"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  help)
    _apictl_change_help
    ;;
  registry)
    _apictl_change_registry
    ;;
  esac
}

function _apictl_change_help {
  _arguments \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_change_registry {
  _arguments \
    '(-h --help)'{-h,--help}'[help for registry]' \
    '(-c --key-file)'{-c,--key-file}'[Credentials file]:' \
    '(-p --password)'{-p,--password}'[Password of the given user]:' \
    '--password-stdin[Prompt for password of the given user in the stdin]' \
    '(-R --registry-type)'{-R,--registry-type}'[Registry type: DOCKER_HUB | AMAZON_ECR |GCR | HTTP]:' \
    '(-r --repository)'{-r,--repository}'[Repository name or URI]:' \
    '(-u --username)'{-u,--username}'[Username of the repository]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}


function _apictl_change-status {
  local -a commands

  _arguments -C \
    '(-h --help)'{-h,--help}'[help for change-status]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "api:Change Status of an API"
      "help:Help about any command"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  api)
    _apictl_change-status_api
    ;;
  help)
    _apictl_change-status_help
    ;;
  esac
}

function _apictl_change-status_api {
  _arguments \
    '(-a --action)'{-a,--action}'[Action to be taken to change the status of the API]:' \
    '(-e --environment)'{-e,--environment}'[Environment of which the API state should be changed]:' \
    '(-h --help)'{-h,--help}'[help for api]' \
    '(-n --name)'{-n,--name}'[Name of the API to be state changed]:' \
    '(-r --provider)'{-r,--provider}'[Provider of the API]:' \
    '(-v --version)'{-v,--version}'[Version of the API to be state changed]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_change-status_help {
  _arguments \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}


function _apictl_delete {
  local -a commands

  _arguments -C \
    '(-h --help)'{-h,--help}'[help for delete]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "api:Delete API"
      "api-product:Delete API Product"
      "app:Delete App"
      "help:Help about any command"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  api)
    _apictl_delete_api
    ;;
  api-product)
    _apictl_delete_api-product
    ;;
  app)
    _apictl_delete_app
    ;;
  help)
    _apictl_delete_help
    ;;
  esac
}

function _apictl_delete_api {
  _arguments \
    '(-e --environment)'{-e,--environment}'[Environment from which the API should be deleted]:' \
    '(-h --help)'{-h,--help}'[help for api]' \
    '(-n --name)'{-n,--name}'[Name of the API to be deleted]:' \
    '(-r --provider)'{-r,--provider}'[Provider of the API to be deleted]:' \
    '(-v --version)'{-v,--version}'[Version of the API to be deleted]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_delete_api-product {
  _arguments \
    '(-e --environment)'{-e,--environment}'[Environment from which the API Product should be deleted]:' \
    '(-h --help)'{-h,--help}'[help for api-product]' \
    '(-n --name)'{-n,--name}'[Name of the API Product to be deleted]:' \
    '(-r --provider)'{-r,--provider}'[Provider of the API Product to be deleted]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_delete_app {
  _arguments \
    '(-e --environment)'{-e,--environment}'[Environment from which the Application should be deleted]:' \
    '(-h --help)'{-h,--help}'[help for app]' \
    '(-n --name)'{-n,--name}'[Name of the Application to be deleted]:' \
    '(-o --owner)'{-o,--owner}'[Owner of the Application to be deleted]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_delete_help {
  _arguments \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}


function _apictl_export {
  local -a commands

  _arguments -C \
    '(-h --help)'{-h,--help}'[help for export]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "api-product:Export API Product"
      "help:Help about any command"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  api-product)
    _apictl_export_api-product
    ;;
  help)
    _apictl_export_help
    ;;
  esac
}

function _apictl_export_api-product {
  _arguments \
    '(-e --environment)'{-e,--environment}'[Environment to which the API Product should be exported]:' \
    '--format[File format of exported archive (json or yaml)]:' \
    '(-h --help)'{-h,--help}'[help for api-product]' \
    '(-n --name)'{-n,--name}'[Name of the API Product to be exported]:' \
    '(-r --provider)'{-r,--provider}'[Provider of the API Product]:' \
    '(-v --version)'{-v,--version}'[Version of the API Product to be exported]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_export_help {
  _arguments \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_export-api {
  _arguments \
    '(-e --environment)'{-e,--environment}'[Environment to which the API should be exported]:' \
    '--format[File format of exported archive(json or yaml)]:' \
    '(-h --help)'{-h,--help}'[help for export-api]' \
    '(-n --name)'{-n,--name}'[Name of the API to be exported]:' \
    '--preserveStatus[Preserve API status when exporting. Otherwise API will be exported in CREATED status]' \
    '(-r --provider)'{-r,--provider}'[Provider of the API]:' \
    '(-v --version)'{-v,--version}'[Version of the API to be exported]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_export-apis {
  _arguments \
    '(-e --environment)'{-e,--environment}'[Environment from which the APIs should be exported]:' \
    '--force[Clean all the previously exported APIs of the given target tenant, in the given environment if any, and to export APIs from beginning]' \
    '--format[File format of exported archives(json or yaml)]:' \
    '(-h --help)'{-h,--help}'[help for export-apis]' \
    '--preserveStatus[Preserve API status when exporting. Otherwise API will be exported in CREATED status]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_export-app {
  _arguments \
    '(-e --environment)'{-e,--environment}'[Environment to which the Application should be exported]:' \
    '(-h --help)'{-h,--help}'[help for export-app]' \
    '(-n --name)'{-n,--name}'[Name of the Application to be exported]:' \
    '(-o --owner)'{-o,--owner}'[Owner of the Application to be exported]:' \
    '--withKeys[Export keys for the application ]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_get-keys {
  _arguments \
    '(-e --environment)'{-e,--environment}'[Key generation environment]:' \
    '(-h --help)'{-h,--help}'[help for get-keys]' \
    '(-n --name)'{-n,--name}'[API or API Product to generate keys]:' \
    '(-r --provider)'{-r,--provider}'[Provider of the API or API Product]:' \
    '(-t --token)'{-t,--token}'[Token endpoint URL of Environment]:' \
    '(-v --version)'{-v,--version}'[Version of the API or API Product]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_help {
  _arguments \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}


function _apictl_import {
  local -a commands

  _arguments -C \
    '(-h --help)'{-h,--help}'[help for import]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "api-product:Import API Product"
      "help:Help about any command"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  api-product)
    _apictl_import_api-product
    ;;
  help)
    _apictl_import_help
    ;;
  esac
}

function _apictl_import_api-product {
  _arguments \
    '(-e --environment)'{-e,--environment}'[Environment from the which the API Product should be imported]:' \
    '(-f --file)'{-f,--file}'[Name of the API Product to be imported]:' \
    '(-h --help)'{-h,--help}'[help for api-product]' \
    '--import-apis[Import dependent APIs associated with the API Product]' \
    '--preserve-provider[Preserve existing provider of API Product after importing]' \
    '--skipCleanup[Leave all temporary files created during import process]' \
    '--update-api-product[Update an existing API Product or create a new API Product]' \
    '--update-apis[Update existing dependent APIs associated with the API Product]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_import_help {
  _arguments \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_import-api {
  _arguments \
    '(-e --environment)'{-e,--environment}'[Environment from the which the API should be imported]:' \
    '(-f --file)'{-f,--file}'[Name of the API to be imported]:' \
    '(-h --help)'{-h,--help}'[help for import-api]' \
    '--params[Provide a API Manager params file]:' \
    '--preserve-provider[Preserve existing provider of API after importing]' \
    '--skipCleanup[Leave all temporary files created during import process]' \
    '--update[Update an existing API or create a new API]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_import-app {
  _arguments \
    '(-e --environment)'{-e,--environment}'[Environment from the which the Application should be imported]:' \
    '(-f --file)'{-f,--file}'[Name of the ZIP file of the Application to be imported]:' \
    '(-h --help)'{-h,--help}'[help for import-app]' \
    '(-o --owner)'{-o,--owner}'[Name of the target owner of the Application as desired by the Importer]:' \
    '--preserveOwner[Preserves app owner]' \
    '--skipCleanup[Leave all temporary files created during import process]' \
    '--skipKeys[Skip importing keys of the Application]' \
    '(-s --skipSubscriptions)'{-s,--skipSubscriptions}'[Skip subscriptions of the Application]' \
    '--update[Update the Application if it is already imported]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_init {
  _arguments \
    '(-d --definition)'{-d,--definition}'[Provide a YAML definition of API]:' \
    '(-f --force)'{-f,--force}'[Force create project]' \
    '(-h --help)'{-h,--help}'[help for init]' \
    '--initial-state[Provide the initial state of the API; Valid states: [CREATED PUBLISHED]]:' \
    '--oas[Provide an OpenAPI specification file for the API]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}


function _apictl_install {
  local -a commands

  _arguments -C \
    '(-h --help)'{-h,--help}'[help for install]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "api-operator:Install API Operator"
      "help:Help about any command"
      "wso2am-operator:Install WSO2AM Operator"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  api-operator)
    _apictl_install_api-operator
    ;;
  help)
    _apictl_install_help
    ;;
  wso2am-operator)
    _apictl_install_wso2am-operator
    ;;
  esac
}

function _apictl_install_api-operator {
  _arguments \
    '(-f --from-file)'{-f,--from-file}'[Path to API Operator directory]:' \
    '(-h --help)'{-h,--help}'[help for api-operator]' \
    '(-c --key-file)'{-c,--key-file}'[Credentials file]:' \
    '(-p --password)'{-p,--password}'[Password of the given user]:' \
    '--password-stdin[Prompt for password of the given user in the stdin]' \
    '(-R --registry-type)'{-R,--registry-type}'[Registry type: DOCKER_HUB | AMAZON_ECR |GCR | HTTP]:' \
    '(-r --repository)'{-r,--repository}'[Repository name or URI]:' \
    '(-u --username)'{-u,--username}'[Username of the repository]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_install_help {
  _arguments \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_install_wso2am-operator {
  _arguments \
    '(-f --from-file)'{-f,--from-file}'[Path to wso2am-operator directory]:' \
    '(-h --help)'{-h,--help}'[help for wso2am-operator]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}


function _apictl_list {
  local -a commands

  _arguments -C \
    '(-h --help)'{-h,--help}'[help for list]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "api-products:Display a list of API Products in an environment"
      "apis:Display a list of APIs in an environment"
      "apps:Display a list of Applications in an environment specific to an owner"
      "envs:Display the list of environments"
      "help:Help about any command"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  api-products)
    _apictl_list_api-products
    ;;
  apis)
    _apictl_list_apis
    ;;
  apps)
    _apictl_list_apps
    ;;
  envs)
    _apictl_list_envs
    ;;
  help)
    _apictl_list_help
    ;;
  esac
}

function _apictl_list_api-products {
  _arguments \
    '(-e --environment)'{-e,--environment}'[Environment to be searched]:' \
    '--format[Pretty-print API Products using Go Templates. Use "{{ jsonPretty . }}" to list all fields]:' \
    '(-h --help)'{-h,--help}'[help for api-products]' \
    '(-l --limit)'{-l,--limit}'[Maximum number of API Products to return]:' \
    '(-q --query)'{-q,--query}'[Query pattern]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_list_apis {
  _arguments \
    '(-e --environment)'{-e,--environment}'[Environment to be searched]:' \
    '--format[Pretty-print apis using Go Templates. Use "{{ jsonPretty . }}" to list all fields]:' \
    '(-h --help)'{-h,--help}'[help for apis]' \
    '(-l --limit)'{-l,--limit}'[Maximum number of apis to return]:' \
    '(-q --query)'{-q,--query}'[Query pattern]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_list_apps {
  _arguments \
    '(-e --environment)'{-e,--environment}'[Environment to be searched]:' \
    '--format[Pretty-print outputusing Go templates. Use "{{jsonPretty .}}" to list all fields]:' \
    '(-h --help)'{-h,--help}'[help for apps]' \
    '(-l --limit)'{-l,--limit}'[Maximum number of applications to return]:' \
    '(-o --owner)'{-o,--owner}'[Owner of the Application]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_list_envs {
  _arguments \
    '--format[Pretty-print environments using go templates]:' \
    '(-h --help)'{-h,--help}'[help for envs]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_list_help {
  _arguments \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_login {
  _arguments \
    '(-h --help)'{-h,--help}'[help for login]' \
    '(-p --password)'{-p,--password}'[Password for login]:' \
    '--password-stdin[Get password from stdin]' \
    '(-u --username)'{-u,--username}'[Username for login]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_logout {
  _arguments \
    '(-h --help)'{-h,--help}'[help for logout]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}


function _apictl_remove {
  local -a commands

  _arguments -C \
    '(-h --help)'{-h,--help}'[help for remove]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "env:Remove Environment from Config file"
      "help:Help about any command"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  env)
    _apictl_remove_env
    ;;
  help)
    _apictl_remove_help
    ;;
  esac
}

function _apictl_remove_env {
  _arguments \
    '(-h --help)'{-h,--help}'[help for env]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_remove_help {
  _arguments \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_set {
  _arguments \
    '--export-directory[Path to directory where APIs should be saved]:' \
    '(-h --help)'{-h,--help}'[help for set]' \
    '--http-request-timeout[Timeout for HTTP Client]:' \
    '(-m --mode)'{-m,--mode}'[If mode is set to "k8s", apictl is capable of executing Kubectl commands. For example "apictl get pods" -> "kubectl get pods". To go back to the default mode, set the mode to "default"]:' \
    '(-t --token-type)'{-t,--token-type}'[Type of the token to be generated]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}


function _apictl_uninstall {
  local -a commands

  _arguments -C \
    '(-h --help)'{-h,--help}'[help for uninstall]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "api-operator:Uninstall API Operator"
      "help:Help about any command"
      "wso2am-operator:Uninstall WSO2AM Operator"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  api-operator)
    _apictl_uninstall_api-operator
    ;;
  help)
    _apictl_uninstall_help
    ;;
  wso2am-operator)
    _apictl_uninstall_wso2am-operator
    ;;
  esac
}

function _apictl_uninstall_api-operator {
  _arguments \
    '--force[Force uninstall API Operator]' \
    '(-h --help)'{-h,--help}'[help for api-operator]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_uninstall_help {
  _arguments \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_uninstall_wso2am-operator {
  _arguments \
    '--force[Force uninstall WSO2AM Operator]' \
    '(-h --help)'{-h,--help}'[help for wso2am-operator]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}


function _apictl_update {
  local -a commands

  _arguments -C \
    '(-h --help)'{-h,--help}'[help for update]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "api:handle APIs in kubernetes cluster "
      "help:Help about any command"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  api)
    _apictl_update_api
    ;;
  help)
    _apictl_update_help
    ;;
  esac
}

function _apictl_update_api {
  _arguments \
    '(*-f *--from-file)'{\*-f,\*--from-file}'[Path to swagger file]:' \
    '(-h --help)'{-h,--help}'[help for api]' \
    '(-m --mode)'{-m,--mode}'[Property to override the deploying mode. Available modes: privateJet, sidecar]:' \
    '(-n --name)'{-n,--name}'[Name of the API]:' \
    '--namespace[namespace of API]:' \
    '--replicas[replica set]:' \
    '(-v --version)'{-v,--version}'[Property to override the existing docker image with same name and version]:' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_update_help {
  _arguments \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}


function _apictl_vcs {
  local -a commands

  _arguments -C \
    '(-h --help)'{-h,--help}'[help for vcs]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "deploy:Deploys projects to the specified environment"
      "help:Help about any command"
      "status:Shows the list of projects that are ready to deploy"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  deploy)
    _apictl_vcs_deploy
    ;;
  help)
    _apictl_vcs_help
    ;;
  status)
    _apictl_vcs_status
    ;;
  esac
}

function _apictl_vcs_deploy {
  _arguments \
    '(-e --environment)'{-e,--environment}'[Name of the environment to deploy the project(s)]:' \
    '(-h --help)'{-h,--help}'[help for deploy]' \
    '--skipRollback[Specifies whether rolling back to the last successful revision during an error situation should be skipped]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_vcs_help {
  _arguments \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_vcs_status {
  _arguments \
    '(-e --environment)'{-e,--environment}'[Name of the environment to check the project(s) status]:' \
    '(-h --help)'{-h,--help}'[help for status]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

function _apictl_version {
  _arguments \
    '(-h --help)'{-h,--help}'[help for version]' \
    '(-k --insecure)'{-k,--insecure}'[Allow connections to SSL endpoints without certs]' \
    '--verbose[Enable verbose mode]'
}

