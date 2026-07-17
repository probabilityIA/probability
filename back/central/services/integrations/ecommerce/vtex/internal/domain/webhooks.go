package domain

const (
	HookKeyHeader  = "X-VTEX-HOOK-KEY"
	HookFilterType = "FromWorkflow"
)

var WebhookOrderStates = []string{
	"payment-pending",
	"ready-for-handling",
	"canceled",
}

type HookConfig struct {
	URL      string
	Statuses []string
	HasKey   bool
}

type WebhookItem struct {
	ID       string
	Address  string
	Statuses []string
	IsOurs   bool
}
