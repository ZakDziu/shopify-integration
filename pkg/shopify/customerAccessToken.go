package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CustomerAccessTokenResponse struct {
	Data struct {
		CustomerAccessTokenCreate struct {
			CustomerAccessToken struct {
				AccessToken string `json:"accessToken"`
			} `json:"customerAccessToken"`
			CustomerUserErrors []interface{} `json:"customerUserErrors"`
		} `json:"customerAccessTokenCreate"`
	} `json:"data"`
}

func (s *Shopify) GetCustomerAccessToken(email, password string) (string, error) {
	url := fmt.Sprintf("https://%v.myshopify.com/api/2023-01/graphql.json", s.StoreName)

	queryCustomerToken := fmt.Sprintf(`mutation {
		customerAccessTokenCreate(input: {
			email: "%v",
			password: "%v"
		}) {
			customerAccessToken {
				accessToken
			}
			customerUserErrors {
				message
			}
		}
	}`, email, password)

	request := map[string]string{
		"query": queryCustomerToken,
	}

	reqB, err := json.Marshal(request)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqB))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Shopify-Storefront-Access-Token", s.StorefrontAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response CustomerAccessTokenResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response.Data.CustomerAccessTokenCreate.CustomerAccessToken.AccessToken, nil
}
