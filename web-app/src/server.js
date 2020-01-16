import * as sapper from '@sapper/server';

import compression from 'compression';
import polka from 'polka';

import { loadEnv } from './_server-utils';


polka()
  .use(
    compression({ threshold: 0 }),
    sapper.middleware({
      session (req, res) {
        return {
          config: loadEnv(),
          user: {
          },
        };
      },
    }),
  ).listen(process.env.PORT, (err) => {
    if (err) console.log('error', err);
  });
