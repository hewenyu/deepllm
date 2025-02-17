package weather

import (
	"context"
	"time"

	"github.com/hewenyu/deepllm/internal/data"
)

// WeatherAgent specializes in weather-based activity recommendations
type WeatherAgent struct {
	store *data.Store
}

// NewWeatherAgent creates a new weather advisor agent
func NewWeatherAgent(store *data.Store) *WeatherAgent {
	return &WeatherAgent{
		store: store,
	}
}

// WeatherAdvice contains weather-based recommendations
type WeatherAdvice struct {
	Weather        *data.DailyForecast `json:"weather"`
	Suitable       []string            `json:"suitable_activities"`   // 适合的活动
	Unsuitable     []string            `json:"unsuitable_activities"` // 不适合的活动
	Precautions    []string            `json:"precautions"`           // 注意事项
	IndoorOptions  []string            `json:"indoor_options"`        // 室内备选
	OutdoorOptions []string            `json:"outdoor_options"`       // 室外备选
}

// GetAdvice provides weather-based activity recommendations
func (a *WeatherAgent) GetAdvice(ctx context.Context, date time.Time) (*WeatherAdvice, error) {
	forecast := a.store.GetWeatherForecast()
	if forecast == nil || len(forecast.DailyForecasts) == 0 {
		return nil, nil
	}

	// Find matching forecast
	var todayForecast *data.DailyForecast
	for _, f := range forecast.DailyForecasts {
		if f.Date == date.Format("2006-01-02") {
			todayForecast = &f
			break
		}
	}

	if todayForecast == nil {
		return nil, nil
	}

	advice := &WeatherAdvice{
		Weather: todayForecast,
	}

	// Generate activity recommendations based on weather conditions
	a.generateActivityRecommendations(todayForecast, advice)

	return advice, nil
}

// generateActivityRecommendations generates activity recommendations based on weather
func (a *WeatherAgent) generateActivityRecommendations(forecast *data.DailyForecast, advice *WeatherAdvice) {
	// Check weather conditions
	isRainy := contains([]string{"小雨", "中雨", "大雨"}, forecast.Weather.Day)
	isSunny := contains([]string{"晴", "多云"}, forecast.Weather.Day)
	isHot := forecast.Temperature.Max > 30
	isCold := forecast.Temperature.Min < 10
	isWindy := false
	if speed := forecast.Wind.Speed; speed.Max > 30 {
		isWindy = true
	}
	isPoorAirQuality := forecast.AirQuality.AQI > 150

	// Generate recommendations
	if isRainy {
		advice.Unsuitable = append(advice.Unsuitable,
			"户外徒步",
			"西湖游船",
			"户外摄影",
		)
		advice.Precautions = append(advice.Precautions,
			"携带雨具",
			"注意路滑",
			"避免湿鞋",
		)
		advice.IndoorOptions = append(advice.IndoorOptions,
			"博物馆参观",
			"茶馆品茶",
			"室内购物",
		)
	}

	if isSunny && !isHot {
		advice.Suitable = append(advice.Suitable,
			"西湖游船",
			"灵隐寺参观",
			"龙井茶园",
		)
		advice.Precautions = append(advice.Precautions,
			"做好防晒",
			"补充水分",
		)
		advice.OutdoorOptions = append(advice.OutdoorOptions,
			"徒步西湖",
			"城市观光",
			"公园漫步",
		)
	}

	if isHot {
		advice.Unsuitable = append(advice.Unsuitable,
			"长时间户外活动",
			"徒步登山",
		)
		advice.Precautions = append(advice.Precautions,
			"避免中午外出",
			"防暑降温",
			"及时补水",
		)
		advice.IndoorOptions = append(advice.IndoorOptions,
			"博物馆",
			"商场购物",
			"室内娱乐",
		)
	}

	if isCold {
		advice.Precautions = append(advice.Precautions,
			"注意保暖",
			"准备厚衣物",
		)
		advice.IndoorOptions = append(advice.IndoorOptions,
			"温泉体验",
			"室内景点",
		)
	}

	if isWindy {
		advice.Unsuitable = append(advice.Unsuitable,
			"登高望远",
			"游船活动",
		)
		advice.Precautions = append(advice.Precautions,
			"注意防风",
			"避免易飞物品",
		)
	}

	if isPoorAirQuality {
		advice.Unsuitable = append(advice.Unsuitable,
			"户外运动",
			"长时间户外活动",
		)
		advice.Precautions = append(advice.Precautions,
			"佩戴口罩",
			"减少户外活动",
		)
		advice.IndoorOptions = append(advice.IndoorOptions,
			"室内活动为主",
			"避免剧烈运动",
		)
	}

	// Always provide some indoor options
	if len(advice.IndoorOptions) == 0 {
		advice.IndoorOptions = []string{
			"博物馆参观",
			"茶馆品茶",
			"特色餐厅",
			"商场购物",
		}
	}

	// Always provide some outdoor options if weather permits
	if len(advice.OutdoorOptions) == 0 && !isRainy && !isPoorAirQuality && !isWindy {
		advice.OutdoorOptions = []string{
			"西湖景区游览",
			"灵隐寺参观",
			"城市观光",
		}
	}
}

// Helper function to check if a string is in a slice
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
