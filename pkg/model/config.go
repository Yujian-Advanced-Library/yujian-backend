package model

import "strconv"

type DBConfig struct {
	UserName string
	PassWord string
	Host     string
	Port     int
	DBName   string
	Params   string
}

// CreateDsn 生成数据库连接字符串。
func (config *DBConfig) CreateDsn() string {
	dsn := config.UserName + ":" + config.PassWord + "@tcp(" + config.Host + ":" + strconv.Itoa(config.Port) + ")/" + config.DBName
	if len(config.Params) > 0 {
		dsn += "?" + config.Params
	}
	return dsn
}

type LogConfig struct {
	FileName    string // 日志文件名
	ErrFileName string // 日志文件名
	LogLevel    string // 日志级别
}

type ServerConfig struct {
	Port int
}

type ESConfig struct {
	Addresses []string
	Username  string
	Password  string
}

type AppConfig struct {
	DB     *DBConfig
	Log    *LogConfig
	Server *ServerConfig
	ES     *ESConfig
}
