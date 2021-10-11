# Go RESTful API 

The project uses the following Go packages which can be easily replaced with your own favorite ones
since their usages are mostly localized and abstracted. 

* Routing: [ozzo-routing](https://github.com/go-ozzo/ozzo-routing)
* Database access: [ozzo-dbx](https://github.com/go-ozzo/ozzo-dbx)
* Database migration: [golang-migrate](https://github.com/golang-migrate/migrate)
* Data validation: [ozzo-validation](https://github.com/go-ozzo/ozzo-validation)
* Logging: [zap](https://github.com/uber-go/zap)
* JWT: [jwt-go](https://github.com/dgrijalva/jwt-go)

## Getting Started

If this is your first time encountering Go, please follow [the instructions](https://golang.org/doc/install) to
install Go on your computer. The kit requires **Go 1.13 or above**.

[Docker](https://www.docker.com/get-started) is also needed if you want to try the kit without setting up your
own database server. The kit requires **Docker 17.05 or higher** for the multi-stage build support.

After installing Go and Docker, run the following commands to start experiencing this starter kit:

```shell
# start a PostgreSQL database server in a Docker container
make db-start

# seed the database with some test data
make testdata

# run the RESTful API server
make run
```
At this time, you have a RESTful API server running at `http://localhost:8080`.

### Updating Database Schema
The starter kit uses [database migration](https://en.wikipedia.org/wiki/Schema_migration) to manage the changes of the 
database schema over the whole project development phase. The following commands are commonly used with regard to database
schema changes:

```shell
# Execute new migrations made by you or other team members.
# Usually you should run this command each time after you pull new code from the code repo. 
make migrate

# Create a new database migration.
# In the generated `migrations/*.up.sql` file, write the SQL statements that implement the schema changes.
# In the `*.down.sql` file, write the SQL statements that revert the schema changes.
make migrate-new

# Revert the last database migration.
# This is often used when a migration has some issues and needs to be reverted.
make migrate-down

# Clean up the database and rerun the migrations from the very beginning.
# Note that this command will first erase all data and tables in the database, and then
# run all migrations. 
make migrate-reset
```

### Managing Configurations
The application configuration is represented in `internal/config/config.go`. When the application starts,
it loads the configuration from a configuration file as well as environment variables. The path to the configuration 
file is specified via the `-config` command line argument which defaults to `./config/local.yml`. Configurations
specified in environment variables should be named with the `APP_` prefix and in upper case. When a configuration
is specified in both a configuration file and an environment variable, the latter takes precedence. 

The `config` directory contains the configuration files named after different environments. For example,
`config/local.yml` corresponds to the local development environment and is used when running the application 
via `make run`.

Do not keep secrets in the configuration files. Provide them via environment variables instead. For example,
you should provide `Config.DSN` using the `APP_DSN` environment variable. Secrets can be populated from a secret
storage (e.g. HashiCorp Vault) into environment variables in a bootstrap script (e.g. `cmd/server/entryscript.sh`).


### API Documentation  ###

# User Creation Flow

1. Create User
   POST /v1/user
   Input Body:
   requester_user_email: logged-in user, who is creating the user
   email_address: email of the user to be a created 
   role: super_admin/admin/visitor
   name: name of user
   country: country of user
super_admin can create admin and visitors
admin can create visitors

2. Update User
   PUT /v1/user/<email>
   email (in the path variable): email of the user to be updated

Input Body:
requester_user_email: logged-in user, who is updating the user
role: super_admin/admin/visitor of the user to be updated
name: name of user
country: country of user

super_admin can update admin and visitors
admin can update visitors

3. Delete User
   DELETE /v1/user/<email>
   email (in the path variable): email of the user to be updated

Input Body:
requester_user_email: logged-in user, who is updating the user

super_admin can delete admin and visitors
admin can delete visitors

4. GET User
   GET /v1/user/<email>
   email (in the path variable): email of the user


## User Signup Flow

1. User SignUp
   POST /userSignup/<email>/<code>
   path variables:
   email: email of the user to be signed up
   code: a unique code which will be used for authentication

2. User Confirm Email
   GET /userEmailConfirm/<email>/<code>
   path variables:
   email: email of the user to be confirm signin
   code: code which is sent in the User Signup email

# How to use?

1. User 'User SignUp' API to send email along with the unique code
2. User will receive the email, and the link will route to 'User Confirm Email' API, hence setting the 'is_auth' as true for the user
3. Now hit the 'GET User' API to check if authentication was successful or not by checking the value of 'is_auth' field.
4. If User SignUp request is hit again for same user, confirm user step needs to be repeated as it will set 'is_auth' as false again for the user.
5. You can change the html template in the email by changing the file : ??


## Idea Creation Flow

1. Create Idea




