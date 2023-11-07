package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"shopify-integration/pkg/authmiddleware"
	"shopify-integration/pkg/logger"
	"shopify-integration/pkg/model"
)

type AuthHandler struct {
	api *api
}

func NewAuthHandler(a *api) *AuthHandler {
	return &AuthHandler{
		api: a,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	user := &model.AuthUser{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		logger.Errorf("Login.ShouldBindJSON", err)
		c.JSON(http.StatusBadRequest, model.ErrInvalidBody)

		return
	}

	if !user.IsValid() {
		logger.Errorf("Login.Empty email or pass", err)
		c.JSON(http.StatusBadRequest, model.ErrInvalidBody)

		return
	}

	customerAccessToken, err := h.api.shopify.GetCustomerAccessToken(user.Email, user.Password)
	if err != nil {
		logger.Errorf("Login.GetCustomerAccessToken", err)
		c.JSON(http.StatusBadRequest, model.ErrInvalidBody)

		return
	}

	customerInfo, err := h.api.shopify.GetCustomerInfo(customerAccessToken)
	if err != nil {
		logger.Errorf("Login.GetCustomerAccessToken", err)
		c.JSON(http.StatusBadRequest, model.ErrUnhealthy)

		return
	}

	tokens, err := h.api.auth.CreateTokens(customerInfo.Data.Customer.Id)
	if err != nil {
		logger.Errorf("Login.CreateTokens", err)
		c.JSON(http.StatusBadRequest, model.ErrUnhealthy)

		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	oldTokens := authmiddleware.Tokens{}
	err := c.ShouldBindJSON(&oldTokens)
	if err != nil {
		logger.Errorf("Refresh.ShouldBindJSON", err)
		c.JSON(http.StatusBadRequest, model.ErrInvalidBody)

		return
	}

	newTokens, err := h.api.auth.Refresh(oldTokens)
	if err != nil {
		logger.Errorf("Refresh.Refresh", err)
		c.JSON(http.StatusUnauthorized, model.ErrUnauthorized)

		return
	}

	c.JSON(http.StatusOK, newTokens)
}
