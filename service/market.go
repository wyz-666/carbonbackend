package service

import (
	"carbonbackend/app/handlers/request"
	"carbonbackend/app/handlers/response"
	"carbonbackend/db"
	"carbonbackend/db/model"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func MarketSubmit(m *request.ReqMarket, s string) error {
	cli := db.Get()
	if s == "CEA" {
		ceaMarket := model.CEAMarket{
			Date:         m.Date,
			LowerPrice:   m.LowerPrice,
			HigherPrice:  m.HigherPrice,
			ClosingPrice: m.ClosingPrice,
		}
		err := cli.Create(&ceaMarket).Error
		if err != nil {
			log.Printf("[ERROR] Submit cea market error: %v", err)
			return err
		}
	} else {
		ccerMarket := model.CCERMarket{
			Date:         m.Date,
			ClosingPrice: m.ClosingPrice,
		}
		err := cli.Create(&ccerMarket).Error
		if err != nil {
			log.Printf("[ERROR] Submit ccer market error: %v", err)
			return err
		}
	}
	return nil
}

func GECStatsSubmit(e *request.ReqGECExpectation) error {
	cli := db.Get()
	res := model.GECMonthExpectation{
		Product:    e.Product,
		Type:       e.Type,
		Date:       e.Date,
		Price:      e.Price,
		PriceIndex: e.PriceIndex,
	}
	err := cli.Create(&res).Error
	if err != nil {
		log.Printf("[ERROR] GEC Submit stats error: %v", err)
		return err
	}
	return nil
}

func StatsSubmit(e *request.ReqExpectation) error {
	cli := db.Get()
	if e.Product == "CEA" {
		if e.Type == "month" {
			e := model.CEAMonthExpectation{
				Date:             e.Date,
				LowerPrice:       e.LowerPrice,
				HigherPrice:      e.HigherPrice,
				MidPrice:         e.MidPrice,
				LowerPriceIndex:  math.Round((e.LowerPrice/40.00)*10000) / 100,
				HigherPriceIndex: math.Round((e.HigherPrice/44.32)*10000) / 100,
				MidPriceIndex:    math.Round((e.MidPrice/42.16)*10000) / 100,
			}
			err := cli.Create(&e).Error
			if err != nil {
				log.Printf("[ERROR] Submit stats error: %v", err)
				return err
			}
		} else {
			e := model.CEAYearExpectation{
				Date:             e.Date,
				LowerPrice:       e.LowerPrice,
				HigherPrice:      e.HigherPrice,
				MidPrice:         e.MidPrice,
				LowerPriceIndex:  math.Round((e.LowerPrice/53.45)*10000) / 100,
				HigherPriceIndex: math.Round((e.HigherPrice/58.26)*10000) / 100,
				MidPriceIndex:    math.Round((e.MidPrice/55.86)*10000) / 100,
			}
			err := cli.Create(&e).Error
			if err != nil {
				log.Printf("[ERROR] Submit stats error: %v", err)
				return err
			}
		}
	} else {
		e := model.CCERMonthExpectation{
			Date:             e.Date,
			LowerPrice:       e.LowerPrice,
			HigherPrice:      e.HigherPrice,
			MidPrice:         e.MidPrice,
			LowerPriceIndex:  math.Round((e.LowerPrice/39.78)*10000) / 100,
			HigherPriceIndex: math.Round((e.HigherPrice/41.57)*10000) / 100,
			MidPriceIndex:    math.Round((e.MidPrice/40.68)*10000) / 100,
		}
		err := cli.Create(&e).Error
		if err != nil {
			log.Printf("[ERROR] Submit stats error: %v", err)
			return err
		}
	}
	return nil
}

func GetCEAMonthScoreList(nowTime time.Time) ([]response.ResUserScore, error) {
	// 1. 查出下月 CEA 的所有报价
	cli := db.Get()
	var quotes []model.MonthQuotation
	// now := time.Now()
	y, m, _ := nowTime.Date()
	month := m
	year := y
	monthStr := fmt.Sprintf("%d年%d月\n", year, month)
	if err := cli.
		Where("applicableTime = ? AND product = ?", monthStr, "CEA").
		Find(&quotes).Error; err != nil {
		return nil, err
	}
	priceList, err := getCEAPriceList(cli, monthStr)
	if err != nil {
		return nil, err
	}
	result := make([]response.ResUserScore, 0)
	for i := range quotes {
		low, err := strconv.ParseFloat(quotes[i].LowerPrice, 64)
		if err != nil {
			return nil, err
		}
		high, err := strconv.ParseFloat(quotes[i].HigherPrice, 64)
		if err != nil {
			return nil, err
		}
		var user model.User
		if err := cli.
			Where("uuid = ?", quotes[i].Uuid).
			First(&user).Error; err != nil {
			return nil, err
		}
		log.Println("user:", user.UserID)
		mid := math.Round(((low+high)/2)*100) / 100
		var countScore int
		var maxDist float64
		maxDist = 0
		countScore = 0
		for j := range priceList {
			// lowToday, err := strconv.ParseFloat(priceList[j].LowerPrice, 64)
			// if err != nil {
			// 	return nil, err
			// }
			// if (lowToday > low) && (lowToday < high) {
			// 	countScore++
			// 	log.Println("lowToday:", lowToday)
			// }
			// highToday, err := strconv.ParseFloat(priceList[j].HigherPrice, 64)
			// if err != nil {
			// 	return nil, err
			// }
			// if (highToday > low) && (highToday < high) {
			// 	countScore++
			// 	log.Println("highToday:", highToday)
			// }
			closingToday, err := strconv.ParseFloat(priceList[j].ClosingPrice, 64)
			if err != nil {
				return nil, err
			}
			if (closingToday < low) || (closingToday > high) {
				countScore++
				log.Println("closingToday:", closingToday)
			}
			if math.Abs(mid-closingToday) > maxDist {
				maxDist = math.Abs(mid - closingToday)
			}
			// if math.Abs(highToday-high) > maxDist {
			// 	maxDist = math.Abs(highToday - high)
			// }
		}
		log.Println("count:", countScore)
		userScore := response.ResUserScore{
			Uuid:        user.Uuid,
			UserID:      user.UserID,
			CompanyName: user.CompanyName,
			UserName:    user.UserName,
			Phone:       user.Phone,
			Email:       user.Email,
			Score:       float64(countScore) + math.Round((maxDist/mid)*100)/100,
		}
		result = append(result, userScore)
	}
	return result, nil
}

func GetCCERMonthScoreList(nowTime time.Time) ([]response.ResUserScore, error) {
	// 1. 查出下月 CCER 的所有报价
	cli := db.Get()
	var quotes []model.MonthQuotation
	// now := time.Now()
	y, m, _ := nowTime.Date()
	month := m
	year := y
	monthStr := fmt.Sprintf("%d年%d月\n", year, month)
	if err := cli.
		Where("applicableTime = ? AND product = ?", monthStr, "CCER").
		Find(&quotes).Error; err != nil {
		return nil, err
	}
	priceList, err := getCCERPriceList(cli, monthStr)
	if err != nil {
		return nil, err
	}
	result := make([]response.ResUserScore, 0)
	for i := range quotes {
		low, err := strconv.ParseFloat(quotes[i].LowerPrice, 64)
		if err != nil {
			return nil, err
		}
		high, err := strconv.ParseFloat(quotes[i].HigherPrice, 64)
		if err != nil {
			return nil, err
		}
		var user model.User
		if err := cli.
			Where("uuid = ?", quotes[i].Uuid).
			First(&user).Error; err != nil {
			return nil, err
		}
		log.Println("user:", user.UserID)
		mid := math.Round(((low+high)/2)*100) / 100
		var countScore int
		var maxDist float64
		maxDist = 0
		countScore = 0
		for j := range priceList {
			closingToday, err := strconv.ParseFloat(priceList[j].ClosingPrice, 64)
			if err != nil {
				return nil, err
			}
			if (closingToday < low) || (closingToday > high) {
				countScore++
				log.Println("closingToday:", closingToday)
			}
			if math.Abs(mid-closingToday) > maxDist {
				maxDist = math.Abs(mid - closingToday)
			}
		}
		log.Println("count:", countScore)
		userScore := response.ResUserScore{
			Uuid:        user.Uuid,
			UserID:      user.UserID,
			CompanyName: user.CompanyName,
			UserName:    user.UserName,
			Phone:       user.Phone,
			Email:       user.Email,
			Score:       float64(countScore) + math.Round((maxDist/mid)*100)/100,
		}
		result = append(result, userScore)
	}
	return result, nil
}

// 得到当月价格数据
func getCEAPriceList(db *gorm.DB, time string) ([]model.CEAMarket, error) {
	monthTrim := strings.TrimSpace(time) // → "2025年4月"

	// 2. 拆成 年 和 月
	parts := strings.SplitN(monthTrim, "年", 2)     // ["2025", "4月"]
	yearStr := parts[0]                            // "2025"
	monthPart := strings.TrimSuffix(parts[1], "月") // "4"
	monInt, err := strconv.Atoi(monthPart)
	if err != nil {
		return nil, err
	}
	// 3. 构造 “YYYY/MM/” 前缀
	prefix := fmt.Sprintf("%s/%02d/", yearStr, monInt) // "2025/04/"

	// 4A. 简单 LIKE 查询
	var markets []model.CEAMarket
	if err := db.Where("`date` LIKE ?", prefix+"%").
		Find(&markets).Error; err != nil {
		return nil, err
	}
	return markets, nil
}

func getCCERPriceList(db *gorm.DB, time string) ([]model.CCERMarket, error) {
	monthTrim := strings.TrimSpace(time) // → "2025年4月"

	// 2. 拆成 年 和 月
	parts := strings.SplitN(monthTrim, "年", 2)     // ["2025", "4月"]
	yearStr := parts[0]                            // "2025"
	monthPart := strings.TrimSuffix(parts[1], "月") // "4"
	monInt, err := strconv.Atoi(monthPart)
	if err != nil {
		return nil, err
	}
	// 3. 构造 “YYYY/MM/” 前缀
	prefix := fmt.Sprintf("%s/%02d/", yearStr, monInt) // "2025/04/"

	// 4A. 简单 LIKE 查询
	var markets []model.CCERMarket
	if err := db.Where("`date` LIKE ?", prefix+"%").
		Find(&markets).Error; err != nil {
		return nil, err
	}
	return markets, nil
}
