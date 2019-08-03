package httpclienttools

/**
 *@description http客户端工具类
 *@author	wupeng364@outlook.com
*/
import (
	"os"
	"io"
	"time"
	"io/ioutil"
	"net/http"
	"strings"
	"encoding/json"
	"mime/multipart"
	"bytes"
)
const(
	_defaultTimeout = 30	// 单位s
)

func BuildUrlsWithMap(url string, params map[string]string) (string) {
	result := url
	p_len  := len(params)
	if params != nil && p_len > 0 {
		result += "?"
		for key, val := range params { 
			result += key+"="+val
			if( p_len > 1 ){
				result += "&"
			}
			p_len--
		}
	}
	return result
}
func BuildUrlsWithArrays(url string, params [][]string) (string) {
	result := url
	p_len  := len(params)
	if params != nil && p_len > 0 {
		result += "?"
		for i:=0; i<p_len; i++ {
			if len(params[i]) >= 2 {
				result += params[i][0]+"="+params[i][1]
				if( i < p_len-1 ){
					result += "&"
				}
			}
		}
	}
	return result
}
// Get
func Get(url string, params map[string]string, headers map[string]string)(*http.Response, error){
    return DoFormUrlEncoded("GET", url, params, headers, _defaultTimeout)
}
// Get Body
func GetResponseBody(url string, params map[string]string, headers map[string]string)(string, error){
    resp, err := Get(url, params, headers)
    defer func(){
	    if nil != resp {
		    resp.Body.Close()
	    }
    }()
    if nil != err {
	    return "", err
    }
    // red response 
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
       return "", err
    }
    return string(body), nil
}
// Post
func Post(url string, params map[string]string, headers map[string]string)(*http.Response, error){
    return DoFormUrlEncoded("POST", url, params, headers, _defaultTimeout)
}
// http Do Form Url Encoded
func DoFormUrlEncoded( reqType, url string, params, headers map[string]string, timeout int64)(*http.Response, error){
    
	// build query
	query := ""
	p_len  := len(params)
	if params != nil && p_len > 0 {
		for key, val := range params { 
			query += key+"="+val
			if( p_len > 1 ){
				query += "&"
			}
			p_len--
		}
	}
	
	// build request method
    req, err := http.NewRequest(reqType, url, strings.NewReader(query))
    if err != nil {
        return nil, err
    }
    
	// set headers 
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    if( headers !=nil && len(headers) > 0 ){
    	for key, val := range headers{
		    req.Header.Set(key, val)
    	}
    }
    
	// do request
	client := &http.Client{}
	if timeout > -1 {
		client.Timeout = time.Second*time.Duration( timeout )
	}
	
    return client.Do(req)
}

// http Do Form Url Encoded
func PostJson( url string, params interface{}, headers map[string]string) (*http.Response, error) {
    
	// build query
	query, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	
	// build request method
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
    if err != nil {
        return nil, err
    }
    
	// set headers 
    req.Header.Set("Content-Type", "application/json;charset=utf-8")
    if( headers !=nil && len(headers) > 0 ){
    	for key, val := range headers{
		    req.Header.Set(key, val)
    	}
    }
    
	// do request
	client := &http.Client{}
    return client.Do(req)
}
// 
func PostFile(url, filePath string) (*http.Response, error) {
	return PostMultiFile(url, filePath, "")
}
// http post multi file
func PostMultiFile(url, filePath, paramName string) (*http.Response, error) {
    body_buf := bytes.NewBufferString("")
    body_writer := multipart.NewWriter(body_buf)
	if len(paramName) == 0 {
		paramName = "file"
	}
    _, err := body_writer.CreateFormFile(paramName, filePath)
    if err != nil {
        return nil, err
    }

    fh, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }

    boundary := body_writer.Boundary()
    close_buf := bytes.NewBufferString("\r\n--"+boundary+"--\r\n")

    request_reader := io.MultiReader(body_buf, fh, close_buf)
    fi, err := fh.Stat( )
    if err != nil {
        return nil, err
    }
    req, err := http.NewRequest("POST", url, request_reader)
    if err != nil {
        return nil, err
    }

    req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
    req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())

    client := &http.Client{}
    return  client.Do(req)
}