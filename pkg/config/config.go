package config

import (
	"github.com/spf13/viper"
	"log" // use default log before we init logger
	"runtime/debug"

	"yujian-backend/pkg/model"
)

var Config model.AppConfig

// initDBConfig 初始化数据库配置。
func initDBConfig() {
	dbConfig := Config.DB
	dbConfig.Host = viper.GetString("db.host")
	dbConfig.Port = viper.GetInt("db.port")
	dbConfig.UserName = viper.GetString("db.username")
	dbConfig.PassWord = viper.GetString("db.password")
	dbConfig.DBName = viper.GetString("db.dbname")
	dbConfig.Params = viper.GetString("db.params")
}

func initLogConfig() {
	logConfig := Config.Log
	logConfig.FileName = viper.GetString("log.filename")
	logConfig.ErrFileName = viper.GetString("log.errFilename")
	logConfig.LogLevel = viper.GetString("log.loglevel")
}

func initServerConfig() {
	serverConfig := Config.Server
	serverConfig.Port = viper.GetInt("server.port")
}

func initESConfig() {
	esConfig := Config.ES
	esConfig.Addresses = viper.GetStringSlice("es.addresses")
	esConfig.Username = viper.GetString("es.username")
	esConfig.Password = viper.GetString("es.password")
}

func InitConfig() {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			log.Fatalf("Init Config failed. Recovered from panic: %v\n", r)
		}
	}()

	Config = model.AppConfig{
		DB:     &model.DBConfig{},
		ES:     &model.ESConfig{},
		Log:    &model.LogConfig{},
		Server: &model.ServerConfig{},
	}

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

	initESConfig()

	initLogConfig()

	initServerConfig()

	for _, v := range viper.AllKeys() {
		log.Printf("%s = %v\n", v, viper.Get(v))
	}
}
