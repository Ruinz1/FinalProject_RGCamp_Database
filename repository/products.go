package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return ProductRepository{db}
}

func (p *ProductRepository) AddProduct(product model.Product) error {
	err := p.db.Create(&product).Error
	if err != nil {
		return err
	}
	
	return nil
}

func (p *ProductRepository) ReadProducts() ([]model.Product, error) {
	products := []model.Product{}

	err := p.db.Table("products").Where("deleted_at IS NULL").
	Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p *ProductRepository) DeleteProduct(id uint) error {
	product := model.Product{}

	err := p.db.Where("id = ?", id).Delete(&product).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductRepository) UpdateProduct(id uint, product model.Product) error {
	
	err := p.db.Table("products").Where("id = ?", id).
	Updates(product).Error
	if err != nil {
		return err
	}
	
	return nil // TODO: replace this
}
