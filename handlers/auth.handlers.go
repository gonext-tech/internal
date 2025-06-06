package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/queries"
	"github.com/gonext-tech/internal/views/auth_views"
	"github.com/gonext-tech/internal/views/components"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const (
	auth_sessions_key string = "authenticate-sessions"
	auth_key          string = "authenticated"
	user_id_key       string = "user_id"
	email_key         string = "email"
	name_key          string = "name"
	tzone_key         string = "time_zone"
)

/********** Handlers for Auth Views **********/

type AuthService interface {
	GetALL(queries.InvoiceQueryParams) ([]models.Admin, models.Meta, error)
	CreateUser(u models.Admin, passwordConfirm string) error
	CheckEmail(email string) (models.Admin, error)
	// GetUserById(id int) (services.User, error)
}

func NewAuthHandler(us AuthService) *AuthHandler {

	return &AuthHandler{
		AdminServices: us,
	}
}

type AuthHandler struct {
	AdminServices AuthService
}

func (ah *AuthHandler) HomeHandler(c echo.Context) error {
	sess, _ := session.Get(auth_sessions_key, c)
	if auth, ok := sess.Values[auth_key].(bool); !ok || !auth {
		fromProtected = false

		return c.Redirect(http.StatusSeeOther, "/login")
	}
	return c.Redirect(http.StatusSeeOther, "/project")
	// homeView := auth_views.Home(fromProtected)
	// isError = false
	// return renderView(c, auth_views.HomeIndex(
	// 	"| Home",
	// 	"",
	// 	fromProtected,
	// 	isError,
	// 	getFlashmessages(c, "error"),
	// 	getFlashmessages(c, "success"),
	// 	homeView,
	// ))
}

func (ah *AuthHandler) RegisterHandler(c echo.Context) error {
	registerView := auth_views.Register(fromProtected)
	isError = false

	if c.Request().Method == "POST" {
		user := models.Admin{
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
			//Username: c.FormValue("username"),
		}
		passwordConfirm := c.FormValue("password_confirm")

		err := ah.AdminServices.CreateUser(user, passwordConfirm)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				err = errors.New("the email is already in use")
				setFlashmessages(c, "error", fmt.Sprintf(
					"something went wrong: %s",
					err,
				))

				return c.Redirect(http.StatusSeeOther, "/register")
			}

			setFlashmessages(c, "error", fmt.Sprintf(
				"something went wrong: %s",
				err,
			))

			return c.Redirect(http.StatusSeeOther, "/register")
		}

		setFlashmessages(c, "success", "You have successfully registered!!")

		return c.Redirect(http.StatusSeeOther, "/login")
	}

	return renderView(c, auth_views.RegisterIndex(
		"| Register",
		"",
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		registerView,
	))
}

func (ah *AuthHandler) SearchUser(c echo.Context) error {
	var queries queries.InvoiceQueryParams
	queries.SetDefaults()
	queries.Limit = 5
	users, _, err := ah.AdminServices.GetALL(queries)
	if err != nil {
		isError = false
		errorMsg = "can't fetch projects"
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
	}

	return renderView(c, components.UserResult(users, false))
}

func (ah *AuthHandler) LoginHandler(c echo.Context) error {
	loginView := auth_views.Login(fromProtected)
	isError = false

	if c.Request().Method == "POST" {
		// obtaining the time zone from the POST request of the login form
		tzone := ""
		if len(c.Request().Header["X-Timezone"]) != 0 {
			tzone = c.Request().Header["X-Timezone"][0]
		}
		// Authentication goes here
		user, err := ah.AdminServices.CheckEmail(c.FormValue("email"))
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				setFlashmessages(c, "error", "There is no user with that email")

				return c.Redirect(http.StatusSeeOther, "/login")
			}

			setFlashmessages(c, "error", fmt.Sprintf(
				"something went wrong: %s",
				err,
			))
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(c.FormValue("password")),
		)

		if err != nil {
			setFlashmessages(c, "error", "Incorrect password")

			return c.Redirect(http.StatusSeeOther, "/login")
		}

		// Get Session and setting Cookies
		sess, err := session.Get(auth_sessions_key, c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   360000, // in seconds
			HttpOnly: true,
		}

		// Set user as authenticated, their username,
		// their ID and the client's time zone
		sess.Values = map[interface{}]interface{}{
			auth_key:    true,
			user_id_key: user.ID,
			email_key:   user.Email,
			name_key:    user.Name,
			tzone_key:   tzone,
		}

		if err := sess.Save(c.Request(), c.Response()); err != nil {
			setFlashmessages(c, "error", "can't login now")
		}

		setFlashmessages(c, "success", "You have successfully logged in!!")

		return c.Redirect(http.StatusSeeOther, "/")
	}

	return renderView(c, auth_views.LoginIndex(
		"| Login",
		"",
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		loginView,
	))
}

func (ah *AuthHandler) LogoutHandler(c echo.Context) error {
	// Get the session
	sess, err := session.Get(auth_sessions_key, c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	// Clear the session

	sess.Options.MaxAge = -1
	for k := range sess.Values {
		delete(sess.Values, k)
	}
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}
	// Redirect to the login page after logout
	return c.Redirect(http.StatusSeeOther, "/login")
}

func (ah *AuthHandler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(auth_sessions_key, c)
		// Optional: Force session clear for debugging
		if auth, ok := sess.Values[auth_key].(bool); !ok || !auth {
			fromProtected = false
			for k := range sess.Values {
				delete(sess.Values, k)
			}
			sess.Save(c.Request(), c.Response())
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		if userId, ok := sess.Values[user_id_key].(int); ok && userId != 0 {
			c.Set(user_id_key, userId) // set the user_id in the context
		}

		if email, ok := sess.Values[email_key].(string); ok && len(email) != 0 {
			c.Set(email_key, email) // set the username in the context
		}
		if name, ok := sess.Values[name_key].(string); ok && len(name) != 0 {
			c.Set(email_key, name) //set the name in the context
		}

		if tzone, ok := sess.Values[tzone_key].(string); ok && len(tzone) != 0 {
			c.Set(tzone_key, tzone) // set the client's time zone in the context
		}

		fromProtected = true

		return next(c)
	}
}
