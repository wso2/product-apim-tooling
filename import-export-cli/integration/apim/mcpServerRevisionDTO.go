package apim

// MCPServerRevision : MCP Server Revision DTO
type MCPServerRevision struct {
	DisplayName    string                        `json:"displayName"`
	ID             string                        `json:"id"`
	Description    string                        `json:"description"`
	DeploymentInfo []MCPServerRevisionDeployment `json:"deploymentInfo"`
}

// MCPServerRevisionDeployment : MCP Server Revision Deployment DTO
type MCPServerRevisionDeployment struct {
	RevisionUUID       string `json:"revisionUuid"`
	Name               string `json:"name"`
	VHost              string `json:"vhost"`
	DisplayOnDevportal bool   `json:"displayOnDevportal"`
}

// MCPServerRevisionList : MCP Server Revisions List DTO
type MCPServerRevisionList struct {
	Count string              `json:"count"`
	List  []MCPServerRevision `json:"list"`
}
