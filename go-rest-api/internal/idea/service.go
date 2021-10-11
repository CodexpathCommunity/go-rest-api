package idea

import (
	"context"
	"github.com/qiangxue/go-rest-api/internal/errors"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"github.com/qiangxue/go-rest-api/internal/user"
	"time"
)

// Service encapsulates usecase logic for ideas.
type Service interface {
	Get(ctx context.Context, id string) (Idea, error)
	Create(ctx context.Context, input CreateIdeaRequest) (Idea, error)
	Update(ctx context.Context, id string, req UpdateIdeaRequest) (Idea, error)
	Delete(ctx context.Context, id string) (Idea, error)
	Count(ctx context.Context) (int, error)
	Query(ctx context.Context, getIdeaRequest GetIdeaRequest) ([]Idea, error)
	Vote(ctx context.Context, voteIdeaRequest VoteIdeaRequest) (Idea, error)
}

// idea represents the data about an idea.
type Idea struct {
	entity.Idea
}

// CreateIdeaRequest represents an idea creation request.
type CreateIdeaRequest struct {
	AuthorEmail string `json:"author_email"`
	Summary     string `json:"summary"`
	Content     string `json:"content"`
	Media       []string `json:"media"`
	MediaTypes  []string `json:"media_types"`
	Tags        []string `json:"tags"`
	Issues      []string `json:"issues"`
}

// UpdateIdeaRequest represents an idea update request.
type UpdateIdeaRequest struct {
	RequesterUserEmail string   `json:"requester_user_email"`
	IP          string          `json:"ip"`
	Summary     string          `json:"summary"`
	Content     string          `json:"content"`
	Media       []string        `json:"media"`
	MediaTypes  []string        `json:"media_types"`
	Tags        []string        `json:"tags"`
	Issues      []string        `json:"issues"`
	BadFlag     bool            `json:"bad_flag"`
	Enabled     bool            `json:"enabled"`
}

type GetIdeaRequest struct {
	IdeaId          string     `json:"idea_id"`
	MediaType       string     `json:"media_type"`
	MinPopularity    int       `json:"min_popularity"`
	MaxPopularity    int       `json:"max_popularity"`
	TopPopularNumber int       `json:"top_popular_number"`
	IncludeMedia     bool      `json:"include_media"`
	IncludeSummary   bool      `json:"include_summary"`
	IncludeContent  bool       `json:"include_content"`
	PageSize         int       `json:"page_size"`
	PageNumber       int       `json:"page_number"`
}

type VoteIdeaRequest struct {
	IdeaId             string     `json:"idea_id"`
	RequesterUserEmail string     `json:"requester_user_email"`
}


type service struct {
	repo   Repository
	logger log.Logger
	userService user.UserService
}

// NewService creates a new idea service.
func NewService(repo Repository, logger log.Logger, user user.UserService) Service {
	return service{repo, logger, user}
}

// Get returns the idea with the specified the idea ID.
func (s service) Get(ctx context.Context, id string) (Idea, error) {
	idea, err := s.repo.Get(ctx, id)
	if err != nil {
		return Idea{}, err
	}
	return Idea{idea}, nil
}

// Create creates a new idea.
func (s service) Create(ctx context.Context, req CreateIdeaRequest) (Idea, error) {
	id := entity.GenerateID()
	now := time.Now()

	author, err2 := s.userService.GetUser(ctx, req.AuthorEmail)
	if err2 !=nil{
		return Idea{}, err2
	}
	if !(author.Role == "admin" ||  author.Role == "super_admin"){
		return Idea{}, errors.InternalServerError("Requester User doesn't have permission to create ideas")
	}

	err := s.repo.Create(ctx, entity.Idea{
		ID:       		  id,
		AuthorEmail:      req.AuthorEmail,
		Summary:          req.Summary,
		Media:            req.Media,
		Tags:             req.Tags,
		Issues:           req.Issues,
		Content:		  req.Content,
		MediaTypes: 	  req.MediaTypes,
		CreatedAt:        now,
		UpdatedAt:        now,
	})
	if err != nil {
		return Idea{}, err
	}
	return s.Get(ctx, id)
}


//Update updates the idea with the specified ID.
func (s service) Update(ctx context.Context, id string, req UpdateIdeaRequest) (Idea, error) {

	idea, err := s.Get(ctx, id)
	if err != nil {
		return Idea{}, err
	}

	author, err2 := s.userService.GetUser(ctx, req.RequesterUserEmail)
	if err2 !=nil{
		return Idea{}, errors.InternalServerError("Requester User doesn't exist")
	}
	if !(author.Role == "admin" ||  author.Role == "super_admin" || idea.AuthorEmail == author.ID){
		return Idea{}, errors.InternalServerError("Requester User doesn't have permission to edit ideas")
	}

	idea.Issues = req.Issues
	idea.Tags = req.Tags
	idea.Media = req.Media
	idea.Summary = req.Summary
	idea.BadFlag = req.BadFlag
	idea.Enabled = req.Enabled
	idea.UpdatedAt = time.Now()

	var IPExists bool
	IPExists = false
	for i := range idea.IssuesIPs  {
		if idea.VotersIds[i] == req.IP {
			IPExists = true
			break
		}
	}

	if !IPExists{
		idea.IssuesIPs = append(idea.IssuesIPs, req.IP)
	}

	if len(idea.IssuesIPs) >= 3 {
		idea.BadFlag = true
	}

	if err := s.repo.Update(ctx, idea.Idea); err != nil {
		return idea, err
	}
	return idea, nil
}


// Delete deletes the idea with the specified ID.
func (s service) Delete(ctx context.Context, id string) (Idea, error) {
	idea, err := s.Get(ctx, id)
	if err != nil {
		return Idea{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return Idea{}, err
	}
	return idea, nil
}

// Count returns the number of ideas.
func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}


func (s service) Query(ctx context.Context, getIdeaRequest GetIdeaRequest) ([]Idea, error) {
	items, err := s.repo.Query(ctx, getIdeaRequest)
	if err != nil {
		return nil, err
	}
	result := []Idea{}
	for _, item := range items {
		result = append(result, Idea{item})
	}
	return result, nil
}

func (s service) Vote(ctx context.Context, voteIdeaRequest VoteIdeaRequest) (Idea, error) {

	idea, err := s.Get(ctx, voteIdeaRequest.IdeaId)
	if err != nil {
		return Idea{}, errors.InternalServerError("Idea doesn't exist")
	}

	_, err2 := s.userService.GetUser(ctx, voteIdeaRequest.RequesterUserEmail)
	if err2 !=nil{
		return Idea{}, errors.InternalServerError("Requester User doesn't exist")
	}

	if voteIdeaRequest.RequesterUserEmail == idea.AuthorEmail  {
		return Idea{}, errors.InternalServerError("User cannot vote on it's own idea")
	}

	var voterExists bool
	voterExists = false
	for i := range idea.VotersIds  {
		if idea.VotersIds[i] == voteIdeaRequest.RequesterUserEmail {
			voterExists = true
			break
		}
	}

	if voterExists {
		return Idea{}, errors.InternalServerError("User has already voted on this idea")
	}

	idea.Votes++
	idea.VotersIds = append(idea.VotersIds,  voteIdeaRequest.RequesterUserEmail )

	var input user.UpdateUserRequest
	input.RequesterUserEmail = voteIdeaRequest.RequesterUserEmail
	input.IncreaseScore = true
	s.userService.UpdateUser(ctx,voteIdeaRequest.RequesterUserEmail, input, true)

	if err := s.repo.Update(ctx, idea.Idea); err != nil {
		return idea, err
	}
	return idea, nil
}

