package repositories

import (
	"errors"
	"gin-freemarket/models"

	"gorm.io/gorm"
)

type IAuthRepository interface {
	CreateUser(newUser models.User)error
	FindUser(email string)(*models.User,error)
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB)IAuthRepository{
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user models.User) error{
	result:=r.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
} 

func(r *AuthRepository)FindUser(email string)(*models.User,error){
	var user models.User //検索結果を格納する変数を定義
	result := r.db.First(&user,"email=?",email)//第二引数はwhere句、第三引数は?に入る値
	if result.Error != nil{
		if result.Error.Error()=="record not found"{
			return nil,errors.New("User not found")
		}
		return nil, result.Error
	}
	return &user,nil
}