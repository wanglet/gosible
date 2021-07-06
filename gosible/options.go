package gosible

type PrivilegeEscalationOptions struct{}

type ConnectionOptions struct{}

type PlaybookOptions struct{}

type AnsibleOptions struct {
	ModuleName string `argument:"-m" json:"moduleName"`
	ModuleArgs string `argument:"-a" json:"moduleArgs,omitempty"`
	Forks      int    `argument:"-f" json:"forks,omitempty"`
	Limit      string `argument:"-l" json:"limit,omitempty"`
}
