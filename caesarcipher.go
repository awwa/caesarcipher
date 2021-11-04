// シーザー暗号化、復号化
package main

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

// 文字列変換テーブル
var TBL = []rune("abcdefghijklmnopqrstuvwxyz")

// 復号化の手がかりワード
var CLUES = []string{"this", "the", "that"}

// 入力エラー
var (
	// 長さエラー
	ErrLength = errors.New("invalid length")
	// 文字種エラー
	ErrChar = errors.New("invalid char")
	// 手がかりを含まないエラー
	ErrNoClue = errors.New("no clue word")
	// 不正なシフト量
	ErrShift = errors.New("invalid shift value")
)

// Encryptはシーザー暗号化を行う
// 概要:
//   入力値の各文字列を変換テーブル上で右方向にシフトして結果を返す
// 入力:
//   in: 平文の文字列（a-z .\r\nのみ入力可、80文字以内）
//   sh: 右方向へのシフト量（正値のみ入力可）
// 出力:
//   enc: 暗号化文字
//   err: 正常時: nil, エラー: 不正な入力
func Encrypt(in string, sh int) (enc string, err error) {
	// 入力チェック
	err = assert(in)
	if sh < 0 {
		err = ErrShift
	}
	if err != nil {
		return
	}
	// 右方向文字シフト
	enc = shift(in, sh)
	return
}

// Decryptはシーザー復号化を行う
// 概要:
//   入力値の各文字列を変換テーブル上で左方向にシフトして結果を返す
// 入力:
//   in: 暗号化文字列(a-z .\r\nのみ入力可、80文字以内、元の文字に手がかりワードthis,the,thatいずれかを含む)
// 出力:
//   dec: 平文
//   sh:  入力に適用されていた右方向へのシフト量
//   err: 正常時: nil, エラー: 不正な入力
// 戦略:
//   暗号化されていても各文字間のインデクスの差は変わらないので、入力文字列からこの条件にマッチする箇所を探し出し、シフト量を算出する
// 例:
//   "this": []int{12, 25, 16}
//   "uijt": []int{12, 25, 16}
func Decrypt(in string) (dec string, sh int, err error) {
	// 入力チェック
	err = assert(in)
	if err != nil {
		return
	}
	var hit bool = false // 入力文字列中に手がかりワードが見つかったことを示すフラグ
	// 1. 入力文字列にて、ある位置の文字とその隣の文字のインデクス値の差を算出し配列を生成する(subin)
	subin := subStr(in)
	// 手がかりワードごとにループ
	for i := 0; i < len(CLUES); i++ {
		// 2. 手がかりワードにて、ある位置の文字とその隣の文字のインデクス値の差を計算し配列を生成する(subclue)
		subclue := subStr(CLUES[i])
		// 3. subin内にsubclueに一致する順序の値をサーチする
		for j := 0; j < len(subin)-len(subclue)+1; j++ {
			// 4. 見つかったらその位置の文字と手がかりワードの先頭文字のインデクス値の差がシフト量(sh)
			if reflect.DeepEqual(subin[j:j+1], subclue[0:len(subclue)-1]) {
				sh = subtract([]rune(in)[j], []rune(CLUES[i])[0])
				hit = true
				break
			}
		}
	}
	// 手がかりワードが見つからないエラー
	if !hit {
		err = ErrNoClue
		return
	}
	// 5. シフト量分左にシフトすれば復号化できる
	dec = shift(in, -sh)
	return
}

// assertは入力文字列の共通チェック処理
// 入力:
//   チェック対象文字列
// 出力:
//   err: 正常時: nil, エラー: 不正な入力
// 詳細:
// - 80文字以内
// - アルファベット小文字
// - ピリオド、半角スペース、改行
// https://regexper.com/#%5Ba-z%5C.%20%5Cr%5Cn%5D%7B1%2C80%7D
func assert(in string) (err error) {
	if regexp.MustCompile(`[^a-z\. \r\n]`).MatchString(in) {
		err = ErrChar
	} else if len(in) > 80 {
		err = ErrLength
	}
	return
}

// shiftは入力値の各文字列を変換テーブル上でシフトして結果を返す
// 入力:
//   in: 文字列（内部向け関数のため、a-z .\r\nのみ入力される想定でエラー処理なし）
//   sh: シフト量（正値:右方向シフト、負値:左方向シフト）
// 出力:
//   out: シフト後の文字列
func shift(in string, sh int) (out string) {
	for _, v := range in {
		// 特定記号の場合、変換せずにそのまま返す
		// if regexp.MustCompile(`[\. \r\n]`).MatchString(string(v)) { // 遅い TODO 削除
		if v == '.' || v == ' ' || v == '\r' || v == '\n' { // 速い
			out += string(v)
			continue
		}
		// 変換テーブル内のインデクス値
		i := indexOf(TBL, v)
		// シフト後に変換テーブルの範囲を超えたらループして調整
		len := len(TBL)
		var ii int = (i + sh) % len
		if ii < 0 {
			ii += len
		}
		if ii > len {
			ii -= len
		}
		out += string(TBL[ii])
	}
	return
}

// subtractは1文字単位で変換テーブル内のインデクス差分計算を行う
// 入力:
//   left: 左側文字列
//   right: 右側文字列
// 出力:
//   out: インデクス差分(被減数<減数の場合、変換テーブル長を加えてループする)
func subtract(left rune, right rune) (out int) {
	l := indexOf(TBL, left)
	r := indexOf(TBL, right)
	out = l - r
	if out < 0 {
		out += len(TBL)
	}
	return
}

// subStrは文字列内の各文字で変換テーブル内のインデクス差分を計算する
// 入力:
//   in: 対象文字列
// 出力:
//   subi: インデクス差分値の配列
// func subStr(in string) (subin []int) {	// 遅い TODO 削除
func subStr(in string) []int { // ほんのり速い
	subin := make([]int, 0, 79) // ほんのり速い
	for i := range in {
		if i > len(in)-2 {
			break
		}
		subin = append(subin, subtract([]rune(in)[i], []rune(in)[i+1]))
	}
	// return
	return subin
}

// indexは指定した文字が最初に見つかった配列内のインデクスを返す
// 入力:
//   target: rune配列
//   searchChar: 検索対象文字
// 出力:
//   見つかった場合: 配列内インデクス、見つからなかった場合: -1
func indexOf(target []rune, searchChar rune) int {
	for i, v := range target {
		if v == searchChar {
			return i
		}
	}
	return -1
}

func main() {
	in := "xlmw mw xli tmgxyvi xlex m xsso mr xli xvmt."
	// in := "適当な文字列"
	fmt.Printf("in : '%s'\n", in)
	out, sh, err := Decrypt(in)
	fmt.Printf("out: '%s'\n", out)
	fmt.Printf("sh : %d\n", sh)
	fmt.Printf("err: %v\n", err)
}
