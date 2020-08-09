// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 文件异步操作, 一般用于批量操作或后台任务

package asynctask

import (
	"encoding/json"
	"errors"
	"gofs/base/httpserver"
	"gofs/base/signature"
	"gofs/data/filemanage"
	"gofs/data/filepermission"
	"gutils/mloader"
	"gutils/strtool"
	"net/http"
	"time"
)

// MoveFileTokenObject 移动Token保存对象
type MoveFileTokenObject struct {
	CountIndex    int64  // 已处理的个数
	ErrorString   string // 错误信息
	Src           string // 当前正在处理的源路径
	Dst           string // 当前正在处理的目标路径
	IsSrcExist    bool   // 源路径是否存在
	IsDstExist    bool   // 目标路径是否存在
	IsReplace     bool   // 是否替换, 单次中断执行指令, 读取后设为false
	IsReplaceAll  bool   // 是否替换, 单次API执行指令, 设置后后续中断时自动替换
	IsIgnore      bool   // 是否忽略错误, 单次中断执行指令, 读取后设为false
	IsIgnoreAll   bool   // 是否忽略错误, 单次API执行指令, 设置后后续中断时自动替换
	IsComplete    bool   // 是否执行完毕
	IsDiscontinue bool   // 是否已中断操作
}

// MoveFile 移动|文件夹
type MoveFile struct {
	fm   *filemanage.FileManager
	sg   *signature.Signature
	fpms *filepermission.FPmsManager
}

// Name 动作名字
func (task MoveFile) Name() string {
	return "MoveFile"
}

// Init 初始化对象
func (task *MoveFile) Init(mctx *mloader.Loader) AsyncTask {
	task.fm = mctx.GetModuleByTemplate(task.fm).(*filemanage.FileManager)
	task.sg = mctx.GetModuleByTemplate(task.sg).(*signature.Signature)
	task.fpms = mctx.GetModuleByTemplate(task.fpms).(*filepermission.FPmsManager)
	return task
}

// Execute 动作执行, 返回一个tooken
func (task *MoveFile) Execute(r *http.Request) (string, error) {
	qSrcPath := r.FormValue("srcPath")
	qDstPath := r.FormValue("dstPath")
	qReplace := strtool.String2Bool(r.FormValue("replace"))
	qIgnore := strtool.String2Bool(r.FormValue("ignore"))

	if len(qSrcPath) == 0 {
		return "", errors.New("srcPath parameter not found")
	}
	if len(qDstPath) == 0 {
		return "", errors.New("dstPath parameter not found")
	}
	userID := task.sg.GetUserID4Request(r)
	if !task.fpms.HashPermission(userID, qSrcPath, filepermission.WRITE) {
		return "", ErrorPermissionInsufficient
	}
	if !task.fpms.HashPermission(userID, qDstPath, filepermission.WRITE) {
		return "", ErrorPermissionInsufficient
	}
	// 异步处理, 返回一个Token用于查询进度
	token := task.fm.AskToken("", &MoveFileTokenObject{
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
	})
	go func(token string) {
		moveDirErr := task.fm.DoMove(qSrcPath, qDstPath, qReplace, qIgnore, func(s_src, s_dst string, moveErr *filemanage.MoveError) error {
			// 获取令牌数据, 不存在则说明已经销毁
			// 并保持刷新token的有效性, 除非终止操作否则都继续
			tokenBody, tokenErr := task.getTokenObject(token)
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
						reMoveErr := task.fm.DoMove(s_src, s_dst, true, false, func(s, d string, mErr *filemanage.MoveError) error {
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
					tokenBody, tokenErr := task.getTokenObject(token)
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
							moveCopyErr := task.fm.DoMove(s_src, s_dst, true, false, func(s, d string, mErr *filemanage.MoveError) error {
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
		tokenBody := task.fm.GetToken(token)
		if nil != tokenBody {
			tokenBody := tokenBody.(*MoveFileTokenObject)
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
	return token, nil
}

// Status 查询动作状态, 在内部返回数据
func (task *MoveFile) Status(w http.ResponseWriter, r *http.Request) {
	qToken := r.FormValue("token")
	task.fm.RefreshToken(qToken)
	tokenBody, tokenErr := task.getTokenObject(qToken)
	// fmt.Println("Token: ", r.Method, qToken, tokenBody)
	if nil == tokenBody || nil != tokenErr {
		httpserver.SendError(w, ErrorOprationExpires)
		return
	}
	qOperation := r.FormValue("operation")
	// 用于获取令牌信息
	if len(qOperation) == 0 {
		bt, _ := json.Marshal(tokenBody)
		httpserver.SendSuccess(w, string(bt))

		// 用于操作|中断
	} else {
		switch qOperation {
		// 忽略单个 错误
		case "ignore":
			tokenBody.ErrorString = ""
			tokenBody.IsIgnore = true
			break
		// 为后续的 错误 执行忽略
		case "ignoreall":
			tokenBody.ErrorString = ""
			tokenBody.IsIgnoreAll = true
			break
		// 覆盖单个 已存在 错误
		case "replace":
			tokenBody.ErrorString = ""
			tokenBody.IsReplace = true
			break
		// 每次都覆盖 已存在 错误
		case "replaceall":
			tokenBody.ErrorString = ""
			tokenBody.IsReplaceAll = true
			break
		// 立即中断操作
		case "discontinue":
			tokenBody.ErrorString = ""
			tokenBody.IsComplete = true
			tokenBody.IsDiscontinue = true
			// task.fm.RemoveToken(qToken)
			break
		default:
			httpserver.SendError(w, ErrorOprationFailed)
			return
		}
		httpserver.SendSuccess(w, "")
	}
}

// getTokenObject 获取文件传输Token对象
func (task *MoveFile) getTokenObject(token string) (*MoveFileTokenObject, error) {
	tokenBody := task.fm.GetToken(token)
	// 并保持刷新token的有效性, 除非终止操作否则都继续
	if nil == tokenBody {
		return nil, ErrorOprationExpires
	}
	ftt, ok := tokenBody.(*MoveFileTokenObject)
	if ok {
		return ftt, nil
	}
	return nil, ErrorOprationExpires
}
