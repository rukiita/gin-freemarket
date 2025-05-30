package controllers

import (
	"gin-freemarket/dto"
	"gin-freemarket/models"
	"gin-freemarket/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IItemController interface {
	FindAll(ctx *gin.Context)
	FindById(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	DeleteById(ctx *gin.Context)
}

type ItemController struct{
	service services.IItemService
}

func NewItemController(service services.IItemService)IItemController{
	return &ItemController{service:service}
}

func(c *ItemController)FindAll(ctx *gin.Context){
	//全件取得をサービスに任せる
	items, err :=c.service.FindAll()
	if err != nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":"Unexpected error"})
		return
	}
	//Gin のコンテキストを使って、HTTP ステータスコード 200 OK と、
	// 取得した商品データ (items) を含む JSON レスポンスをクライアントに返します。
	ctx.JSON(http.StatusOK,gin.H{"data":items})
}

func(c *ItemController)FindById(ctx *gin.Context){
	user,exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	userId := user.(*models.User).ID
	//IDをパースして準備する
	itemId,err := strconv.ParseUint(ctx.Param("id"),10,64)
	if err != nil{
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"Invalid id"})
		return
	}
	//IDで検索するのをサービスに任せる
	item,err := c.service.FindById(uint(itemId),userId)
	if err != nil{
		if err.Error()== "Item no foud"{
			ctx.JSON(http.StatusNotFound,gin.H{"error":err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":"Unexpected error"})
		return
	}
	//Gin のコンテキストを使って、HTTP ステータスコード 200 OK と、
	// 取得した商品データ (items) を含む JSON レスポンスをクライアントに返します。
	ctx.JSON(http.StatusOK,gin.H{"data":item})
}

func(c *ItemController)Create(ctx *gin.Context){ 
	user,exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	userId := user.(*models.User).ID
	//updateする際の型の変数を用意
	var input dto.CreateItemInput
	//HTTP リクエストのボディに含まれる JSON データを、変数 input のアドレス 
	// (&input) が指す構造体にマッピング（バインド）しようとします。
	if err := ctx.ShouldBindJSON(&input);err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	//エラー出なければ一件作成するのをサービスに任せる
	newItem,err:=c.service.Create(input,userId)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	//Gin のコンテキストを使って、HTTP ステータスコード 201 OK と、
	//作成したデータ を含む JSON レスポンスをクライアントに返します。
	ctx.JSON(http.StatusCreated,gin.H{"data":newItem})
}

func(c *ItemController) Update(ctx *gin.Context){
	user,exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	userId := user.(*models.User).ID

	itemId,err:=strconv.ParseUint(ctx.Param("id"),10,64)
	if err != nil{
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"Invalid id"})
		return
	}
	var input dto.UpdateItemInput
	if err := ctx.ShouldBindJSON(&input);err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	
	updatedItem,err:=c.service.Update(uint(itemId),input,userId)
	if err != nil{
		if err.Error()== "Item no foud"{
			ctx.JSON(http.StatusNotFound,gin.H{"error":err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":"Unexpected error"})
		return
	}
	ctx.JSON(http.StatusOK,gin.H{"data":updatedItem})
}

func(c *ItemController) DeleteById(ctx *gin.Context){
	user,exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	userId := user.(*models.User).ID

	itemId,err:=strconv.ParseUint(ctx.Param("id"),10,64)
	if err != nil{
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"Invalid id"})
		return
	}
	err =c.service.DeleteById(uint(itemId),userId)
	if err != nil{
		if err.Error()== "Item no foud"{
			ctx.JSON(http.StatusNotFound,gin.H{"error":err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":"Unexpected error"})
		return
	}
	ctx.Status(http.StatusOK)
}