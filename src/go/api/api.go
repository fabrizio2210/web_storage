package main

import (
    "os"

    "Api/routes"
    "Api/rediswrapper"

    "github.com/gin-gonic/gin"
)


func setupRouter() *gin.Engine {
  r := gin.Default()
  if os.Getenv("STORAGE_PATH") != "" {
    routes.StoragePath = os.Getenv("STORAGE_PATH")
  }
  routes.StorageRoute(r)
  return r
}

func main() {
  rediswrapper.RedisClient = rediswrapper.ConnectRedis(os.Getenv("REDIS_HOST") + ":6379")
  router := setupRouter()
  router.Run("0.0.0.0:5000")
}
