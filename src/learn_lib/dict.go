package learn_lib

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Dict struct {
}

const STATUS_SUCCESS int = 200
const STATUS_FAILD int = 500

type DictResult struct {
	Word        string
	KK          string
	KK_mp3      string
	IPA         string
	IPA_mp3     string
	Description string
	Status      int
}

func ExistsFile(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (this Dict) See(word string) *DictResult {
	md5code := fmt.Sprintf("%x", md5.Sum([]byte(word)))
	fpath := "/tmp/dict"
	if !ExistsFile(fpath) {
		os.Mkdir(fpath, 755)
	}
	fpath = fmt.Sprintf("%s/%s.json", fpath, md5code)
	var ret *DictResult
	if ExistsFile(fpath) {
		handler, err := os.Open(fpath)
		defer handler.Close()
		if err != nil {
			panic(err)
		}
		bs, err := ioutil.ReadAll(handler)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(bs, &ret)
		return ret
	}
	ret = this.NetSee(word)
	marshal, err := json.Marshal(ret)
	if err != nil {
		panic(err)
	}
	handler, err := os.Create(fpath)
	defer handler.Close()
	_, err = handler.Write(marshal)
	if err != nil {
		panic(err)
	}
	return ret

}

func (this Dict) NetSee(word string) *DictResult {

	dict := new(DictResult)
	dict.Word = word

	client := &http.Client{}
	url := fmt.Sprintf("http://www.iciba.com/%s", word)
	response, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	html := string(bytes)
	var ipa_html string
	var kk_html string
	i1 := 0
	i2 := 0
	i1 = strings.Index(html, "<span>英")
	if i1 != -1 {
		//		fmt.Println(html[i1:])
		i1 += len("<span>英")
		i2 = i1 + strings.Index(html[i1:], "</span>")
		i2 = i2 + strings.Index(html[i2+len("</span>"):], "</span>")
		ipa_html = html[i1:i2]

		i2 = strings.Index(html[i1:], "</span>")
		dict.IPA = html[i1 : i1+i2]

		i1 = strings.Index(ipa_html, "sound('")
		if i1 != -1 {
			i1 += len("sound('")
			i2 = i1 + strings.Index(ipa_html[i1:], "')")
			dict.IPA_mp3 = ipa_html[i1:i2]
		}
	}
	fmt.Println("1")
	i1 = strings.Index(html, "<span>美")
	if i1 != -1 {
		i1 += len("<span>美")
		i2 = i1 + strings.Index(html[i1:], "</span>")
		i2 = i2 + strings.Index(html[i2+len("</span>"):], "</span>")
		kk_html = html[i1:i2]

		i2 = strings.Index(html[i1:], "</span>")
		dict.KK = html[i1 : i1+i2]
		i1 = strings.Index(kk_html, "sound('")
		if i1 != -1 {
			i1 += len("sound('")
			i2 = i1 + strings.Index(kk_html[i1:], "')")
			dict.KK_mp3 = kk_html[i1:i2]
		}
	}

	i1 = strings.Index(html, "<ul class=\"base-list switch_part\" class=\"\">")

	if i1 != -1 {
		i1 += len("<ul class=\"base-list switch_part\" class=\"\">")
		i2 = strings.Index(html[i1:], "</ul>")
		dict.Description = html[i1 : i1+i2]
		cmp, _ := regexp.Compile("\\s*")
		dict.Description = cmp.ReplaceAllString(dict.Description, "")
		dict.Description = strings.Replace(dict.Description, "</li>", "\n", -1)
		cmp, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
		dict.Description = cmp.ReplaceAllString(dict.Description, "")
	}
	dict.Status = STATUS_SUCCESS
	return dict

}
