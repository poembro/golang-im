package db

import (
    "fmt"
    "golang-im/config"
    "golang-im/pkg/logger"

    "github.com/go-redis/redis"
    //"github.com/jinzhu/gorm"

    _ "github.com/go-sql-driver/mysql"
)

var (
    //DB       *gorm.DB
    RedisCli *redis.Client
)

// InitMysql 初始化MySQL
func InitMysql(dataSource string) {
    /*
        logger.Logger.Info("init mysql")
        var err error
        DB, err = gorm.Open("mysql", dataSource)
        if err != nil {
            panic(err)
        }
        DB.SingularTable(true)
        DB.LogMode(true)
        logger.Logger.Info("init mysql ok")
    */
}

/*
// InitMysql 初始化MySQL分布式节点
func InitMysqlMulti() {
    logger.Logger.Info("init mysql Multi")

    conf := map[string]string{
        "db_user1" : "root:123456@tcp(192.168.3.111:3306)/gim?charset=utf8&parseTime=true",
        "db_user2" : "root:123456@tcp(192.168.3.111:3306)/gim?charset=utf8&parseTime=true",
        "db_user3" : "root:123456@tcp(192.168.3.111:3306)/gim?charset=utf8&parseTime=true",
    }
    // 定义个全局变量 conns
    conns := make(map[string]*gorm.DB)

    for k, v := range conf {
        DB, err := gorm.Open("mysql", v)
        if err != nil {
            panic(err)
        }
        DB.SingularTable(true)
        DB.LogMode(true)
        conns[k] = DB
    }

    logger.Logger.Info("init mysql Multi ok")
}
*/

// InitRedis 初始化Redis
func InitRedis(addr, password string) {
    logger.Logger.Info("init redis")
    RedisCli = redis.NewClient(&redis.Options{
        Addr:     addr,
        DB:       0,
        Password: password,
    })

    _, err := RedisCli.Ping().Result()
    if err != nil {
        panic(err)
    }

    logger.Logger.Info("init redis ok")
}

// InitByTest 初始化数据库配置，仅用在单元测试
func InitByTest() {
    fmt.Println("init db")
    logger.Target = logger.Console
    logger.Init()

    InitMysql(config.Logic.MySQL)
    InitRedis(config.Global.RedisIP, config.Global.RedisPassword)
}
