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
	"gofs/comm/httpserver"
	"gofs/data/filemanage"
	"gofs/service/restful/signature"
	"gutils/hstool"
	"gutils/mloader"
	"gutils/strtool"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// FileAPI 文件API接口模块
type FileAPI struct {
	fm *filemanage.FileManager
	hs *httpserver.HTTPServer
	sg *signature.Signature
}

// ModuleOpts 模块加载器接口实现, 返回模块信息&配置Name:
func (fsapi *FileAPI) ModuleOpts() mloader.Opts {
	return mloader.Opts{
		Name:        "FileAPI",
		Version:     1.0,
		Description: "文件管理API接口模块",
		OnReady: func(mctx *mloader.Loader) {
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

	// 单个注册Handle
	err = fsapi.hs.AddHandlerFunc(baseurl+"/upload/:"+`[\S]+`, fsapi.Upload)
	if err != nil {
		panic(err)
	}
	err = fsapi.hs.AddHandlerFunc(baseurl+"/download/:"+`[\S]+`, fsapi.Download)
	if err != nil {
		panic(err)
	}
	err = fsapi.hs.AddHandlerFunc(baseurl+"/openfile/:"+`[\S]+`, fsapi.OpenFile)
	if err != nil {
		panic(err)
	}

	// 注册Api签名拦截器
	fsapi.hs.AddIgnoreFilter(baseurl + "/download/:" + `[\S]+`)
	fsapi.hs.AddIgnoreFilter(baseurl + "/openfile/:" + `[\S]+`)
	fsapi.hs.AddIgnoreFilter(baseurl + "/upload/:" + `[\S]+`)
	fsapi.hs.AddURLFilter(baseurl+"/:"+`[\S]+`, fsapi.sg.RestfulAPIFilter)

	fmt.Println("   > FileApiModule http registered end")
}

// RoutList 向 Server Router 中注册下列处理器, 实现接口 httpserver.Registrar
func (fsapi *FileAPI) RoutList() httpserver.StructRegistrar {
	return httpserver.StructRegistrar{
		IsToLower: true,
		BasePath:  baseurl,
		HandlerFunc: []hstool.HandlersFunc{
			fsapi.TransferToken,
			fsapi.BatchOperationTokenStauts,
			fsapi.List,
			fsapi.Del,
			fsapi.ReName,
			fsapi.NewFolder,
			fsapi.MoveAsync,
			fsapi.CopyAsync,
			fsapi.Info,
			// fsapi.Upload,
			// fsapi.Download,
		},
	}
}

// TransferToken 传输令牌申请
func (fsapi *FileAPI) TransferToken(w http.ResponseWriter, r *http.Request) {
	qpath := strtool.Parse2UnixPath(r.FormValue("path"))
	if len(qpath) == 0 {
		httpserver.SendError(w, ErrorOprationFailed)
		return
	}
	if strings.Index(qpath, "../") > -1 {
		httpserver.SendError(w, ErrorFileNotExist)
		return
	}
	// if !fsapi.fm.IsFile(qpath) {
	// 	httpserver.SendError(w, ErrorFileNotExist)
	// 	return
	// }
	httpserver.SendSuccess(w, fsapi.fm.AskToken(TokenType.transfer, &FileTransferToken{
		FilePath: qpath,
	}))
}

// BatchOperationTokenStauts Token信息查询 Get用于查询|Post用于操作(ErrorOperation)
func (fsapi *FileAPI) BatchOperationTokenStauts(w http.ResponseWriter, r *http.Request) {
	qToken := r.FormValue("token")
	fsapi.fm.RefreshToken(qToken)
	tokenBody, tokenErr := getFileBatchOperationTokenObject(fsapi, qToken)
	// fmt.Println("Token: ", r.Method, qToken, tokenBody)
	if nil == tokenBody || nil != tokenErr {
		httpserver.SendError(w, ErrorOprationExpires)
		return
	}
	// Get 用于获取令牌信息
	if r.Method == "GET" {
		bt, _ := json.Marshal(tokenBody)
		httpserver.SendSuccess(w, string(bt))

		// POST 用于操作|中断
	} else if r.Method == "POST" {
		qOperation := r.FormValue("operation")
		if len(qOperation) == 0 {
			httpserver.SendError(w, ErrorOprationFailed)
		} else {
			switch qOperation {
			// 忽略单个 错误
			case ErrorOperation.ignore:
				tokenBody.ErrorString = ""
				tokenBody.IsIgnore = true
				break
			// 为后续的 错误 执行忽略
			case ErrorOperation.ignoreall:
				tokenBody.ErrorString = ""
				tokenBody.IsIgnoreAll = true
				break
			// 覆盖单个 已存在 错误
			case ErrorOperation.replace:
				tokenBody.ErrorString = ""
				tokenBody.IsReplace = true
				break
			// 每次都覆盖 已存在 错误
			case ErrorOperation.replaceall:
				tokenBody.ErrorString = ""
				tokenBody.IsReplaceAll = true
				break
			// 立即中断操作
			case ErrorOperation.discontinue:
				tokenBody.ErrorString = ""
				tokenBody.IsComplete = true
				tokenBody.IsDiscontinue = true
				// fsapi.fm.RemoveToken(qToken)
				break
			default:
				httpserver.SendError(w, ErrorOprationFailed)
				return
			}
			httpserver.SendSuccess(w, "")
		}
	}
}

// List 查询路径下的列表以及基本信息
func (fsapi *FileAPI) List(w http.ResponseWriter, r *http.Request) {
	qpath := r.FormValue("path")
	qSort := r.FormValue("sort")
	qAsc := r.FormValue("asc")
	// fmt.Println(r.URL, r.Body, qpath)
	if len(qpath) == 0 {
		httpserver.SendError(w, errors.New("'Path' parameter not found"))
		return
	}
	if len(qAsc) == 0 {
		qAsc = "true"
	}
	res, err := fsapi.fm.GetDirListInfo(qpath)
	if err != nil {
		httpserver.SendError(w, err)
		return
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
		httpserver.SendError(w, errors.New("'Path' parameter not found"))
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

// MoveAsync 移动文件|文件夹 - 异步处理, 返回Token用于查询进度
func (fsapi *FileAPI) MoveAsync(w http.ResponseWriter, r *http.Request) {
	qSrcPath := r.FormValue("srcPath")
	qDstPath := r.FormValue("dstPath")
	qReplace := strtool.String2Bool(r.FormValue("replace"))
	qIgnore := strtool.String2Bool(r.FormValue("ignore"))
	if len(qSrcPath) == 0 {
		httpserver.SendError(w, errors.New("'srcPath' parameter not found"))
		return
	}
	if len(qDstPath) == 0 {
		httpserver.SendError(w, errors.New("'dstPath' parameter not found"))
		return
	}
	// 异步处理, 返回一个Token用于查询进度
	moveTokenObject := &FileBatchOperationTokenObject{
		CountIndex:   1,
		ErrorString:  "",
		Src:          qSrcPath,
		Dst:          qDstPath,
		IsSrcExist:   true,
		IsDstExist:   false,
		IsReplace:    false,
		IsReplaceAll: qReplace,
		IsIgnore:     false,
		IsIgnoreAll:  qIgnore,
	}
	token := fsapi.fm.AskToken(TokenType.moveDir, moveTokenObject)
	go func(token string) {
		moveDirErr := fsapi.fm.DoMove(qSrcPath, qDstPath, qReplace, qIgnore, func(s_src, s_dst string, moveErr *filemanage.MoveError) error {
			// 获取令牌数据, 不存在则说明已经销毁
			// 并保持刷新token的有效性, 除非终止操作否则都继续
			tokenBody, tokenErr := getFileBatchOperationTokenObject(fsapi, token)
			if nil != tokenErr {
				return tokenErr
			}
			if tokenBody.IsDiscontinue {
				return ErrorDiscontinue
			}
			tokenBody.CountIndex = int64(1)
			tokenBody.IsSrcExist = false
			tokenBody.IsDstExist = false
			tokenBody.ErrorString = ""
			tokenBody.Src = s_src
			tokenBody.Dst = s_dst

			// 如果遇到错误了
			if nil != moveErr {
				// 判断是否是目标位置已经存在的错误, 如果是的话需要选择是否覆盖他
				if moveErr.DstIsExist {
					// 查找之前是否设置了 替换全部错误
					if tokenBody.IsReplaceAll {
						// 先删除然后再替换, 如果覆盖操作没有出现问题
						reMoveErr := fsapi.fm.DoMove(s_src, s_dst, true, false, func(s, d string, mErr *filemanage.MoveError) error {
							if nil != mErr {
								return errors.New(mErr.ErrorString)
							}
							return nil
						})
						if nil == reMoveErr {
							return nil
						}
						tokenBody.ErrorString = reMoveErr.Error()
					}
					// 如果设置了自动覆盖, 但是任然出错, 则判断是否忽略错误选项
					if tokenBody.IsIgnoreAll {
						tokenBody.ErrorString = ""
						return nil // 跳过这个文件
					}
				} else {
					// 不是路径重复类错误
					// 如果是其他错误就不管了, 暂时无法处理只能选择 忽略|暂停
					// 查找之前是否设置了 忽略全部错误
					if tokenBody.IsIgnoreAll {
						return nil // 跳过这个文件
					}
				}

				// 到此说明 没有设置自动覆盖和自动忽略
				tokenBody.IsSrcExist = moveErr.SrcIsExist
				tokenBody.IsDstExist = moveErr.DstIsExist
				if len(tokenBody.ErrorString) == 0 {
					tokenBody.ErrorString = moveErr.ErrorString // 设置错误, 等待客户端获取, 等待操作
				}
				for {
					tokenBody, tokenErr := getFileBatchOperationTokenObject(fsapi, token)
					if nil != tokenErr {
						return tokenErr
					}
					if tokenBody.IsDiscontinue {
						return ErrorDiscontinue
					}
					// fmt.Println("for: ", tokenBody)
					// 选择了忽略|忽略全部
					if tokenBody.IsIgnore || tokenBody.IsIgnoreAll {
						if tokenBody.IsIgnore {
							tokenBody.IsIgnore = false // 一次性的
						}
						return nil
					}
					// 选择了覆盖|覆盖全部
					if tokenBody.IsReplace || tokenBody.IsReplaceAll {
						if tokenBody.IsReplace {
							tokenBody.IsReplace = false // 一次性的
						}
						if moveErr.SrcIsExist {
							moveCopyErr := fsapi.fm.DoMove(s_src, s_dst, true, false, func(s, d string, mErr *filemanage.MoveError) error {
								if nil != mErr {
									return errors.New(mErr.ErrorString)
								}
								return nil
							})
							if nil != moveCopyErr {
								tokenBody.ErrorString = moveCopyErr.Error()
							} else {
								return nil
							}
						}
					}
					time.Sleep(time.Duration(100) * time.Millisecond) // 休眠100ms
				}
			}
			return nil
		})

		// 到这里如果没有错误就是成功了
		tokenBody := fsapi.fm.GetToken(token)
		if nil != tokenBody {
			tokenBody := tokenBody.(*FileBatchOperationTokenObject)
			if nil != moveDirErr {
				tokenBody.ErrorString = moveDirErr.Error()
			} else {
				tokenBody.ErrorString = ""
			}
			tokenBody.IsComplete = true
			tokenBody.IsDiscontinue = ErrorDiscontinue.Error() == tokenBody.ErrorString
			// fmt.Println("copyDirErr: ", copyDirErr)
		}
	}(token)
	httpserver.SendSuccess(w, token)
}

// CopyAsync 拷贝文件|文件夹 - 异步操作, 返回Token用于查询进度
func (fsapi *FileAPI) CopyAsync(w http.ResponseWriter, r *http.Request) {
	qSrcPath := r.FormValue("srcPath")
	qDstPath := r.FormValue("dstPath")
	qReplace := strtool.String2Bool(r.FormValue("replace"))
	qIgnore := strtool.String2Bool(r.FormValue("ignore"))
	if len(qSrcPath) == 0 {
		httpserver.SendError(w, errors.New("'qSrcPath' parameter not found"))
		return
	}
	if len(qDstPath) == 0 {
		httpserver.SendError(w, errors.New("'qDstPath' parameter not found"))
		return
	}
	// 异步处理, 返回一个Token用于查询进度
	copyTokenObject := &FileBatchOperationTokenObject{
		ErrorString:  "",
		Src:          qSrcPath,
		Dst:          qDstPath,
		IsSrcExist:   true,
		IsDstExist:   false,
		IsReplace:    false,
		IsReplaceAll: qReplace,
		IsIgnore:     false,
		IsIgnoreAll:  qIgnore,
	}
	token := fsapi.fm.AskToken(TokenType.copyDir, copyTokenObject)
	go func(token string) {
		// 这里面已经不属于一个会话, 使用令牌保存数据
		copyDirErr := fsapi.fm.DoCopy(qSrcPath, qDstPath, qReplace, qIgnore, func(s_src, s_dst string, copyErr *filemanage.CopyError) error {
			// 获取令牌数据, 不存在则说明已经销毁
			// 并保持刷新token的有效性, 除非终止操作否则都继续
			tokenBody, tokenErr := getFileBatchOperationTokenObject(fsapi, token)
			if nil != tokenErr {
				return tokenErr
			}
			if tokenBody.IsDiscontinue {
				return ErrorDiscontinue
			}
			tokenBody.CountIndex = tokenBody.CountIndex + 1
			tokenBody.IsSrcExist = false
			tokenBody.IsDstExist = false
			tokenBody.ErrorString = ""
			tokenBody.Src = s_src
			tokenBody.Dst = s_dst
			// 如果遇到错误了
			if nil != copyErr {
				// 判断是否是目标位置已经存在的错误, 如果是的话需要选择是否覆盖他
				if copyErr.DstIsExist {

					// 查找之前是否设置了 替换全部错误
					if tokenBody.IsReplaceAll {
						// 先删除然后再替换, 如果覆盖操作没有出现问题
						reCopyErr := fsapi.fm.DoCopy(s_src, s_dst, true, false, func(s, d string, cErr *filemanage.CopyError) error {
							if nil != cErr {
								return errors.New(cErr.ErrorString)
							}
							return nil
						})
						if nil == reCopyErr {
							return nil
						}
						tokenBody.ErrorString = reCopyErr.Error()
					}
					// 如果设置了自动覆盖, 但是任然出错, 则判断是否忽略错误选项
					if tokenBody.IsIgnoreAll {
						tokenBody.ErrorString = ""
						return nil // 跳过这个文件
					}
				} else {
					// 不是路径重复类错误
					// 如果是其他错误就不管了, 暂时无法处理只能选择 忽略|暂停
					// 查找之前是否设置了 忽略全部错误
					if tokenBody.IsIgnoreAll {
						return nil // 跳过这个文件
					}
				}

				// 到此说明 没有设置自动覆盖和自动忽略
				tokenBody.IsSrcExist = copyErr.SrcIsExist
				tokenBody.IsDstExist = copyErr.DstIsExist
				if len(tokenBody.ErrorString) == 0 {
					tokenBody.ErrorString = copyErr.ErrorString // 设置错误, 等待客户端获取, 等待操作
				}
				for {
					tokenBody, tokenErr := getFileBatchOperationTokenObject(fsapi, token)
					if nil != tokenErr {
						return tokenErr
					}
					if tokenBody.IsDiscontinue {
						return ErrorDiscontinue
					}
					// fmt.Println("for: ", tokenBody)
					// 选择了忽略|忽略全部
					if tokenBody.IsIgnore || tokenBody.IsIgnoreAll {
						if tokenBody.IsIgnore {
							tokenBody.IsIgnore = false // 一次性的
						}
						return nil
					}
					// 选择了覆盖|覆盖全部
					if tokenBody.IsReplace || tokenBody.IsReplaceAll {
						if tokenBody.IsReplace {
							tokenBody.IsReplace = false // 一次性的
						}
						if copyErr.SrcIsExist {
							reCopyErr := fsapi.fm.DoCopy(s_src, s_dst, true, false, func(s, d string, cErr *filemanage.CopyError) error {
								if nil != cErr {
									return errors.New(cErr.ErrorString)
								}
								return nil
							})
							if nil != reCopyErr {
								tokenBody.ErrorString = reCopyErr.Error()
							} else {
								return nil
							}
						}
					}
					time.Sleep(time.Duration(100) * time.Millisecond) // 休眠100ms
				}
			}
			return nil
		})
		// 到这里如果没有错误就是成功了
		tokenBody := fsapi.fm.GetToken(token)
		if nil != tokenBody {
			tokenBody := tokenBody.(*FileBatchOperationTokenObject)
			if nil != copyDirErr {
				tokenBody.ErrorString = copyDirErr.Error()
			} else {
				tokenBody.ErrorString = ""
			}
			tokenBody.IsComplete = true
			tokenBody.IsDiscontinue = ErrorDiscontinue.Error() == tokenBody.ErrorString
			// fmt.Println("copyDirErr: ", copyDirErr)
		}
	}(token)
	httpserver.SendSuccess(w, token)
}

// ReName 重命名
func (fsapi *FileAPI) ReName(w http.ResponseWriter, r *http.Request) {
	qSrcPath := r.FormValue("path")
	qName := r.FormValue("name")
	if !fsapi.fm.IsExist(qSrcPath) {
		httpserver.SendError(w, ErrorFileNotExist)
		return
	}
	if len(qName) == 0 {
		httpserver.SendError(w, ErrorNewNameIsEmpty)
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
	qSrcPath := r.FormValue("path")
	qSrcPath = strtool.Parse2UnixPath(qSrcPath)
	if !fsapi.fm.IsExist(strtool.GetPathParent(qSrcPath)) {
		httpserver.SendError(w, ErrorParentFolderNotExist)
		return
	}
	rnmErr := fsapi.fm.DoNewFolder(qSrcPath)
	if nil == rnmErr {
		httpserver.SendSuccess(w, "")
	} else {
		httpserver.SendError(w, rnmErr)
	}
}

// Info Info
func (fsapi *FileAPI) Info(w http.ResponseWriter, r *http.Request) {

}

// NameSearch NameSearch
func (fsapi *FileAPI) NameSearch(w http.ResponseWriter, r *http.Request) {

}

// Upload 文件上传, 支持Form和Body上传方式
// 参数: Header("Save-Path", ["Formname-File"])
func (fsapi *FileAPI) Upload(w http.ResponseWriter, r *http.Request) {
	token := strtool.GetPathName(r.URL.Path)
	if len(token) > 0 {
		tkobj, _ := getFileTransferTokenObject(fsapi, token)
		if tkobj != nil && len(tkobj.FilePath) > 0 {
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
					err = fsapi.fm.DoWrite(tkobj.FilePath, p)
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
				err := fsapi.fm.DoWrite(tkobj.FilePath, r.Body)
				if nil != err {
					httpserver.SendError(w, err)
				} else {
					httpserver.SendSuccess(w, "")
				}
			} else {
				httpserver.SendError(w, err)
			}
		}
		// 使token无效
		fsapi.fm.RemoveToken(token)
	}
	httpserver.SendError(w, ErrorOprationExpires)
}

// OpenFile 打开
func (fsapi *FileAPI) OpenFile(w http.ResponseWriter, r *http.Request) {
	token := strtool.GetPathName(r.URL.Path)
	index := strings.Index(token, ".")
	if index > -1 {
		token = token[:index]
	}
	tokenObject, err := getFileTransferTokenObject(fsapi, token)
	if nil != err || nil == tokenObject {
		httpserver.SendError(w, ErrorOprationExpires)
		return
	}
	doSendStream(fsapi.fm, token, w, r)
}

// Download 下载
func (fsapi *FileAPI) Download(w http.ResponseWriter, r *http.Request) {
	token := strtool.GetPathName(r.URL.Path)
	tokenObject, err := getFileTransferTokenObject(fsapi, token)
	if nil != err || nil == tokenObject {
		httpserver.SendError(w, ErrorOprationExpires)
		return
	}
	fileName := strtool.GetPathName(tokenObject.FilePath)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	doSendStream(fsapi.fm, token, w, r)
}

// getFileBatchOperationTokenObject 获取批文件量操作Token对象
func getFileBatchOperationTokenObject(fsapi *FileAPI, token string) (*FileBatchOperationTokenObject, error) {
	tokenBody := fsapi.fm.GetToken(token)
	// 并保持刷新token的有效性, 除非终止操作否则都继续
	if nil == tokenBody {
		return nil, ErrorOprationExpires
	}
	// fmt.Println("tokenBody: ", tokenBody)
	return tokenBody.(*FileBatchOperationTokenObject), nil
}

// getFileTransferTokenObject 获取文件传输Token对象
func getFileTransferTokenObject(fsapi *FileAPI, token string) (*FileTransferToken, error) {
	tokenBody := fsapi.fm.GetToken(token)
	// 并保持刷新token的有效性, 除非终止操作否则都继续
	if nil == tokenBody {
		return nil, ErrorOprationExpires
	}
	ftt, ok := tokenBody.(*FileTransferToken)
	if ok {
		return ftt, nil
	}
	return nil, ErrorOprationExpires
}

// doSendStream 发送数据流, 支持分段
func doSendStream(fm *filemanage.FileManager, token string, w http.ResponseWriter, r *http.Request) {
	tBody := fm.GetToken(token)
	if nil == tBody {
		httpserver.SendError(w, ErrorOprationExpires)
		return
	}
	FTTObject, ok := tBody.(*FileTransferToken)
	if !ok {
		httpserver.SendError(w, ErrorOprationExpires)
		return
	}
	// 刷新token, 使其不过期
	fm.RefreshToken(token)
	// 校验
	if !fm.IsFile(FTTObject.FilePath) {
		httpserver.SendError(w, ErrorFileNotExist)
		return
	}
	size, err := fm.GetFileSize(FTTObject.FilePath)
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
	fr, err := fm.DoRead(FTTObject.FilePath, start)
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
