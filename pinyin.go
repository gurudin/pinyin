package pinyin

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var (
	// punctuations 标点符号
	punctuations = []string{
		// 逗号
		"，", ",",
		// 句号
		"。", ".",
		// 感叹号
		"！", "!",
		// 问号
		"？", "?",
		// 冒号
		"：", ":",
		// 分号
		"；", ";",
		// 左/右单引号
		"‘", " '", "’", " '",
		// 左/右双引号
		"“", ` "`, "”", ` "`,
		// 左/右直角引号
		"「", " [", "」", " ]",
		"『", " [", "』", " ]",
		// 左/右括号
		"（", " (", "）", " )",
		"〔", " [", "〕", " ]",
		"【", " [", "】", " ]",
		"{", " {", "}", " }",
		// 省略号
		"……", "...",
		// 破折号
		"——", "-",
		// 连接号
		"—", "-",
		// 左/右斜杆
		"/", " /", "\\", " \\",
		// 波浪线
		"～", "~",
		// 书名号
		"《", " <", "》", " >",
		"〈", " <", "〉", " >",
		// 间隔号
		"·", " ·",
		// 顿号
		"、", ",",
	}
	// finals 韵母表
	finals = []string{
		// a
		"a1", "ā", "a2", "á", "a3", "ǎ", "a4", "à",
		// o
		"o1", "ō", "o2", "ó", "o3", "ǒ", "o4", "ò",
		// e
		"e1", "ē", "e2", "é", "e3", "ě", "e4", "è",
		// i
		"i1", "ī", "i2", "í", "i3", "ǐ", "i4", "ì",
		// u
		"u1", "ū", "u2", "ú", "u3", "ǔ", "u4", "ù",
		// v
		"v1", "ǖ", "v2", "ǘ", "v3", "ǚ", "v4", "ǜ",

		// ai
		"ai1", "āi", "ai2", "ái", "ai3", "ǎi", "ai4", "ài",
		// ei
		"ei1", "ēi", "ei2", "éi", "ei3", "ěi", "ei4", "èi",
		// ui
		"ui1", "uī", "ui2", "uí", "ui3", "uǐ", "ui4", "uì",
		// ao
		"ao1", "āo", "ao2", "áo", "ao3", "ǎo", "ao4", "ào",
		// ou
		"ou1", "ōu", "ou2", "óu", "ou3", "ǒu", "ou4", "òu",
		// iu
		"iu1", "īu", "iu2", "íu", "iu3", "ǐu", "iu4", "ìu",

		// ie
		"ie1", "iē", "ie2", "ié", "ie3", "iě", "ie4", "iè",
		// ve
		"ue1", "üē", "ue2", "üé", "ue3", "üě", "ue4", "üè",
		// er
		"er1", "ēr", "er2", "ér", "er3", "ěr", "er4", "èr",

		// an
		"an1", "ān", "an2", "án", "an3", "ǎn", "an4", "àn",
		// en
		"en1", "ēn", "en2", "én", "en3", "ěn", "en4", "èn",
		// in
		"in1", "īn", "in2", "ín", "in3", "ǐn", "in4", "ìn",
		// un/vn
		"un1", "ūn", "un2", "ún", "un3", "ǔn", "un4", "ùn",

		// ang
		"ang1", "āng", "ang2", "áng", "ang3", "ǎng", "ang4", "àng",
		// eng
		"eng1", "ēng", "eng2", "éng", "eng3", "ěng", "eng4", "èng",
		// ing
		"ing1", "īng", "ing2", "íng", "ing3", "ǐng", "ing4", "ìng",
		// ong
		"ong1", "ōng", "ong2", "óng", "ong3", "ǒng", "ong4", "òng",
	}
)

// ConvertResult 转换后字符串
type ConvertResult string

// 字典
type (
	dictDir     [2]string
	surNamesDir [1]string
)

// Config request conifg.
type Config struct {
	Dict       dictDir
	Surnames   surNamesDir
	OutputType string
}

// InitConfig 初始化配置
var InitConfig Config

func init() {
	var dictDir = dictDir{
		"dict/words_0.dict",
		"dict/words_1.dict",
	}
	var surNamesDir = surNamesDir{
		"dict/surnames.dict",
	}
	InitConfig = Config{
		Dict:       dictDir,
		Surnames:   surNamesDir,
		OutputType: "none",
	}
}

// Convert 字符串转换拼音.
// strs: 转换字符串
// delimiter: 分隔符
func Convert(strs string, delimiter string) string {
	s := InitConfig.romanize(strs, false)

	return s
}

func (c *Config) prepare(s string) string {
	re := regexp.MustCompile(`[a-zA-Z0-9_-]+`)
	s = re.ReplaceAllStringFunc(s, func(repl string) string {
		return "\t" + repl
	})

	re = regexp.MustCompile(`~[^\p{Han}\p{P}\p{Z}\p{M}\p{N}\p{L}\t]~u`)

	return re.ReplaceAllString(s, "")
}

func (c *Config) romanize(s string, surnames bool) string {
	s = c.prepare(s)

	if surnames {
		for i := 0; i < len(c.Surnames); i++ {
			if !isChineseChar(s) {
				break
			}

			s = charToPinyin(s, c.Surnames[i])
		}
	}

	if isChineseChar(s) {
		for i := 0; i < len(c.Dict); i++ {
			if !isChineseChar(s) {
				break
			}

			s = charToPinyin(s, c.Dict[i])
		}
	}

	return s
}

func charToPinyin(s string, path string) string {
	file, _ := os.Open(path)

	defer file.Close()

	br := bufio.NewReader(file)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		tmp := strings.Split(string(a), ":")
		s = strings.Replace(s, tmp[0], tmp[1], -1)
	}

	s = strings.TrimSpace(s)
	s = strings.Replace(strings.Replace(s, "  ", " ", -1), "\t", " ", -1)

	return s
}

func isChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}

	return false
}
