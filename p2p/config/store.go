package config

const DefaultDataStoreDirectory = "data"

type DataStore struct {

	StorageMax string
	StorageGCWatermark int64 // percentage * StorageMax
	GCPeriod string
	Spec map[string]interface{}
	HashOnRead bool
	BloomFilterSize int

}

func DataStorePath(configroot string) (string, error) {
	return Path(configroot, DefaultDataStoreDirectory)
}