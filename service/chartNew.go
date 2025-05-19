package service

import (
	"carbonbackend/app/handlers/response"
	"carbonbackend/db"
	"carbonbackend/db/model"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
	// "gorm.io/gorm"
)

// GetStatsNextMonth 按用户类型和公司类型分组计算 CEA 下月报价平均值
func GetStatsNextMonth(nowTime time.Time, product string) ([]response.CategoryStats, error) {
	// 1. 查出下月 CEA 的所有报价
	cli := db.Get()

	// now := time.Now()
	y, m, _ := nowTime.Date()
	nextMonth := m + 1
	nextYear := y
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}
	nextMonthStr := fmt.Sprintf("%d年%d月\n", nextYear, nextMonth)
	var quotes []model.MonthQuotation
	if err := cli.
		Where("applicableTime = ? AND product = ?", nextMonthStr, product).
		Find(&quotes).Error; err != nil {
		return nil, err
	}

	// 2. 收集所有 UserID 并批量查 User 表
	uuids := make([]string, 0, len(quotes))
	for _, q := range quotes {
		uuids = append(uuids, q.Uuid)
	}
	var users []model.User
	if err := cli.
		Where("uuid IN ?", uuids).
		Find(&users).Error; err != nil {
		return nil, err
	}
	// 建一个 map[userId]User 方便后续查
	userMap := make(map[string]model.User, len(users))
	for _, u := range users {
		userMap[u.Uuid] = u
		log.Println("UserID", u.Uuid)
	}

	// 3. 在 Go 里合并数据并分组统计
	type key struct {
		t     int
		cType string
	}
	sums := make(map[key]struct {
		sumL, sumH float64
		count      int
	})

	for _, q := range quotes {
		// 价格字段是 string，要先转 float
		low, err1 := strconv.ParseFloat(q.LowerPrice, 64)
		high, err2 := strconv.ParseFloat(q.HigherPrice, 64)
		if err1 != nil || err2 != nil {
			continue
		}
		log.Println("Test_Uuid:", q.Uuid)
		u, ok := userMap[q.Uuid]
		if !ok {
			continue
		}
		k := key{t: u.Type, cType: u.CompanyType}
		log.Println("Test_key:", k)
		tmp := sums[k]
		tmp.sumL += low
		tmp.sumH += high
		tmp.count++
		sums[k] = tmp
	}

	// 4. 构造最终结果
	var stats []response.CategoryStats
	for k, v := range sums {
		avgL := v.sumL / float64(v.count)
		avgH := v.sumH / float64(v.count)
		stats = append(stats, response.CategoryStats{
			Type:        k.t,
			CompanyType: k.cType,
			AvgLower:    math.Round(avgL*100) / 100,
			AvgHigher:   math.Round(avgH*100) / 100,
			AvgBoth:     math.Round(((avgL+avgH)/2)*100) / 100,
		})
	}

	return stats, nil
}

// GetYearStatsNextMonth 按用户类型和公司类型分组计算 CEA 下月报价平均值
func GetYearStatsNextMonth(nowTime time.Time) ([]response.CategoryStats, error) {
	// 1. 查出下月 CEA 的所有报价
	cli := db.Get()

	// now := time.Now()
	y, m, _ := nowTime.Date()
	nextMonth := m + 1
	nextYear := y
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}
	nextMonthStr := fmt.Sprintf("%d年%d月\n", nextYear, nextMonth)
	var quotes []model.YearQuotation
	if err := cli.
		Where("submitTime = ? AND product = ?", nextMonthStr, "CEA").
		Find(&quotes).Error; err != nil {
		return nil, err
	}

	// 2. 收集所有 UserID 并批量查 User 表
	uuids := make([]string, 0, len(quotes))
	for _, q := range quotes {
		uuids = append(uuids, q.Uuid)
	}
	var users []model.User
	if err := cli.
		Where("uuid IN ?", uuids).
		Find(&users).Error; err != nil {
		return nil, err
	}
	// 建一个 map[userId]User 方便后续查
	userMap := make(map[string]model.User, len(users))
	for _, u := range users {
		userMap[u.Uuid] = u
		log.Println("UserID", u.Uuid)
	}

	// 3. 在 Go 里合并数据并分组统计
	type key struct {
		t     int
		cType string
	}
	sums := make(map[key]struct {
		sumL, sumH float64
		count      int
	})

	for _, q := range quotes {
		// 价格字段是 string，要先转 float
		low, err1 := strconv.ParseFloat(q.LowerPrice, 64)
		high, err2 := strconv.ParseFloat(q.HigherPrice, 64)
		if err1 != nil || err2 != nil {
			continue
		}
		log.Println("Test_Uuid:", q.Uuid)
		u, ok := userMap[q.Uuid]
		if !ok {
			continue
		}
		k := key{t: u.Type, cType: u.CompanyType}
		log.Println("Test_key:", k)
		tmp := sums[k]
		tmp.sumL += low
		tmp.sumH += high
		tmp.count++
		sums[k] = tmp
	}

	// 4. 构造最终结果
	var stats []response.CategoryStats
	for k, v := range sums {
		avgL := v.sumL / float64(v.count)
		avgH := v.sumH / float64(v.count)
		stats = append(stats, response.CategoryStats{
			Type:        k.t,
			CompanyType: k.cType,
			AvgLower:    math.Round(avgL*100) / 100,
			AvgHigher:   math.Round(avgH*100) / 100,
			AvgBoth:     math.Round(((avgL+avgH)/2)*100) / 100,
		})
	}

	return stats, nil
}

func GetMonthlyCEAExpectation() ([]response.MonthlyPriceStats, error) {
	cli := db.Get()
	var expectation []model.CEAMonthExpectation
	err := cli.Find(&expectation).Error
	if err != nil {
		return nil, err
	}
	var results []response.MonthlyPriceStats
	for i := range expectation {
		result := response.MonthlyPriceStats{
			Month:     expectation[i].Date,
			AvgLow:    expectation[i].LowerPrice,
			AvgHigh:   expectation[i].HigherPrice,
			AvgPrice:  expectation[i].MidPrice,
			LowIndex:  expectation[i].LowerPriceIndex,
			HighIndex: expectation[i].HigherPriceIndex,
			MidIndex:  expectation[i].MidPriceIndex,
		}
		results = append(results, result)
	}

	resultFit := AddFitPriceToStats(results)
	return resultFit, nil
}

func GetMonthlyCCERExpectation() ([]response.MonthlyPriceStats, error) {
	cli := db.Get()
	var expectation []model.CCERMonthExpectation
	err := cli.Find(&expectation).Error
	if err != nil {
		return nil, err
	}
	var results []response.MonthlyPriceStats
	for i := range expectation {
		result := response.MonthlyPriceStats{
			Month:     expectation[i].Date,
			AvgLow:    expectation[i].LowerPrice,
			AvgHigh:   expectation[i].HigherPrice,
			AvgPrice:  expectation[i].MidPrice,
			LowIndex:  expectation[i].LowerPriceIndex,
			HighIndex: expectation[i].HigherPriceIndex,
			MidIndex:  expectation[i].MidPriceIndex,
		}
		results = append(results, result)
	}

	resultFit := AddFitPriceToStats(results)
	return resultFit, nil
}

func GetYearlyCEAExpectation() ([]response.MonthlyPriceStats, error) {
	cli := db.Get()
	var expectation []model.CEAYearExpectation
	err := cli.Find(&expectation).Error
	if err != nil {
		return nil, err
	}
	var results []response.MonthlyPriceStats
	for i := range expectation {
		result := response.MonthlyPriceStats{
			Month:     expectation[i].Date,
			AvgLow:    expectation[i].LowerPrice,
			AvgHigh:   expectation[i].HigherPrice,
			AvgPrice:  expectation[i].MidPrice,
			LowIndex:  expectation[i].LowerPriceIndex,
			HighIndex: expectation[i].HigherPriceIndex,
			MidIndex:  expectation[i].MidPriceIndex,
		}
		results = append(results, result)
	}

	resultFit := AddFitPriceToStats(results)
	return resultFit, nil
}
