package main

import (
	"bufio"
	"log"
	"os"
	"time"

	"bytes"
	"fmt"

	"github.com/cleanshaven/gologging"
	"github.com/cleanshaven/noaaweather"
	Notification "github.com/cleanshaven/notification"
	"github.com/cleanshaven/persistantweatherblock/config"
	"github.com/cleanshaven/persistantweatherblock/weathertemplates"
	"github.com/spf13/viper"
)

type ForecastToOutput struct {
	Forecast noaaweather.PeriodJson
	Alert    noaaweather.AlertFeatureJson
}

var currentWeather ForecastToOutput

func main() {
	gologging.SetupSyslog("persistantweatherblock")

	log.Println("starting up weather block")
	configChangedChannel := make(chan struct{})
	defer close(configChangedChannel)
	err := config.SetupConfig(configChangedChannel)
	if err != nil {
		log.Println(err)
		return
	}
	weathertemplates.Initialize()

	ticker := time.NewTicker(time.Duration(config.MyConfig.Block.Interval) * time.Minute)
	mouseButtonPressedChannel := make(chan struct{})
	defer close(mouseButtonPressedChannel)
	go func() {

		reader := bufio.NewReader(os.Stdin)
		for {
			reader.ReadString('\n')
			mouseButtonPressedChannel <- struct{}{}
		}
	}()
	refreshWeather()
	for {
		select {
		case <-configChangedChannel:
			configChanged()
			ticker.Reset(time.Duration(config.MyConfig.Block.Interval) * time.Minute)

		case <-ticker.C:
			refreshWeather()

		case <-mouseButtonPressedChannel:
			displayDetail()
		}
	}

}

func configChanged() error {
	err := config.SetConfigFile()
	if err != nil {
		log.Println(err)
		return err
	}

	err = viper.Unmarshal(&config.MyConfig)
	if err != nil {
		log.Println(err)
		return err
	}
	weathertemplates.ParseTemplates()
	refreshWeather()
	return nil

}

func refreshWeather() error {
	forecast, alerts, err := noaaweather.GetWeather(config.MyConfig.Location.Latitude,
		config.MyConfig.Location.Longitude)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	currentWeather.Forecast = forecast.Properties.Periods[0]
	if len(alerts.Features) > 0 {
		currentWeather.Alert = alerts.Features[0]
	}
	var tpl bytes.Buffer

	err = weathertemplates.WeatherBarTemplate.Execute(&tpl, currentWeather)
	fmt.Print(tpl.String())
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func getIcon(url string) (filename string, err error) {
	filename = ""
	err = nil

	filename, err = noaaweather.GetIcon(url)
	return
}

func displayDetail() {
	err := refreshWeather()
	if err != nil {
		return
	}

	var tpl bytes.Buffer
	err = weathertemplates.DetailTemplate.Execute(&tpl, currentWeather)

	iconFile, err := getIcon(currentWeather.Forecast.Icon)
	alertNotification(config.MyConfig.DetailPopup.Title, tpl.String(), iconFile, config.MyConfig.DetailPopup.TimeToShow)

}

func alertNotification(summary, description, icon string, expireTime int) error {
	return Notification.Notify("I3Weather", summary, description, icon, expireTime)
}
