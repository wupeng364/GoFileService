package usermanage
import(
	"time"
	"errors"
)
const(
	cfc_db_type = "usermanage.dbType"
	UserType_Admin  = "0"
	UserType_Normal = "1"
)

var Error_ConnIsNil = errors.New("The data source is empty")
var Error_UserIdIsNil = errors.New("The userId is empty")
var Error_UserNotExist = errors.New("User does not exist")
var Error_UserNameIsNil = errors.New("The userName is empty")

// 用户表存储的结构
type UserInfo struct{
	UserType	int
	UserId 		string
	Username	string
	Userpwd		string
	Cttime		time.Time
}