package handlers

import (
	"calhacks/api/db"
	"calhacks/api/lib"
	"calhacks/api/models"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ScanRequest struct {
	Owner string `json:"owner" binding:"required"`
	Repo  string `json:"repo" binding:"required"`
}

func StartScanHandler(c *gin.Context) {
	var req ScanRequest

	// Bind the JSON body to the struct
	if err := c.BindJSON(&req); err != nil {
		// If there is an error, return a 400 Bad Request response
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	// Remove the "Bearer " prefix from the token
	token = token[len("Bearer "):]

	var repository models.Repo
	result := db.DB.Where("owner = ? AND name = ?", req.Owner, req.Repo).First(&repository)

	fmt.Println(result.Error)

	if result.Error == gorm.ErrRecordNotFound {
		fmt.Println("Creating new repo")
		repository = models.Repo{
			Owner:        req.Owner,
			Name:         req.Repo,
			ScanComplete: false,
		}
		db.DB.Create(&repository)
	} else {
		fmt.Println("Updating repo")
		db.DB.Model(&repository).Update("scan_complete", false)
	}

	lib.ProcessRepoFiles(req.Owner, req.Repo, token)

	c.JSON(http.StatusOK, gin.H{
		"status": "Scan started",
	})
}

func GetScanStatusHandler(c *gin.Context) {

	var req ScanRequest

	// Bind the JSON body to the struct
	if err := c.BindJSON(&req); err != nil {
		// If there is an error, return a 400 Bad Request response
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	var repo models.Repo
	result := db.DB.Where("owner = ? AND name = ?", req.Owner, req.Repo).First(&repo)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Repository not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"scanComplete": repo.ScanComplete,
	})
}

func GetScanResultsHandler(c *gin.Context) {
	var req ScanRequest

	// Bind the JSON body to the struct
	if err := c.BindJSON(&req); err != nil {
		// If there is an error, return a 400 Bad Request response
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	var repo models.Repo
	result := db.DB.Where("owner = ? AND name = ?", req.Owner, req.Repo).First(&repo)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Repository not found",
		})
		return
	}

	var languages []models.Language

	error := db.DB.
		Joins("JOIN repo_languages ON languages.id = repo_languages.language_id").
		Where("repo_languages.repository_id = ? AND languages.supported = ?", repo.ID, true).
		Find(&languages).Error

	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch scan results",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"languages": languages,
	})
}

func GetInstallHandler(c *gin.Context) {
	var req ScanRequest

	// Bind the JSON body to the struct
	if err := c.BindJSON(&req); err != nil {
		// If there is an error, return a 400 Bad Request response
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	// Remove the "Bearer " prefix from the token
	token = token[len("Bearer "):]

	lib.MakeAction(req.Owner, req.Repo, token)

	c.JSON(http.StatusOK, gin.H{
		"status": "Installation started",
	})
}
