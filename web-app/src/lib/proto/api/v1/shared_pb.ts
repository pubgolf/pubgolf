// Shared contains objects used across methods in multiple services.

// @generated by protoc-gen-es v1.2.0 with parameter "target=ts"
// @generated from file api/v1/shared.proto (package api.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from enum api.v1.ScoringCategory
 */
export enum ScoringCategory {
  /**
   * @generated from enum value: SCORING_CATEGORY_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: SCORING_CATEGORY_PUB_GOLF_NINE_HOLE = 1;
   */
  PUB_GOLF_NINE_HOLE = 1,

  /**
   * @generated from enum value: SCORING_CATEGORY_PUB_GOLF_FIVE_HOLE = 2;
   */
  PUB_GOLF_FIVE_HOLE = 2,

  /**
   * @generated from enum value: SCORING_CATEGORY_PUB_GOLF_CHALLENGES = 3;
   */
  PUB_GOLF_CHALLENGES = 3,
}
// Retrieve enum metadata with: proto3.getEnumType(ScoringCategory)
proto3.util.setEnumType(ScoringCategory, "api.v1.ScoringCategory", [
  { no: 0, name: "SCORING_CATEGORY_UNSPECIFIED" },
  { no: 1, name: "SCORING_CATEGORY_PUB_GOLF_NINE_HOLE" },
  { no: 2, name: "SCORING_CATEGORY_PUB_GOLF_FIVE_HOLE" },
  { no: 3, name: "SCORING_CATEGORY_PUB_GOLF_CHALLENGES" },
]);

/**
 * @generated from message api.v1.Color
 */
export class Color extends Message<Color> {
  /**
   * @generated from field: float r = 1;
   */
  r = 0;

  /**
   * @generated from field: float g = 2;
   */
  g = 0;

  /**
   * @generated from field: float b = 3;
   */
  b = 0;

  /**
   * @generated from field: float a = 4;
   */
  a = 0;

  constructor(data?: PartialMessage<Color>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v1.Color";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "r", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
    { no: 2, name: "g", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
    { no: 3, name: "b", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
    { no: 4, name: "a", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Color {
    return new Color().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Color {
    return new Color().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Color {
    return new Color().fromJsonString(jsonString, options);
  }

  static equals(a: Color | PlainMessage<Color> | undefined, b: Color | PlainMessage<Color> | undefined): boolean {
    return proto3.util.equals(Color, a, b);
  }
}

/**
 * @generated from message api.v1.Venue
 */
export class Venue extends Message<Venue> {
  /**
   * Global ID for the venue in ULID format (26 characters, base32), not to be confused with the venue key.
   *
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: string name = 2;
   */
  name = "";

  /**
   * Address string suitable for display or using for a mapping query.
   *
   * @generated from field: string address = 3;
   */
  address = "";

  /**
   * @generated from field: string image_url = 4;
   */
  imageUrl = "";

  constructor(data?: PartialMessage<Venue>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v1.Venue";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "address", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "image_url", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Venue {
    return new Venue().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Venue {
    return new Venue().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Venue {
    return new Venue().fromJsonString(jsonString, options);
  }

  static equals(a: Venue | PlainMessage<Venue> | undefined, b: Venue | PlainMessage<Venue> | undefined): boolean {
    return proto3.util.equals(Venue, a, b);
  }
}

/**
 * @generated from message api.v1.Player
 */
export class Player extends Message<Player> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: api.v1.PlayerData data = 2;
   */
  data?: PlayerData;

  constructor(data?: PartialMessage<Player>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v1.Player";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "data", kind: "message", T: PlayerData },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Player {
    return new Player().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Player {
    return new Player().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Player {
    return new Player().fromJsonString(jsonString, options);
  }

  static equals(a: Player | PlainMessage<Player> | undefined, b: Player | PlainMessage<Player> | undefined): boolean {
    return proto3.util.equals(Player, a, b);
  }
}

/**
 * PlayerData contains the user-editable fields for a player.
 *
 * @generated from message api.v1.PlayerData
 */
export class PlayerData extends Message<PlayerData> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: api.v1.ScoringCategory scoring_category = 2;
   */
  scoringCategory = ScoringCategory.UNSPECIFIED;

  constructor(data?: PartialMessage<PlayerData>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v1.PlayerData";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "scoring_category", kind: "enum", T: proto3.getEnumType(ScoringCategory) },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PlayerData {
    return new PlayerData().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PlayerData {
    return new PlayerData().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PlayerData {
    return new PlayerData().fromJsonString(jsonString, options);
  }

  static equals(a: PlayerData | PlainMessage<PlayerData> | undefined, b: PlayerData | PlainMessage<PlayerData> | undefined): boolean {
    return proto3.util.equals(PlayerData, a, b);
  }
}

/**
 * @generated from message api.v1.ScoreBoard
 */
export class ScoreBoard extends Message<ScoreBoard> {
  /**
   * @generated from field: repeated api.v1.ScoreBoard.ScoreBoardEntry scores = 1;
   */
  scores: ScoreBoard_ScoreBoardEntry[] = [];

  constructor(data?: PartialMessage<ScoreBoard>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v1.ScoreBoard";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "scores", kind: "message", T: ScoreBoard_ScoreBoardEntry, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ScoreBoard {
    return new ScoreBoard().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ScoreBoard {
    return new ScoreBoard().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ScoreBoard {
    return new ScoreBoard().fromJsonString(jsonString, options);
  }

  static equals(a: ScoreBoard | PlainMessage<ScoreBoard> | undefined, b: ScoreBoard | PlainMessage<ScoreBoard> | undefined): boolean {
    return proto3.util.equals(ScoreBoard, a, b);
  }
}

/**
 * @generated from enum api.v1.ScoreBoard.ScoreStatus
 */
export enum ScoreBoard_ScoreStatus {
  /**
   * @generated from enum value: SCORE_STATUS_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * SCORE_STATUS_PENDING indicates the player has not yet submitted/finalized their score for this round.
   *
   * @generated from enum value: SCORE_STATUS_PENDING = 1;
   */
  PENDING = 1,

  /**
   * SCORE_STATUS_FINALIZED indicates that a player's score is "locked in" as of a given milestone.
   *
   * @generated from enum value: SCORE_STATUS_FINALIZED = 2;
   */
  FINALIZED = 2,

  /**
   * SCORE_STATUS_INCOMPLETE indicates that a player's score is in an invalid or non-comparable state (e.g. they have dropped out of the event).
   *
   * @generated from enum value: SCORE_STATUS_INCOMPLETE = 3;
   */
  INCOMPLETE = 3,

  /**
   * SCORE_STATUS_NON_SCORING indicates that a player's score will not be counted towards the overall leaderboard.
   *
   * @generated from enum value: SCORE_STATUS_NON_SCORING = 4;
   */
  NON_SCORING = 4,
}
// Retrieve enum metadata with: proto3.getEnumType(ScoreBoard_ScoreStatus)
proto3.util.setEnumType(ScoreBoard_ScoreStatus, "api.v1.ScoreBoard.ScoreStatus", [
  { no: 0, name: "SCORE_STATUS_UNSPECIFIED" },
  { no: 1, name: "SCORE_STATUS_PENDING" },
  { no: 2, name: "SCORE_STATUS_FINALIZED" },
  { no: 3, name: "SCORE_STATUS_INCOMPLETE" },
  { no: 4, name: "SCORE_STATUS_NON_SCORING" },
]);

/**
 * @generated from message api.v1.ScoreBoard.ScoreBoardEntry
 */
export class ScoreBoard_ScoreBoardEntry extends Message<ScoreBoard_ScoreBoardEntry> {
  /**
   * @generated from field: optional string entity_id = 1;
   */
  entityId?: string;

  /**
   * @generated from field: string label = 2;
   */
  label = "";

  /**
   * @generated from field: int32 score = 3;
   */
  score = 0;

  /**
   * display_score_signed indicates that non-zero scores should be displayed with an explicit +/-.
   *
   * @generated from field: bool display_score_signed = 4;
   */
  displayScoreSigned = false;

  /**
   * rank is a display value indicating the ranking of the score. May be omitted in the case of ties, so ordering should be done based on the index of the `ScoreBoardEntry` in the repated field `Scoreboard.scores`.
   *
   * @generated from field: optional uint32 rank = 5;
   */
  rank?: number;

  /**
   * icon_key is an SF-Symbol name (e.g. "heart.fill").
   *
   * @generated from field: optional string icon_key = 6;
   */
  iconKey?: string;

  /**
   * @generated from field: optional api.v1.Color icon_color = 7;
   */
  iconColor?: Color;

  /**
   * @generated from field: api.v1.ScoreBoard.ScoreStatus status = 8;
   */
  status = ScoreBoard_ScoreStatus.UNSPECIFIED;

  /**
   * @generated from field: optional string status_details = 9;
   */
  statusDetails?: string;

  constructor(data?: PartialMessage<ScoreBoard_ScoreBoardEntry>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v1.ScoreBoard.ScoreBoardEntry";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "entity_id", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
    { no: 2, name: "label", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "score", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
    { no: 4, name: "display_score_signed", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 5, name: "rank", kind: "scalar", T: 13 /* ScalarType.UINT32 */, opt: true },
    { no: 6, name: "icon_key", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
    { no: 7, name: "icon_color", kind: "message", T: Color, opt: true },
    { no: 8, name: "status", kind: "enum", T: proto3.getEnumType(ScoreBoard_ScoreStatus) },
    { no: 9, name: "status_details", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ScoreBoard_ScoreBoardEntry {
    return new ScoreBoard_ScoreBoardEntry().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ScoreBoard_ScoreBoardEntry {
    return new ScoreBoard_ScoreBoardEntry().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ScoreBoard_ScoreBoardEntry {
    return new ScoreBoard_ScoreBoardEntry().fromJsonString(jsonString, options);
  }

  static equals(a: ScoreBoard_ScoreBoardEntry | PlainMessage<ScoreBoard_ScoreBoardEntry> | undefined, b: ScoreBoard_ScoreBoardEntry | PlainMessage<ScoreBoard_ScoreBoardEntry> | undefined): boolean {
    return proto3.util.equals(ScoreBoard_ScoreBoardEntry, a, b);
  }
}

