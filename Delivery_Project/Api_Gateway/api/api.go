package api

import (
	"api-gateway-service/api/handler"
	"api-gateway-service/config"

	_ "api-gateway-service/api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func SetUpApi(r *gin.Engine, h *handler.Handler, cfg config.Config) {
	r.Use(customCORSMiddleware())
	r.Use(MaxAllowed(500))

	v1 := r.Group("/v1")

	//Product service
	// // product api
	v1.POST("/product", h.CreateProduct)
	v1.GET("/product", h.GetListProduct)
	v1.GET("/product/:id", h.GetProduct)
	v1.PUT("/product/:id", h.UpdateProduct)
	v1.DELETE("/product/:id", h.DeleteProduct)

	// category api
	v1.POST("/category", h.CreateCategory)
	v1.GET("/category", h.GetListCategory)
	v1.GET("/category/:id", h.GetCategory)
	v1.PUT("/category/:id", h.UpdateCategory)
	v1.DELETE("/category/:id", h.DeleteCategory)

	//Order service
	// // product api
	v1.POST("/order", h.CreateOrder)
	v1.GET("/order", h.GetListOrder)
	v1.GET("/order/:id", h.GetOrder)
	v1.PUT("/order/:id", h.UpdateOrder)
	v1.DELETE("/order/:id", h.DeleteOrder)

	// delivery_tariff api
	v1.POST("/delivery_tariff", h.CreateDeliveryTariff)
	v1.GET("/delivery_tariff", h.GetListDeliveryTariff)
	v1.GET("/delivery_tariff/:id", h.GetDeliveryTariff)
	v1.PUT("/delivery_tariff/:id", h.UpdateDeliveryTariff)
	v1.DELETE("/delivery_tariff/:id", h.DeleteDeliveryTariff)

	//User service
	// // branch api
	v1.POST("/branch", h.CreateBranch)
	v1.GET("/branch", h.GetListBranch)
	v1.GET("/branch/:id", h.GetBranch)
	v1.PUT("/branch/:id", h.UpdateBranch)
	v1.DELETE("/branch/:id", h.DeleteBranch)

	// user api
	v1.POST("/user", h.CreateUser)
	v1.GET("/user", h.GetListUser)
	v1.GET("/user/:id", h.GetUser)
	v1.PUT("/user/:id", h.UpdateUser)
	v1.DELETE("/user/:id", h.DeleteUser)

	// client api
	v1.POST("/client", h.CreateClients)
	v1.GET("/client", h.GetListClients)
	v1.GET("/client/:id", h.GetClients)
	v1.PUT("/client/:id", h.UpdateClients)
	v1.DELETE("/client/:id", h.DeleteClients)

	// courier api
	v1.POST("/courier", h.CreateCourier)
	v1.GET("/courier", h.GetListCourier)
	v1.GET("/courier/:id", h.GetCourier)
	v1.PUT("/courier/:id", h.UpdateCourier)
	v1.DELETE("/courier/:id", h.DeleteCourier)

	// Logic api
	v1.GET("/logic", h.GetCourierOrders)
	v1.PUT("/logic/:id", h.UpdateOrderStatus)
	v1.GET("/branch/active", h.GetListActiveBranch)
	v1.GET("/courier/active-orders/list", h.CourierGetOrder)
	v1.GET("/courier/delete_order/:id", h.DeleteCourierInOrder)
	v1.GET("/courier/get_order/:id", h.GetCourierOrders)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

func MaxAllowed(n int) gin.HandlerFunc {
	var countReq int64
	sem := make(chan struct{}, n)
	acquire := func() {
		sem <- struct{}{}
		countReq++
	}

	release := func() {
		select {
		case <-sem:
		default:
		}
		countReq--
	}

	return func(c *gin.Context) {
		acquire()       // before request
		defer release() // after request

		c.Set("sem", sem)
		c.Set("count_request", countReq)

		c.Next()
	}
}

func customCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE, HEAD")
		c.Header("Access-Control-Allow-Headers", "Platform-Id, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Max-Age", "3600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
