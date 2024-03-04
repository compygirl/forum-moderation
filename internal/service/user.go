package service

import (
	"errors"
	"fmt"
	"forum/internal/database"
	"forum/internal/models"
	"forum/internal/web/handlers/helpers"
	"log"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

type UserServiceImpl struct {
	repo database.UserRepoInterface
}

func CreateNewUserService(repo database.UserRepoInterface) *UserServiceImpl {
	usrSrvc := UserServiceImpl{repo: repo}
	return &usrSrvc
}

func (userObj *UserServiceImpl) CreateUser(user *models.User) (int, int, error) {
	var id int64
	var err error
	if err = userObj.isUserParamsValid(user); err != nil {
		return http.StatusBadRequest, -1, err
	}
	emailUser, err := userObj.repo.GetUserByEmail(user.Email)
	if err != nil {
		if err.Error() != errors.New("element with EMAIL not found").Error() {
			return http.StatusInternalServerError, -1, err
		}
	}

	accUser, err := userObj.repo.GetUserByUsername(user.Username)
	if err != nil {
		if err.Error() != errors.New("element with USERNAME not found").Error() {
			return http.StatusInternalServerError, -1, err
		}
	}

	if emailUser != nil || accUser != nil {
		return http.StatusBadRequest, -1, errors.New("username or email was already used")
	} else {
		id, err = userObj.repo.CreateUserRepo(user)
		if err != nil {
			return http.StatusBadRequest, -1, err
		}
	}
	return http.StatusOK, int(id), nil
}

func (userObj *UserServiceImpl) Login(email, password string, admin bool) (*models.Session, error) {
	// fmt.Println("Logining...: ", admin)
	user := &models.User{}
	var err error

	if user, err = userObj.repo.GetUserByEmail(email); err != nil {
		log.Printf("Login: GetUserByEmail: %v", err)
		return nil, errors.New("Provided Email is Incorrect or doesn't exist")
	}
	// fmt.Println("USERNAME: ", user.Username)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("Provided Password is Incorrect")
	}

	role, err := userObj.repo.GetUserRole(user.UserUserID)
	// fmt.Println("USER ID: ", user.UserID, "  role:  ", role)
	if err != nil {
		return nil, errors.New("Some error with query to get user role")
	}

	if admin {
		if role != "admin" {
			return nil, errors.New("You do not have Admin access!")
		}
	} else {
		if role == "admin" {
			return nil, errors.New("You should select the admin role when logining")
		}
	}

	session := &models.Session{
		UserID:  user.UserUserID,
		Token:   uuid.New().String(),
		ExpTime: time.Now().Add(10 * time.Minute),
	}

	if err := userObj.repo.DeleteSessionByUserID(user.UserUserID); err != nil {
		return nil, errors.New("Error deleting the Session by USER ID")
	}

	err = userObj.repo.CreateSession(session)
	if err != nil {
		return nil, err
	}
	// fmt.Println("Reaching the end of the Login")
	return session, nil
}

func (userObj *UserServiceImpl) isUserParamsValid(user *models.User) error {
	if err := userObj.isUserEmailValid(user); err != nil {
		return err
	}
	if err := userObj.isUserNameValid(user); err != nil {
		return err
	}
	if err := userObj.isPasswordValid(user); err != nil {
		return err
	}
	return nil
}

func (userObj *UserServiceImpl) isUserEmailValid(user *models.User) error {
	regexPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(regexPattern, user.Email)
	if !matched {
		return errors.New("Invalid Email")
	}
	return nil
}

func (userObj *UserServiceImpl) isUserNameValid(user *models.User) error {
	if user.Username == "" || len(user.Username) < 2 {
		return errors.New("Invalid Username")
	}
	return nil
}

func (userObj *UserServiceImpl) isPasswordValid(user *models.User) error {
	if len(user.Password) < 2 {
		return errors.New("Weak Password (length cannot less than 5)")
	}
	return nil
}

// ! Эта функция на самом деле проверяет на пустоту куки, по этому должна называтся по другому
// ! IsUserLoggedIn - Это название обозначает что она проверяет куки на существование пользоватля
func (userObj *UserServiceImpl) IsUserLoggedIn(r *http.Request) bool {
	cookie := helpers.SessionCookieGet(r)
	return cookie != nil && cookie.Value != ""
}

func (userObj *UserServiceImpl) Logout(token string) error {
	return userObj.repo.DeleteSessionByToken(token)
}

func (userObj *UserServiceImpl) IsTokenExist(token string) bool {
	if session, err := userObj.repo.GetSessionByToken(token); session == nil || err != nil {
		return false
	}
	return true
}

func (userObj *UserServiceImpl) GetSession(token string) (*models.Session, error) {
	session, err := userObj.repo.GetSessionByToken(token)
	if session == nil || err != nil {
		return nil, err
	}
	return session, nil
}

// prop: (token:string), returns: ( expires, err )
func (userObj *UserServiceImpl) ExtendSessionTimeout(token string) (time.Time, error) {
	session, err := userObj.repo.GetSessionByToken(token)
	if session == nil || err != nil {
		fmt.Println("ExtendSessionTimeout: Problem with getting session")
		return time.Time{}, err
	}
	session.ExpTime = session.ExpTime.Add(1 * time.Minute)
	if err = userObj.repo.UpdateSession(session); err != nil {
		fmt.Println("ExtendSessionTimeout: Problem with update session")

		return time.Time{}, err
	}

	return session.ExpTime, nil
}

func (userObj *UserServiceImpl) GetUserByUserID(userID int) (*models.User, error) {
	user, err := userObj.repo.GetUserByUserID(userID)
	if user == nil || err != nil {
		return nil, fmt.Errorf("GetUserByUserID: %w", err)
	}
	return user, nil
}

func (userObj *UserServiceImpl) GoogleAuthorization(googleUser *models.GoogleLoginUserData) (*models.Session, error) {
	user, err := userObj.repo.GetUserByEmail(googleUser.Email)
	// var userID int
	if err != nil {
		// If the user does not exist, create a new user record
		user = &models.User{
			Username: googleUser.Email,
			Email:    googleUser.Email,
			Password: "dummypassword",
		}
		_, user.UserUserID, err = userObj.CreateUser(user)

		if err != nil && err.Error() != errors.New("element with EMAIL not found").Error() {
			return nil, err
		}
	}

	session := &models.Session{
		UserID:  user.UserUserID,
		Token:   uuid.New().String(),
		ExpTime: time.Now().Add(10 * time.Minute),
	}
	if err := userObj.repo.DeleteSessionByUserID(user.UserUserID); err != nil {
		return nil, errors.New("Error deleting the Session by USER ID")
	}

	err = userObj.repo.CreateSession(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (userObj *UserServiceImpl) GitHubAuthorization(githubUser *models.GitHubLoginUserData) (*models.Session, error) {
	if githubUser.Login == "" {
		githubUser.Login = githubUser.Email
	} else if githubUser.Email == "" {
		githubUser.Email = githubUser.Login + "@gmail.com"
	}
	user, err := userObj.repo.GetUserByEmail(githubUser.Email)
	// var userID int
	if err != nil {
		// If the user does not exist, create a new user record
		user = &models.User{
			Username: githubUser.Login,
			Email:    githubUser.Email,
			Password: "dummypassword",
		}
		_, user.UserUserID, err = userObj.CreateUser(user)
		if err != nil && err.Error() != errors.New("element with EMAIL not found").Error() {
			return nil, err
		}
	}
	session := &models.Session{
		UserID:  user.UserUserID,
		Token:   uuid.New().String(),
		ExpTime: time.Now().Add(10 * time.Minute),
	}
	if err := userObj.repo.DeleteSessionByUserID(user.UserUserID); err != nil {
		return nil, errors.New("Error deleting the Session by USER ID")
	}
	err = userObj.repo.CreateSession(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (userObj *UserServiceImpl) ChangeUserRole(newRole string, userID int) error {
	err := userObj.repo.ChangeUserRole(newRole, userID)
	if err != nil {
		return err
	}
	return nil
}

func (userObj *UserServiceImpl) GetUsersByRole(role string) ([]*models.User, error) {
	users, err := userObj.repo.GetUserByRole(role)
	if err != nil {
		return nil, err
	}
	return users, nil
}
