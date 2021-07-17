package config

import (
	"github.com/spf13/viper"
)

//Config struct using viper
type Config struct {
	v *viper.Viper
}

//New Create a new config
func New() *Config {
	c := Config{
		v: viper.New(),
	}
	c.v.SetDefault("API_PORT", 8080)

	c.v.SetEnvPrefix("")
	c.v.AutomaticEnv()

	return &c
}

//GetAPIPort gets the main API port
func (c *Config) APIPort() int {
	return c.v.GetInt("API_PORT")
}
