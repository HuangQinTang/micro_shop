package model

type User struct {
	ID           int64  `gorm:"primary_key;not_null;auto_increment" json:"id"`
	UserName     string `gorm:"unique_index;not_null" json:"user_name"` //账号
	FirstName    string `json:"first_name"`                             //昵称
	HashPassword string `json:"hash_password"`
}
