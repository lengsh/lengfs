package web

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lengsh/findme/user"
	"github.com/lengsh/findme/utils"
	"github.com/lengsh/lengfs/lfs"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// static/lengfs/Node/Date/domain/filename.xyz
//   |------|     |    |     |         |
// Parent       Inode  |  user-domain  |
//        Pnode    Create-date       fileName
//

const pathSep = string(os.PathSeparator)

func lfs_router_register() {
//	PthSep := string(os.PathSeparator)
	lengfs := pathSep + lfs.LNode.Pnode + pathSep // such as  "/lengfs/"

	if !strings.HasPrefix(lfs.LNode.Parent, pathSep) {
		err := "lfs.LNode.Parent must be absolute path: begin with '/', but now it's = " + lfs.LNode.Parent
		logs.Error(err)
		panic(err)
	}
	lfsDir := lfs.LNode.Parent + lengfs
	fmt.Println("fs path:", lfsDir, "  --> http file dir locate !")
	http.Handle(lengfs, http.StripPrefix(lengfs, http.FileServer(http.Dir(lfsDir))))
	/*****
	  lfs web command
	  ******/
	http.HandleFunc(lfs.URL_COMMAND_USER_TEST, upload_test)
	http.HandleFunc(lfs.URL_COMMAND_USER_UPLOAD, upload)
	http.HandleFunc(lfs.URL_COMMAND_PSYNC, pathSync)
	http.HandleFunc(lfs.URL_COMMAND_PATH_INFO, pathInfo)
	http.HandleFunc(lfs.URL_COMMAND_PEER_UPLOAD, peerUpload)
	http.HandleFunc(lfs.URL_COMMAND_DEFAULT, lfsStat)
	/**********/
}

type imglist struct {
	Thumbnail string
	Original  string
}

func getTemplate(view string, r *http.Request) (*template.Template, error) {
	var viewRoot string
	if !isMobile(r) {
		viewRoot = utils.ServerConfig.WebDir + "pview/"
	} else {
		viewRoot = utils.ServerConfig.WebDir + "mview/"
	}
	return template.New(view).Funcs(template.FuncMap{
		"UserName": func() template.HTML {
			return template.HTML("lengsh")
		},
		"Scrumb": func() template.HTML {
			return template.HTML(utils.CreateScrumb())
		},
	"IsContains": func(fn string, sub string) bool {
			return  strings.Contains(fn, sub)
		},
		"DateFormat": func(t time.Time, format string) template.HTML {
			return template.HTML(t.Format(format))
		},
	}).ParseFiles(viewRoot+view, viewRoot+"header.gtpl", viewRoot+"footer.gtpl")
}

func upload_test(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.RequestURI, "/favicon.ico") {
		return
	}

	if !user.IsValid(r) {
		http.Redirect(w, r, "/login?goback="+lfs.URL_COMMAND_USER_UPLOAD, http.StatusFound)
		return
	}

	node := lfs.LNode
	u, _ := user.GetUser(r)
	if usr, ok := u["user"]; ok {
		node.Domain = strings.Replace(usr.(string), "@", ".", 10)
	}
	r.ParseForm()
	//   fmt.Println(r.Form)
	//  fmt.Println(r.PostForm)
       //var data = make(map[string] interface{}, 1)
root := lfs.LNode.Parent + "/" + lfs.LNode.Pnode + "/" + lfs.LNode.Inode + "/" + time.Now().UTC().Add(8*time.Hour).Format("20060102") + "/"
	//data["root"] = root
	//data["data"] = getDatefiles( root )
	data := getDatefiles(root)
	t, er := getTemplate("uptest.gtpl", r)
	if er != nil {
		logs.Error(er)
		return
	}
	err := t.Execute(w, data)
	if err != nil {
		logs.Error(err.Error())
	}
}
func upload(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.RequestURI, "/favicon.ico") {
		return
	}

	if !user.IsValid(r) {
		http.Redirect(w, r, "/login?goback="+lfs.URL_COMMAND_USER_UPLOAD, http.StatusFound)
		return
	}

	node := lfs.LNode
	u, _ := user.GetUser(r)
	if usr, ok := u["user"]; ok {
		node.Domain = strings.Replace(usr.(string), "@", ".", 10)
	}

	node.Date = time.Now().Format("20060102")

	data := imglist{}
	if r.Method == "POST" {
		if url, nail, ok := node.UserUploadFile(r); ok {
			data.Original = url
			data.Thumbnail = nail
		} else {
			logs.Debug("fail to upload file")
		}
		if gotos := r.FormValue("goto"); len(gotos) > 0 { // has goto url
			gourl := gotos + "?original=" + data.Original + "&thumbnail=" + data.Thumbnail
			http.Redirect(w, r, gourl, http.StatusFound)
			return
		}

	}
	/*
		var viewRoot string
		if !isMobile(r) {
			viewRoot = utils.ServerConfig.WebDir + "pview/"
		} else {
			viewRoot = utils.ServerConfig.WebDir + "mview/"
		}
		fmt.Println(viewRoot)
		t, _ := template.New("upload.gtpl").Funcs(template.FuncMap{
			"UserName": func() template.HTML {
				return template.HTML(user.UserNameScript(r))
			},
		}).ParseFiles(viewRoot+"upload.gtpl", viewRoot+"header.gtpl", viewRoot+"footer.gtpl")
	*/

	t, er := getTemplate("upload.gtpl", r)
	if er != nil {
		logs.Error(er)
		return
	}
	err := t.Execute(w, data)
	if err != nil {
		logs.Error(err.Error())
	}

}

func lfsStat(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.RequestURI, "/favicon.ico") {
		return
	}
	/*
		view, err := statView()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.Write(view)
	*/
	data := map[string]interface{}{"Node": lfs.LNode, "Stat": lfs.LfsStat}
	t, er := getTemplate("stat.gtpl", r)
	if er != nil {
		logs.Error(er)
		return
	}
	err := t.Execute(w, data)
	if err != nil {
		logs.Error(err.Error())
	}
}

func peerUpload(w http.ResponseWriter, r *http.Request) {
	/*
	           if !user.IsValid(r) {
	   		fmt.Println( http.StatusText(http.StatusNonAuthoritativeInfo))
	   		w.WriteHeader(http.StatusNonAuthoritativeInfo)
	   		return
	   	}
	*/
	fmt.Println("peer upload request now!")
	if r.Method == "POST" {
		node := lfs.LNode
		node.PeerUploadFile(r)
	}
}

func pathInfo(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() //解析url传递的参数，对于POST则解析响应包的主体（request body）
	//   date=?&inode=?&.scrumb=?
        data := getDefaultData(w,r)

	beCommand := r.FormValue(lfs.LFS_SYNC_PathInfoKey)
	fdate := r.FormValue(lfs.LFS_SYNC_FilePathKey)
	inode := r.FormValue(lfs.LFS_SYNC_FileInodeKey)
	if len(fdate) <= 1 || len(inode) < 1 {
		logs.Debug("error to get param by ", lfs.LFS_SYNC_FilePathKey)
		w.WriteHeader(http.StatusInternalServerError)
	}

	rest := lfs.LNode.PathInfo(fdate, inode)
	if beCommand == lfs.LFS_SYNC_PathInfoValue {
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		w.Write([]byte(rest))
		return
	}
	flist := strings.Split(rest, lfs.LFS_FILENAMESEPARATOR)
	for k, v := range flist {
		fmt.Println(k, " = ", v)
	}
  root :=  lfs.LNode.Pnode + "/" + lfs.LNode.Inode + "/" + time.Now().UTC().Add(8*time.Hour).Format("20060102") + "/"
 // data := map[string]interface{}{ "data":flist, "root": root}
 data["data"] = flist
 data["root"] = root

	t, er := getTemplate("pathinfo.gtpl", r)
	if er != nil {
		logs.Error(er)
		return
	}
	err := t.Execute(w, data)
	if err != nil {
		logs.Error(err.Error())
	}
	/*
	   logs.Debug(rest)
	   w.Header().Set("Content-Type", "text/html;charset=utf-8")
	   w.Write([]byte(rest))
	*/
}

func pathSync(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() //解析url传递的参数，对于POST则解析响应包的主体（request body）
	fn := r.FormValue(lfs.LFS_SYNC_FilePathKey)
	lfs.LNode.SyncPathFile(fn)
	// data := map[string]interface{}{ "Node":lfs.LNode,"Stat":lfs.LfsStat}
	t, er := getTemplate("sync.gtpl", r)
	if er != nil {
		logs.Error(er)
		return
	}
	err := t.Execute(w, nil)
	if err != nil {
		logs.Error(err.Error())
	}
}

func isMobile(r *http.Request) bool {

	if agent, ok := r.Header["User-Agent"]; ok {
		//  " /Android|webOS|iPhone|iPod|BlackBerry/i.test(navigator.userAgent))"
		if strings.Contains(agent[0], "Android") || strings.Contains(agent[0], "iPhone") {
			return true
		}
	}
	return false
}

func getDatefiles( root string ) []imglist {
	var files []imglist
	// root := lfs.LNode.Parent + "/" + lfs.LNode.Pnode + "/" + lfs.LNode.Inode + "/" + time.Now().UTC().Add(8*time.Hour).Format("20060102") + "/"
	// root := LNode.Parent+  "./lengfs/0/20190528/"
	fmt.Println("walk root = ", root)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".thumbnail.") {
			f := strings.Replace(path, lfs.LNode.Parent, "", 1)
			img := imglist{Thumbnail: f, Original: strings.Replace(f, ".thumbnail", "", 1)}
			files = append(files, img)
		}
		return nil
	})
	if err != nil {
		logs.Error(err)
	}
	fmt.Println(files)
	return files
}
