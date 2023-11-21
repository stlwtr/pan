package file

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/stlwtr/pan/conf"
	"github.com/stlwtr/pan/utils"
	"github.com/stlwtr/pan/utils/httpclient"
)

// https://pan.baidu.com/api/sharedownload?sign=&timestamp=&clienttype=0&app_id=250528&web=1&dp-logid=55234600809338140145

const (
	WebMetasUri = "/api/sharedownload?sign=&timestamp=&clienttype=0&app_id=%s&web=1&dp-logid=17724300763630380079"
)

type WebFile struct {
	AccessToken string
	Appid       string
	UK          string
	GID         string
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
	}
}

/**
 *
 */
func NewWebFileClient(accessToken string, appid, uk, gid string) *WebFile {
	return &WebFile{
		AccessToken: accessToken,
		Appid:       appid,
		UK:          uk,
		GID:         gid,
	}
}

// 通过FsID获取文件信息
func (f *WebFile) Metas(fsIDs []uint64) (WebMetasResponse, error) {
	ret := WebMetasResponse{}

	fsIDsByte, err := utils.MarshalJSON(fsIDs)
	if err != nil {
		return ret, err
	}

	// v := url.Values{}
	// v.Add("app_id", string(f.Appid))
	// query := v.Encode()
	requestUrl := conf.OpenApiDomain + fmt.Sprintf(WebMetasUri, f.Appid)
	log.Println("requestUrl:", requestUrl)
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Mobile Safari/537.36",
		// "Accept-Encoding":  "gzip, deflate, br",
		"Accept":           "application/json, text/plain, */*",
		"X-Requested-With": "XMLHttpRequest",
		"Content-Type":     "application/x-www-form-urlencoded",
		"Origin":           "https://pan.baidu.com",
		"Referer":          "https://pan.baidu.com/disk/main?from=homeFlow",
		"Cookie":           "BAIDUID=B19DCDB85FC78E8C2A3B0B670816EA47:FG=1; BIDUPSID=B19DCDB85FC78E8CD828AC48FA6DABEE; PSTM=1699403304; newlogin=1; csrfToken=bn0lCF1kGBXvd0bSX4K6XEiV; RT=\"z=1&dm=baidu.com&si=0568efd4-e58f-48bd-8974-3fa09ea2b48a&ss=lp67pcx7&sl=2&tt=i7&bcn=https%3A%2F%2Ffclog.baidu.com%2Flog%2Fweirwood%3Ftype%3Dperf&ld=1yon&nu=9y8m6cy&cl=166s&ul=20ad&hd=20ky\"; BDUSS=JZeGlaNG1zOUViRElid2dVN0pjQX5-U2oxb2NxeTljM3ZMSnRsSFRBZlNQNEpsRVFBQUFBJCQAAAAAAAAAAAEAAACaLtGZeWFuZ2xpc2h1YW5faGUAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANKyWmXSslplTG; STOKEN=395e00420f0a0a941b431b5ed9c9d1c6a66959e84d73b2395c11ba0710fed197; ndut_fmt=863C835426EABE4824876960C7D69C255950A615F0529DE4824CC4C4B20C08C1; PANPSC=11942364391563931001%3A0GmL5sexpYY1JSw%2Btun8DDuv496EomYGJZzphDaPU5WERqtn4XIbyQjn4HR56dePCnOAuwe5%2BoHiueYauXAiTv6egEyo8Awm1%2FI1E4WnvHFk%2FzCt4mbzqKKwET9hWE%2FHP1LKwX8sNjJfZeVSwpGLRrBu2GS532flr%2Fo2LS0zEd5CNDWYUdN3xbnkjkK9qX6SfW3xVeVF9lU%3D; ab_sr=1.0.1_YTA1NDRmZmJkZGQyYjU4OTczY2M5OGM0OGYyMGFlYzhmZGU2YTgzM2EyNWIyMDIxZjIxYTY5MWZkNjVlNDUyMGZiMTIyOGVjNmZlY2EyOTZhNGM2ZDMxNzY5MjIxNzczNjVjODEwM2QwYWY1MmUyMTU1YTM1ZDhmZjM0NzQ0NjAwNGU1YjAxN2E5YTBiYTI0ZjY4MjQyNThlZmYwZTUxYmNhZGE2YTNiYTU1M2NkOGQzODViZjYxOWFjYjk2Yzk1",
	}

	// 	encrypt: 0
	// uk: 782095374
	// product: mbox
	// timestamp:
	// sign:
	// primaryid: 8292282171460833016
	// fid_list: [58578524529993]
	// extra: {"type":"group","gid":"394046314352271309"}

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
	// body = "encrypt=0&uk=782095374&product=mbox&timestamp=&sign=&primaryid=8292282171460833016&fid_list=%5B561191839888988%5D&extra=%7B%22type%22%3A%22group%22%2C%22gid%22%3A%22394046314352271309%22%7D"
	log.Println("reqBody:", reqBody)
	headers["User-Agent"] = httpclient.GetRandomUserAgent()
	resp, err := httpclient.Post(requestUrl, headers, reqBody)
	// bodyReader, _ := gzip.NewReader(io.Reader(resp.Body))
	// log.Println("resp header", resp.Header)
	body_str, _ := utils.UnescapeUnicode(resp.Body)
	log.Println("resp body:", body_str)
	log.Println("resp header:", resp.Header)
	if err != nil {
		log.Println("httpclient.Get failed, err:", err)
		return ret, err
	}

	if resp.StatusCode != 200 {
		return ret, errors.New(fmt.Sprintf("HttpStatusCode is not equal to 200, httpStatusCode[%d], respBody[%s]", resp.StatusCode, body_str))
	}

	if err := utils.UnmarshalJSON([]byte(body_str), &ret); err != nil {
		return ret, err
	}

	if ret.ErrorCode != 0 { //错误码不为0
		return ret, errors.New(fmt.Sprintf("error_code:%d", ret.ErrorCode))
	}

	ret.RequestIDStr = strconv.Itoa(ret.RequestID)
	if err != nil {
		return ret, err
	}

	return ret, nil
}
