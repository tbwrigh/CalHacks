package lib

import (
	"bytes"
	"calhacks/api/db"
	"calhacks/api/models"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GitHubUser struct {
	Login  string `json:"login"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar_url"`
}

type GitHubRepo struct {
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Description   string `json:"description"`
	Language      string `json:"language"`
	DefaultBranch string `json:"default_branch"`
	Owner         struct {
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

func GetRepo(owner, repo, token string) (*GitHubRepo, error) {
	client := &http.Client{}

	// Create the GitHub API request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo), nil)
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
	var repoData GitHubRepo
	if err := json.NewDecoder(resp.Body).Decode(&repoData); err != nil {
		return nil, err
	}

	return &repoData, nil
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

type BranchRef struct {
	Ref string `json:"ref"`
	SHA string `json:"sha"`
}

// Function to get the latest commit SHA of a branch
func GetBranchSHA(owner, repo, branch, token string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs/heads/%s", owner, repo, branch)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch branch sha: %s", resp.Status)
	}

	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	sha := result["object"].(map[string]interface{})["sha"].(string)
	return sha, nil
}

// Function to create a new branch based on the SHA of an existing branch
func CreateBranch(owner, repo, newBranch, sha, token string) error {
	fmt.Println("Creating branch...")
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs", owner, repo)

	newRef := BranchRef{
		Ref: "refs/heads/" + newBranch,
		SHA: sha,
	}

	fmt.Println("Created new ref")
	fmt.Println(newRef)

	body, err := json.Marshal(newRef)
	if err != nil {
		return err
	}

	fmt.Println("Marshalled body")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	fmt.Println("Created request")

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "AutoLockAI")

	fmt.Println("Set headers")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	fmt.Println("Sent request")

	defer resp.Body.Close()

	fmt.Println("Closed response body")

	if resp.StatusCode != 201 {
		return fmt.Errorf("failed to create branch: %s", resp.Status)
	}

	fmt.Println("Created branch")

	return nil
}

func GeneratePRBranch(owner, repo, token string) (string, error) {
	fmt.Println("Generating PR branch...")

	// get repo info
	repoData, err := GetRepo(owner, repo, token)
	if err != nil {
		return "", err
	}

	fmt.Println("Got repo info")

	fmt.Println("Getting latest commit SHA...")

	// Get the latest commit SHA of the main branch
	sha, err := GetBranchSHA(owner, repo, repoData.DefaultBranch, token)
	if err != nil {
		return "", err
	}

	fmt.Println("Got latest commit SHA")

	fmt.Println("Creating new branch...")

	// Create a new branch based on the main branch
	newBranch := "autolock-" + uuid.New().String()
	if err := CreateBranch(owner, repo, newBranch, sha, token); err != nil {
		return "", err
	}

	fmt.Println("Created new branch")

	return newBranch, nil
}

func getFileSHA(url, accessToken string) (string, error) {
	// Set up the HTTP GET request to fetch the file
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create GET request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// If the file exists, return its SHA
	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		var fileData map[string]interface{}
		if err := json.Unmarshal(body, &fileData); err != nil {
			return "", fmt.Errorf("failed to parse GET response: %v", err)
		}
		if sha, exists := fileData["sha"].(string); exists {
			return sha, nil
		}
	}

	// If the file does not exist, return an empty SHA for new files
	return "", nil
}

type FileChange struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Message string `json:"message"`
	NewFile bool   `json:"new_file"`
}

type FileContentRequest struct {
	Message string `json:"message"`
	Content string `json:"content"`
	Branch  string `json:"branch"`
	Sha     string `json:"sha,omitempty"` // Optional SHA for updating
}

func CreateCommit(owner string, repo string, branch string, change FileChange, token string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, change.Path)

	encodedContent := base64.StdEncoding.EncodeToString([]byte(change.Content))

	reqBody := FileContentRequest{
		Message: change.Message,
		Content: encodedContent,
		Branch:  branch,
	}

	if !change.NewFile {
		sha, err := getFileSHA(url, token)
		if err != nil {
			return err
		}
		reqBody.Sha = sha
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send PUT request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the file was created or updated successfully
	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		fmt.Println("File created/updated successfully!")
		return nil
	}

	// If the request failed, read the error response
	body_resp, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("failed to create/update file: %s, %s", resp.Status, string(body_resp))

}

func PRBranch(owner string, repo string, branch string, token string) error {
	repository, err := GetRepo(owner, repo, token)
	if err != nil {
		return fmt.Errorf("failed to get repository")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repo)

	reqBody := map[string]interface{}{
		"title": "Auto-locking dependencies",
		"head":  branch,
		"base":  repository.DefaultBranch,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the PR was created successfully
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("PR created successfully!")
		return nil
	}

	// If the request failed, read the error response
	body_resp, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("failed to create PR: %s, %s", resp.Status, string(body_resp))
}

func CreatePR(owner string, repo string, changes []FileChange, token string) error {
	// step 1 make a new branch
	fmt.Println("Creating PR branch...")
	newBranch, err := GeneratePRBranch(owner, repo, token)
	if err != nil {
		return err
	}

	fmt.Println("Created PR branch")

	fmt.Println("Creating commits...")

	// step 2 create a new commit
	for _, change := range changes {
		err := CreateCommit(owner, repo, newBranch, change, token)
		if err != nil {
			return err
		}
	}

	fmt.Println("Created commits")

	fmt.Println("Creating PR...")

	// step 3 create a new PR
	err = PRBranch(owner, repo, newBranch, token)
	if err != nil {
		return err
	}

	fmt.Println("Created PR")

	return nil
}

func MakeAction(owner, repo, token string) error {
	// create PR in go routine
	go func() {
		fmt.Println("Creating PR...")

		_, filePath, _, _ := runtime.Caller(0) // Get the current file's directory
		dir := filepath.Dir(filePath)

		fmt.Println("Set dir")

		path := filepath.Join(dir, "resources/autolock.yml")

		fmt.Println("Set path")
		fmt.Println(path)

		// Read the file contents
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("Error reading file")
			fmt.Println(err)
			return
		}

		fmt.Println("Read file")

		// replace "%LANGUAGES%" in content with the supported languages
		var languages []models.Language
		result := db.DB.Joins("JOIN repo_languages ON languages.id = repo_languages.language_id").
			Joins("JOIN repos ON repos.id = repo_languages.repository_id").
			Where("repos.owner = ? AND repos.name = ? AND languages.supported = ?", owner, repo, true).
			Find(&languages)

		if result.Error != nil {
			return
		}

		fmt.Println("Got languages")

		// replace the placeholder with the supported languages
		langs := make([]string, len(languages))
		for i, lang := range languages {
			langs[i] = "'" + strings.ToLower(lang.Name) + "'"
		}

		fmt.Println("Got langs")

		content = []byte(strings.ReplaceAll(string(content), "%LANGUAGES%", strings.Join(langs, ",")))

		fmt.Println("Replaced content")

		fileChange := FileChange{
			Path:    ".github/workflows/autolock.yml",
			Content: string(content),
			Message: "Add autolock workflow",
			NewFile: true,
		}

		fmt.Println("Created file change")

		CreatePR(owner, repo, []FileChange{fileChange}, token)

		fmt.Println("Created PR")

		// update the repo to have install started
		var repository models.Repo
		result = db.DB.Where("owner = ? AND name = ?", owner, repo).First(&repository)
		if result.Error != nil {
			return
		}
		db.DB.Model(&repository).Update("install_started", true)
	}()

	return nil
}
