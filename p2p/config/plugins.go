package config

type Plugins struct {
	Plugins map[string]Plugin
}

type Plugin struct {
	Disabled bool
	Config interface{}
}

