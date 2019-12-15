import * as sapper from '@sapper/server';

import compression from 'compression';
import dotenv from 'dotenv';
import polka from 'polka';

// Attempt to load .env from the monorepo root, since the dev server gets run
// from the web-app directory. If this file doesn't exist, we're running inside
// Docker and have actual env vars already loaded.
dotenv.config({path: '../.env'})

polka()
  .use(
    compression({ threshold: 0 }),
    sapper.middleware({
      session: (req, res) => ({
        config: {
          PUBGOLF_ENV: process.env.PUBGOLF_ENV,
          API_HOST_EXTERNAL: process.env.$API_HOST_EXTERNAL,
        }
      }),
    }),
).listen(process.env.PORT, err => {
  if (err) console.log('error', err);
});
