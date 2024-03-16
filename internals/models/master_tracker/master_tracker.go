package mastertracker

// Record represents a single record in the master
type Record struct {
	FileName        string
	DataKeeperNode  string
	IsDataNodeAlive bool
}

// Master represents the master data structure containing records
type Master struct {
	Records map[string]Record // Map with file path as key
}

// NewMaster creates a new Master instance'
func NewMaster() *Master {
	return &Master{
		Records: make(map[string]Record),
	}
}

// AddRecord adds a new record to the master
func (m *Master) AddRecord(fileName, dataKeeperNode string, isDataNodeAlive bool) {
	m.Records[fileName] = Record{
		FileName:        fileName,
		DataKeeperNode:  dataKeeperNode,
		IsDataNodeAlive: isDataNodeAlive,
	}
}

// RemoveRecord removes a record from the master
func (m *Master) RemoveRecord(fileName string) {
	delete(m.Records, fileName)
}

// GetRecord returns a record from the master
func (m *Master) GetRecord(fileName string) Record {
	return m.Records[fileName]
}

// GetAllRecords returns all records from the master
func (m *Master) GetAllRecords() map[string]Record {
	return m.Records
}

// UpdateRecord updates a record in the master
func (m *Master) UpdateRecord(fileName, dataKeeperNode string, isDataNodeAlive bool) {
	m.Records[fileName] = Record{
		FileName:        fileName,
		DataKeeperNode:  dataKeeperNode,
		IsDataNodeAlive: isDataNodeAlive,
	}
}

// RemoveDataNodeRecords removes all records associated with a data node
func (m *Master) RemoveDataNodeRecords(dataNode string) {
	for fileName, record := range m.Records {
		if record.DataKeeperNode == dataNode {
			delete(m.Records, fileName)
		}
	}
}

// GetRecordsByDataNode returns all records associated with a data node
func (m *Master) GetRecordsByDataNode(dataNode string) map[string]Record {
	records := make(map[string]Record)
	for fileName, record := range m.Records {
		if record.DataKeeperNode == dataNode {
			records[fileName] = record
		}
	}
	return records
}

// GetAliveRecords returns all records associated with alive data nodes
func (m *Master) GetAliveRecords() map[string]Record {
	records := make(map[string]Record)
	for fileName, record := range m.Records {
		if record.IsDataNodeAlive {
			records[fileName] = record
		}
	}
	return records
}
