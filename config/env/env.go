package env

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type appConfig struct {
	// database information
	MongoURI string `mapstructure:"mongo_uri"`
	DBName   string `mapstructure:"db_name"`
	Port     int64  `mapstructure:"port"`

	// jwt key information
	AccessTokenKey        string        `mapstructure:"access_token_key"`
	AccessTokenExpiredIn  time.Duration `mapstructure:"access_token_expired_in"`
	AccessTokenMaxAge     int64         `mapstructure:"access_token_max_age"`
	RefreshTokenKey       string        `mapstructure:"refresh_token_key"`
	RefreshTokenExpiredIn time.Duration `mapstructure:"refresh_token_expired_in"`
	RefreshTokenMaxAge    int64         `mapstructure:"refresh_token_max_age"`

	// redis information
	ClientOrigin string `mapstructure:"client_origin"`
	RedisUrl     string `mapstructure:"redis_url"`

	// google client info
	GoogleOauthClientID    string `mapstructure:"google_oauth_client_id"`
	GoogleOauthSecret      string `mapstructure:"google_oauth_secret"`
	GoogleOauthRedirectUrl string `mapstructure:"google_oauth_redirect_url"`

	// AWS Service
	AWSRegion      string `mapstructure:"aws_region"`
	AWSAccessKeyID string `mapstructure:"aws_access_key_id"`
	AWSSecretKey   string `mapstructure:"aws_secret_key"`
	AWSBucketName  string `mapstructure:"aws_bucket_name"`
}

var (
	AppConfig *appConfig
)

func LoadConfig(configPath string) error {
	v := viper.New()

	if configPath == "" {
		return fmt.Errorf("failed to load config file. please set config")
	}

	v.AddConfigPath(configPath)
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	if err := v.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("failed to unmarshall app config: %v", err)
	}

	return nil
}
