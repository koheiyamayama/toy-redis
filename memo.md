# expireを実装する
- Setのインターフェイスを拡張する
- Setで受け取ったint秒後にキャッシュデータを削除する
  - https://pkg.go.dev/time#After おそらくtime.Afterを使うと良さそう

1エントリーに対して1gorutineを起動する。
gorutineはexpireが更新されるたびに、再起動されるイメージ。

# 負荷試験と監視を実装する
Grafana系サービスを使って負荷試験と監視をする。
負荷試験はとりあえず、疎通できたので、あとは監視システムを導入してモニタリングできるようにする。
- goruntime
- 各コマンドの処理回数とそれぞれの処理時間

prometheusの概念を理解していないからそのあたりやりつつ導入を進める。
registryやcollectorとか知らんわ。

# コネクションを使い回せるようにしたい
c4fd100
今の実装だとSETした後に一度コネクションをクライアント側でcloseして、再度つなぎ直してGETコマンドの実行とコネクションからレスポンスの読み込みを行っている。

本当はSETの後のコネクションcloseをせずにコネクションを使い回せるようにしたい。

今はなんかSET,GETコマンドを書き込んだ後にレスポンスが表示されない。原因そうなのは
- クライアント
  - conn.Readがio待ち
- サーバー
  - conn.Writeがio待ち

のどちらかっぽさそう。

って思ってたが、ロジックが悪かった。
クライアント側でconn.Writeを2回実行していると、
```
000100000SEThogehogehogehoge\rfugafuga2times\n000100000GEThogehogehogehoge\n
```
っていうストリームができる。
整形するとこうなる。
```
SET hogehogehogehoge fugafuga2times
GET hogehogehogehoge
```

処理を見ると、こんな感じになっている。
```golang
		r := bufio.NewReader(conn)
		b, err := r.ReadBytes('\n')
		b = b[:len(b)-1]
```
こうすると、ストリームからSETコマンドしか読み込まない。
無限ループしてもGETコマンドを読み込まないので、どうもReadBytesは全て読みだしてから区切っている？
このあたりどうするか考えないといけない。
コネクションに書き込まれたコマンドを全て読みだして実行するっていう感じにしないといけない。

# io.ReadAllの処理が終わらない
0c401c496fccb9640de17695fb862f474044caad で./cmd/client/main.goを実行すると、main.handleConn内のio.ReadAllの処理が終わらない。
これを調査する。

まず、io.ReadAllの中身はこうなってた。
```golang
package io
// ReadAll reads from r until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because ReadAll is
// defined to read from src until EOF, it does not treat an EOF from Read
// as an error to be reported.
func ReadAll(r Reader) ([]byte, error) {
  // 512byteのスライスを確保する
	b := make([]byte, 0, 512)
	for {
    // len(b)はスライスの最後の要素のindexを表す
    // cap(b)はスライスの確保領域の最後のindexを返す
    // つまり、スライス内の空き領域にバイトを読み込もうとしている
		n, err := r.Read(b[len(b):cap(b)])

    // nはReaderから読み込んだバイト数
    // len(b)はスライスの最後の要素のindexを表す
    // それまでに読み込んだバイトと新しく読み込んだバイトを合体する?
		b = b[:len(b)+n]
		if err != nil {
			if err == EOF {
				err = nil
			}
      // EOFの場合は全てを読み込んだということなので、読み込んだバイトを返す
      // EOF以外の場合は読み込めたバイトとエラーを返す
			return b, err
		}

    // スライスが要素で満たされた場合、キャパシティを追加している
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
	}
}
```

どうもr.ReadでEOFが返るまで読み込むっぽい。
つまり、EOFが返ってないから処理が無限ループしているっぽい。
net.ConnがEOFを返すのはcloseされたコネクションに対して操作しようとした時っぽい。
https://castaneai.hatenablog.com/entry/2020/01/09/193539

じゃあ、どうするかっていうとそもそもbufioを使う実装が多そう？
https://github.com/venilnoronha/tcp-echo-server/blob/master/main.go#L45
