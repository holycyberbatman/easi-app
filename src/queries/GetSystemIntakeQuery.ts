import { gql } from '@apollo/client';

export default gql`
  query GetSystemIntake($id: UUID!) {
    systemIntake(id: $id) {
      id
      adminLead
      businessNeed
      businessSolution
      businessOwner {
        component
        name
      }
      contract {
        contractor
        endDate {
          day
          month
          year
        }
        hasContract
        startDate {
          day
          month
          year
        }
        vehicle
      }
      costs {
        isExpectingIncrease
        expectedIncreaseAmount
      }
      currentStage
      grbDate
      grtDate
      governanceTeams {
        isPresent
        teams {
          acronym
          collaborator
          key
          label
          name
        }
      }
      isso {
        isPresent
        name
      }
      fundingSource {
        fundingNumber
        isFunded
        source
      }
      lcid
      needsEaSupport
      productManager {
        component
        name
      }
      requester {
        component
        email
        name
      }
      requestName
      requestType
      status
      submittedAt
      businessCase {
        id
        alternativeASolution {
          acquisitionApproach
          cons
          costSavings
          hasUi
          hosting {
            cloudServiceType
            location
            type
          }
          pros
          security {
            isApproved
            isBeingReviewed
          }
          summary
          title
        }
        alternativeBSolution {
          acquisitionApproach
          cons
          costSavings
          hasUi
          hosting {
            cloudServiceType
            location
            type
          }
          pros
          security {
            isApproved
            isBeingReviewed
          }
          summary
          title
        }
        asIsSolution {
          cons
          costSavings
          pros
          summary
          title
        }
        businessNeed
        businessOwner {
          name
        }
        cmsBenefit
        euaUserId
        lifecycleCostLines {
          businessCaseId
          cost
          id
          phase
          solution
          year
        }
        preferredSolution {
          acquisitionApproach
          cons
          costSavings
          hasUi
          hosting {
            cloudServiceType
            location
            type
          }
          pros
          security {
            isApproved
            isBeingReviewed
          }
          summary
          title
        }
        priorityAlignment
        requestName
        requester {
          name
          phoneNumber
        }
        successIndicators
        status
      }
    }
  }
`;
