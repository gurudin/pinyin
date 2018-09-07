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
	punctuations = map[string]string{
		"，":  ",",
		"。":  ".",
		"！":  "!",
		"？":  "?",
		"：":  ":",
		"；":  ";",
		"‘":  "'",
		"’":  "'",
		"“":  `"`,
		"”":  `"`,
		"「":  "[",
		"」":  "]",
		"『":  "[",
		"』":  "]",
		"（":  "(",
		"）":  ")",
		"〔":  "[",
		"〕":  "]",
		"【":  "[",
		"】":  "]",
		"{":  "{",
		"}":  "}",
		"……": "...",
		"——": "-",
		"—":  "-",
		"～":  "~",
		"《":  "<",
		"》":  ">",
		"〈":  "<",
		"〉":  ">",
		"、":  ",",
	}
	// replacements 声调对应
	replacements = map[string][]string{
		"üē": {"ue", "1"},
		"üé": {"ue", "2"},
		"üě": {"ue", "3"},
		"üè": {"ue", "4"},
		"ā":  {"a", "1"},
		"ē":  {"e", "1"},
		"ī":  {"i", "1"},
		"ō":  {"o", "1"},
		"ū":  {"u", "1"},
		"ǖ":  {"v", "1"},
		"á":  {"a", "2"},
		"é":  {"e", "2"},
		"í":  {"i", "2"},
		"ó":  {"o", "2"},
		"ú":  {"u", "2"},
		"ǘ":  {"v", "2"},
		"ǎ":  {"a", "3"},
		"ě":  {"e", "3"},
		"ǐ":  {"i", "3"},
		"ǒ":  {"o", "3"},
		"ǔ":  {"u", "3"},
		"ǚ":  {"v", "3"},
		"à":  {"a", "4"},
		"è":  {"e", "4"},
		"ì":  {"i", "4"},
		"ò":  {"o", "4"},
		"ù":  {"u", "4"},
		"ǜ":  {"v", "4"},
	}
)

// ConvertResult 转换后字符串
type ConvertResult string

// 字典
type (
	dictDir     [6]string
	surNamesDir [1]string
)

// Config request conifg.
type Config struct {
	Dict      dictDir
	Surnames  surNamesDir
	Delimiter string
}

// InitConfig 初始化配置
var InitConfig Config

func init() {
	var dictDir = dictDir{
		os.Getenv("GOPATH") + "/src/github.com/gurudin/pinyin/dict/words_0.dict",
		os.Getenv("GOPATH") + "/src/github.com/gurudin/pinyin/dict/words_1.dict",
		os.Getenv("GOPATH") + "/src/github.com/gurudin/pinyin/dict/words_2.dict",
		os.Getenv("GOPATH") + "/src/github.com/gurudin/pinyin/dict/words_3.dict",
		os.Getenv("GOPATH") + "/src/github.com/gurudin/pinyin/dict/words_4.dict",
		os.Getenv("GOPATH") + "/src/github.com/gurudin/pinyin/dict/words_5.dict",
	}
	var surNamesDir = surNamesDir{
		os.Getenv("GOPATH") + "/src/github.com/gurudin/pinyin/dict/surnames.dict",
	}

	InitConfig = Config{
		Dict:     dictDir,
		Surnames: surNamesDir,
	}
}

// Convert 字符串转换拼音.
// strs: 转换字符串
func Convert(strs string) []string {
	result := tran(strs, false)

	return result.None()
}

// UnicodeConvert 字符串转换拼音.
// strs: 转换字符串
func UnicodeConvert(strs string) []string {
	result := tran(strs, false)

	return result.Unicode()
}

// ASCIIConvert 字符串转换拼音.
// strs: 转换字符串
func ASCIIConvert(strs string) []string {
	result := tran(strs, false)

	return result.ASCII()
}

// Name 翻译姓名
func Name(strs string) *ConvertResult {
	return tran(strs, true)
}

// Permalink 生成带分隔符的拼音字符串
func Permalink(strs, delimiter string) string {
	InitConfig.Delimiter = delimiter
	resStrs := tran(strs, false).None()

	return tranDelimiter(resStrs)
}

// Abbr 获取首字母带分隔符的拼音字符串
func Abbr(strs, delimiter string) string {
	InitConfig.Delimiter = delimiter
	resStrArr := tran(strs, false).None()

	var result []string
	for _, v := range resStrArr {
		result = append(result, string([]byte(v)[:1]))
	}

	return strings.Join(result, InitConfig.Delimiter)
}

func tranDelimiter(split []string) string {
	s := strings.Join(split, InitConfig.Delimiter)

	return s
}

func tran(strs string, surnames bool) *ConvertResult {
	s := InitConfig.romanize(strs, surnames)
	cr := ConvertResult(s)

	return &cr
}

// None 不带声调输出
// output: [pin, yin]
func (r *ConvertResult) None() []string {
	s := string(*r)

	for key, value := range replacements {
		s = strings.Replace(s, key, value[0], -1)
	}

	return strings.Split(s, " ")
}

// Unicode 输出
// output: [pīn, yīn]
func (r *ConvertResult) Unicode() []string {
	return strings.Split(string(*r), " ")
}

// ASCII 输出
// output: [pin1, yin1]
func (r *ConvertResult) ASCII() []string {
	s := string(*r)
	split := strings.Split(s, " ")

	for key, value := range replacements {
		for i := 0; i < len(split); i++ {
			tmpRep := strings.Replace(split[i], key, value[0], -1)
			if split[i] != tmpRep {
				split[i] = tmpRep + value[1]
			}
		}
	}

	return split
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

	s = c.Punctuations(s)

	s = strings.TrimSpace(s)
	s = strings.Replace(strings.Replace(s, "  ", " ", -1), "\t", " ", -1)

	return s
}

// Punctuations 转换标点符号
func (c *Config) Punctuations(s string) string {
	for k, v := range punctuations {
		s = strings.Replace(s, k, " "+v, -1)
	}

	return s
}

// 转换为字符串数组
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
