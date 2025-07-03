package handlers

import (
	"carbonbackend/app/handlers/request"
	"carbonbackend/app/handlers/response"
	"carbonbackend/db"
	"carbonbackend/db/model"
	"carbonbackend/service"
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func MultiSubmit(c *gin.Context) {
	log.Println("################## Multi Quotation Submit ##################")
	var quotations []request.ReqQuotation
	if err := c.ShouldBind(&quotations); err != nil {
		glog.Errorln(err.Error())
		response.MakeFail(c, http.StatusNotAcceptable, err.Error())
		return
	}
	var err error
	for _, quotation := range quotations {
		switch quotation.FormName {
		case "SemiMonth":
			err = service.AddSemiMonth(&quotation)
		case "Month":
			err = service.AddMonth(&quotation)
		case "Year":
			err = service.AddYear(&quotation)
		default:
			response.MakeFail(c, http.StatusBadRequest, "form name is wrong!")
			return
		}
		if err != nil {
			response.MakeFail(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	log.Println("multi quotation submit successful")
	response.MakeSuccess(c, http.StatusOK, "successfully submit multi quotation!")
	return
}

func SemiMonthSubmit(c *gin.Context) {
	log.Println("################## SemiMonth Quotation Submit ##################")
	var reqQuotation request.ReqQuotation
	if err := c.ShouldBind(&reqQuotation); err != nil {
		glog.Errorln(err.Error())
		response.MakeFail(c, http.StatusNotAcceptable, err.Error())
		return
	}
	log.Println("semimonth quotation submit user:")
	log.Println(reqQuotation.Uuid)
	err := service.AddSemiMonth(&reqQuotation)
	if err != nil {
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("semimonth quotation submit successful")
	response.MakeSuccess(c, http.StatusOK, "successfully submit semimonth quotation!")
	return
}

func MonthSubmit(c *gin.Context) {
	log.Println("################## Month Quotation Submit ##################")
	var reqQuotation request.ReqQuotation
	if err := c.ShouldBind(&reqQuotation); err != nil {
		glog.Errorln(err.Error())
		response.MakeFail(c, http.StatusNotAcceptable, err.Error())
		return
	}
	log.Println("month quotation submit user:")
	log.Println(reqQuotation.Uuid)
	err := service.AddMonth(&reqQuotation)
	if err != nil {
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("month quotation submit successful")
	response.MakeSuccess(c, http.StatusOK, "successfully submit month quotation!")
	return
}

func YearSubmit(c *gin.Context) {
	log.Println("################## Year Quotation Submit ##################")
	var reqQuotation request.ReqQuotation
	if err := c.ShouldBind(&reqQuotation); err != nil {
		glog.Errorln(err.Error())
		response.MakeFail(c, http.StatusNotAcceptable, err.Error())
		return
	}
	log.Println("year quotation submit user:")
	log.Println(reqQuotation.Uuid)
	err := service.AddYear(&reqQuotation)
	if err != nil {
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("year quotation submit successful")
	response.MakeSuccess(c, http.StatusOK, "successfully submit year quotation!")
	return
}

func ApproveQuotation(c *gin.Context) {
	log.Println("################## Approve Quotation ##################")
	var reqApproveQuotation request.ReqApproveQuotation
	if err := c.ShouldBind(&reqApproveQuotation); err != nil {
		glog.Errorln("approve quotation error")
		response.MakeFail(c, http.StatusNotAcceptable, "approve quotation failure!")
		return
	}
	var modelType interface{}
	switch reqApproveQuotation.Type {
	case "semimonth":
		modelType = &model.SemiMonthQuotation{}
	case "month":
		modelType = &model.MonthQuotation{}
	case "year":
		modelType = &model.YearQuotation{}
	default:
		response.MakeFail(c, http.StatusBadRequest, "quotation type error!")
		return
	}
	err := service.ApproveQuotation(reqApproveQuotation.QID, modelType)
	if err != nil {
		glog.Errorln("approve quotation error!")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	response.MakeSuccess(c, http.StatusOK, "approve quotation successfully")

}
func SemiMonthPublish(c *gin.Context) {
	log.Println("################## Publish SemiMonth Quotation ##################")
	timestr := c.Query("time")
	t, err := time.Parse("2006-01-02 15:04:05", timestr)
	res, err := service.GetApprovedSemimonthQuotations(t)
	if err != nil {
		glog.Errorln("publish semimonth quotation error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("publish semimonth quotation successful")
	response.MakeSuccess(c, http.StatusOK, res)
	return
}
func MonthPublish(c *gin.Context) {
	log.Println("################## Publish Month Quotation ##################")
	// timestr := c.Query("time")
	// t, err := time.Parse("2006-01-02 15:04:05", timestr)
	// res, err := service.GetApprovedMonthQuotations()
	res, err := service.GetAllApprovedMonthQuotations()
	if err != nil {
		glog.Errorln("publish month quotation error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("publish month quotation successful")
	response.MakeSuccess(c, http.StatusOK, res)
	return
}
func YearPublish(c *gin.Context) {
	log.Println("################## Publish Year Quotation ##################")
	// timestr := c.Query("time")
	// t, err := time.Parse("2006-01-02 15:04:05", timestr)
	// res, err := service.GetApprovedYearQuotations()
	res, err := service.GetAllApprovedYearQuotations()
	if err != nil {
		glog.Errorln("publish year quotation error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("publish year quotation successful")
	response.MakeSuccess(c, http.StatusOK, res)
	return
}

func GetApprovingSemiMonthQuotations(c *gin.Context) {
	log.Println("################## Get Approving SemiMonth Quotations ##################")
	var quotations []model.SemiMonthQuotation

	cli := db.Get()
	err := cli.Where("approved = ?", false).Find(&quotations).Error
	if err != nil {
		glog.Errorln("get all approving semimonth quotations error!")
		response.MakeFail(c, http.StatusBadRequest, "get all approving semimonth quotations error")
		return
	}
	log.Println("get all approving semimonth quotations successful")
	response.MakeSuccess(c, http.StatusOK, quotations)
	return
}

func GetApprovingMonthQuotations(c *gin.Context) {
	log.Println("################## Get Approving Month Quotations ##################")
	var quotations []model.MonthQuotation

	cli := db.Get()
	err := cli.Where("approved = ?", false).Find(&quotations).Error
	if err != nil {
		glog.Errorln("get all approving month quotations error!")
		response.MakeFail(c, http.StatusBadRequest, "get all approving month quotations error")
		return
	}
	log.Println("get all approving month quotations successful")
	response.MakeSuccess(c, http.StatusOK, quotations)
	return
}

func GetApprovingYearQuotations(c *gin.Context) {
	log.Println("################## Get Approving Month Quotations ##################")
	var quotations []model.YearQuotation

	cli := db.Get()
	err := cli.Where("approved = ?", false).Find(&quotations).Error
	if err != nil {
		glog.Errorln("get all approving year quotations error!")
		response.MakeFail(c, http.StatusBadRequest, "get all approving year quotations error")
		return
	}
	log.Println("get all approving year quotations successful")
	response.MakeSuccess(c, http.StatusOK, quotations)
	return
}

func AdminGetMonthQuotation(c *gin.Context) {
	log.Println("################## Admin Get Month Quotation ##################")
	timestr := c.Query("nowTime")
	product := c.Query("product")
	t, err := time.Parse("2006-01-02 15:04:05", timestr)
	if err != nil {
		c.JSON(400, gin.H{"msg": "nowTime 格式错误", "error": err.Error()})
		return
	}
	y, m, _ := t.Date()
	nextMonth := m + 1
	nextYear := y
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}
	nextMonthStr := fmt.Sprintf("%d年%d月\n", nextYear, nextMonth)
	var quoes []model.MonthQuotation
	cli := db.Get()
	err = cli.Where("applicableTime = ? AND product = ?", nextMonthStr, product).Find(&quoes).Error
	if err != nil {
		glog.Errorln("[ERROR]Admin Get Month Quotation error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	var res []model.MonthQuotationWithUser
	res, err = service.AttachUserNameToMonthQuotations(cli, quoes)
	if err != nil {
		glog.Errorln("[ERROR]Admin Get Month Quotation with User error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("Admin Get Month Quotation successful")
	response.MakeSuccess(c, http.StatusOK, res)
	return
}

func AdminGetYearQuotation(c *gin.Context) {
	log.Println("################## Admin Get Year Quotation ##################")
	timestr := c.Query("time")
	product := c.Query("product")
	t, err := time.Parse("2006-01-02 15:04:05", timestr)
	y, m, _ := t.Date()
	nextMonth := m + 1
	nextYear := y
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}
	nextMonthStr := fmt.Sprintf("%d年%d月\n", nextYear, nextMonth)
	var quoes []model.YearQuotation
	cli := db.Get()
	err = cli.Where("submitTime = ? AND product = ?", nextMonthStr, product).Find(&quoes).Error
	if err != nil {
		glog.Errorln("[ERROR]Admin Get Year Quotation error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	var res []model.YearQuotationWithUser
	res, err = service.AttachUserNameToYearQuotations(cli, quoes)
	if err != nil {
		glog.Errorln("[ERROR]Admin Get Year Quotation with User error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("Admin Get Year Quotation successful")
	response.MakeSuccess(c, http.StatusOK, res)
	return
}

func GetCEAMonthExpectation(c *gin.Context) {
	log.Println("################## Get CEA Month Expectation ##################")
	var res []model.CEAMonthExpectation
	cli := db.Get()
	err := cli.Find(&res).Error
	if err != nil {
		glog.Errorln("[ERROR]Get CEA Month Expectation")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("Get CEA Month Expectation successful")
	response.MakeSuccess(c, http.StatusOK, res)
	return
}

func GetCEAYearExpectation(c *gin.Context) {
	log.Println("################## Get CEA Year Expectation ##################")
	var res []model.CEAYearExpectation
	cli := db.Get()
	err := cli.Find(&res).Error
	if err != nil {
		glog.Errorln("[ERROR]Get CEA Year Expectation")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("Get CEA Year Expectation successful")
	response.MakeSuccess(c, http.StatusOK, res)
	return
}

func GetCCERMonthExpectation(c *gin.Context) {
	log.Println("################## Get CCER Month Expectation ##################")
	var res []model.CCERMonthExpectation
	cli := db.Get()
	err := cli.Find(&res).Error
	if err != nil {
		glog.Errorln("[ERROR]Get CCER Month Expectation")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("Get CCER Month Expectation successful")
	response.MakeSuccess(c, http.StatusOK, res)
	return
}

func GetGECMonthExpectation(c *gin.Context) {
	log.Println("################## Get GEC Month Expectation ##################")
	var res []model.GECMonthExpectation
	cli := db.Get()
	err := cli.Find(&res).Error
	if err != nil {
		glog.Errorln("[ERROR]Get GEC Month Expectation")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("Get GEC Month Expectation successful")
	response.MakeSuccess(c, http.StatusOK, res)
	return
}
