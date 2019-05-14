package web

import (
	"net/http"
"github.com/lengsh/findme/user"
)

func Router() {
	//	http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir(utils.ServerConfig.WebDir+"doc"))))
	lfs_router_register()

	http.HandleFunc("/", lfsStat)
	http.HandleFunc("/hello/", hello)
	http.HandleFunc("/login", login)
	http.HandleFunc("/hello/", hello)
http.HandleFunc("/api/userinfo/v1/", user.JwtAuth(userinfo))

}
