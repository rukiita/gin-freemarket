package services

import (
	"fmt"
	"gin-freemarket/models"
	"gin-freemarket/repositories"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Signup(email string, password string) error
	Login(email string,password string) (*string,error)
	GetUserFromToken(tokenString string)(*models.User,error)
}

type AuthService struct {
	repository repositories.IAuthRepository
}

func NewAuthService(repository repositories.IAuthRepository)IAuthService{
	return &AuthService{repository: repository}
}

func (s *AuthService) Signup(email string, password string) error{
	//クリプトデフォルトコストはハッシュ化のコストを指定するパラメータで、
	// セキュリティとパフォーマンスのバランスをとるために使用されます
	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user:=models.User{
		Email:email,
		Password:string(hashedPassword),
	}
	return s.repository.CreateUser(user)
}

func(s *AuthService)Login(email string,password string)(*string,error){
	foundUser,err := s.repository.FindUser(email)
	if err != nil {
		return nil, err
	}
	//第一引数にbyteのスライスに変換したパスワード、第二引数は
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password),[]byte(password))
	if err != nil {
		return nil, err
	}

	token, err := CreateToken(foundUser.ID,foundUser.Email)
	return token,nil
}

func CreateToken(userId uint,email string)(*string,error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"sub":userId,
		"email":email,
		"exp":time.Now().Add(time.Hour).Unix(),
	})

	tokenString,err:=token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil{
		return nil,err
	}
	return &tokenString,nil
}

func (s *AuthService) GetUserFromToken(tokenString string)(*models.User,error){
	token, err := jwt.Parse(tokenString,func(token *jwt.Token) (interface{}, error) {
		if _,ok:= token.Method.(*jwt.SigningMethodHMAC);!ok{
			return nil,fmt.Errorf("Unexpected signing method:%v",token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")),nil
	})
	if err != nil {
		return nil,err
	}
	var user *models.User
	if claims,ok:= token.Claims.(jwt.MapClaims); ok{
		if float64(time.Now().Unix())>claims["exp"].(float64){
			return nil,jwt.ErrTokenExpired
		}
		user,err = s.repository.FindUser(claims["email"].(string))
	}
	if err != nil{
		return nil,err
	}
	return user, nil
}