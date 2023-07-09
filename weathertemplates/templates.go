package weathertemplates

import (
	"log"
	"text/template"

	"github.com/cleanshaven/persistantweatherblock/config"
)

var WeatherBarTemplate, DetailTemplate *template.Template

func Initialize() {
	WeatherBarTemplate = template.New("WeatherBar")
	DetailTemplate = template.New("Detail")
	ParseTemplates()
}

func ParseTemplates() error {
	_, err := WeatherBarTemplate.Parse(config.MyConfig.Block.Template)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = DetailTemplate.Parse(config.MyConfig.DetailPopup.Template)
	if err != nil {
		log.Println(err)
	}
	return err
}
