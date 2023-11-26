package file

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/stlwtr/pan/conf"
	"github.com/stlwtr/pan/utils"
	"github.com/stlwtr/pan/utils/httpclient"
)

// https://pan.baidu.com/api/sharedownload?sign=&timestamp=&clienttype=0&app_id=250528&web=1&dp-logid=55234600809338140145
// https://pan.baidu.com/mbox/msg/shareinfo?from_uk=1102905323160&msg_id=3938873226778728013&type=2&num=50&page=1&fs_id=449076463385100&gid=394046314352271309&limit=50&desc=1&clienttype=0&app_id=250528&web=1&dp-logid=35518200857989140084
const (
	WebMetasUri = "/api/sharedownload?sign=&timestamp=&clienttype=0&app_id=%s&web=1&dp-logid=17724300763630380079"
	WebListUri  = "/mbox/msg/shareinfo?from_uk=%s&msg_id=%s&type=2&num=%d&page=%d&fs_id=%s&gid=%s&limit=%d&desc=1&clienttype=0&app_id=%s&web=1&dp-logid=83769400584490230081"
)

type WebFile struct {
	AccessToken string
	Appid       string
	UK          string
	GID         string
	MID         string
}

type WebMetasResponse struct {
	ErrorCode    int   `json:"errno"`
	ServerTime   int64 `json:"server_time"`
	RequestID    int   `json:"request_id"`
	RequestIDStr string
	List         []struct {
		FsID        uint64            `json:"fs_id"`
		Path        string            `json:"path"`
		Category    int               `json:"category"`
		FileName    string            `json:"server_filename"`
		IsDir       int               `json:"isdir"`
		Size        int               `json:"size"`
		Md5         string            `json:"md5"`
		PathMd5     string            `json:"path_md5"`
		DLink       string            `json:"dlink"`
		Thumbs      map[string]string `json:"thumbs"`
		ServerCtime int               `json:"server_ctime"`
		ServerMtime int               `json:"server_mtime"`
	} `json:"list"`
}

type WebListResponse struct {
	ErrorCode  int   `json:"errno"`
	ServerTime int64 `json:"server_time"`
	RequestID  int   `json:"request_id"`
	HasMore    int   `json:"has_more"`
	Records    []struct {
		Category    int    `json:"category"`
		FsID        uint64 `json:"fs_id"`
		Path        string `json:"path"`
		IsDir       int    `json:"isdir"`
		FileName    string `json:"server_filename"`
		ServerCtime int    `json:"server_ctime"`
		ServerMtime int    `json:"server_mtime"`
		Size        int    `json:"size"`
	} `json:"records"`
}

/**
 *
 */
func NewWebFileClient(accessToken string, appid, uk, gid, mid string) *WebFile {
	return &WebFile{
		AccessToken: accessToken,
		Appid:       appid,
		UK:          uk,
		GID:         gid,
		MID:         mid,
	}
}

//

func (f *WebFile) Headers() map[string]string {
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Mobile Safari/537.36",
		// "Accept-Encoding":  "gzip, deflate, br",
		"Accept": "application/json, text/plain, */*",
		// "Accept-Language":  "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"X-Requested-With": "XMLHttpRequest",
		"Content-Type":     "application/x-www-form-urlencoded",
		"Origin":           "https://pan.baidu.com",
		"Referer":          "https://pan.baidu.com/disk/main",
		"Cookie":           "BAIDUID=B19DCDB85FC78E8C2A3B0B670816EA47:FG=1; BIDUPSID=B19DCDB85FC78E8CD828AC48FA6DABEE; PSTM=1699403304; newlogin=1; csrfToken=bn0lCF1kGBXvd0bSX4K6XEiV; RT=\"z=1&dm=baidu.com&si=0568efd4-e58f-48bd-8974-3fa09ea2b48a&ss=lp67pcx7&sl=2&tt=i7&bcn=https%3A%2F%2Ffclog.baidu.com%2Flog%2Fweirwood%3Ftype%3Dperf&ld=1yon&nu=9y8m6cy&cl=166s&ul=20ad&hd=20ky\"; BDUSS=JZeGlaNG1zOUViRElid2dVN0pjQX5-U2oxb2NxeTljM3ZMSnRsSFRBZlNQNEpsRVFBQUFBJCQAAAAAAAAAAAEAAACaLtGZeWFuZ2xpc2h1YW5faGUAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANKyWmXSslplTG; STOKEN=395e00420f0a0a941b431b5ed9c9d1c6a66959e84d73b2395c11ba0710fed197; ndut_fmt=863C835426EABE4824876960C7D69C255950A615F0529DE4824CC4C4B20C08C1; PANPSC=11942364391563931001%3A0GmL5sexpYY1JSw%2Btun8DDuv496EomYGJZzphDaPU5WERqtn4XIbyQjn4HR56dePCnOAuwe5%2BoHiueYauXAiTv6egEyo8Awm1%2FI1E4WnvHFk%2FzCt4mbzqKKwET9hWE%2FHP1LKwX8sNjJfZeVSwpGLRrBu2GS532flr%2Fo2LS0zEd5CNDWYUdN3xbnkjkK9qX6SfW3xVeVF9lU%3D; ab_sr=1.0.1_YTA1NDRmZmJkZGQyYjU4OTczY2M5OGM0OGYyMGFlYzhmZGU2YTgzM2EyNWIyMDIxZjIxYTY5MWZkNjVlNDUyMGZiMTIyOGVjNmZlY2EyOTZhNGM2ZDMxNzY5MjIxNzczNjVjODEwM2QwYWY1MmUyMTU1YTM1ZDhmZjM0NzQ0NjAwNGU1YjAxN2E5YTBiYTI0ZjY4MjQyNThlZmYwZTUxYmNhZGE2YTNiYTU1M2NkOGQzODViZjYxOWFjYjk2Yzk1",
	}
	return headers
}

// 通过FsID获取文件信息
func (f *WebFile) Metas(fsIDs []uint64) (WebMetasResponse, error) {
	ret := WebMetasResponse{}

	fsIDsByte, err := utils.MarshalJSON(fsIDs)
	if err != nil {
		return ret, err
	}

	requestUrl := conf.OpenApiDomain + fmt.Sprintf(WebMetasUri, f.Appid)

	headers := f.Headers()

	b := url.Values{}
	b.Add("encrypt", "0")
	b.Add("uk", f.UK)
	b.Add("product", "mbox")
	b.Add("timestamp", "")
	b.Add("primaryid", "8292282171460833016")
	b.Add("sign", "")
	b.Add("fid_list", string(fsIDsByte))
	if len(f.GID) > 0 {
		extra := map[string]string{}
		extra["type"] = "group"
		extra["gid"] = f.GID
		b_extra, _ := utils.MarshalJSON(extra)
		b.Add("extra", string(b_extra))
	}

	reqBody := b.Encode()
	headers["User-Agent"] = httpclient.GetRandomUserAgent()
	resp, err := httpclient.Post(requestUrl, headers, reqBody)
	body_str, _ := utils.UnescapeUnicode(resp.Body)
	if err != nil {
		log.Println("httpclient.Get failed, err:", err)
		return ret, err
	}

	if resp.StatusCode != 200 {
		return ret, fmt.Errorf(fmt.Sprintf("HttpStatusCode is not equal to 200, httpStatusCode[%d], respBody[%s]", resp.StatusCode, body_str))
	}

	if err := utils.UnmarshalJSON([]byte(body_str), &ret); err != nil {
		return ret, err
	}

	if ret.ErrorCode != 0 { //错误码不为0
		return ret, fmt.Errorf(fmt.Sprintf("error_code:%d", ret.ErrorCode))
	}

	ret.RequestIDStr = strconv.Itoa(ret.RequestID)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

// 通过FsID获取子目录信息，page从1开始
func (f *WebFile) List(fsid string, page int) (WebListResponse, error) {
	ret := WebListResponse{}

	numPerPage := 1000
	requestUrl := conf.OpenApiDomain + fmt.Sprintf(WebListUri, f.UK, f.MID, numPerPage, page, fsid, f.GID, numPerPage, f.Appid)

	log.Println("requestUrl = ", requestUrl)

	headers := f.Headers()

	headers["User-Agent"] = httpclient.GetRandomUserAgent()
	resp, err := httpclient.Post(requestUrl, headers, "")

	if resp.StatusCode != 200 {
		return ret, fmt.Errorf(fmt.Sprintf("HttpStatusCode is not equal to 200, httpStatusCode[%d], respBody[%s]", resp.StatusCode, string(resp.Body)))
	}

	log.Println("resp = ", string(resp.Body))

	if err := utils.UnmarshalJSON(resp.Body, &ret); err != nil {
		return ret, err
	}

	if ret.ErrorCode != 0 { //错误码不为0
		return ret, fmt.Errorf(fmt.Sprintf("error_code:%d", ret.ErrorCode))
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}

// 通过FsID获取所有子目录信息
func (f *WebFile) Listall(fsid string) (WebListResponse, error) {
	ret := WebListResponse{}

	numPerPage := 1000
	page := 1
	// WebListUri  = "/mbox/msg/shareinfo?from_uk=%s&msg_id=%s&type=2&num=%d&page=%d&fs_id=%s&gid=%s&limit=%d&desc=1&clienttype=0&app_id=%s&web=1&dp-logid=83769400584490230081"
	requestUrl := conf.OpenApiDomain + fmt.Sprintf(WebListUri, f.UK, f.MID, numPerPage, page, fsid, f.GID, numPerPage, f.Appid)

	headers := f.Headers()

	headers["User-Agent"] = httpclient.GetRandomUserAgent()
	resp, err := httpclient.Post(requestUrl, headers, "")

	if resp.StatusCode != 200 {
		return ret, fmt.Errorf(fmt.Sprintf("HttpStatusCode is not equal to 200, httpStatusCode[%d], respBody[%s]", resp.StatusCode, string(resp.Body)))
	}

	if err := utils.UnmarshalJSON(resp.Body, &ret); err != nil {
		return ret, err
	}

	if ret.ErrorCode != 0 { //错误码不为0
		return ret, fmt.Errorf(fmt.Sprintf("error_code:%d", ret.ErrorCode))
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}

// 通过FsID，本地路径localPath进行比对，输出未下载到本地的文件
func (f *WebFile) Compare(fsid string, localPath string) (WebListResponse, error) {
	ret := WebListResponse{}

	numPerPage := 100000
	//
	requestUrl := conf.OpenApiDomain + fmt.Sprintf(WebListUri, f.UK, numPerPage, 0, fsid, f.GID, numPerPage, f.Appid)

	headers := f.Headers()

	headers["User-Agent"] = httpclient.GetRandomUserAgent()
	resp, err := httpclient.Post(requestUrl, headers, "")

	if resp.StatusCode != 200 {
		return ret, fmt.Errorf(fmt.Sprintf("HttpStatusCode is not equal to 200, httpStatusCode[%d], respBody[%s]", resp.StatusCode, string(resp.Body)))
	}

	if err := utils.UnmarshalJSON(resp.Body, &ret); err != nil {
		return ret, err
	}

	if ret.ErrorCode != 0 { //错误码不为0
		return ret, fmt.Errorf(fmt.Sprintf("error_code:%d", ret.ErrorCode))
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}
