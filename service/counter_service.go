package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// JsonResult 返回结构
type JsonResult struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// IndexHandler 计数器接口
func IndexHandler(c *gin.Context) {
	data, err := getIndex()
	if err != nil {
		c.String(http.StatusInternalServerError, "内部错误")
		return
	}
	c.String(http.StatusOK, data)
}

// GET: /api/user
func GetUserInfo(c *gin.Context) {
	res := &JsonResult{}
	counter, err := getFirstUserInfo()
	if err != nil {
		res.Code = -1
		res.ErrorMsg = err.Error()
	} else {
		res.Data = counter
	}
	c.JSON(http.StatusOK, res)
}

// GET: /api/count
func GetCounterHandler(c *gin.Context) {
	res := &JsonResult{}
	counter, err := getCurrentCounter()
	if err != nil {
		res.Code = -1
		res.ErrorMsg = err.Error()
	} else {
		res.Data = counter.Count
	}
	c.JSON(http.StatusOK, res)
}

// POST: /api/count
func PostCounterHandler(c *gin.Context) {
	res := &JsonResult{}
	count, err := modifyCounter(c)
	if err != nil {
		res.Code = -1
		res.ErrorMsg = err.Error()
	} else {
		res.Data = count
	}
	c.JSON(http.StatusOK, res)
}

// modifyCounter 更新计数，自增或者清零
func modifyCounter(c *gin.Context) (int32, error) {
	action, err := getAction(c)
	if err != nil {
		return 0, err
	}

	var count int32
	if action == "inc" {
		count, err = upsertCounter()
	} else if action == "clear" {
		err = clearCounter()
		count = 0
	} else {
		err = fmt.Errorf("参数 action : %s 错误", action)
	}
	return count, err
}

// 获取action参数
func getAction(c *gin.Context) (string, error) {
	var body struct {
		Action string `json:"action"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return "", fmt.Errorf("解析action失败: %v", err)
	}
	if body.Action == "" {
		return "", fmt.Errorf("缺少 action 参数")
	}
	return body.Action, nil
}

// upsertCounter 更新或插入计数器
func upsertCounter() (int32, error) {
	currentCounter, err := getCurrentCounter()
	var count int32
	createdAt := time.Now()

	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	} else if err == gorm.ErrRecordNotFound {
		count = 1
	} else {
		count = currentCounter.Count + 1
		createdAt = currentCounter.CreatedAt
	}

	counter := &model.CounterModel{
		Id:        1,
		Count:     count,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}

	err = dao.Imp.UpsertCounter(counter)
	return counter.Count, err
}

func clearCounter() error {
	return dao.Imp.ClearCounter(1)
}

// getCurrentCounter 查询当前计数器
func getCurrentCounter() (*model.CounterModel, error) {
	counter, err := dao.Imp.GetCounter(1)
	if err != nil {
		return nil, err
	}

	return counter, nil
}

func getFirstUserInfo() (*model.UserInfoModel, error) {
	userInfo, err := dao.UserInfoImp.GetUser()

	if err != nil {
		return nil, err
	}

	return userInfo, nil

}

// getIndex 获取主页
func getIndex() (string, error) {
	b, err := ioutil.ReadFile("./index.html")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
