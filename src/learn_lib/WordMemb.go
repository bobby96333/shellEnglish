package learn_lib

import (
	"io/ioutil"
	"os"
	"strings"
)

type WordMem struct {
	fpath string
	words map[string]bool
}

func (this *WordMem) Init(fpath string) {
	//init
	this.words = make(map[string]bool)
	this.fpath = fpath
	//file load
	_, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		return
	}
	file, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	//split data format
	str := string(bs)
	for _, word := range strings.Split(str, ",") {
		if _, has := this.words[word]; !has {
			this.words[word] = true
		}
	}
}
func (this *WordMem) Exists(word string) bool {
	_, has := this.words[word]
	return has
}

func (this *WordMem) Append(word string) {
	this.words[word] = true
}

func (this *WordMem) Flush() {

	file, err := os.Create(this.fpath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for word, _ := range this.words {
		file.Write([]byte(word + ","))
	}

}
