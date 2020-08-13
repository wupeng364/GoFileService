// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 文件API接口模块, 文件的新建、删除、移动、复制等操作

package fileapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"gofs/base/httpserver"
	"gofs/base/signature"
	"gofs/data/filemanage"
	"gofs/data/filepermission"
	"gofs/service/restful/fileapi/asynctask"
	"gutils/hstool"
	"gutils/mloader"
	"gutils/strtool"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
)

// FileAPI 文件API接口模块
type FileAPI struct {
	fm   *filemanage.FileManager
	hs   *httpserver.HTTPServer
	sg   *signature.Signature
	ast  *asynctask.AsyncTasks
	fpms *filepermission.FPmsManager
}

// ModuleOpts 模块加载器接口实现, 返回模块信息&配置Name:
func (fsapi *FileAPI) ModuleOpts() mloader.Opts {
	return mloader.Opts{
		Name:        "FileAPI",
		Version:     1.0,
		Description: "文件管理API接口模块",
		OnReady: func(mctx *mloader.Loader) {
			mctx.Load(new(asynctask.AsyncTasks))
			fsapi.fpms = mctx.GetModuleByTemplate(fsapi.fpms).(*filepermission.FPmsManager)
			fsapi.ast = mctx.GetModuleByTemplate(fsapi.ast).(*asynctask.AsyncTasks)
			fsapi.fm = mctx.GetModuleByTemplate(fsapi.fm).(*filemanage.FileManager)
			fsapi.hs = mctx.GetModuleByTemplate(fsapi.hs).(*httpserver.HTTPServer)
			fsapi.sg = mctx.GetModuleByTemplate(fsapi.sg).(*signature.Signature)

		},
		OnInit: fsapi.init,
	}
}

// 向 HttpServerModule 中注册服务地址
func (fsapi *FileAPI) init() {
	// 批量注册Handle
	err := fsapi.hs.AddRegistrar(fsapi)
	if err != nil {
		panic(err)
	}

	// 注册Api签名拦截器
	fsapi.hs.AddURLFilter(baseurl+"/:"+`[\S]+`, fsapi.sg.RestfulAPIFilter)

	fmt.Println("   > FileApiModule http registered end")
}

// RoutList 向 Server Router 中注册下列处理器, 实现接口 httpserver.Registrar
func (fsapi *FileAPI) RoutList() httpserver.StructRegistrar {
	fsapi.hs.AddHandlerFunc(baseurl+"/stream/:"+`[\S]+`, fsapi.Stream)
	return httpserver.StructRegistrar{
		IsToLower: true,
		BasePath:  baseurl,
		HandlerFunc: []hstool.HandlersFunc{
			fsapi.Info,
			fsapi.List,
			fsapi.Del,
			fsapi.ReName,
			fsapi.NewFolder,
			fsapi.AsyncExec,
			fsapi.AsyncExecToken,
			fsapi.Stream,
			fsapi.StreamToken,
		},
	}
}

// checkPermision 检查权限
func (fsapi *FileAPI) checkPermision(userID, path string, permission int64) bool {

	return fsapi.fpms.HashPermission(userID, path, permission)
}

// Info 获取文件|文件夹信息
func (fsapi *FileAPI) Info(w http.ResponseWriter, r *http.Request) {
	qpath := path.Clean(r.FormValue("path"))
	if !fsapi.checkPermision(fsapi.sg.GetUserID4Request(r), qpath, filepermission.VISIBLE) {
		httpserver.SendError(w, ErrorPermissionInsufficient)
		return
	}
	if len(qpath) == 0 {
		httpserver.SendError(w, errors.New("path is empty"))
		return
	}
	fs, err := fsapi.fm.GetPathDriver(qpath)
	if nil != err {
		httpserver.SendError(w, err)
		return
	}
	if is, err := fs.IsExist(qpath); !is || nil != err {
		httpserver.SendError(w, errors.New(qpath+" is not exist"))
		return
	}
	isFile, err := fs.IsFile(qpath)
	if nil != err {
		httpserver.SendError(w, err)
		return
	}
	info := filemanage.FsInfo{
		Path:   qpath,
		CtTime: (func() int64 { res, _ := fs.GetModifyTime(qpath); return res })(),
		IsFile: isFile,
		FileSize: (func() int64 {
			if !isFile {
				return 0
			}
			res, _ := fs.GetFileSize(qpath)
			return res
		})(),
	}
	json, err := json.Marshal(info)
	if nil != err {
		httpserver.SendError(w, err)
		return
	}
	httpserver.SendSuccess(w, string(json))
}

// List 查询路径下的列表以及基本信息
func (fsapi *FileAPI) List(w http.ResponseWriter, r *http.Request) {
	qpath := r.FormValue("path")
	qSort := r.FormValue("sort")
	qAsc := r.FormValue("asc")
	// fmt.Println(r.URL, r.Body, qpath)
	if len(qpath) == 0 {
		httpserver.SendError(w, errors.New("path is empty"))
		return
	}
	if len(qAsc) == 0 {
		qAsc = "true"
	}
	userID := fsapi.sg.GetUserID4Request(r)
	if !fsapi.checkPermision(userID, qpath, filepermission.VISIBLECHILD) {
		httpserver.SendError(w, ErrorPermissionInsufficient)
		return
	}
	list, err := fsapi.fm.GetDirListInfo(qpath)
	if err != nil {
		httpserver.SendError(w, err)
		return
	}
	// 如果当前或上级路径有可见以上权限, 则文件默认可见
	canVisible := fsapi.checkPermision(userID, qpath, filepermission.VISIBLE)
	res := make([]filemanage.FsInfo, 0)
	if len(list) > 0 {
		for i := 0; i < len(list); i++ {
			if list[i].IsFile && canVisible {
				res = append(res, list[i])
				continue
			}
			if fsapi.checkPermision(userID, list[i].Path, filepermission.VISIBLECHILD) {
				res = append(res, list[i])
			}
		}
	}
	json, err := json.Marshal(filemanage.FileListSorter{
		Asc:       strtool.String2Bool(qAsc),
		SortField: qSort,
	}.Sort(res))
	if err != nil {
		httpserver.SendError(w, err)
	}
	httpserver.SendSuccess(w, string(json))
}

// Del 批量|单个删除文件|文件夹
func (fsapi *FileAPI) Del(w http.ResponseWriter, r *http.Request) {
	qpath := r.FormValue("path")
	if len(qpath) == 0 {
		httpserver.SendError(w, errors.New("path is empty"))
		return
	}
	if !fsapi.checkPermision(fsapi.sg.GetUserID4Request(r), qpath, filepermission.WRITE) {
		httpserver.SendError(w, ErrorPermissionInsufficient)
		return
	}
	if !fsapi.fm.IsExist(qpath) {
		httpserver.SendError(w, ErrorFileNotExist)
		return
	}
	err := fsapi.fm.DoDelete(qpath)
	if nil == err {
		httpserver.SendSuccess(w, "")
	} else {
		httpserver.SendError(w, err)
	}
}

// ReName 重命名
func (fsapi *FileAPI) ReName(w http.ResponseWriter, r *http.Request) {
	qSrcPath := r.FormValue("path")
	qName := r.FormValue("name")
	if len(qName) == 0 {
		httpserver.SendError(w, ErrorNewNameIsEmpty)
		return
	}
	if !fsapi.checkPermision(fsapi.sg.GetUserID4Request(r), qSrcPath, filepermission.WRITE) {
		httpserver.SendError(w, ErrorPermissionInsufficient)
		return
	}
	if !fsapi.fm.IsExist(qSrcPath) {
		httpserver.SendError(w, ErrorFileNotExist)
		return
	}
	rnmErr := fsapi.fm.DoRename(qSrcPath, qName)
	if nil == rnmErr {
		httpserver.SendSuccess(w, "")
	} else {
		httpserver.SendError(w, rnmErr)
	}
}

// NewFolder 新建文件夹
func (fsapi *FileAPI) NewFolder(w http.ResponseWriter, r *http.Request) {
	qPath := r.FormValue("path")
	qPath = strtool.Parse2UnixPath(qPath)
	if !fsapi.checkPermision(fsapi.sg.GetUserID4Request(r), qPath, filepermission.WRITE) {
		httpserver.SendError(w, ErrorPermissionInsufficient)
		return
	}
	if !fsapi.fm.IsExist(strtool.GetPathParent(qPath)) {
		httpserver.SendError(w, ErrorParentFolderNotExist)
		return
	}
	rnmErr := fsapi.fm.DoNewFolder(qPath)
	if nil == rnmErr {
		httpserver.SendSuccess(w, "")
	} else {
		httpserver.SendError(w, rnmErr)
	}
}

// StreamToken 传输令牌申请
func (fsapi *FileAPI) StreamToken(w http.ResponseWriter, r *http.Request) {
	qdata := r.FormValue("data")
	qtype := r.FormValue("type")
	if len(qtype) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if qtype == "download" || qtype == "stream" {
		if !fsapi.checkPermision(fsapi.sg.GetUserID4Request(r), qdata, filepermission.READ) {
			httpserver.SendError(w, ErrorPermissionInsufficient)
			return
		}

	} else if qtype == "upload" {
		if !fsapi.checkPermision(fsapi.sg.GetUserID4Request(r), qdata, filepermission.WRITE) {
			httpserver.SendError(w, ErrorPermissionInsufficient)
			return
		}

	} else {
		httpserver.SendError(w, ErrorPermissionInsufficient)
		return
	}
	httpserver.SendSuccess(w, fsapi.fm.AskToken(qtype, &StreamToken{
		Type: qtype,
		Data: qdata,
	}))
}

// parseFileSteamToken 获取文件传输Token对象
func (fsapi *FileAPI) parseFileSteamToken(token string) (*StreamToken, error) {
	tokenBody := fsapi.fm.GetToken(token)
	// 并保持刷新token的有效性, 除非终止操作否则都继续
	if nil == tokenBody {
		return nil, ErrorOprationExpires
	}
	ftt, ok := tokenBody.(*StreamToken)
	if ok {
		return ftt, nil
	}
	return nil, ErrorOprationExpires
}

// Stream 文件上传|下载
func (fsapi *FileAPI) Stream(w http.ResponseWriter, r *http.Request) {
	qMethod := strings.ToLower(r.Method)
	if qMethod == "post" || qMethod == "put" {
		fsapi.upload(w, r)
	} else if qMethod == "get" {
		fsapi.download(w, r)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// Upload 文件上传, 支持Form和Body上传方式
func (fsapi *FileAPI) upload(w http.ResponseWriter, r *http.Request) {
	TObject, err := fsapi.parseFileSteamToken(r.FormValue("token"))
	if nil != err || nil == TObject {
		httpserver.SendError(w, err)
		return
	}
	qPath := TObject.Data
	if len(qPath) == 0 {
		httpserver.SendError(w, errors.New("path is empty"))
		return
	}
	reqMethod := strings.ToLower(r.Method)
	if reqMethod != "post" && reqMethod != "put" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	{
		mr, err := r.MultipartReader()
		if err == nil {
			pName := r.Header.Get(headerFormNameFile)
			if len(pName) == 0 {
				pName = defaultFormNameFile
			}
			hasfile := false
			for {
				p, err := mr.NextPart()
				if nil == p || err == io.EOF {
					break
				}
				if p.FormName() != pName {
					continue
				}
				hasfile = true
				err = fsapi.fm.DoWrite(qPath, p)
				if nil != err {
					httpserver.SendError(w, err)
				} else {
					httpserver.SendSuccess(w, "")
				}
				break
			}
			if !hasfile {
				httpserver.SendError(w, errors.New("File not found from the form"))
				return
			}
		} else if nil != err && err == http.ErrNotMultipart {
			err := fsapi.fm.DoWrite(qPath, r.Body)
			if nil != err {
				httpserver.SendError(w, err)
			} else {
				httpserver.SendSuccess(w, "")
			}
		} else {
			httpserver.SendError(w, err)
		}
	}
}

// Download 下载
func (fsapi *FileAPI) download(w http.ResponseWriter, r *http.Request) {
	qToken := r.FormValue("token")
	TObject, err := fsapi.parseFileSteamToken(qToken)
	if nil != err || nil == TObject {
		httpserver.SendError(w, err)
		return
	}
	// 刷新token, 使其不过期 - 60s不访问则过期
	fsapi.fm.RefreshToken(qToken)
	if TObject.Type == "download" {
		fileName := strtool.GetPathName(TObject.Data)
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	fsapi.doSendStream(w, r, TObject.Data)
}

// doSendStream 发送数据流, 支持分段
func (fsapi *FileAPI) doSendStream(w http.ResponseWriter, r *http.Request, path string) {
	// 校验
	if !fsapi.fm.IsFile(path) {
		httpserver.SendError(w, ErrorFileNotExist)
		return
	}
	size, err := fsapi.fm.GetFileSize(path)
	if err != nil {
		httpserver.SendError(w, err)
		return
	}
	qRange := r.FormValue("Range")
	if len(qRange) == 0 {
		qRange = r.Header.Get("Range")
	}
	start := int64(0)
	end := int64(size)
	if len(qRange) > 0 {
		temp := qRange[strings.Index(qRange, "=")+1:]
		index := strings.Index(temp, "-")
		if index > -1 {
			start, err = strconv.ParseInt(temp[0:strings.Index(temp, "-")], 10, 64)
			if nil != err {
				start = 0
			}
			end, err = strconv.ParseInt(temp[strings.Index(temp, "-")+1:], 10, 64)
			if nil != err || end == 0 {
				end = size
			}
		}
	}
	fr, err := fsapi.fm.DoRead(path, start)
	defer func() {
		if nil != fr {
			fr.Close()
		}
	}()
	if nil != err {
		httpserver.SendError(w, err)
		return
	}
	{
		stransSize := end - start
		w.Header().Set("Content-Length", strconv.Itoa(int(stransSize)))
		if len(qRange) > 0 {
			w.Header().Set("Content-Range", "bytes "+strconv.Itoa(int(start))+"-"+strconv.Itoa(int(end-1))+"/"+strconv.Itoa(int(size)))
			w.WriteHeader(http.StatusPartialContent)
		}
		//
		for {
			if stransSize == 0 || stransSize < 0 {
				if end == size {
					// fm.RemoveToken(token)
				}
				break
			}
			buf := make([]byte, 4096)
			n, err := fr.Read(buf)
			if err != nil && err != io.EOF {
				httpserver.SendError(w, err)
				break
			}
			if 0 == n {
				break
			}
			if n > int(stransSize) {
				n = int(stransSize)
			}
			wn, err := w.Write(buf[:n])
			if nil != err {
				break
			}
			stransSize = stransSize - int64(wn)
		}
	}
}

// AsyncExec 发起一个异步操作, 返回一个可以查询的tooken
func (fsapi *FileAPI) AsyncExec(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("func")
	executor, err := fsapi.ast.GetTaskObject(name)
	if nil != err {
		httpserver.SendError(w, err)
		return
	}
	token, err := executor.Execute(r)
	if nil != err {
		httpserver.SendError(w, err)
		return
	}
	httpserver.SendSuccess(w, token)
}

// AsyncExecToken 查询由AsyncExec返回的token状态
func (fsapi *FileAPI) AsyncExecToken(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("func")
	executor, err := fsapi.ast.GetTaskObject(name)
	if nil != err {
		httpserver.SendError(w, err)
		return
	}
	executor.Status(w, r)
}
