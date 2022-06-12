package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/k3a/html2text"
)

type Connected_Status struct {
	User           string
	User_Hased     string
	Rank_Id        string
	Rank_Id_Hashed string
}

var Connected Connected_Status
var templatesDir = os.Getenv("TEMPLATES_DIR")

//#------------------------------------------------------------------------------------------------------------# ↓ Return to [Select_Page] ↓

//Return to page Selected need Path (string)
func Return_To_Page(w http.ResponseWriter, r *http.Request, Path string) {
	template.Must(template.ParseFiles(filepath.Join(templatesDir, Path))).Execute(w, " ")
}

//#------------------------------------------------------------------------------------------------------------# ↓ Return error html ↓

//Send Http error method
func Send_Error(w http.ResponseWriter, r *http.Request) {
	Return_To_Page(w, r, "../static/templates/managed_pages/404.html")
	// http.Error(w, "Method is not supported.", http.StatusBadRequest) //<-- Print [error] Method is not supported
	// fmt.Fprint(w, http.StatusBadRequest)

}

//#------------------------------------------------------------------------------------------------------------# ↓ Home Page ↓

//Home page
func forum(w http.ResponseWriter, r *http.Request) {
	type Statement_of_user struct {
		User string
		Rank string
	}
	if r.Method == "GET" {
		//<<< --- Check rank

		var (
			_, statement, User = Check_Cookie(w, r)
			pos                = Statement_of_user{}
		)
		pos.User = User
		pos.Rank = statement

		template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/forum.html"))).Execute(w, pos)

		//<<< --- Check rank
		Return_To_Page(w, r, "../static/templates/forum.html")

	} else if r.Method == "POST" {

	} else {

		Send_Error(w, r)

		return
	}
}
func edit_desc(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("")
	type Statement_of_user struct {
		User     string
		Rank     string
		Desc     string
		Descedit string
	}
	//<<< --- Check rank

	var (
		_, statement, User = Check_Cookie(w, r)
		pos                = Statement_of_user{}
	)

	//<<< --- Check rank
	if statement != "4" {

		if r.Method == "GET" {
			var (
				result        = Select_column("profil", "user", User) //Rows
				result_profil = Return_Profil(result)
			)
			//<<<<
			pos.Descedit = html2text.HTML2Text(result_profil[6])
			pos.Desc = result_profil[6]
			pos.Rank = statement
			pos.User = User

			//<<<<
			template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/managed_pages/edit_desc_profile.html"))).Execute(w, pos)

		} else if r.Method == "POST" {
			if query == "send" {
				desc_edit := r.Form["description"][0]
				fmt.Println(len(desc_edit))
				if len(desc_edit) > 2000 || len(desc_edit) == 0 {
					fmt.Fprint(w, "<script> window.alert('Description too long.'); </script>")
					fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/profil/edit"; </script>`)
					return
				}
				Update_Field("profil", "Desc", "user", User, desc_edit)
				fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/profil?=`+User+`"; </script>`)

			}

		} else {

			Send_Error(w, r)

			return
		}
	} else {
		fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
		return
	}
}

//#------------------------------------------------------------------------------------------------------------# ↓ Pages Selection & init http_serv ↓
//Server Http
func httpServ() {
	fs := http.FileServer(http.Dir("../static")) // <- ce qu'on envoie en static vers le serv
	http.Handle("/", fs)
	http.HandleFunc("/profil/edit", edit_desc)
	http.HandleFunc("/profil", profil)
	http.HandleFunc("/forum", forum)
	http.HandleFunc("/validation_mail", Validation_URLbyMail)
	http.HandleFunc("/resend_mail", Resend_Mail)
	http.HandleFunc("/admin_panel", Admin_Panel)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/valide_password", valide_password_page)
	http.HandleFunc("/reset_password_page", reset_password_page)
	http.HandleFunc("/create_post", Create_Post)
	fmt.Println("Started https serv successfully on http://localhost:1010")
	http.ListenAndServe(":1010", nil)

}
