const COOKIE_NAME = 'user';
const COOKIE_AGE = 3600 * 24 * 3; // 3 days in seconds

export async function post (req, res) {
  // Parse the request body
  let rawBody = '';
  req.on('data', (chunk) => {
    rawBody += chunk;
  });
  req.on('end', async () => {
    res.cookies.set(COOKIE_NAME, rawBody, {
      maxAge: COOKIE_AGE,
      httpOnly: true, // Protect the cookie from javascript access
    });
    res.end('OK');
  });
}

export function get (req, res) {
  const user = req.cookies.get(COOKIE_NAME);
  if (user) {
    res.writeHead(200, {
      'Content-Type': 'application/json',
    }).end(user);
  } else {
    res.writeHead(401, 'NO AUTHENTICATION FOUND').end();
  }
}
