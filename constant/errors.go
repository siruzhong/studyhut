package constant

import "errors"

// 错误类型
var (
	ErrMemberNoExist             = errors.New("用户不存在")
	ErrMemberExist               = errors.New("用户已存在")
	ErrMemberDisabled            = errors.New("用户被禁用")
	ErrMemberEmailEmpty          = errors.New("用户邮箱不能为空")
	ErrMemberEmailExist          = errors.New("用户邮箱已被使用")
	ErrMemberDescriptionTooLong  = errors.New("用户描述必须小于500字")
	ErrMemberEmailFormatError    = errors.New("邮箱格式不正确")
	ErrMemberPasswordFormatError = errors.New("密码必须在6-50个字符之间")
	ErrMemberAccountFormatError  = errors.New("账号只能由英文字母数字组成，且在3-50个字符")
	ErrMemberRoleError           = errors.New("用户权限不正确")
	ErrorMemberPasswordError     = errors.New("用户密码错误")
	ErrMemberAuthMethodInvalid   = errors.New("不支持此认证方式")
	ErrLDAPConnect               = errors.New("无法连接到LDAP服务器")
	ErrLDAPFirstBind             = errors.New("第一次LDAP绑定失败")
	ErrLDAPSearch                = errors.New("LDAP搜索失败")
	ErrLDAPUserNotFoundOrTooMany = errors.New("LDAP用户不存在或者多于一个")
	ErrDataNotExist              = errors.New("数据不存在")
	ErrInvalidParameter          = errors.New("无效参数")
	ErrPermissionDenied          = errors.New("拒绝访问")
	ErrCommentClosed             = errors.New("评论已关闭")
	ErrCommentContentNotEmpty    = errors.New("评论内容不能为空")
)
