# io.ReadAllの処理が終わらない
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
