package handlers

import (
	"log"
	"net/http"
	"wxcloudrun-golang/app/handlers/request"
	"wxcloudrun-golang/app/handlers/response"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/service"
	"wxcloudrun-golang/utils/crypto"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func Register(c *gin.Context) {

	log.Println("################## User Register ##################")
	var reqUser request.ReqUser
	if err := c.ShouldBind(&reqUser); err != nil || !checkRegisterParams(&reqUser) {
		log.Printf("[ERROR] User registration error: %v", err)
		response.MakeFail(c, http.StatusNotAcceptable, "user registration failure!")
		return
	}
	log.Println("request user parameters:")
	log.Println(reqUser)
	err := service.AddUser(&reqUser)
	if err != nil {
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}

	log.Println("add new user successful!")
	response.MakeSuccess(c, http.StatusOK, "successfully register the user!")
	return
}

func checkRegisterParams(reqUser *request.ReqUser) bool {
	if reqUser.Name == "" || reqUser.Account == "" || reqUser.Type == 0 {
		log.Printf("[ERROR] Missing user registration parameters")
		return false
	}
	ps := reqUser.Password
	if len(ps) < 9 {
		log.Printf("[ERROR] password len is < 9")
		return false
	}
	return true
}

func Login(c *gin.Context) {
	log.Println("################## User Login ##################")
	var reqLogin request.ReqLogin
	if err := c.ShouldBind(&reqLogin); err != nil || !checkLoginParams(&reqLogin) {
		response.MakeFail(c, http.StatusBadRequest, "login parameters error!")
		return
	}
	reqPwdHash := crypto.CalculateSHA256(reqLogin.Password, "FDUCPIF")
	uuid, userID, pwdHash, userType, approved, err := service.QueryUserInfo(reqLogin.Account)
	if err != nil {
		log.Printf("[ERROR] query fsims password hash error: %v", err)
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}

	if userType == 2 && !approved {
		log.Println("[ERROR] 该用户资格审批尚未通过!")
		response.MakeFail(c, http.StatusBadRequest, "该用户资格审批尚未通过!")
		return
	}

	// check password
	if reqPwdHash != pwdHash {
		log.Println("[ERROR] password incorrect!")
		response.MakeFail(c, http.StatusBadRequest, "password incorrecrt!")
		return
	}

	// create jwt token with uuid
	token, err := service.CreateJwtToken(uuid, userID, userType)
	if err != nil {
		log.Println("[ERROR] create jwt token error!")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Println(reqLogin.Account, " login successful!")
	log.Println("token:", token)

	resLogin := response.ResLogin{
		Uuid:     uuid,
		UserID:   userID,
		UserType: userType,
		Token:    token,
	}
	response.MakeSuccessAdmin(c, http.StatusOK, "登录成功", resLogin)
}

func ApproveUser(c *gin.Context) {
	log.Println("################## Approve VIP User ##################")
	var reqApprove request.ReqApproveUser
	if err := c.ShouldBind(&reqApprove); err != nil {
		log.Printf("[ERROR] approve user error: %v", err)
		response.MakeFail(c, http.StatusNotAcceptable, "approve user failure!")
		return
	}
	err := service.ApproveUser(reqApprove.Uuid)
	if err != nil {
		log.Printf("[ERROR] approve vip user error:%v", err)
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	response.MakeSuccess(c, http.StatusOK, "approve vip user successfully")

}

func GetAllUser(c *gin.Context) {
	log.Println("################## Get All User ##################")
	var users []model.User

	cli := db.Get()
	err := cli.Find(&users).Error
	if err != nil {
		log.Printf("[ERROR] query all users error:%v", err)
		response.MakeFail(c, http.StatusBadRequest, "query all user error")
		return
	}
	log.Println("query all users successful")
	response.MakeSuccess(c, http.StatusOK, users)
	return
}

func GetAllApprovingUser(c *gin.Context) {
	log.Println("################## Get All Approving User ##################")
	var users []model.User

	cli := db.Get()
	err := cli.Where("approved = ?", false).Find(&users).Error
	if err != nil {
		log.Printf("[ERROR] query all approving users error:%v", err)
		response.MakeFail(c, http.StatusBadRequest, "query all approving user error")
		return
	}
	log.Println("query all approving users successful")
	response.MakeSuccess(c, http.StatusOK, users)
	return
}

func UpToVipByAdmin(c *gin.Context) {
	log.Println("################## Upgrade VIP User By Admin ##################")
	uuid := c.Query("uuid")
	err := service.UpToVipByAdmin(uuid)
	if err != nil {
		log.Printf("[ERROR] upgrade vip user error:%v", err)
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	response.MakeSuccess(c, http.StatusOK, "upgrade vip user successfully")
}

func DownToCommonByAdmin(c *gin.Context) {
	log.Println("################## Down User By Admin ##################")
	uuid := c.Query("uuid")
	err := service.DownToCommonByAdmin(uuid)
	if err != nil {
		glog.Errorln("dowm vip user error!")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	response.MakeSuccess(c, http.StatusOK, "down vip user successfully")
}

func ApplyToVip(c *gin.Context) {
	log.Println("################## Apply VIP User ##################")
	var reqApprove request.ReqApproveUser
	if err := c.ShouldBind(&reqApprove); err != nil {
		glog.Errorln("apply user error")
		response.MakeFail(c, http.StatusNotAcceptable, "apply user failure!")
		return
	}
	err := service.ApplyToVip(reqApprove.Uuid)
	if err != nil {
		glog.Errorln("apply vip user error!")
		response.MakeFail(c, http.StatusBadRequest, err.Error())
		return
	}
	response.MakeSuccess(c, http.StatusOK, "apply vip user successfully")
}

func checkLoginParams(reqLogin *request.ReqLogin) bool {
	if reqLogin.Account == "" || reqLogin.Password == "" {
		glog.Errorln("Missing login parameters")
		return false
	}
	return true
}
