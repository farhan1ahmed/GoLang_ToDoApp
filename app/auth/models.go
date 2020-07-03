package auth

import "github.com/jinzhu/gorm"

var tokendb *gorm.DB

type AuthApp struct {
	DB *gorm.DB
}

type BlackListToken struct {
	TokenVal string `gorm:"notnull"`
}

func (aApp *AuthApp) InitBlackListModel() {
	tokendb = aApp.DB
	tokendb.AutoMigrate(&BlackListToken{})
}
