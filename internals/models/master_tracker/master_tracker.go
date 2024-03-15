package mastertracker

// Record represents a single record in the master
type Record struct {
	FileName       string
	DataKeeperNode string
	IsDataNodeAlive bool
}

// Master represents the master data structure containing records
type Master struct {
	Records map[string]Record // Map with file path as key
}
