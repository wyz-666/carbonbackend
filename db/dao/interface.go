package dao

import (
	"carbonbackend/db/model"
)

// CounterInterface 计数器数据模型接口
type CounterInterface interface {
	GetCounter(id int32) (*model.CounterModel, error)
	UpsertCounter(counter *model.CounterModel) error
	ClearCounter(id int32) error
}

// CounterInterfaceImp 计数器数据模型实现
type CounterInterfaceImp struct{}

// Imp 实现实例
var Imp CounterInterface = &CounterInterfaceImp{}

type UserInfoInterface interface {
	GetUser() (*model.UserInfoModel, error)
}

// CounterInterfaceImp 计数器数据模型实现
type UserInfoInterfaceImp struct{}

// Imp 实现实例
var UserInfoImp UserInfoInterface = &UserInfoInterfaceImp{}
