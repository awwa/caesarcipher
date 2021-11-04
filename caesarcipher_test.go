package main

import (
	"reflect"
	"testing"
)

// 暗号化処理のテストケース構造体
type encryptTest struct {
	in  string
	sh  int
	enc string
	err error
}

// 暗号化処理のテストケース
var encrypttests = []encryptTest{
	{"z", -1, "", ErrShift},
	{"ABC", 4, "", ErrChar},   // 非対象文字テスト
	{"あいうえお", 5, "", ErrChar}, // 非ASCII文字テスト
	{"abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghija", 31, "", ErrLength}, // 81文字の境界テスト
	{"123456789012345678901234567890123456789012345678901234567890123456789012345678901", 31, "", ErrChar},   // 不正な文字種＋文字数制限オーバー
	{"a", 0, "a", nil},
	{"a", 1, "b", nil},
	{"a", 25, "z", nil},
	{"a", 26, "a", nil},
	{"a", 27, "b", nil},
	{"z", 1, "a", nil},
	{"z", 26, "z", nil},
	{"z", 52, "z", nil},
	{"a", 34, "i", nil},
	{"hoge", 3, "krjh", nil},
	{"this is a pen", 8, "bpqa qa i xmv", nil},
	{"this is a pen.\r\n", 61, "cqrb rb j ynw.\r\n", nil},
	{"is this a pen", 8, "qa bpqa i xmv", nil},       // 先頭にthisを含むテスト
	{"pen is this", 8, "xmv qa bpqa", nil},           // 末尾にthisを含むテスト
	{"the pen is yours", 8, "bpm xmv qa gwcza", nil}, // 先頭にtheを含むテスト
	{"yours is pen the", 8, "gwcza qa xmv bpm", nil}, // 末尾にtheを含むテスト
	{"that is a pen", 8, "bpib qa i xmv", nil},       // 先頭にthatを含むテスト
	{"pen is that", 8, "xmv qa bpib", nil},           // 末尾にthatを含むテスト
	{
		"abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzab",
		8,
		"ijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghij",
		nil,
	}, // 80文字の境界テスト
}

// 暗号化処理テスト
func TestEncrypt(t *testing.T) {
	for i := range encrypttests {
		test := &encrypttests[i]
		actual, err := Encrypt(test.in, test.sh)
		if test.enc != actual {
			t.Errorf("Test failed: Encrypt('%s', %d) = '%s', %v want '%s', %v",
				test.in, test.sh, actual, test.err, test.enc, err)
		}
	}
}

// 復号化処理のテストケース構造体
type decryptTest struct {
	in  string
	dec string
	sh  int
	err error
}

// 復号化処理のテストケース
var decrypttests = []decryptTest{
	{"krjh", "", 0, ErrNoClue},
	{"uijt jt b qfo", "this is a pen", 8, nil},
	{"this is a pen", "this is a pen", 0, nil},
	{"bpqa qa i xmv", "this is a pen", 8, nil},
	{"cqrb rb j ynw.\r\n", "this is a pen.\r\n", 8, nil},
	{"", "", 0, ErrNoClue},                           // 非対称文字テスト
	{"あいうえお", "", 0, ErrChar},                        // 非ASCII文字テスト
	{"qa bpqa i xmv", "is this a pen", 8, nil},       // 先頭にthisを含むテスト
	{"xmv qa bpqa", "pen is this", 8, nil},           // 末尾にthisを含むテスト
	{"bpm xmv qa gwcza", "the pen is yours", 8, nil}, // 先頭にtheを含むテスト
	{"gwcza qa xmv bpm", "yours is pen the", 8, nil}, // 末尾にtheを含むテスト
	{"bpib qa i xmv", "that is a pen", 8, nil},       // 先頭にthatを含むテスト
	{"xmv qa bpib", "pen is that", 8, nil},           // 末尾にthatを含むテスト
	{"ijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghij", "", 8, ErrNoClue}, // 80文字の境界テスト
}

// 暗号化処理テスト
func TestDecrypt(t *testing.T) {
	for i := range decrypttests {
		test := &decrypttests[i]
		actual, sh, err := Decrypt(test.in)
		if test.dec != actual || test.err != err {
			t.Errorf("Test failed: Decrypt('%s') = '%s', %d, '%v' want '%s', %d, '%v'",
				test.in, actual, sh, err, test.dec, test.sh, test.err)
		}
	}
}

// 入力チェックのテストケース構造体
type assertTest struct {
	in  string
	err error
}

// 入力チェック処理のテストケース
// TableDrivenTests
// https://golang.org/src/strconv/atoi_test.go
var asserttests = []assertTest{
	{"a", nil},
	{"jldks kfjke klekl lskdj.\nhoge furga", nil},
	{"abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghij", nil}, // 80文字の境界テスト
	{"A", ErrChar},
	{"A\tH&3kdkwlkd", ErrChar},
	{"abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghija", ErrLength}, // 81文字の境界テスト
	{"123456789012345678901234567890123456789012345678901234567890123456789012345678901", ErrChar},   // 不正な文字種＋文字数制限オーバー
}

// 入力チェック処理テスト
func TestAssert(t *testing.T) {
	for i := range asserttests {
		test := &asserttests[i]
		actual := assert(test.in)
		if test.err != actual {
			t.Errorf("Test failed: Assert('%s') = %v want %v",
				test.in, actual, test.err)
		}
	}
}

// シフト処理のテストケース構造体
type shiftTest struct {
	in     string
	sh     int
	expect string
}

// シフト処理のテストケース
var shifttests = []shiftTest{
	{"a", 1, "b"},
	{"a", 8, "i"},
	{"a", 25, "z"},
	{"a", 26, "a"},
	{"a", 27, "b"},
	{"a", 51, "z"},
	{"a", 52, "a"},
	{"a", 53, "b"},
	{"a", -1, "z"},
	{"a", -8, "s"},
	{"a", -25, "b"},
	{"a", -26, "a"},
	{"a", -27, "z"},
	{"t", 1, "u"},
	{"t", 8, "b"},
	{"t", 25, "s"},
	{"t", 26, "t"},
	{"t", 27, "u"},
	{"t", 51, "s"},
	{"t", 52, "t"},
	{"t", 53, "u"},
	{"t", -1, "s"},
	{"t", -8, "l"},
	{"t", -25, "u"},
	{"t", -26, "t"},
	{"t", -27, "s"},
	{".", 10, "."},
	{" ", 10, " "},
	{"\r", 10, "\r"},
	{"\n", 10, "\n"},
}

// シフト処理のテストケース構造体
func TestShift(t *testing.T) {
	for i := range shifttests {
		test := &shifttests[i]
		actual := shift(test.in, test.sh)
		if actual != test.expect {
			t.Errorf("Test failed: Shift('%s', %d) = %v want %v",
				test.in, test.sh, actual, test.expect)
		}
	}
}

// インデクス差分計算処理テストケース
type subtractTest struct {
	minuend  rune
	divament rune
	sub      int
}

// インデクス差分計算処理のテストケース
var substracttests = []subtractTest{
	{'t', 'h', 12}, // t=19 h=7
	{'h', 'i', 25}, // h=7  i=8
	{'i', 's', 16}, // i=8  s=18
	{'b', 'p', 12}, // b=2  p=15
	{'p', 'q', 25}, // p=15 q=16
	{'q', 'a', 16}, // q=16 a=0
	{'t', 'h', 12}, // t=19 h=7
	{'t', 'e', 15}, // t=19 e=4
	{'b', 'p', 12}, // b=1 p=15
	{'b', 'm', 15}, // b=1 m=12
}

// インデクス差分計算処理のテスト
func TestSubtract(t *testing.T) {
	for i := range substracttests {
		test := &substracttests[i]
		actual := subtract(test.minuend, test.divament)
		if actual != test.sub {
			t.Errorf("Test failed: Subtract(%v, %v) = %d want %d",
				test.minuend, test.divament, actual, test.sub)
		}
	}
}

// 1文字単位のインデクス差分計算処理テストケース構造体
type substrTest struct {
	in    string
	subin []int
}

// 1文字単位のインデクス差分計算処理テストケース
var substrtests = []substrTest{
	{"this", []int{12, 25, 16}},
	{"the", []int{12, 3}},
	{"that", []int{12, 7, 7}},
}

// 1文字単位のインデクス差分計算処理テスト
func TestSubStr(t *testing.T) {
	for i := range substrtests {
		test := &substrtests[i]
		actual := subStr(test.in)
		if !reflect.DeepEqual(actual, test.subin) {
			t.Errorf("Test failed: subStr('%s') = %v want %v",
				test.in, actual, test.subin)
		}
	}
}

// 配列内で指定文字のインデクスを調べる処理テストケース構造体
type indexOfTest struct {
	target     []rune
	searchChar rune
	index      int
}

// 配列内で指定文字のインデクスを調べる処理テストケース
var indextests = []indexOfTest{
	{[]rune("abcde"), 'a', 0},
	{[]rune("abcde"), 'c', 2},
	{[]rune("abcde"), 'e', 4},
	{[]rune("abcde"), 'f', -1},
}

// 1文字単位のインデクス差分計算処理テスト
func TestIndex(t *testing.T) {
	for i := range indextests {
		test := &indextests[i]
		actual := indexOf(test.target, test.searchChar)
		if actual != test.index {
			t.Errorf("Test failed: index('%v', '%v') = %d want %d",
				test.target, test.searchChar, actual, test.index)
		}
	}
}

func BenchmarkAssert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		assert("xlmw mw xli tmgxyvi xlex m xsso mr xli xvmt.")
	}
}

func BenchmarkSubtract(b *testing.B) {
	for i := 0; i < b.N; i++ {
		subtract('x', 'j')
	}
}

func BenchmarkSubStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		subStr("xlmw mw xli tmgxyvi xlex m xsso mr xli xvmt.")
	}
}

func BenchmarkShift(b *testing.B) {
	for i := 0; i < b.N; i++ {
		shift("xlmw mw xli tmgxyvi xlex m xsso mr xli xvmt.", -8)
	}
}

func BenchmarkIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		indexOf([]rune("abcdefghijklmnopqrstuvwxyz"), 't')
	}
}

func BenchmarkDecrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Decrypt("xlmw mw xli tmgxyvi xlex m xsso mr xli xvmt.")
	}
}
