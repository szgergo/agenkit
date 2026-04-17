package manifest

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

type Manifest struct {
	Version    int                  `yaml:"version"`
	Mode       Mode                 `yaml:"mode"`
	Vars       map[string]string    `yaml:"vars"`
	McpServers map[string]MCPServer `yaml:"mcp_servers"`
	Providers  map[string]Provider  `yaml:"providers"`
}

type MCPServer struct {
	Command           string                       `yaml:"command"`
	Args              []string                     `yaml:"args"`
	Env               map[string]string            `yaml:"env"`
	Providers         []string                     `yaml:"providers"`
	ProviderOverrides map[string]ProviderOverrides `yaml:"provider_overrides"`
}

type ProviderOverrides struct {
	Args []string `yaml:"args"`
}

type Provider struct {
	ConfigPath string `yaml:"config_path"`
	Enabled    *bool  `yaml:"enabled"`
}

type Mode string

const (
	ModeMerge   Mode = "merge"
	ModeReplace Mode = "replace"
	ModeExtend  Mode = "extend"
)

func (m *Mode) UnmarshalYAML(value *yaml.Node) error {
	var s string
	err := value.Decode(&s)

	if err != nil {
		return err
	}

	switch Mode(s) {
	case ModeExtend, ModeMerge, ModeReplace:
		*m = Mode(s)
		return nil
	}
	return fmt.Errorf("Invalid mode: %s valid values: %s,%s,%s", s, ModeExtend, ModeMerge, ModeReplace)

}
