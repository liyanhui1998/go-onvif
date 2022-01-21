package hikvision

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"
)

/*
	说明:
		海康摄像头http接口获取抓拍图像信息
	参数:
		url,username,passwd
	返回:
		二进制数据(图像)
*/
func DowloadHttpSnapshotImage(url, username, passwd string) ([]byte, error) {
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
