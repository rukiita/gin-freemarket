package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	//ユーザーが削除された際、そのユーザーが登録した商品も自動的に削除されるようになる。
	items []Item `gorm:"constraint:OnDelete:CASCADE`
}