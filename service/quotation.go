package service

import (
	"fmt"
	"log"
	"time"
	"wxcloudrun-golang/app/handlers/request"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"

	"errors"

	"github.com/golang/glog"
	"gorm.io/gorm"
)

func AddSemiMonth(quotation *request.ReqQuotation) error {
	var resTime string
	cli := db.Get()
	resTime, qid, err := getSemiMonthTimeAndID(quotation.NowTime, cli)
	if err != nil {
		return err
	}
	semiMonthQuotation := model.SemiMonthQuotation{
		QID:            qid,
		Uuid:           quotation.Uuid,
		Product:        quotation.Product,
		Type:           quotation.Type,
		LowerPrice:     quotation.LowerPrice,
		HigherPrice:    quotation.HigherPrice,
		Price:          quotation.Price,
		TxVolume:       quotation.TxVolume,
		ApplicableTime: resTime,
		Approved:       true,
	}

	// 先查是否重复提交
	var count int64
	cli.Model(&model.SemiMonthQuotation{}).
		Where("userId = ? AND product = ? AND type = ? AND applicableTime = ?",
			quotation.Uuid, quotation.Product, quotation.Type, resTime).
		Count(&count)

	if count > 0 {
		return errors.New("不能重复报价：该用户本期已对该产品报价")
	}

	err = cli.Create(&semiMonthQuotation).Error
	if err != nil {
		glog.Errorln("Submit semimonth quotation error: %v", err)
		return err
	}
	return nil

}

func AddMonth(quotation *request.ReqQuotation) error {
	var resTime string
	cli := db.Get()
	resTime, qid, err := getMonthTimeAndID(quotation.NowTime, cli)
	if err != nil {
		return err
	}
	monthQuotation := model.MonthQuotation{
		QID:            qid,
		Uuid:           quotation.Uuid,
		Product:        quotation.Product,
		Type:           quotation.Type,
		LowerPrice:     quotation.LowerPrice,
		HigherPrice:    quotation.HigherPrice,
		Price:          quotation.Price,
		TxVolume:       quotation.TxVolume,
		ApplicableTime: resTime,
		Approved:       true,
	}

	// 先查是否重复提交
	var count int64
	cli.Model(&model.MonthQuotation{}).
		Where("uuid = ? AND product = ? AND type = ? AND applicableTime = ?",
			quotation.Uuid, quotation.Product, quotation.Type, resTime).
		Count(&count)

	if count > 0 {
		return errors.New("不能重复报价：该用户本期已对该产品报价")
	}

	err = cli.Create(&monthQuotation).Error
	if err != nil {
		glog.Errorln("Submit month quotation error: %v", err)
		return err
	}
	return nil
}

func AddYear(quotation *request.ReqQuotation) error {
	var resTime string
	cli := db.Get()
	submitTime, _, err := getMonthTimeAndID(quotation.NowTime, cli)
	resTime, qid, err := getYearTimeAndID(quotation.NowTime, cli)

	start := quotation.NowTime.Truncate(24 * time.Hour)
	monthStart := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.Local)
	nextMonthStart := monthStart.AddDate(0, 1, 0)

	if err != nil {
		return err
	}
	yearQuotation := model.YearQuotation{
		QID:            qid,
		Uuid:           quotation.Uuid,
		Product:        quotation.Product,
		Type:           quotation.Type,
		LowerPrice:     quotation.LowerPrice,
		HigherPrice:    quotation.HigherPrice,
		Price:          quotation.Price,
		TxVolume:       quotation.TxVolume,
		ApplicableTime: resTime,
		SubmitTime:     submitTime,
		Approved:       true,
	}

	var count int64
	cli.Model(&model.YearQuotation{}).
		Where("userId = ? AND product = ? AND created_at >= ? AND created_at < ?",
			quotation.Uuid, quotation.Product, monthStart, nextMonthStart).
		Count(&count)

	if count > 0 {
		return errors.New("您本月已对该产品提交过报价")
	}

	err = cli.Create(&yearQuotation).Error
	if err != nil {
		glog.Errorln("Submit year quotation error: %v", err)
		return err
	}
	return nil
}

func getSemiMonthTimeAndID(nowTime time.Time, db *gorm.DB) (string, string, error) {
	// now := time.Now()
	// year, month, day := now.Date()
	year, month, day := nowTime.Date()

	dateStr := nowTime.Format("20060102")
	// 查询今天已有多少条报价记录
	var count int64
	// err := db.Model(&model.SemiMonthQuotation{}).
	// 	Where("DATE(created_at) = ?", dateStr).
	// 	Count(&count).Error
	err := db.Model(&model.SemiMonthQuotation{}).Count(&count).Error
	if err != nil {
		return "", "", err
	}
	// 生成 QID（递增 +1）
	seq := count + 1
	qid := fmt.Sprintf("Q%s-%06d", dateStr, seq)

	// 如果今天是14号，返回当月16号到月底
	if day == 14 {
		// 获取当月最后一天
		firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local)
		lastDay := firstOfNextMonth.AddDate(0, 0, -1).Day()

		result := fmt.Sprintf("%04d/%02d/16-%04d/%02d/%02d", year, month, year, month, lastDay)
		return result, qid, nil

	}

	// 如果今天是29号，或特殊月份最后一天（28或27号）
	firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local)
	lastDay := firstOfNextMonth.AddDate(0, 0, -1).Day()

	if (day == 28) || (lastDay < 29 && day == lastDay-1) {
		// 下一个月
		nextMonth := month + 1
		nextYear := year
		if nextMonth > 12 {
			nextMonth = 1
			nextYear++
		}

		result := fmt.Sprintf("%04d/%02d/01-%04d/%02d/15", nextYear, nextMonth, nextYear, nextMonth)
		return result, qid, nil
	}
	return "", "", errors.New("当前时间不可提交半月度报价")
}

func getMonthTimeAndID(nowTime time.Time, db *gorm.DB) (string, string, error) {
	// now := time.Now()
	// year, month, day := now.Date()
	year, month, day := nowTime.Date()
	hour, _, _ := nowTime.Clock()
	log.Println("test hour:", hour)

	dateStr := nowTime.Format("20060102")
	// 查询今天已有多少条报价记录
	var count int64

	err := db.Model(&model.MonthQuotation{}).Count(&count).Error
	if err != nil {
		return "", "", err
	}
	// 生成 QID（递增 +1）
	seq := count + 1
	qid := fmt.Sprintf("Q%s-%06d", dateStr, seq)
	// 获取当前月的最后一天
	// firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local)
	// lastDay := firstOfNextMonth.AddDate(0, 0, -1).Day()

	// 判断是否为28号，或是特殊月的最后一天（28或27）
	// if day == 28 || (lastDay < 29 && day == lastDay-1) {
	if (day < 26 && day > 21) || (day == 26 && hour < 12) {
		// 下一个月
		nextMonth := month + 1
		nextYear := year
		if nextMonth > 12 {
			nextMonth = 1
			nextYear++
		}
		result := fmt.Sprintf("%d年%d月\n", nextYear, nextMonth)
		return result, qid, nil
	}
	return "", "", errors.New("当前时间不可提交月度报价")
}

func getYearTimeAndID(nowTime time.Time, db *gorm.DB) (string, string, error) {
	// now := time.Now()
	// year, month, day := now.Date()
	year, month, day := nowTime.Date()
	hour, _, _ := nowTime.Clock()

	dateStr := nowTime.Format("20060102")
	// 查询今天已有多少条报价记录
	var count int64
	// err := db.Model(&model.YearQuotation{}).
	// 	Where("DATE(created_at) = ?", dateStr).
	// 	Count(&count).Error
	err := db.Model(&model.YearQuotation{}).Count(&count).Error
	if err != nil {
		return "", "", err
	}
	// 生成 QID（递增 +1）
	seq := count + 1
	qid := fmt.Sprintf("Q%s-%06d", dateStr, seq)
	// 获取当前月的最后一天
	// firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local)
	// lastDay := firstOfNextMonth.AddDate(0, 0, -1).Day()

	// 判断是否为29号，或是特殊月的最后一天（28或27）
	// if day == 28 || (lastDay < 29 && day == lastDay-1) {
	// 下一个月
	if (day < 26 && day > 21) || (day == 26 && hour < 12) {
		nextMonth := month + 1
		nextYear := year
		if nextMonth > 12 {
			nextMonth = 1
			nextYear++
		}
		result := fmt.Sprintf("%d年12月\n", nextYear)
		return result, qid, nil
	}
	return "", "", errors.New("当前时间不可提交月度报价")
}

func ApproveQuotation(qID string, model interface{}) error {
	cli := db.Get()
	err := cli.Model(model).Where("qid = ?", qID).Update("approved", true).Error
	if err != nil {
		return err
	}
	return nil
}

func GetApprovedSemimonthQuotations(t time.Time) ([]model.SemiMonthQuotation, error) {
	// now := time.Now()
	// year, month, day := now.Date()
	year, month, day := t.Date()

	// 获取本月最后一天
	firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local)
	lastDay := firstOfNextMonth.AddDate(0, 0, -1).Day()

	var targetPeriod string

	if day == 15 {
		// 15号，查本月16号到月底
		targetPeriod = fmt.Sprintf("%04d/%02d/16-%04d/%02d/%02d", year, month, year, month, lastDay)
	} else if day == 29 || (lastDay < 29 && day == lastDay) {
		// 29号或特殊月最后一天，查下个月1-15号
		nextMonth := month + 1
		nextYear := year
		if nextMonth > 12 {
			nextMonth = 1
			nextYear++
		}
		targetPeriod = fmt.Sprintf("%04d/%02d/01-%04d/%02d/15", nextYear, nextMonth, nextYear, nextMonth)
	} else {
		// 不是公示日
		return nil, errors.New("不是公示日")
	}

	// 查询
	var result []model.SemiMonthQuotation
	cli := db.Get()
	err := cli.Where("applicableTime = ? AND approved = ?", targetPeriod, true).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetApprovedMonthQuotations() ([]model.MonthQuotation, error) {
	now := time.Now()
	year, month, _ := now.Date()
	// year, month, day := t.Date()

	// 获取本月最后一天
	// firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local)
	// lastDay := firstOfNextMonth.AddDate(0, 0, -1).Day()

	var targetPeriod string

	// if day == 29 || (lastDay < 29 && day == lastDay) {
	// 29号或特殊月最后一天，查下个月1-15号
	nextMonth := month + 1
	nextYear := year
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}
	targetPeriod = fmt.Sprintf("%d年%d月\n", nextYear, nextMonth)
	// } else {
	// 不是公示日
	// return nil, errors.New("不是公示日")
	// }

	// 查询
	var result []model.MonthQuotation
	cli := db.Get()
	err := cli.Where("applicableTime = ? AND approved = ?", targetPeriod, true).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetApprovedYearQuotations() ([]model.YearQuotation, error) {
	now := time.Now()
	year, month, _ := now.Date()
	// year, month, day := t.Date()
	// 获取本月最后一天
	// firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local)
	// lastDay := firstOfNextMonth.AddDate(0, 0, -1).Day()

	// var start, end time.Time

	// if day == 29 || (lastDay < 29 && day == lastDay) {
	// // 本月第一天
	// start = time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	// // 下月第一天（即本月最后一秒的下一秒）
	// end = start.AddDate(0, 1, 0)
	nextMonth := month + 1
	nextYear := year
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}
	nowMonth := fmt.Sprintf("%d年%d月\n", nextYear, nextMonth)
	var result []model.YearQuotation
	cli := db.Get()
	err := cli.Where("submitTime = ? AND approved = ?", nowMonth, true).
		Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
	// }
	// else {
	// 	// 不是公示日
	// 	return nil, errors.New("不是公示日")
	// }

}
