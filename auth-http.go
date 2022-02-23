/*
 * @Author: YanHui Li
 * @Date: 2022-02-16 11:23:30
 * @LastEditTime: 2022-02-23 17:02:21
 * @LastEditors: YanHui Li
 * @Description:
 * @FilePath: \go-onvif\auth-http.go
 *
 */
package onvif

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

//通过http base 认证方式获取设备快照,返回图片二进制
func HttpBaseAuthSnapshotImage(url, username, passwd string) ([]byte, error) {
	httpClient := &http.Client{}
	/* 生成需要访问的http.Request信息 */
	if reqest, err := http.NewRequest("GET", url, nil); err == nil {
		/*
			鉴权方法
				在http请求头中添加
					Authorization  值为 "Basic " + Base64("name:passwd")
				例如:
					admin:123qweasdZXC Base64 后 YWRtaW46MTIzcXdlYXNkWlhD
					则 Authorization 的值为 "Base64 YWRtaW46MTIzcXdlYXNkWlhD"
		*/
		reqest.Header.Add("Authorization", "Base64 "+base64.StdEncoding.EncodeToString([]byte(username+":"+passwd)))
		if response, err := httpClient.Do(reqest); err == nil {
			defer response.Body.Close()
			if imageBytes, err := io.ReadAll(response.Body); err == nil {
				return imageBytes, nil
			} else {
				return nil, err
			}
		}
	} else {
		return nil, err
	}
	return nil, errors.New("unknown error")
}

//通过http digest 认证方式获取设备快照,返回图片二进制
func HttpDigestAuthGetSnapshotImage(url, username, password string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		return nil, fmt.Errorf("recieved status code '%d' auth skipped", resp.StatusCode)
	}
	digestParts := digestParts(resp)
	digestParts["uri"] = url
	digestParts["method"] = "GET"
	digestParts["username"] = username
	digestParts["password"] = password
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", getDigestAuthrization(digestParts))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		return body, nil
	}
	return nil, fmt.Errorf("response status code '%v'", resp.StatusCode)
}

func digestParts(resp *http.Response) map[string]string {
	result := map[string]string{}
	if len(resp.Header["Www-Authenticate"]) > 0 {
		wantedHeaders := []string{"nonce", "realm", "qop"}
		responseHeaders := strings.Split(resp.Header["Www-Authenticate"][0], ",")
		for _, r := range responseHeaders {
			for _, w := range wantedHeaders {
				if strings.Contains(r, w) {
					result[w] = strings.Split(r, `"`)[1]
				}
			}
		}
	}
	return result
}

func getMD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getCnonce() string {
	b := make([]byte, 8)
	io.ReadFull(rand.Reader, b)
	return fmt.Sprintf("%x", b)[:16]
}

func getDigestAuthrization(digestParts map[string]string) string {
	d := digestParts
	ha1 := getMD5(d["username"] + ":" + d["realm"] + ":" + d["password"])
	ha2 := getMD5(d["method"] + ":" + d["uri"])
	nonceCount := 00000001
	cnonce := getCnonce()
	response := getMD5(fmt.Sprintf("%s:%s:%v:%s:%s:%s", ha1, d["nonce"], nonceCount, cnonce, d["qop"], ha2))
	authorization := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", cnonce="%s", nc="%v", qop="%s", response="%s", algorithm="md5"`,
		d["username"], d["realm"], d["nonce"], d["uri"], cnonce, nonceCount, d["qop"], response)
	return authorization
}
