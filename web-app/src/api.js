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

export class API {
  constructor (eventKey, metadata = {}) {
    this.eventKey = eventKey;
    this.client = new APIPromiseClient('http://127.0.0.1:8080');
    this.metadata = metadata;
  }

  /**
   * @param promise
   *
   * @returns {Object}
   */
  _unWrap (promise) {
    return promise.then(
      instance => instance.toObject(),
      error => {
        if (process.env.NODE_ENV === 'development') {
          console.error(error);
        }

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
   * @param {Object} playerInfo
   *
   * @returns {Promise<RegisterPlayerReply>}
   */
  registerPlayer (playerInfo) {
    const request = new RegisterPlayerRequest();
    request.setEventkey(this.eventKey);
    request.setName(playerInfo.name);
    request.setPhonenumber(`+1${playerInfo.phone}`);
    request.setLeague(playerInfo.league);

    return this._unWrap(this.client.registerPlayer(request, this.metadata));
  }

  /**
   * @param {string} phone
   *
   * @returns {Promise<RequestPlayerLoginReply>}
   */
  requestPlayerLogin (phone) {
    const request = new RequestPlayerLoginRequest();
    request.setEventkey(this.eventKey);
    request.setPhonenumber(`+1${phone}`);

    return this._unWrap(this.client.requestPlayerLogin(
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
  playerLogin (phone, code) {
    const request = new PlayerLoginRequest();
    request.setEventkey(this.eventKey);
    request.setPhonenumber(`+1${phone}`);
    request.setAuthcode(code);

    return this._unWrap(this.client.playerLogin(
      request,
      this.metadata,
    )).then(({ authtoken }) => {
      this.metadata = {
        ...this.metadata,
        authorization: authtoken,
      };
    });
  }

  /**
   * @returns {Promise<GetScheduleReply>}
   */
  getSchedule () {
    const request = new GetScheduleRequest();
    request.setEventkey(this.eventKey);

    return this._unWrap(this.client.getSchedule(request, this.metadata));
  }

  /**
   * @returns {Promise<GetScoresReply>}
   */
  getScores () {
    const request = new GetScoresRequest();
    request.setEventkey(this.eventKey);

    return this._unWrap(this.client.getScores(request, this.metadata));
  }

  /**
   * @param {string} playerId
   * @param {number} venueId
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

    return this._unWrap(this.client.createOrUpdateScore(
      request,
      this.metadata,
    ));
  }
}

export const DEFAULT_CLIENT = new API('sf-2019');
