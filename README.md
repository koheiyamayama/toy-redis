# 目的
- golangを用いておもちゃRedisを作成する

# 目標
- Redisっぽいサーバーを作成する
- クライアントも作成する

# やること
- 以下のコマンドを実装する
  - GET
  - SET
  - DEL
  - EXISTS
  - EXPIRE
- 以下の機能を実装する
  - レプリケーション

# セットアップ
```
$ docker compose up --build
$ cd loadtesting && npm run bundle && cd -
$ ruby loadtesting/generate_testdata.rb
$ make build-k6
$ make run-k6
```

# 環境変数
- TOY_REDIS_LOG_LEVEL
  - log/slogパッケージのログレベルに依存する
  - LevelDebug Level = -4
 	- LevelInfo  Level = 0
	- LevelWarn  Level = 4
	- LevelError Level = 8
