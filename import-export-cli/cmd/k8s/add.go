package k8s

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const K8sAddCmdLiteral = "add"
const k8sAddCmdShortDesc = "Add an API to the kubernetes cluster"
const k8sAddCmdLongDesc = `Add an API either from a Swagger file or project zip to the kubernetes cluster. JSON, YAML and zip formats are accepted.`
const k8sAddCmdExamples = utils.ProjectName + " " + K8sCmdLiteral + " " + K8sAddCmdLiteral + " " + AddApiCmdLiteral + " " + `-n petstore --from-file=./Swagger.json --replicas=1 --namespace=wso2

` + utils.ProjectName + " " + K8sCmdLiteral + " " + K8sAddCmdLiteral + " " + AddApiCmdLiteral + " " + `-n petstore --from-file=./product-apim-tooling/import-export-cli/build/target/apictl/myapi --replicas=1 --namespace=wso2 --override=true`

// K8sAddCmd represents the add command
var K8sAddCmd = &cobra.Command{
	Use:     K8sAddCmdLiteral,
	Short:   k8sAddCmdShortDesc,
	Long:    k8sAddCmdLongDesc,
	Example: k8sAddCmdExamples,
}

func init() {
	K8sCmd.AddCommand(K8sAddCmd)
}
