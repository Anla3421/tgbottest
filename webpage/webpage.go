package webpage

import (
	"fmt"
	"net/http"
	"server/db/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

func init() {
	fmt.Println("init web server")
}
func StartWebServer() {
	fmt.Println("Web initial complete")
	apiserver := gin.New()
	apiserver.GET("/test/:id", test)
	apiserver.Run(":3388")

}

func test(context *gin.Context) {
	id := context.Param("id")
	inNum, err := strconv.Atoi(id)
	if err != nil {
		context.JSON(http.StatusOK, gin.H(map[string]interface{}{"error": "1"}))
		return
	}
	if inNum > 3 || inNum <= 0 {
		context.JSON(http.StatusOK, gin.H(map[string]interface{}{"error": "1"}))
		return
	}
	catch := sql.Websql(id)
	context.JSON(http.StatusOK, catch.Name+"  ")
	context.JSON(http.StatusOK, catch.Text)
}
