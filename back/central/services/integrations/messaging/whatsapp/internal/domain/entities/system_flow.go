package entities

type SystemFlowNodeKind string

const (
	SystemFlowNodeTemplate SystemFlowNodeKind = "template"
	SystemFlowNodeDecision SystemFlowNodeKind = "decision"
)

type SystemFlowNode struct {
	Kind     SystemFlowNodeKind
	Template string
	Title    string
	Detail   string
	State    ConversationState
	Event    string
	Effects  []string
	Branches []SystemFlowBranch
}

type SystemFlowBranch struct {
	Label  string
	Node   *SystemFlowNode
	BackTo string
}

type SystemFlow struct {
	Key         string
	Name        string
	Description string
	Trigger     string
	Root        *SystemFlowNode
}
