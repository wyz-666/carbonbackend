package service

import (
	"carbonbackend/app/handlers/request"
	"carbonbackend/db"
	"carbonbackend/db/model"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"errors"

	"github.com/golang/glog"
	"gorm.io/gorm"
)

func AddSemiMonth(quotation *request.ReqQuotation) error {
	var resTime string
	cli := db.Get()
	// resTime, qid, err := getSemiMonthTimeAndID(quotation.NowTime, cli)
	resTime, qid, err := getSemiMonthTimeAndID(cli)
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
	// resTime, qid, err := getMonthTimeAndID(quotation.NowTime, cli)
	resTime, qid, err := getMonthTimeAndID(cli)
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
	submitTime, _, err := getMonthTimeAndID(cli)
	if err != nil {
		return err
	}
	now := time.Now()
	newTime := now.Add(8 * time.Hour)
	resTime, qid, err := getYearTimeAndID(newTime, cli)

	start := newTime.Truncate(24 * time.Hour)
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

// func getSemiMonthTimeAndID(nowTime time.Time, db *gorm.DB) (string, string, error) {
func getSemiMonthTimeAndID(db *gorm.DB) (string, string, error) {
	now := time.Now()
	newTime := now.Add(8 * time.Hour)
	year, month, day := newTime.Date()
	// year, month, day := nowTime.Date()

	// dateStr := nowTime.Format("20060102")
	dateStr := newTime.Format("20060102")
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

// func getMonthTimeAndID(nowTime time.Time, db *gorm.DB) (string, string, error) {
func getMonthTimeAndID(db *gorm.DB) (string, string, error) {
	now := time.Now()
	newTime := now.Add(8 * time.Hour)
	year, month, day := newTime.Date()
	// year, month, day := nowTime.Date()
	hour, _, _ := newTime.Clock()
	log.Println("test hour:", hour)

	dateStr := newTime.Format("20060102")
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
	if (day < 27 && day > 21) || (day == 27 && hour < 19) {
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
	if (day < 27 && day > 21) || (day == 27 && hour < 19) {
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
	year, month, day := now.Date()
	hour, _, _ := now.Clock()
	// year, month, day := t.Date()

	// 获取本月最后一天
	// firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local)
	// lastDay := firstOfNextMonth.AddDate(0, 0, -1).Day()

	var targetPeriod string

	if (day == 26 && hour > 12) || (day > 26) {

		nextMonth := month + 1
		nextYear := year
		if nextMonth > 12 {
			nextMonth = 1
			nextYear++
		}
		targetPeriod = fmt.Sprintf("%d年%d月\n", nextYear, nextMonth)
	} else {
		//不是公示日
		return nil, errors.New("不是公示日")
	}

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

func AttachUserNameToMonthQuotations(db *gorm.DB, quotes []model.MonthQuotation) ([]model.MonthQuotationWithUser, error) {
	var result []model.MonthQuotationWithUser

	for _, quote := range quotes {
		var user model.User
		if err := db.Where("uuid = ?", quote.Uuid).First(&user).Error; err != nil {
			// 记录找不到时可以跳过，或者记录空用户名
			result = append(result, model.MonthQuotationWithUser{
				MonthQuotation: quote,
				UserName:       user.UserName,
			})
			continue
		}
		result = append(result, model.MonthQuotationWithUser{
			MonthQuotation: quote,
			UserName:       user.UserName,
		})
	}

	return result, nil
}

func AttachUserNameToYearQuotations(db *gorm.DB, quotes []model.YearQuotation) ([]model.YearQuotationWithUser, error) {
	var result []model.YearQuotationWithUser

	for _, quote := range quotes {
		var user model.User
		if err := db.Where("uuid = ?", quote.Uuid).First(&user).Error; err != nil {
			// 记录找不到时可以跳过，或者记录空用户名
			result = append(result, model.YearQuotationWithUser{
				YearQuotation: quote,
				UserName:      user.UserName,
			})
			continue
		}
		result = append(result, model.YearQuotationWithUser{
			YearQuotation: quote,
			UserName:      user.UserName,
		})
	}

	return result, nil
}

func GetAllApprovedMonthQuotations() ([]model.MonthQuotation, error) {
	cli := db.Get()

	// 使用当前时间 +8小时，模拟北京时间（不使用 Location）
	now := time.Now().Add(8 * time.Hour)

	// 查出所有已通过审批的记录
	var all []model.MonthQuotation
	err := cli.Where("approved = ?", true).Find(&all).Error
	if err != nil {
		return nil, err
	}

	var result []model.MonthQuotation

	for _, q := range all {
		appTime := strings.TrimSpace(q.ApplicableTime) // 例如 "2025年4月"

		var y, m int
		_, err := fmt.Sscanf(appTime, "%d年%d月", &y, &m)
		if err != nil {
			continue // 跳过解析失败的数据
		}

		// 能展示的时间是 (applicableTime的上一个月) 的 26日19点（北京时间）
		showYear := y
		showMonth := m - 1
		if showMonth <= 0 {
			showMonth = 12
			showYear--
		}

		canShowTime := time.Date(showYear, time.Month(showMonth), 26, 19, 0, 0, 0, time.UTC).Add(8 * time.Hour)

		if now.After(canShowTime) {
			result = append(result, q)
		}
	}

	// 按照 applicableTime 倒序排列
	sort.Slice(result, func(i, j int) bool {
		var y1, m1, y2, m2 int
		fmt.Sscanf(strings.TrimSpace(result[i].ApplicableTime), "%d年%d月", &y1, &m1)
		fmt.Sscanf(strings.TrimSpace(result[j].ApplicableTime), "%d年%d月", &y2, &m2)
		if y1 != y2 {
			return y1 > y2
		}
		return m1 > m2
	})

	return result, nil
}

func GetAllApprovedYearQuotations() ([]model.YearQuotation, error) {
	cli := db.Get()

	// 使用当前时间 +8小时，模拟北京时间（不使用 Location）
	now := time.Now().Add(8 * time.Hour)

	// 查出所有已通过审批的记录
	var all []model.YearQuotation
	err := cli.Where("approved = ?", true).Find(&all).Error
	if err != nil {
		return nil, err
	}

	var result []model.YearQuotation

	for _, q := range all {
		appTime := strings.TrimSpace(q.SubmitTime) // 例如 "2025年4月"

		var y, m int
		_, err := fmt.Sscanf(appTime, "%d年%d月", &y, &m)
		if err != nil {
			continue // 跳过解析失败的数据
		}

		// 能展示的时间是 (applicableTime的上一个月) 的 26日19点（北京时间）
		showYear := y
		showMonth := m - 1
		if showMonth <= 0 {
			showMonth = 12
			showYear--
		}

		canShowTime := time.Date(showYear, time.Month(showMonth), 26, 19, 0, 0, 0, time.UTC).Add(8 * time.Hour)

		if now.After(canShowTime) {
			result = append(result, q)
		}
	}

	// 按照 applicableTime 倒序排列
	sort.Slice(result, func(i, j int) bool {
		var y1, m1, y2, m2 int
		fmt.Sscanf(strings.TrimSpace(result[i].SubmitTime), "%d年%d月", &y1, &m1)
		fmt.Sscanf(strings.TrimSpace(result[j].SubmitTime), "%d年%d月", &y2, &m2)
		if y1 != y2 {
			return y1 > y2
		}
		return m1 > m2
	})

	return result, nil
}
