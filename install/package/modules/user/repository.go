package user

import (
	"iwogo/modules/user/entity"

	"gorm.io/gorm"
)

type Repository interface {
	Save(user entity.User) (entity.User, error)
	FindByEmail(email string) (entity.User, error)
	FindById(id int) (entity.User, error)
	Update(user entity.User) (entity.User, error)
	AllUser() ([]entity.User, error)
	Delete(id int, user entity.User) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(user entity.User) (entity.User, error) {
	err := r.db.Create(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) AllUser() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Find(&users).Error

	if err != nil {
		return users, err
	}
	return users, nil
}

func (r *repository) FindByEmail(email string) (entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil

}

func (r *repository) FindById(id int) (entity.User, error) {
	var user entity.User

	err := r.db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil

}

func (r *repository) Update(user entity.User) (entity.User, error) {
	err := r.db.Save(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) Delete(id int, user entity.User) (bool, error) {
	err := r.db.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
