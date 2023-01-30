package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return CartRepository{db}
}

func (c *CartRepository) ReadCart() ([]model.JoinCart, error) {
	join := []model.JoinCart{}

	err := c.db.Table("carts").Select("carts.id, products.id as product_id, products.name, carts.quantity, carts.total_price").
	Joins("join products on carts.product_id = products.id").Scan(&join).Error
	if err != nil {
		return nil, err
	}


	return join, nil 
}

func (c *CartRepository) AddCart(product model.Product) error {
	if c.db.Where("product_id = ?", product.ID).First(&model.Cart{}).RowsAffected == 1{
		c.db.Transaction(func(tx *gorm.DB) error {
			total := product.Price - (product.Price * (product.Discount / 100))
			tx.Model(&model.Cart{}).Where("product_id = ?", product.ID).Update("quantity", gorm.Expr("quantity + ?", 1))
			tx.Model(&model.Cart{}).Where("product_id = ?", product.ID).Update("total_price", gorm.Expr("total_price + ?", total))
			tx.Transaction(func(tx *gorm.DB) error {
				tx.Model(&model.Product{}).Where("id = ?", product.ID).Update("stock", gorm.Expr("stock - ?", 1))
				return nil
			})

		return nil
		})
	} else {
		c.db.Transaction(func(tx *gorm.DB) error {
			tx.Create(&model.Cart{ProductID: product.ID, Quantity:  1, TotalPrice:  product.Price - (product.Price * (product.Discount) / 100)})
			tx.Transaction(func(tx *gorm.DB) error {
				tx.Model(&model.Product{}).Where("id = ?", product.ID).Update("stock", gorm.Expr("stock - ?", 1))
			return nil
			})
			return nil
		})

	}
	return nil 
}

func (c *CartRepository) DeleteCart(id uint, productID uint) error {

	c.db.Transaction(func(tx *gorm.DB) error {
		cart := model.Cart{}
		product := model.Product{}

		err := tx.Where("id = ?", id).Delete(&cart).Error
		if err != nil {
			return err
		}

		err = tx.Model(&product).Where("id = ?", productID).Update("stock", gorm.Expr("stock + ?", 1)).Error
		if err != nil {
			return err
		}

		return nil
	})

	return nil 
}

func (c *CartRepository) UpdateCart(id uint, cart model.Cart) error {

	err := c.db.Table("carts").Where("id = ?", id).
	Updates(cart).Error
	if err != nil {
		return err
	}
	
	return nil
}
