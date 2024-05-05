package models

type ObjectInputForm struct {
  Override bool `form:"override"`
  ContentType string `form:"content-type" validate:"required"`
  RequireAuthentication bool `form:"require-authentication"`
}

