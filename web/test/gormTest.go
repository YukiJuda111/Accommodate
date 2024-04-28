package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Student struct {
	gorm.Model
	Name string
}

var Globaldb *gorm.DB

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	Globaldb = db
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := Globaldb.DB()

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		fmt.Println("连接mysql失败: ", err)
		return
	}

	err = Globaldb.AutoMigrate(&Student{})
	if err != nil {
		return
	}

	InsertData()
	DeleteData()

}

func InsertData() {
	Globaldb.Create(&Student{Name: "stu1"})
	Globaldb.Create(&Student{Name: "stu2"})
}

func DeleteData() {
	Globaldb.Where("name = ?", "stu1").Delete(&Student{})
}
