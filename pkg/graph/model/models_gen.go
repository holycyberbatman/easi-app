// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/cmsgov/easi-app/pkg/models"
	"github.com/google/uuid"
)

// Document type of an Accessibility Request document
type AccessibilityRequestDocumentType struct {
	CommonType           models.AccessibilityRequestDocumentCommonType `json:"commonType"`
	OtherTypeDescription *string                                       `json:"otherTypeDescription"`
}

// An edge of an AccessibilityRequestConnection
type AccessibilityRequestEdge struct {
	Cursor string                       `json:"cursor"`
	Node   *models.AccessibilityRequest `json:"node"`
}

// A collection of AccessibilityRequests
type AccessibilityRequestsConnection struct {
	Edges      []*AccessibilityRequestEdge `json:"edges"`
	TotalCount int                         `json:"totalCount"`
}

// Input for adding GRT Feedback
type AddGRTFeedbackInput struct {
	EmailBody string    `json:"emailBody"`
	Feedback  string    `json:"feedback"`
	IntakeID  uuid.UUID `json:"intakeID"`
}

// Response for adding GRT Feedback
type AddGRTFeedbackPayload struct {
	ID *uuid.UUID `json:"id"`
}

// Parameters for actions without additional fields
type BasicActionInput struct {
	Feedback string    `json:"feedback"`
	IntakeID uuid.UUID `json:"intakeId"`
}

// The shape of a solution for a business case
type BusinessCaseAsIsSolution struct {
	Cons        *string `json:"cons"`
	CostSavings *string `json:"costSavings"`
	Pros        *string `json:"pros"`
	Summary     *string `json:"summary"`
	Title       *string `json:"title"`
}

// The shape of a business owner for a business case
type BusinessCaseBusinessOwner struct {
	Name *string `json:"name"`
}

// The shape of a requester for a business case
type BusinessCaseRequester struct {
	Name        *string `json:"name"`
	PhoneNumber *string `json:"phoneNumber"`
}

// The shape of a solution for a business case
type BusinessCaseSolution struct {
	AcquisitionApproach *string                       `json:"acquisitionApproach"`
	Cons                *string                       `json:"cons"`
	CostSavings         *string                       `json:"costSavings"`
	HasUI               *string                       `json:"hasUi"`
	Hosting             *BusinessCaseSolutionHosting  `json:"hosting"`
	Pros                *string                       `json:"pros"`
	Security            *BusinessCaseSolutionSecurity `json:"security"`
	Summary             *string                       `json:"summary"`
	Title               *string                       `json:"title"`
}

// The shape of hosting for a business case solution
type BusinessCaseSolutionHosting struct {
	CloudServiceType *string `json:"cloudServiceType"`
	Location         *string `json:"location"`
	Type             *string `json:"type"`
}

// The shape of security for a business case solution
type BusinessCaseSolutionSecurity struct {
	IsApproved      *bool   `json:"isApproved"`
	IsBeingReviewed *string `json:"isBeingReviewed"`
}

// A date for a contract
type ContractDate struct {
	Day   *string `json:"day"`
	Month *string `json:"month"`
	Year  *string `json:"year"`
}

// Parameters for createAccessibilityRequestDocument
type CreateAccessibilityRequestDocumentInput struct {
	CommonDocumentType           models.AccessibilityRequestDocumentCommonType `json:"commonDocumentType"`
	MimeType                     string                                        `json:"mimeType"`
	Name                         string                                        `json:"name"`
	OtherDocumentTypeDescription *string                                       `json:"otherDocumentTypeDescription"`
	RequestID                    uuid.UUID                                     `json:"requestID"`
	Size                         int                                           `json:"size"`
	URL                          string                                        `json:"url"`
}

// Result of createAccessibilityRequestDocument
type CreateAccessibilityRequestDocumentPayload struct {
	AccessibilityRequestDocument *models.AccessibilityRequestDocument `json:"accessibilityRequestDocument"`
	UserErrors                   []*UserError                         `json:"userErrors"`
}

// Parameters required to create an AccessibilityRequest
type CreateAccessibilityRequestInput struct {
	IntakeID uuid.UUID `json:"intakeID"`
	Name     string    `json:"name"`
}

// Result of CreateAccessibilityRequest
type CreateAccessibilityRequestPayload struct {
	AccessibilityRequest *models.AccessibilityRequest `json:"accessibilityRequest"`
	UserErrors           []*UserError                 `json:"userErrors"`
}

// Parameters for creating a test date
type CreateTestDateInput struct {
	Date      time.Time               `json:"date"`
	RequestID uuid.UUID               `json:"requestID"`
	Score     *int                    `json:"score"`
	TestType  models.TestDateTestType `json:"testType"`
}

// Result of createTestDate
type CreateTestDatePayload struct {
	TestDate   *models.TestDate `json:"testDate"`
	UserErrors []*UserError     `json:"userErrors"`
}

// Parameters required to generate a presigned upload URL
type GeneratePresignedUploadURLInput struct {
	FileName string `json:"fileName"`
	MimeType string `json:"mimeType"`
	Size     int    `json:"size"`
}

// Result of CreateAccessibilityRequest
type GeneratePresignedUploadURLPayload struct {
	URL        *string      `json:"url"`
	UserErrors []*UserError `json:"userErrors"`
}

// Input for issuing a lifecycle id
type IssueLifecycleIDInput struct {
	ExpiresAt time.Time `json:"expiresAt"`
	Feedback  string    `json:"feedback"`
	IntakeID  uuid.UUID `json:"intakeId"`
	Lcid      *string   `json:"lcid"`
	NextSteps *string   `json:"nextSteps"`
	Scope     string    `json:"scope"`
}

// A collection of Systems
type SystemConnection struct {
	Edges      []*SystemEdge `json:"edges"`
	TotalCount int           `json:"totalCount"`
}

// An edge of an SystemConnection
type SystemEdge struct {
	Cursor string         `json:"cursor"`
	Node   *models.System `json:"node"`
}

// A business owner for a system intake
type SystemIntakeBusinessOwner struct {
	Component *string `json:"component"`
	Name      *string `json:"name"`
}

// A collaborator for an intake
type SystemIntakeCollaborator struct {
	Acronym      *string `json:"acronym"`
	Collaborator *string `json:"collaborator"`
	Key          *string `json:"key"`
	Label        *string `json:"label"`
	Name         *string `json:"name"`
}

// A contract for a system intake
type SystemIntakeContract struct {
	Contractor  *string       `json:"contractor"`
	EndDate     *ContractDate `json:"endDate"`
	HasContract *string       `json:"hasContract"`
	StartDate   *ContractDate `json:"startDate"`
	Vehicle     *string       `json:"vehicle"`
}

// costs for a system intake
type SystemIntakeCosts struct {
	ExpectedIncreaseAmount *string `json:"expectedIncreaseAmount"`
	IsExpectingIncrease    *string `json:"isExpectingIncrease"`
}

// A funding source for a system intake
type SystemIntakeFundingSource struct {
	FundingNumber *string `json:"fundingNumber"`
	IsFunded      *bool   `json:"isFunded"`
	Source        *string `json:"source"`
}

// governanceTeam for an intake
type SystemIntakeGovernanceTeam struct {
	IsPresent *bool                       `json:"isPresent"`
	Teams     []*SystemIntakeCollaborator `json:"teams"`
}

// An isso for a system intake
type SystemIntakeIsso struct {
	IsPresent *bool   `json:"isPresent"`
	Name      *string `json:"name"`
}

// A note on a system intake
type SystemIntakeNote struct {
	Author    *SystemIntakeNoteAuthor `json:"author"`
	Content   string                  `json:"content"`
	CreatedAt time.Time               `json:"createdAt"`
	ID        uuid.UUID               `json:"id"`
}

// The author of a system intake note
type SystemIntakeNoteAuthor struct {
	Eua  string `json:"eua"`
	Name string `json:"name"`
}

// A product manager for a system intake
type SystemIntakeProductManager struct {
	Component *string `json:"component"`
	Name      *string `json:"name"`
}

// A requester for a system intake
type SystemIntakeRequester struct {
	Component *string `json:"component"`
	Email     *string `json:"email"`
	Name      string  `json:"name"`
}

// Parameters required to update the admin lead for an intake
type UpdateSystemIntakeAdminLeadInput struct {
	AdminLead string    `json:"adminLead"`
	ID        uuid.UUID `json:"id"`
}

// Result of UpdateSystemIntake mutations
type UpdateSystemIntakePayload struct {
	SystemIntake *models.SystemIntake `json:"systemIntake"`
	UserErrors   []*UserError         `json:"userErrors"`
}

// Parameters required to update the grt and grb dates for an intake
type UpdateSystemIntakeReviewDatesInput struct {
	GrbDate *time.Time `json:"grbDate"`
	GrtDate *time.Time `json:"grtDate"`
	ID      uuid.UUID  `json:"id"`
}

// Parameters for editing a test date
type UpdateTestDateInput struct {
	Date     time.Time               `json:"date"`
	ID       uuid.UUID               `json:"id"`
	Score    *int                    `json:"score"`
	TestType models.TestDateTestType `json:"testType"`
}

// Result of editTestDate
type UpdateTestDatePayload struct {
	TestDate   *models.TestDate `json:"testDate"`
	UserErrors []*UserError     `json:"userErrors"`
}

// UserError represents application-level errors that are the result of
// either user or application developer error.
type UserError struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

// A user role associated with a job code
type Role string

const (
	// A 508 Tester
	RoleEasi508Tester Role = "EASI_508_TESTER"
	// A 508 request owner
	RoleEasi508User Role = "EASI_508_USER"
	// A member of the GRT
	RoleEasiGovteam Role = "EASI_GOVTEAM"
	// A generic EASi user
	RoleEasiUser Role = "EASI_USER"
)

var AllRole = []Role{
	RoleEasi508Tester,
	RoleEasi508User,
	RoleEasiGovteam,
	RoleEasiUser,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleEasi508Tester, RoleEasi508User, RoleEasiGovteam, RoleEasiUser:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
