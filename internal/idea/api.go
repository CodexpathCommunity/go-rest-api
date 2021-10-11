package idea

import (
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/qiangxue/go-rest-api/internal/errors"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"net/http"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, logger log.Logger) {
	res := resource{service, logger}

	r.Post("/idea", res.create)
	r.Get("/idea/<id>", res.get)
	r.Put("/idea/<id>", res.update)
	r.Delete("/idea/<id>", res.delete)
	r.Post("/getIdeas", res.query)
	r.Post("/voteAnIdea", res.vote)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	idea, err := r.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(idea)
}

func (r resource) query(c *routing.Context) error {
	ctx := c.Request.Context()

	var input GetIdeaRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	ideas, err := r.service.Query(ctx, input)

	if err != nil {
		return err
	}

	return c.Write(ideas)
}

func (r resource) vote(c *routing.Context) error {
	ctx := c.Request.Context()

	var input VoteIdeaRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	ideas, err := r.service.Vote(ctx, input)

	if err != nil {
		return err
	}

	return c.Write(ideas)
}

func (r resource) create(c *routing.Context) error {
	var input CreateIdeaRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	idea, err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(idea, http.StatusCreated)
}

func (r resource) update(c *routing.Context) error {
	var input UpdateIdeaRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	idea, err := r.service.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		return err
	}

	return c.Write(idea)
}

func (r resource) delete(c *routing.Context) error {
	idea, err := r.service.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(idea)
}
