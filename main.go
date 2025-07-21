package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin" // Removed cors import as it's no longer needed
)

// ArticleRequest defines the structure for the incoming HTML content.
type ArticleRequest struct {
	HTMLContent string `json:"htmlContent"`
}

func main() {
	router := gin.Default()

	// Serve the static HTML file from the "static" directory.
	// Make sure your index.html is inside a folder named "static"
	router.StaticFS("/", http.Dir("./static"))

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	outputDir := "articles"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.Mkdir(outputDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
			return
		}
	}

	filename := fmt.Sprintf("article-%d.html", time.Now().UnixNano())
	filePath := filepath.Join(outputDir, filename)

	// Construct a complete HTML document for the saved article
	fullHTMLContent := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Article - %s</title>
    <link href="https://cdn.jsdelivr.net/npm/quill@2.0.2/dist/quill.snow.css" rel="stylesheet">
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
</html>`, filename, req.HTMLContent)

	if err := os.WriteFile(filePath, []byte(fullHTMLContent), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save article"})
		return
	}

	log.Printf("Article saved: %s", filePath)
	c.JSON(http.StatusOK, gin.H{"message": "Article published!", "filename": filename})
}
