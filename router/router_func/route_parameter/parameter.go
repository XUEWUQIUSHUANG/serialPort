package parameter

import (

	datastruct "go_code/serialPort/data_struct"
	"go_code/serialPort/sql"
	"net/http"
	"github.com/gin-gonic/gin"
)

type ParameterCtrl struct{}

func (ctrl *ParameterCtrl) Get_Parameter(c *gin.Context) {
	cookie, _ := c.Cookie("frost_cookie")
	hasusesr := sql.Check_Users(cookie)
	if hasusesr {
		result := sql.Get_Parameter()
		c.JSON(http.StatusOK, gin.H{"data": result})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": "未登录", "hasuser": false})
	}
}

func (ctrl *ParameterCtrl) Set_Parameter(c *gin.Context) {
	var parameter datastruct.EnvironmentData
	c.ShouldBind(&parameter)
	sql.Update_Parameter(parameter)
	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}


