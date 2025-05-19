package handlers

import (
	"log"
	"net/http"
	"wxcloudrun-golang/app/handlers/request"
	"wxcloudrun-golang/app/handlers/response"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/service"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func SellerTxSubmit(c *gin.Context) {
	log.Println("################## Seller Transaction Submit ##################")
	var reqSeller request.ReqTransaction
	if err := c.ShouldBind(&reqSeller); err != nil {
		glog.Errorln(err.Error())
		response.MakeFail(c, http.StatusNotAcceptable, err.Error())
		return
	}
	log.Println("Seller:")
	log.Println(reqSeller.Uuid)
	err := service.SellerTxSubmit(&reqSeller)
	if err != nil {
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("seller transaction submit successful")
	response.MakeSuccess(c, http.StatusOK, "successfully submit seller transaction!")
	return
}

func BuyerTxSubmit(c *gin.Context) {
	log.Println("################## Buyer Transaction Submit ##################")
	var reqBuyer request.ReqTransaction
	if err := c.ShouldBind(&reqBuyer); err != nil {
		glog.Errorln(err.Error())
		response.MakeFail(c, http.StatusNotAcceptable, err.Error())
		return
	}
	log.Println("Buyer:")
	log.Println(reqBuyer.Uuid)
	err := service.BuyerTxSubmit(&reqBuyer)
	if err != nil {
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("buyer transaction submit successful")
	response.MakeSuccess(c, http.StatusOK, "successfully submit buyer transaction!")
	return
}

func SellerTxPublish(c *gin.Context) {
	log.Println("################## Publish Seller Transaction ##################")
	res, err := service.GetSellerTx()
	if err != nil {
		glog.Errorln("publish seller transaction error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("publish seller transaction successful")
	response.MakeSuccess(c, http.StatusOK, res)
	return
}

func BuyerTxPublish(c *gin.Context) {
	log.Println("################## Publish Buyer Transaction ##################")
	res, err := service.GetBuyerTx()
	if err != nil {
		glog.Errorln("publish buyer transaction error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("publish buyer transaction successful")
	response.MakeSuccess(c, http.StatusOK, res)
	return
}

func GetAllBuyerTx(c *gin.Context) {
	log.Println("################## Get All Buyer Tx ##################")
	var txs []model.BuyerTx
	cli := db.Get()
	err := cli.Find(&txs).Error
	if err != nil {
		glog.Errorln("GetAllBuyerTx error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("GetAllBuyerTx successful")
	response.MakeSuccess(c, http.StatusOK, txs)
	return
}

func GetAllSellerTx(c *gin.Context) {
	log.Println("################## Get All Seller Tx ##################")
	var txs []model.SellerTx
	cli := db.Get()
	err := cli.Find(&txs).Error
	if err != nil {
		glog.Errorln("GetAllSellerTx error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("GetAllSellerTx successful")
	response.MakeSuccess(c, http.StatusOK, txs)
	return
}

func GetNotion(c *gin.Context) {
	log.Println("################## Get Notion ##################")
	var notions []model.Notition
	cli := db.Get()
	err := cli.Find(&notions).Error
	if err != nil {
		glog.Errorln("get notition error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("get notition successful")
	response.MakeSuccess(c, http.StatusOK, notions)
	return
}

func GetNotionByState(c *gin.Context) {
	log.Println("################## Get Notion By State##################")
	var notions []model.Notition
	cli := db.Get()
	err := cli.Where("state = ?", 1).Find(&notions).Error
	if err != nil {
		glog.Errorln("get notition by state error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("get notition by state successful")
	response.MakeSuccess(c, http.StatusOK, notions)
	return
}

func FixSellerTx(c *gin.Context) {
	log.Println("################## Fix Seller Tx ##################")
	tid := c.Query("tid")
	cli := db.Get()
	err := cli.Model(model.SellerTx{}).Where("tid = ?", tid).Update("state", 2).Error
	if err != nil {
		glog.Errorln("Fix Seller Tx error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("Fix Seller Tx successful")
	response.MakeSuccess(c, http.StatusOK, "Fix Seller Tx successful")
	return
}

func FixBuyerTx(c *gin.Context) {
	log.Println("################## Fix Buyer Tx ##################")
	tid := c.Query("tid")
	cli := db.Get()
	err := cli.Model(model.BuyerTx{}).Where("tid = ?", tid).Update("state", 2).Error
	if err != nil {
		glog.Errorln("Fix Buyer Tx error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("Fix Buyer Tx successful")
	response.MakeSuccess(c, http.StatusOK, "Fix Buyer Tx successful")
	return
}

func FixNotion(c *gin.Context) {
	log.Println("################## Fix Notion ##################")
	nid := c.Query("nid")
	cli := db.Get()
	err := cli.Model(model.Notition{}).Where("nid = ?", nid).Update("state", 2).Error
	if err != nil {
		glog.Errorln("Fix Notion by state error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("Fix Notion by state successful")
	response.MakeSuccess(c, http.StatusOK, "Fix Notion by state successful")
	return
}

func SubmitNotition(c *gin.Context) {
	log.Println("################## Submit Notition ##################")
	var reqNotition request.ReqNotition
	if err := c.ShouldBind(&reqNotition); err != nil {
		glog.Errorln(err.Error())
		response.MakeFail(c, http.StatusNotAcceptable, err.Error())
		return
	}
	log.Println("Seller:")
	log.Println(reqNotition.Uuid)
	err := service.SubmitNotition(&reqNotition)
	if err != nil {
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("Submit Notition successful")
	response.MakeSuccess(c, http.StatusOK, "Submit Notition transaction!")
	return
}

func SubmitBoard(c *gin.Context) {
	log.Println("################## Submit Board ##################")
	var reqBoard request.ReqBoard
	if err := c.ShouldBind(&reqBoard); err != nil {
		glog.Errorln(err.Error())
		response.MakeFail(c, http.StatusNotAcceptable, err.Error())
		return
	}
	board := model.Board{
		Date:    reqBoard.Date,
		Content: reqBoard.Content,
	}
	cli := db.Get()
	err := cli.Create(&board).Error
	if err != nil {
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("Submit Board successful")
	response.MakeSuccess(c, http.StatusOK, "Submit Board transaction!")
	return
}

func GetLatestBoard(c *gin.Context) {
	log.Println("################## Get Latest Board ##################")
	var board model.Board
	cli := db.Get()
	err := cli.Order("id desc").First(&board).Error
	if err != nil {
		glog.Errorln("get latest board error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("get latest board successful")
	response.MakeSuccess(c, http.StatusOK, board)
	return
}

func GetAllBoards(c *gin.Context) {
	log.Println("################## Get All Board ##################")
	var board []model.Board
	cli := db.Get()
	err := cli.Find(&board).Error
	if err != nil {
		glog.Errorln("get all board error")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("get all board successful")
	response.MakeSuccess(c, http.StatusOK, board)
	return
}
