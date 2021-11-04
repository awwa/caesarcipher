# caesarcipher
シーザー暗号化、復号化する関数

## 実装方針

- マジックナンバーの利用を避ける
  - 計算処理がわかりづらくなる
- 計算とログ出力は分離する
  - Unit testするため
  - プログラムの再利用性を高めるため
- エラーは関数の戻り値でerrを返すことで表現
  - [Golangの標準的なエラー処理戦略](https://golang.org/doc/effective_go)に従う
- 第三者が見て理解できるコメントをつける
  - 誰でもメンテナンスができるようにするため

## 実行環境

Node.js `v1.17.2` にて動作確認

## 準備

```
$ git clone git@github.com:awwa/caesarcipher.git
$ cd caesarcipher
```

## テスト

```
$ go test

PASS
ok      caesarcipher    0.493s
```

## ベンチマーク

shift()とSubStr()を改善して処理速度を上げてみた(TODO箇所参照)

```
$ go test --bench .
goos: darwin
goarch: amd64
pkg: caesarcipher
cpu: Intel(R) Core(TM) i5-6360U CPU @ 2.00GHz
BenchmarkAssert-4         364064              3289 ns/op
BenchmarkSubtract-4     44459637                23.23 ns/op
BenchmarkSubStr-4          73123             14851 ns/op
BenchmarkShift-4          312428              4130 ns/op
BenchmarkIndex-4        75046777                16.08 ns/op
BenchmarkDecrypt-4         27697             38651 ns/op
PASS
ok      caesarcipher    10.599s
```

## 実行

```
$ go run .
in : 'xlmw mw xli tmgxyvi xlex m xsso mr xli xvmt.'
out: 'this is the picture that i took in the trip.'
sh : 4
err: <nil>
```
