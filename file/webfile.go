package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/stlwtr/pan/conf"
	"github.com/stlwtr/pan/utils/httpclient"
)

// https://pan.baidu.com/api/sharedownload?sign=&timestamp=&clienttype=0&app_id=250528&web=1&dp-logid=55234600809338140145

const (
	WebMetasUri = "/api/sharedownload?sign=&timestamp=&clienttype=0&dp-logid=55234600809338140145&web=1"
)

type WebFile struct {
	AccessToken string
	Appid       string
	UK          string
	GID         string
}

type WebMetasResponse struct {
	ErrorCode    int    `json:"errno"`
	ServerTime   string `json:"server_time"`
	RequestID    int
	RequestIDStr string `json:"request_id"`
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

	fsIDsByte, err := json.Marshal(fsIDs)
	if err != nil {
		return ret, err
	}

	v := url.Values{}
	v.Add("app_id", string(f.Appid))
	query := v.Encode()
	log.Println("query:", query)
	requestUrl := conf.OpenApiDomain + WebMetasUri + "&" + query
	headers := map[string]string{}
	headers["Cookie"] = "BIDUPSID=2C0095CAF5313259404BEC23BD4CECFC; PSTM=1677252161; secu=1; PANWEB=1; MCITY=-42%3A; BAIDUID=0DEC33605BE8F36EA05775C0FC771A0F:FG=1;"

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
	b.Add("uk", strconv.Itoa(f.UK))
	b.Add("product", "mbox")
	b.Add("timestamp", "")
	b.Add("sign", "")
	b.Add("fid_list", string(fsIDsByte))
	if len(f.GID) > 0 {
		extra := map[string]string{}
		extra["type"] = "group"
		extra["gid"] = f.GID
		b_extra, _ := json.Marshal(extra)
		b.Add("extra", string(b_extra))
	}

	body := b.Encode()
	resp, err := httpclient.Post(requestUrl, headers, body)
	if err != nil {
		log.Println("httpclient.Get failed, err:", err)
		return ret, err
	}

	if resp.StatusCode != 200 {
		return ret, errors.New(fmt.Sprintf("HttpStatusCode is not equal to 200, httpStatusCode[%d], respBody[%s]", resp.StatusCode, string(resp.Body)))
	}

	if err := json.Unmarshal(resp.Body, &ret); err != nil {
		return ret, err
	}

	if ret.ErrorCode != 0 { //错误码不为0
		return ret, errors.New(fmt.Sprintf("error_code:%d", ret.ErrorCode))
	}

	ret.RequestID, err = strconv.Atoi(ret.RequestIDStr)
	if err != nil {
		return ret, err
	}

	return ret, nil
}
