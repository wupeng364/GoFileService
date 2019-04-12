package tokenmanager

/**
 *@description 简单令牌管理器[生成一个自动销毁的令牌]
 *@author	wupeng364@outlook.com
*/
import(
	"common/stringtools"
	"runtime"
	"time"
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
}
// 
func(tm *TokenManager) Init( )*TokenManager{
	tm.tokenMap = make(map[string]*TokenObject)
	// 定期清理
	go func( P_TM map[string]*TokenObject ){
		for{
			if len(P_TM) >= 0 {
				now := time.Now( ).UnixNano( ) / 1e6
				for key, val := range P_TM {
					// fmt.Println(len(P_TM), val.expired - now, key, val)
					if val.expired <= now {
						delete(P_TM, key)
					}
				}
			}
			runtime.Gosched( )
			time.Sleep(time.Duration(1)*time.Second)
		}		
	}(tm.tokenMap)
	return tm
}
// 
func(tm *TokenManager) AskToken( tb *TokenObject ) string{
	token := stringtools.GetUUID( )
	tb.regtime = time.Now().UnixNano( ) / 1e6
	tb.expired = tb.regtime + int64(tb.Second*1000)
	tm.tokenMap[token] = tb
	return token
}
// 
func(tm *TokenManager) DestroyToken( tk string ){
	if _, ok := tm.tokenMap[tk]; ok {
		delete(tm.tokenMap, tk)
	}
}
// 
func(tm *TokenManager) GetTokenInfo( tk string ) (*TokenObject, bool) {
	val,ok := tm.tokenMap[tk]
	return val, ok
}
// 
func(tm *TokenManager) RefreshToken( tk string ){
	val,ok := tm.tokenMap[tk]
	if ok {
		used := val.expired - val.regtime
		val.regtime = time.Now().UnixNano( ) / 1e6
		val.expired = val.regtime + used
	}
}

