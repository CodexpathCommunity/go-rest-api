package user

import (
	"context"
	"github.com/mailgun/mailgun-go/v3"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/internal/errors"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"time"
)

// Service encapsulates usecase logic for users.
type UserService interface {
	CreateUser(ctx context.Context, input CreateUserRequest) (User, error)
	GetUser(ctx context.Context, email string) (User, error)
	UpdateUser(ctx context.Context, email string, input UpdateUserRequest, bypassAuth bool) (User, error)
	DeleteUser(ctx context.Context, email string, input DeleteUserRequest) (User, error)
	CheckPermission(ctx context.Context, requesterEmail string, role string) (bool, error)
	UserSignUp(ctx context.Context, email string, code string) (mes string, id string, err error)
	AuthenticateUser(ctx context.Context, email string, code string) (User, error)
}

// User represents the data about a User.
type User struct {
	entity.Users
}

// CreateUserRequest represents an user creation request.
type CreateUserRequest struct {
	// the logged in user, who is making the request for createUser
	RequesterUserEmail string `json:"requester_user_email"`
	EmailAddress   string	  `json:"email_address"`
	Role           string     `json:"role"`
	Name           string     `json:"name"`
	Country        string     `json:"country"`
}

// UpdateUserRequest represents an user update request.
type UpdateUserRequest struct {
	RequesterUserEmail string `json:"requester_user_email"`
	Role               string `json:"role"`
	Name               string `json:"name"`
	Country            string `json:"country"`

	// for internal use only
	IncreaseScore      bool   `json:"increase_score"`
	AuthCode           string `json:"auth_code"`
	Authenticate       bool   `json:"authenticate"`
	ResetAuth          bool   `json:"reset_auth"`
}

type DeleteUserRequest struct {
	RequesterUserEmail string `json:"requester_user_email"`
	EmailAddress   string	  `json:"email_address"`
}

type userService struct {
	repo   UsersRepository
	logger log.Logger
}

func (s userService) UserSignUp(ctx context.Context, email string, code string) (mes string, id string, err error) {

	mg := mailgun.NewMailgun(
		//"YOUR_DOMAIN_NAME", // Domain name
		"sandbox25224b1d21a0489d823328614e0bf07c.mailgun.org",
		//"YOUR_API_KEY",     // API Key
		"3b802a0169e414f8991b84f250717afb-dbdfb8ff-9e66809e",
	)

	m :=  mg.NewMessage(
		/* From */ "Excited User <mailgun@sandbox25224b1d21a0489d823328614e0bf07c.mailgun.org>",
		/* Subject */ "Please confirm your email address!!",
		/* Body */ "Hey, Please confirm your email address by clicking on the below link !!",
		/* To */ email,
	)

	// Create a new template
	err = mg.CreateTemplate(ctx, &mailgun.Template{
		Name: "my-template7",
		Version: mailgun.TemplateVersion{
			Template: `
<!DOCTYPE html>
<head>
<meta name="viewport" content="width=device-width" />
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />

<style type="text/css">
* {
  margin: 0;
  font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
  box-sizing: border-box;
  font-size: 14px;
}

img {
  max-width: 100%;
}

body {
}

/* Let's make sure all tables have defaults */
table td {
  vertical-align: top;
}

/* -------------------------------------
    BODY & CONTAINER
------------------------------------- */
body {
  background-color: #f6f6f6;
}

.body-wrap {
  background-color: #f6f6f6;
  width: 100%;
}

.container {
  display: block !important;
  max-width: 600px !important;
  margin: 0 auto !important;
  /* makes it centered */
  clear: both !important;
}

.content {
  max-width: 600px;
  margin: 0 auto;
  display: block;
  padding: 20px;
}

.div-centre
{
  width: 500px;
  display: block;
  margin-left: auto;
  margin-right: auto;
}

/* -------------------------------------
    HEADER, FOOTER, MAIN
------------------------------------- */
.main {
  padding-top: 20px;
  padding-left: 20px;
  background-color: #fff;
  border: 1px solid #e9e9e9;
  border-radius: 3px;
}

.content-wrap {
  padding: 20px;
}

.logo-small {
  padding: 10px;
}

#logo-bg {
  background-color: grey;
}

.padded a>img {
margin-top: 25px;
}

.content-block {
  padding: 0 0 20px;
}

.header {
  width: 100%;
  margin-bottom: 20px;
}

.footer {
  width: 100%;
  clear: both;
  color: #999;
  padding: 20px;
}
.footer p, .footer a, .footer td {
  color: #999;
  font-size: 12px;
}

/* -------------------------------------
    TYPOGRAPHY
------------------------------------- */
h1, h2, h3 {
  font-family: "Helvetica Neue", Helvetica, Arial, "Lucida Grande", sans-serif;
  color: #000;
  margin: 40px 0 0;
  line-height: 1.2em;
  font-weight: 400;
}

h1 {
  font-size: 32px;
  font-weight: 500;
  /* 1.2em * 32px = 38.4px, use px to get airier line-height also in Thunderbird, and Yahoo!, Outlook.com, AOL webmail clients */
  /*line-height: 38px;*/
}

h2 {
  font-size: 24px;
  /* 1.2em * 24px = 28.8px, use px to get airier line-height also in Thunderbird, and Yahoo!, Outlook.com, AOL webmail clients */
  /*line-height: 29px;*/
}

h3 {
  font-size: 18px;
  /* 1.2em * 18px = 21.6px, use px to get airier line-height also in Thunderbird, and Yahoo!, Outlook.com, AOL webmail clients */
  /*line-height: 22px;*/
}

h4 {
  font-size: 14px;
  font-weight: 600;
}

p, ul, ol {
  margin-bottom: 10px;
  font-weight: normal;
}
p li, ul li, ol li {
  margin-left: 5px;
  list-style-position: inside;
}

/* -------------------------------------
    LINKS & BUTTONS
------------------------------------- */
a {
  color: #348eda;
  text-decoration: underline;
}

.btn-primary {
  text-decoration: none;
  color: #FFF;
  background-color: #348eda;
  border: solid #348eda;
  border-width: 10px 20px;
  line-height: 2em;
  /* 2em * 14px = 28px, use px to get airier line-height also in Thunderbird, and Yahoo!, Outlook.com, AOL webmail clients */
  /*line-height: 28px;*/
  font-weight: bold;
  text-align: center;
  cursor: pointer;
  display: inline-block;
  border-radius: 5px;
  text-transform: capitalize;
}

/* -------------------------------------
    OTHER STYLES THAT MIGHT BE USEFUL
------------------------------------- */
.last {
  margin-bottom: 0;
}

.first {
  margin-top: 0;
}

.aligncenter {
  text-align: center;
}

.alignright {
  text-align: right;
}

.alignleft {
  text-align: left;
}

.clear {
  clear: both;
}

/* -------------------------------------
    ALERTS
    Change the class depending on warning email, good email or bad email
------------------------------------- */
.alert {
  font-size: 16px;
  color: #fff;
  font-weight: 500;
  padding: 5px;
  text-align: center;
  border-radius: 3px 3px 0 0;
}
.alert a {
  color: #fff;
  text-decoration: none;
  font-weight: 500;
  font-size: 16px;
}
.alert.alert-warning {
  background-color: #FF9F00;
}
.alert.alert-bad {
  background-color: #D0021B;
}
.alert.alert-good {
  background-color: #68B90F;
}

/* -------------------------------------
    INVOICE
    Styles for the billing table
------------------------------------- */
.invoice {
  margin: 40px auto;
  text-align: left;
  width: 80%;
}
.invoice td {
  padding: 5px 0;
}
.invoice .invoice-items {
  width: 100%;
}
.invoice .invoice-items td {
  border-top: #eee 1px solid;
}
.invoice .invoice-items .total td {
  border-top: 2px solid #333;
  border-bottom: 2px solid #333;
  font-weight: 700;
}

/* -------------------------------------
    RESPONSIVE AND MOBILE FRIENDLY STYLES
------------------------------------- */
@media only screen and (max-width: 640px) {
  body {
    padding: 0 !important;
  }

  h1, h2, h3, h4 {
    font-weight: 800 !important;
    margin: 20px 0 5px !important;
  }

  h1 {
    font-size: 22px !important;
  }

  h2 {
    font-size: 18px !important;
  }

  h3 {
    font-size: 16px !important;
  }

  .container {
    padding: 0 !important;
    width: 100% !important;
  }

  .content {
    padding: 0 !important;
  }

  .content-wrap {
    padding: 10px !important;
  }

  .invoice {
    width: 100% !important;
  }
}
</style>

<title>Actionable emails e.g. reset password</title>
<link href="styles.css" media="all" rel="stylesheet" type="text/css" />
</head>

<body itemscope itemtype="http://schema.org/EmailMessage">

<table class="body-wrap">
	<tr>
		<td></td>
		<td class="container" width="600">
			<div class="content">
				<table class="main logo-small logo-bg" width="100%" cellpadding="0" cellspacing="0" itemprop="action" itemscope itemtype="http://schema.org/ConfirmAction">
					<tr>
						<td class="content-wrap">
							<meta itemprop="name" content="Confirm Email"/>
							<table width="100%" cellpadding="0" cellspacing="0">
								<tr>
									<td class="content-block">
										Please confirm your email address by clicking the link below.
									</td>
								</tr>
								<tr>
									<td class="content-block" itemprop="handler" itemscope itemtype="http://schema.org/HttpActionHandler">
										<a href="http://localhost:8080/v1/userEmailConfirm/{{.email}}/{{.code}}" class="btn-primary" itemprop="url">Confirm email address</a>
									</td>
								</tr>
								<tr>
									<td class="content-block">
										&mdash; Danderdee
									</td>
								</tr>
							</table>
						</td>
					</tr>
				</table>
				<div class="footer">
					<table width="100%">
					<tr>
					<td class="aligncenter content-block">
</td>
</tr>

						<tr>
						<td class="aligncenter">
						Copyright © 2020 Danderdee Limited. All Rights Reserved.
						</td>
						</tr>
						<tr>
						<td class="aligncenter">
						Danderdee Technologies Limited is a company registered in Scotland. Company number: SC653922 
						</td>
						</tr>
						<tr>
						<td class="aligncenter">
						Registered Office: 9 Watling Street, Dumfries, Scotland, DG1 3HQ.
						</td>
						</tr>
						<tr>
						<td class="aligncenter content-block padded">
						<div class="padded">
						<a href="https://www.facebook.com/DanderdeeWifi">
						</a>
						<a href="https://twitter.com/danderdee">
						</a>
						</div>
						</td>
						</tr>									
					</table>
				</div></div>
		</td>
		<td></td>
	</tr>
</table>

</body>
</html>
`,
            Engine:   mailgun.TemplateEngineGo,
			Tag:      "v1",
		},
	})

	// Give time for template to show up in the system.
	time.Sleep(time.Second * 2)

	m.SetTemplate("my-template7")
	m.AddVariable("code",code)
	m.AddVariable("email",email)

	msg, id, err := mg.Send(ctx, m)

	if err!=nil{
		return msg, id, err
	}

	var input UpdateUserRequest
	input.AuthCode = code
	s.UpdateUser(ctx, email, input, true)

	return msg, id, err
}

func (s userService) AuthenticateUser(ctx context.Context, email string, code string) (User, error) {

	user, errRqu := s.GetUser(ctx, email)
	if errRqu != nil {
		return User{}, errors.InternalServerError("User to be authenticated doesn't exists")
	}

	if user.AuthCode == code{
		var input UpdateUserRequest
		input.Authenticate = true
		user,errReq := s.UpdateUser(ctx, email, input, true)
		return user, errReq
	}else{
		return User{}, errors.InternalServerError("Auth code doesn't match")
	}

}

// NewService creates a new user service.
func NewUserService(repo UsersRepository, logger log.Logger) UserService {
	return userService{repo, logger}
}

var SUPER_ADMIN = "super_admin"
var ADMIN = "admin"
var VISITOR = "visitor"
var roles = []string{SUPER_ADMIN, ADMIN, VISITOR}

// constants
var permissions = map[string][]string{
	SUPER_ADMIN: {SUPER_ADMIN, ADMIN, VISITOR},
	ADMIN: []string{VISITOR},
	VISITOR: []string{},
}

// Create creates a new user.
func (s userService) CreateUser(ctx context.Context, req CreateUserRequest) (User, error) {
	now := time.Now()

	isPermitted, errMsg := s.CheckPermission(ctx, req.RequesterUserEmail, req.Role)

	if !isPermitted {
		return User{}, errMsg
	}

	err := s.repo.CreateUser(ctx, entity.Users{
		ID:           req.EmailAddress,
		Role:         req.Role,
		Name:         req.Name,
		Country:      req.Country,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return User{}, err
	}
	return s.GetUser(ctx, req.EmailAddress)
}

// Get returns the user with the specified the user email.
func (s userService) GetUser(ctx context.Context, email string) (User, error) {
	user, err := s.repo.GetUser(ctx, email)
	if err != nil {
		return User{}, err
	}
	return User{user}, nil
}


// Update updates the user
func (s userService) UpdateUser(ctx context.Context, email string, req UpdateUserRequest, bypassAuth bool) (User, error) {

		user, errRqu := s.GetUser(ctx, email)
		if errRqu != nil {
			return User{}, errors.InternalServerError("User to be updated doesn't exists")
		}

	if !bypassAuth {
		isPermitted, errMsg := s.CheckPermission(ctx, req.RequesterUserEmail, req.Role)

		if !isPermitted {
			return User{}, errMsg
		}
		user.Name = req.Name
		user.Country = req.Country
		user.Role = req.Role
	}else {
		if req.IncreaseScore {
			user.Score++
		}
		if req.ResetAuth {
			user.AuthCode = req.AuthCode
			user.IsAuth = false
		}
		if req.Authenticate {
			user.IsAuth = true
		}
	}
		user.UpdatedAt = time.Now()

	if err := s.repo.UpdateUser(ctx, user.Users); err != nil {
		return user, err
	}
	return user, nil
}

// Delete deletes the user with the specified ID.
func (s userService) DeleteUser(ctx context.Context, email string, req DeleteUserRequest) (User, error) {

	user, err :=  s.GetUser(ctx, email)
	if err != nil  {
		return User{}, errors.InternalServerError("Requester User doesn't exists")
	}

	isPermitted, errMsg := s.CheckPermission(ctx, req.RequesterUserEmail, user.Role)

	if !isPermitted {
		return User{}, errMsg
	}

	if err = s.repo.DeleteUser(ctx, email); err != nil {
		return User{}, err
	}
	return user, nil
}


func (s userService) CheckPermission(ctx context.Context, requesterEmail string, role string) (bool, error) {

	requesterUser, errRqu :=  s.GetUser(ctx, requesterEmail)

	if errRqu != nil  {
		return false, errors.InternalServerError("Requester User doesn't exists")
	}

	var roleExists bool
	roleExists = false
	for i := range roles {
		if roles[i] == role {
			roleExists = true
			break
		}
	}

	if !roleExists {
		return false, errors.InternalServerError("This role doesn't exists in the system : " + role)
	}

	var rolePermitted bool
	rolePermitted = false
	var permList []string
	permList = permissions[requesterUser.Role]
	for i := range permList {
		if permList[i] == role {
			rolePermitted = true
			break
		}
	}

	if !rolePermitted {
		return false, errors.InternalServerError("Requester User doesn't have required permission")
	}

	return true, errors.NotFound("")
}