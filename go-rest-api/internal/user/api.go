package user

import (
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/qiangxue/go-rest-api/internal/errors"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"net/http"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service UserService, logger log.Logger) {
	res := resource{service, logger}

	r.Get("/user/<email>", res.get)
	r.Post("/user", res.create)
	r.Put("/user/<email>", res.update)
	r.Delete("/user/<email>", res.delete)
	r.Post("/userSignup/<email>/<code>", res.UserSignUp)
	r.Get("/userEmailConfirm/<email>/<code>", res.AuthenticateUser)
}

type resource struct {
	service UserService
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	user, err := r.service.GetUser(c.Request.Context(), c.Param("email"))
	if err != nil {
		return errors.NotFound("user doesn't exists")
	}
	return c.Write(user)
}

func (r resource) create(c *routing.Context) error {
	var input CreateUserRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	user, err := r.service.CreateUser(c.Request.Context(), input)
	if err != nil {
		return c.WriteWithStatus(err, http.StatusInternalServerError)
	}

	return c.WriteWithStatus(user, http.StatusCreated)
}


func (r resource) update(c *routing.Context) error {
	var input UpdateUserRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	user, err := r.service.UpdateUser(c.Request.Context(), c.Param("email"), input, false)
	if err != nil {
		return err
	}

	return c.Write(user)
}

func (r resource) delete(c *routing.Context) error {
	var input DeleteUserRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	user, err := r.service.DeleteUser(c.Request.Context(), c.Param("email"), input)
	if err != nil {
		return err
	}

	return c.Write(user)
}


func (r resource) UserSignUp(c *routing.Context) error {
	user, _, err := r.service.UserSignUp(c.Request.Context(), c.Param("email"), c.Param("code"))
	if err != nil {
		return err
	}
	return c.Write(user)
}


func (r resource) AuthenticateUser(c *routing.Context) error {
	user, err := r.service.AuthenticateUser(c.Request.Context(), c.Param("email"), c.Param("code"))
	if err != nil {
		return err
	}
	return c.Write(user)
}
