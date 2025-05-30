package middlewares

import (
	"gin-freemarket/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//1.リクエストヘッダーの検証
//2.トークンを元にユーザー情報の取得
//3.取得したユーザー情報をリクエストへセットする
func AuthMiddleware(authService services.IAuthService) gin.HandlerFunc{
	return func(ctx *gin.Context){
	header := ctx.GetHeader("Authorization")
	if header == ""{
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !strings.HasPrefix(header,"Bearer"){
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return 
	}
	tokenString := strings.TrimPrefix(header,"Bearer ")
	user,err := authService.GetUserFromToken(tokenString)
	if err != nil{
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return 
	}
	if !strings.HasPrefix(header,"Bearer"){
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return 
	}
	//リクエストの生存期間中に製造されるキーと値のペアを設定
	ctx.Set("user",user)
	ctx.Next()
  }
}