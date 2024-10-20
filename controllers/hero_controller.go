package controllers

import (
	"fmt"
	"ml-master-data/config"
	"ml-master-data/models"
	"ml-master-data/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetAllHeroes(c *gin.Context) {
	var heroes []models.Hero

	if err := config.DB.Find(&heroes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, heroes)
}

func CreateHero(c *gin.Context) {
	// Mengambil nama hero dari FormValue
	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hero name is required"})
		return
	}

	file, err := c.FormFile("hero_image")
	var heroImagePath string

	// Menetapkan path default jika tidak ada gambar
	if err != nil {
		heroImagePath = "https://placehold.co/400x600"
	} else {
		// Mendapatkan ekstensi file
		ext := strings.ToLower(filepath.Ext(file.Filename))

		// Validasi ekstensi file
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			return
		}

		// Membuat nama file unik
		newFileName := utils.GenerateUniqueFileName("hero") + ext
		heroImagePath = fmt.Sprintf("public/images/%s", newFileName)

		// Menyimpan file yang diupload
		if err := c.SaveUploadedFile(file, heroImagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save hero image"})
			return
		}

		heroImagePath = os.Getenv("BASE_URL") + "/" + heroImagePath
	}

	// Membuat objek hero baru
	hero := models.Hero{
		Name:      name,
		HeroImage: heroImagePath,
	}

	// Menyimpan hero ke database
	if err := config.DB.Create(&hero).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mengembalikan respon dengan hero yang baru dibuat
	c.JSON(http.StatusCreated, hero)
}

func GetHeroByID(c *gin.Context) {
	heroID := c.Param("heroID")

	var hero models.Hero
	if err := config.DB.First(&hero, heroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		return
	}

	c.JSON(http.StatusOK, hero)
}

func UpdateHero(c *gin.Context) {
	heroID := c.Param("heroID")
	if heroID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hero ID is required"})
		return
	}

	var hero models.Hero
	if err := config.DB.First(&hero, heroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		return
	}

	// Mengambil nama hero dari FormValue
	name := c.PostForm("name")

	// Mengambil file gambar dari FormFile
	file, err := c.FormFile("hero_image")

	// Memperbarui nama jika ada
	if name != "" {
		hero.Name = name
	}

	// Memeriksa jika ada gambar baru yang diupload
	if err == nil {

		// Validasi ekstensi file
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			return
		}

		// Cek apakah file gambar lama ada di sistem
		if hero.HeroImage != "" && hero.HeroImage != "https://placehold.co/400x600" {
			// Cek apakah file gambar lama ada di sistem
			hero.HeroImage = strings.Replace(hero.HeroImage, os.Getenv("BASE_URL")+"/", "", 1)
			if _, err := os.Stat(hero.HeroImage); err == nil {
				// Jika file ada, hapus file gambar lama dari folder images
				if err := os.Remove(hero.HeroImage); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove old image"})
					return
				}
			} else if os.IsNotExist(err) {
				// Jika file tidak ada, berikan pesan peringatan (opsional)
				c.JSON(http.StatusNotFound, gin.H{"warning": "Old image not found, skipping deletion"})
			}
		}

		// Membuat nama file unik
		newFileName := utils.GenerateUniqueFileName("hero") + ext
		heroImagePath := fmt.Sprintf("public/images/%s", newFileName)

		// Menyimpan file yang diupload
		if err := c.SaveUploadedFile(file, heroImagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new image"})
			return
		}
		hero.HeroImage = os.Getenv("BASE_URL") + "/" + heroImagePath // Perbarui dengan path gambar baru
	}

	// Simpan perubahan ke database
	if err := config.DB.Save(&hero).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kembalikan response sukses
	c.JSON(http.StatusOK, hero)
}
