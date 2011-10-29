
// 人人网图片上传
package main

import (
	"../curl/_obj/curl"
	"time"
	"regexp"
	"fmt"
)


func getUploadUrl() string {
	page := ""
	easy := curl.EasyInit()
	defer easy.Cleanup()

	u := "http://3g.renren.com/album/wuploadphoto.do"
	easy.Setopt(curl.OPT_URL, u)
	easy.Setopt(curl.OPT_COOKIEFILE, "./cookie.jar")
	easy.Setopt(curl.OPT_COOKIEJAR, "./cookie.jar")
	easy.Setopt(curl.OPT_VERBOSE, true)
	easy.Setopt(curl.OPT_WRITEFUNCTION, func(ptr []byte, size uintptr, _ interface{}) uintptr {
		page += string(ptr)
		return size
	})
	easy.Perform()
	// extract url from
	// <form enctype="multipart/form-data" action="http://3g.renren.com/album/wuploadphoto.do?type=3&amp;sid=zv3tiXTZr6Cu1rj5dhgX_X"
	pattern, _ := regexp.Compile(`action="(.*?)"`)

	if matches := pattern.FindStringSubmatch(page); len(matches) == 2 {
		return matches[1]
	}
	return ""
}


func main() {
	// init the curl session

	easy := curl.EasyInit()
	defer easy.Cleanup()

	posturl := getUploadUrl()

	easy.Setopt(curl.OPT_URL, posturl)

	easy.Setopt(curl.OPT_PORT, 80)
	easy.Setopt(curl.OPT_VERBOSE, true)

	// save cookie and load cookie
	easy.Setopt(curl.OPT_COOKIEFILE, "./cookie.jar")
	easy.Setopt(curl.OPT_COOKIEJAR, "./cookie.jar")

	// disable HTTP/1.1 Expect: 100-continue
	easy.Setopt(curl.OPT_HTTPHEADER, []string{"Expect:"})

	form := curl.NewForm()
	form.Add("albumid", "452618633")
	form.AddFile("theFile", "path_to_your_img.jpg")
	form.Add("description", "我就尝试下这段代码靠谱不。。最后一张")
	form.Add("post", "上传照片")

	easy.Setopt(curl.OPT_HTTPPOST, form)

	// print upload progress
	easy.Setopt(curl.OPT_NOPROGRESS, false)
	easy.Setopt(curl.OPT_PROGRESSFUNCTION, func (_ interface{}, dltotal float64, dlnow float64, ultotal float64, ulnow float64) int {
		fmt.Printf("Download %3.2f%%, Uploading %3.2f%%\r", dlnow/dltotal * 100, ulnow/ultotal * 100)
		return 0
	})

	if err := easy.Perform(); err != nil {
		println("ERROR: ", err.String(), err)
	}

	time.Sleep(1000000000)			// wait gorotine

}
