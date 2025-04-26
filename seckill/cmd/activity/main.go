package main

import(
	"log"

	"github.com/spf13/viper"

	"Redrock/seckill/internal/activity/config"
	"Redrock/seckill/internal/activity/service"
	"Redrock/seckill/internal/pkg/database"
	"Redrock/seckill/internal/pkg/redis"
	activity "Redrock/seckill/kitex_gen/activity/internalactivityservice"
)

func main(){

	// 读取配置
	viper.SetConfigName("activity")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/activity/config")
	viper.AutomaticEnv() // 自动读取环境变量，方便后续调试

	err := viper.ReadInConfig()
	if err != nil{
		log.Fatalf("读取活动配置文件失败：%v", err)
	}

	var config config.Config
	err = viper.Unmarshal(&config)
	if err != nil{
		log.Fatalf("解析活动配置文件失败：%v", err)
	}

	// 连接数据库
	if err := database.InitDB(&config.Database); err != nil{
		log.Fatalf("初始化连接数据库失败：%v", err)
	}
	defer database.CloseDB()

	// 连接Redis
	if err := redis.InitRedis(&config.Redis); err != nil{
		log.Fatalf("初始化连接Redis失败：%v", err)
	}
	defer redis.CloseRedis()

	// TODO:启动kitex服务
	svr := activity.NewServer(new(service.InternalActivityServiceImpl))

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}

	select{}
}
