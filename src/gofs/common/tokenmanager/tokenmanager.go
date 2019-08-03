package tokenmanager

/**
 *@description 简单令牌管理器[生成一个自动销毁的令牌]
 *@author	wupeng364@outlook.com
*/
import(
	//"runtime"
	"time"
	"sync"
	"os"
	"strings"
	"strconv"
	"math/rand"
    "crypto/md5"
	"encoding/hex"
    "encoding/binary"
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
func(this *TokenManager) Init( )*TokenManager{
	this.tokenMap = make(map[string]*TokenObject)
	this.tokenLock = new(sync.RWMutex)
	
	// 定期清理
	go func( ){
		for{
			if len(this.tokenMap) >= 0 {
				this.tokenLock.Lock() // 读取时锁定map, 防止中途修改
				for key, val := range this.tokenMap {
					if val.expired == -1 {
						continue
					}
					now := time.Now( ).UnixNano( )
					if val.expired <= now {
						// fmt.Println("remove: ", key, val.expired - now, val)
						delete(this.tokenMap, key);
					}
				}
				this.tokenLock.Unlock()
			}
			//runtime.Gosched( )
			time.Sleep(time.Duration(1)*time.Second)
		}		
	}( )
	return this
}
// 生成令牌
// Second=-1时, 不会自动销毁内存中的信息
func(this *TokenManager) AskToken( tb *TokenObject ) string{
	token := getUUID( )
	tb.regtime = time.Now().UnixNano( )
	if tb.Second > -1 {
		tb.expired = tb.regtime + tb.Second*int64(time.Second)
	}else{
		tb.expired = -1
	}
	this.tokenLock.Lock()
	defer this.tokenLock.Unlock()
	this.tokenMap[token] = tb
	return token
}
// 销毁令牌
func(this *TokenManager) DestroyToken( tk string ){
	this.tokenLock.Lock()
	defer this.tokenLock.Unlock()
	
	if _, ok := this.tokenMap[tk]; ok {
		delete(this.tokenMap, tk)
	}
}
// 获取令牌信息
func(this *TokenManager) GetTokenInfo( tk string ) (*TokenObject, bool) {
	this.tokenLock.RLock()
	defer this.tokenLock.RUnlock()
	
	val,ok := this.tokenMap[tk]
	if ok {
		now := time.Now( ).UnixNano( )
		if val.expired <= now {
			return nil, false
		}
	}
	return val, ok
}
// 刷新|重置令牌过期时间
func(this *TokenManager) RefreshToken( tk string ){
	this.tokenLock.Lock()
	defer this.tokenLock.Unlock()
	
	val,ok := this.tokenMap[tk]
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


// >================== uuid start ==================<
// 获取机器唯一标识=主机名+进程ID+随机数=>MD5
var _machineid []byte = func( )(machHash []byte) {
	host, _ := os.Hostname()
    machine := strings.Join([]string{
        host,
        strconv.FormatInt(int64(os.Getpid( )), 10),
        hex.EncodeToString(func(n uint32)[]byte {
		    uintByte := make([]byte, 4, 4)
		    binary.BigEndian.PutUint32(uintByte, n)
		    return uintByte
		}( uint32(rand.Int31( )) ) ),
    }, ",")
    
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(machine))
    machHash = md5Ctx.Sum(nil)
    return machHash
}( )
// 获取唯一ID
func getUUID( ) string {
	baseId := func( ) []byte {
	    id := make([]byte, 0, 36)
	    id = append(id, _machineid[0:4]...)
	    id = append(id, func (i uint64) []byte {
		    var buf = make([]byte, 8, 8)
		    binary.BigEndian.PutUint64(buf, i)
		    return buf
		}(uint64(time.Now().UnixNano( )))...)
	    id = append(id, func(n uint32)[]byte {
			    uintByte := make([]byte, 4, 4)
			    binary.BigEndian.PutUint32(uintByte, n)
			    return uintByte
		 }(uint32(rand.Int31()))...)
	    return id
	}( )
	
    gid := make([]byte, 0, 36)
    id := []byte(hex.EncodeToString(baseId))
    gid = append(gid, id[0:8]...)
    gid = append(gid, '-')
    gid = append(gid, id[8:12]...)
    gid = append(gid, '-')
    gid = append(gid, id[12:16]...)
    gid = append(gid, '-')
    gid = append(gid, id[16:20]...)
    gid = append(gid, '-')
    gid = append(gid, id[20:]...)
    return string(gid)
}
// >================== uuid end ==================<