package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp" // Required for filename sanitization
	"time"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday" // Recommended for HTML sanitization
)

// ArticleRequest defines the structure for the incoming JSON payload.
type ArticleRequest struct {
	HTMLContent string `json:"htmlContent"`
	ArticleName string `json:"articleName"` // Field for the custom filename
}

func main() {
	router := gin.Default()

	// Serve the static HTML file from the "static" directory (e.g., index.html)
	// Access via http://localhost:8080/
	router.StaticFile("/", "./static/index.html")
	router.StaticFS("/static", http.Dir("./static"))

	// Serve the saved HTML articles from the "articles" directory.
	// Access articles via http://localhost:8080/articles/your-article-name.html
	router.StaticFS("/articles", http.Dir("./articles"))

	// API endpoint to save the article
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
	fullHTMLContent := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <link href="https://cdn.jsdelivr.net/npm/quill@2.0.2/dist/quill.snow.css" rel="stylesheet">
    <link href="https://fonts.googleapis.com/css2?family=Montserrat:ital,wght@0,100..900;1,100..900&family=Lato:ital,wght@0,100;0,300;0,400;0,700;0,900;1,100;1,300;1,400;1,700;1,900&display=swap" rel="stylesheet">
    <style>
        body {
            font-family: 'Montserrat', sans-serif;
            padding: 2rem;
            background-color: #f0f2f5;
        }
        .article-content {
            max-width: 800px;
            margin: auto;
            padding: 2rem;
            background-color: #fff;
            border-radius: 8px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
    </style>
</head>
<body>
    <div class="article-content">
        %s
    </div>
</body>
</html>`, req.ArticleName, sanitizedHTMLContent) // Use original ArticleName for title, sanitized HTML for content

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
