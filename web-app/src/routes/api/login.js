import crypto from 'crypto';

import Cookies from 'cookies';

// FIXME: as-is, this will lead to a memory leak,
// FIXME: need to add a mechanism to expire entries
const cache = new Map();

export async function post (req, res) {
  const cookies = new Cookies(req, res);
  // Parse the request body
  let rawBody = '';
  req.on('data', (chunk) => {
    rawBody += chunk;
  });
  req.on('end', async () => {
    const user = JSON.parse(rawBody);
    const buf = crypto.randomBytes(256);
    const key = buf.toString('hex');

    cache.set(key, user);

    cookies.set('user', key, {
      maxAge: 3600,
      path: '/api/login',
      httpOnly: true,
      // overwrite: true,
    });
    res.end('OK');
  });
}

export function get (req, res) {
  const cookies = new Cookies(req, res);
  const key = cookies.get('sessionId');
  console.log(key, cache);
  res.end(cache.get(key));
}
