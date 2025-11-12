package response

// 业务状态码定义
const (
	// 通用状态码
	CodeSuccess         = 200   // 成功
	CodeBadRequest      = 400   // 错误的请求
	CodeUnauthorized    = 401   // 未授权
	CodeForbidden       = 403   // 禁止访问
	CodeNotFound        = 404   // 未找到
	CodeTooManyRequests = 429   // 请求过多（限流）
	CodeInternalError   = 500   // 内部错误
	CodeValidateError   = 10001 // 参数验证错误

	// 用户相关 (20xxx)
	CodeUserNotFound       = 20001 // 用户不存在
	CodeUserAlreadyExists  = 20002 // 用户已存在
	CodeUserPasswordError  = 20003 // 密码错误
	CodeUserDisabled       = 20004 // 用户已禁用
	CodeUserTokenInvalid   = 20005 // Token 无效
	CodeUserTokenExpired   = 20006 // Token 过期
	CodeUserPermissionDeny = 20007 // 权限不足

	// 数据库相关 (30xxx)
	CodeDatabaseError  = 30001 // 数据库错误
	CodeRecordNotFound = 30002 // 记录不存在
	CodeRecordExists   = 30003 // 记录已存在

	// 缓存相关 (40xxx)
	CodeCacheError = 40001 // 缓存错误

	// 消息队列相关 (50xxx)
	CodeKafkaError = 50001 // Kafka 错误

	// 第三方服务相关 (60xxx)
	CodeThirdPartyError = 60001 // 第三方服务错误
)

// CodeMessage 状态码对应的默认消息
var CodeMessage = map[int]string{
	CodeSuccess:         "success",
	CodeBadRequest:      "bad request",
	CodeUnauthorized:    "unauthorized",
	CodeForbidden:       "forbidden",
	CodeNotFound:        "not found",
	CodeTooManyRequests: "too many requests",
	CodeInternalError:   "internal server error",
	CodeValidateError:   "validate error",

	CodeUserNotFound:       "user not found",
	CodeUserAlreadyExists:  "user already exists",
	CodeUserPasswordError:  "password error",
	CodeUserDisabled:       "user disabled",
	CodeUserTokenInvalid:   "token invalid",
	CodeUserTokenExpired:   "token expired",
	CodeUserPermissionDeny: "permission denied",

	CodeDatabaseError:  "database error",
	CodeRecordNotFound: "record not found",
	CodeRecordExists:   "record already exists",

	CodeCacheError: "cache error",
	CodeKafkaError: "kafka error",

	CodeThirdPartyError: "third party service error",
}

// GetMessage 获取状态码对应的消息
func GetMessage(code int) string {
	if msg, ok := CodeMessage[code]; ok {
		return msg
	}
	return "unknown error"
}
