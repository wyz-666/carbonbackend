package dao

import (
	"carbonbackend/db"
	"carbonbackend/db/model"
)

const tableName = "Counters"

// ClearCounter 清除Counter
func (imp *CounterInterfaceImp) ClearCounter(id int32) error {
	cli := db.Get()
	return cli.Table(tableName).Delete(&model.CounterModel{Id: id}).Error
}

// UpsertCounter 更新/写入counter
func (imp *CounterInterfaceImp) UpsertCounter(counter *model.CounterModel) error {
	cli := db.Get()
	return cli.Table(tableName).Save(counter).Error
}

// GetCounter 查询Counter
func (imp *CounterInterfaceImp) GetCounter(id int32) (*model.CounterModel, error) {
	var err error
	var counter = new(model.CounterModel)

	cli := db.Get()
	err = cli.Table(tableName).Where("id = ?", id).First(counter).Error

	return counter, err
}

// GetUser implements UserInfoInterface.
func (userInfoImp *UserInfoInterfaceImp) GetUser() (*model.UserInfoModel, error) {
	var err error
	var user = new(model.UserInfoModel)

	cli := db.Get()
	err = cli.Table("UserInfo").First(user).Error

	return user, err
}
