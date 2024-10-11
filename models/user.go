package models

type User struct {
	UserID   uint   `gorm:"primaryKey;autoIncrement" json:"user_id"`
	Username string `gorm:"unique" json:"username"`
	Password string `json:"-"`
}
