package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	defaultLatitude       = "41.26"
	defaultLongitude      = "-81.86"
	defaultBlockTemplate  = `{{.Forecast.Temperature}}°F {{.Forecast.ShortForecast}} {{.Forecast.ProbabilityOfPrecipitation.Value}}% <span foreground="red">{{.Alert.Properties.Event}}</span>`
	defaultTimeToShow     = 30
	defaultAlertTemplate  = `{{.Alert.Properties.Event}}`
	defaultDetailTemplate = `{{.Alert.Properties.Description}}`
	defaultAlertTitle     = "Weather Alert"
	defaultDetailTitle    = "Weather Details"

	latitudeVariable              = "location.latitude"
	longitudeVariable             = "location.longitude"
	blockTemplateVariable         = "block.template"
	alertPopupTemplateVariable    = "alertPopup.template"
	alertPopupTimeToShowVariable  = "alertPopup.timeToShow"
	detailPopupTemplateVariable   = "detailPopup.template"
	detailPopupTimeToShowVariable = "detailPopup.timeToShow"
	buttonVariable                = "BLOCK_BUTTON"
	defaultBlockButton            = "-1"
	alertTitleVariable            = "alertPopup.title"
	detailTitleVariable           = "detailPopup.title"
	intervalVariable              = "block.interval"
	defaultInterval               = 30

	configName = "config.yml"
)

type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type Block struct {
	Template string `json:"template"`
	Interval int    `json:"interval"`
	Button   int    `json:"button"`
}

type Popup struct {
	Template   string `json:"template"`
	TimeToShow int    `json:"timeToShow"`
	Title      string `json:"title"`
}

type Config struct {
	CheckInterval int      `json:"interval"`
	Location      Location `json:"location"`
	Block         Block    `json:"block"`
	AlertPopup    Popup    `json:"alertPopup"`
	DetailPopup   Popup    `json:"detailPopup"`
}

var MyConfig Config

var ConfigLocation string = `$HOME/.config/i3persistantweatherblock`

func SetupConfig(changeChannel chan<- struct{}) error {
	setDefaults()
	setEnv()

	err := SetConfigFile()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&MyConfig)
	if err != nil {
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		changeChannel <- struct{}{}
	})
	return nil
}

func setDefaults() {
	viper.SetDefault(latitudeVariable, defaultLatitude)
	viper.SetDefault(longitudeVariable, defaultLongitude)
	viper.SetDefault(blockTemplateVariable, defaultBlockTemplate)
	viper.SetDefault(alertPopupTemplateVariable, defaultAlertTemplate)
	viper.SetDefault(alertPopupTimeToShowVariable, defaultTimeToShow)
	viper.SetDefault(detailPopupTemplateVariable, defaultDetailTemplate)
	viper.SetDefault(detailPopupTimeToShowVariable, defaultTimeToShow)
	viper.SetDefault(buttonVariable, defaultBlockButton)
	viper.SetDefault(alertTitleVariable, defaultAlertTitle)
	viper.SetDefault(detailTitleVariable, defaultDetailTitle)
	viper.SetDefault(intervalVariable, defaultInterval)
}

func setEnv() {
	viper.AutomaticEnv()
}

func SetConfigFile() error {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ConfigLocation)
	err := viper.ReadInConfig()
	return err
}

func IsButtonPressed() bool {
	return viper.GetString(buttonVariable) != "-1"
}
