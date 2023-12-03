package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type MiscController struct{}

func (m *MiscController) Ping(c *gin.Context) {
	website := c.PostForm("website")

	var userInput string
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		userInput = "ping -n 4 " + website
		cmd = exec.Command("cmd", "/C", userInput)
	} else {
		userInput = "ping -c 4 " + website
		cmd = exec.Command("sh", "-c", userInput)
	}
	output, err := cmd.CombinedOutput()

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": "Ошибка ping"})
		return
	}

	result := string(output)

	c.HTML(200, "ping.html", gin.H{
		"Result": result,
	})
}

func (m *MiscController) Files(c *gin.Context) {
	filename := c.Query("name")
	if filename == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Ошибка"})
		return
	}

	basePath := "./uploads/"
	cleanedPath := CleanPath(filepath.Join(basePath, filename))

	if !isFileWithinDirectory(cleanedPath, basePath) {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Ошибка"})
		return
	}

	c.File(cleanedPath)
}

func isFileWithinDirectory(filePath, basePath string) bool {
	relPath, err := filepath.Rel(basePath, filePath)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(relPath, ".."+string(filepath.Separator))
}

func CleanPath(path string) string {
	if path == "" {
		return ""
	}
	path = filepath.Clean(path)

	if !filepath.IsAbs(path) {
		path = filepath.Clean(string(os.PathSeparator) + path)
		path, _ = filepath.Rel(string(os.PathSeparator), path)
	}

	return filepath.Clean(path)
}
