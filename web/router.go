package web

import (
	"github.com/lengsh/findme/user"
	"net/http"
)

func Router() {
	//	http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir(utils.ServerConfig.WebDir+"doc"))))
	lfs_router_register()

	http.HandleFunc("/", lfsStat)
	//	http.HandleFunc("/hello/", hello)
	http.HandleFunc("/login", login)
	http.HandleFunc("/html/", hellohtml)
	http.HandleFunc("/api/userinfo/v1/", user.JwtAuth(userinfo))

}
