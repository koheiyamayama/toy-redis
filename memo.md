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