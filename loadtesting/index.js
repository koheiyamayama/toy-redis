import tcp from 'k6/x/tcp';
import { check, sleep } from 'k6';
import { sha256 } from 'js-sha256';


export const options = {
  vus: 100,
  duration: '5m',
};

const key_salt = 'toy-redis'
const key_cardinality = 100000

export default function () {
  let conn = tcp.connect('localhost:9999');
  tcp.write(conn, '000100000SEThoge\rfuga\r100\n');
  tcp.close(conn);

  conn = tcp.connect('localhost:9999');
  tcp.write(conn, '000100000GEThoge\n');
  tcp.close(conn);

  conn = tcp.connect('localhost:9999');
  tcp.write(conn, '000100000GEThoge\n');

  let res = String.fromCharCode(...tcp.read(conn, 1024))
  check (res, {
    'verify value': (res) => res.includes('fuga')
  });
  tcp.close(conn);
  sleep(1)
}

function getTestData() {
  // id = Sha256(key_salt + rand(0..key_cardinality))
  // select key, value from testdata where key = id;
  // return key, value
  const hash = sha256.create()
  hash.update(`${key_salt}${randomIntFromInterval(1, key_cardinality)}`)
  const key = hash.hex()

  console.log(key)
}

function randomIntFromInterval(min, max) { // min and max included 
  return Math.floor(Math.random() * (max - min + 1) + min)
}
