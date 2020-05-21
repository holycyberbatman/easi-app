import { LifecyclePhase } from 'types/estimatedLifecycle';

export type EstimatedLifecycleCostLines = {
  year1: LifecyclePhase[];
  year2: LifecyclePhase[];
  year3: LifecyclePhase[];
  year4: LifecyclePhase[];
  year5: LifecyclePhase[];
};

// Base Solution
export type BusinessCaseSolution = {
  title: string;
  summary: string;
  pros: string;
  cons: string;
  estimatedLifecycleCost: EstimatedLifecycleCostLines;
  costSavings: string;
};

// Proposed Solution
export type ProposedBusinessCaseSolution = BusinessCaseSolution & {
  acquisitionApproach: string;
};

// Business Case Form Model
export type BusinessCaseModel = {
  requestName: string;
  requester: {
    name: string;
    phoneNumber: string;
  };
  businessOwner: {
    name: string;
  };
  businessNeed: string;
  cmsBenefit: string;
  priorityAlignment: string;
  successIndicators: string;
  asIsSolution: BusinessCaseSolution;
  preferredSolution: ProposedBusinessCaseSolution;
  alternativeA: ProposedBusinessCaseSolution;
  alternativeB?: ProposedBusinessCaseSolution;
};

// Redux store type for business case
export type BusinessCaseState = {
  form: BusinessCaseModel;
  isLoading: boolean | null;
  error: any;
};
