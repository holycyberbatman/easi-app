/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

import { SystemIntakeRequestType, SystemIntakeStatus, LifecycleCostPhase, LifecycleCostSolution, LifecycleCostYear, BusinessCaseStatus } from "./../../types/graphql-global-types";

// ====================================================
// GraphQL query operation: GetSystemIntake
// ====================================================

export interface GetSystemIntake_systemIntake_businessOwner {
  __typename: "SystemIntakeBusinessOwner";
  component: string | null;
  name: string | null;
}

export interface GetSystemIntake_systemIntake_contract_endDate {
  __typename: "ContractDate";
  day: string | null;
  month: string | null;
  year: string | null;
}

export interface GetSystemIntake_systemIntake_contract_startDate {
  __typename: "ContractDate";
  day: string | null;
  month: string | null;
  year: string | null;
}

export interface GetSystemIntake_systemIntake_contract {
  __typename: "SystemIntakeContract";
  contractor: string | null;
  endDate: GetSystemIntake_systemIntake_contract_endDate | null;
  hasContract: string | null;
  startDate: GetSystemIntake_systemIntake_contract_startDate | null;
  vehicle: string | null;
}

export interface GetSystemIntake_systemIntake_costs {
  __typename: "SystemIntakeCosts";
  isExpectingIncrease: string | null;
  expectedIncreaseAmount: string | null;
}

export interface GetSystemIntake_systemIntake_governanceTeams_teams {
  __typename: "SystemIntakeCollaborator";
  acronym: string | null;
  collaborator: string | null;
  key: string | null;
  label: string | null;
  name: string | null;
}

export interface GetSystemIntake_systemIntake_governanceTeams {
  __typename: "SystemIntakeGovernanceTeam";
  isPresent: boolean | null;
  teams: GetSystemIntake_systemIntake_governanceTeams_teams[] | null;
}

export interface GetSystemIntake_systemIntake_isso {
  __typename: "SystemIntakeISSO";
  isPresent: boolean | null;
  name: string | null;
}

export interface GetSystemIntake_systemIntake_fundingSource {
  __typename: "SystemIntakeFundingSource";
  fundingNumber: string | null;
  isFunded: boolean | null;
  source: string | null;
}

export interface GetSystemIntake_systemIntake_productManager {
  __typename: "SystemIntakeProductManager";
  component: string | null;
  name: string | null;
}

export interface GetSystemIntake_systemIntake_requester {
  __typename: "SystemIntakeRequester";
  component: string | null;
  email: string | null;
  name: string;
}

export interface GetSystemIntake_systemIntake_businessCase_alternativeASolution_hosting {
  __typename: "BusinessCaseSolutionHosting";
  cloudServiceType: string | null;
  location: string | null;
  type: string | null;
}

export interface GetSystemIntake_systemIntake_businessCase_alternativeASolution_security {
  __typename: "BusinessCaseSolutionSecurity";
  isApproved: boolean | null;
  isBeingReviewed: string | null;
}

export interface GetSystemIntake_systemIntake_businessCase_alternativeASolution {
  __typename: "BusinessCaseSolution";
  acquisitionApproach: string | null;
  cons: string | null;
  costSavings: string | null;
  hasUi: string | null;
  hosting: GetSystemIntake_systemIntake_businessCase_alternativeASolution_hosting | null;
  pros: string | null;
  security: GetSystemIntake_systemIntake_businessCase_alternativeASolution_security | null;
  summary: string | null;
  title: string | null;
}

export interface GetSystemIntake_systemIntake_businessCase_alternativeBSolution_hosting {
  __typename: "BusinessCaseSolutionHosting";
  cloudServiceType: string | null;
  location: string | null;
  type: string | null;
}

export interface GetSystemIntake_systemIntake_businessCase_alternativeBSolution_security {
  __typename: "BusinessCaseSolutionSecurity";
  isApproved: boolean | null;
  isBeingReviewed: string | null;
}

export interface GetSystemIntake_systemIntake_businessCase_alternativeBSolution {
  __typename: "BusinessCaseSolution";
  acquisitionApproach: string | null;
  cons: string | null;
  costSavings: string | null;
  hasUi: string | null;
  hosting: GetSystemIntake_systemIntake_businessCase_alternativeBSolution_hosting | null;
  pros: string | null;
  security: GetSystemIntake_systemIntake_businessCase_alternativeBSolution_security | null;
  summary: string | null;
  title: string | null;
}

export interface GetSystemIntake_systemIntake_businessCase_asIsSolution {
  __typename: "BusinessCaseAsIsSolution";
  cons: string | null;
  costSavings: string | null;
  pros: string | null;
  summary: string | null;
  title: string | null;
}

export interface GetSystemIntake_systemIntake_businessCase_businessOwner {
  __typename: "BusinessCaseBusinessOwner";
  name: string | null;
}

export interface GetSystemIntake_systemIntake_businessCase_lifecycleCostLines {
  __typename: "EstimatedLifecycleCost";
  businessCaseId: UUID;
  cost: number | null;
  id: UUID;
  phase: LifecycleCostPhase | null;
  solution: LifecycleCostSolution | null;
  year: LifecycleCostYear | null;
}

export interface GetSystemIntake_systemIntake_businessCase_preferredSolution_hosting {
  __typename: "BusinessCaseSolutionHosting";
  cloudServiceType: string | null;
  location: string | null;
  type: string | null;
}

export interface GetSystemIntake_systemIntake_businessCase_preferredSolution_security {
  __typename: "BusinessCaseSolutionSecurity";
  isApproved: boolean | null;
  isBeingReviewed: string | null;
}

export interface GetSystemIntake_systemIntake_businessCase_preferredSolution {
  __typename: "BusinessCaseSolution";
  acquisitionApproach: string | null;
  cons: string | null;
  costSavings: string | null;
  hasUi: string | null;
  hosting: GetSystemIntake_systemIntake_businessCase_preferredSolution_hosting | null;
  pros: string | null;
  security: GetSystemIntake_systemIntake_businessCase_preferredSolution_security | null;
  summary: string | null;
  title: string | null;
}

export interface GetSystemIntake_systemIntake_businessCase_requester {
  __typename: "BusinessCaseRequester";
  name: string | null;
  phoneNumber: string | null;
}

export interface GetSystemIntake_systemIntake_businessCase {
  __typename: "BusinessCase";
  id: UUID;
  alternativeASolution: GetSystemIntake_systemIntake_businessCase_alternativeASolution | null;
  alternativeBSolution: GetSystemIntake_systemIntake_businessCase_alternativeBSolution | null;
  asIsSolution: GetSystemIntake_systemIntake_businessCase_asIsSolution | null;
  businessNeed: string | null;
  businessOwner: GetSystemIntake_systemIntake_businessCase_businessOwner | null;
  cmsBenefit: string | null;
  euaUserId: string;
  lifecycleCostLines: GetSystemIntake_systemIntake_businessCase_lifecycleCostLines[] | null;
  preferredSolution: GetSystemIntake_systemIntake_businessCase_preferredSolution | null;
  priorityAlignment: string | null;
  requestName: string | null;
  requester: GetSystemIntake_systemIntake_businessCase_requester | null;
  successIndicators: string | null;
  status: BusinessCaseStatus;
}

export interface GetSystemIntake_systemIntake {
  __typename: "SystemIntake";
  id: UUID;
  adminLead: string | null;
  businessNeed: string | null;
  businessSolution: string | null;
  businessOwner: GetSystemIntake_systemIntake_businessOwner | null;
  contract: GetSystemIntake_systemIntake_contract | null;
  costs: GetSystemIntake_systemIntake_costs | null;
  currentStage: string | null;
  grbDate: Time | null;
  grtDate: Time | null;
  governanceTeams: GetSystemIntake_systemIntake_governanceTeams | null;
  isso: GetSystemIntake_systemIntake_isso | null;
  fundingSource: GetSystemIntake_systemIntake_fundingSource | null;
  lcid: string | null;
  needsEaSupport: boolean | null;
  productManager: GetSystemIntake_systemIntake_productManager | null;
  requester: GetSystemIntake_systemIntake_requester;
  requestName: string | null;
  requestType: SystemIntakeRequestType;
  status: SystemIntakeStatus;
  submittedAt: Time | null;
  businessCase: GetSystemIntake_systemIntake_businessCase | null;
}

export interface GetSystemIntake {
  systemIntake: GetSystemIntake_systemIntake | null;
}

export interface GetSystemIntakeVariables {
  id: UUID;
}
