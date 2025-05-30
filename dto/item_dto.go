package dto

type CreateItemInput struct {
	Name        string `json:"name" binding:"required,min=2"`
	Price       uint	`json:"price" binding:"required,min=1,max=999999"`
	Description string	`json:"descroption"`
}

type UpdateItemInput struct {
	//フィールドが存在しない場合には更新しない使用にしたい
	Name 		*string `json:"name" binding:"omitnil,min=2"`
	Price 		*uint	`json:"price binding:"omitnil.min=1,max=999999"`
	Description *string	`json:"description`
	SoldOut 	*bool	`json:"soldOut`
}