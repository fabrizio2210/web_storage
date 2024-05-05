package routes

import (
  "Api/controllers"

  "github.com/gin-gonic/gin"
)

func StorageRoute(router *gin.Engine) {
  // router.DELETE("/api/photo/:photoId", controllers.DeleteObject())
  router.GET(   "/storage/:objectId", controllers.GetObject())
  router.POST(  "/storage/:objectId",      controllers.PostNewObject())
}
