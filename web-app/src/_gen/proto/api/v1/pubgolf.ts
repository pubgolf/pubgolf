/* eslint-disable */
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import * as Long from "long";
import { grpc } from "@improbable-eng/grpc-web";
import { BrowserHeaders } from "browser-headers";

export const protobufPackage = "api.v1";

export interface ClientVersionRequest {
  clientVersion: number;
}

export interface ClientVersionResponse {
  versionStatus: ClientVersionResponse_VersionStatus;
}

export enum ClientVersionResponse_VersionStatus {
  VERSION_STATUS_UNSPECIFIED = 0,
  VERSION_STATUS_OK = 1,
  VERSION_STATUS_OUTDATED = 2,
  VERSION_STATUS_INCOMPATIBLE = 3,
  UNRECOGNIZED = -1,
}

export function clientVersionResponse_VersionStatusFromJSON(
  object: any
): ClientVersionResponse_VersionStatus {
  switch (object) {
    case 0:
    case "VERSION_STATUS_UNSPECIFIED":
      return ClientVersionResponse_VersionStatus.VERSION_STATUS_UNSPECIFIED;
    case 1:
    case "VERSION_STATUS_OK":
      return ClientVersionResponse_VersionStatus.VERSION_STATUS_OK;
    case 2:
    case "VERSION_STATUS_OUTDATED":
      return ClientVersionResponse_VersionStatus.VERSION_STATUS_OUTDATED;
    case 3:
    case "VERSION_STATUS_INCOMPATIBLE":
      return ClientVersionResponse_VersionStatus.VERSION_STATUS_INCOMPATIBLE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ClientVersionResponse_VersionStatus.UNRECOGNIZED;
  }
}

export function clientVersionResponse_VersionStatusToJSON(
  object: ClientVersionResponse_VersionStatus
): string {
  switch (object) {
    case ClientVersionResponse_VersionStatus.VERSION_STATUS_UNSPECIFIED:
      return "VERSION_STATUS_UNSPECIFIED";
    case ClientVersionResponse_VersionStatus.VERSION_STATUS_OK:
      return "VERSION_STATUS_OK";
    case ClientVersionResponse_VersionStatus.VERSION_STATUS_OUTDATED:
      return "VERSION_STATUS_OUTDATED";
    case ClientVersionResponse_VersionStatus.VERSION_STATUS_INCOMPATIBLE:
      return "VERSION_STATUS_INCOMPATIBLE";
    default:
      return "UNKNOWN";
  }
}

const baseClientVersionRequest: object = { clientVersion: 0 };

export const ClientVersionRequest = {
  encode(
    message: ClientVersionRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.clientVersion !== 0) {
      writer.uint32(8).uint32(message.clientVersion);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ClientVersionRequest {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseClientVersionRequest } as ClientVersionRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.clientVersion = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ClientVersionRequest {
    const message = { ...baseClientVersionRequest } as ClientVersionRequest;
    message.clientVersion =
      object.clientVersion !== undefined && object.clientVersion !== null
        ? Number(object.clientVersion)
        : 0;
    return message;
  },

  toJSON(message: ClientVersionRequest): unknown {
    const obj: any = {};
    message.clientVersion !== undefined &&
      (obj.clientVersion = Math.round(message.clientVersion));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ClientVersionRequest>, I>>(
    object: I
  ): ClientVersionRequest {
    const message = { ...baseClientVersionRequest } as ClientVersionRequest;
    message.clientVersion = object.clientVersion ?? 0;
    return message;
  },
};

const baseClientVersionResponse: object = { versionStatus: 0 };

export const ClientVersionResponse = {
  encode(
    message: ClientVersionResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.versionStatus !== 0) {
      writer.uint32(8).int32(message.versionStatus);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ClientVersionResponse {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseClientVersionResponse } as ClientVersionResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.versionStatus = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ClientVersionResponse {
    const message = { ...baseClientVersionResponse } as ClientVersionResponse;
    message.versionStatus =
      object.versionStatus !== undefined && object.versionStatus !== null
        ? clientVersionResponse_VersionStatusFromJSON(object.versionStatus)
        : 0;
    return message;
  },

  toJSON(message: ClientVersionResponse): unknown {
    const obj: any = {};
    message.versionStatus !== undefined &&
      (obj.versionStatus = clientVersionResponse_VersionStatusToJSON(
        message.versionStatus
      ));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ClientVersionResponse>, I>>(
    object: I
  ): ClientVersionResponse {
    const message = { ...baseClientVersionResponse } as ClientVersionResponse;
    message.versionStatus = object.versionStatus ?? 0;
    return message;
  },
};

/** PubGolfService is the API server which handles all scorekeeping, scheduling and account management for pub golf. */
export interface PubGolfService {
  /** ClientVersion indicates to the server that a client of a given version is attempting to connect, and allows the server to respond with a "soft" or "hard" upgrade notification. */
  ClientVersion(
    request: DeepPartial<ClientVersionRequest>,
    metadata?: grpc.Metadata
  ): Promise<ClientVersionResponse>;
}

export class PubGolfServiceClientImpl implements PubGolfService {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.ClientVersion = this.ClientVersion.bind(this);
  }

  ClientVersion(
    request: DeepPartial<ClientVersionRequest>,
    metadata?: grpc.Metadata
  ): Promise<ClientVersionResponse> {
    return this.rpc.unary(
      PubGolfServiceClientVersionDesc,
      ClientVersionRequest.fromPartial(request),
      metadata
    );
  }
}

export const PubGolfServiceDesc = {
  serviceName: "api.v1.PubGolfService",
};

export const PubGolfServiceClientVersionDesc: UnaryMethodDefinitionish = {
  methodName: "ClientVersion",
  service: PubGolfServiceDesc,
  requestStream: false,
  responseStream: false,
  requestType: {
    serializeBinary() {
      return ClientVersionRequest.encode(this).finish();
    },
  } as any,
  responseType: {
    deserializeBinary(data: Uint8Array) {
      return {
        ...ClientVersionResponse.decode(data),
        toObject() {
          return this;
        },
      };
    },
  } as any,
};

interface UnaryMethodDefinitionishR
  extends grpc.UnaryMethodDefinition<any, any> {
  requestStream: any;
  responseStream: any;
}

type UnaryMethodDefinitionish = UnaryMethodDefinitionishR;

interface Rpc {
  unary<T extends UnaryMethodDefinitionish>(
    methodDesc: T,
    request: any,
    metadata: grpc.Metadata | undefined
  ): Promise<any>;
}

export class GrpcWebImpl {
  private host: string;
  private options: {
    transport?: grpc.TransportFactory;

    debug?: boolean;
    metadata?: grpc.Metadata;
  };

  constructor(
    host: string,
    options: {
      transport?: grpc.TransportFactory;

      debug?: boolean;
      metadata?: grpc.Metadata;
    }
  ) {
    this.host = host;
    this.options = options;
  }

  unary<T extends UnaryMethodDefinitionish>(
    methodDesc: T,
    _request: any,
    metadata: grpc.Metadata | undefined
  ): Promise<any> {
    const request = { ..._request, ...methodDesc.requestType };
    const maybeCombinedMetadata =
      metadata && this.options.metadata
        ? new BrowserHeaders({
            ...this.options?.metadata.headersMap,
            ...metadata?.headersMap,
          })
        : metadata || this.options.metadata;
    return new Promise((resolve, reject) => {
      grpc.unary(methodDesc, {
        request,
        host: this.host,
        metadata: maybeCombinedMetadata,
        transport: this.options.transport,
        debug: this.options.debug,
        onEnd: function (response) {
          if (response.status === grpc.Code.OK) {
            resolve(response.message);
          } else {
            const err = new Error(response.statusMessage) as any;
            err.code = response.status;
            err.metadata = response.trailers;
            reject(err);
          }
        },
      });
    });
  }
}

type Builtin =
  | Date
  | Function
  | Uint8Array
  | string
  | number
  | boolean
  | undefined;

export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin
  ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & Record<
        Exclude<keyof I, KeysOfUnion<P>>,
        never
      >;

// If you get a compile-error about 'Constructor<Long> and ... have no overlap',
// add '--ts_proto_opt=esModuleInterop=true' as a flag when calling 'protoc'.
if (util.Long !== Long) {
  util.Long = Long as any;
  configure();
}
