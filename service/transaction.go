package service

import (
	"fmt"
	"time"
	"wxcloudrun-golang/app/handlers/request"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"

	"github.com/golang/glog"
	"gorm.io/gorm"
)

func SellerTxSubmit(tx *request.ReqTransaction) error {
	cli := db.Get()
	tid, err := generateTID(cli, &model.SellerTx{})
	if err != nil {
		return err
	}
	sellerTx := model.SellerTx{
		TID:      tid,
		Uuid:     tx.Uuid,
		Project:  tx.Project,
		Type:     tx.Type,
		Price:    tx.Price,
		TxVolume: tx.TxVolume,
		State:    1,
	}
	err = cli.Create(&sellerTx).Error
	if err != nil {
		glog.Errorln("Submit seller tx error: %v", err)
		return err
	}
	return nil
}

func BuyerTxSubmit(tx *request.ReqTransaction) error {
	cli := db.Get()
	tid, err := generateTID(cli, &model.BuyerTx{})
	if err != nil {
		return err
	}
	buyerTx := model.BuyerTx{
		TID:      tid,
		Uuid:     tx.Uuid,
		Project:  tx.Project,
		Type:     tx.Type,
		Price:    tx.Price,
		TxVolume: tx.TxVolume,
		State:    1,
	}
	err = cli.Create(&buyerTx).Error
	if err != nil {
		glog.Errorln("Submit buyer tx error: %v", err)
		return err
	}
	return nil
}

func GetSellerTx() ([]model.SellerTx, error) {
	var result []model.SellerTx
	cli := db.Get()
	err := cli.Where("state = ?", 1).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetBuyerTx() ([]model.BuyerTx, error) {
	var result []model.BuyerTx
	cli := db.Get()
	err := cli.Where("state = ?", 1).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func SubmitNotition(n *request.ReqNotition) error {
	cli := db.Get()
	var user model.User
	err := cli.Where("uuid = ?", n.Uuid).First(&user).Error
	if err != nil {
		return err
	}
	nid, err := generateNID(cli)
	if err != nil {
		return err
	}
	notition := model.Notition{
		NID:         nid,
		TID:         n.Tid,
		Uuid:        n.Uuid,
		Type:        n.Type,
		UserName:    user.UserName,
		CompanyName: user.CompanyName,
		Phone:       user.Phone,
		Email:       user.Email,
		State:       1,
	}
	err = cli.Create(&notition).Error
	if err != nil {
		glog.Errorln("Submit notition error: %v", err)
		return err
	}
	return nil
}

func generateTID(db *gorm.DB, model interface{}) (string, error) {
	now := time.Now()
	dateStr := now.Format("20060102")

	var count int64
	err := db.Model(model).Count(&count).Error
	if err != nil {
		return "", err
	}

	seq := count + 1
	qid := fmt.Sprintf("Q%s-%06d", dateStr, seq)
	return qid, nil
}

func generateNID(db *gorm.DB) (string, error) {
	now := time.Now()
	dateStr := now.Format("20060102")

	var count int64
	err := db.Model(model.Notition{}).Count(&count).Error
	if err != nil {
		return "", err
	}
	seq := count + 1
	nid := fmt.Sprintf("N%s-%06d", dateStr, seq)
	return nid, nil
}
