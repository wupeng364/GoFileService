package sqlite
/**
 *@description 用户管理模块
 *提供系统数据库库管理
 *@author	wupeng364@outlook.com
*/
import (
	"fmt"
	"strings"
	"gofs/common/moduleloader"
	"path/filepath"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteModule struct{
	dbSource string
	mctx *moduleloader.Loader
}

// 返回模块信息
func (this *SqliteModule)ModuleOpts( )(moduleloader.Opts){
	return moduleloader.Opts{
		Name: "SqliteModule",
		Version: 1.0,
		Description: "Sqlite模块",
		OnReady: func (mctx *moduleloader.Loader){
			this.mctx = mctx
		},
		OnInit: this.onMInit,
	}
}

// 每次启动加载模块执行一次
func (this *SqliteModule)onMInit( ) {
	path, err := filepath.Abs(strings.Replace(cfg_db_path, "{name}", this.mctx.GetLoaderName( ), -1))
	if nil != err {
		panic(err)
	}
	this.dbSource = path;
	
	fmt.Println("   > SqliteModule dbSource="+ this.dbSource)
}
// ==============================================================================================
// 获取一个数据库连接
func (this *SqliteModule)Open( )(*sql.DB, error){
	return sql.Open("sqlite3", this.dbSource)
}