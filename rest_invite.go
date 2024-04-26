package main

import "github.com/gin-gonic/gin"

func invite(c *gin.Context) {
	if !authenticateBackend(c) {
		return
	}

	code := getNewInviteCode()

	c.JSON(201, gin.H{
		"inviteCode": code,
	})
}
