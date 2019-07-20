package stringtools

/**
 *@description 字符处理工具
 *@author	wupeng364@outlook.com
*/
import (
	"io"
	"os"
	"time"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
    "crypto/md5"
	"encoding/hex"
    "encoding/binary"
    mrand "math/rand"
)
var ClerPath = path.Clean
var machineId []byte = GetMachineHash( )


func Uint32ToHexString(n uint32) string {
    return ByteToHexString(Uint32ToBytes(n))
}
func Uint32ToBytes(n uint32) []byte {
    uintByte := make([]byte, 4, 4)
    binary.BigEndian.PutUint32(uintByte, n)
    return uintByte
}
func Uint64ToBytes(i uint64) []byte {
    var buf = make([]byte, 8, 8)
    binary.BigEndian.PutUint64(buf, i)
    return buf
}
func BytesToInt64(buf []byte) int64 {
    return int64(binary.BigEndian.Uint64(buf))
}
func ByteToHexString(n []byte) string {
    return hex.EncodeToString(n)
}
func Rand() uint32 {
    return uint32(mrand.Int31())
}
// Int => Int64 => string
func Int2String( i int64 ) string{
	return strconv.FormatInt(i, 10)
}
// 范围最后一个'/'前的文字
func GetParentPath( s_path string )string{
	if strings.Index(s_path, "\\") > -1 {
		return s_path[:strings.LastIndex(s_path, "\\")]
	}else{
		return s_path[:strings.LastIndex(s_path, "/")]
	}
}
// 范围最后一个'/'后的文字
func GetPathName( s_path string )string{
	if strings.Index(s_path, "\\") > -1 {
		return s_path[strings.LastIndex(s_path, "\\")+1:]
	}else{
		return s_path[strings.LastIndex(s_path, "/")+1:]
	}
}
// 读取文字
func ReadString(src io.Reader)string{
	if nil == src {
		return ""
	}
	buf := make([]byte, 0)
	for {
		buf_temp := make([]byte, 1024)
		nr, er := src.Read(buf_temp)
		if nr > 0 {
			buf = append(buf, buf_temp[:nr]...)
		}
		if er != nil {
			if er != io.EOF {
				return ""
			}else{
				break
			}
		}
	}
	return string(buf)
}
// 去除数组中指定的字符, 返回新的数组
func ArrayClear( strs []string, str string ){
}
// 删除路径后面 /, 把\转换为/
func UnixPathClear( str string ) string{
	if len(strings.Replace(str, " ", "", -1) ) == 0 {
		return ""
	}
	return ClerPath(strings.Replace(str, "\\", "/", -1))
}
// 
func ReplaceAll(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}
// 获取函数名称
func GetFunctionName(i interface{}, seps ...rune) string {
    fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer( )).Name()
    // 用 seps 进行分割
    fields := strings.FieldsFunc(fn, func(sep rune) bool {
        for _, s := range seps {
            if sep == s {
                return true
            }
        }
        return false
    })

    // fmt.Println(fields)

    if size := len(fields); size > 0 {
        return fields[size-1]
    }
    return ""
}
func Contains(s, sep string) bool {
	n := len(sep)
	c := sep[0]
	for i := 0; i+n <= len(s); i++ {
		if s[i] == c && s[i:i+n] == sep {
			return true
		}
	}
	return false
}
func ParseBool(str string) bool {
    switch str {
    case "1", "t", "T", "true", "TRUE", "True":
        return true
    case "0", "f", "F", "false", "FALSE", "False":
        return false
    }
    return false
}
func FormatBool(b bool) string {
    if b {
        return "true"
    }
    return "false"
}
func GetUUID( ) string {
    gid := make([]byte, 0, 36)
    id := []byte(ByteToHexString(createBaseId()))
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
func GetSimpleUUID() string {
    return ByteToHexString(createBaseId())
}
func createBaseId( ) []byte {
    id := make([]byte, 0, 36)
    id = append(id, machineId[0:4]...)
    id = append(id, Uint64ToBytes(uint64(time.Now().UnixNano( )))...)
    id = append(id, Uint32ToBytes(Rand())...)
    return id
}
// 获取机器唯一标识-主机名+进程ID+随机数=>MD5
func GetMachineHash() (machHash []byte) {
    machine := strings.Join([]string{
        GetHostname( ),
        Getpid( ),
        Uint32ToHexString(Rand( )),
    }, ",")
    
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(machine))
    machHash = md5Ctx.Sum(nil)
    return
}
func GetHostname( )string {
    host, err := os.Hostname()
    if err != nil {
        host = ""
    }
    return host
}
func Getpid( )string {
  return Int2String(int64(os.Getpid( )))
}
func GetTimestamp( ) int64{
	return time.Now().UnixNano( ) / 1e6
}
// 字符转MD5
func String2MD5( str string )string{
	md5Ctx := md5.New( )
    md5Ctx.Write([]byte(str))
    return hex.EncodeToString(md5Ctx.Sum(nil))
}