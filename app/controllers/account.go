package controllers

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AccountDetails(c *gin.Context) {

	cf := c.MustGet("cloudflare").(*cloudflare.API)
	var pageOpts cloudflare.PaginationOptions
	pageOpts.Page = 1
	pageOpts.PerPage = 100
	accounts, info, _ := cf.Accounts(c, pageOpts)
	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
		"info":     info,
	})
}
func UserDetails(c *gin.Context) {

	cf := c.MustGet("cloudflare").(*cloudflare.API)
	userDetails, _ := cf.UserDetails(c)
	c.JSON(http.StatusOK, gin.H{
		"userDetails": userDetails,
	})

}