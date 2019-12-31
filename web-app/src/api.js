import { StatusCode } from 'grpc-web';


const Cookies = require('js-cookie');

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

// Wrapper for handling cookies so it's abstracted from the rest of the application
export const getCookieJar = () => { // eslint-disable-line arrow-body-style
  return {
    get (name) {
      return Cookies.get(name);
    },
    set (name, value, attributes) {
      Cookies.set(name, value, attributes);
    },
    remove (name, attributes) {
      Cookies.remove(name, attributes);
    },
  };
};

const TOKEN_COOKIE = 'pg-token';


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


class API {
  constructor (host, metadata = {}) {
    this.client = new APIPromiseClient(host);
    this.metadata = metadata;
  }

  isLoggedIn () {
    return Boolean(this.metadata && this.metadata.authorization);
  }

  _logIn (token) {
    if (!token) return;

    this.metadata = {
      ...this.metadata,
      authorization: token,
    };
    this._cookieJar.set(TOKEN_COOKIE, token);
  }

  _logOut () {
    delete this.metadata.authorization;
    this._cookieJar.remove(TOKEN_COOKIE);
  }

  /**
   * @param {Object} playerInfo
   *
   * @returns {Promise<RegisterPlayerReply>}
   */
  registerPlayer ({
    eventKey,
    name,
    phoneNumber,
    league,
  }) {
    const request = new RegisterPlayerRequest();
    request.setEventkey(eventKey);
    request.setName(name);
    request.setPhonenumber(`+1${phoneNumber}`);
    request.setLeague(league);

    return _unWrap(this.client.registerPlayer(request, this.metadata));
  }

  /**
   * @param {string} phone
   *
   * @returns {Promise<RequestPlayerLoginReply>}
   */
  requestPlayerLogin ({ eventKey, phoneNumber }) {
    const request = new RequestPlayerLoginRequest();
    request.setEventkey(eventKey);
    request.setPhonenumber(`+1${phoneNumber}`);

    return _unWrap(this.client.requestPlayerLogin(
      request,
      this.metadata,
    ));
  }

  /**
   * @param {string} phone
   * @param {number} code
   *
   * @returns {Promise<PlayerLoginReply>}
   */
  playerLogin ({ eventKey, phoneNumber, authCode }) {
    const request = new PlayerLoginRequest();
    request.setEventkey(eventKey);
    request.setPhonenumber(`+1${phoneNumber}`);
    request.setAuthcode(authCode);

    return _unWrap(this.client.playerLogin(
      request,
      this.metadata,
    )).then(({ authtoken }) => ({
      token: authtoken,
    }));
  }

  /**
   * @returns {Promise<GetScheduleReply>}
   */
  getSchedule () {
    const request = new GetScheduleRequest();
    request.setEventkey(this.eventKey);

    return _unWrap(this.client.getSchedule(
      request,
      this.metadata,
    )).then(response => response, (error) => {
      if (error.code === StatusCode.PERMISSION_DENIED) {
        this._logOut();
      }
      throw error;
    });
  }

  /**
   * @returns {Promise<GetScoresReply>}
   */
  getScores () {
    const request = new GetScoresRequest();
    request.setEventkey(this.eventKey);

    return _unWrap(this.client.getScores(request, this.metadata));
  }

  /**
   * @param {string} playerId
   * @param {string} venueId
   * @param {number} strokes
   *
   * @returns {Promise<CreateOrUpdateScoreReply>}
   */
  createOrUpdateScore ({ playerId, venueId, strokes }) {
    const request = new CreateOrUpdateScoreRequest();
    request.setEventkey(this.eventKey);
    request.setPlayerid(playerId);
    request.setVenueid(venueId);
    request.setStrokes(strokes);

    return _unWrap(this.client.createOrUpdateScore(
      request,
      this.metadata,
    ));
  }
}

const CACHE = new Map();

export function getAPI (session) {
  const {
    config: { API_HOST_EXTERNAL: host },
    user: { token },
  } = session;
  const cacheKey = `${host} ${token}`;

  if (CACHE.has(cacheKey)) {
    return CACHE.get(cacheKey);
  }
  CACHE.set(cacheKey, new API(host, { authorization: token }));

  return CACHE.get(cacheKey);
}
