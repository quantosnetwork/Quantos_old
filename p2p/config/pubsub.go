package config

type PubsubConfig struct {

	Router string

	// DisableSigning disables message signing. Message signing is *enabled*
	// by default.
	DisableSigning bool

	// Enable pubsub
	Enabled Flag `json:",omitempty"`
}
