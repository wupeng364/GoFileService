// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 文件预览接口模块

package preview

import (
	"encoding/json"
	"fmt"
	"gofs/base/httpserver"
	"gofs/base/signature"
	"gofs/data/filemanage"
	"gutils/hstool"
	"gutils/mloader"
	"gutils/strtool"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Preview 文件预览模块
type Preview struct {
	fm *filemanage.FileManager
	hs *httpserver.HTTPServer
	sg *signature.Signature
}

// ModuleOpts 模块加载器接口实现, 返回模块信息&配置Name:
func (preview *Preview) ModuleOpts() mloader.Opts {
	return mloader.Opts{
		Name:        "Preview",
		Version:     1.0,
		Description: "文件预览模块",
		OnReady: func(mctx *mloader.Loader) {
			preview.fm = mctx.GetModuleByTemplate(preview.fm).(*filemanage.FileManager)
			preview.hs = mctx.GetModuleByTemplate(preview.hs).(*httpserver.HTTPServer)
			preview.sg = mctx.GetModuleByTemplate(preview.sg).(*signature.Signature)
		},
		OnInit: preview.init,
	}
}

// 向 HttpServerModule 中注册服务地址
func (preview *Preview) init() {
	// 批量注册Handle
	err := preview.hs.AddRegistrar(preview)
	if err != nil {
		panic(err)
	}

	// 注册Api签名拦截器
	preview.hs.AddURLFilter(baseurl+"/:"+`[\S]+`, preview.sg.RestfulAPIFilter)

	fmt.Println("   > PreviewModule http registered end")
}

// RoutList 向 Server Router 中注册下列处理器, 实现接口 httpserver.Registrar
func (preview *Preview) RoutList() httpserver.StructRegistrar {
	return httpserver.StructRegistrar{
		IsToLower: true,
		BasePath:  baseurl,
		HandlerFunc: []hstool.HandlersFunc{
			preview.Status,
			preview.Asktoken,
			preview.TokenDatas,
		},
	}
}

// Status 保持会话用
func (preview *Preview) Status(w http.ResponseWriter, r *http.Request) {
	httpserver.SendSuccess(w, "")
}

// Asktoken 预览令牌申请
func (preview *Preview) Asktoken(w http.ResponseWriter, r *http.Request) {
	qpath := r.FormValue("path")
	if len(qpath) == 0 {
		httpserver.SendError(w, ErrorOprationFailed)
		return
	}
	if !preview.fm.IsFile(qpath) {
		httpserver.SendError(w, ErrorFileNotExist)
		return
	}
	token := strtool.GetUUID()
	err := preview.sg.SetSessionAttr4Request(r, token, qpath)
	if nil != err {
		httpserver.SendError(w, err)
		return
	}
	httpserver.SendSuccess(w, token)
}

// TokenDatas Token信息查询
func (preview *Preview) TokenDatas(w http.ResponseWriter, r *http.Request) {
	qToken := r.FormValue("token")
	qType := r.FormValue("type")
	tData, err := preview.sg.GetSessionAttr4Request(r, qToken)
	if nil != err {
		httpserver.SendError(w, err)
		return
	}
	if len(tData) == 0 {
		httpserver.SendError(w, ErrorFileNotExist)
		return
	}
	switch qType {
	case "list":
		{
			prentPath := tData
			if !preview.fm.IsExist(prentPath) || preview.fm.IsFile(prentPath) {
				prentPath = strtool.GetPathParent(prentPath)
			}
			if !preview.fm.IsExist(prentPath) {
				httpserver.SendError(w, ErrorFileNotExist)
				return
			}
			fList, err := preview.fm.GetDirListInfo(prentPath)
			if err != nil {
				httpserver.SendError(w, err)
				return
			}
			json, err := json.Marshal(PInfo{
				Path: tData,
				PeerDatas: filemanage.FileListSorter{
					Asc:       true,
					SortField: "Path",
				}.Sort(fList),
			})
			if err != nil {
				httpserver.SendError(w, err)
				return
			}
			httpserver.SendSuccess(w, string(json))
		}
		break
	case "stream":
		{
			qPath := r.FormValue("path")
			if len(qPath) == 0 {
				qPath = tData
			}
			preview.doSendStream(w, r, qPath)
		}
		break
	default:
		httpserver.SendError(w, ErrorOprationFailed)
	}

}

// doSendStream 发送数据流, 支持分段
func (preview *Preview) doSendStream(w http.ResponseWriter, r *http.Request, path string) {
	// 校验
	if !preview.fm.IsFile(path) {
		httpserver.SendError(w, ErrorFileNotExist)
		return
	}
	size, err := preview.fm.GetFileSize(path)
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
	fr, err := preview.fm.DoRead(path, start)
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
