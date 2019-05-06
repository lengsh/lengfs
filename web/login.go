package web

import (
	"bytes"
	"fmt"
	"github.com/lengsh/findme/user"
	"github.com/nilslice/jwt"
	"html/template"
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
		view, err := loginView(r.URL.RequestURI(), eusr)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write(view)

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
			errView, err := ErrorMessage("title", "message")
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, r.URL.String(), http.StatusFound)
				return
			}
			w.Write(errView)
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

var startAdminHTML = `<!doctype html>
<html lang="en">
    <head>
        <title> Hello, Logo  </title>
        <link rel="stylesheet" href="/static/dashboard/css/admin.css" />
        <link rel="stylesheet" href="/static/dashboard/css/materialize.min.css" />
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
    </head>
    <body class="grey lighten-4">
       <div class="navbar-fixed">
            <nav class="grey darken-2">
            <div class="nav-wrapper">
                <a class="brand-logo" href="/">Logo</a>

                <ul class="right">
                    <li><a href="/logout">Logout</a></li>
                </ul>
            </div>
            </nav>
        </div>

        <div class="admin-ui row">`

var endAdminHTML = `
        </div>
        <footer class="row">
            <div class="col s12">
                <p class="center-align"> © 2019 <a target="_blank" href="http://www.lengsh.cn/">lengsh</a>. <京ICP备19007810号> </p>
            </div>
        </footer>
    </body>
</html>`

var loginAdminHTML = `
<div class="init col s5">
<div class="card">
<div class="card-content">
    <div class="card-title">Welcome {{.usr.Email}}!</div>
    <blockquote>请输入登录账号和密码.</blockquote>
    <form method="post" action="{{.action}}" class="row">
        <div class="input-field col s12">
            <input placeholder="Enter your email address e.g. you@example.com" class="validate required" type="email" id="email" name="email"/>
            <label for="email" class="active">Email</label>
        </div>
        <div class="input-field col s12">
            <input placeholder="Enter your password" class="validate required" type="password" id="password" name="password"/>
            <a href="/recover">Forgot password?</a>            
            <label for="password" class="active">Password</label>  
        </div>
        <button class="btn waves-effect waves-light right">Log in</button>
    </form>
</div>
</div>
</div>
<script>
    $(function() {
        $('.nav-wrapper ul.right').hide();
    });
</script>
`

func loginView(postUrl string, email string) ([]byte, error) {
	html := startAdminHTML + loginAdminHTML + endAdminHTML
	buf := &bytes.Buffer{}
	usr := user.User{Email: email}
	data := map[string]interface{}{
		"usr":    usr,
		"action": postUrl,
	}

	tmpl := template.Must(template.New("login").Parse(html))
	err := tmpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ErrorMessage is a generic error message container, similar to Error500() and
// others in this package, ecxept it expects the caller to provide a title and
// message to describe to a view why the error is being shown
func ErrorMessage(title, message string) ([]byte, error) {

	var errMessageHTML = `
 <div class="error-page eMsg col s6">
 <div class="card">
 <div class="card-content">
     <div class="card-title"><b>Error:&nbsp;</b>%s</div>
     <blockquote>%s</blockquote>
 </div>
 </div>
 </div>
 `

	eHTML := fmt.Sprintf(errMessageHTML, title, message)
	return []byte(eHTML), nil
}
