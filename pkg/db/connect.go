package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"yujian-backend/pkg/config"

	"yujian-backend/pkg/log"
	"yujian-backend/pkg/model"
)

func InitDB() {
	db := createConnect(config.Config.DB)
	if err := db.AutoMigrate(&model.UserDO{}, &model.PostDO{}, &model.PostCommentDO{}, &model.BookInfoDO{}, &model.BookCommentDO{}, &model.UserRecommendRecordDO{}); err != nil {
		log.GetLogger().Fatalf("failed to migrate database: %s", err)
	} else {
		log.GetLogger().Info("Successfully migrated database...")
	}

	userRepository = UserRepository{DB: db}
	postRepository = PostRepository{DB: db}
	bookRepository = BookRepository{DB: db}
	recommendRepository = RecommendRepository{DB: db}
}

func createConnect(config *model.DBConfig) *gorm.DB {
	logger := log.GetLogger()
	db, err := gorm.Open(mysql.Open(config.CreateDsn()), &gorm.Config{})
	if err != nil {
		logger.Fatalf("failed to connect database: %s", err)
		return nil
	} else {
		logger.Info("Connected to database")
		return db
	}
}
