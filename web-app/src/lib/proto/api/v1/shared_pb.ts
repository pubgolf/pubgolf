// Shared contains objects used across methods in multiple services.

// @generated by protoc-gen-es v1.2.0 with parameter "target=ts"
// @generated from file api/v1/shared.proto (package api.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { proto3 } from "@bufbuild/protobuf";

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
