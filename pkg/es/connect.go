package es

import (
	"yujian-backend/pkg/config"
	"yujian-backend/pkg/log"

	"github.com/elastic/go-elasticsearch/v8"
)

var esClient *elasticsearch.Client

func InitESClient() {
	esConfig := config.Config.ES
	cfg := elasticsearch.Config{
		Addresses: esConfig.Addresses,
		Username:  esConfig.Username,
		Password:  esConfig.Password,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.GetLogger().Error("创建ES客户端失败: %v", err)
		return
	} else {
		log.GetLogger().Info("创建ES客户端成功")
		esClient = client
	}
}
