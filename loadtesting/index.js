import tcp from 'k6/x/tcp';
import { check, sleep } from 'k6';

export const options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  let conn = tcp.connect('localhost:9999');
  tcp.write(conn, '000100000SEThoge\rfuga\n');
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
