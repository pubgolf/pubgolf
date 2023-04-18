// Admin defines the admin API service for the game management UI.

// @generated by protoc-gen-es v1.2.0 with parameter "target=ts"
// @generated from file api/v1/admin.proto (package api.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import { Player, PlayerData } from "./shared_pb.js";

/**
 * @generated from message api.v1.AdminServiceCreatePlayerRequest
 */
export class AdminServiceCreatePlayerRequest extends Message<AdminServiceCreatePlayerRequest> {
  /**
   * @generated from field: string event_key = 1;
   */
  eventKey = "";

  /**
   * @generated from field: api.v1.PlayerData player_data = 2;
   */
  playerData?: PlayerData;

  constructor(data?: PartialMessage<AdminServiceCreatePlayerRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v1.AdminServiceCreatePlayerRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "event_key", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "player_data", kind: "message", T: PlayerData },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): AdminServiceCreatePlayerRequest {
    return new AdminServiceCreatePlayerRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): AdminServiceCreatePlayerRequest {
    return new AdminServiceCreatePlayerRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): AdminServiceCreatePlayerRequest {
    return new AdminServiceCreatePlayerRequest().fromJsonString(jsonString, options);
  }

  static equals(a: AdminServiceCreatePlayerRequest | PlainMessage<AdminServiceCreatePlayerRequest> | undefined, b: AdminServiceCreatePlayerRequest | PlainMessage<AdminServiceCreatePlayerRequest> | undefined): boolean {
    return proto3.util.equals(AdminServiceCreatePlayerRequest, a, b);
  }
}

/**
 * @generated from message api.v1.AdminServiceCreatePlayerResponse
 */
export class AdminServiceCreatePlayerResponse extends Message<AdminServiceCreatePlayerResponse> {
  /**
   * @generated from field: api.v1.Player player = 1;
   */
  player?: Player;

  constructor(data?: PartialMessage<AdminServiceCreatePlayerResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v1.AdminServiceCreatePlayerResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "player", kind: "message", T: Player },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): AdminServiceCreatePlayerResponse {
    return new AdminServiceCreatePlayerResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): AdminServiceCreatePlayerResponse {
    return new AdminServiceCreatePlayerResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): AdminServiceCreatePlayerResponse {
    return new AdminServiceCreatePlayerResponse().fromJsonString(jsonString, options);
  }

  static equals(a: AdminServiceCreatePlayerResponse | PlainMessage<AdminServiceCreatePlayerResponse> | undefined, b: AdminServiceCreatePlayerResponse | PlainMessage<AdminServiceCreatePlayerResponse> | undefined): boolean {
    return proto3.util.equals(AdminServiceCreatePlayerResponse, a, b);
  }
}

/**
 * @generated from message api.v1.UpdatePlayerRequest
 */
export class UpdatePlayerRequest extends Message<UpdatePlayerRequest> {
  /**
   * @generated from field: string player_id = 1;
   */
  playerId = "";

  /**
   * @generated from field: api.v1.PlayerData player_data = 2;
   */
  playerData?: PlayerData;

  constructor(data?: PartialMessage<UpdatePlayerRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v1.UpdatePlayerRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "player_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "player_data", kind: "message", T: PlayerData },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): UpdatePlayerRequest {
    return new UpdatePlayerRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): UpdatePlayerRequest {
    return new UpdatePlayerRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): UpdatePlayerRequest {
    return new UpdatePlayerRequest().fromJsonString(jsonString, options);
  }

  static equals(a: UpdatePlayerRequest | PlainMessage<UpdatePlayerRequest> | undefined, b: UpdatePlayerRequest | PlainMessage<UpdatePlayerRequest> | undefined): boolean {
    return proto3.util.equals(UpdatePlayerRequest, a, b);
  }
}

/**
 * @generated from message api.v1.UpdatePlayerResponse
 */
export class UpdatePlayerResponse extends Message<UpdatePlayerResponse> {
  /**
   * @generated from field: api.v1.Player player = 1;
   */
  player?: Player;

  constructor(data?: PartialMessage<UpdatePlayerResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v1.UpdatePlayerResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "player", kind: "message", T: Player },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): UpdatePlayerResponse {
    return new UpdatePlayerResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): UpdatePlayerResponse {
    return new UpdatePlayerResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): UpdatePlayerResponse {
    return new UpdatePlayerResponse().fromJsonString(jsonString, options);
  }

  static equals(a: UpdatePlayerResponse | PlainMessage<UpdatePlayerResponse> | undefined, b: UpdatePlayerResponse | PlainMessage<UpdatePlayerResponse> | undefined): boolean {
    return proto3.util.equals(UpdatePlayerResponse, a, b);
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
   * @generated from field: repeated api.v1.Player players = 1;
   */
  players: Player[] = [];

  constructor(data?: PartialMessage<ListPlayersResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v1.ListPlayersResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "players", kind: "message", T: Player, repeated: true },
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

