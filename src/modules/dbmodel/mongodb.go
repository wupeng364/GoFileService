package dbmodel
/**
 *@description Mongodb实例
 *@author	wupeng364@outlook.com
*/
import (
	//"fmt"
	"errors"
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	// "time"
)

type MongoDB struct{
	ms 		*mgo.Session
	dbLib 	string
}

func ( md *MongoDB )InitDB( dbAddr, dbPort, dbLib, dbUser, dbPW string) error{
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
			Addrs:[]string{dbAddr+":"+dbPort,}, 
			Source:dbLib, 
			Username:dbUser, 
			Password:dbPW })
	if err != nil {
		return err
	}
	md.ms = session
	md.dbLib = dbLib
	// md.db = session.DB(dbLib)
	return nil
}

func ( md *MongoDB )Insert(col string, params []interface{}) (err error){
	con, sion, err := getConnection(md, col)
	defer func( ){
		if sion != nil {
			sion.Close( )
		}
	}( )
	if err != nil {
		return
	}
	_len := len(params)
	if _len == 0 {
		err = errors.New("params is nil")
		return
	} else if _len == 1{
		err = con.Insert( params[0] )
	}else{
		for _, temp := range params{
			if err != nil {
				return
			}
			err = con.Insert( temp )
		}
	}
	return
}

func ( md *MongoDB )Update(col string, params []interface{}) error{
	con, sion, err := getConnection(md, col)
	defer func( ){
		if sion != nil {
			sion.Close( )
		}
	}( )
	if err != nil {
		return err
	}
	_, err = con.UpdateAll(params[0], params[1])
	return err
}

func ( md *MongoDB )Delete(col string, params interface{}) error{
	con, sion, err := getConnection(md, col)
	defer func( ){
		if sion != nil {
			sion.Close( )
		}
	}( )
	if err != nil {
		return err
	}
	_, err = con.RemoveAll(params)
	return err
}
func ( md *MongoDB )DelTable(col string) error{
	con, sion, err := getConnection(md, col)
	defer func( ){
		if sion != nil {
			sion.Close( )
		}
	}( )
	if err != nil {
		return err
	}
	return con.DropCollection( )
}

func ( md *MongoDB )Query(col string, params []interface{})( res interface{}, err error){
	con, sion, err := getConnection(md, col)
	defer func( ){
		if sion != nil {
			sion.Close( )
		}
	}( )
	if err != nil {
		return
	}
	_len := len(params)
	var result []interface{}
	query := con.Find(params[0])
	if _len > 1{
		query = query.Select(params[1])
	}
	if _len > 2{
		query = query.Skip(params[2].(int))
	}
	if _len > 3{
		query = query.Limit(params[2].(int))
	}
	query.All(&result)
	//fmt.Println( params, result )
	return result, nil
}

func (md *MongoDB)GetCon( col string )(interface{}, error){
	con, _, err := getConnection(md, col)
	return con, err
}
// ===================
// 获取一个连接
func getConnection( md *MongoDB, col string )(*mgo.Collection, *mgo.Session, error){
	if len(col) ==0{
		return nil, nil, errors.New("collection is empty!")
	}
	ms := md.ms.Copy( )
	ms.SetMode(mgo.Monotonic, true)
	return ms.DB( md.dbLib ).C( col ), ms, nil
}
