package tokenmanager

/**
 *@description 简单令牌管理器[生成一个自动销毁的令牌]
 *@author	wupeng364@outlook.com
*/
import(
	//"time"
	"fmt"
	"testing"
)
var tokenManager *TokenManager
var tokens []string

func init(){
	fmt.Println("tokenManager Test init ...")
	tokenManager = &TokenManager{}
	tokenManager.Init()
	tokens = make([]string, 0)
}
// 生成令牌
func TestAskToken( t *testing.T ){
	for i:=0; i<=100000; i++ {
		if i <10{
			// 每次休眠一秒, 前面申请的部分就会过期
			// time.Sleep(time.Duration(2)*time.Second)
		}
		var tb *TokenObject
		tb = &TokenObject{Second:int64(60), TypeStr:"tokenManager_Test", TokenBody: i,}
		token := tokenManager.AskToken(tb)
		fmt.Println("AskToken: ", token, tb)
		tokens = append(tokens, token)
	}
}
// 获取令牌信息
func TestGetTokenInfo( t *testing.T ){
	if len(tokens) > 0 {
		for i, val := range tokens{
			tb, ok := tokenManager.GetTokenInfo(val)
			fmt.Println("GetTokenInfo: ", i, val, ok, tb)
		}
	}
}
// 刷新|重置令牌过期时间
func TestRefreshToken( t *testing.T ){
	 if len(tokens) > 0 {
		for _, val := range tokens{
			tb_old, ok := tokenManager.GetTokenInfo(val)
			if ok {
				expired_old := tb_old.expired
				tokenManager.RefreshToken(val)
				tb_new, _ := tokenManager.GetTokenInfo(val)
				fmt.Println("RefreshToken: ", val, tb_new.expired - expired_old)
			}
		}
	}
}