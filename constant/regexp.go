package constant

const (
	RegexpEmail   = `(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`
	RegexpAccount = `^[a-zA-z0-9_]{2,50}$` // 允许用户名中出现点号
)
