package controllers

import (
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"super-supply-chain/configs"
	"super-supply-chain/utils"
)

func LoadStatic(c *gin.Engine) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if !configs.IsDev() {
		staticPath := filepath.Join(filepath.Dir(pwd), "frontend/dist")
		utils.Logger.Info("Load static path: " + staticPath)
		c.Static("/super-supply-chain", staticPath)
	}
}
