package config

type DNS struct {
	Resolvers map[string]string
	MaxCacheTTL int64 `json:", omitempty"`
}
