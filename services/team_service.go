package services

import (
	"fmt"
	"log"
	"os"
	"strings"

	"ml-master-data/models"

	"gorm.io/gorm"
)

// Fungsi untuk menghapus Team beserta semua relasi terkait termasuk Match
func DeleteTeam(db *gorm.DB, team models.Team) error {
	// Mulai transaksi
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 1. Hapus Player terkait
	if err := tx.Where("team_id = ?", team.TeamID).Delete(&models.Player{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus Player: %w", err)
	}

	// 2. Hapus Coach terkait
	if err := tx.Where("team_id = ?", team.TeamID).Delete(&models.Coach{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus Coach: %w", err)
	}

	// 3. Temukan semua Match yang melibatkan tim tersebut
	var matches []models.Match
	if err := tx.Where("team_a_id = ? OR team_b_id = ?", team.TeamID, team.TeamID).Find(&matches).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menemukan Match terkait: %w", err)
	}

	// 4. Hapus semua Match dan relasi terkait menggunakan fungsi deleteMatch
	for _, match := range matches {
		if err := DeleteMatch(db, match.MatchID); err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal menghapus Match ID %d: %w", match.MatchID, err)
		}
	}

	// 5. Hapus MatchTeamDetail terkait
	if err := tx.Where("team_id = ?", team.TeamID).Delete(&models.MatchTeamDetail{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus MatchTeamDetail: %w", err)
	}

	// 6. Hapus Goldlaner, Explaner, TrioMid, GameResult, LordResult, dan TurtleResult terkait
	if err := tx.Where("team_id = ?", team.TeamID).Delete(&models.Goldlaner{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus Goldlaner: %w", err)
	}

	if err := tx.Where("team_id = ?", team.TeamID).Delete(&models.Explaner{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus Explaner: %w", err)
	}

	if err := tx.Where("team_id = ?", team.TeamID).Delete(&models.TrioMid{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus TrioMid: %w", err)
	}

	if err := tx.Where("team_id = ?", team.TeamID).Delete(&models.GameResult{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus GameResult: %w", err)
	}

	if err := tx.Where("team_id = ?", team.TeamID).Delete(&models.LordResult{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus LordResult: %w", err)
	}

	if err := tx.Where("team_id = ?", team.TeamID).Delete(&models.TurtleResult{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus TurtleResult: %w", err)
	}

	// 7. Hapus Team itu sendiri
	if err := tx.Delete(&models.Team{}, team.TeamID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menghapus Team: %w", err)
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi: %w", err)
	}

	if team.Image != "" && team.Image != "https://placehold.co/400x600" && strings.HasPrefix(team.Image, os.Getenv("BASE_URL")) {
		team.Image = strings.Replace(team.Image, os.Getenv("BASE_URL")+"/", "", 1)
		// Cek apakah file Image lama ada di sistem
		if _, err := os.Stat(team.Image); err == nil {
			// Jika file ada, hapus file Image lama dari folder images
			if err := os.Remove(team.Image); err != nil {
				return fmt.Errorf("gagal menghapus gambar lama: %w", err)
			}
		} else if os.IsNotExist(err) {
			// Jika file tidak ada, lanjutkan ke tahap selanjutnya
			log.Printf("File gambar lama tidak ditemukan: %s", team.Image)
		}
	}

	log.Printf("Team dengan ID %d dan semua data terkait telah dihapus.", team.TeamID)
	return nil
}
