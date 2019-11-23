// FIXME: these should really come from a config file eventually
const HOST_NAMES = {
  dev: 'http://127.0.0.1:8080',
  prod: 'https://api.pubgolf.co',
  staging: 'https://api-staging.pubgolf.co',
};

export function getHost () {
  return HOST_NAMES[process.env.PUBGOLF_ENV] || HOST_NAMES.prod;
}

export async function get (req, res) {
  res.setHeader('Content-Type', 'application/json');
  res.end(JSON.stringify({ host: getHost() }));
}
