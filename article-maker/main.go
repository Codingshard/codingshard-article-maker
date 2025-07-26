package main

import (
	"fmt"           //fmt for formatting
	"log"           // log for logging and error debugging
	"net/http"      // net/http for HTTP server and request handling
	"os"            // os for file operations
	"path/filepath" // path/filepath for file path manipulation
	"regexp"        // Required for filename sanitization
	"time"          // time for generating unique timestamps

	"github.com/gin-gonic/gin"           // Gin framework for building web applications
	"github.com/microcosm-cc/bluemonday" // Recommended for HTML sanitization
)

// ArticleRequest defines the structure for the incoming JSON payload.
type ArticleRequest struct {
	HTMLContent string `json:"htmlContent"` // Field for the HTML content
	ArticleName string `json:"articleName"` // Field for the custom filename (user input)
}

func main() {
	router := gin.Default() // Create Gin router instance

	// Serve the static HTML file from the "static" directory
	router.StaticFile("/", "./static/index.html")
	router.StaticFS("/static", http.Dir("./static"))

	// Serve the saved HTML articles from the "articles" directory.
	// Access articles via http://localhost:8080/articles/your-article-name.html
	router.StaticFS("/articles", http.Dir("./articles"))

	// API POST endpoint to save the article.R
	router.POST("/save-article", saveArticleHandler)

	// Start the server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// saveArticleHandler processes the request to save the HTML content.
func saveArticleHandler(c *gin.Context) {
	var req ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// --- 1. Sanitize the incoming HTML content (HIGHLY RECOMMENDED FOR SECURITY) ---
	// Create a strict policy for user-generated content (UGC)
	p := bluemonday.UGCPolicy()
	sanitizedHTMLContent := p.Sanitize(req.HTMLContent)

	// --- 2. Ensure the output directory exists ---
	outputDir := "articles"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.Mkdir(outputDir, 0755); err != nil { // 0755 for rwx for owner, rx for others
			log.Printf("Failed to create directory %s: %v", outputDir, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article directory"})
			return
		}
	}

	// --- 3. Sanitize and generate the filename from ArticleName ---
	baseName := req.ArticleName
	if baseName == "" {
		// Fallback to a timestamp if no name is provided
		baseName = fmt.Sprintf("article-%d", time.Now().UnixNano())
	} else {
		// Replace non-alphanumeric/non-space characters with hyphens
		reNonAlphaNumeric := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
		sanitizedName := reNonAlphaNumeric.ReplaceAllString(baseName, "-")

		// Replace spaces with hyphens
		reSpaces := regexp.MustCompile(`\s+`)
		sanitizedName = reSpaces.ReplaceAllString(sanitizedName, "-")

		// Remove any leading/trailing hyphens
		sanitizedName = regexp.MustCompile(`^-|-$`).ReplaceAllString(sanitizedName, "")

		// Convert to lowercase for consistent URLs
		// sanitizedName = strings.ToLower(sanitizedName) // Uncomment if you want lowercase filenames

		// Append a unique timestamp to prevent collisions, even if names are similar
		baseName = fmt.Sprintf("%s-%d", sanitizedName, time.Now().UnixNano())
	}

	filename := baseName + ".html"
	filePath := filepath.Join(outputDir, filename)

	// --- 4. Construct the complete HTML document for saving ---
	articleTemplate, err := os.ReadFile("article-template.html")
	if err != nil {
		log.Printf("file reading failed with error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read article template"})
		return
	}
	fullHTMLContent := fmt.Sprintf(string(articleTemplate), req.ArticleName, sanitizedHTMLContent) // Use original ArticleName for title, sanitized HTML for content

	// --- 5. Write the HTML content to the file ---
	if err := os.WriteFile(filePath, []byte(fullHTMLContent), 0644); err != nil { // 0644 for rw for owner, r for others
		log.Printf("Failed to save article to %s: %v", filePath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save article content"})
		return
	}

	log.Printf("Article saved: %s", filePath)

	// --- 6. Send success response ---
	// You can also return the full URL to the article if you wish
	articleURL := fmt.Sprintf("/articles/%s", filename)
	c.JSON(http.StatusOK, gin.H{
		"message":    "Article published!",
		"filename":   filename,
		"articleURL": articleURL, // Added for convenience
	})
}
