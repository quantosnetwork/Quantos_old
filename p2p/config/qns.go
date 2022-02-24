package config

type Qns struct {
	RepublishPeriod string
	RecordLifetime string
	ResolveCacheSize int
	UsePubsub Flag `json:",omitempty"`
}

