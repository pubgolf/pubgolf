import * as sapper from '@sapper/server';

import Cookies from 'cookies';
import compression from 'compression';
import polka from 'polka';

import { loadEnv } from './_server-utils';


polka()
  .use(
    Cookies.express(),
    compression({ threshold: 0 }),
    sapper.middleware({
      session (req, res) {
        const userCookie = req.cookies.get('user');
        return {
          config: loadEnv(),
          user: userCookie ? JSON.parse(userCookie) : {
          },
        };
      },
    }),
  ).listen(process.env.PORT, (err) => {
    if (err) console.log('error', err);
  });
