package config

import "github.com/spf13/viper"

type Configurations struct {
	Database   DatabaseConfigurations
	PokemonAPI PokemonAPIConfigurations
	DiscordAPI DiscordAPIConfigurations
}

type DatabaseConfigurations struct {
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	CardCount  int
}

type PokemonAPIConfigurations struct {
	APIKey    string
	PageCount int
}

type DiscordAPIConfigurations struct {
	APIKey string
}

func LoadConfig(configPath string) (Configurations, error) {
	var config Configurations

	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
