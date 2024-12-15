package config

import (
	"log" // use default log before we init logger

	"github.com/spf13/viper"

	"yujian-backend/pkg/model"
)

var Config model.AppConfig

// initDBConfig 初始化数据库配置。
func initDBConfig() {
	dbConfig := Config.DB
	dbConfig.Host = viper.GetString("db.host")
	dbConfig.UserName = viper.GetString("db.username")
	dbConfig.PassWord = viper.GetString("db.password")
	dbConfig.DBName = viper.GetString("db.dbname")
	dbConfig.Charset = viper.GetString("db.charset")
	dbConfig.TimeZone = viper.GetString("db.timezone")
}

func initLogConfig() {
	logConfig := Config.Log
	logConfig.FileName = viper.GetString("log.filename")
	logConfig.LogLevel = viper.GetString("log.loglevel")
}

func initServerConfig() {
	serverConfig := Config.Server
	serverConfig.Port = viper.GetString("server.port")
}

func initESConfig() {
	esConfig := Config.ES
	esConfig.Addresses = viper.GetStringSlice("es.addresses")
	esConfig.Username = viper.GetString("es.username")
	esConfig.Password = viper.GetString("es.password")
}

func InitConfig() {
	// 初始化 viper
	viper.SetConfigName("config")  // 配置文件名称（不带扩展名）
	viper.SetConfigType("yaml")    // 配置文件类型
	viper.AddConfigPath("config/") // 配置文件路径

	// 允许使用环境变量
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	initDBConfig()

	initLogConfig()

	initServerConfig()
}
