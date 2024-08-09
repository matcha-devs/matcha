package main

import (
	"embed"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/matcha-devs/matcha/internal"
)

var (
	//go:embed all:internal/templates all:public
	content        embed.FS
	publicServer   = http.FileServer(http.FS(content))
	templateServer = template.Must(template.ParseFS(content, "internal/templates/*.go.html"))
	surfacePages   = map[string]struct{}{"signup": {}, "login": {}, "reset-password": {}}
)

func getPublic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "max-age=3600")
	publicServer.ServeHTTP(w, r)
}

func getIndex(w http.ResponseWriter, _ *http.Request) {
	err := templateServer.ExecuteTemplate(w, "index.go.html", nil)
	if err != nil {
		log.Println("Error executing template index.go.html -", err)
	}
}

func setSessionCookie(w http.ResponseWriter, id uint64) {
	log.Println("Issued cookie for id:", id)
	http.SetCookie(
		w, &http.Cookie{
			Name:     "c_user_id",
			Value:    strconv.FormatUint(id, 10),
			Path:     "/",
			Expires:  time.Now().Add(20 * time.Minute),
			MaxAge:   20 * 60,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
	)
}

func postSignup(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("first_name")
	// TODO(@seoyoungcho213): Validate user data way better here.
	middleName := r.FormValue("middle_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	// TODO(@seoyoungcho213): read about mail.ParseAddress' return type and use it to add first+last name formatting too
	if _, err := mail.ParseAddress("<" + email + ">"); err != nil {
		log.Println("Error adding user {"+email+"} to database -", err)
		if _, err := io.WriteString(w, err.Error()); err != nil {
			log.Println("Error writing signup email format error -", err)
		}
		return
	}
	password := r.FormValue("psw")
	if password != r.FormValue("psw-repeat") {
		log.Println("Passwords didnt match.")
		if _, err := io.WriteString(w, "passwords do not match"); err != nil {
			log.Println("Error writing signup failure -", err)
		}
		return
	}
	dateOfBirth := r.FormValue("date_of_birth")
	// TODO(@FaaizMemonPurdue): Add API call timeouts.
	id, err := matcha.database.AddUser(firstName, middleName, lastName, email, password, dateOfBirth)
	if err != nil {
		log.Println("Error adding user {"+email+"} to database -", err)
		if _, err := io.WriteString(w, err.Error()); err != nil {
			log.Println("Error writing server error -", err)
		}
		return
	}
	setSessionCookie(w, id)
	w.Header().Set("HX-Redirect", "/dashboard")
}

func postLogout(w http.ResponseWriter, _ *http.Request) {
	setSessionCookie(w, 0)
	w.Header().Set("HX-Redirect", "/")
}

func postLogin(w http.ResponseWriter, r *http.Request) {
	// TODO(@FaaizMemonPurdue): Add API call timeouts.
	id, err := matcha.database.AuthenticateLogin(r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		log.Println("Login failed -", err)
		if _, err := io.WriteString(w, err.Error()); err != nil {
			log.Println("Error writing login failure -", err)
		}
		return
	}
	setSessionCookie(w, id)
	w.Header().Set("HX-Redirect", "/dashboard")
}

func postResetPassword(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	// Validate the email format
	if _, err := mail.ParseAddress("<" + email + ">"); err != nil {
		log.Println("Invalid email format -", err)
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	return
}

func postDeleteUser(w http.ResponseWriter, r *http.Request) {
	// TODO(@FaaizMemonPurdue): Add API call timeouts.
	email := r.FormValue("email")
	id, err := matcha.database.AuthenticateLogin(email, r.FormValue("password"))
	if err != nil {
		if _, err := io.WriteString(w, err.Error()); err != nil {
			log.Println("User failed to validate delete request -", err)
		}
		return
	}

	// TODO(@FaaizMemonPurdue): Add API call timeouts.
	err = matcha.database.DeleteUser(id)
	if err != nil {
		log.Println("Delete User failed -", err)
		if _, err = io.WriteString(w, "internal server error"); err != nil {
			log.Println("Error writing deleted user internal server error -", err)
			return
		}
	}
	postLogout(w, r)
}

func checkLoginStatus(w http.ResponseWriter, r *http.Request) (user *internal.User) {
	cookie, err := r.Cookie("c_user_id")
	if errors.Is(err, http.ErrNoCookie) {
		log.Println("Client has no session cookie -", err)
		http.Error(w, "Unauthorized login session.", http.StatusUnauthorized)
		return nil
	} else if err != nil {
		http.Error(w, "Unauthorized login session.", http.StatusUnauthorized)
		log.Println("Error getting session cookie -", err)
		return nil
	}
	id, err := strconv.ParseUint(cookie.Value, 10, 64)
	if err != nil {
		log.Println("Failed to convert user id -", err)
		http.Error(w, "Invalid login session.", http.StatusBadRequest)
		return nil
	} else if id < 1 {
		log.Println("Invalid user id:", id)
		http.Error(w, "Invalid login session.", http.StatusBadRequest)
		return nil
	}
	// TODO(@FaaizMemonPurdue): Add API call timeouts.
	return matcha.database.GetUser(id)
}

func getPage(w http.ResponseWriter, r *http.Request) {
	pageData := struct {
		PageName string
		User     *internal.User
	}{
		PageName: strings.TrimLeft(r.URL.Path, "/"),
	}
	if _, exists := surfacePages[pageData.PageName]; !exists {
		pageData.User = checkLoginStatus(w, r)
		if pageData.User == nil {
			return
		}
	}
	if err := templateServer.ExecuteTemplate(w, pageData.PageName+".go.html", pageData); err != nil {
		log.Println("Error executing template", pageData.PageName, "-", err)
	}
}
