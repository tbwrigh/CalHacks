package lib

import (
	"calhacks/api/db"
	"calhacks/api/models"

	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type GitHubUser struct {
	Login  string `json:"login"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar_url"`
}

type GitHubRepo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Owner       struct {
		Login string `json:"login"`
	} `json:"owner"`
}

type File struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

func GetGitHubUser(token string) (*GitHubUser, error) {
	client := &http.Client{}

	// Create the GitHub API request
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	// Add the Authorization header with the Bearer token
	req.Header.Add("Authorization", "Bearer "+token)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	// Parse the response body
	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserRepos(token string) ([]GitHubRepo, error) {
	client := &http.Client{}

	// Create the GitHub API request
	req, err := http.NewRequest("GET", "https://api.github.com/user/repos", nil)
	if err != nil {
		return nil, err
	}

	// Add the Authorization header with the Bearer token
	req.Header.Add("Authorization", "Bearer "+token)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	// Parse the response body
	var repos []GitHubRepo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return repos, nil
}

func FetchRepoFiles(owner, repo, token, path string) ([]File, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add the access token to the request headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch files: %s", resp.Status)
	}

	var files []File
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}

	return files, nil
}

func GetFileExtensionsRecursive(owner, repo, token, path string, extensions map[string]bool) error {
	files, err := FetchRepoFiles(owner, repo, token, path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.Type == "file" {
			// Extract the file extension
			if dotIdx := strings.LastIndex(file.Name, "."); dotIdx != -1 {
				ext := strings.ToLower(file.Name[dotIdx+1:])
				extensions[ext] = true
			}
		} else if file.Type == "dir" {
			// Recursively fetch files in the subdirectory
			if err := GetFileExtensionsRecursive(owner, repo, token, file.Path, extensions); err != nil {
				return err
			}
		}
	}

	return nil
}

// ProcessRepoFiles processes the GitHub repository files in a subroutine, updating the database.
func ProcessRepoFiles(owner, repo, token string) {
	// Use a goroutine to handle this asynchronously
	go func() {
		extensions := make(map[string]bool)
		var wg sync.WaitGroup

		// Start fetching files recursively
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := GetFileExtensionsRecursive(owner, repo, token, "", extensions)
			if err != nil {
				fmt.Println("Error fetching repo files:", err)
				return
			}
		}()

		// Wait for all the file fetching to be complete
		wg.Wait()

		// Process each extension and update the database
		for ext := range extensions {
			var language models.Language
			result := db.DB.Where("extension = ?", ext).First(&language)

			// If the language does not exist, insert it
			if result.Error == gorm.ErrRecordNotFound {
				newLanguage := models.Language{
					Name:      ext,
					Extension: ext,
					Supported: false,
				}
				db.DB.Create(&newLanguage)
			}

			var repository models.Repo
			result = db.DB.Where("owner = ? AND name = ?", owner, repo).First(&repository)

			if result.Error != nil {
				fmt.Println("Error fetching repo:", result.Error)
				return
			}

			// Check if the repository-language relationship exists
			var repoLanguage models.RepoLanguage
			result = db.DB.Where("repository_id = ? AND language_id = ?", repository.ID, language.ID).First(&repoLanguage)

			// If the relationship does not exist, insert it
			if result.Error == gorm.ErrRecordNotFound {
				repoLanguage = models.RepoLanguage{
					RepositoryID: repository.ID,
					LanguageID:   language.ID,
				}
				db.DB.Create(&repoLanguage)
			}
		}

		// After processing, mark the repository as having the scan complete
		db.DB.Model(&models.Repo{}).Where("owner = ? AND name = ?", owner, repo).Update("scan_complete", true)

		fmt.Printf("Scan complete for repository: %s/%s\n", owner, repo)
	}()
}
