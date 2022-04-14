package api

import (
	"log"
	"net/http"
	"sales-api/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	headerKey      = "authorization"
	headerType     = "bearer"
	authPayloadKey = "authoization_payload"
)

func authMiddleware(tm token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(headerKey)
		if authHeader == "" {
			log.Println("[ERR] Authorization header is absent")

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, genericResponse{
				Success: false,
				Message: "Unauthorized",
			})
			return
		}

		f := strings.Fields(authHeader)
		if len(f) < 2 {
			log.Println("[ERR] Authorization header format is invalid")

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, genericResponse{
				Success: false,
				Message: "Unauthorized",
			})
			return
		}

		authHeaderType := strings.ToLower(f[0])
		if !strings.EqualFold(authHeaderType, headerType) {
			log.Printf("[ERR] Unsupported authorization header type %s", authHeaderType)

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, genericResponse{
				Success: false,
				Message: "Unauthorized",
			})
			return
		}

		accessToken := f[1]
		payload, err := tm.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, genericResponse{
				Success: false,
				Message: "Unauthorized",
			})
			return
		}

		ctx.Set(authPayloadKey, payload)
		ctx.Next()
	}
}
