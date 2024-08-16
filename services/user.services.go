package services

import (
	"errors"

	"github.com/ramyjaber1/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func NewUserServices(u models.User, uStore *gorm.DB) *UserServices {
	return &UserServices{
		User: u,
		DB:   uStore,
	}
}

type UserServices struct {
	User models.User
	DB   *gorm.DB
}

func (us *UserServices) GetALL(limit, page int, orderBy, sortBy, searchTerm, status string) ([]models.User, models.Meta, error) {
	var users []models.User
	query := us.DB
	totalQuery := query

	if searchTerm != "" {
		searchTermWithWildcard := "%" + searchTerm + "%"
		query = query.Where("name LIKE ?", searchTermWithWildcard)
		totalQuery = query
	}

	if status != "" {
		query = query.Where("status = ?", status)
		totalQuery = totalQuery.Where("status = ?", status)
	}
	offset := (page - 1) * limit
	query.Order(sortBy + " " + orderBy).Offset(offset).Limit(limit).Find(&users)
	totalRecords := int64(0)
	totalQuery.Model(&us.User).Count(&totalRecords)
	lastPage := int64(0)
	if limit > 0 {
		lastPage = (totalRecords + int64(limit) - 1) / int64(limit)
	}
	meta := models.Meta{
		CurrentPage: page,
		TotalCount:  int(totalRecords),
		LastPage:    int(lastPage),
		Limit:       limit,
	}

	return users, meta, nil
}

func (us *UserServices) CreateUser(u models.User, passwordConfirm string) error {
	if u.Password != passwordConfirm {
		return errors.New("password not match")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	if result := us.DB.Create(&u); result.Error != nil {
		return result.Error
	}
	return nil
}

func (us *UserServices) CheckEmail(email string) (models.User, error) {
	result := us.DB.Where("email = ?", email).First(&us.User)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return us.User, nil
}

/* func (us *UserServices) GetUserById(id int) (User, error) {

	query := `SELECT id, email, password, username FROM users
		WHERE id = ?`

	stmt, err := us.UserStore.Db.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	us.User.ID = id
	err = stmt.QueryRow(
		us.User.ID,
	).Scan(
		&us.User.ID,
		&us.User.Email,
		&us.User.Password,
		&us.User.Username,
	)
	if err != nil {
		return User{}, err
	}

	return us.User, nil
} */
