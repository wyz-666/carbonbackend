package handlers

import (
	"carbonbackend/app/handlers/request"
	"carbonbackend/app/handlers/response"
	"carbonbackend/db"
	"carbonbackend/db/model"
	"carbonbackend/service"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func MarketSubmit(c *gin.Context) {
	log.Println("################## CEA Market Submit ##################")
	var reqMarket request.ReqMarket
	if err := c.ShouldBind(&reqMarket); err != nil {
		log.Printf("[ERROR] : %v", err)
		response.MakeFail(c, http.StatusNotAcceptable, err.Error())
		return
	}

	err := service.MarketSubmit(&reqMarket, reqMarket.Product)
	if err != nil {
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("market submit successful")
	response.MakeSuccess(c, http.StatusOK, "successfully market submit!")
	return
}

// 预测结果提交
func StatsSubmit(c *gin.Context) {
	log.Println("################## Stats Submit ##################")
	var reqExpectation request.ReqExpectation
	if err := c.ShouldBind(&reqExpectation); err != nil {
		log.Printf("[ERROR] : %v", err)
		response.MakeFail(c, http.StatusNotAcceptable, err.Error())
		return
	}
	err := service.StatsSubmit(&reqExpectation)
	if err != nil {
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("stats submit successful")
	response.MakeSuccess(c, http.StatusOK, "successfully stats submit!")
	return
}

func GECStatsSubmit(c *gin.Context) {
	log.Println("################## GEC Stats Submit ##################")
	var reqExpectation request.ReqGECExpectation
	if err := c.ShouldBind(&reqExpectation); err != nil {
		log.Printf("[ERROR] : %v", err)
		response.MakeFail(c, http.StatusNotAcceptable, err.Error())
		return
	}
	err := service.GECStatsSubmit(&reqExpectation)
	if err != nil {
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("GEC stats submit successful")
	response.MakeSuccess(c, http.StatusOK, "successfully GEC stats submit!")
	return
}

// func CCERMarketSubmit(c *gin.Context) {
// 	log.Println("################## CCER Market Submit ##################")
// 	var reqMarket request.ReqMarket
// 	if err := c.ShouldBind(&reqMarket); err != nil {
// 		log.Printf("[ERROR] : %v", err)
// 		response.MakeFail(c, http.StatusNotAcceptable, err.Error())
// 		return
// 	}
// 	err := service.MarketSubmit(&reqMarket, "CCER")
// 	if err != nil {
// 		response.MakeFail(c, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	log.Println("CCER market submit successful")
// 	response.MakeSuccess(c, http.StatusOK, "successfully CCER market submit!")
// 	return
// }

func GetCCERMarket(c *gin.Context) {
	var ms []model.CCERMarket
	cli := db.Get()
	err := cli.Find(&ms).Error
	if err != nil {
		log.Printf("[ERROR] get all ccer market error!")
		response.MakeFail(c, http.StatusBadRequest, "get all ccer market error")
		return
	}
	log.Println("get ccer market successful")
	response.MakeSuccess(c, http.StatusOK, ms)
	return
}

func GetCEAMarket(c *gin.Context) {
	var ms []model.CEAMarket
	cli := db.Get()
	err := cli.Find(&ms).Error
	if err != nil {
		log.Printf("[ERROR] get all ccer market error!")
		response.MakeFail(c, http.StatusBadRequest, "get all ccer market error")
		return
	}
	log.Println("get ccer market successful")
	response.MakeSuccess(c, http.StatusOK, ms)
	return
}

// 得分排序
func GetCEAMonthScore(c *gin.Context) {
	log.Println("################## Get CEA Month Score ##################")
	nowTimeStr := c.Query("nowTime")
	t, err := time.Parse("2006-01-02 15:04:05", nowTimeStr)
	if err != nil {
		c.JSON(400, gin.H{"msg": "nowTime 格式错误", "error": err.Error()})
		return
	}
	result, err := service.GetCEAMonthScoreList(t)
	if err != nil {
		log.Printf("Get CEA CEA Month Score error!")
		response.MakeFail(c, http.StatusBadRequest, "Get CEA Month Score error")
		return
	}
	log.Println("Get CEA Month Score successfully")
	response.MakeSuccess(c, http.StatusOK, result)
	return
}

func GetCCERMonthScore(c *gin.Context) {
	log.Println("################## Get CCER Month Score ##################")
	nowTimeStr := c.Query("nowTime")
	t, err := time.Parse("2006-01-02 15:04:05", nowTimeStr)
	if err != nil {
		c.JSON(400, gin.H{"msg": "nowTime 格式错误", "error": err.Error()})
		return
	}
	result, err := service.GetCCERMonthScoreList(t)
	if err != nil {
		log.Printf("Get CCER Month Score error!")
		response.MakeFail(c, http.StatusBadRequest, "Get CCER Month Score error")
		return
	}
	log.Println("Get CCER Month Score successfully")
	response.MakeSuccess(c, http.StatusOK, result)
	return
}
