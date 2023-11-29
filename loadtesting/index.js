import tcp from 'k6/x/tcp';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';
import randomInteger from 'random-int';

export const options = {
  vus: 30,
  duration: '5m',
};

const data = new SharedArray("testdata", function() {
  return JSON.parse(open("../testdata.json")).data
})

export default function () {
  const kv = data[randomInteger(0, data.length-1)]

  let conn = tcp.connect('localhost:9999');
  tcp.write(conn, Set(kv.key, kv.value, randomInteger(300, 360)));
  tcp.close(conn);

  sleep(1)

  conn = tcp.connect('localhost:9999');
  tcp.write(conn, Get(kv.key));

  // toy-redisのレスポンスにはprefixとしてデータ型を表す記号が1文字付与されるので、テストデータの長さ + 記号分の長さ(=1)分をtcpストリームから読み込む必要がある。
  let res = String.fromCharCode(...tcp.read(conn, kv.value.length+1))
  check(res, {
    'verify value': (res) => res.includes(kv.value),
    'not found key': (res) => res.includes("not found key")
  });
  tcp.close(conn);
  sleep(1)
}

function Set(key, value, exp) {
  return `000100000SET${key}\r${value}\r${exp}\n`
}

function Get(key) {
  return `000100000GET${key}\n`
}
