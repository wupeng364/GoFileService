package tokenmanager

/**
 *@description 简单令牌管理器[生成一个自动销毁的令牌]
 *@author	wupeng364@outlook.com
*/
import(
	"common/stringtools"
	"runtime"
	"time"
	"sync"
	//"fmt"
)
type TokenObject struct{
	regtime   int64
	expired   int64
	Second	  int64
	TypeStr	  string
	TokenBody interface{}
}
// 
type TokenManager struct{
	tokenMap map[string]*TokenObject
	tokenLock *sync.RWMutex
}
// 初始化-启动一个管理线程, 负责令牌的生命周期
func(tm *TokenManager) Init( )*TokenManager{
	tm.tokenMap = make(map[string]*TokenObject)
	tm.tokenLock = new(sync.RWMutex)
	
	// 定期清理
	go func( ){
		for{
			if len(tm.tokenMap) >= 0 {
				tm.tokenLock.Lock() // 读取时锁定map, 防止中途修改
				for key, val := range tm.tokenMap {
					if val.expired == -1 {
						continue
					}
					now := time.Now( ).UnixNano( )
					if val.expired <= now {
						// fmt.Println("remove: ", key, val.expired - now, val)
						delete(tm.tokenMap, key);
					}
				}
				tm.tokenLock.Unlock()
			}
			runtime.Gosched( )
			time.Sleep(time.Duration(1)*time.Nanosecond)
		}		
	}( )
	return tm
}
// 生成令牌
// Second=-1时, 不会自动销毁内存中的信息
func(tm *TokenManager) AskToken( tb *TokenObject ) string{
	token := stringtools.GetUUID( )
	tb.regtime = time.Now().UnixNano( )
	if tb.Second > -1 {
		tb.expired = tb.regtime + tb.Second*int64(time.Second)
	}else{
		tb.expired = -1
	}
	tm.tokenLock.Lock()
	defer tm.tokenLock.Unlock()
	tm.tokenMap[token] = tb
	return token
}
// 销毁令牌
func(tm *TokenManager) DestroyToken( tk string ){
	tm.tokenLock.Lock()
	defer tm.tokenLock.Unlock()
	
	if _, ok := tm.tokenMap[tk]; ok {
		delete(tm.tokenMap, tk)
	}
}
// 获取令牌信息
func(tm *TokenManager) GetTokenInfo( tk string ) (*TokenObject, bool) {
	tm.tokenLock.RLock()
	defer tm.tokenLock.RUnlock()
	
	val,ok := tm.tokenMap[tk]
	if ok {
		now := time.Now( ).UnixNano( )
		if val.expired <= now {
			return nil, false
		}
	}
	return val, ok
}
// 刷新|重置令牌过期时间
func(tm *TokenManager) RefreshToken( tk string ){
	tm.tokenLock.Lock()
	defer tm.tokenLock.Unlock()
	
	val,ok := tm.tokenMap[tk]
	if ok {
		now := time.Now( ).UnixNano( )
		if val.expired <= now {
			return 
		}
		used := val.expired - val.regtime
		val.regtime = time.Now().UnixNano( )
		val.expired = val.regtime + used
	}
}

