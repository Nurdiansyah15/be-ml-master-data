package services

import (
	"log"

	"ml-master-data/models"

	"gorm.io/gorm"
)

// Fungsi untuk menghapus Match dan semua relasi terkait
func DeleteMatch(db *gorm.DB, matchID uint) error {
	// Mulai transaksi
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 1. Hapus semua Game terkait dengan Match ini
	var games []models.Game
	if err := tx.Where("match_id = ?", matchID).Find(&games).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Panggil DeleteGame untuk setiap Game terkait
	for _, game := range games {
		if err := DeleteGame(db, game); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 5. Hapus HeroPick dan HeroPickGame terkait
	if err := tx.Where("hero_pick_id IN (SELECT hero_pick_id FROM hero_picks WHERE match_team_detail_id IN (SELECT match_team_detail_id FROM match_team_details WHERE match_id = ?))", matchID).
		Delete(&models.HeroPickGame{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("match_team_detail_id IN (SELECT match_team_detail_id FROM match_team_details WHERE match_id = ?)", matchID).
		Delete(&models.HeroPick{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 6. Hapus HeroBan dan HeroBanGame terkait

	if err := tx.Where("hero_ban_id IN (SELECT hero_ban_id FROM hero_bans WHERE match_team_detail_id IN (SELECT match_team_detail_id FROM match_team_details WHERE match_id = ?))", matchID).
		Delete(&models.HeroBanGame{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("match_team_detail_id IN (SELECT match_team_detail_id FROM match_team_details WHERE match_id = ?)", matchID).
		Delete(&models.HeroBan{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 4. Hapus FlexPick terkait
	if err := tx.Where("match_team_detail_id IN (SELECT match_team_detail_id FROM match_team_details WHERE match_id = ?)", matchID).
		Delete(&models.FlexPick{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 7. Hapus PriorityPick dan PriorityBan terkait
	if err := tx.Where("match_team_detail_id IN (SELECT match_team_detail_id FROM match_team_details WHERE match_id = ?)", matchID).
		Delete(&models.PriorityPick{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("match_team_detail_id IN (SELECT match_team_detail_id FROM match_team_details WHERE match_id = ?)", matchID).
		Delete(&models.PriorityBan{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 3. Hapus CoachMatch terkait
	if err := tx.Where("match_team_detail_id IN (SELECT match_team_detail_id FROM match_team_details WHERE match_id = ?)", matchID).
		Delete(&models.CoachMatch{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 8. Hapus PlayerMatch terkait
	if err := tx.Where("match_team_detail_id IN (SELECT match_team_detail_id FROM match_team_details WHERE match_id = ?)", matchID).
		Delete(&models.PlayerMatch{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. Hapus MatchTeamDetail terkait
	if err := tx.Where("match_id = ?", matchID).Delete(&models.MatchTeamDetail{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 9. Hapus Match itu sendiri
	if err := tx.Delete(&models.Match{}, matchID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return err
	}

	log.Printf("Match with ID %d and all related records have been deleted.", matchID)
	return nil
}
