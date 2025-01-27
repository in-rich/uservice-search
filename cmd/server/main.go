package main

import (
	"fmt"
	"github.com/in-rich/lib-go/deploy"
	"github.com/in-rich/lib-go/monitor"
	search_pb "github.com/in-rich/proto/proto-go/search"
	"github.com/in-rich/uservice-search/config"
	"github.com/in-rich/uservice-search/migrations"
	"github.com/in-rich/uservice-search/pkg/dao"
	"github.com/in-rich/uservice-search/pkg/handlers"
	"github.com/in-rich/uservice-search/pkg/services"
	"github.com/rs/zerolog"
	"os"
)

func getLogger() monitor.GRPCLogger {
	if deploy.IsReleaseEnv() {
		return monitor.NewGCPGRPCLogger(zerolog.New(os.Stdout), "uservice-search")
	}

	return monitor.NewConsoleGRPCLogger()
}

func main() {
	logger := getLogger()

	logger.Info("Starting server")
	db, closeDB, err := deploy.OpenDB(config.App.Postgres.DSN)
	if err != nil {
		logger.Fatal(err, "failed to connect to database")
	}
	defer closeDB()

	logger.Info("Running migrations")
	if err := migrations.Migrate(db); err != nil {
		logger.Fatal(err, "failed to migrate")
	}

	depCheck := deploy.DepsCheck{
		Dependencies: func() map[string]error {
			return map[string]error{
				"Postgres": db.Ping(),
			}
		},
		Services: deploy.DepCheckServices{
			"SearchNotes":          {"Postgres"},
			"UpsertNotes":          {"Postgres"},
			"SearchMessages":       {"Postgres"},
			"UpsertMessages":       {"Postgres"},
			"CreateTeamMeta":       {"Postgres"},
			"DeleteTeamMetaMember": {"Postgres"},
			"SearchReminders":      {"Postgres"},
			"UpsertReminders":      {"Postgres"},
		},
	}

	createNoteDAO := dao.NewCreateNoteRepository(db)
	deleteNoteDAO := dao.NewDeleteNoteRepository(db)
	searchNotesDAO := dao.NewSearchNotesRepository(db)
	updateNoteDAO := dao.NewUpdateNoteRepository(db)
	createMessageDAO := dao.NewCreateMessageRepository(db)
	deleteMessageDAO := dao.NewDeleteMessageRepository(db)
	searchMessageDAO := dao.NewSearchMessagesRepository(db)
	updateMessageDAO := dao.NewUpdateMessageRepository(db)
	createTeamMetaDAO := dao.NewCreateTeamMetaRepository(db)
	deleteTeamMetaDAO := dao.NewDeleteTeamMetaRepository(db)
	deleteTeamMetaMember := dao.NewDeleteTeamMetaMemberRepository(db)
	createReminderDAO := dao.NewCreateReminderRepository(db)
	deleteReminderDAO := dao.NewDeleteReminderRepository(db)
	searchRemindersDAO := dao.NewSearchRemindersRepository(db)
	updateReminderDAO := dao.NewUpdateReminderRepository(db)

	searchNotesService := services.NewSearchNotesService(searchNotesDAO)
	upsertNoteService := services.NewUpsertNoteService(updateNoteDAO, createNoteDAO, deleteNoteDAO)
	searchMessagesService := services.NewSearchMessagesService(searchMessageDAO)
	upsertMessageService := services.NewUpsertMessageService(updateMessageDAO, createMessageDAO, deleteMessageDAO)
	createTeamMetaService := services.NewCreateTeamMetaService(createTeamMetaDAO)
	deleteTeamMetaService := services.NewDeleteTeamMetaService(deleteTeamMetaDAO, deleteTeamMetaMember)
	searchRemindersService := services.NewSearchRemindersService(searchRemindersDAO)
	upsertReminderService := services.NewUpsertReminderService(updateReminderDAO, createReminderDAO, deleteReminderDAO)

	searchNotesHandler := handlers.NewSearchNotesHandler(searchNotesService)
	upsertNoteHandler := handlers.NewUpsertNoteHandler(upsertNoteService)
	searchMessagesHandler := handlers.NewSearchMessagesHandler(searchMessagesService)
	upsertMessageHandler := handlers.NewUpsertMessageHandler(upsertMessageService)
	createTeamMetaHandler := handlers.NewCreateTeamMetaHandler(createTeamMetaService)
	deleteTeamMetaHandler := handlers.NewDeleteTeamMetaHandler(deleteTeamMetaService)
	searchRemindersHandler := handlers.NewSearchRemindersHandler(searchRemindersService)
	upsertReminderHandler := handlers.NewUpsertReminderHandler(upsertReminderService)

	logger.Info(fmt.Sprintf("Starting to listen on port %v", config.App.Server.Port))
	listener, server, health := deploy.StartGRPCServer(logger, config.App.Server.Port, depCheck)
	defer deploy.CloseGRPCServer(listener, server)
	go health()

	search_pb.RegisterSearchNotesServer(server, searchNotesHandler)
	search_pb.RegisterUpsertNoteServer(server, upsertNoteHandler)
	search_pb.RegisterSearchMessagesServer(server, searchMessagesHandler)
	search_pb.RegisterUpsertMessageServer(server, upsertMessageHandler)
	search_pb.RegisterCreateTeamMetaServer(server, createTeamMetaHandler)
	search_pb.RegisterDeleteTeamMetaServer(server, deleteTeamMetaHandler)
	search_pb.RegisterSearchRemindersServer(server, searchRemindersHandler)
	search_pb.RegisterUpsertReminderServer(server, upsertReminderHandler)

	logger.Info("Server started")
	if err := server.Serve(listener); err != nil {
		logger.Fatal(err, "failed to serve")
	}
}
