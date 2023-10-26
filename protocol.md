# Structure
- version 4byte
- command 8byte
- value 1mb


# GET
## command
0001,00000GET,key

## result
ok,type,value

# SET
## command
可変長な値の区切り文字どうする？
今のままだとkeyとvalueの区切りができない？
0001,00000SET,len(key)key,len(value)value,expire

## result
ok
error

