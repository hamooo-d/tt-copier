package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	BanksNames    map[string]string   `mapstructure:"banksNames"`
	Database      DatabaseConfig      `mapstructure:"database"`
	Dests         DestsConfig         `mapstructure:"dests"`
	LogPath       string              `mapstructure:"log_path"`
	Env           string              `mapstructure:"env"`
	AfterDate     string              `mapstructure:"after_date"`
	SFTP          SFTPConfig          `mapstructure:"sftp"`
	FilesPrefixes FilesPrefixesConfig `mapstructure:"files_prefixes"`
	SourceList    []string            `mapstructure:"source_list"`
}

type SFTPConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	DBPath string `mapstructure:"db_path"`
}

type DestsConfig struct {
	BankDest string `mapstructure:"bank_dest"`
}

type FilesPrefixesConfig struct {
	BankFilesPrefixes []string `mapstructure:"bankFilesPrefixes"`
	TTFilesPrefixes   []string `mapstructure:"TTFilesPrefixes"`
}

func LoadConfig(configPath string) (*Config, error) {
	var config Config

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
