package monitor

type CollInfo struct {
	Size           int `json:"size" bson:"size"`
	Count          int `json:"count" bson:"count"`
	AvgObjSize     int `json:"avgObjSize" bson:"avgObjSize"`
	StorageSize    int `json:"storageSize" bson:"storageSize"`
	TotalIndexSize int `json:"totalIndexSize" bson:"totalIndexSize"`
}
