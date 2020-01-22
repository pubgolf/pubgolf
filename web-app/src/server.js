import * as sapper from '@sapper/server';

import compression from 'compression';
import polka from 'polka';

import { loadEnv } from './_server-utils';
import { getAPI } from './api';


polka()
  .use(
    compression({ threshold: 0 }),
    sapper.middleware({
      session (req, res) {
        const api = getAPI({
          config: loadEnv(),
          user: {
          },
        });
        api.playerLogin({});
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
