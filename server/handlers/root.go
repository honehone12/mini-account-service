package handlers

import (
	"mini-account-service/db/controller"
	"mini-account-service/server/context"
	"mini-account-service/server/quick"
	"mini-account-service/server/session"
	"mini-account-service/server/session/onetime"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RegistrationForm struct {
	Email    string `form:"email" validate:"required,email,max=64"`
	Password string `form:"password" validate:"required,alphanum,min=8,max=64"`
}

type LoginForm struct {
	Email    string `form:"email" validate:"required,email,max=64"`
	Password string `form:"password" validate:"required,alphanum,min=8,max=64"`
}

type RootResponse struct {
	Name    string
	Version string
}

type AuthorizeResponse struct {
	Uuid      string
	OneTimeId string
}

func Root(c echo.Context) error {
	ctx := c.(*context.Context)
	return c.JSON(http.StatusOK, RootResponse{
		Name:    ctx.Name,
		Version: ctx.Version,
	})
}

func Register(c echo.Context) error {
	formData := &RegistrationForm{}
	if err := quick.ProcessFormData(c, formData); err != nil {
		c.Logger().Warn(err)
		return quick.BadRequest()
	}

	cc := c.(*context.Context)
	ctrl := cc.User()
	uuid, err := ctrl.Create(formData.Email, formData.Password)
	if err != nil {
		c.Logger().Error(err)
		return quick.ServiceError()
	}

	//
	// or move data to redis when the user logged in
	formValues := url.Values{"uuid": {uuid}}
	url := cc.GamedataService + "/jewel/init"
	res, err := http.PostForm(url, formValues)
	if err != nil {
		c.Logger().Error(err)
		return quick.ServiceError()
	} else if res.StatusCode != http.StatusOK {
		c.Logger().Error(quick.ErrorStatusCodeNotOk(res.StatusCode))
		return quick.ServiceError()
	}
	defer res.Body.Close()
	//

	return c.NoContent(http.StatusOK)
}

func Login(c echo.Context) error {
	formData := &LoginForm{}
	if err := quick.ProcessFormData(c, formData); err != nil {
		c.Logger().Warn(err)
		return quick.BadRequest()
	}

	ctrl := c.(*context.Context).User()
	uuid, err := ctrl.VerifyPassword(formData.Email, formData.Password)
	if err == gorm.ErrRecordNotFound || err == controller.ErrorInvalidPassword {
		c.Logger().Error(err)
		return quick.BadRequest()
	} else if err != nil {
		c.Logger().Error(err)
		return quick.ServiceError()
	}

	err = session.Set(c, uuid)
	if err != nil {
		c.Logger().Error(err)
		return quick.ServiceError()
	}
	return c.NoContent(http.StatusOK)
}

func Authorize(c echo.Context) error {
	sess, err := session.RequireSession(c)
	if err != nil {
		c.Logger().Warn(err)
		return quick.BadRequest()
	}

	onetimeId, err := onetime.New()
	if err != nil {
		c.Logger().Error(err)
		return quick.ServiceError()
	}

	formValues := url.Values{
		"uuid": {sess.Uuid},
		"id":   {onetimeId},
	}
	url := c.(*context.Context).GamedataService + "/session/set"
	res, err := http.PostForm(url, formValues)
	if err != nil {
		c.Logger().Error(err)
		return quick.ServiceError()
	} else if res.StatusCode != http.StatusOK {
		c.Logger().Error(quick.ErrorStatusCodeNotOk(res.StatusCode))
		return quick.ServiceError()
	}
	defer res.Body.Close()

	return c.JSON(http.StatusOK, AuthorizeResponse{
		Uuid:      sess.Uuid,
		OneTimeId: onetimeId,
	})
}
