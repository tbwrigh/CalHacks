package handlers

import (
	"net/http"

	"calhacks/api/db"
	"calhacks/api/lib"
	"calhacks/api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MeHandler handles the /me endpoint, expects an Authorization Bearer token
func MeHandler(c *gin.Context) {
	// Get the Authorization header
	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	// Remove the "Bearer " prefix from the token
	token = token[len("Bearer "):]

	githubUser, err := lib.GetGitHubUser(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info from GitHub"})
		return
	}

	// Check if a user exists in the database with this email
	var user models.User
	result := db.DB.Where("github_username = ?", githubUser.Login).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// If user does not exist, insert a new user
			user = models.User{
				GithubUsername: githubUser.Login,
				HasAccess:      false,
			}
			db.DB.Create(&user)
		} else {
			// If some other database error occurs
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	// Return the user (either found or newly created)
	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"user":      user.GithubUsername,
		"hasAccess": user.HasAccess,
	})
}

func RepoHandler(c *gin.Context) {
	// Get the Authorization header
	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	// Remove the "Bearer " prefix from the token
	token = token[len("Bearer "):]

	repos, err := lib.GetUserRepos(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user repos from GitHub"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"repos": repos,
	})
}
