package config

import "github.com/spf13/viper"

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSourceName  string `mapstructure:"DB_SOURCE_NAME"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	var c Config

	err := viper.ReadInConfig()
	if err != nil {
		return c, err
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return c, err
	}

	return c, nil
}
