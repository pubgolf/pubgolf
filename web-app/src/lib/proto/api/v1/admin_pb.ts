// Admin defines the admin API service for the game management UI.

// @generated by protoc-gen-es v1.2.0 with parameter "target=ts"
// @generated from file api/v1/admin.proto (package api.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from enum api.v1.ScoringCategory
 */
export enum ScoringCategory {
  /**
   * @generated from enum value: UNKNOWN = 0;
   */
  UNKNOWN = 0,

  /**
   * @generated from enum value: PUB_GOLF_NINE_HOLE = 1;
   */
  PUB_GOLF_NINE_HOLE = 1,

  /**
   * @generated from enum value: PUB_GOLF_FIVE_HOLE = 2;
   */
  PUB_GOLF_FIVE_HOLE = 2,

  /**
   * @generated from enum value: PUB_GOLF_CHALLENGES = 3;
   */
  PUB_GOLF_CHALLENGES = 3,
}
// Retrieve enum metadata with: proto3.getEnumType(ScoringCategory)
proto3.util.setEnumType(ScoringCategory, "api.v1.ScoringCategory", [
  { no: 0, name: "UNKNOWN" },
  { no: 1, name: "PUB_GOLF_NINE_HOLE" },
  { no: 2, name: "PUB_GOLF_FIVE_HOLE" },
  { no: 3, name: "PUB_GOLF_CHALLENGES" },
]);

/**
 * @generated from message api.v1.CreatePlayerRequest
 */
export class CreatePlayerRequest extends Message<CreatePlayerRequest> {
  /**
   * @generated from field: string event_key = 1;
   */
  eventKey = "";

  /**
   * @generated from field: api.v1.CreatePlayerRequest.PlayerInfo player = 2;
   */
  player?: CreatePlayerRequest_PlayerInfo;

  constructor(data?: PartialMessage<CreatePlayerRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v1.CreatePlayerRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "event_key", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "player", kind: "message", T: CreatePlayerRequest_PlayerInfo },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CreatePlayerRequest {
    return new CreatePlayerRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CreatePlayerRequest {
    return new CreatePlayerRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CreatePlayerRequest {
    return new CreatePlayerRequest().fromJsonString(jsonString, options);
  }

  static equals(a: CreatePlayerRequest | PlainMessage<CreatePlayerRequest> | undefined, b: CreatePlayerRequest | PlainMessage<CreatePlayerRequest> | undefined): boolean {
    return proto3.util.equals(CreatePlayerRequest, a, b);
  }
}

/**
 * @generated from message api.v1.CreatePlayerRequest.PlayerInfo
 */
export class CreatePlayerRequest_PlayerInfo extends Message<CreatePlayerRequest_PlayerInfo> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: optional api.v1.ScoringCategory scoring_category = 2;
   */
  scoringCategory?: ScoringCategory;

  constructor(data?: PartialMessage<CreatePlayerRequest_PlayerInfo>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v1.CreatePlayerRequest.PlayerInfo";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "scoring_category", kind: "enum", T: proto3.getEnumType(ScoringCategory), opt: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CreatePlayerRequest_PlayerInfo {
    return new CreatePlayerRequest_PlayerInfo().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CreatePlayerRequest_PlayerInfo {
    return new CreatePlayerRequest_PlayerInfo().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CreatePlayerRequest_PlayerInfo {
    return new CreatePlayerRequest_PlayerInfo().fromJsonString(jsonString, options);
  }

  static equals(a: CreatePlayerRequest_PlayerInfo | PlainMessage<CreatePlayerRequest_PlayerInfo> | undefined, b: CreatePlayerRequest_PlayerInfo | PlainMessage<CreatePlayerRequest_PlayerInfo> | undefined): boolean {
    return proto3.util.equals(CreatePlayerRequest_PlayerInfo, a, b);
  }
}

/**
 * @generated from message api.v1.CreatePlayerResponse
 */
export class CreatePlayerResponse extends Message<CreatePlayerResponse> {
  /**
   * @generated from field: string player_id = 1;
   */
  playerId = "";

  constructor(data?: PartialMessage<CreatePlayerResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v1.CreatePlayerResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "player_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CreatePlayerResponse {
    return new CreatePlayerResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CreatePlayerResponse {
    return new CreatePlayerResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CreatePlayerResponse {
    return new CreatePlayerResponse().fromJsonString(jsonString, options);
  }

  static equals(a: CreatePlayerResponse | PlainMessage<CreatePlayerResponse> | undefined, b: CreatePlayerResponse | PlainMessage<CreatePlayerResponse> | undefined): boolean {
    return proto3.util.equals(CreatePlayerResponse, a, b);
  }
}

/**
 * @generated from message api.v1.ListPlayersRequest
 */
export class ListPlayersRequest extends Message<ListPlayersRequest> {
  /**
   * @generated from field: string event_key = 1;
   */
  eventKey = "";

  constructor(data?: PartialMessage<ListPlayersRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v1.ListPlayersRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "event_key", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ListPlayersRequest {
    return new ListPlayersRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ListPlayersRequest {
    return new ListPlayersRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ListPlayersRequest {
    return new ListPlayersRequest().fromJsonString(jsonString, options);
  }

  static equals(a: ListPlayersRequest | PlainMessage<ListPlayersRequest> | undefined, b: ListPlayersRequest | PlainMessage<ListPlayersRequest> | undefined): boolean {
    return proto3.util.equals(ListPlayersRequest, a, b);
  }
}

/**
 * @generated from message api.v1.ListPlayersResponse
 */
export class ListPlayersResponse extends Message<ListPlayersResponse> {
  /**
   * @generated from field: repeated api.v1.ListPlayersResponse.PlayerInfo players = 1;
   */
  players: ListPlayersResponse_PlayerInfo[] = [];

  constructor(data?: PartialMessage<ListPlayersResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v1.ListPlayersResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "players", kind: "message", T: ListPlayersResponse_PlayerInfo, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ListPlayersResponse {
    return new ListPlayersResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ListPlayersResponse {
    return new ListPlayersResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ListPlayersResponse {
    return new ListPlayersResponse().fromJsonString(jsonString, options);
  }

  static equals(a: ListPlayersResponse | PlainMessage<ListPlayersResponse> | undefined, b: ListPlayersResponse | PlainMessage<ListPlayersResponse> | undefined): boolean {
    return proto3.util.equals(ListPlayersResponse, a, b);
  }
}

/**
 * @generated from message api.v1.ListPlayersResponse.PlayerInfo
 */
export class ListPlayersResponse_PlayerInfo extends Message<ListPlayersResponse_PlayerInfo> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: string name = 2;
   */
  name = "";

  /**
   * @generated from field: optional api.v1.ScoringCategory scoring_category = 3;
   */
  scoringCategory?: ScoringCategory;

  constructor(data?: PartialMessage<ListPlayersResponse_PlayerInfo>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v1.ListPlayersResponse.PlayerInfo";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "scoring_category", kind: "enum", T: proto3.getEnumType(ScoringCategory), opt: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ListPlayersResponse_PlayerInfo {
    return new ListPlayersResponse_PlayerInfo().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ListPlayersResponse_PlayerInfo {
    return new ListPlayersResponse_PlayerInfo().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ListPlayersResponse_PlayerInfo {
    return new ListPlayersResponse_PlayerInfo().fromJsonString(jsonString, options);
  }

  static equals(a: ListPlayersResponse_PlayerInfo | PlainMessage<ListPlayersResponse_PlayerInfo> | undefined, b: ListPlayersResponse_PlayerInfo | PlainMessage<ListPlayersResponse_PlayerInfo> | undefined): boolean {
    return proto3.util.equals(ListPlayersResponse_PlayerInfo, a, b);
  }
}

