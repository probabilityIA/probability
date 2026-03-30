package entities

type MonitoringUser struct {
	ID       uint
	Email    string
	Name     string
	ScopeID  *uint
	IsActive bool
}

func (u *MonitoringUser) IsPlatformScope() bool {
	return u.ScopeID != nil && *u.ScopeID == 1
}
