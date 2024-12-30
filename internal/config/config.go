package config

import (
	"im-system/internal/model/db"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Logger Global logger
var Logger *logrus.Logger

// RedisClient Redis client
var RedisClient *redis.Client

// JWTSecret JWT密钥
var JWTSecret string

// Config 配置结构体
type Config struct {
	Server struct {
		HTTPPort string `yaml:"http_port"`
		GRPCPort string `yaml:"grpc_port"`
	} `yaml:"server"`
	MySQL struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"mysql"`
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
	JWTSecret string `yaml:"jwt_secret"` // JWT 密钥
}

// LoadConfig 加载配置文件
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	// 全局赋值
	JWTSecret = config.JWTSecret
	return &config, nil
}

// LoadDatabaseConfig 加载数据库配置
func LoadDatabaseConfig(cfg *Config) (string, error) {
	dsn := cfg.MySQL.Username + ":" + cfg.MySQL.Password + "@tcp(" + cfg.MySQL.Host + ":" + cfg.MySQL.Port + ")/" + cfg.MySQL.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	return dsn, nil
}

// InitDB 初始化数据库连接
func InitDB(cfg *Config) error {
	dsn, err := LoadDatabaseConfig(cfg)
	if err != nil {
		return err
	}

	dbInstance, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	// 全局赋值
	db.DB = dbInstance
	Logger.Info("成功连接到MySQL数据库")
	return nil
}

// InitLogger 初始化日志配置
func InitLogger() {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetLevel(logrus.DebugLevel)

	// 获取当前日期并格式化为字符串
	currentDate := time.Now().Format("2006-01-02")
	logFileName := "log/system_" + currentDate + ".log"

	// 确保日志目录存在
	if err := os.MkdirAll("log", os.ModePerm); err != nil {
		Logger.Fatalf("无法创建日志目录: %v", err)
	}

	// 创建日志文件
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.Fatalf("无法打开日志文件: %v", err)
	}
	Logger.SetOutput(logFile)
}

// InitRedis 初始化 Redis 客户端
func InitRedis(cfg *Config) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// todo: 测试连接
	if err := RedisClient.Ping(RedisClient.Context()).Err(); err != nil {
		Logger.Fatalf("无法连接到Redis: %v", err)
	}
}
