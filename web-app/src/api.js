import { StatusCode } from 'grpc-web';

import {
  capFirst,
  mapEntries,
} from './utils';

const {
  CreateOrUpdateScoreRequest,
  GetScheduleRequest,
  GetScoresRequest,
  RegisterPlayerRequest,
  RequestPlayerLoginRequest,
  PlayerLoginRequest,
  League,
} = require('./proto/pubgolf_pb');
const {
  APIPromiseClient,
} = require('./proto/pubgolf_grpc_web_pb');


export const LEAGUE = League;

// Optional override messages for GRPC status codes
// https://github.com/grpc/grpc/blob/master/doc/statuscodes.md
const MESSAGES = {
  [StatusCode.UNKNOWN]: 'Unexpected server error',
  [StatusCode.INVALID_ARGUMENT]: 'Invalid input',
};


/**
 * @param promise
 *
 * @returns {Object}
 */
function _unWrap (promise) {
  return promise.then(
    instance => instance.toObject(),
    (error) => {
      console.error(error);

      // Replace the top-level message with something user-presentable.
      // The original message is still preserved in error.metadata
      if (error.code && error.code in MESSAGES) {
        error.message = MESSAGES[error.code];
      }
      throw error;
    },
  );
}

/**
 * Generate a wrapper for a given gRPC method
 * @param {APIPromiseClient} client - The gRPC client that will be used to
 *                                    make the requests
 * @param {Object} metadata - Metadata to add to each request
 * @param {string} methodName - The name of the method to call on `client`
 * @param {Class} RequestClass - The class of request that that method uses
 *
 * @returns {function(Object): Object}
 *           - A function that takes in a plain object of the method's request
 *             type and returns a plain object of the methods return type
 */
function rpcMethod (client, metadata, methodName, RequestClass) {
  return (params) => {
    const request = new RequestClass();
    // Populate the request instance from the given params
    Object.entries(params).forEach(([key, value]) => {
      request[`set${capFirst(key)}`](value);
    });

    // Make the request then turn the response into a plain object
    return _unWrap(client[methodName](request, metadata));
  };
}

/**
 * Build an object of a bunch of gRPC methods
 * @param {APIPromiseClient} client - The gRPC client that will be used to
 *                                    make the requests
 * @param {Object} metadata - Metadata to add to each request
 * @param {Object.<string,Class>} methods - a map of method names to classes
 *         TODO: This was the simplest I could get the config part...
 *
 * @returns {Object.<string, function(Object): Object>}
 */
function buildMethods (client, metadata, methods) {
  return mapEntries(methods, ([methodName, RequestClass]) => (
    [methodName, rpcMethod(client, metadata, methodName, RequestClass)]
  ));
}

/**
 * @typedef {Object} APIWrapper - Expected inputs and outputs for each api method
 *
 * @property {function({
 *   eventKey: string,
 *   name: string,
 *   phoneNumber: string,
 *   league: string,
 * }): Promise<void>} registerPlayer - Create a player with the given data
 *
 * @property {function({
 *   eventKey: string,
 *   phoneNumber: string,
 * }): Promise<void>} requestPlayerLogin - Start authenticating as given player
 *
 * @property {function({
 *   eventKey: string,
 *   phoneNumber: string,
 *   authCode: string,
 * }): Promise<PlayerLoginReply.AsObject>} playerLogin - Log in as the given player
 *
 * @property {function({
 *   eventKey: string,
 * }): Promise<GetScheduleReply.AsObject>} getSchedule - Get the schedule of stops for an event
 *
 * @property {function({
 *   eventKey: string,
 * }): Promise<GetScoresReply.AsObject>} getScores - Get groups of current scores for the event
 *
 * @property {function({
 *   venueId: string,
 *   playerId: string,
 *   strokes: number,
 * }): Promise<void>} createOrUpdateScore - Submit a score for approval
 */

/**
 * Build the API wrapper
 * @param {string} host - The hostname of the gRPC server
 * @param {Object} metadata - gRPC metadata
 *
 * @returns {APIWrapper}
 */
function buildAPIWrapper (host, metadata = {}) {
  const client = new APIPromiseClient(host);
  return buildMethods(client, metadata, {
    registerPlayer: RegisterPlayerRequest,
    requestPlayerLogin: RequestPlayerLoginRequest,
    playerLogin: PlayerLoginRequest,
    getSchedule: GetScheduleRequest,
    getScores: GetScoresRequest,
    createOrUpdateScore: CreateOrUpdateScoreRequest,
  });
}

const CACHE = new Map();

/**
 *
 * @param {Object} session
 *
 * @returns {APIWrapper}
 */
export function getAPI (session) {
  const {
    config: { API_HOST_EXTERNAL: host },
    user: { authToken },
  } = session;
  const cacheKey = `${host} ${authToken}`;

  if (CACHE.has(cacheKey)) {
    return CACHE.get(cacheKey);
  }
  CACHE.set(cacheKey, buildAPIWrapper(host, { authorization: authToken }));

  return CACHE.get(cacheKey);
}
