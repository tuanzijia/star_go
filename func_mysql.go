package starGo

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Mysql struct {
	db            *gorm.DB
	connectionStr string
}

func NewMysql(connection string) *Mysql {
	mysql := new(Mysql)
	db, err := gorm.Open("mysql", connection)
	if err != nil {
		ErrorLog("连接mysql出错,错误信息:%v", err)
		panic(fmt.Errorf("连接mysql出错,错误信息:%v", err))
	}

	mysql.db = db
	mysql.connectionStr = connection
	mysqlCfg = mysql

	return mysql
}

func (m *Mysql) GetDb() *gorm.DB {
	return m.db
}

func (m *Mysql) GetConnectionStr() string {
	return m.connectionStr
}

func (m *Mysql) RegisterTableModel(model interface{}) error {
	return m.db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(model).Error

}

func (m *Mysql) RegisterTableModelForTableName(tableName string, model interface{}) error {
	return m.db.Table(tableName).AutoMigrate(model).Error
}
