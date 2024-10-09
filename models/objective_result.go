package models

type ObjectiveResult struct {
	ObjectiveResultID uint   `gorm:"primaryKey;autoIncrement"`
	GameID            uint   `gorm:"not null"`
	TeamID            uint   `gorm:"not null"`
	ObjectiveType     string `gorm:"size:20"`
	ObjectiveNumber   int
	SetupResult       string `gorm:"size:20"`
	InitiateResult    string `gorm:"size:20"`
	FinalResult       string `gorm:"size:20"`
}
