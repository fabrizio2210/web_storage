package routes

import (
  "Api/controllers"

  "github.com/gin-gonic/gin"
)

var StoragePath = "storage"

func StorageRoute(router *gin.Engine) {
  router.DELETE(StoragePath + "/*objectId", controllers.DeleteObject())
  router.GET(   StoragePath + "/*objectId", controllers.GetObject())
  router.POST(  StoragePath + "/*objectId", controllers.PostObject())
}
