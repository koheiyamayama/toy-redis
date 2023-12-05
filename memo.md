# docker-composeで作ってたダッシュボードが全て消えた
理由はわからない。たぶん、俺の操作ミス。
そろそろ
- 開発環境はログ垂れ流し + 1バイナリで開発する。
- 本番環境はGrafana + 1バイナリで動作確認する。

上記の環境を作っていくことにする。

## Grafana
master-nodeにGrafanaを入れる。
https://grafana.com/docs/grafana/latest/setup-grafana/installation/debian/

## Prometheus
worker-node14にPrometheusを入れる。
```
cat /etc/prometheus/prometheus.yml
# my global config
global:
  scrape_interval: 2s
  evaluation_interval: 30s

scrape_configs:
  - job_name: toy-redis
    metrics_path: /metrics
    honor_labels: false
    honor_timestamps: true
    scheme: http
    follow_redirects: true
    body_size_limit: 15MB
    sample_limit: 1500
    target_limit: 30
    label_limit: 30
    label_name_length_limit: 200
    label_value_length_limit: 200
    static_configs:
      - targets: ["192.168.1.12:8080"]
```

## Loki
worker-node-11にLokiを入れる

## toy-redis
worker-node-12にtoy-redisとpromtailを入れる
```
cat /etc/promtail/promtail-local-config.yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /var/log/positions.yaml # This location needs to be writeable by Promtail.

clients:
  - url: http://192.168.1.11:3100/loki/api/v1/push

scrape_configs:
  - job_name: toy-redis
    pipeline_stages:
    static_configs:
      - targets:
          - localhost
        labels:
          job: toy-redis # A `job` label is fairly standard in prometheus and useful for linking metrics and logs.
          __path__: /var/log/toy-redis.log # The path matching uses a third party library: https://github.com/bmatcuk/doublestar
```

手順はOpenAIの言う通りにしたら上手く行った。
https://chat.openai.com/share/f853d5ad-172c-4a46-b0b4-6d512d01b0ca

# expireを過ぎてもエントリ数が0にならない
gorutineとほぼ同じペースで減っていくが、ある一定のところからエントリ数が減らなくなる。
何が起きているのか調査が必要。

# 100vu負荷時の話
シナリオ上はSET, GETの順でリクエストしているが、サーバーログを見るとGET,SET担ってしまうケースが多発して色々とだめそう。
何が原因なのか調査する。

# Expireコマンドを実装する
toy-redisのexpireについて
entryがストアから削除される。その時間を設定できる。
entryがセットされてから削除されるまでの秒数。
expireは更新可能で、更新すると、更新後から削除されるまでの秒数を設定できる。

時間軸的にはこういう感じ。
---set---expired
---set---expire---expired
---set---expire---expire---...
---set(key1)---set(key1)---expired

setで同じキーを更新したときもexpireを更新する必要がある。

expireコマンドの仕様としてkeyがない場合どうするか？
- not found keyエラーを返す
- ~何もしない~

# setのexpireを実装する
- Setのインターフェイスを拡張する
- Setで受け取ったint秒後にキャッシュデータを削除する
  - https://pkg.go.dev/time#After おそらくtime.Afterを使うと良さそう

1エントリーに対して1gorutineを起動する。
gorutineはexpireが更新されるたびに、再起動されるイメージ。

# 負荷試験と監視を実装する
Grafana系サービスを使って負荷試験と監視をする。
負荷試験はとりあえず、疎通できたので、あとは監視システムを導入してモニタリングできるようにする。
- goruntine
- 各コマンドの処理回数とそれぞれの処理時間

prometheusの概念を理解していないからそのあたりやりつつ導入を進める。
registryやcollectorとか知らんわ。

## 負荷試験
- コマンドの実行時間
- 負荷試験結果

## toy-redisサーバー
- メトリクス
  - ~各コマンドの実行回数~
    - metrics type count
  - 各コマンドの実行時間の分布
    - 99%tile
    - 95%tile
    - 90%tile
    - 50%tile
  - ~エントリーの総数(時系列)~
    - gauge
  - ~エントリーの増加率~
    - count,rate
  - ロック待ちの時間
    - read,writeそれぞれあると良さそう？
    - 調べる
  - ~go runtime~
    - promauto
    - prometheusのセットアップは完了したので、grafanaで必要なメトリクスを可視化する
- ロギング
  - Grafana Dashboardを作成する
    - クエリを書く
  - toy-redisからPromtail経由でLokiにログを送りつける
    - toy-redisで./tmp/toy-redis.logファイルを作成
      - func NewJSONHandler(w io.Writer, opts *HandlerOptions) *JSONHandler の第一引数に作成したファイルを渡すだけで良さそう。あとはPromtailでログを加工してLokiに送りつけるだけっぽい。
    - Promtailで./tmp/toy-redisを設定
    - PromtailでLokiへログを送付する
    - GrafanaからLogQLを使ってビジュアライズ
- プロファイル
  - Grafana Pyroscope使ってみる

- 負荷試験
  - GET,SET,EXPIREを1.適度に流す 2.GETを極端に叩く 3.SETを極端に叩く 4. EXPIREを極端に流す の4パターンを用意し、それぞれの負荷に高低付ける。
  - キーのカーディナリティ、キーとバリューのサイズも幅を持たせる。

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
