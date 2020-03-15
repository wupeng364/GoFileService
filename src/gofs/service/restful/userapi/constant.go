// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package userapi

import (
	"errors"
)

const (
	baseurl = "/userapi"
)

// ErrorUserIDIsNil ErrorUserIdIsNil
var ErrorUserIDIsNil = errors.New("User ID cannot be empty")

// ErrorUserNameIsNil ErrorUserNameIsNil
var ErrorUserNameIsNil = errors.New("User name cannot be empty")

// ErrorPwdIsError ErrorPwdIsError
var ErrorPwdIsError = errors.New("User passwork is error")

// ErrorNotSupport ErrorNotSupport
var ErrorNotSupport = errors.New("This type is not supported")

// ErrorParamsNotEmpty ErrorParamsNotEmpty
var ErrorParamsNotEmpty = errors.New("Property name or content cannot be empty")