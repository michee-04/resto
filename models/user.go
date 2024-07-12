package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/michee-04/resto/database"
	"github.com/michee-04/resto/utils"
	"gorm.io/gorm"
)

var Db *gorm.DB

type User struct {
	UserId        string    `gorm:"not null; unique; column:user_id; primary_key"`
	Username      string    `gorm:"column:username; not null;unique" json:"username"`
	Email         string    `gorm:"not null;unique;column:email" json:"email"`
	EmailVerify   bool      `gorm:"column:email_verify" json:"email_verify"`
	Password      string    `gorm:"column:password" json:"password"`
	IsAdmin       bool      `gorm:"column:is_admin" json:"is_admin"`
	Avatar        string    `gorm:"column:avatar" json:"avatar"`
	Token         string    `gorm:"column:token" json:"token"`
	TokenAccount  string    `gorm:"column:token_account" json:"token_account"`
	TokenPassword string    `gorm:"column:token_password" json:"token_password"`
	ExpiryToken   time.Time `gorm:"column:expiry_token" json:"expiry_token"`
}

func init() {
	database.ConnectDB()
	Db = database.GetBD()
	Db.Migrator().DropTable(&User{})
	Db.AutoMigrate(&User{})
}

func (u *User) CreateUser() *User {
	u.UserId = uuid.New().String()
	Db.Create(&u)
	return u
}

func GetUser() []User {
	var u []User
	Db.Find(&u)
	return u
}

func GetUserById(Id string) (*User, *gorm.DB) {
	var u User
	db := Db.Where("user_id=?", Id).First(&u)
	return &u, db
}

func GetUserByEmail(email string) (*User, error) {
	var u User
	if err := Db.Where("email=?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func DeleteUser(Id string) (*User, error) {
	var u User
	Db.Where("user_id=?", Id).Delete(&u)
	return &u, nil
}

func (u *User) Logout() error {
	u.Token = ""
	return Db.Save(&u).Error
}

func FindUserByToken(t string) (*User, error) {
	var u User
	if err := Db.Where("token_account=?", t).First(&u).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &u, nil
}

func FindUserPasswordToken(token string) (*User, error) {
	var u User
	if err := Db.Where("token_password = ?", token).First(&u).Error; err != nil {
		fmt.Println("Error finding user by token:", err)
		return nil, fmt.Errorf("token password not found")
	}

	// Vérifiez si le token a expiré
	if time.Now().After(u.ExpiryToken) {
		return nil, fmt.Errorf("reset token has expired")
	}
	return &u, nil
}

func (u *User) Verify() error {
	u.EmailVerify = true
	u.TokenAccount = ""
	Db.Save(&u)
	return nil
}

func (u *User) GeneratePasswordToken() error {
	token := uuid.New().String()
	u.TokenPassword = token
	fmt.Println("Generated token:", token)
	return Db.Save(&u).Error
}

func (u *User) UpdatePassword(newPassword string) error {
	hashedPassword, err := utils.HashedPassword(newPassword)
	if err != nil {
		return err
	}

	err = Db.Transaction(func(tx *gorm.DB) error {
		u.Password = hashedPassword
		u.TokenPassword = ""
		u.ExpiryToken = time.Time{}

		if err := tx.Save(&u).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("New hased Password: ", hashedPassword)
	return nil
}
