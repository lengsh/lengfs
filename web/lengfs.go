package web

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lengsh/findme/user"
	"github.com/lengsh/findme/utils"
	"github.com/lengsh/lengfs/lfs"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"
)

// static/lengfs/Node/Date/domain/filename.xyz
//   |------|     |    |     |         |
// Parent       Inode  |  user-domain  |
//        Pnode    Create-date       fileName
//

func lfs_router_register() {
	PthSep := string(os.PathSeparator)
	lengfs := PthSep + lfs.LNode.Pnode + PthSep //  "/lengfs/"
	if len(utils.ServerConfig.WebDir) < 1 {
		utils.ServerConfig.WebDir = "." + PthSep
	}
	lengfs_Local_dir := ""
	if strings.HasSuffix(utils.ServerConfig.WebDir, PthSep) {
		lengfs_Local_dir = utils.ServerConfig.WebDir
	} else {
		lengfs_Local_dir = utils.ServerConfig.WebDir + PthSep
	}

	logs.Debug("1: fs dir:", lengfs_Local_dir)

	if strings.HasPrefix(lfs.LNode.Parent, ".") {
		lengfs_Local_dir = lengfs_Local_dir + strings.Replace(lfs.LNode.Parent, "./", "", 1) //   "./static/lengfs/"
		logs.Debug(lfs.LNode.Parent, "\n", lengfs_Local_dir)
	} else {
		lengfs_Local_dir = lengfs_Local_dir + lfs.LNode.Parent //   "./static/lengfs/"
	}
	logs.Debug("2: fs dir:", lengfs_Local_dir)
	lengfs_Local_dir += lengfs
	logs.Debug("3: fs dir:", lengfs_Local_dir)
	http.Handle(lengfs, http.StripPrefix(lengfs, http.FileServer(http.Dir(lengfs_Local_dir))))
	/*****
	  lfs web command
	  ******/
	http.HandleFunc(lfs.URL_COMMAND_USER_UPLOAD, upload)
	http.HandleFunc(lfs.URL_COMMAND_PSYNC, pathSync)
	http.HandleFunc(lfs.URL_COMMAND_PATH_INFO, pathInfo)
	http.HandleFunc(lfs.URL_COMMAND_PEER_UPLOAD, peerUpload)
	http.HandleFunc(lfs.URL_COMMAND_DEFAULT, lfsStat)
	/**********/
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

	fview := ""
	if r.Method == "POST" {
		if ok, url := node.UserUploadFile(w, r); ok {
			logs.Debug(url)
			fview = url
		} else {
			logs.Debug("fail to upload file")
		}
	}
	view, err := uploadView(fview)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(view)
}

func lfsStat(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.RequestURI, "/favicon.ico") {
		return
	}

	view, err := statView()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(view)
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
		node.PeerUploadFile(w, r)
	}
}

func pathInfo(w http.ResponseWriter, r *http.Request) {
	lfs.LNode.PathInfo(w, r)
}

func pathSync(w http.ResponseWriter, r *http.Request) {
	lfs.LNode.SyncPathFile(w, r)
}

func uploadView(fn string) ([]byte, error) {
	html := `
<HTML>
<head>
	<meta charset="utf-8" />
	<title>lengfs</title>
</head>
<body>
 <div align=center>
    <form enctype="multipart/form-data" action="/lfs/upload/" method="post" >
         <input type="file" name="file" />
         <input type="submit" value="upload" />
	 <input type="hidden" name=".scrumb" value="` + utils.CreateScrumb() + `"  />
</form>
</div>
`
	if len(fn) > 0 {
		str := ` <div><h1>Successful upload file's URL</h1> `
		str2 := `<a href="` + fn + `">` + fn + "</a> </div>"
		html = html + str + str2
	}

	html += `
<BR>
<a href="/lfs/">lengfs Home</a>
</body>
</HTML>
 `
	buf := &bytes.Buffer{}
	tmpl := template.Must(template.New("upload").Parse(html))
	err := tmpl.Execute(buf, nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func statView() ([]byte, error) {
	fHts := "<div><h1> curent lfs's dir path</h1></div>"
	fHts += "<a href='/" + lfs.LNode.Pnode + "/'>" + lfs.LNode.Pnode + "</a><BR>"
	stat := "<div><H1>current INode = " + lfs.LNode.Inode + "</H1></div>"
	stat += "<div>Start = " + lfs.GetLfsStatStart().Format("2006-01-02 15:04:05") + "<BR>" + "modtime = " + lfs.GetLfsStatModTime().Format("2006-01-02 15:04:05") + "<BR></div>"

	str := `
<HTML>
<head>
	<meta charset="utf-8" />
	<title>lengfs</title>
</head>
<body>
   <div align=left>
   <H1>Upload Web</h1>
    <form enctype="multipart/form-data" action="/lfs/upload/" method="post" >
         <input type="file" name="file" />
         <input type="submit" value="upload" />
	 <input type="hidden" name=".scrumb" value="` + utils.CreateScrumb() + `"  />
</form>
  <BR>
  <div align=left>  
     <H1>All lfs's COMMAND:</h1>
        <a href="/lfs/">lengfs</a><BR>
        <a href="/login">login</a><BR>
        <a href="/lfs/upload/">lfs/upload</a><BR>
        <a href="/lfs/psync/?date=20190430">/lfs/psync/?date=20190430</a><BR>
        <a href="/lfs/pathinfo/?date=20190430&inode=0">/lfs/pathinfo/?date=20190430&inode=0</a><BR>
  </div>
`
	str2 := stat + `</body> </HTML> `
	str = str + fHts
	html := ""
	html = str + str2
	buf := &bytes.Buffer{}
	tmpl := template.Must(template.New("upload").Parse(html))
	err := tmpl.Execute(buf, nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
