package lfs

import (
	"bytes"
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lengsh/findme/utils"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const LFS_PATH_SEPARATOR = string(os.PathSeparator)
const LFS_FILENAMESEPARATOR = "\001\002<BR>"
const LFS_POST_FileNameKey = "file"
const LFS_POST_FilePathKey = "source"
const LFS_SYNC_FilePathKey = "date"
const LFS_SYNC_FileInodeKey = "inode"
const AUTH_SCRUMB_KEY = ".scrumb"
const URL_COMMAND_PEER_UPLOAD = "/lfs/peerupload/"
const URL_COMMAND_USER_UPLOAD = "/lfs/upload/"
const URL_COMMAND_PATH_INFO = "/lfs/pathinfo/"
const URL_COMMAND_PSYNC = "/lfs/psync/"
const URL_COMMAND_USER_TEST = "/lfs/uptest/"
const URL_COMMAND_DEFAULT = "/lfs/"
const LFS_SYNC_PathInfoKey = "pathinfo"
const LFS_SYNC_PathInfoValue = "Yes"

// static/lengfs/Node/Date/domain/filename.xyz
//    |-----|     |     |    |          |
// Parent  Pnode Inode  |  user-domain  |
//                 Create-date        fileName
//

type Node struct {
	Parent string /*    "static"     */
	Pnode  string /*    "lengfs"     */
	Inode  string /*    "0"          */
	Date   string /*    "20190501"   */
	Domain string /*    "lengsh"     */
	Queues string /*    "localhost:8081;localhost:8080"    */
}

type lfsStat struct {
	IsRunning    bool
	ModTime      time.Time
	StartTime    time.Time
	TaskPosition time.Time
	Quantity     int
	Lock         sync.RWMutex
}

var LNode = Node{}
var LfsStat = lfsStat{IsRunning: false, StartTime: time.Now().UTC().Add(8 * time.Hour), ModTime: time.Now().UTC().Add(8 * time.Hour), Quantity: 0}

func (r Node) String() string {
	return fmt.Sprintf("Parent=%s;\nPnode=%s;\nInode=%s;\nDate=%s;\nDomain=%s;\nQueues=%s;", r.Parent, r.Pnode, r.Inode, r.Date, r.Domain, r.Queues)
}

func (node Node) UserUploadFile(r *http.Request) (string, string, bool) {
	r.ParseMultipartForm(10 << 20)
	scrumb := r.FormValue(AUTH_SCRUMB_KEY)
	if !utils.CheckScrumb(scrumb) {
		logs.Debug("No Permiss: scrumb is error!....  ", AUTH_SCRUMB_KEY, " = ", scrumb)
		return "", "", false
	}
       node.Date = time.Now().UTC().Add(8*time.Hour).Format("20060102")
	if fn,nail, ok := saveFile2Node(r, node, true); ok {
		go peerSync(LNode.Parent+nail)
		go peerSync(LNode.Parent+fn)
		return fn, nail, ok
	}
	return "", "", false
}

func (node Node) PeerUploadFile(r *http.Request) (string, string, bool) {
	r.ParseMultipartForm(10 << 20)
	path := r.FormValue(LFS_POST_FilePathKey)

	scrumb := r.FormValue(AUTH_SCRUMB_KEY)
	if !utils.CheckScrumb(scrumb) {
		logs.Debug("No Authorizilaiotn.... ")
	} else {
		if len(path) > 0 {
			logs.Debug(path)
			if ok, node1 := string2Node(path); ok {
				logs.Debug("upload from(client) Node =", node1)
				return saveFile2Node(r, node1, false)
			}
		}
	}
	return "", "", false
}

func PeerSyncFile( fname string)  {
  // fname == "/lengfs/0/20190611/abc.jpeg"
  //   node.Date = time.Now().UTC().Add(8*time.Hour).Format("20060102")
    // 
    logs.Debug(LNode.Parent + fname)
    peerSync(LNode.Parent+fname)
}

func NotifyLfs2Sync(svr string, fname string)  {
logs.Debug(svr+ "/lfs/psync/?"+LFS_POST_FileNameKey+"="+fname + "&"+AUTH_SCRUMB_KEY+"="+utils.CreateScrumb())
resp, err := http.Get( svr+"/lfs/psync/?"+LFS_POST_FileNameKey+"="+fname + "&"+AUTH_SCRUMB_KEY+"="+utils.CreateScrumb())
	if err != nil {
		logs.Debug(err)
	}
	defer resp.Body.Close()
}

func (node Node) SyncPathFile(fn string) error { //    w http.ResponseWriter, r *http.Request) {
	if len(fn) > 0 {
		logs.Debug("SyncPathFile: ", LFS_SYNC_FilePathKey, " = ", fn)
		base_format := "20060102"
		_, err := time.Parse(base_format, fn)
		if err == nil {
			syncByDatePath(fn)
			return nil
		}
	}
	return fmt.Errorf("something is wrong!")
}

func (node Node) PathInfo(fdate, inode string) (string,error) {
	//      w http.ResponseWriter, r *http.Request) {
	logs.Debug("PathInfo  date = ", fdate)
	logs.Debug("PathInfo by inode = ", inode, "; current server's Inode=", LNode.Inode)
	err, rest := getFileByDatePath(inode, fdate)
	if err != nil {
		fmt.Println(err)
		//  w.WriteHeader(http.StatusInternalServerError)
		//	return
		rest = err.Error()
		return rest,err
	}
	return rest,nil
}

func saveFile2Node(r *http.Request, node Node, bCreate bool) (string, string, bool) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	// 	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile(LFS_POST_FileNameKey)
	if err != nil {
		logs.Debug("Error Retrieving the File")
		logs.Debug(err)
		return "", "", false
	}
	defer file.Close()
	logs.Debug("Uploaded File: ", handler.Filename)
	logs.Debug("File Size: ", handler.Size)
	logs.Debug("MIME Header: ", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tpath, err := getCurrentPath(node)
	if err != nil {
		logs.Debug(err)
		return "", "", false
	}
	outPath := tpath + "/" + handler.Filename
	if _, err := os.Stat(outPath); err == nil {
		logs.Debug("file is exist! : ", outPath)
		return "", "", false
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		logs.Debug(err)
		return "", "", false
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, file); err != nil {
		logs.Debug(err)
		return "", "", false
	}
	fi, err := outFile.Stat()
	if err != nil {
		logs.Debug(err)
		return "", "", false
	}

	if fi.Size() != handler.Size {
		logs.Debug("(error)file uncomplete")
		return "", "", false
	}
	surl := strings.Replace(outPath, node.Parent, "", 1)
	surl = filepath.Clean( surl)
       nail := surl 
       if bCreate {
	nail, _ = createThumbnail(outPath)
	}
logs.Debug("Successfully Uploaded File: ", handler.Filename, " ; save as ", outPath)
	return surl, nail, true
}

func mkdir(fp string) (string, error) {
	if _, err := os.Stat(fp); err == nil {
		logs.Debug(fp, " is exist!")
		return fp, nil
	} else {
		logs.Debug("create file:", fp)
		err := os.MkdirAll(fp, 0755)
		if err != nil {
			return "", err
		}
		return fp, nil
	}
}

/*
 create path by (r Node),

*/
func getCurrentPath(r Node) (string, error) {

	f1 := r.Parent + "/" + r.Pnode + "/" + r.Inode + "/" + r.Date + "/" + r.Domain
	if _, err := os.Stat(f1); err == nil {
		logs.Debug(f1, " is exist!")
		return f1, nil
	}

	f1 = r.Parent
	f2, err := mkdir(f1)
	if err != nil {
		logs.Debug(err)
		return "", err
	}

	f1 = f2 + "/" + r.Pnode
	f2, err = mkdir(f1)
	if err != nil {
		logs.Debug(err)
		return "", err
	}

	f1 = f2 + "/" + r.Inode
	f2, err = mkdir(f1)
	if err != nil {
		logs.Debug(err)
		return "", err
	}

	f1 = f2 + "/" + r.Date
	f2, err = mkdir(f1)
	if err != nil {
		logs.Debug(err)
		return "", err
	}

	f1 = f2 + "/" + r.Domain
	f2, err = mkdir(f1)
	if err != nil {
		logs.Debug(err)
		return "", err
	}
	return f2, nil
}

// exists returns whether the given file or directory exists
/*
func isExists(path string)  error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	logs.Debug(err)
	return false, err
}
*/

func peerSync(fn string) {
       logs.Debug("try to peerSync file = ",fn)
        slist := strings.Split(LNode.Queues, ";")
	ok := true
	for _, v := range slist {
		logs.Debug("peerSync to ", v, ": ", fn)
		if err := localPostFile(fn, "http://"+v+URL_COMMAND_PEER_UPLOAD); err != nil {
			ok = false
		}
	}
	if !ok {
		logs.Debug("fail to peerSync:", fn, "; server queues =", LNode.Queues )
	}
}

func localPostFile(filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	ffn := filename
	sv := strings.Split(filename, "/")
	if len(sv) > 0 {
		ffn = sv[len(sv)-1]
	}
	//fmt.Println( "sync local-file = ", filename)
	//        postfilename := strings.Replace(filename, LNode.Parent, "", 1)
	//fmt.Println( "sync lengfs-file = ", postfilename)
	scrumb := utils.CreateScrumb()

	bodyWriter.WriteField(AUTH_SCRUMB_KEY, scrumb)
	bodyWriter.WriteField(LFS_POST_FilePathKey, filename)
	// this step is very important

	fileWriter, err := bodyWriter.CreateFormFile(LFS_POST_FileNameKey, ffn) // filename)
	if err != nil {
		logs.Debug("error writing to buffer")
		return err
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		logs.Debug("error opening file")
		return err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logs.Debug(resp.Status)
	logs.Debug(string(resp_body))
	if resp.StatusCode == http.StatusOK {
		logs.Debug("peer to upload successful!")
		return nil
	} else {
		logs.Debug("peer to upload fail! ", resp.Status)
		return fmt.Errorf(resp.Status)
	}
}

func string2Node(files string) (bool, Node) {

	logs.Debug("parse path:", files)
	// ./static/lengfs/0/20190430/lengsh/a.jpg
	// ./static/lengfs/0/20190430/lengsh
	sv := strings.Split(files, string(filepath.Separator)) // "/")
	iNum := len(sv)
	if iNum < 5 {
		return false, LNode
	}

	n := Node{}
	n.Parent = LNode.Parent
	n.Pnode = sv[iNum-5]
	n.Inode = sv[iNum-4]
	n.Date = sv[iNum-3]
	n.Domain = sv[iNum-2]

	_, err := time.Parse("20060102", n.Date)

	logs.Debug("!!!!!!! ", files, ", string2Node (pay atention to Domain):", n)
	if err != nil {
		return false, LNode
	}
	return true, n
}

func syncByDatePath(datePath string) (error, string) {
	ret := ""
	reterr := fmt.Errorf("") //  nil
	if len(LNode.Queues) < 1 {
		logs.Debug("Queues is null,  exit!")
		return fmt.Errorf("Node.Queues is nuthing"), ret
	}
	sv := strings.Split(LNode.Queues, ";")
	for _, v := range sv {
		if len(v) <= 1 {
			continue
		}
		logs.Debug("try get info from", v)
		ftext, err := getRemotePathFileList(v, LNode.Inode, datePath)
		if err != nil {
			logs.Error(err)
			reterr = err
			continue
		}
		ret = ftext
		flist := strings.Split(ftext, LFS_FILENAMESEPARATOR) //  PATH_SEPARATOR)
		logs.Debug("\n\n\n\n remote file list:")
		for _, v := range flist {
			fmt.Println(v)
		}

		PthSep := string(os.PathSeparator)
		dirPth := LNode.Parent + PthSep + LNode.Pnode + PthSep + LNode.Inode + PthSep + datePath
		dir, err := ioutil.ReadDir(dirPth)
		if err != nil {
			logs.Debug("cant sync by datepath", err)
			reterr = err
			return err, ret
		}
		// suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
		for _, fi := range dir {
			if fi.IsDir() { // domain 目录
				sDomain := fi.Name()
				logs.Debug(fi.Name())
				fdir := dirPth + PthSep + fi.Name()
				logs.Debug("sync local path + domain = ", fdir)
				dir1, err := ioutil.ReadDir(fdir)
				if err != nil {
					reterr = err
					continue
				}
				for _, ff := range dir1 {
					if ff.IsDir() { // 忽略目录
					} else {

						fn := fdir + PthSep + ff.Name()
						logs.Debug(fn)
						bRun := true
						for _, fv := range flist {
							if fv == sDomain+PthSep+ff.Name() {
								bRun = false
								logs.Debug(fv, " is Exist, need't to sync!")
							}
						}
						if bRun {

							err := localPostFile(fn, "http://"+v+URL_COMMAND_PEER_UPLOAD)
							if err != nil {
								logs.Debug(err)
								reterr = err
							}
							logs.Debug("sync to remote server", fn)
							time.Sleep(2 * time.Second)
						}
					}
				}
			} else {
				//	flist = flist + sFILENAMESEPARATOR + fi.Name()
				//  忽略文件
			}
		}
	}

	if len(reterr.Error()) <= 0 {
		return nil, ret
	}
	return reterr, ret
}

func getRemotePathFileList(svr string, inode string, date string) (string, error) {
	resp, err := http.Get("http://" + svr + URL_COMMAND_PATH_INFO + "?" + LFS_SYNC_FilePathKey + "=" + date + "&" + LFS_SYNC_FileInodeKey + "=" + inode + "&" + LFS_SYNC_PathInfoKey + "=" + LFS_SYNC_PathInfoValue)
	if err != nil {
		// handle error
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return "", err
	}

	res := string(body)
	logs.Debug(res)
	return res, nil
}

func getFileByDatePath(inode, dPath string) (error, string) {
	PthSep := string(os.PathSeparator)
	if strings.Contains(dPath, PthSep) {
		return fmt.Errorf("data Path is error"), ""
	}
	dirPath := LNode.Parent + PthSep + LNode.Pnode + PthSep + inode + PthSep + dPath
        logs.Debug("get  date_path (dir files) = ", dirPath)
        dirPath = filepath.Clean(dirPath)
        logs.Debug("filepath.Clean (dirPath) = ", dirPath)
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err, ""
	}
	flist := ""
	for _, fi := range dir {
		if fi.IsDir() { // domain 目录
			logs.Debug(fi.Name())
			sDomain := fi.Name()
			fdir := dirPath + PthSep + fi.Name()
			fdir = filepath.Clean(fdir)
			logs.Debug("get date_path + domain = ", fdir)
			dir1, err := ioutil.ReadDir(fdir)
			if err != nil {
				continue
			}
			for _, ff := range dir1 {
				if ff.IsDir() { // 忽略目录
				} else {
					logs.Debug(ff.Name())
					if len(flist) > 0 {
						flist = flist + LFS_FILENAMESEPARATOR + sDomain + PthSep + ff.Name()
					} else {
						flist = sDomain + PthSep + ff.Name()
					}
				}
			}
		} else {
			//	flist = flist + sFILENAMESEPARATOR + fi.Name()
			//  忽略文件
		}
	}
	return nil, flist
}

func JobWatch(ctx context.Context, iT int) {
	if iT < 60*5 {
		iT = 60*5
	}
	for {
		select {
		case <-ctx.Done():
			logs.Debug("lfs background sync deamon优雅退出...")
			//logs.Debug("graceful quit ...")
			return
		case <-time.After(time.Duration(iT) * time.Second):
			logs.Debug(iT, " Seconds scan jobs ...(", time.Now().Format("2006-01-02 15:04:05"), ")")
			//if be should to do task ,then do !
			LfsStat.Lock.Lock()
			if !LfsStat.IsRunning {
				LfsStat.IsRunning = true
				go backgroundSync()
			} else {
				logs.Debug("backgroudSync is running...., exit current times")
			}
			LfsStat.Lock.Unlock()
		}
	}
}

func releaseBackSyncLock() {
	LfsStat.Lock.Lock()
	if LfsStat.IsRunning {
		LfsStat.IsRunning = false
	}
	LfsStat.Lock.Unlock()
}

func backgroundSync() error {

	LfsStatModify()
	logs.Debug("backgroundSync, current iNode = ", LNode.Inode)
	//  get remote inode info by Node.Queues
	//  check current inode info
	//  if local_inode_info != remote_inode_info
	//  peerSync again
	//   case <-time.After(time.Duration( 600 ) * time.Second):
	//   break;

	defer releaseBackSyncLock()

	PthSep := string(os.PathSeparator)
	dirPth := LNode.Parent + PthSep + LNode.Pnode + PthSep + LNode.Inode
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return err
	}
	// suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			logs.Debug("backgroundSync path = ", fi.Name())
			syncByDatePath(fi.Name()) //   dirPth + PthSep + fi.Name())
		} else {
		}
	}
	return nil
}

func GetLfsStatQuantity() int {
	LfsStat.Lock.RLock()
	i := LfsStat.Quantity
	LfsStat.Lock.RUnlock()
	return i
}

func GetLfsStatStart() time.Time {
	LfsStat.Lock.RLock()
	i := LfsStat.StartTime
	LfsStat.Lock.RUnlock()
	return i
}

func GetLfsStatModTime() time.Time {
	LfsStat.Lock.RLock()
	i := LfsStat.ModTime
	LfsStat.Lock.RUnlock()
	return i
}

func LfsStatStart() {
	LfsStat.Lock.RLock()
	LfsStat.StartTime = time.Now().UTC().Add(8 * time.Hour)
	LfsStat.Lock.RUnlock()
}

func LfsStatModify() {
	LfsStat.Lock.Lock()
	LfsStat.ModTime = time.Now().UTC().Add(8 * time.Hour)
	LfsStat.Lock.Unlock()
}
func LfsStatQuantityIncr() {
	LfsStat.Lock.Lock()
	LfsStat.Quantity++
	LfsStat.Lock.Unlock()
}

func createThumbnail(fn string) (string, error) {
	fws := path.Base(fn) //获取文件名带后缀
	fex := path.Ext(fws) //获取文件后缀
	fex1 := strings.ToLower(fex)
	if fex1 != ".png" && fex1 != ".gif" && fex1 != ".jpg" && fex1 != ".jpeg" {
		return "", fmt.Errorf("error file type")
	}
	file, err := os.Open(fn)
	if err != nil {
		logs.Error(err)
		return "", err
	}
	defer file.Close()

	//:= jpeg.Decode(file)
	img, _, err := image.Decode(file)
	if err != nil {
		logs.Error(err)
		return "", err
	}
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Thumbnail(320, 1200, img, resize.Lanczos3)
	nailName := strings.TrimSuffix(fn, fex)
	nailName = nailName + ".thumbnail" + fex
	out, err := os.Create(nailName)
	if err != nil {
		logs.Error(err)
		return "", err
	}
	defer out.Close()
	// write new image to file
	switch fex1 {
	case ".gif":
		gif.Encode(out, m, nil)
	case ".jpeg":
		fallthrough
	case ".jpg":
		jpeg.Encode(out, m, nil)
	case ".png":
		png.Encode(out, m)
	}

	surl := strings.Replace(nailName, LNode.Parent, "", 1)
	return surl, nil
}
