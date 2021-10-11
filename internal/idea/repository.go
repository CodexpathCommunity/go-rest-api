package idea

import (
	"context"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/pkg/dbcontext"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"strconv"
)

// Repository encapsulates the logic to access ideas from the data source.
type Repository interface {

	// Get returns the idea with the specified idea ID.
	Get(ctx context.Context, id string) (entity.Idea, error)

	// Count returns the number of ideas.
	Count(ctx context.Context) (int, error)

	// Create saves a new idea in the storage.
	Create(ctx context.Context, idea entity.Idea) error

	// Update updates the album with given ID in the storage.
	Update(ctx context.Context, idea entity.Idea) error

	// Delete removes the idea with given ID from the storage.
	Delete(ctx context.Context, id string) error

	Query(ctx context.Context, getIdeaRequest GetIdeaRequest) ([]entity.Idea, error)
}

// repository persists ideas in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new album repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the idea with the specified ID from the database.
func (r repository) Get(ctx context.Context, id string) (entity.Idea, error) {
	var idea entity.Idea
	err := r.db.With(ctx).Select().Model(id, &idea)
	return idea, err
}

// Create saves a new idea record in the database.
// It returns the ID of the newly inserted idea record.
func (r repository) Create(ctx context.Context, idea entity.Idea) error {
	return r.db.With(ctx).Model(&idea).Insert()
}

// Update saves the changes to an idea in the database.
func (r repository) Update(ctx context.Context, idea entity.Idea) error {
	return r.db.With(ctx).Model(&idea).Update()
}

// Delete deletes an idea with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id string) error {
	idea, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&idea).Delete()
}

// Count returns the number of the idea records in the database.
func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("idea").Row(&count)
	return count, err
}

func (r repository) Query(ctx context.Context, getIdeaRequest GetIdeaRequest) ([]entity.Idea, error) {
	var ideas []entity.Idea

	var pageSize, pageOffset int

	// TODO : Specify the max length here
	pageSize = 10

	if getIdeaRequest.PageSize!=0 {
		pageSize = getIdeaRequest.PageSize
	}
	pageOffset = getIdeaRequest.PageNumber * pageSize

	var queryString string

	if getIdeaRequest.TopPopularNumber != 0 {
		queryString = "select * from idea order by votes DESC LIMIT "+ strconv.Itoa(getIdeaRequest.TopPopularNumber)
	}else {
		queryString = "select id, author_email,tags,bad_flag,enabled,issues,votes,created_at,updated_at "
		if getIdeaRequest.IncludeSummary {
			queryString = queryString + ", summary"
		}
		if getIdeaRequest.IncludeContent {
			queryString = queryString + ", content"
		}
		if getIdeaRequest.IncludeMedia {
			queryString = queryString + ", media, media_types"
		}
		queryString = queryString + " from idea "

		if getIdeaRequest.IdeaId != "" || getIdeaRequest.MediaType != "" || getIdeaRequest.MinPopularity != 0 || getIdeaRequest.MaxPopularity != 0 {
			queryString = queryString + " where "
			var andReqd bool
			andReqd = false
			if getIdeaRequest.IdeaId != "" {
				queryString = queryString + " id = '" + getIdeaRequest.IdeaId + "'"
				andReqd = true
			}
			if getIdeaRequest.MediaType != "" {
				if andReqd {
					queryString = queryString + " and "
				}
				queryString = queryString + " '" + getIdeaRequest.MediaType + "' = any(media_types)"
				andReqd = true
			}
			if getIdeaRequest.MinPopularity != 0 {
				if andReqd {
					queryString = queryString + " and "
				}
				queryString = queryString + " votes >=" + strconv.Itoa(getIdeaRequest.MinPopularity)
				andReqd = true
			}
			if getIdeaRequest.MaxPopularity != 0 {
				if andReqd {
					queryString = queryString + " and "
				}
				queryString = queryString + " votes <=" + strconv.Itoa(getIdeaRequest.MaxPopularity)
			}
		}

		queryString = queryString + "limit " + strconv.Itoa(pageSize) + " offset " +
			strconv.Itoa(pageOffset)
	}

	err := r.db.With(ctx).
		NewQuery(queryString).
		All(&ideas)

	return ideas, err
}
