package services

import (
	"gin-freemarket/dto"
	"gin-freemarket/models"
	"gin-freemarket/repositories"
)

//Serviceのインタフェースを定義
type IItemService interface {
	FindAll() (*[]models.Item,error)
	FindById(itemId uint,userId uint)(*models.Item,error)
	Create(createItemInput dto.CreateItemInput,userId uint)(*models.Item,error)
	Update(itemId uint,UpdateItemInput dto.UpdateItemInput,userId uint)(*models.Item,error)
	DeleteById(itemId uint,userId uint)error
}

type ItemService struct{
	repository repositories.IItemRepository
}

func NewItemService(repository repositories.IItemRepository) IItemService {
	//具体的な実装であるItemServiceのインスタンスを作成しつつ、
	//その方をインタフェースであるIItemRepositoryとして返す。
	return &ItemService{repository: repository}
}

func(s *ItemService)FindAll()(*[]models.Item,error){
	return s.repository.FindAll()
}

func(s *ItemService)FindById(itemId uint,userId uint)(*models.Item,error){
	return s.repository.FindById(itemId)
}

func(s *ItemService)Create(CreateItemInput dto.CreateItemInput,userId uint)(*models.Item,error){
	newItem := models.Item{
		Name: CreateItemInput.Name,
		Price: CreateItemInput.Price,
		Description: CreateItemInput.Description,
		SoldOut: false,
		UserID: userId,
	}
	return s.repository.Create(newItem)
}

func(s *ItemService)Update(itemId uint, updateItemInput dto.UpdateItemInput,userId uint)(*models.Item,error){
	targetItem,err := s.FindById(itemId,userId)
	if err != nil {
		return nil,err
	}
	if updateItemInput.Name != nil{
		targetItem.Name = *updateItemInput.Name
	}
	if updateItemInput.Price != nil{
		targetItem.Price = *updateItemInput.Price
	}
	if updateItemInput.Description != nil{
		targetItem.Description = *updateItemInput.Description
	}
	if updateItemInput.SoldOut != nil{
		targetItem.SoldOut = *updateItemInput.SoldOut
	}
	return s.repository.Update(*targetItem)
}

func(s *ItemService)DeleteById(itemId uint,userId uint)error{
	return s.repository.DeleteById(itemId)
}