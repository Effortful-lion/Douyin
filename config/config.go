// config.go
// 将配置文件加载到内存中，方便使用
package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

var Config = &AppConfig{}

// app应用配置
type AppConfig struct {
	MysqlConfig    MysqlConfig    `mapstructure:"MYSQL"`
	RedisConfig    RedisConfig    `mapstructure:"Redis"` // 改为 Redis
	GinConfig      GinConfig      `mapstructure:"Gin"`   // 改为 Gin
	JaegerConfig   JaegerConfig   `mapstructure:"Jaeger"`
	EtcdConfig     EtcdConfig     `mapstructure:"Etcd"`
	ServiceConfig  ServiceConfig  `mapstructure:"Service"`
	RabbitMQConfig RabbitMQConfig `mapstructure:"RabbitMQ"`
	AliyunConfig   AliyunConfig   `mapstructure:"AliYun"` // 改为 AliYun
	Domain         Domain         `mapstructure:"Domain"`
	User           User           `mapstructure:"User"`
}

// mysql配置
type MysqlConfig struct {
	DBHost     string `mapstructure:"DBHost"`
	DBPort     int    `mapstructure:"DBPort"` // 类型改为 int，匹配 YAML
	DBUser     string `mapstructure:"DBUser"`
	DBPassword string `mapstructure:"DBPassWord"` // 改为 DBPassWord
	DBName     string `mapstructure:"DBName"`
	DBCharset  string `mapstructure:"Charset"` // 改为 Charset
}

// redis配置
type RedisConfig struct {
	RedisHost     string `mapstructure:"RedisHost"` // 改为 RedisHost
	RedisPort     int    `mapstructure:"RedisPort"` // 改为 RedisPort，类型 int
	RedisPassword string `mapstructure:"RedisPassword"`
}

// gin配置
type GinConfig struct {
	AppMode string `mapstructure:"AppMode"`
	AppHost string `mapstructure:"HttpHost"` // 改为 HttpHost
	AppPort string `mapstructure:"HttpPort"` // 改为 HttpPort
}

// etcd配置
type EtcdConfig struct {
	EtcdHost string `mapstructure:"EtcdHost"`
	EtcdPort int    `mapstructure:"EtcdPort"` // 类型 int
}

// jaeger配置
type JaegerConfig struct {
	JaegerHost string `mapstructure:"JaegerHost"`
	JaegerPort int    `mapstructure:"JaegerPort"` // 类型 int
}

// rabbitmq配置
type RabbitMQConfig struct {
	RabbitMQ         string `mapstructure:"RabbitMQ"`
	RabbitMQHost     string `mapstructure:"RabbitMQHost"`
	RabbitMQPort     int    `mapstructure:"RabbitMQPort"` // 类型 int
	RabbitMQUser     string `mapstructure:"RabbitMQUser"`
	RabbitMQPassword string `mapstructure:"RabbitMQPassWord"` // 改为 RabbitMQPassWord
}

// aliyun配置
type AliyunConfig struct {
	Bucket    string `mapstructure:"Bucket"` // 改为 Bucket
	AccessKey string `mapstructure:"AccessKey"`
	SecretKey string `mapstructure:"SecretKey"`
	Domain    string `mapstructure:"Domain"`
}

// 服务配置(记录多个服务节点)
type ServiceConfig struct {
	TestServiceAddress     string `mapstructure:"TestServiceAddress"`
	UserServiceAddress     string `mapstructure:"UserServiceAddress"`
	VideoServiceAddress    string `mapstructure:"VideoServiceAddress"`
	FavoriteServiceAddress string `mapstructure:"FavoriteServiceAddress"`
	CommentServiceAddress  string `mapstructure:"CommentServiceAddress"`
	RelationServiceAddress string `mapstructure:"RelationServiceAddress"`
	MessageServiceAddress  string `mapstructure:"MessageServiceAddress"`
}

type Domain struct {
	TestServiceDomain     string `mapstructure:"TestServiceDomain"`
	UserServiceDomain     string `mapstructure:"UserServiceDomain"`
	VideoServiceDomain    string `mapstructure:"VideoServiceDomain"`
	FavoriteServiceDomain string `mapstructure:"FavoriteServiceDomain"`
	CommentServiceDomain  string `mapstructure:"CommentServiceDomain"`
	RelationServiceDomain string `mapstructure:"RelationServiceDomain"`
	MessageServiceDomain  string `mapstructure:"MessageServiceDomain"`
}

type User struct {
	Avatar     string `mapstructure:"Avatar"`
	Background string `mapstructure:"Background"`
	Signature  string `mapstructure:"Signature"`
}

func InitConfig() {
	// 使用 filepath.Join 自动处理路径分隔符
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		panic(err)
	}
	// 监控配置文件，当发生变化时，重新加载配置文件
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(&Config); err != nil {
			panic(err)
		}
		log.Println("配置文件已更新")
	})
	log.Println("配置文件加载成功")
}

// func InitConfig() (err error) {
// 	// 从etcd获取配置文件
// 	if err = discovery.GetConfig("config.yaml"); err != nil {}
// 	// 读取配置文件

// 	// 将配置文件中的内容映射到结构体中

// 	return nil
// }
