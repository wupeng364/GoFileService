package sqlite
/**
 *@description 用户管理模块
 *提供系统数据库库管理
 *@author	wupeng364@outlook.com
*/
import (
	"fmt"
	"gofs/common/moduletools"
	"gofs/modules/common/config"
	"path/filepath"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteModule struct{
	dbSource string
	cfg *config.ConfigModule
}

// 返回模块信息
func (this *SqliteModule)MInfo( )(moduletools.ModuleInfo)	{
	return moduletools.ModuleInfo{
		"SqliteModule",
		1.0,
		"Sqlite模块",
	}
}

// 模块安装, 一个模块只初始化一次
func (this *SqliteModule)OnMSetup( ref moduletools.ReferenceModule ) {
	
}
// 模块升级, 一个版本执行一次
func (this *SqliteModule)OnMUpdate( ref moduletools.ReferenceModule ) {
	
}

// 每次启动加载模块执行一次
func (this *SqliteModule)OnMInit( ref moduletools.ReferenceModule ) {
	this.cfg = ref(this.cfg).(*config.ConfigModule)
	path, err := filepath.Abs(cfg_db_path)
	if nil != err {
		panic(err)
	}
	this.dbSource = path;
	
	fmt.Println("   > SqliteModule dbSource="+ this.dbSource)
}

// 系统执行销毁时执行
func (this *SqliteModule)OnMDestroy( ref moduletools.ReferenceModule ) {
	
}

// ==============================================================================================
// 获取一个数据库连接
func (this *SqliteModule)Open( )(*sql.DB, error){
	return sql.Open("sqlite3", this.dbSource)
}