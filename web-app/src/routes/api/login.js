import Cookies from 'cookies';

import { loadEnv } from 'src/_server-utils';
import { getAPI } from 'src/api';

const COOKIE_NAME = 'user';
const COOKIE_AGE = 3600 * 24 * 3; // 3 days in seconds

export async function post (req, res) {
  const api = getAPI({
    config: loadEnv(),
    user: {},
  });
  const cookies = new Cookies(req, res);
  // Parse the request body
  let rawBody = '';
  req.on('data', (chunk) => {
    rawBody += chunk;
  });
  req.on('end', async () => {
    const loginDetails = JSON.parse(rawBody);

    const user = JSON.stringify(await api.playerLogin(loginDetails));

    cookies.set(COOKIE_NAME, user, {
      maxAge: COOKIE_AGE,
      path: '/api/login',
      httpOnly: true,
    });
    res.end(user);
  });
}

export function get (req, res) {
  const cookies = new Cookies(req, res);
  const user = cookies.get(COOKIE_NAME);
  // TODO: do something if the user doesn't exist
  console.log('Restore User:', user);
  res.end(user);
}
