# golang-pinyin

Go语言 中文转拼音类库

### 安装
```
go get github.com/gurudin/pinyin
```

## 使用
---

### 引入
```
import "github.com/gurudin/pinyin"
```

### 转换数组
```
/* 转换拼音 */
pinyin.Convert("go语言汉字转换拼音")
// -> [go yu yan han zi zhuan huan pin yin]

/* 转换拼音有声调*/
pinyin.UnicodeConvert("go语言汉字转换拼音")
// -> [go yǔ yán hàn zì zhuǎn huàn pīn yīn]

/* 转换拼音数字声调 */
pinyin.ASCIIConvert("go语言汉字转换拼音")
// -> [go yu3 yan2 han4 zi4 zhuan3 huan4 pin1 yin1]
```

### 姓名翻译
```
/* 转换拼音 */
pinyin.Name("冒顿单于").None()
// -> [mo du chan yu]

/* 转换拼音有声调 */
pinyin.Name("冒顿单于").Unicode()
// -> [mò dú chán yú]

/* 转换拼音数字声调 */
pinyin.Name("冒顿单于").ASCII()
// -> [mo4 du2 chan2 yu2]
```

### 生成带分隔符的拼音字符串
```
pinyin.Permalink("go语言汉字转换拼音", "-")
// -> go-yu-yan-han-zi-zhuan-huan-pin-yin
```