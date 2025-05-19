package service

import (
	"errors"
	"strings"
	"time"
	"wxcloudrun-golang/app/handlers/request"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils/crypto"

	"github.com/golang/glog"
	"gorm.io/gorm"
)

const (
	UUID_PREFIX      = "CPIF"
	ADMIN_USER_TYPE  = 1
	VIP_USER_TYPE    = 2
	COMMON_USER_TYPE = 3
)

func QueryUserInfo(account string) (uuid, userId, password string, userType int, Approved bool, err error) {
	var user model.User
	cli := db.Get()
	err = cli.Where("account = ?", account).First(&user).Error
	if err != nil {
		return "", "", "", 0, false, err
	}
	return user.Uuid, user.UserID, user.PasswordHash, user.Type, user.Approved, nil

}

func AddUser(user *request.ReqUser) error {
	//添加判断唯一性的内容（uuid唯一，用户名唯一，公司唯一）
	userId, err := generateUserid(user.Account, user.Type)
	if err != nil {
		return err
	}
	passwordHash := crypto.CalculateSHA256(user.Password, "FDUCPIF")
	approved := true
	if user.Type == 2 {
		approved = false
	}
	cli := db.Get()
	// 检查 UserID 是否唯一
	var count int64
	// cli.Model(&model.User{}).Where("userId = ?", uuid).Count(&count)
	// if count > 0 {
	// 	return errors.New("用户ID已存在")
	// }
	uuid, err := generateUuid(cli, user.Account)
	if err != nil {
		return err
	}
	// 检查 Account 是否唯一
	cli.Model(&model.User{}).Where("account = ?", user.Account).Count(&count)
	if count > 0 {
		return errors.New("账号已存在")
	}

	// 检查 CompanyName 是否唯一
	// cli.Model(&model.User{}).Where("company_name = ?", user.Company).Count(&count)
	// if count > 0 {
	// 	return errors.New("公司已经注册")
	// }
	addUser := model.User{
		Uuid:         uuid,
		UserID:       userId,
		UserName:     user.Name,
		CompanyName:  user.Company,
		CompanyType:  user.CompanyType,
		Type:         user.Type,
		Phone:        user.Phone,
		Email:        user.Email,
		Account:      user.Account,
		PasswordHash: passwordHash,
		Approved:     approved,
	}

	err = cli.Create(&addUser).Error
	if err != nil {
		glog.Errorln("create new user error: %v", err)
		return err
	}
	return nil
}

func ApproveUser(uuid string) error {
	cli := db.Get()
	err := cli.Model(&model.User{}).Where("uuid = ?", uuid).Update("approved", true).Error
	if err != nil {
		return err
	}
	return nil
}

func UpToVipByAdmin(uuid string) error {
	cli := db.Get()
	// 查询该用户
	var user model.User
	if err := cli.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		return err
	}
	// 替换 UserID 前缀
	if strings.HasPrefix(user.UserID, "CPIFB") {
		user.UserID = strings.Replace(user.UserID, "CPIFB", "CPIFA", 1)
	}
	// 设置已通过
	user.Approved = true
	user.Type = 2
	return cli.Save(&user).Error
}

func DownToCommonByAdmin(uuid string) error {
	cli := db.Get()
	// 查询该用户
	var user model.User
	if err := cli.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		return err
	}
	// 替换 UserID 前缀
	if strings.HasPrefix(user.UserID, "CPIFA") {
		user.UserID = strings.Replace(user.UserID, "CPIFA", "CPIFB", 1)
	}
	// 设置已通过
	user.Approved = true
	user.Type = 3
	return cli.Save(&user).Error
}

func ApplyToVip(uuid string) error {
	cli := db.Get()
	// 查询该用户
	var user model.User
	if err := cli.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		return err
	}
	// 设置已通过
	user.Approved = false
	return cli.Save(&user).Error
}

func generateUserid(account string, userType int) (string, error) {
	accountHash := crypto.CalculateSHA256(account, "userId")

	var userId string
	switch userType {
	case ADMIN_USER_TYPE:
		userId = UUID_PREFIX + "0" + accountHash
	case VIP_USER_TYPE:
		userId = UUID_PREFIX + "A" + accountHash
	case COMMON_USER_TYPE:
		userId = UUID_PREFIX + "B" + accountHash
	default:
		return "", errors.New("user type is not exist")
	}
	return userId, nil
}

func generateUuid(db *gorm.DB, account string) (string, error) {
	for i := 0; i < 5; i++ {
		now := time.Now()
		str := now.Format("2006-01-02 15:04:05")
		uuidHash := crypto.CalculateSHA256(account, str)
		uuid := UUID_PREFIX + uuidHash
		var count int64
		err := db.Model(&model.User{}).Where("uuid = ?", uuid).Count(&count).Error
		if err != nil {
			return "", nil
		}
		if count == 0 {
			return uuid, nil
		}
	}
	return "", errors.New("无法生成唯一 UUID，请重试")
}
