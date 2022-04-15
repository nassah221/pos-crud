package api

import (
	"fmt"
	db "sales-api/db/sqlc"
	"sales-api/token"
	"time"

	"github.com/gin-gonic/gin"
)

// todo: token secret and duration should be read in from the env
const (
	tokenSecret   = "01234567890123456789012345678901"
	tokenDuration = time.Hour * 10
)

type Server struct {
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(store *db.Store) (*Server, error) {
	tm, err := token.NewJWTMaker(tokenSecret)
	if err != nil {
		return nil, fmt.Errorf("unable to create token maker: %v", err)
	}

	s := &Server{
		store:      store,
		tokenMaker: tm,
	}
	s.setupRoutes()

	return s, nil
}

func (s *Server) setupRoutes() {
	r := gin.Default()

	r.GET("/cashiers/:id", s.GetCashier)
	r.GET("/cashiers", s.ListCashiers)
	r.GET("/cashiers/:id/passcode", s.GetPassword)

	r.POST("/cashiers", s.CreateCashier)
	r.POST("/cashiers/:id/login", s.LoginUser)
	r.POST("/products", s.CreateProduct)
	r.POST("/categories", s.CreateCategory)
	r.POST("/payments", s.CreatePayment)

	r.PUT("/cashiers/:id", s.UpdateCashier)
	r.PUT("/products/:id", s.UpdateProduct)
	r.PUT("/categories/:id", s.UpdateCategory)
	r.PUT("/payments/:id", s.UpdatePayment)

	r.DELETE("/cashiers/:id", s.DeleteCashier)
	r.DELETE("/products/:id", s.DeleteProduct)
	r.DELETE("/categories/:id", s.DeleteCategory)
	r.DELETE("/payments/:id", s.DeleteCategory)

	authRoutes := r.Group("/").Use(authMiddleware(s.tokenMaker))
	authRoutes.GET("/products", s.SearchProducts)
	authRoutes.GET("/products/:id", s.GetProduct)
	authRoutes.GET("/categories", s.ListCategories)
	authRoutes.GET("/categories/:id", s.GetCategory)
	authRoutes.GET("/payments", s.ListPayments)
	authRoutes.GET("/payments/:id", s.GetPayment)
	authRoutes.GET("/solds", s.GetSold)
	authRoutes.GET("/orders", s.ListAllOrderDetails)
	authRoutes.GET("/orders/:id", s.GetOrderDetails)
	authRoutes.GET("/revenues", s.GetRevenue)

	authRoutes.POST("/orders", s.CreateOrder)
	authRoutes.POST("/orders/subtotal", s.SubtotalOrder)

	s.router = r
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
