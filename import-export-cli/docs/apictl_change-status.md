## apictl change-status

Change Status of an API, MCP Server or Product

### Synopsis

Change the lifecycle status of an API, MCP Server or API Product in an environment

```
apictl change-status [flags]
```

### Examples

```
apictl change-status api -a Publish -n TwitterAPI -v 1.0.0 -r admin -e dev
apictl change-status api -a Publish -n FacebookAPI -v 2.1.0 -e production
apictl change-status mcp-server -a Publish -n WeatherMCPServer -v 1.0.0 -r admin -e dev
apictl change-status api-product -a Publish -n SocialMediaProduct -v 1.0.0 -r admin -e dev
```

### Options

```
  -h, --help   help for change-status
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications and Managing WSO2 Micro Integrator
* [apictl change-status api](apictl_change-status_api.md)	 - Change Status of an API
* [apictl change-status api-product](apictl_change-status_api-product.md)	 - Change Status of an API Product
* [apictl change-status mcp-server](apictl_change-status_mcp-server.md)	 - Change Status of an MCP Server

