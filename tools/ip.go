package tools

import (
	"os/exec"

	"github.com/gin-gonic/gin"
)

func GetIP(e *gin.Engine) {
	e.GET("/ip", getip)
}

func getip(r *gin.Context) {
	ip := r.Request.Header.Get("X-Forwarded-For")
	cmd := exec.Command("python", "./getip.py", ip)
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	address := string(out)
	r.String(200, address+""+ip)
}
