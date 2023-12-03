package controllers

import (
	"NotVulnApp/db"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"html/template"
	"net/http"
	"path/filepath"
)

type UserController struct{}

func (u *UserController) Profile(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")

	favoriteFilename := ""
	err := db.GetDB().QueryRow("SELECT filename FROM users WHERE username = ?", username).Scan(&favoriteFilename)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Ошибка базы данных"})
		return
	}

	if username == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	escapedUsername := template.HTMLEscapeString(username.(string))
	c.HTML(http.StatusOK, "profile.html", gin.H{"username": template.HTML(escapedUsername), "filename": favoriteFilename})
}

func (u *UserController) Delete(c *gin.Context) {
	username := c.Query("user")

	session := sessions.Default(c)
	cookieUsername := session.Get("username")

	if cookieUsername != username || cookieUsername == nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Ошибка"})
		return
	}

	_, err := db.GetDB().Exec("DELETE FROM users WHERE username = ?", username)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Ошибка базы данных"})
		return
	}

	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/")
}

func (u *UserController) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Ошибка при получении файла")
		return
	}

	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Ошибка"})
		return
	}

	fileExt := filepath.Ext(file.Filename)
	rndFilename := uuid.NewString() + fileExt

	err = c.SaveUploadedFile(file, "./uploads/"+rndFilename)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Ошибка при сохранении файла"})
		return
	}

	_, err = db.GetDB().Exec("UPDATE users SET filename = ? WHERE username = ?", rndFilename, username)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Ошибка базы данных"})
		return
	}

	c.Redirect(http.StatusFound, "/profile")
}
