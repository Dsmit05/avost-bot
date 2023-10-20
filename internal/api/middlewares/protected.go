package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Dsmit05/avost-bot/internal/api/cryptography"
	"github.com/Dsmit05/avost-bot/internal/logger"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

var (
	ErrUserRole = errors.New("the user has a negative role")
)

type ProtectedMiddleware struct {
	auth cryptography.AccessToken
}

func NewProtectedMiddleware(auth cryptography.AccessToken) *ProtectedMiddleware {
	return &ProtectedMiddleware{auth: auth}
}

func (o *ProtectedMiddleware) AuthMiddleware(c *gin.Context) {
	email, role, err := o.parseAuthHeader(c)
	if err != nil {
		c.AbortWithError(http.StatusExpectationFailed, err) //nolint:errcheck
		return
	}

	logger.Infof("AuthMiddleware", zap.String("email", email), zap.String("role", role), zap.String("path", c.FullPath()))
}

func (o *ProtectedMiddleware) parseAuthHeader(c *gin.Context) (email, role string, err error) {
	header := c.GetHeader("Authorizations")
	if header == "" || len(header) > 250 {
		return "", "", fmt.Errorf("bad header")
	}

	return o.auth.ParseToken(header)
}

func CheckAccessRights(c *gin.Context, roles ...string) error {
	roleFromContext, ok := c.Get("role")
	if !ok {
		return ErrUserRole
	}

	for _, role := range roles {
		if roleFromContext == role {
			return nil
		}
	}

	return ErrUserRole
}
