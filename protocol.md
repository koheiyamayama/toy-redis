# Structure
- version 4byte
- command 8byte
- value 1mb


# GET
キーに対する値を取得する。

## command
0001,00000GET,key

## result
ok,type,value

# SET
キーに対して値をセットする。
expにセットされた秒数後に削除される。
expはuin32型。

## command
- protocol
  - 000100000SETkey\rvalue\rexpire\n
- cli
  - set key value exp

## result
ok
error

# EXPIRE
キーがセットされたエントリーのexpを更新する。

## command
- protocol
  - 000100EXPIREkey\rexpire\n
- cli
  - expire key exp
## result
ok
