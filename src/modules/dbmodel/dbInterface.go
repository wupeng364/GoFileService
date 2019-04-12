package dbmodel

/**
 * 数据库操作接口
 */
type dbInterface interface{

	InitDB( dbAddr, dbPort, dbLib, dbUser, dbPW string) error
	Insert(sql string, params []interface{}) error
	Update(sql string, params []interface{}) error
	Delete(sql string, params interface{}) error
	Query(sql string, params []interface{})(interface{}, error)
	GetCon( col string ) (interface{}, error)
	DelTable(sql string) error
}