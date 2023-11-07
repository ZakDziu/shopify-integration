package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CustomerInfoResponse struct {
	Data struct {
		Customer struct {
			Id               string      `json:"id"`
			FirstName        string      `json:"firstName"`
			LastName         string      `json:"lastName"`
			AcceptsMarketing bool        `json:"acceptsMarketing"`
			Email            string      `json:"email"`
			Phone            interface{} `json:"phone"`
		} `json:"customer"`
	} `json:"data"`
}

func (s *Shopify) GetCustomerInfo(customerAccessToken string) (*CustomerInfoResponse, error) {
	url := fmt.Sprintf("https://%v.myshopify.com/api/2023-01/graphql.json", s.StoreName)

	queryCustomerInfo := fmt.Sprintf(`query { customer(customerAccessToken: "%v") { id firstName lastName acceptsMarketing email phone } }`, customerAccessToken)

	request := map[string]string{
		"query": queryCustomerInfo,
	}

	reqB, err := json.Marshal(request)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqB))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Shopify-Storefront-Access-Token", s.StorefrontAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response *CustomerInfoResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
