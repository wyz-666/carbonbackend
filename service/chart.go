package service

import (
	"wxcloudrun-golang/app/handlers/response"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

func GetMonthlyCEAStats() ([]response.MonthlyPriceStats, error) {
	cli := db.Get()

	var data []model.MonthQuotation
	err := cli.Where("product = ?", "CEA").Find(&data).Error
	if err != nil {
		return nil, err
	}

	// 临时按月份分组
	grouped := make(map[string][]model.MonthQuotation)
	for _, item := range data {
		grouped[item.ApplicableTime] = append(grouped[item.ApplicableTime], item)
	}
	result, err := MonthlyAvg1(grouped, 44.32, 40.00, 40.68)
	if err != nil {
		return nil, err
	}
	resultFit := AddFitPriceToStats(result)
	return resultFit, nil
}

func GetMonthlyCCERStats() ([]response.MonthlyPriceStats, error) {
	cli := db.Get()

	var data []model.MonthQuotation
	err := cli.Where("product = ?", "CCER").Find(&data).Error
	if err != nil {
		return nil, err
	}

	// 临时按月份分组
	grouped := make(map[string][]model.MonthQuotation)
	for _, item := range data {
		grouped[item.ApplicableTime] = append(grouped[item.ApplicableTime], item)
	}
	result, err := MonthlyAvg1(grouped, 41.57, 39.78, 40.68)

	if err != nil {
		return nil, err
	}
	resultFit := AddFitPriceToStats(result)
	return resultFit, nil
}

func GetGECMonthlyStatsByType() ([]response.GECMonthlyPriceStats, error) {
	cli := db.Get()
	var data []model.MonthQuotation
	err := cli.Where("product = ?", "GEC").Find(&data).Error
	if err != nil {
		return nil, err
	}

	// 分组：type + month 作为键
	grouped := make(map[string][]model.MonthQuotation)
	for _, item := range data {
		key := item.Type + "|" + item.ApplicableTime
		grouped[key] = append(grouped[key], item)
	}
	result, err := MonthlyAvg2(grouped)
	if err != nil {
		return nil, err
	}
	return result, nil
}
