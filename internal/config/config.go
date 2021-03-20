package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type LevelGoSetting struct {
	URI string
}
type Setting struct {
	LevelGo LevelGoSetting
}

func Read() *Setting {
	var settings Setting
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.SetConfigType("json")
	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	if err := v.Unmarshal(&settings); err != nil {
		fmt.Println(err)
	}
	return &settings
}
