package dbmodel
/**
 *@description 数据库模块[此模块未使用]
 *@author	wupeng364@outlook.com
*/

import (
	"fmt"
	"errors"
	"modules/configmodel"
	Gmd "common/gomodule"
)

type DatabaseModel struct{
	db 	dbInterface
	cfg *configmodel.ConfigModel
}

// 返回模块信息
func (dbm *DatabaseModel)MInfo( )(*Gmd.ModelInfo)	{
	return &Gmd.ModelInfo{
		dbm,
		"DatabaseModel",
		1.0,
		"数据库模块",
	}
}
// 模块安装, 一个模块只初始化一次
func (dbm *DatabaseModel)MSetup( ) {
	
}
// 模块升级, 一个版本执行一次
func (dbm *DatabaseModel)MUpdate( ) {
	
}

// 每次启动加载模块执行一次
func (dbm *DatabaseModel)OnMInit( getPointer func(m interface{})interface{}  ) {

	dbm.cfg = getPointer(dbm.cfg).(*configmodel.ConfigModel)
	dbType := dbm.cfg.GetConfig(cfg_db_type)
	dbAddr := dbm.cfg.GetConfig(cfg_db_addr) 
	dbPort := dbm.cfg.GetConfig(cfg_db_port)
	dbLib  := dbm.cfg.GetConfig(cfg_db_db)
	dbUser := dbm.cfg.GetConfig(cfg_db_user)
	dbPW   := dbm.cfg.GetConfig(cfg_db_pwd)
	
	// 初始化数据库接口
	dbm.doInitInterface( dbType, dbAddr, dbPort, dbLib, dbUser, dbPW )
}

// 系统执行销毁时执行
func (dbm *DatabaseModel)OnMDestroy( ) {
	
}

// ==============================================================================================
func (dbm *DatabaseModel)Insert(sql string, params ...interface{}) error{
	return dbm.db.Insert(sql, params)
}
func (dbm *DatabaseModel)Update(sql string, params ...interface{}) error{
	return dbm.db.Update(sql, params)
}
func (dbm *DatabaseModel)Query(sql string, params ...interface{})(res interface{}, err error){
	res, err = dbm.db.Query(sql, params)
	return res, err
}
func (dbm *DatabaseModel)Delete(sql string, params interface{}) error{
	return dbm.db.Delete(sql, params)
}
func (dbm *DatabaseModel)DelTable(sql string)(err error){
	return dbm.db.DelTable(sql)
}
// ==============================================================================================
func (dbm *DatabaseModel)doInitInterface( dbType, dbAddr, dbPort, dbLib, dbUser, dbPW string  ){
	if dbType == "MongoDB" {
		dbm.db = &MongoDB{}
	}
	if dbm.db == nil {
		panic(errors.New("No support method for "+dbType+" type database"))
	}
	fmt.Println("doInitInterface: ", dbType, dbAddr, dbPort, dbLib, dbUser, dbPW)
	err := dbm.db.InitDB(dbAddr, dbPort, dbLib, dbUser, dbPW)
	if err != nil {
		panic( err )
	}
}