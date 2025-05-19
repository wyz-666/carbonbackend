package service

import (
	"log"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"wxcloudrun-golang/app/handlers/response"
	"wxcloudrun-golang/db/model"

	"gonum.org/v1/gonum/stat"
)

func MonthlyAvg1(grouped map[string][]model.MonthQuotation, highIndex, lowIndex, midIndex float64) ([]response.MonthlyPriceStats, error) {
	var result []response.MonthlyPriceStats
	for month, items := range grouped {
		var sumLow, sumHigh, sumPrice float64
		var count int

		for _, q := range items {
			low, err1 := strconv.ParseFloat(q.LowerPrice, 64)
			high, err2 := strconv.ParseFloat(q.HigherPrice, 64)
			if err1 != nil || err2 != nil {
				continue // 忽略格式错误的记录
			}
			sumLow += low
			sumHigh += high
			price := (low + high) / 2
			sumPrice += price
			count++
		}

		if count > 0 {
			rate1 := sumLow / float64(count)
			rate2 := sumHigh / float64(count)
			rate3 := sumPrice / float64(count)
			rate4 := ((sumLow / float64(count)) / lowIndex) * 100
			rate5 := ((sumHigh / float64(count)) / highIndex) * 100
			rate6 := ((sumPrice / float64(count)) / midIndex) * 100
			result = append(result, response.MonthlyPriceStats{
				Month:     month,
				AvgLow:    math.Round(rate1*100) / 100,
				AvgHigh:   math.Round(rate2*100) / 100,
				AvgPrice:  math.Round(rate3*100) / 100,
				LowIndex:  math.Round(rate4*100) / 100,
				HighIndex: math.Round(rate5*100) / 100,
				MidIndex:  math.Round(rate6*100) / 100,
			})
		}
	}

	// 可选：按月份排序（假设格式为 2025年3月）
	sort.Slice(result, func(i, j int) bool {
		return result[i].Month < result[j].Month
	})
	return result, nil
}

func MonthlyAvg2(grouped map[string][]model.MonthQuotation) ([]response.GECMonthlyPriceStats, error) {
	var result []response.GECMonthlyPriceStats

	for key, items := range grouped {
		parts := strings.SplitN(key, "|", 2)
		if len(parts) != 2 {
			continue
		}
		tp, month := parts[0], parts[1]

		var sumPrice float64
		var count int

		for _, q := range items {
			price, err := strconv.ParseFloat(q.Price, 64)
			if err != nil {
				continue
			}
			sumPrice += price
			count++
		}

		if count > 0 {
			result = append(result, response.GECMonthlyPriceStats{
				Type:     tp,
				Month:    month,
				AvgPrice: sumPrice / float64(count),
			})
		}
	}

	// 按类型 + 月份排序（可选）
	sort.Slice(result, func(i, j int) bool {
		if result[i].Type == result[j].Type {
			return result[i].Month < result[j].Month
		}
		return result[i].Type < result[j].Type
	})
	return result, nil
}

func AddFitPriceToStats(data []response.MonthlyPriceStats) []response.MonthlyPriceStats {
	n := len(data)
	if len(data) < 2 {
		return data
	}
	// 按时间排序，保证 X 序列单调
	sort.Slice(data, func(i, j int) bool {
		return parseMonthNumber(data[i].Month) < parseMonthNumber(data[j].Month)
	})
	// 构建 X、Y 数组
	base := parseMonthNumber(data[0].Month)
	X := make([]float64, n)
	Y := make([]float64, n)
	for i, d := range data {
		X[i] = float64(parseMonthNumber(d.Month) - base)
		Y[i] = d.AvgPrice
		// log.Println("i:", i)
		// log.Println("X:", X[i])
		// log.Println("Y:", Y[i])
	}

	var intercept, slope float64
	intercept, slope = stat.LinearRegression(X, Y, nil, false)
	log.Println("a:", slope)
	log.Println("b:", intercept)
	// 填回 data.FitPrice，并做 NaN 检查 + 四舍五入两位
	for i := range data {
		y := slope*X[i] + intercept
		if math.IsNaN(y) {
			y = 0
		}
		data[i].FitPrice = math.Round(y*100) / 100
	}

	return data
}

var monthRe = regexp.MustCompile(`^([0-9]{4})年([0-9]{1,2})月`)

// parseMonthNumber 解析 "YYYY年M月" 返回自定义月份索引（年*12 + 月-1）
func parseMonthNumber(monthStr string) int {
	s := strings.TrimSpace(monthStr)
	m := monthRe.FindStringSubmatch(s)
	if len(m) != 3 {
		return 0
	}
	year, _ := strconv.Atoi(m[1])
	mon, _ := strconv.Atoi(m[2])
	// 按年*12 + (月-1) 计算连续月份索引
	return year*12 + (mon - 1)
}
