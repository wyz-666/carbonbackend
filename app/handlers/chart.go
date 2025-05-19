package handlers

import (
	"log"
	"net/http"
	"time"
	"wxcloudrun-golang/app/handlers/response"
	"wxcloudrun-golang/service"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// func GetMonthlyCCERStats(c *gin.Context) {
// 	log.Println("################## GetMonthlyCCERStats ##################")
// 	result, err := service.GetMonthlyCCERStats()

// 	if err != nil {
// 		glog.Errorln("get monthly ccer stats error!")
// 		response.MakeFail(c, http.StatusBadRequest, "get  monthly ccer stats error")
// 		return
// 	}
// 	log.Println("get monthly ccer stats successfully")
// 	response.MakeSuccess(c, http.StatusOK, result)
// 	return
// }

// func GetMonthlyCEAStats(c *gin.Context) {
// 	log.Println("################## GetMonthlyCEAStats ##################")
// 	result, err := service.GetMonthlyCEAStats()

// 	if err != nil {
// 		glog.Errorln("get monthly cea stats error!")
// 		response.MakeFail(c, http.StatusBadRequest, "get  monthly cea stats error")
// 		return
// 	}
// 	log.Println("get monthly cea stats successfully")
// 	response.MakeSuccess(c, http.StatusOK, result)
// 	return
// }

func GetMonthlyCCERExpectation(c *gin.Context) {
	log.Println("################## GetMonthlyCCERExpectation ##################")
	result, err := service.GetMonthlyCCERExpectation()
	if err != nil {
		glog.Errorln("[ERROR]get monthly CCER expectation error!")
		response.MakeFail(c, http.StatusBadRequest, "get  monthly CCER expectation error")
		return
	}
	log.Println("get monthly CCER expectation successfully")
	response.MakeSuccess(c, http.StatusOK, result)
	return
}

func GetMonthlyCEAExpectation(c *gin.Context) {
	log.Println("################## GetMonthlyCEAExpectation ##################")
	result, err := service.GetMonthlyCEAExpectation()
	if err != nil {
		glog.Errorln("[ERROR]get monthly cea expectation error!")
		response.MakeFail(c, http.StatusBadRequest, "get  monthly cea expectation error")
		return
	}
	log.Println("get monthly cea expectation successfully")
	response.MakeSuccess(c, http.StatusOK, result)
	return
}

func GetYearlyCEAExpectation(c *gin.Context) {
	log.Println("################## GetYearlyCEAExpectation ##################")
	// result, err := service.GetMonthlyCEAStats()
	result, err := service.GetYearlyCEAExpectation()
	if err != nil {
		glog.Errorln("[ERROR]get yearly cea expectation error!")
		response.MakeFail(c, http.StatusBadRequest, "get  yearly cea expectation error")
		return
	}
	log.Println("get yearly cea expectation successfully")
	response.MakeSuccess(c, http.StatusOK, result)
	return
}

func GetGECMonthlyStatsByType(c *gin.Context) {
	log.Println("################## GetMonthlyCCERStats ##################")
	result, err := service.GetGECMonthlyStatsByType()

	if err != nil {
		glog.Errorln("get monthly gec stats error!")
		response.MakeFail(c, http.StatusBadRequest, "get  monthly gec stats error")
		return
	}
	log.Println("get monthly gec stats successfully")
	response.MakeSuccess(c, http.StatusOK, result)
	return
}

func GetMonthGroupAVG(c *gin.Context) {
	log.Println("################## Get Month Group AVG ##################")
	nowTimeStr := c.Query("nowTime")
	product := c.Query("product")
	t, err := time.Parse("2006-01-02 15:04:05", nowTimeStr)
	if err != nil {
		c.JSON(400, gin.H{"msg": "nowTime 格式错误", "error": err.Error()})
		return
	}
	result, err := service.GetStatsNextMonth(t, product)
	if err != nil {
		log.Printf("[ERROR]Get Month Group AVG error!")
		response.MakeFail(c, http.StatusBadRequest, "Get Month Group AVG error")
		return
	}
	log.Println("Get Month Group AVG successfully")
	response.MakeSuccess(c, http.StatusOK, result)
	return
}

func GetYearGroupAVG(c *gin.Context) {
	log.Println("################## Get Year Group AVG ##################")
	nowTimeStr := c.Query("nowTime")
	t, err := time.Parse("2006-01-02 15:04:05", nowTimeStr)
	if err != nil {
		c.JSON(400, gin.H{"msg": "nowTime 格式错误", "error": err.Error()})
		return
	}
	result, err := service.GetYearStatsNextMonth(t)
	if err != nil {
		log.Printf("[ERROR]Get Year Group AVG error!")
		response.MakeFail(c, http.StatusBadRequest, "Get Year Group AVG error")
		return
	}
	log.Println("Get Year Group AVG successfully")
	response.MakeSuccess(c, http.StatusOK, result)
	return
}
