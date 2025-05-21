package router

import (
	parameter "go_code/serialPort/router/router_func/route_parameter"
	users "go_code/serialPort/router/router_func/route_users"

	"github.com/gin-gonic/gin"
)

func Init_listen(r *gin.Engine) {

	parameter := &parameter.ParameterCtrl{}
	users := &users.UsersCtrl{}

	api := r.Group("/api")
	{
		api.POST("/login", users.Login)
		api.POST("/register", users.Register)
		api.GET("/checkusers", users.Check_Users)
		api.GET("/cache", users.Get_Cache)
		api.POST("/changepassword", users.Change_Password)

		api.GET("/parameters", parameter.Get_Parameter)
		api.POST("/parameters", parameter.Set_Parameter)
		
	}
}
