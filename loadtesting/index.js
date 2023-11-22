import tcp from 'k6/x/tcp';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';
import randomInteger from 'random-int';

export const options = {
  vus: 1,
  duration: '5m',
};

const data = new SharedArray("testdata", function() {
  return JSON.parse(open("../testdata.json")).data
})

export default function () {
  const kv = data[randomInteger(0, data.length-1)]

  let conn = tcp.connect('localhost:9999');
  console.log("set: ",Set(kv.key, kv.value, randomInteger(300, 360)))
  tcp.write(conn, Set(kv.key, kv.value, randomInteger(300, 360)));
  tcp.close(conn);

  sleep(1)

  conn = tcp.connect('localhost:9999');
  console.log("get: ", Get(kv.key))
  tcp.write(conn, Get(kv.key));

  let res = String.fromCharCode(...tcp.read(conn, kv.value.length))
  check (res, {
    'verify value': (res) => res.includes(kv.value)
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
