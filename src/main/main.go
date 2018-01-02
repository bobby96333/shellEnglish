package main

import (
	"container/list"
	"flag"
	"fmt"
	"io/ioutil"
	"learn_lib"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const LISTEN_TYPE_IPA string = "ipa"
const LISTEN_TYPE_KK string = "kk"

var listen_type = LISTEN_TYPE_KK
var delWords *learn_lib.WordMem

func main_init() {
	delWords = &learn_lib.WordMem{}
	delWords.Init("/tmp/dict/del.txt")
}

func main() {
	fmt.Println("welcome to bobby leaning...")
	main_init()
	listen_type_arg := flag.String("listen", "kk", "kk or ipa")
	flag.Usage = falg_usage
	flag.Parse()
	listen_type = *listen_type_arg

	fmt.Println("listen type:", listen_type)

	dict := new(learn_lib.Dict)
	var words *list.List
	words = list.New()

	if len(flag.Args()) > 0 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			panic(err)
		}
		bs, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
		str := string(bs)

		//cut text

		start := strings.Index(str, "#start")
		if start != -1 {
			start += len("#start")
			str = str[start:]
		}

		end := strings.Index(str, "#end")
		if end != -1 {
			str = str[:end]
		}

		//replace first upper case to lower case
		regcmp, _ := regexp.Compile("(?m)^([A-Z])[a-z]")
		indexs := regcmp.FindAllStringSubmatchIndex(str, -1)
		for _, indexRange := range indexs {
			stri := indexRange[2]
			wordStr := strings.ToLower(string(str[stri]))
			bs := []byte(str)
			bs[stri] = wordStr[0]
			str = string(bs)
		}

		//replace wrap to space
		str = strings.Replace(str, "\n", " ", -1)
		str = strings.Replace(str, "\r", "", -1)
		//get words

		words_slice := strings.Split(str, " ")

		for i1, word := range words_slice {
			//去空白
			if strings.TrimSpace(word) == "" {
				continue
			}
			//去己删除的
			if delWords.Exists(word) {
				continue
			}

			//检查是不是单词
			if !learn_lib.IsEnglishWord(word) {
				continue
			}

			repeat := false
			i2 := 0
			//去去重
			for e := words.Front(); e != nil; e = e.Next() {
				word2 := e.Value.(string)
				if i1 != i2 && word2 == word {
					repeat = true
					break
				}
				i2++
			}
			if repeat {
				continue
			}
			words.PushBack(word)
		}

	} else {
		fmt.Println("not found file arg")
		return
	}
	fmt.Println("inited")

	for {
		if words.Len() == 0 {
			break
		}
		for e := words.Front(); e != nil; e = e.Next() {
			word := e.Value.(string)
			if strings.TrimSpace(word) == "" {
				words.Remove(e)
			}
			dictr := dict.See(word)
			fmt.Println(dictr.Description)
			if dictr.KK_mp3 != "" {
				listen(dictr.KK_mp3)
			}
			var input string
		reinput:
			fmt.Scanln(&input)
			input = strings.ToLower(input)
			if input == ":past" {
				continue
			} else if input == ":del" {
				delWords.Append(word)
				delWords.Flush()
				olde := e
				if e.Prev() != nil {
					e = e.Prev()
				}
				words.Remove(olde)
				continue
			} else if input == ":info" {
				fmt.Printf("%+v\n", dictr)
				goto reinput
			} else if input == ":kk" {
				fmt.Println(dictr.KK)
				goto reinput
			} else if input == ":ipa" {
				fmt.Println(dictr.IPA)
				goto reinput
			} else if input == ":listen" {
				if listen_type == LISTEN_TYPE_IPA {
					listen(dictr.IPA_mp3)
				} else if listen_type == LISTEN_TYPE_KK {
					listen(dictr.KK_mp3)
				}
				goto reinput
			} else if input == ":all" {

				for e := words.Front(); e != nil; e = e.Next() {
					word := e.Value.(string)
					fmt.Print(word, "  ")
				}
				fmt.Print("\n")
				goto reinput
			} else if input == ":list" {

				for e := words.Front(); e != nil; e = e.Next() {
					word := e.Value.(string)
					dictr := dict.See(word)
					wordinfo := word
					if listen_type == LISTEN_TYPE_IPA {
						wordinfo += " " + dictr.IPA
					} else if listen_type == LISTEN_TYPE_KK {
						wordinfo += " " + dictr.KK
					}
					wordinfo += "\n   " + strings.Replace(dictr.Description, "\n", "\n   ", -1)
					fmt.Println(wordinfo)
				}
				fmt.Print("\n")
				goto reinput
			}

			if strings.TrimSpace(strings.ToLower(input)) == strings.TrimSpace(strings.ToLower(dictr.Word)) {
				fmt.Println(":right")
				continue
			}
			fmt.Println(dictr.Word, "*****************")
		}
	}

}

func falg_usage() {

	fmt.Println(`命令参数:
	--help 查看帮助
	--listen  kk是美语发音 ipa是英式发音
输入指令：
	:del 删除当前单词
	:kk 显示美音音标
	:ipa 显示英式音标
	:listen 重读
	:all  显示还有多少单词
	:list 显示词汇列表，可打印用于后期记忆
文本标记：
	#begin  取文本中的开始位置
	#end    取广西的结束位置
		
`)
}
func listen(url string) {
	if url == "" {
		fmt.Println("no found voice url")
		return
	}
	cmd := exec.Command("mplayer", url)
	cmd.Wait()
	cmd.Output()
}
