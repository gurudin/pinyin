package pinyin

import (
	"fmt"
	"testing"
)

func TestConvert(t *testing.T) {
	fmt.Println(Convert("拼音"))
}

func TestUnicodeConvert(t *testing.T) {
	fmt.Println(UnicodeConvert("拼音"))
}

func TestASCIIConvert(t *testing.T) {
	fmt.Println(ASCIIConvert("拼音"))
}

func TestName(t *testing.T) {
	fmt.Println(Name("冒顿单于").None())
	fmt.Println(Name("冒顿单于").Unicode())
	fmt.Println(Name("冒顿单于").ASCII())
}

func TestPermalink(t *testing.T) {
	fmt.Println(Permalink("go语言汉字转换拼音", "-"))
}

func TestAbbr(t *testing.T) {
	fmt.Println(Abbr("获取首字母带分隔符的拼音字符串", ""))
	fmt.Println(Abbr("获取首字母带分隔符的拼音字符串", "-"))
}
