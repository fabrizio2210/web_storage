package controllers

import (
  "context"
  "fmt"
  "log"
  "encoding/json"
  "net/http"
  "io/ioutil"
  "os"
  "time"

  "Api/models"
  "Api/responses"
  "Api/rediswrapper"

  "github.com/gin-gonic/gin"
  "github.com/go-playground/validator/v10"
)


type StoredObject struct {
  Bytes []byte
  ContentType string
  RequireAuthentication bool
}

var validate = validator.New()

func GetObject() gin.HandlerFunc {
  return func(c *gin.Context) {
    _, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectId := c.Param("objectId")
    marshalledObject := rediswrapper.Get(objectId)
    if marshalledObject == nil {
      c.JSON(
        http.StatusNotFound,
        responses.Response{
          Status: http.StatusNotFound,
        },
      )
      return
    }
    var object StoredObject
    err := json.Unmarshal(marshalledObject, &object)
    if err != nil {
      log.Printf("Error in unmarshalling: %s", err)
      c.JSON(
        http.StatusBadRequest,
        responses.Response{
          Status: http.StatusBadRequest,
          Message: fmt.Sprintf("Error in unmarshalling: %s", err),
        },
      )
      return
    }

    if object.RequireAuthentication {
      if (len(c.Request.Header["Api-Key"]) == 0) || (c.Request.Header["Api-Key"][0] != os.Getenv("API_KEY")) {
        c.JSON(
          http.StatusUnauthorized,
          responses.Response{
            Status: http.StatusUnauthorized,
            Message: "Not authorized.",
          },
        )
        return
      }
    }

    c.Data(http.StatusOK, object.ContentType, object.Bytes)
  }
}

func PostObject() gin.HandlerFunc {
  return func(c *gin.Context) {
    _, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Do not authorize if wrong API_KEY.
    if (len(c.Request.Header["Api-Key"]) == 0) || (c.Request.Header["Api-Key"][0] != os.Getenv("API_KEY")) {
      c.JSON(
        http.StatusUnauthorized,
        responses.Response{
          Status: http.StatusUnauthorized,
          Message: "Not authorized.",
        },
      )
      return
    }

    objectId := c.Param("objectId")

    var data models.ObjectInputForm
    if (! maybeGetForm(c, &data)) {
      log.Printf("Wrong parsing of the post data.")
      return
    }

    form, err := c.MultipartForm()
    if err != nil {
      log.Printf("No multipart form found: %v", err.Error())
      c.JSON(
        http.StatusBadRequest,
        responses.Response{
          Status: http.StatusBadRequest,
          Message: "Error: no multipart found in the request.",
        },
      )
      return
    }

    files := form.File["file"]
    if (len(files) != 1) {
      log.Printf("Number of files is different from 1: %d", len(files))
      c.JSON(
        http.StatusBadRequest,
        responses.Response{
          Status: http.StatusBadRequest,
          Message: fmt.Sprintf("Error: expected 1 file, found %d in the request.", len(files)),
        },
      )
      return
    }

    if data.Override == false && rediswrapper.Get(objectId) != nil {
      c.JSON(
        http.StatusBadRequest,
        responses.Response{
          Status: http.StatusBadRequest,
          Message: "Error: trying to override an " +
                   "object when the \"override\" parameter is false.",
        },
      )
      return
    }

    fl, _ := files[0].Open()
    binary, _ := ioutil.ReadAll(fl)
    object := StoredObject{
      Bytes: binary,
      ContentType: data.ContentType,
      RequireAuthentication: data.RequireAuthentication,
    }
    marshalledObject, err := json.Marshal(object)
    if err != nil {
      c.JSON(
        http.StatusBadRequest,
        responses.Response{
          Status: http.StatusBadRequest,
          Message: fmt.Sprintf("Error in marshalling: %s.", err),
        },
      )
       return
    }
    rediswrapper.Store(objectId, marshalledObject)
    log.Printf("Object upload:%v", objectId)
  }
}

func maybeGetForm(c *gin.Context, data any) bool {
  if err := c.Bind(data); err != nil {
    log.Printf("Error in parsing: %v", err.Error())
    c.JSON(
      http.StatusBadRequest,
      responses.Response{
        Status: http.StatusBadRequest,
        Message: "error",
        Data: map[string]interface{}{"event": err.Error()},
      },
    )
    return false
  }

  if err := validate.Struct(data); err != nil {
    log.Printf("Error in validation: %v", err.Error())
    c.JSON(
      http.StatusBadRequest,
      responses.Response{
        Status: http.StatusBadRequest,
        Message: "error",
        Data: map[string]interface{}{"event": err.Error()},
      },
    )
    return false
  }
  return true
}

func DeleteObject() gin.HandlerFunc {
  return func(c *gin.Context) {
    _, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Do not authorize if wrong API_KEY.
    if (len(c.Request.Header["Api-Key"]) == 0) || (c.Request.Header["Api-Key"][0] != os.Getenv("API_KEY")) {
      c.JSON(
        http.StatusUnauthorized,
        responses.Response{
          Status: http.StatusUnauthorized,
          Message: "Not authorized.",
        },
      )
      return
    }

    objectId := c.Param("objectId")

    rediswrapper.Delete(objectId)
    log.Printf("Object deleted:%v", objectId)
  }
}
