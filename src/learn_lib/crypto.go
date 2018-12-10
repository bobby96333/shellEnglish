package learn_lib

import (
	"crypto/sha1"
	"fmt"
)

type Crypto struct{


}

func Sha1(txt string) string{

	hash:=sha1.New();
	hash.Write([]byte(txt))
	bs:=hash.Sum(nil)
	return fmt.Sprintf("%x", bs)
}