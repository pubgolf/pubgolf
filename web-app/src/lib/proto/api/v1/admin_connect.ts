// Admin defines the admin API service for the game management UI.

// @generated by protoc-gen-connect-es v0.8.6 with parameter "target=ts"
// @generated from file api/v1/admin.proto (package api.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { AdminServiceCreatePlayerRequest, AdminServiceCreatePlayerResponse, CreateStageScoreRequest, CreateStageScoreResponse, DeleteStageScoreRequest, DeleteStageScoreResponse, ListEventStagesRequest, ListEventStagesResponse, ListPlayersRequest, ListPlayersResponse, ListStageScoresRequest, ListStageScoresResponse, UpdatePlayerRequest, UpdatePlayerResponse, UpdateStageScoreRequest, UpdateStageScoreResponse } from "./admin_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * AdminService allows administering events with enhanced permissions.
 *
 * Players
 *
 * @generated from service api.v1.AdminService
 */
export const AdminService = {
  typeName: "api.v1.AdminService",
  methods: {
    /**
     * CreatePlayer creates a new player profile for a given event.
     *
     * @generated from rpc api.v1.AdminService.CreatePlayer
     */
    createPlayer: {
      name: "CreatePlayer",
      I: AdminServiceCreatePlayerRequest,
      O: AdminServiceCreatePlayerResponse,
      kind: MethodKind.Unary,
    },
    /**
     * UpdatePlayer modifies the player's profile and settings for a given event.
     *
     * @generated from rpc api.v1.AdminService.UpdatePlayer
     */
    updatePlayer: {
      name: "UpdatePlayer",
      I: UpdatePlayerRequest,
      O: UpdatePlayerResponse,
      kind: MethodKind.Unary,
    },
    /**
     * ListPlayers returns all players for a given event.
     *
     * @generated from rpc api.v1.AdminService.ListPlayers
     */
    listPlayers: {
      name: "ListPlayers",
      I: ListPlayersRequest,
      O: ListPlayersResponse,
      kind: MethodKind.Unary,
    },
    /**
     * ListEventStages returns a full schedule for an event.
     *
     * @generated from rpc api.v1.AdminService.ListEventStages
     */
    listEventStages: {
      name: "ListEventStages",
      I: ListEventStagesRequest,
      O: ListEventStagesResponse,
      kind: MethodKind.Unary,
    },
    /**
     * CreateStageScore sets the score and adjustments for a given pair of player and stage IDs.
     *
     * @generated from rpc api.v1.AdminService.CreateStageScore
     */
    createStageScore: {
      name: "CreateStageScore",
      I: CreateStageScoreRequest,
      O: CreateStageScoreResponse,
      kind: MethodKind.Unary,
    },
    /**
     * CreateStageScore updates the score and adjustments for a player/stage pair, based on their IDs.
     *
     * @generated from rpc api.v1.AdminService.UpdateStageScore
     */
    updateStageScore: {
      name: "UpdateStageScore",
      I: UpdateStageScoreRequest,
      O: UpdateStageScoreResponse,
      kind: MethodKind.Unary,
    },
    /**
     * ListStageScores returns all sets of (scores, adjustments[]) for an event, ordered chronologically by event stage, then chronologically by score creation time.
     *
     * @generated from rpc api.v1.AdminService.ListStageScores
     */
    listStageScores: {
      name: "ListStageScores",
      I: ListStageScoresRequest,
      O: ListStageScoresResponse,
      kind: MethodKind.Unary,
    },
    /**
     * DeleteStageScore removes all scoring data for a player/stage pair.
     *
     * @generated from rpc api.v1.AdminService.DeleteStageScore
     */
    deleteStageScore: {
      name: "DeleteStageScore",
      I: DeleteStageScoreRequest,
      O: DeleteStageScoreResponse,
      kind: MethodKind.Unary,
    },
  }
} as const;

