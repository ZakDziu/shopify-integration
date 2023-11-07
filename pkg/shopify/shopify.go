package shopify

import "shopify-integration/pkg/config"

type Shopify struct {
	StoreName             string
	StorefrontAccessToken string
}

type Interface interface {
	GetCustomerAccessToken(email, password string) (string, error)
	GetCustomerInfo(customerAccessToken string) (*CustomerInfoResponse, error)
}

func NewShopify(config config.ShopifyConfig) *Shopify {
	return &Shopify{
		StoreName:             config.StoreName,
		StorefrontAccessToken: config.StoreFrontAccessToken,
	}
}
