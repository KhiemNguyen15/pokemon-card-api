package config

import "github.com/spf13/viper"

type Configurations struct {
	Database DatabaseConfigurations
}

type DatabaseConfigurations struct {
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
}

func LoadConfig(configPath string) (Configurations, error) {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return Configurations{}, err
	}

	var config Configurations

	err := viper.Unmarshal(&config)
	if err != nil {
		return Configurations{}, err
	}

	return config, nil
}
