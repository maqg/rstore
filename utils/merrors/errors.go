package merrors

const (
	ERR_OCT_SUCCESS = iota
	ERR_DB_ERR
	ERR_NOT_ENOUGH_PARAS
	ERR_TOO_MANY_PARAS
	ERR_UNACCP_PARAS
	ERR_CMD_ERR
	ERR_COMMON_ERR
	ERR_SEGMENT_NOT_EXIST
	ERR_SEGMENT_ALREADY_EXIST
	ERR_TIMEOUT
	ERR_SYSCALL_ERR
	ERR_SYSTEM_ERR
	ERR_NO_SUCH_API
	ERR_NOT_IMPLEMENTED

	// User
	ERR_USER_NOT_EXIST
	ERR_USER_ALREADY_EXIST
	ERR_PASSWORD_DONT_MATCH
	ERR_USER_NOT_LOGIN
	ERR_USER_GROUPS_NOT_EMPTY

	// User Group
	ERR_USERGROUP_NOT_EXIST
	ERR_USERGROUP_ALREADY_EXIST
	ERR_USERGROUP_USERS_NOT_EMPTY
)

var GErrors = map[int]string{
	ERR_OCT_SUCCESS:           "Command Success",
	ERR_DB_ERR:                "Database Error",
	ERR_NOT_ENOUGH_PARAS:      "No Enough Paras",
	ERR_TOO_MANY_PARAS:        "Too Many Paras",
	ERR_UNACCP_PARAS:          "Unaccept Paras",
	ERR_CMD_ERR:               "Command Error",
	ERR_COMMON_ERR:            "Common Error",
	ERR_SEGMENT_NOT_EXIST:     "Segment Not Exist",
	ERR_SEGMENT_ALREADY_EXIST: "Segment Already Exist",
	ERR_TIMEOUT:               "Timeout Error",
	ERR_SYSCALL_ERR:           "System Call Error",
	ERR_SYSTEM_ERR:            "System Error",
	ERR_NO_SUCH_API:           "No Such API",
	ERR_NOT_IMPLEMENTED:       "Function not Implemented",

	// User
	ERR_USER_NOT_EXIST:        "User Not Exist",
	ERR_USER_ALREADY_EXIST:    "User Already Exist",
	ERR_PASSWORD_DONT_MATCH:   "User And Password Not Match",
	ERR_USER_NOT_LOGIN:        "User Not Login",
	ERR_USER_GROUPS_NOT_EMPTY: "Groups under Account must be empty",

	// User group
	ERR_USERGROUP_NOT_EXIST:       "User Group Not Exist",
	ERR_USERGROUP_ALREADY_EXIST:   "User Group Already Exist",
	ERR_USERGROUP_USERS_NOT_EMPTY: "Users under Group must be empty",
}

var GErrorsCN = map[int]string{
	ERR_OCT_SUCCESS:           "操作成功",
	ERR_DB_ERR:                "数据库错误",
	ERR_NOT_ENOUGH_PARAS:      "参数不足",
	ERR_TOO_MANY_PARAS:        "太多参数",
	ERR_UNACCP_PARAS:          "参数不合法",
	ERR_CMD_ERR:               "命令执行错误",
	ERR_COMMON_ERR:            "通用错误",
	ERR_SEGMENT_NOT_EXIST:     "对象不存在",
	ERR_SEGMENT_ALREADY_EXIST: "对象已存在",
	ERR_TIMEOUT:               "超时错误",
	ERR_SYSCALL_ERR:           "系统调用错误",
	ERR_SYSTEM_ERR:            "系统错误",
	ERR_NO_SUCH_API:           "无此API",
	ERR_NOT_IMPLEMENTED:       "功能未实现",

	// User
	ERR_USER_NOT_EXIST:        "用户不存在",
	ERR_USER_ALREADY_EXIST:    "用户已经存在",
	ERR_PASSWORD_DONT_MATCH:   "用户和密码不匹配",
	ERR_USER_NOT_LOGIN:        "用户未登录",
	ERR_USER_GROUPS_NOT_EMPTY: "该账号下的用户组不为空",

	// User group
	ERR_USERGROUP_NOT_EXIST:       "用户组不存在",
	ERR_USERGROUP_ALREADY_EXIST:   "用户组已经存在",
	ERR_USERGROUP_USERS_NOT_EMPTY: "该用户组下的用户不为空",
}

type MirageError struct {
	ErrorNo  int    `json:"no"`
	ErrorMsg string `json:"msg"`
}

func NewError(code int, message string) *MirageError {
	return &MirageError{
		ErrorNo:  code,
		ErrorMsg: message,
	}
}

func GetMsg(errorNo int) string {
	return GErrors[errorNo]
}

func GetMsgCN(errorNo int) string {
	return GErrorsCN[errorNo]
}
