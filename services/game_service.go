package services

import (
	"fmt"
	"log"
	"os"
	"strings"

	"ml-master-data/models"

	"gorm.io/gorm"
)

func DeleteGame(db *gorm.DB, game models.Game) error {
	// Mulai transaksi
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Hapus TrioMidHero terkait
	if err := tx.Where("trio_mid_id IN (SELECT trio_mid_id FROM trio_mids WHERE game_id = ?)", game.GameID).Delete(&models.TrioMidHero{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Hapus TrioMid terkait
	if err := tx.Where("game_id = ?", game.GameID).Delete(&models.TrioMid{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Hapus Goldlaner terkait
	if err := tx.Where("game_id = ?", game.GameID).Delete(&models.Goldlaner{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Hapus Explaner terkait
	if err := tx.Where("game_id = ?", game.GameID).Delete(&models.Explaner{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Hapus LordResult terkait
	if err := tx.Where("game_id = ?", game.GameID).Delete(&models.LordResult{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Hapus TurtleResult terkait
	if err := tx.Where("game_id = ?", game.GameID).Delete(&models.TurtleResult{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Hapus GameResult terkait
	if err := tx.Where("game_id = ?", game.GameID).Delete(&models.GameResult{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// update heropick
	var heroPickGames []models.HeroPickGame
	if err := tx.Where("game_number = ?", game.GameNumber).Find(&heroPickGames).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, heroPickGame := range heroPickGames {
		if heroPickGame.IsPicked {
			heroPick := models.HeroPick{}

			if err := tx.Where("hero_pick_id = ?", heroPickGame.HeroPickID).First(&heroPick).Error; err != nil {
				tx.Rollback()
				return err
			}

			fmt.Println("heroPick", heroPick)

			heroPick.Total = heroPick.Total - 1
			if err := tx.Save(&heroPick).Error; err != nil {
				tx.Rollback()
				return err
			}

			fmt.Println("heroPick after", heroPick)
		}
	}

	// Hapus HeroPickGame terkait
	if err := tx.Where("game_number = ?", game.GameNumber).Delete(&models.HeroPickGame{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// update heroban
	var heroBanGames []models.HeroBanGame
	if err := tx.Where("game_number = ?", game.GameNumber).Find(&heroBanGames).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, heroBanGame := range heroBanGames {
		if heroBanGame.IsBanned {
			heroBan := models.HeroBan{}

			if err := tx.Where("hero_ban_id = ?", heroBanGame.HeroBanID).First(&heroBan).Error; err != nil {
				tx.Rollback()
				return err
			}

			heroBan.Total = heroBan.Total - 1
			if err := tx.Save(&heroBan).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Hapus HeroBanGame terkait
	if err := tx.Where("game_number = ?", game.GameNumber).Delete(&models.HeroBanGame{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Hapus Game itu sendiri
	if err := tx.Delete(&models.Game{}, game.GameID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return err
	}

	if game.FullDraftImage != "" && game.FullDraftImage != "https://placehold.co/400x600" && strings.HasPrefix(game.FullDraftImage, os.Getenv("BASE_URL")) {
		game.FullDraftImage = strings.Replace(game.FullDraftImage, os.Getenv("BASE_URL")+"/", "", 1)
		// Cek apakah file Image lama ada di sistem
		if _, err := os.Stat(game.FullDraftImage); err == nil {
			// Jika file ada, hapus file Image lama dari folder images
			if err := os.Remove(game.FullDraftImage); err != nil {
				return fmt.Errorf("gagal menghapus gambar lama: %w", err)
			}
		} else if os.IsNotExist(err) {
			// Jika file tidak ada, lanjutkan ke tahap selanjutnya
			log.Printf("File gambar lama tidak ditemukan: %s", game.FullDraftImage)
		}
	}

	log.Printf("Game with ID %d and all related records have been deleted.", game.GameID)
	return nil
}
