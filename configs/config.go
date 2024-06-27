package configs

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type conf struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
	JwtSecret     string `mapstructure:"JWT_SECRET"`
	JwtExpiresIn  int    `mapstructure:"JWT_EXPIRES_IN"`
	TokenAuth     *jwtauth.JWTAuth
}

func LoadConfig(configFilePath string) (*conf, error) {
	var config *conf

	viper.SetConfigFile("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(configFilePath)
	viper.SetConfigFile("configs/.env")
	// allow override .env file with system environment variables
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	config.TokenAuth = jwtauth.New("HS256", []byte(config.JwtSecret), nil)
	return config, err
}
