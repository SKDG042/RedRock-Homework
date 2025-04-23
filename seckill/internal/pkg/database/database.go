package database

import(
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	mysqldriver "github.com/go-sql-driver/mysql"
)

var DB *gorm.DB

func isDatabaseNotExistsError(err error) bool{
	// 类型断言，将err转换为*mysqldriver.MySQLError类型
	if mysqlErr, ok := err.(*mysqldriver.MySQLError); ok{
		return mysqlErr.Number == 1049
	}
	return false
}

func connectDB(config *DatabaseConfig) (*gorm.DB, error){
	db,err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.Charset,
		config.ParseTime,
		config.Loc,
	)), &gorm.Config{})

	if err != nil{
		return nil, fmt.Errorf("连接数据库失败：%w", err)
	}

	return db, nil
}

func createDB(config *DatabaseConfig) error{
	db,err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
	)), &gorm.Config{})
	
	if err != nil{
		return fmt.Errorf("连接数据库失败：%w", err)
	}

	err =db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.DBName)).Error
	if err != nil{
		return fmt.Errorf("创建数据库失败：%w", err)
	}

	sqlDB,err := db.DB()
	if err != nil{
		return fmt.Errorf("获取原始数据库连接失败：%w", err)
	}
	
	err = sqlDB.Close()
	if err != nil{
		return fmt.Errorf("关闭原始数据库连接失败：%w", err)
	}

	return nil
}

// 当数据库不存在时会先创建数据库
func InitDB(config *DatabaseConfig) error{
	var err error
	
	DB,err = connectDB(config)
	if err != nil{
		if isDatabaseNotExistsError(err){
			err = createDB(config)
			if err != nil{
				return fmt.Errorf("数据库不存在且创建失败：%w", err)
			}
			DB,err = connectDB(config)
			if err != nil{
				return fmt.Errorf("数据库连接失败：%w", err)
			}
		}else{
			return fmt.Errorf("数据库存在但连接失败：%w", err)
		}
	}

	return nil
}

func GetDB() *gorm.DB{
	if DB == nil{
		panic("你还没有初始化数据库")
	}

	return DB
}

func CloseDB(){
	DB = GetDB()

	if DB !=nil{
		sqlDB,err := DB.DB()
		if err != nil{
			log.Printf("获取数据库连接失败：%v", err)
		}
		if err := sqlDB.Close(); err != nil{
			log.Printf("关闭数据库连接失败：%v", err)
		}else{
			log.Println("数据库连接成功关闭")
		}
	}
}

// 可变参数... 表示接受不定数量的参数,
func MigrateDB(models ...any) error{
	if DB == nil{
		return fmt.Errorf("数据库未初始化")
	}

	err := DB.AutoMigrate(models...)
	if err != nil{
		return fmt.Errorf("自动迁移数据库失败：%w", err)
	}

	log.Println("数据库迁移成功")

	return nil
}
