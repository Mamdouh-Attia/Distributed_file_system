package mastertracker

// Record represents a single record in the master
type Record struct {
	FileName       string
	DataKeeperNode string
	FilePath       string
	IsDataNodeAlive bool
}

// Master represents the master data structure containing records
type Mastertracker struct {
	Records []Record
}	

