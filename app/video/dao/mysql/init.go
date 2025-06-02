package mysql

import (
	"Douyin/config"
	"Douyin/model"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

var db *gorm.DB

func InitMysql() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Config.MysqlConfig.DBUser,
		config.Config.MysqlConfig.DBPassword,
		config.Config.MysqlConfig.DBHost,
		config.Config.MysqlConfig.DBPort,
		config.Config.MysqlConfig.DBName)
	// 连接数据库
	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}
	DB, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 微秒 精度，MySQL 5.6 之前的数据库不支持（采用秒级精度）
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	db = DB
	log.Println("数据库连接成功")
	err = migration()
	if err != nil {
		panic(err)
	}
}

// 创建一个有 context 的 DB对象，可以控制请求的生命周期
func NewDBClient(ctx context.Context) *gorm.DB {
	DB := db
	return DB.WithContext(ctx)
}

// 自动迁移表结构
func migration() error {
	err := db.Set(`gorm:table_options`, "charset=utf8mb4").
		AutoMigrate(&model.User{}, &model.Follow{}, &model.Video{}, &model.Favorite{}, &model.Comment{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("表结构迁移成功")
	return err
}
