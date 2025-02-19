package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/guregu/null"
	_ "github.com/lib/pq" // required for postgres driver in sql

	"github.com/cmsgov/easi-app/pkg/models"
)

func (s GraphQLTestSuite) TestFetchSystemIntakeQuery() {
	ctx := context.Background()
	projectName := "Big Project"
	businessOwner := "Firstname Lastname"
	businessOwnerComponent := "OIT"

	intake, intakeErr := s.store.CreateSystemIntake(ctx, &models.SystemIntake{
		ProjectName:            null.StringFrom(projectName),
		Status:                 models.SystemIntakeStatusINTAKESUBMITTED,
		RequestType:            models.SystemIntakeRequestTypeNEW,
		BusinessOwner:          null.StringFrom(businessOwner),
		BusinessOwnerComponent: null.StringFrom(businessOwnerComponent),
	})
	s.NoError(intakeErr)

	var resp struct {
		SystemIntake struct {
			ID            string
			RequestName   string
			Status        string
			RequestType   string
			BusinessOwner struct {
				Name      string
				Component string
			}
			BusinessOwnerComponent string
			BusinessNeed           *string
			BusinessCase           *string
		}
	}

	// TODO we're supposed to be able to pass variables as additional arguments using client.Var()
	// but it wasn't working for me.
	s.client.MustPost(fmt.Sprintf(
		`query {
			systemIntake(id: "%s") {
				id
				requestName
				status
				requestType
				businessOwner {
					name
					component
				}
				businessNeed
			}
		}`, intake.ID), &resp)

	s.Equal(intake.ID.String(), resp.SystemIntake.ID)
	s.Equal(projectName, resp.SystemIntake.RequestName)
	s.Equal("INTAKE_SUBMITTED", resp.SystemIntake.Status)
	s.Equal("NEW", resp.SystemIntake.RequestType)
	s.Equal(businessOwner, resp.SystemIntake.BusinessOwner.Name)
	s.Equal(businessOwnerComponent, resp.SystemIntake.BusinessOwner.Component)
	s.Nil(resp.SystemIntake.BusinessNeed)
	s.Nil(resp.SystemIntake.BusinessCase)
}

func (s GraphQLTestSuite) TestFetchSystemIntakeWithNotesQuery() {
	ctx := context.Background()
	projectName := "Big Project"
	businessOwner := "Firstname Lastname"
	businessOwnerComponent := "OIT"

	intake, intakeErr := s.store.CreateSystemIntake(ctx, &models.SystemIntake{
		ProjectName:            null.StringFrom(projectName),
		Status:                 models.SystemIntakeStatusINTAKESUBMITTED,
		RequestType:            models.SystemIntakeRequestTypeNEW,
		BusinessOwner:          null.StringFrom(businessOwner),
		BusinessOwnerComponent: null.StringFrom(businessOwnerComponent),
	})
	s.NoError(intakeErr)

	_, noteErr := s.store.CreateNote(ctx, &models.Note{
		SystemIntakeID: intake.ID,
		AuthorEUAID:    "QQQQ",
		AuthorName:     null.StringFrom("Author Name"),
		Content:        null.StringFrom("a clever remark"),
	})
	s.NoError(noteErr)

	var resp struct {
		SystemIntake struct {
			ID    string
			Notes []struct {
				ID     string
				Author struct {
					Name string
					EUA  string
				}
				Content   string
				CreatedAt string
			}
		}
	}

	// TODO we're supposed to be able to pass variables as additional arguments using client.Var()
	// but it wasn't working for me.
	s.client.MustPost(fmt.Sprintf(
		`query {
			systemIntake(id: "%s") {
				id
				notes {
					id
					createdAt
					content
					author {
						name
						eua
					}
				}
			}
		}`, intake.ID), &resp)

	s.Equal(intake.ID.String(), resp.SystemIntake.ID)

	s.Len(resp.SystemIntake.Notes, 1)

	s.NotEmpty(resp.SystemIntake.Notes[0].ID)
	s.Equal("a clever remark", resp.SystemIntake.Notes[0].Content)
	s.Equal("QQQQ", resp.SystemIntake.Notes[0].Author.EUA)
	s.Equal("Author Name", resp.SystemIntake.Notes[0].Author.Name)
	s.NotEmpty(resp.SystemIntake.Notes[0].CreatedAt)
}

func (s GraphQLTestSuite) TestFetchSystemIntakeWithContractMonthAndYearQuery() {
	ctx := context.Background()
	contracStartMonth := "10"
	contractStartYear := "2002"
	contractEndMonth := "08"
	contractEndYear := "2020"
	projectName := "My cool project"

	intake, intakeErr := s.store.CreateSystemIntake(ctx, &models.SystemIntake{
		ProjectName:        null.StringFrom(projectName),
		Status:             models.SystemIntakeStatusINTAKESUBMITTED,
		RequestType:        models.SystemIntakeRequestTypeNEW,
		ContractStartMonth: null.StringFrom(contracStartMonth),
		ContractStartYear:  null.StringFrom(contractStartYear),
		ContractEndMonth:   null.StringFrom(contractEndMonth),
		ContractEndYear:    null.StringFrom(contractEndYear),
	})
	s.NoError(intakeErr)

	var resp struct {
		SystemIntake struct {
			ID       string
			Contract struct {
				EndDate struct {
					Day   string
					Month string
					Year  string
				}
				StartDate struct {
					Day   string
					Month string
					Year  string
				}
			}
		}
	}

	// TODO we're supposed to be able to pass variables as additional arguments using client.Var()
	// but it wasn't working for me.
	s.client.MustPost(fmt.Sprintf(
		`query {
			systemIntake(id: "%s") {
				id
				contract {
					endDate {
						day
						month
						year
					}
					startDate {
						day
						month
						year
					}
				}
			}
		}`, intake.ID), &resp)

	s.Equal(intake.ID.String(), resp.SystemIntake.ID)

	startDate := resp.SystemIntake.Contract.StartDate
	s.Equal("", startDate.Day)
	s.Equal(contracStartMonth, startDate.Month)
	s.Equal(contractStartYear, startDate.Year)

	endDate := resp.SystemIntake.Contract.EndDate
	s.Equal("", endDate.Day)
	s.Equal(contractEndMonth, endDate.Month)
	s.Equal(contractEndYear, endDate.Year)
}

func (s GraphQLTestSuite) TestFetchSystemIntakeWithContractDatesQuery() {
	ctx := context.Background()
	projectName := "My cool project"
	contractStartDate, _ := time.Parse("2006-1-2", "2002-8-24")
	contractEndDate, _ := time.Parse("2006-1-2", "2020-10-31")

	intake, intakeErr := s.store.CreateSystemIntake(ctx, &models.SystemIntake{
		ProjectName:       null.StringFrom(projectName),
		Status:            models.SystemIntakeStatusINTAKESUBMITTED,
		RequestType:       models.SystemIntakeRequestTypeNEW,
		ContractStartDate: &contractStartDate,
		ContractEndDate:   &contractEndDate,
	})
	s.NoError(intakeErr)

	var resp struct {
		SystemIntake struct {
			ID       string
			Contract struct {
				EndDate struct {
					Day   string
					Month string
					Year  string
				}
				StartDate struct {
					Day   string
					Month string
					Year  string
				}
			}
		}
	}

	// TODO we're supposed to be able to pass variables as additional arguments using client.Var()
	// but it wasn't working for me.
	s.client.MustPost(fmt.Sprintf(
		`query {
			systemIntake(id: "%s") {
				id
				contract {
					endDate {
						day
						month
						year
					}
					startDate {
						day
						month
						year
					}
				}
			}
		}`, intake.ID), &resp)

	s.Equal(intake.ID.String(), resp.SystemIntake.ID)

	startDate := resp.SystemIntake.Contract.StartDate
	s.Equal("24", startDate.Day)
	s.Equal("8", startDate.Month)
	s.Equal("2002", startDate.Year)

	endDate := resp.SystemIntake.Contract.EndDate
	s.Equal("31", endDate.Day)
	s.Equal("10", endDate.Month)
	s.Equal("2020", endDate.Year)
}

func (s GraphQLTestSuite) TestFetchSystemIntakeWithNoCollaboratorsQuery() {
	ctx := context.Background()
	projectName := "My cool project"

	intake, intakeErr := s.store.CreateSystemIntake(ctx, &models.SystemIntake{
		ProjectName:                 null.StringFrom(projectName),
		Status:                      models.SystemIntakeStatusINTAKESUBMITTED,
		RequestType:                 models.SystemIntakeRequestTypeNEW,
		EACollaboratorName:          null.StringFrom(""),
		OITSecurityCollaboratorName: null.StringFrom(""),
		TRBCollaboratorName:         null.StringFrom(""),
	})
	s.NoError(intakeErr)

	var resp struct {
		SystemIntake struct {
			ID              string
			GovernanceTeams struct {
				IsPresent bool
				Teams     []struct {
					Acronym      string
					Collaborator string
					Key          string
					Label        string
					Name         string
				}
			}
		}
	}

	// TODO we're supposed to be able to pass variables as additional arguments using client.Var()
	// but it wasn't working for me.
	s.client.MustPost(fmt.Sprintf(
		`query {
			systemIntake(id: "%s") {
				id
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
			}
		}`, intake.ID), &resp)

	s.Equal(intake.ID.String(), resp.SystemIntake.ID)
	s.False(resp.SystemIntake.GovernanceTeams.IsPresent)
	s.Nil(resp.SystemIntake.GovernanceTeams.Teams)
}

func (s GraphQLTestSuite) TestFetchSystemIntakeWithCollaboratorsQuery() {
	ctx := context.Background()
	projectName := "My cool project"
	eaName := "My EA Rep"
	oitName := "My OIT Rep"
	trbName := "My TRB Rep"

	intake, intakeErr := s.store.CreateSystemIntake(ctx, &models.SystemIntake{
		ProjectName:                 null.StringFrom(projectName),
		Status:                      models.SystemIntakeStatusINTAKESUBMITTED,
		RequestType:                 models.SystemIntakeRequestTypeNEW,
		EACollaboratorName:          null.StringFrom(eaName),
		OITSecurityCollaboratorName: null.StringFrom(oitName),
		TRBCollaboratorName:         null.StringFrom(trbName),
	})
	s.NoError(intakeErr)

	var resp struct {
		SystemIntake struct {
			ID              string
			GovernanceTeams struct {
				IsPresent bool
				Teams     []struct {
					Acronym      string
					Collaborator string
					Key          string
					Label        string
					Name         string
				}
			}
		}
	}

	// TODO we're supposed to be able to pass variables as additional arguments using client.Var()
	// but it wasn't working for me.
	s.client.MustPost(fmt.Sprintf(
		`query {
			systemIntake(id: "%s") {
				id
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
			}
		}`, intake.ID), &resp)

	s.Equal(intake.ID.String(), resp.SystemIntake.ID)
	s.True(resp.SystemIntake.GovernanceTeams.IsPresent)
	s.Equal(eaName, resp.SystemIntake.GovernanceTeams.Teams[0].Collaborator)
	s.Equal(oitName, resp.SystemIntake.GovernanceTeams.Teams[1].Collaborator)
	s.Equal(trbName, resp.SystemIntake.GovernanceTeams.Teams[2].Collaborator)
}

func (s GraphQLTestSuite) TestFetchSystemIntakeWithActionsQuery() {
	ctx := context.Background()

	intake, intakeErr := s.store.CreateSystemIntake(ctx, &models.SystemIntake{
		ProjectName: null.StringFrom("Test Project"),
		Status:      models.SystemIntakeStatusINTAKESUBMITTED,
		RequestType: models.SystemIntakeRequestTypeNEW,
	})
	s.NoError(intakeErr)

	action1, action1Err := s.store.CreateAction(ctx, &models.Action{
		IntakeID:       &intake.ID,
		ActionType:     models.ActionTypeSUBMITINTAKE,
		ActorName:      "First Actor",
		ActorEmail:     "first.actor@example.com",
		ActorEUAUserID: "ACT1",
	})
	s.NoError(action1Err)

	action2, action2Err := s.store.CreateAction(ctx, &models.Action{
		IntakeID:       &intake.ID,
		ActionType:     models.ActionTypePROVIDEFEEDBACKNEEDBIZCASE,
		ActorName:      "Second Actor",
		ActorEmail:     "second.actor@example.com",
		ActorEUAUserID: "ACT2",
		Feedback:       null.StringFrom("feedback for action two"),
	})
	s.NoError(action2Err)

	var resp struct {
		SystemIntake struct {
			ID      string
			Actions []struct {
				ID    string
				Type  string
				Actor struct {
					Name  string
					Email string
				}
				Feedback  *string
				CreatedAt string
			}
		}
	}

	// TODO we're supposed to be able to pass variables as additional arguments using client.Var()
	// but it wasn't working for me.
	s.client.MustPost(fmt.Sprintf(
		`query {
			systemIntake(id: "%s") {
				id
				actions {
					id
					type
					actor {
						name
						email
					}
					feedback
					createdAt
				}
			}
		}`, intake.ID), &resp)

	s.Equal(2, len(resp.SystemIntake.Actions))

	respAction1 := resp.SystemIntake.Actions[0]
	s.Equal(action1.ID.String(), respAction1.ID)
	s.Nil(respAction1.Feedback)
	s.Equal(action1.CreatedAt.UTC().Format(time.RFC3339), respAction1.CreatedAt)
	s.Equal("SUBMIT_INTAKE", respAction1.Type)
	s.Equal("First Actor", respAction1.Actor.Name)
	s.Equal("first.actor@example.com", respAction1.Actor.Email)

	respAction2 := resp.SystemIntake.Actions[1]
	s.Equal(action2.ID.String(), respAction2.ID)
	s.Equal("feedback for action two", *respAction2.Feedback)
	s.Equal(action2.CreatedAt.UTC().Format(time.RFC3339), respAction2.CreatedAt)
	s.Equal("PROVIDE_FEEDBACK_NEED_BIZ_CASE", respAction2.Type)
	s.Equal("Second Actor", respAction2.Actor.Name)
	s.Equal("second.actor@example.com", respAction2.Actor.Email)
}

func (s GraphQLTestSuite) TestIssueLifecycleIDWithPassedLCID() {
	ctx := context.Background()
	projectName := "My cool project"

	intake, intakeErr := s.store.CreateSystemIntake(ctx, &models.SystemIntake{
		ProjectName: null.StringFrom(projectName),
		Status:      models.SystemIntakeStatusINTAKESUBMITTED,
		RequestType: models.SystemIntakeRequestTypeNEW,
	})
	s.NoError(intakeErr)

	var resp struct {
		IssueLifecycleID struct {
			SystemIntake struct {
				ID                string
				Lcid              string
				LcidExpiresAt     string
				LcidScope         string
				DecisionNextSteps string
				Status            string
			}
		}
	}

	// TODO we're supposed to be able to pass variables as additional arguments using client.Var()
	// but it wasn't working for me.
	s.client.MustPost(fmt.Sprintf(
		`mutation {
			issueLifecycleId(input: {
				intakeId: "%s",
				expiresAt: "2021-03-18T00:00:00Z",
				scope: "Your scope",
				feedback: "My feedback",
				lcid: "123456A"
				nextSteps: "Your next steps"
			}) {
				systemIntake {
					id
					lcid
					lcidExpiresAt
					lcidScope
					decisionNextSteps
					status
				}
			}
		}`, intake.ID), &resp)

	s.Equal(intake.ID.String(), resp.IssueLifecycleID.SystemIntake.ID)

	respIntake := resp.IssueLifecycleID.SystemIntake
	s.Equal(respIntake.LcidExpiresAt, "2021-03-18T00:00:00Z")
	s.Equal(respIntake.Lcid, "123456A")
}

func (s GraphQLTestSuite) TestIssueLifecycleIDSetNewLCID() {
	ctx := context.Background()
	projectName := "My cool project"

	intake, intakeErr := s.store.CreateSystemIntake(ctx, &models.SystemIntake{
		ProjectName: null.StringFrom(projectName),
		Status:      models.SystemIntakeStatusINTAKESUBMITTED,
		RequestType: models.SystemIntakeRequestTypeNEW,
	})
	s.NoError(intakeErr)

	var resp struct {
		IssueLifecycleID struct {
			SystemIntake struct {
				ID                string
				Lcid              string
				LcidExpiresAt     string
				LcidScope         string
				DecisionNextSteps string
				Status            string
			}
		}
	}

	// TODO we're supposed to be able to pass variables as additional arguments using client.Var()
	// but it wasn't working for me.
	s.client.MustPost(fmt.Sprintf(
		`mutation {
			issueLifecycleId(input: {
				intakeId: "%s",
				expiresAt: "2021-03-18T00:00:00Z",
				scope: "Your scope",
				feedback: "My feedback",
				lcid: ""
				nextSteps: "Your next steps"
			}) {
				systemIntake {
					id
					lcid
					lcidExpiresAt
					lcidScope
					decisionNextSteps
					status
				}
			}
		}`, intake.ID), &resp)

	s.Equal(intake.ID.String(), resp.IssueLifecycleID.SystemIntake.ID)

	respIntake := resp.IssueLifecycleID.SystemIntake
	s.Equal(respIntake.LcidExpiresAt, "2021-03-18T00:00:00Z")
	s.Equal(respIntake.Lcid, "654321B")
}
