import { StatusCode } from 'grpc-web';


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

function rpcMethod (methodName, RequestClass) {
  return function (params) {
    const request = new RequestClass();
    Object.entries(params).forEach(([key, value]) => {
      const setMethodName = `set${key[0].toUpperCase()}${key.slice(1).toLowerCase()}`;
      request[setMethodName](value);
    });

    return _unWrap(this.client[methodName](request, this.metadata));
  };
}

function buildMethods (methods) {
  return Object.entries(methods).reduce((acc, [methodName, RequestClass]) => {
    acc[methodName] = rpcMethod(methodName, RequestClass);
    return acc;
  }, {});
}


class API {
  constructor (host, metadata = {}) {
    this.client = new APIPromiseClient(host);
    this.metadata = metadata;

    Object.assign(this, buildMethods({
      registerPlayer: RegisterPlayerRequest,
      requestPlayerLogin: RequestPlayerLoginRequest,
      playerLogin: PlayerLoginRequest,
      getSchedule: GetScheduleRequest,
      getScores: GetScoresRequest,
      createOrUpdateScore: CreateOrUpdateScoreRequest,
    }));
  }
}

const CACHE = new Map();

/**
 *
 * @param {Object} session
 *
 * @returns {API}
 */
export function getAPI (session) {
  const {
    config: { API_HOST_EXTERNAL: host },
    user: { authtoken },
  } = session;
  const cacheKey = `${host} ${authtoken}`;

  if (CACHE.has(cacheKey)) {
    return CACHE.get(cacheKey);
  }
  CACHE.set(cacheKey, new API(host, { authorization: authtoken }));

  return CACHE.get(cacheKey);
}
