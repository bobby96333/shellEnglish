package learn_lib

import (
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type Sentence struct {
	ch chan *SentenceResult
}

type SentenceResult struct{

	Sentence string
	LocalMedia string
	HasDownload bool
	callback func(info *SentenceResult)
}


func (this *SentenceResult)download(){

	socket5proxy,err:= proxy.SOCKS5("tcp","127.0.0.1:1080",nil,proxy.Direct)
	if(err!=nil){
		panic(err)
	}
	fhandler,err:= os.Create(this.LocalMedia)
	if(err!=nil) {
		panic(err)
	}
	defer fhandler.Close()

	mediaUrl:="http://translate.google.com/translate_tts"
	//ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q="+urlWord+"&tl=En-us
	params:=url.Values{}
	params.Add("ie","UTF-8")
	params.Add("total","1")
	params.Add("idx","0")
	params.Add("textlen", string(len(this.Sentence)))
	params.Add("client","tw-ob")
	params.Add("q",this.Sentence)
	params.Add("tl","En-us")
	mediaUrl+="?"+params.Encode()
	//urli,err:=url.Parse(mediaUrl)
	if (err!=nil) {
		panic(err)
	}
	httpTransport:=&http.Transport{}
	httpClient:=&http.Client{Transport:httpTransport}
	httpTransport.Dial=socket5proxy.Dial
	if resp,err:=httpClient.Get(mediaUrl);err==nil {
		defer resp.Body.Close()
		bs,err:=ioutil.ReadAll(resp.Body)
		if(err!=nil){
			panic(err)
		}
		ioutil.WriteFile(this.LocalMedia,bs,666);

	}else{
		panic(err)
	}
	this.HasDownload=true
	go this.callback(this)
}


func (this *Sentence) Init() {
	this.ch=make(chan *SentenceResult,10)
	go this.reduce()
}
func (this *Sentence) reduce() {
	var task *SentenceResult
	for{
		task = <- this.ch
		sha1 := Sha1(task.Sentence)
		task.LocalMedia="/tmp/shellEnglish_"+sha1
		if(!ExistsFile(task.LocalMedia)){
			//download
			task.download()
		}else{
			task.HasDownload=true
			go task.callback(task)
		}
	}
}


func (this *Sentence) Seek(sentence string,callback func(info *SentenceResult)) *SentenceResult{

	ret:=new (SentenceResult)
	ret.Sentence=sentence
	ret.HasDownload=false
	ret.callback=callback
	this.ch <- ret
	return ret

}
