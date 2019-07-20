package userapi
import(
	"errors"
)

var Error_UserIdIsNil = errors.New("User ID cannot be empty")
var Error_UserNameIsNil = errors.New("User name cannot be empty")
var Error_PwdIsError  = errors.New("User passwork is error")
var Error_NotSupport  = errors.New("This type is not supported")
var Error_ParamsNotEmpty  = errors.New("Property name or content cannot be empty")
