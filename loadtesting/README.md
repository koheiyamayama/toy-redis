[k6](https://k6.io/)と拡張である[xk6-tcp](https://github.com/NAlexandrov/xk6-tcp)を利用する

```
# k6バイナリを拡張込みでbuildする
$ xk6 build master \
  --with github.com/NAlexandrov/xk6-tcp
# 負荷シナリオを実施する  
$ k6 run index.js
```
