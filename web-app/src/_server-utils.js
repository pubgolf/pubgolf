import dotenv from 'dotenv';


/**
 * return relevant items from the .env file
 * @returns {{PUBGOLF_ENV: string, API_HOST_EXTERNAL: string}}
 */
export function loadEnv () {
  // Attempt to load .env from the monorepo root, since the dev server gets run
  // from the web-app directory. If this file doesn't exist, we're running inside
  // Docker and have actual env vars already loaded.
  dotenv.config({
    path: '../.env',
  });
  return {
    PUBGOLF_ENV: process.env.PUBGOLF_ENV,
    API_HOST_EXTERNAL: process.env.API_HOST_EXTERNAL,
  };
}
