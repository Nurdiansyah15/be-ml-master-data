package controllers

import (
	"fmt"
	"log"
	"ml-master-data/config"
	"ml-master-data/models"
	"ml-master-data/services"
	"ml-master-data/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

// GetAllHeroes godoc
// @Summary Get all heroes
// @Description Get all heroes data
// @Tags Hero
// @Produce json
// @Security Bearer
// @Success 200 {array} models.Hero
// @Router /heroes [get]
func GetAllHeroes(c *gin.Context) {
	var heroes []models.Hero

	if err := config.DB.Find(&heroes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, heroes)
}

// CreateHero godoc
// @Summary Create a hero
// @Description Create a hero and save its image
// @Tags Hero
// @Produce json
// @Security Bearer
// @Param name formData string true "Hero name"
// @Param image formData file true "Hero image"
// @Success 201 {object} models.Hero
// @Router /heroes [post]
func CreateHero(c *gin.Context) {
	// Mengambil nama hero dari FormValue
	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hero name is required"})
		return
	}

	file, err := c.FormFile("image")
	var heroImagePath string

	// Menetapkan path default jika tidak ada gambar
	if err != nil {
		heroImagePath = "https://placehold.co/400x600"
	} else {
		// Memeriksa ukuran file
		if file.Size > 500*1024 { // 500 KB
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size must not exceed 500 KB"})
			return
		}

		// Mendapatkan ekstensi file
		ext := strings.ToLower(filepath.Ext(file.Filename))

		// Validasi ekstensi file
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			return
		}

		// Inisialisasi Cloudinary
		cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
		if err != nil {
			log.Printf("Cloudinary initialization error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Cloudinary"})
			return
		}

		// Buka dan unggah file baru ke Cloudinary
		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
			return
		}
		defer fileContent.Close()

		newFileName := utils.GenerateUniqueFileName("hero")
		uploadResp, err := cld.Upload.Upload(c, fileContent, uploader.UploadParams{
			PublicID: fmt.Sprintf("heroes/%s", newFileName),
			Folder:   "heroes",
		})
		if err != nil {
			log.Printf("Upload to Cloudinary failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload hero image"})
			return
		}

		heroImagePath = uploadResp.SecureURL
	}

	// Membuat objek hero baru
	hero := models.Hero{
		Name:  name,
		Image: heroImagePath,
	}

	// Menyimpan hero ke dalam database
	if err := config.DB.Create(&hero).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mengembalikan response hero yang baru dibuat
	c.JSON(http.StatusCreated, hero)
}

// GetHeroByID godoc
// @Summary Get a hero by ID
// @Description Get a hero data by ID
// @Tags Hero
// @Produce json
// @Security Bearer
// @Param heroID path string true "Hero ID"
// @Success 200 {object} models.Hero
// @Router /heroes/{heroID} [get]
func GetHeroByID(c *gin.Context) {
	heroID := c.Param("heroID")

	var hero models.Hero
	if err := config.DB.First(&hero, heroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		return
	}

	c.JSON(http.StatusOK, hero)
}

// UpdateHero godoc
// @Summary Update a hero
// @Description Update a hero and save its image
// @Tags Hero
// @Produce json
// @Security Bearer
// @Param heroID path string true "Hero ID"
// @Param name formData string false "Hero name"
// @Param image formData file false "Hero image"
// @Success 200 {object} models.Hero
// @Router /heroes/{heroID} [put]
func UpdateHero(c *gin.Context) {
	heroID := c.Param("heroID")
	if heroID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hero ID is required"})
		return
	}

	// Mencari hero berdasarkan ID
	var hero models.Hero
	if err := config.DB.First(&hero, heroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		return
	}

	// Mengambil nama hero dari FormValue
	name := c.PostForm("name")
	if name != "" {
		hero.Name = name
	}

	// Menangani file gambar jika ada
	file, err := c.FormFile("image")
	if err == nil {
		// Memeriksa ukuran file
		if file.Size > 500*1024 { // 500 KB
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size must not exceed 500 KB"})
			return
		}

		// Menghapus gambar lama dari Cloudinary jika ada
		if hero.Image != "" && hero.Image != "https://placehold.co/400x600" {
			publicID := utils.ExtractPublicID(hero.Image)
			cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
			if err != nil {
				log.Printf("Cloudinary initialization error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Cloudinary"})
				return
			}

			// Menghapus gambar lama
			_, err = cld.Upload.Destroy(c, uploader.DestroyParams{
				PublicID:   publicID,
				Invalidate: &[]bool{true}[0],
			})
			if err != nil {
				log.Printf("Failed to delete old image: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old image"})
				return
			}
		}

		// Buka dan unggah file baru ke Cloudinary
		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
			return
		}
		defer fileContent.Close()

		newFileName := utils.GenerateUniqueFileName("hero")
		cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
		if err != nil {
			log.Printf("Cloudinary initialization error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Cloudinary"})
			return
		}

		// Mengunggah gambar baru
		uploadResp, err := cld.Upload.Upload(c, fileContent, uploader.UploadParams{
			PublicID: fmt.Sprintf("heroes/%s", newFileName),
			Folder:   "heroes",
		})
		if err != nil {
			log.Printf("Upload to Cloudinary failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload hero image"})
			return
		}

		// Update URL gambar baru ke dalam database
		hero.Image = uploadResp.SecureURL
	}

	// Menyimpan perubahan hero ke database
	if err := config.DB.Save(&hero).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mengembalikan response hero yang telah diperbarui
	c.JSON(http.StatusOK, hero)
}

// DeleteHero godoc
// @Summary Delete a hero
// @Description Delete a hero by ID and remove its image from the system if it exists
// @Tags Hero
// @Produce json
// @Security Bearer
// @Param heroID path string true "Hero ID"
// @Success 200 {string} string "Hero deleted successfully"
// @Failure 400 {string} string "Hero ID is required"
// @Failure 404 {string} string "Hero not found" or "Old image not found, skipping deletion"
// @Failure 500 {string} string "Failed to remove old image" or "Internal server error"
// @Router /heroes/{heroID} [delete]
func DeleteHero(c *gin.Context) {
	heroID := c.Param("heroID")
	if heroID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hero ID is required"})
		return
	}

	hero := models.Hero{}
	if err := config.DB.Where("hero_id = ?", heroID).First(&hero).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hero not found"})
		return
	}

	if err := services.DeleteHero(config.DB, hero); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kembalikan response sukses
	c.JSON(http.StatusOK, gin.H{"message": "Hero deleted successfully"})
}
