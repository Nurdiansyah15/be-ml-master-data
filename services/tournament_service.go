package services

import (
	"log"

	"ml-master-data/models"

	"gorm.io/gorm"
)

// Fungsi untuk menghapus Tournament dan semua relasi terkait
func DeleteTournament(db *gorm.DB, tournamentID uint) error {
	// Mulai transaksi
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 1. Temukan semua Match terkait dengan Tournament ini
	var matches []models.Match
	if err := tx.Where("tournament_id = ?", tournamentID).Find(&matches).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. Panggil DeleteMatch untuk setiap Match terkait
	for _, match := range matches {
		if err := DeleteMatch(db, match.MatchID); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 3. Hapus Tournament itu sendiri
	if err := tx.Delete(&models.Tournament{}, tournamentID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return err
	}

	log.Printf("Tournament with ID %d and all related records have been deleted.", tournamentID)
	return nil
}
