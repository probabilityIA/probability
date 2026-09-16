package response

type SystemFlowNode struct {
	Kind        string             `json:"kind"`
	Template    string             `json:"template,omitempty"`
	Description string             `json:"description,omitempty"`
	Body        string             `json:"body,omitempty"`
	Buttons     []string           `json:"buttons,omitempty"`
	Title       string             `json:"title,omitempty"`
	Detail      string             `json:"detail,omitempty"`
	Effects     []string           `json:"effects,omitempty"`
	Branches    []SystemFlowBranch `json:"branches,omitempty"`
}

type SystemFlowBranch struct {
	Label  string          `json:"label"`
	Node   *SystemFlowNode `json:"node,omitempty"`
	BackTo string          `json:"back_to,omitempty"`
}

type SystemFlow struct {
	Key         string          `json:"key"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Trigger     string          `json:"trigger"`
	Root        *SystemFlowNode `json:"root"`
}
