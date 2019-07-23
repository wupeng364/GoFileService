package sqlite
/**
 *@description 用户管理模块
 *提供系统数据库库管理
 *@author	wupeng364@outlook.com
*/
import (
	"fmt"
	"gofs/common/gomodule"
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
func (s *SqliteModule)MInfo( )(*gomodule.ModuleInfo)	{
	return &gomodule.ModuleInfo{
		s,
		"SqliteModule",
		1.0,
		"Sqlite模块",
	}
}

// 模块安装, 一个模块只初始化一次
func (s *SqliteModule)OnMSetup( ref gomodule.ReferenceModule ) {
	
}
// 模块升级, 一个版本执行一次
func (s *SqliteModule)OnMUpdate( ref gomodule.ReferenceModule ) {
	
}

// 每次启动加载模块执行一次
func (s *SqliteModule)OnMInit( ref gomodule.ReferenceModule ) {
	s.cfg = ref(s.cfg).(*config.ConfigModule)
	path, err := filepath.Abs(cfg_db_path)
	if nil != err {
		panic(err)
	}
	s.dbSource = path;
	
	fmt.Println("   > SqliteModule dbSource="+ s.dbSource)
}

// 系统执行销毁时执行
func (s *SqliteModule)OnMDestroy( ref gomodule.ReferenceModule ) {
	
}

// ==============================================================================================
// 获取一个数据库连接
func (s *SqliteModule)Open( )(*sql.DB, error){
	return sql.Open("sqlite3", s.dbSource)
}