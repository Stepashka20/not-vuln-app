package controllers

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

type MiscController struct{}

func (m *MiscController) Ping(c *gin.Context) {
	website := c.PostForm("website")
	if !isWebsiteValid(website) {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Ошибка"})
		return
	}
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "4", website)
	} else {
		cmd = exec.Command("ping", "-c", "4", website)
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

func isWebsiteValid(website string) bool {
	domainRegex := regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,3})$`)
	return domainRegex.MatchString(website) || net.ParseIP(website) != nil
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
