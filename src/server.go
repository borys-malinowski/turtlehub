package main

import (
	_ "database/sql"
	"fmt"
	_ "github.com/bdwilliams/go-jsonify/jsonify"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/microcosm-cc/bluemonday"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func main() {
	router := gin.Default()
	go router.Use(static.Serve("/", static.LocalFile("public", false)))
	go router.GET("/api/yt-downloader", ytDownloader)
	router.Run()
}

func ytDownloader(context *gin.Context) {
	sanitizer := bluemonday.UGCPolicy()
	link := sanitizer.Sanitize(context.Query("link"))
	uniqueName := randomString(20, charset)
	command := exec.Command("annie", "-o", "./videos", "-O", uniqueName, link)
	command.Run()
	context.Header("Content-Disposition", "attachment; filename="+uniqueName+".mp4")
	context.File("./videos/" + uniqueName + ".mp4")
	err := os.Remove("./videos/" + uniqueName + ".mp4")
	checkError(err)
}

func randomString(length uint, charset string) string {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	uniqueName := make([]byte, length)
	for i := range uniqueName {
		uniqueName[i] = charset[seed.Intn(len(charset))]
	}
	return string(uniqueName)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
