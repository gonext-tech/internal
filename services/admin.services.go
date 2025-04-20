package services

import (
	"errors"
	"time"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/queries"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func NewAdminService(uStore *gorm.DB) *AdminServices {
	return &AdminServices{
		DB: uStore,
	}
}

type AdminServices struct {
	DB *gorm.DB
}

func (us *AdminServices) GetALL(query queries.InvoiceQueryParams) ([]models.Admin, models.Meta, error) {
	var admins []models.Admin
	var totalCount int64

	dbQuery := us.buildInvoicesURL(query)
	dbQuery.Session(&gorm.Session{}).Model(&models.Invoice{}).Count(&totalCount)
	offset := (query.Page - 1) * query.Limit
	dbQuery.Order(query.SortBy + " " + query.OrderBy).Offset(offset).Limit(query.Limit).Find(&admins)
	lastPage := int64(0)
	if query.Limit > 0 {
		lastPage = (totalCount + int64(query.Limit) - 1) / int64(query.Limit)
	}
	meta := models.Meta{
		CurrentPage: query.Page,
		TotalCount:  int(totalCount),
		LastPage:    int(lastPage),
		Limit:       query.Limit,
	}

	return admins, meta, nil
}

func (us *AdminServices) CreateUser(u models.Admin, passwordConfirm string) error {
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

func (us *AdminServices) CheckEmail(email string) (models.Admin, error) {
	var admin models.Admin
	result := us.DB.Where("email = ?", email).First(&admin)
	if result.Error != nil {
		return admin, result.Error
	}
	return admin, nil
}

func (us *AdminServices) GetID(id string) (*models.Admin, error) {
	var admin models.Admin
	result := us.DB.Where("id = ?", id).First(&admin)
	if result.Error != nil {
		return nil, result.Error
	}
	return &admin, nil
}

func (us *AdminServices) Create(u *models.Admin) error {
	return us.DB.Create(&u).Error
}

func (us *AdminServices) Update(u *models.Admin) error {

	return us.DB.Model(&u).Updates(u).Error

}

func (us *AdminServices) Delete(u *models.Admin) error {
	now := time.Now()
	u.Status = "NOT_ACTIVE"
	u.Deleted = true
	u.DeletedAt = &now
	return us.DB.Model(&u).Updates(u).Error
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

func (us *AdminServices) buildInvoicesURL(query queries.InvoiceQueryParams) *gorm.DB {
	dbQuery := us.DB

	if query.SearchTerm != "" {
		term := "%" + query.SearchTerm + "%"
		dbQuery = dbQuery.Where("name LIKE ?", term)
	}
	if query.Status != "" {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}
	if query.InvoiceType != "" {
		dbQuery = dbQuery.Where("invoice_type = ?", query.InvoiceType)
	}
	return dbQuery

}
