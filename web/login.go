package web

import (
	"fmt"
	"github.com/lengsh/findme/user"
	"github.com/nilslice/jwt"
	//	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

func login(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.RequestURI, "/favicon.ico") {
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, r.URL.String(), http.StatusFound)
		return
	}
	goback := r.FormValue("goback")
	eusr := ""
	if usr, err := user.GetUser(r); err == nil {
		for k, v := range usr {
			fmt.Println(k, " = ", v)
		}
		eusr = fmt.Sprintf("%s", usr["user"])
	}

	switch r.Method {
	case http.MethodGet:
		data := map[string]string{"user": eusr}
		t, er := getTemplate("login.gtpl", r)
		if er != nil {
			fmt.Println(er)
			return
		}
		err := t.Execute(w, data)
		if err != nil {
			fmt.Println(err.Error())
		}

		/*
			view, err := loginView(r.URL.RequestURI(), eusr)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html")
			w.Write(view)
		*/
	case http.MethodPost:
		// check email & password
		j := strings.ToLower(r.FormValue("email"))
		if len(j) <= 0 {
			log.Println("no email ")
			http.Redirect(w, r, r.URL.String(), http.StatusFound)
			return
		}

		//  这里应该是从DB中取出用户注册数据
		/*
				        if  usr,err :=  db.Get(email); err != nil {
					IsExist := false
					}
			                if !IsUser(usr, password ){

					}
		*/

		//  新用户，则生成并保持基础数据
		////////////////////////////////////////
		//
		//              <@-----@>
		//              (      )
		//               ======
		//   SET PASSWORD == "12345678901234567890"
		//
		//
		////////////////////////////////////////////////////
		usr, err := user.New(j, r.FormValue("password"))
		if err != nil || r.FormValue("password") != "12345678901234567890" {
			log.Println(err)
			http.Redirect(w, r, r.URL.String(), http.StatusFound)
			return
		} else {
			// 	db.Save(usr)
		}

		// create new token
		week := time.Now().Add(time.Hour * 24 * 7)
		claims := map[string]interface{}{
			"exp":  week,
			"user": usr.Email,
		}
		token, err := jwt.New(claims)
		if err != nil {
			log.Println(err)
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, r.URL.String(), http.StatusFound)
				return
			}
			w.Write([]byte("Error!"))
			return
		}
		// add it to cookie +1 week expiration
		http.SetCookie(w, &http.Cookie{
			Name:    "_token",
			Value:   token,
			Expires: week,
			Path:    "/",
		})

		if len(goback) > 0 {
			http.Redirect(w, r, r.URL.Scheme+r.URL.Host+goback, http.StatusFound)
		} else {

			http.Redirect(w, r, r.URL.Scheme+r.URL.Host, http.StatusFound)
		}
		//http.Redirect(w, r, strings.TrimSuffix(r.URL.String(), "/login"), http.StatusFound)
	}
}
