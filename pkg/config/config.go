package config

import (
	"github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload" // By design
)

type Configs struct {
	Server  ServerConfig
	Keys    Path
	Shopify ShopifyConfig
}

type Path struct {
	AccessKey  string `env:"HASH_KEY_ACCESS"`
	RefreshKey string `env:"HASH_KEY_REFRESH"`
}

type ServerConfig struct {
	ServerPort  string   `env:"SERVER_PORT"`
	ReadTimeout Duration `env:"READ_TIMEOUT"`
}

type ShopifyConfig struct {
	StoreFrontAccessToken string `env:"STOREFRONT_ACCESS_TOKEN"`
	StoreName             string `env:"STORE_NAME"`
}

func New() (*Configs, error) {
	var config Configs
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
