package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func WebStart() {
	router := gin.Default()

	router.GET("/list", listHandle)
	router.GET("/get_status/:uid", getStatusHandle)
	router.GET("/get_apdata/:uid", getApDataHandle)
	router.GET("/get_devinfo/:uid", getDevInfoHandle)
	router.GET("/reboot/:uid", rebootHandle)
	router.GET("/reset/:uid", resetHandle)
	router.POST("/switch_config/:uid", switchConfigHandle)
	router.POST("/wifi_config/:uid", wifiConfigHandle)
	router.POST("/upgrade_config/:uid", upgradeConfigHandle)

	log.Println(router.Run(":8000"))
}

func listHandle(ctx *gin.Context) {
	res := make([]*ElinkSession, 0)
	SessionMap.Range(func(key, value interface{}) bool {
		res = append(res, value.(*ElinkSession))
		return true
	})

	ctx.JSON(http.StatusOK, res)
}

func getStatusHandle(ctx *gin.Context) {
	uid := ctx.Param("uid")

	if val, ok := SessionMap.Load(uid); ok {
		err := val.(*ElinkSession).getStatus()
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": ""})
		}
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "not found the session"})
	}
}

func getApDataHandle(ctx *gin.Context) {
	uid := ctx.Param("uid")

	if val, ok := SessionMap.Load(uid); ok {
		err := val.(*ElinkSession).getAPData()
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": ""})
		}
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "not found the session"})
	}
}

func getDevInfoHandle(ctx *gin.Context) {
	uid := ctx.Param("uid")

	if val, ok := SessionMap.Load(uid); ok {
		err := val.(*ElinkSession).getDevInfo()
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": ""})
		}
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "not found the session"})
	}
}

func switchConfigHandle(ctx *gin.Context) {
	uid := ctx.Param("uid")

	cfg := SwitchSet{}
	err := ctx.BindJSON(&cfg)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": "json decode error"})
		return
	}

	if val, ok := SessionMap.Load(uid); ok {
		err := val.(*ElinkSession).switchConfig(cfg)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": ""})
		}
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "not found the session"})
	}
}
func wifiConfigHandle(ctx *gin.Context) {
	uid := ctx.Param("uid")

	cfg := WiFiSet{}
	err := ctx.BindJSON(&cfg)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": "json decode error"})
		return
	}

	if val, ok := SessionMap.Load(uid); ok {
		err := val.(*ElinkSession).wifiConfig(cfg)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": ""})
		}
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "not found the session"})
	}
}

func upgradeConfigHandle(ctx *gin.Context) {
	uid := ctx.Param("uid")

	cfg := UpgradeSet{}
	err := ctx.BindJSON(&cfg)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": "json decode error"})
		return
	}

	if val, ok := SessionMap.Load(uid); ok {
		err := val.(*ElinkSession).upgradeConfig(cfg)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": ""})
		}
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "not found the session"})
	}
}

func rebootHandle(ctx *gin.Context) {
	uid := ctx.Param("uid")

	if val, ok := SessionMap.Load(uid); ok {
		err := val.(*ElinkSession).reboot()
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": ""})
		}
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "not found the session"})
	}
}

func resetHandle(ctx *gin.Context) {
	uid := ctx.Param("uid")

	if val, ok := SessionMap.Load(uid); ok {
		err := val.(*ElinkSession).reset()
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": ""})
		}
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "not found the session"})
	}
}
