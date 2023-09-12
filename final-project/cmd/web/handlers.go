package main

import "net/http"

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.gohtml", nil)
}

func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}

func (app *Config) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context())

	//parse form post

	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)

	}

	//get email and password from form post

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.Models.User.GetByEmail(email)

	if err != nil {
		app.Session.Put(r.Context(), "error", "invalid credentials.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	//check password

	validPassword, err := user.PasswordMatches(password)

	if !validPassword {
		msg := Message{
			To:      "younesbouchbouk.py@gmail.com",
			Data:    "Invalid Password",
			subject: "Invalid Email address",
		}

		app.sendEmail(msg)
		app.Session.Put(r.Context(), "error", "invalid credentials.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	//okay , we can og the user in

	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)

	app.Session.Put(r.Context(), "flash", "successful login")
	//resirect the user

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {
	// clean up session
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}

func (app *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {
	// create a user

	// send an activation email

	// subscbribe the user to an account
}

func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	// validate url

	// generate an invoice

	// send an email with attachments

	// send an email with the invoice attached
}
