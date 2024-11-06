package main

import (
	"log"
	"strings"

	"github.com/rgaiffe/rss-parser/internal/app/parser"
	"github.com/spf13/viper"
)

// Load config file locate on ./config/config.yaml
func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./configs")

	replacer := strings.NewReplacer(".", "_", "-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AllowEmptyEnv(true)

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := parser.StartParser(); err != nil {
		log.Fatal(err)
	}
}
