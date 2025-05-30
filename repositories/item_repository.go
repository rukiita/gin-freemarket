package repositories

import (
	"errors"
	"gin-freemarket/models"

	"gorm.io/gorm"
)

// リポジトリのインタフェース定義
type IItemRepository interface {
	FindAll() (*[]models.Item, error)
	FindById(itemId uint) (*models.Item, error)
	Create(newItem models.Item) (*models.Item, error)
	Update(updateItem models.Item) (*models.Item, error)
	DeleteById(itemId uint) error
}

// リポジトリの実装
type ItemMemoryRepository struct {
	items []models.Item
}

// リポジトリのインスタンスを作成するメソッド
func NewItemMemoryRepository(items []models.Item) IItemRepository {
	return &ItemMemoryRepository{items: items}
	//具体的な実装であるItemMemoryRepositoryインスタンスを作成しつつ、その型を
	//インタフェースであるIItemRepositoryで返す。
}

func (r *ItemMemoryRepository) FindAll() (*[]models.Item, error) {
	return &r.items, nil
}

func (r *ItemMemoryRepository) FindById(itemId uint) (*models.Item, error) {
	for _, v := range r.items {
		if v.ID == itemId {
			return &v, nil
		}
	}
	return nil, errors.New("Item no found")
}

func (r *ItemMemoryRepository) Create(newItem models.Item) (*models.Item, error) {
	newItem.ID = uint(len(r.items) + 1)
	r.items = append(r.items, newItem)
	return &newItem, nil
}

func (r *ItemMemoryRepository) Update(updateItem models.Item) (*models.Item, error) {
	for i, v := range r.items {
		if v.ID == updateItem.ID {
			r.items[i] = updateItem
			return &r.items[i], nil
		}
	}
	return nil, errors.New("unexpected error")
}

func (r *ItemMemoryRepository) DeleteById(itemId uint) error {
	for i, v := range r.items {
		if v.ID == itemId {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return nil
		}
	}
	return errors.New("Item not found")
}


type itemRepository struct {
	db *gorm.DB
}

//メモリの場合と同様にdbの場合もIItemRepository型のインスタンスを作れる
func NewItemRepository(db *gorm.DB) IItemRepository {
	return &itemRepository{db: db}
}

// FindAll implements IItemRepository.
func (r *itemRepository) FindAll() (*[]models.Item, error) {
	var items []models.Item
	result := r.db.Find(&items)
	if result.Error != nil {
		return nil,result.Error
	}
	return &items,nil
}

// FindById implements IItemRepository.
func (r *itemRepository) FindById(itemId uint) (*models.Item, error) {
	var item models.Item
	result:=r.db.First(&item,itemId)
	if result.Error !=nil{
		if result.Error.Error() == "record not found"{
			return nil,errors.New("Item not found")
		}
		return nil,result.Error
	}
	return &item, nil
}

// Create implements IItemRepository.
func (r *itemRepository) Create(newItem models.Item) (*models.Item, error) {
	result := r.db.Create(&newItem)
	//エラーはgormの構造体のエラーフィールドに格納される
	if result.Error != nil {
		return nil,result.Error
	}
	return &newItem, nil
}

// Update implements IItemRepository.
func (r *itemRepository) Update(updateItem models.Item) (*models.Item, error) {
	result := r.db.Save(&updateItem)
	if result.Error !=nil{
		return nil,result.Error
	}
	return &updateItem,nil
}

// DeleteById implements IItemRepository.
func (r *itemRepository) DeleteById(itemId uint) error {
	deleteItem,err := r.FindById(itemId)
	if err!=nil {
		return err
	}
	result := r.db.Delete(&deleteItem)
	if result.Error!=nil{
		return result.Error
	}
	return nil
}

