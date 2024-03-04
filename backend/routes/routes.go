package routes

import (
	"github.com/AmanPr33tS1ngh/go-Ecomm.git/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine)  {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/admin/add_product", controllers.ProductViewerAdmin())
	incomingRoutes.POST("/users/product_view", controllers.SearchProduct())
	incomingRoutes.POST("/users/search", controllers.SearchProductByAdmin())
}