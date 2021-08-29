# 功能说明

`func ConnectDb() (*mongo.Client, error)`

`func ClearCollection(table string) error` 

`func DeleteCollection(table string) error` 

`func CreateCollection(table string) error`

`func InitCollection(table string) error`

`func FindOne(db string, table string, data interface{}, filter bson.M) error`

`func FindMany(db string, table string, filter bson.M) (*mongo.Cursor, error)`

`func InsertOne(db string, table string, data interface{}) error`

`func UpdateOne(db string, table string, filter bson.M, data interface{}) error`

data 字段需要手动指定 bson.M{"$set": __} 命令, 必须是 {"$set": ___ } 类型 !!!
把某个 bson 结构体全修改: UpdateOnd(db, table, filter, bson.M{"$set": data})

`func DeleteOne(db string, table string, filter bson.M) error`

`func DeleteMany(dbName string, collectionName string, filter bson.M) (int, error)` 

第一个 int 表示删除元素个数
