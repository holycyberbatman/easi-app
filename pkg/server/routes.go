package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // pq is required to get the postgres driver into sqlx
	"go.uber.org/zap"

	"github.com/cmsgov/easi-app/pkg/appconfig"
	"github.com/cmsgov/easi-app/pkg/appses"
	"github.com/cmsgov/easi-app/pkg/appvalidation"
	"github.com/cmsgov/easi-app/pkg/cedar/cedareasi"
	"github.com/cmsgov/easi-app/pkg/cedar/cedarldap"
	cedarintake "github.com/cmsgov/easi-app/pkg/cedar/intake"
	"github.com/cmsgov/easi-app/pkg/email"
	"github.com/cmsgov/easi-app/pkg/flags"
	"github.com/cmsgov/easi-app/pkg/graph"
	"github.com/cmsgov/easi-app/pkg/graph/generated"
	"github.com/cmsgov/easi-app/pkg/graph/model"
	"github.com/cmsgov/easi-app/pkg/handlers"
	"github.com/cmsgov/easi-app/pkg/local"
	"github.com/cmsgov/easi-app/pkg/models"
	"github.com/cmsgov/easi-app/pkg/services"
	"github.com/cmsgov/easi-app/pkg/storage"
	"github.com/cmsgov/easi-app/pkg/upload"
)

func (s *Server) routes(
	authorizationMiddleware func(handler http.Handler) http.Handler,
	corsMiddleware func(handler http.Handler) http.Handler,
	traceMiddleware func(handler http.Handler) http.Handler,
	loggerMiddleware func(handler http.Handler) http.Handler) {

	s.router.Use(
		traceMiddleware, // trace all requests with an ID
		loggerMiddleware,
		corsMiddleware,
		// authorizationMiddleware, // is supposed to be authN, not authZ; TODO: those responsibilities should be split out
	)

	// set up handler base
	base := handlers.NewHandlerBase(s.logger)

	// endpoints that dont require authorization go directly on the main router
	s.router.HandleFunc("/api/v1/healthcheck", handlers.NewHealthCheckHandler(base, s.Config).Handle())
	s.router.HandleFunc("/api/graph/playground", playground.Handler("GraphQL playground", "/api/graph/query"))

	// set up Feature Flagging utilities
	ldClient, err := flags.NewLaunchDarklyClient(s.NewFlagConfig())
	if err != nil {
		s.logger.Fatal("Failed to create LaunchDarkly client", zap.Error(err))
	}

	// set up CEDAR client
	publisher := cedarintake.NewClient(
		s.Config.GetString(appconfig.CEDARAPIURL),
		s.Config.GetString(appconfig.CEDARAPIKey),
		ldClient,
	)
	if s.environment.Deployed() {
		s.NewCEDARClientCheck()
		if cerr := publisher.CheckConnection(context.Background()); cerr != nil {
			s.logger.Info("Non-Fatal - Failed to connect to CEDAR Intake API on startup", zap.Error(cerr))
		}
	}

	var cedarEasiClient cedareasi.Client = local.NewCedarEasiClient()
	if !(s.environment.Local() || s.environment.Test()) {
		// check we have all of the configs for CEDAR clients
		s.NewCEDARClientCheck()

		cedarEasiClient = cedareasi.NewTranslatedClient(
			s.Config.GetString(appconfig.CEDARAPIURL),
			s.Config.GetString(appconfig.CEDARAPIKey),
			ldClient,
		)
		if s.environment.Deployed() {
			s.CheckCEDAREasiClientConnection(cedarEasiClient)
		}
	}
	_ = cedarEasiClient

	var cedarLDAPClient cedarldap.Client
	cedarLDAPClient = cedarldap.NewTranslatedClient(
		s.Config.GetString(appconfig.CEDARAPIURL),
		s.Config.GetString(appconfig.CEDARAPIKey),
	)
	if s.environment.Local() || s.environment.Test() {
		cedarLDAPClient = local.NewCedarLdapClient(s.logger)
	}

	// set up Email Client
	sesConfig := s.NewSESConfig()
	sesSender := appses.NewSender(sesConfig)
	emailConfig := s.NewEmailConfig()
	emailClient, err := email.NewClient(emailConfig, sesSender)
	if err != nil {
		s.logger.Fatal("Failed to create email client", zap.Error(err))
	}
	// override email client with local one
	if s.environment.Local() || s.environment.Test() {
		localSender := local.NewSender()
		emailClient, err = email.NewClient(emailConfig, localSender)
		if err != nil {
			s.logger.Fatal("Failed to create email client", zap.Error(err))
		}
	}

	if s.environment.Deployed() {
		s.CheckEmailClient(emailClient)
	}

	// set up S3 client
	s3Config := s.NewS3Config()
	if s.environment.Local() {
		s3Config.IsLocal = true
	}

	s3Client := upload.NewS3Client(s3Config)

	var lambdaClient *lambda.Lambda
	var princeLambdaName string
	lambdaSession := session.Must(session.NewSession())

	princeConfig := s.NewPrinceLambdaConfig()
	princeLambdaName = princeConfig.FunctionName

	if s.environment.Local() || s.environment.Test() {
		endpoint := princeConfig.Endpoint
		lambdaClient = lambda.New(lambdaSession, &aws.Config{Endpoint: &endpoint, Region: aws.String("us-west-2")})
	} else {
		lambdaClient = lambda.New(lambdaSession, &aws.Config{})
	}

	store, storeErr := storage.NewStore(
		s.logger,
		s.NewDBConfig(),
		ldClient,
	)
	if storeErr != nil {
		s.logger.Fatal("Failed to create store", zap.Error(storeErr))
	}

	serviceConfig := services.NewConfig(s.logger, ldClient)

	// set up GraphQL routes
	gql := s.router.PathPrefix("/api/graph").Subrouter()
	gql.Use(authorizationMiddleware) // TODO: see comment at top-level router

	saveAction := services.NewSaveAction(
		store.CreateAction,
		cedarLDAPClient.FetchUserInfo,
	)

	resolver := graph.NewResolver(
		store,
		graph.ResolverService{
			CreateTestDate: services.NewCreateTestDate(
				serviceConfig,
				services.NewAuthorizeHasEASiRole(),
				store.CreateTestDate,
			),
			AddGRTFeedback: services.NewProvideGRTFeedback(
				serviceConfig,
				store.FetchSystemIntakeByID,
				store.UpdateSystemIntake,
				saveAction,
				store.CreateGRTFeedback,
				cedarLDAPClient.FetchUserInfo,
				emailClient.SendSystemIntakeReviewEmail,
			),
			CreateActionUpdateStatus: services.NewCreateActionUpdateStatus(
				serviceConfig,
				store.UpdateSystemIntakeStatus,
				saveAction,
				cedarLDAPClient.FetchUserInfo,
				emailClient.SendSystemIntakeReviewEmail,
				services.NewCloseBusinessCase(
					serviceConfig,
					store.FetchBusinessCaseByID,
					store.UpdateBusinessCase,
				),
			),
			IssueLifecycleID: services.NewUpdateLifecycleFields(
				serviceConfig,
				services.NewAuthorizeRequireGRTJobCode(),
				store.FetchSystemIntakeByID,
				store.UpdateSystemIntake,
				saveAction,
				cedarLDAPClient.FetchUserInfo,
				emailClient.SendIssueLCIDEmail,
				store.GenerateLifecycleID,
			),
			AuthorizeUserIsReviewTeamOrIntakeRequester: services.AuthorizeUserIsIntakeRequesterOrHasGRTJobCode,
			RejectIntake: services.NewUpdateRejectionFields(
				serviceConfig,
				services.NewAuthorizeRequireGRTJobCode(),
				store.FetchSystemIntakeByID,
				store.UpdateSystemIntake,
				saveAction,
				cedarLDAPClient.FetchUserInfo,
				emailClient.SendRejectRequestEmail,
			),
		},
		&s3Client,
	)
	gqlDirectives := generated.DirectiveRoot{HasRole: func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (res interface{}, err error) {
		hasRole, err := services.HasRole(ctx, role)
		if err != nil {
			return nil, err
		}
		if !hasRole {
			return nil, errors.New("not authorized")
		}
		return next(ctx)
	}}
	gqlConfig := generated.Config{Resolvers: resolver, Directives: gqlDirectives}
	graphqlServer := handler.NewDefaultServer(generated.NewExecutableSchema(gqlConfig))
	gql.Handle("/query", graphqlServer)

	// API base path is versioned
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.Use(authorizationMiddleware) // TODO: see comment at top-level router

	systemIntakeHandler := handlers.NewSystemIntakeHandler(
		base,
		services.NewCreateSystemIntake(
			serviceConfig,
			store.CreateSystemIntake,
		),
		services.NewUpdateSystemIntake(
			serviceConfig,
			store.FetchSystemIntakeByID,
			store.UpdateSystemIntake,
			services.AuthorizeUserIsIntakeRequesterOrHasGRTJobCode,
		),
		services.NewFetchSystemIntakeByID(
			serviceConfig,
			store.FetchSystemIntakeByID,
			services.NewAuthorizeHasEASiRole(),
		),
		services.NewArchiveSystemIntake(
			serviceConfig,
			store.FetchSystemIntakeByID,
			store.UpdateSystemIntake,
			services.NewCloseBusinessCase(
				serviceConfig,
				store.FetchBusinessCaseByID,
				store.UpdateBusinessCase,
			),
			services.NewAuthorizeUserIsIntakeRequester(),
			emailClient.SendWithdrawRequestEmail,
		),
	)
	api.Handle("/system_intake/{intake_id}", systemIntakeHandler.Handle())
	api.Handle("/system_intake", systemIntakeHandler.Handle())

	systemIntakesHandler := handlers.NewSystemIntakesHandler(
		base,
		services.NewFetchSystemIntakes(
			serviceConfig,
			store.FetchSystemIntakesByEuaID,
			store.FetchSystemIntakes,
			store.FetchSystemIntakesByStatuses,
			services.NewAuthorizeHasEASiRole(),
		),
	)
	api.Handle("/system_intakes", systemIntakesHandler.Handle())

	businessCaseHandler := handlers.NewBusinessCaseHandler(
		base,
		services.NewFetchBusinessCaseByID(
			serviceConfig,
			store.FetchBusinessCaseByID,
			services.NewAuthorizeHasEASiRole(),
		),
		services.NewCreateBusinessCase(
			serviceConfig,
			store.FetchSystemIntakeByID,
			services.NewAuthorizeUserIsIntakeRequester(),
			store.CreateAction,
			cedarLDAPClient.FetchUserInfo,
			store.CreateBusinessCase,
			store.UpdateSystemIntake,
		),
		services.NewUpdateBusinessCase(
			serviceConfig,
			store.FetchBusinessCaseByID,
			services.NewAuthorizeUserIsBusinessCaseRequester(),
			store.UpdateBusinessCase,
		),
	)
	api.Handle("/business_case/{business_case_id}", businessCaseHandler.Handle())
	api.Handle("/business_case", businessCaseHandler.Handle())

	businessCasesHandler := handlers.NewBusinessCasesHandler(
		base,
		services.NewFetchBusinessCasesByEuaID(
			serviceConfig,
			store.FetchBusinessCasesByEuaID,
			services.NewAuthorizeHasEASiRole(),
		),
	)
	api.Handle("/business_cases", businessCasesHandler.Handle())

	metricsHandler := handlers.NewMetricsHandler(
		base,
		services.NewFetchMetrics(serviceConfig, store.FetchSystemIntakeMetrics),
	)
	api.Handle("/metrics", metricsHandler.Handle())

	actionHandler := handlers.NewActionHandler(
		base,
		services.NewTakeAction(
			store.FetchSystemIntakeByID,
			map[models.ActionType]services.ActionExecuter{
				models.ActionTypeSUBMITINTAKE: services.NewSubmitSystemIntake(
					serviceConfig,
					services.NewAuthorizeUserIsIntakeRequester(),
					store.UpdateSystemIntake,
					func(c context.Context, si *models.SystemIntake) (string, error) {
						// quick adapter to retrofit the new interface to take the place
						// of the old interface
						err := publisher.PublishSnapshot(c, si, nil, nil, nil, nil)
						return "", err
					},
					saveAction,
					emailClient.SendSystemIntakeSubmissionEmail,
				),
				models.ActionTypeNOTITREQUEST: services.NewTakeActionUpdateStatus(
					serviceConfig,
					models.SystemIntakeStatusNOTITREQUEST,
					store.UpdateSystemIntake,
					services.NewAuthorizeRequireGRTJobCode(),
					saveAction,
					cedarLDAPClient.FetchUserInfo,
					emailClient.SendSystemIntakeReviewEmail,
					true,
					services.NewCloseBusinessCase(
						serviceConfig,
						store.FetchBusinessCaseByID,
						store.UpdateBusinessCase,
					),
				),
				models.ActionTypeNEEDBIZCASE: services.NewTakeActionUpdateStatus(
					serviceConfig,
					models.SystemIntakeStatusNEEDBIZCASE,
					store.UpdateSystemIntake,
					services.NewAuthorizeRequireGRTJobCode(),
					saveAction,
					cedarLDAPClient.FetchUserInfo,
					emailClient.SendSystemIntakeReviewEmail,
					false,
					services.NewCloseBusinessCase(
						serviceConfig,
						store.FetchBusinessCaseByID,
						store.UpdateBusinessCase,
					),
				),
				models.ActionTypeREADYFORGRT: services.NewTakeActionUpdateStatus(
					serviceConfig,
					models.SystemIntakeStatusREADYFORGRT,
					store.UpdateSystemIntake,
					services.NewAuthorizeRequireGRTJobCode(),
					saveAction,
					cedarLDAPClient.FetchUserInfo,
					emailClient.SendSystemIntakeReviewEmail,
					false,
					services.NewCloseBusinessCase(
						serviceConfig,
						store.FetchBusinessCaseByID,
						store.UpdateBusinessCase,
					),
				),
				models.ActionTypeSUBMITBIZCASE: services.NewSubmitBusinessCase(
					serviceConfig,
					services.NewAuthorizeUserIsIntakeRequester(),
					store.FetchOpenBusinessCaseByIntakeID,
					appvalidation.BusinessCaseForSubmit,
					saveAction,
					store.UpdateSystemIntake,
					store.UpdateBusinessCase,
					emailClient.SendBusinessCaseSubmissionEmail,
					models.SystemIntakeStatusBIZCASEDRAFTSUBMITTED,
				),
				models.ActionTypeSUBMITFINALBIZCASE: services.NewSubmitBusinessCase(
					serviceConfig,
					services.NewAuthorizeUserIsIntakeRequester(),
					store.FetchOpenBusinessCaseByIntakeID,
					appvalidation.BusinessCaseForSubmit,
					saveAction,
					store.UpdateSystemIntake,
					store.UpdateBusinessCase,
					emailClient.SendBusinessCaseSubmissionEmail,
					models.SystemIntakeStatusBIZCASEFINALSUBMITTED,
				),
				models.ActionTypeBIZCASENEEDSCHANGES: services.NewTakeActionUpdateStatus(
					serviceConfig,
					models.SystemIntakeStatusBIZCASECHANGESNEEDED,
					store.UpdateSystemIntake,
					services.NewAuthorizeRequireGRTJobCode(),
					saveAction,
					cedarLDAPClient.FetchUserInfo,
					emailClient.SendSystemIntakeReviewEmail,
					false,
					services.NewCloseBusinessCase(
						serviceConfig,
						store.FetchBusinessCaseByID,
						store.UpdateBusinessCase,
					),
				),
				models.ActionTypeNOGOVERNANCENEEDED: services.NewTakeActionUpdateStatus(
					serviceConfig,
					models.SystemIntakeStatusNOGOVERNANCE,
					store.UpdateSystemIntake,
					services.NewAuthorizeRequireGRTJobCode(),
					saveAction,
					cedarLDAPClient.FetchUserInfo,
					emailClient.SendSystemIntakeReviewEmail,
					true,
					services.NewCloseBusinessCase(
						serviceConfig,
						store.FetchBusinessCaseByID,
						store.UpdateBusinessCase,
					),
				),
				models.ActionTypeSENDEMAIL: services.NewTakeActionUpdateStatus(
					serviceConfig,
					models.SystemIntakeStatusSHUTDOWNINPROGRESS,
					store.UpdateSystemIntake,
					services.NewAuthorizeRequireGRTJobCode(),
					saveAction,
					cedarLDAPClient.FetchUserInfo,
					emailClient.SendSystemIntakeReviewEmail,
					false,
					services.NewCloseBusinessCase(
						serviceConfig,
						store.FetchBusinessCaseByID,
						store.UpdateBusinessCase,
					),
				),
				models.ActionTypeGUIDERECEIVEDCLOSE: services.NewTakeActionUpdateStatus(
					serviceConfig,
					models.SystemIntakeStatusSHUTDOWNCOMPLETE,
					store.UpdateSystemIntake,
					services.NewAuthorizeRequireGRTJobCode(),
					saveAction,
					cedarLDAPClient.FetchUserInfo,
					emailClient.SendSystemIntakeReviewEmail,
					true,
					services.NewCloseBusinessCase(
						serviceConfig,
						store.FetchBusinessCaseByID,
						store.UpdateBusinessCase,
					),
				),
				models.ActionTypeNOTRESPONDINGCLOSE: services.NewTakeActionUpdateStatus(
					serviceConfig,
					models.SystemIntakeStatusNOGOVERNANCE,
					store.UpdateSystemIntake,
					services.NewAuthorizeRequireGRTJobCode(),
					saveAction,
					cedarLDAPClient.FetchUserInfo,
					emailClient.SendSystemIntakeReviewEmail,
					true,
					services.NewCloseBusinessCase(
						serviceConfig,
						store.FetchBusinessCaseByID,
						store.UpdateBusinessCase,
					),
				),
			},
		),
		services.NewFetchActionsByRequestID(
			services.NewAuthorizeRequireGRTJobCode(),
			store.GetActionsByRequestID,
		),
	)
	api.Handle("/system_intake/{intake_id}/actions", actionHandler.Handle())

	systemIntakeLifecycleIDHandler := handlers.NewSystemIntakeLifecycleIDHandler(
		base,
		services.NewUpdateLifecycleFields(
			serviceConfig,
			services.NewAuthorizeRequireGRTJobCode(),
			store.FetchSystemIntakeByID,
			store.UpdateSystemIntake,
			saveAction,
			cedarLDAPClient.FetchUserInfo,
			emailClient.SendIssueLCIDEmail,
			store.GenerateLifecycleID,
		),
	)
	api.Handle("/system_intake/{intake_id}/lcid", systemIntakeLifecycleIDHandler.Handle())

	systemIntakeRejectionHandler := handlers.NewSystemIntakeRejectionHandler(
		base,
		services.NewUpdateRejectionFields(
			serviceConfig,
			services.NewAuthorizeRequireGRTJobCode(),
			store.FetchSystemIntakeByID,
			store.UpdateSystemIntake,
			saveAction,
			cedarLDAPClient.FetchUserInfo,
			emailClient.SendRejectRequestEmail,
		),
	)
	api.Handle("/system_intake/{intake_id}/reject", systemIntakeRejectionHandler.Handle())

	notesHandler := handlers.NewNotesHandler(
		base,
		services.NewFetchNotes(
			serviceConfig,
			store.FetchNotesBySystemIntakeID,
			services.NewAuthorizeRequireGRTJobCode(),
		),
		services.NewCreateNote(
			serviceConfig,
			store.CreateNote,
			services.NewAuthorizeRequireGRTJobCode(),
		),
	)
	api.Handle("/system_intake/{intake_id}/notes", notesHandler.Handle())

	// File Upload Handlers
	fileUploadHandler := handlers.NewFileUploadHandler(
		base,
		services.NewCreateAccessibilityRequestDocument(
			serviceConfig,
			services.NewAuthorizeRequireGRTJobCode(),
			store.CreateAccessibilityRequestDocument),
		services.NewFetchAccessibilityRequestDocument(
			serviceConfig,
			services.NewAuthorizeRequireGRTJobCode(),
			store.FetchAccessibilityRequestDocumentByID),
	)
	api.Handle("/file_uploads", fileUploadHandler.Handle())

	presignedURLUploadHandler := handlers.NewPresignedURLUploadHandler(
		base,
		services.NewCreateFileUploadURL(
			serviceConfig,
			services.NewAuthorizeRequireGRTJobCode(),
			s3Client,
		),
	)
	api.Handle("/file_uploads/upload_url", presignedURLUploadHandler.Handle())

	presignedURLDownloadHandler := handlers.NewPresignedURLDownloadHandler(
		base,
		services.NewCreateFileDownloadURL(
			serviceConfig,
			services.NewAuthorizeRequireGRTJobCode(),
			s3Client,
		),
	)

	api.Handle("/file_uploads/{file_name}/download_url", presignedURLDownloadHandler.Handle())

	s.router.PathPrefix("/").Handler(handlers.NewCatchAllHandler(
		base,
	).Handle())

	api.Handle("/pdf/generate", handlers.NewPDFHandler(services.NewInvokeGeneratePDF(serviceConfig, lambdaClient, princeLambdaName)).Handle())

	systemsHandler := handlers.NewSystemsHandler(
		base,
		services.NewFetchSystems(
			serviceConfig,
			store.ListSystems,
			services.NewAuthorizeHasEASiRole(),
		),
	)
	api.Handle("/systems", systemsHandler.Handle())

	if ok, _ := strconv.ParseBool(os.Getenv("DEBUG_ROUTES")); ok {
		// useful for debugging route issues
		_ = s.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			pathTemplate, err := route.GetPathTemplate()
			if err == nil {
				fmt.Println("ROUTE:", pathTemplate)
			}
			pathRegexp, err := route.GetPathRegexp()
			if err == nil {
				fmt.Println("Path regexp:", pathRegexp)
			}
			queriesTemplates, err := route.GetQueriesTemplates()
			if err == nil {
				fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
			}
			queriesRegexps, err := route.GetQueriesRegexp()
			if err == nil {
				fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
			}
			methods, err := route.GetMethods()
			if err == nil {
				fmt.Println("Methods:", strings.Join(methods, ","))
			}
			fmt.Println()
			return nil
		})
	}
}
