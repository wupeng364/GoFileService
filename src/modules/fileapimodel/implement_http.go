package fileapimodel
/**
 *@description 文件API接口模块-Http实现
 *@author	wupeng364@outlook.com
*/
import(
	"encoding/json"
	"net/http"
	"errors"
	"fmt"
	hst "common/httpservertools"
	"common/stringtools"
	"modules/filemanage"
	"modules/httpservermodel"
	"time"
	"io"
)
const(
	base_url	  = "/fsapi"
	cfg_http 	  = "http"
	cfg_http_addr = "addr"
	cfg_http_port = "port"
	// 
	header_formname_file  = "Formname-File"			// 用头信息标记Form表单中文件的FormName
	default_formname_file = "file"					// 默认使用这个作为Form表单中文件的FormName
	default_formname_fspath = "Save-Path"			// 默认使用这个作为Form表单中文件保存位置的FormName
	//default_formname_fsname = "filename"			// 默认使用这个作为Form表单中文件名字的FormName
)
/**
 * 文件基础操作网络接口
 */
type imp_http struct{
	fm *filemanage.FileManageModel
	hs *httpservermodel.HttpServerModel
}

// 向 HttpServerModel 中注册服务地址
func (fa imp_http) init(fm *filemanage.FileManageModel, hs *httpservermodel.HttpServerModel){
	fa.fm = fm; fa.hs = hs
	// 批量注册
	res := fa.hs.AddRegistrar(fa)
	if res != nil {
		panic("imp_http AddRegistrar failed")
	}
	// 单个注册
	res = fa.hs.AddHandlerFunc(base_url+"/upload/:"+`[\S]+`, fa.Upload)
	if res != nil {
		panic("imp_http AddHandlerFunc - upload failed")
	}
	res = fa.hs.AddHandlerFunc(base_url+"/download/:"+`[\S]+`, fa.Download)
	if res != nil {
		panic("imp_http AddHandlerFunc - download failed")
	}
	res = fa.hs.AddHandlerFunc(base_url+"/openfile/:"+`[\S]+`, fa.OpenFile)
	if res != nil {
		panic("imp_http AddHandlerFunc - openfile failed")
	}
	fmt.Println("   > FileApiModel http registered end")
}

// 向 Server Router 中注册下列处理器 , 实现接口 httpservertools.Registrar
func (fa imp_http) RoutList( ) hst.StructRegistrar{
	return hst.StructRegistrar{
		true,
		base_url,
		[]hst.HandlersFunc{
			fa.TransferToken,
			fa.BatchOperationTokenStauts,
			fa.List,
			fa.Del,
			fa.DelVer,
			fa.ReName,
			fa.NewFolder,
			fa.MoveAsync,
			fa.CopyAsync,
			fa.Info,
			// fa.Upload,
			// fa.Download,
		},
	}
}
// 传输令牌申请
func (fa imp_http) TransferToken( w http.ResponseWriter, r *http.Request ){
	q_path := r.FormValue("path")
	if len( q_path ) == 0 {
		sendError(w, Error_OprationFailed); return
	}
	if !fa.fm.IsFile(q_path) {
		sendError(w, Error_FileNotExist); return
	}
	sendSuccess(w, fa.fm.AskToken(TokenType.download, &FileTransferToken{
		FilePath: q_path,
	}))
}
// Token信息查询 Get用于查询|Post用于操作(ErrorOperation)
func (fa imp_http) BatchOperationTokenStauts( w http.ResponseWriter, r *http.Request ){
	q_token := r.FormValue("token")
	fa.fm.RefreshToken(q_token)
	tokenErr, tokenBody := getFileBatchOperationTokenObject(fa, q_token)
	// fmt.Println("Token: ", r.Method, q_token, tokenBody)
	if nil == tokenBody || nil != tokenErr {
		sendError(w, Error_OprationExpires); return
	}
	// Get 用于获取令牌信息
	if r.Method == "GET" {
		bt, _ := json.Marshal(tokenBody)
		sendSuccess(w, string(bt))
		
		// POST 用于操作|中断 
	}else if r.Method == "POST" {
		q_operation := r.FormValue("operation")
		if len( q_operation ) == 0 {
			sendError(w, Error_OprationFailed)
		}else {
			switch q_operation {
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
					// fa.fm.RemoveToken(q_token)
				break
				default:
					sendError(w, Error_OprationFailed)
					return
			}
			sendSuccess(w, "")
		}
	}
}
// 查询路径下的列表以及基本信息
func (fa imp_http) List( w http.ResponseWriter, r *http.Request ){
	q_path := r.FormValue("path")
	// fmt.Println(r.URL, r.Body, q_path)
	if len(q_path) == 0 {
		sendError(w, errors.New("'Path' parameter not found"))
		return
	}
	res, err := fa.fm.GetDirListInfo(q_path)
	if err != nil {
		sendError(w, err ) 
		return
	}
	json, err := json.Marshal(res)
	if err != nil {
		sendError(w, err )
	}
	w.Write([]byte(json))
}
// 批量|单个删除文件|文件夹
func (fa imp_http) Del( w http.ResponseWriter, r *http.Request ){
	q_path := r.FormValue("path")
	fmt.Println(r.URL, r.Body, q_path)
	if len(q_path) == 0 {
		sendError(w, errors.New("'Path' parameter not found")); return
	}
	if !fa.fm.IsExist(q_path) {
		sendError(w, Error_FileNotExist); return
	}
	del_err := fa.fm.DoDelete(q_path)
	if nil == del_err {
		sendSuccess(w, "")
	}else{
		sendError(w, del_err)
	}
}
// 移动文件|文件夹 - 异步处理, 返回Token用于查询进度
func (fa imp_http) MoveAsync(w http.ResponseWriter, r *http.Request){
	q_src_path := r.FormValue("srcPath")
	q_dst_path := r.FormValue("dstPath")
	q_replace  := stringtools.ParseBool( r.FormValue("replace") )
	q_ignore   := stringtools.ParseBool( r.FormValue("ignore") )
	if len(q_src_path) == 0 {
		sendError(w, errors.New("'q_src_path' parameter not found"))
		return
	}
	if len(q_dst_path) == 0 {
		sendError(w, errors.New("'q_dst_path' parameter not found"))
		return
	}
	// 异步处理, 返回一个Token用于查询进度
	moveTokenObject := FileBatchOperationTokenObject {
		CountIndex: 1,
		ErrorString:	"",
		Src:			q_src_path,
		Dst:			q_dst_path,
		IsSrcExist:		true,
		IsDstExist:		false,
		IsReplace:		false,
		IsReplaceAll:	q_replace,
		IsIgnore:		false,
		IsIgnoreAll:	q_ignore,
	}
	token := fa.fm.AskToken(TokenType.moveDir, &moveTokenObject)
	go func ( token string) {
		moveDirErr := fa.fm.DoMove(q_src_path, q_dst_path, q_replace, q_ignore, func(count int, s_src, s_dst string, moveErr *filemanage.MoveError)error{
			// 获取令牌数据, 不存在则说明已经销毁
			// 并保持刷新token的有效性, 除非终止操作否则都继续
			tokenErr, tokenBody := getFileBatchOperationTokenObject(fa, token)
			if nil != tokenErr {
				return tokenErr
			}
			if tokenBody.IsDiscontinue {
				return Error_Discontinue
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
				if moveErr.IsDstExist {
					// 查找之前是否设置了 替换全部错误 
					if tokenBody.IsReplaceAll {
						// 先删除然后再替换, 如果覆盖操作没有出现问题
						reMoveErr := fa.fm.DoMove(s_src, s_dst, true, false, func(c int, s, d string, mErr *filemanage.MoveError)error{
							if nil != mErr {
								return errors.New(mErr.ErrorString)
							}
							return nil
						})
						if nil == reMoveErr {
							return nil
						}else{
							tokenBody.ErrorString = reMoveErr.Error( )
						}
					}
					// 如果设置了自动覆盖, 但是任然出错, 则判断是否忽略错误选项
					if tokenBody.IsIgnoreAll {
						tokenBody.ErrorString = ""
						return nil // 跳过这个文件
					}
				}else{
					// 不是路径重复类错误
					// 如果是其他错误就不管了, 暂时无法处理只能选择 忽略|暂停
					// 查找之前是否设置了 忽略全部错误
					if tokenBody.IsIgnoreAll {
						return nil // 跳过这个文件
					}
				}
				
				// 到此说明 没有设置自动覆盖和自动忽略
				tokenBody.IsSrcExist  = moveErr.IsSrcExist
				tokenBody.IsDstExist  = moveErr.IsDstExist
				if len(tokenBody.ErrorString) == 0 {
					tokenBody.ErrorString = moveErr.ErrorString // 设置错误, 等待客户端获取, 等待操作
				}
				for{
					tokenErr, tokenBody := getFileBatchOperationTokenObject(fa, token)
					if nil != tokenErr {
						return tokenErr
					}
					if tokenBody.IsDiscontinue {
						return Error_Discontinue
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
						if moveErr.IsSrcExist {
							moveCopyErr := fa.fm.DoMove(s_src, s_dst, true, false, func(c int, s, d string, mErr *filemanage.MoveError)error{ 
								if nil != mErr {
									return errors.New(mErr.ErrorString)
								}
								return nil
							})
							if nil != moveCopyErr {
								tokenBody.ErrorString = moveCopyErr.Error( )
							}else{
								return nil
							}							
						}
					}
					time.Sleep(time.Duration(100)*time.Millisecond)	// 休眠100ms
				}
			}
			return nil
		})
		
		// 到这里如果没有错误就是成功了
		tokenObject := fa.fm.GetToken(token)
		if nil != tokenObject && nil != tokenObject.TokenBody {
			tokenBody := tokenObject.TokenBody.(*FileBatchOperationTokenObject)
			if nil != moveDirErr {
				tokenBody.ErrorString = moveDirErr.Error( )
			}else{
				tokenBody.ErrorString = ""
			}
			tokenBody.IsComplete = true
			tokenBody.IsDiscontinue = Error_Discontinue.Error() == tokenBody.ErrorString
			// fmt.Println("copyDirErr: ", copyDirErr)
		}
	}( token )
	sendSuccess(w, token)
}
// 拷贝文件|文件夹 - 异步操作, 返回Token用于查询进度
func (fa imp_http) CopyAsync( w http.ResponseWriter, r *http.Request ){
	q_src_path := r.FormValue("srcPath")
	q_dst_path := r.FormValue("dstPath")
	q_replace  := stringtools.ParseBool( r.FormValue("replace") )
	q_ignore   := stringtools.ParseBool( r.FormValue("ignore") )
	if len(q_src_path) == 0 {
		sendError(w, errors.New("'q_src_path' parameter not found"))
		return
	}
	if len(q_dst_path) == 0 {
		sendError(w, errors.New("'q_dst_path' parameter not found"))
		return
	}
	// 异步处理, 返回一个Token用于查询进度
	copyTokenObject := FileBatchOperationTokenObject {
		ErrorString:	"",
		Src:			q_src_path,
		Dst:			q_dst_path,
		IsSrcExist:		true,
		IsDstExist:		false,
		IsReplace:		false,
		IsReplaceAll:	q_replace,
		IsIgnore:		false,
		IsIgnoreAll:	q_ignore,
	}
	token := fa.fm.AskToken(TokenType.copyDir, &copyTokenObject)
	go func ( token string ) {
		// 这里面已经不属于一个会话, 使用令牌保存数据
		copyDirErr := fa.fm.DoCopy(q_src_path, q_dst_path, q_replace, q_ignore, func(count int, s_src, s_dst string, copyErr *filemanage.CopyError)error{
			// 获取令牌数据, 不存在则说明已经销毁
			// 并保持刷新token的有效性, 除非终止操作否则都继续
			tokenErr, tokenBody := getFileBatchOperationTokenObject(fa, token)
			if nil != tokenErr {
				return tokenErr
			}
			if tokenBody.IsDiscontinue {
				return Error_Discontinue
			}
			tokenBody.CountIndex = int64(count)
			tokenBody.IsSrcExist = false
			tokenBody.IsDstExist = false
			tokenBody.ErrorString = ""
			tokenBody.Src = s_src
			tokenBody.Dst = s_dst
			// 如果遇到错误了
			if nil != copyErr {
				// 判断是否是目标位置已经存在的错误, 如果是的话需要选择是否覆盖他
				if copyErr.IsDstExist {
					
					// 查找之前是否设置了 替换全部错误 
					if tokenBody.IsReplaceAll {
						// 先删除然后再替换, 如果覆盖操作没有出现问题
						reCopyErr := fa.fm.DoCopy(s_src, s_dst, true, false, func(c int, s, d string, cErr *filemanage.CopyError)error{
							if nil != cErr {
								return errors.New(cErr.ErrorString)
							}
							return nil
						})
						if nil == reCopyErr {
							return nil
						}else{
							tokenBody.ErrorString = reCopyErr.Error( )
						}
					}
					// 如果设置了自动覆盖, 但是任然出错, 则判断是否忽略错误选项
					if tokenBody.IsIgnoreAll {
						tokenBody.ErrorString = ""
						return nil // 跳过这个文件
					}
				}else{
					// 不是路径重复类错误
					// 如果是其他错误就不管了, 暂时无法处理只能选择 忽略|暂停
					// 查找之前是否设置了 忽略全部错误
					if tokenBody.IsIgnoreAll {
						return nil // 跳过这个文件
					}
				}
				
				// 到此说明 没有设置自动覆盖和自动忽略
				tokenBody.IsSrcExist  = copyErr.IsSrcExist
				tokenBody.IsDstExist  = copyErr.IsDstExist
				if len(tokenBody.ErrorString) == 0 {
					tokenBody.ErrorString = copyErr.ErrorString // 设置错误, 等待客户端获取, 等待操作
				}
				for{
					tokenErr, tokenBody := getFileBatchOperationTokenObject(fa, token)
					if nil != tokenErr {
						return tokenErr
					}
					if tokenBody.IsDiscontinue {
						return Error_Discontinue
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
						if copyErr.IsSrcExist {
							reCopyErr := fa.fm.DoCopy(s_src, s_dst, true, false, func(c int, s, d string, cErr *filemanage.CopyError)error{ 
								if nil != cErr {
									return errors.New(cErr.ErrorString)
								}
								return nil
							})
							if nil != reCopyErr {
								tokenBody.ErrorString = reCopyErr.Error( )
							}else{
								return nil
							}							
						}
					}
					time.Sleep(time.Duration(100)*time.Millisecond)	// 休眠100ms
				}
			}
			return nil
		})
		// 到这里如果没有错误就是成功了
		tokenObject := fa.fm.GetToken(token)
		if nil != tokenObject && nil != tokenObject.TokenBody {
			tokenBody := tokenObject.TokenBody.(*FileBatchOperationTokenObject)
			if nil != copyDirErr {
				tokenBody.ErrorString = copyDirErr.Error( )
			}else{
				tokenBody.ErrorString = ""
			}
			tokenBody.IsComplete = true
			tokenBody.IsDiscontinue = Error_Discontinue.Error() == tokenBody.ErrorString
			// fmt.Println("copyDirErr: ", copyDirErr)
		}
	}( token )
	sendSuccess(w, token)
}
func (fa imp_http) DelVer( w http.ResponseWriter, r *http.Request ){
	
}
// 重命名
func (fa imp_http) ReName( w http.ResponseWriter, r *http.Request ){
	q_src_path := r.FormValue("path")
	q_name := r.FormValue("name")
	if !fa.fm.IsExist( q_src_path ){
		sendError(w, Error_FileNotExist); return;
	}
	if len(q_name) == 0 {
		sendError(w, Error_NewNameIsEmpty); return;
	}
	rnm_err := fa.fm.DoRename(q_src_path, q_name)
	if nil == rnm_err {
		sendSuccess(w, "")
	}else{
		sendError(w, rnm_err);
	}
}
// 新建文件夹
func (fa imp_http) NewFolder( w http.ResponseWriter, r *http.Request ){
	q_src_path := r.FormValue("path")
	q_src_path = stringtools.ClerPath( q_src_path )
	if !fa.fm.IsExist( stringtools.GetParentPath( q_src_path ) ){
		sendError(w, Error_ParentFolderNotExist); return;
	}
	rnm_err := fa.fm.DoNewFolder(q_src_path)
	if nil == rnm_err {
		sendSuccess(w, "")
	}else{
		sendError(w, rnm_err);
	}
}
func (fa imp_http) Info( w http.ResponseWriter, r *http.Request ){
	
}
func (fa imp_http) NameSearch( w http.ResponseWriter, r *http.Request ){
	
}
// 文件上传, 支持Form和Body上传方式
// 参数: Header("Save-Path", ["Formname-File"])
func (fa imp_http) Upload( w http.ResponseWriter, r *http.Request ){
	mr, err := r.MultipartReader( )
	if err == nil{
		p_name := r.Header.Get(header_formname_file)
		if len(p_name) == 0 {
			p_name = default_formname_file;
		}
		dst := "";
		for {
			p, err := mr.NextPart( )
			if nil == p || err == io.EOF {
				break
			}
			// 文件保存位置
			if p.FormName( ) == default_formname_fspath {
				dst = stringtools.ReadString(p)
			}
			if p.FormName( ) != p_name {
				continue
			}
			if len(dst) == 0 {
				sendError(w, errors.New("Cannot get header: Save-Path")); return
			}
			err = fa.fm.DoWrite(dst, p)
			if nil != err {
				sendError(w, err)
			}else{
				sendSuccess(w, "")
			}
			break
		}
		
	}else if nil != err && err == http.ErrNotMultipart {
		dst := r.Header.Get(default_formname_fspath)
		if len(dst) == 0 {
			sendError(w, errors.New("Cannot get header: Save-Path")); return
		}
		err := fa.fm.DoWrite(dst, r.Body)
		if nil != err {
			sendError(w, err)
		}else{
			sendSuccess(w, "")
		}
	}else{
		sendError(w, err)
	}
}
// 打开
func (fa imp_http) OpenFile( w http.ResponseWriter, r *http.Request ){
	token := stringtools.GetPathName( r.URL.Path )
	err, tokenObject := getFileTransferTokenObject(fa, token)
	if nil != err || nil == tokenObject {
		sendError(w, Error_OprationExpires); return
	}
	//fa.fm.RemoveToken(token)
	rd, r_err := fa.fm.DoRead(tokenObject.FilePath)
	if r_err != nil {
		sendError(w, r_err); return
	}
	io.Copy(w, rd)
}
// 下载
func (fa imp_http) Download( w http.ResponseWriter, r *http.Request ){
	token := stringtools.GetPathName( r.URL.Path )
	err, tokenObject := getFileTransferTokenObject(fa, token)
	if nil != err || nil == tokenObject {
		sendError(w, Error_OprationExpires); return
	}
	fa.fm.RemoveToken(token)
	rd, r_err := fa.fm.DoRead(tokenObject.FilePath)
	if r_err != nil {
		sendError(w, r_err); return
	}
	fileSize, fs_size_err := fa.fm.GetFileSize(tokenObject.FilePath)
	if nil != fs_size_err {
		sendError(w, fs_size_err); return
	}
	fileName := stringtools.GetPathName( tokenObject.FilePath )
	w.Header( ).Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header( ).Set("Content-Type", "application/octet-stream")
	w.Header( ).Set("Content-Length", stringtools.Int2String(fileSize))
	io.Copy(w, rd)
}
// =========================================================
type FsApiResponse struct{
	Code int
	Data string
}
func sendSuccess(w http.ResponseWriter, msg string){
	setJson(w)
	w.WriteHeader(http.StatusOK)
	w.Write(parseJson(http.StatusOK, msg))
}
func sendError( w http.ResponseWriter, err error ){
	setJson(w)
	w.WriteHeader(http.StatusBadRequest)
	w.Write( parseJson(http.StatusBadRequest, err.Error( )) )
}
func parseJson( code int, str string) []byte{
	bt, err := json.Marshal(FsApiResponse{Code: code, Data: str})
	if nil != err {
		fmt.Println("parseJson: ", err)
	}
	return bt
}
func setJson( w http.ResponseWriter ){
	w.Header( ).Set("Content-type", "application/json;charset=utf-8")
}
// 获取批文件量操作Token对象
func getFileBatchOperationTokenObject( fa imp_http, token string )(error, *FileBatchOperationTokenObject){
	tokenObject := fa.fm.GetToken(token)
	// 并保持刷新token的有效性, 除非终止操作否则都继续
	if nil == tokenObject || nil == tokenObject.TokenBody {
		return Error_Discontinue, nil
	}
	// fmt.Println("tokenObject: ", tokenObject)
	return nil, tokenObject.TokenBody.(*FileBatchOperationTokenObject)
}
// 获取文件传输Token对象
func getFileTransferTokenObject( fa imp_http, token string )(error, *FileTransferToken){
	tokenObject := fa.fm.GetToken(token)
	// 并保持刷新token的有效性, 除非终止操作否则都继续
	if nil == tokenObject || nil == tokenObject.TokenBody {
		return Error_Discontinue, nil
	}
	return nil, tokenObject.TokenBody.(*FileTransferToken)
}