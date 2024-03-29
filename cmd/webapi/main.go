/*
Webapi is the executable for the main web server.
It builds a web server around APIs from `service/api`.
Webapi connects to external resources needed (database) and starts two web servers: the API web server, and the debug.
Everything is served via the API web server, except debug variables (/debug/vars) and profiler infos (pprof).

Usage:

	webapi [flags]

Flags and configurations are handled automatically by the code in `load-configuration.go`.

Return values (exit codes):

	0
		The program ended successfully (no errors, stopped by signal)

	> 0
		The program ended due to an error

Note that this program will update the schema of the database to the latest version available (embedded in the
executable during the build).
*/
package main

import (
	"context"
	"database/sql"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wasa/service/api"
	"wasa/service/database"
	"wasa/service/globaltime"
	"wasa/service/model"
)

// main is the program entry point. The only purpose of this function is to call run() and set the exit code if there is
// any error

func main() {

	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error: ", err)
		os.Exit(1)
	}
}

// run executes the program. The body of this function should perform the following steps:
// * reads the configuration
// * creates and configure the logger
// * connects to any external resources (like databases, authenticators, etc.)
// * creates an instance of the service/api package
// * starts the principal web server (using the service/api.Router.Handler() for HTTP handlers)
// * waits for any termination event: SIGTERM signal (UNIX), non-recoverable server error, etc.
// * closes the principal web server
func run() error {
	rand.Seed(globaltime.Now().UnixNano())
	// Load Configuration and defaults
	cfg, err := loadConfiguration()
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			return nil
		}
		return err
	}

	// Init logging
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(true)
	if cfg.Debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.Infof("application initializing")

	// Start Database
	logger.Println("initializing database support")
	dbconn, err := sql.Open("sqlite3", cfg.DB.Filename)
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = dbconn.Close()
	}()
	db, err := database.New(dbconn)
	if err != nil {
		logger.WithError(err).Error("error creating AppDatabase")
		return fmt.Errorf("creating AppDatabase: %w", err)
	}

	memdb, errm := database.NewMem()
	if errm != nil {
		logger.WithError(err).Error("error creating MemoryAppDatabase")
		return fmt.Errorf("creating MemoryAppDatabase: %w", errm)
	}

	util := api.NewUtil(logger)
	// Start (main) API server
	logger.Info("initializing API server")
	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Create the API router
	apirouter, err := api.New(api.Config{
		Logger:      logger,
		Database:    db,
		MemDatabase: memdb,
		Util:        util,
	})

	go populateMockData(memdb, logger)

	if err != nil {
		logger.WithError(err).Error("error creating the API server instance")
		return fmt.Errorf("creating the API server instance: %w", err)
	}
	router := apirouter.Handler()

	router, err = registerWebUI(router)
	if err != nil {
		logger.WithError(err).Error("error registering web UI handler")
		return fmt.Errorf("registering web UI handler: %w", err)
	}

	// Apply CORS policy
	router = applyCORSHandler(router)

	// Create the API server
	apiserver := http.Server{
		Addr:              cfg.Web.APIHost,
		Handler:           router,
		ReadTimeout:       cfg.Web.ReadTimeout,
		ReadHeaderTimeout: cfg.Web.ReadTimeout,
		WriteTimeout:      cfg.Web.WriteTimeout,
	}

	// Start the service listening for requests in a separate goroutine
	go func() {
		logger.Infof("API listening on %s", apiserver.Addr)
		serverErrors <- apiserver.ListenAndServe()
		logger.Infof("stopping API server")
	}()

	// Waiting for shutdown signal or POSIX signals
	select {
	case err := <-serverErrors:
		// Non-recoverable server error
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		logger.Infof("signal %v received, start shutdown", sig)

		// Asking API server to shut down and load shed.
		err := apirouter.Close()
		if err != nil {
			logger.WithError(err).Warning("graceful shutdown of apirouter error")
		}

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shut down and load shed.
		err = apiserver.Shutdown(ctx)
		if err != nil {
			logger.WithError(err).Warning("error during graceful shutdown of HTTP server")
			err = apiserver.Close()
		}

		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGTERM: // SIGSTOP
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func populateMockData(memdb database.AppDatabaseMemory, logger *logrus.Logger) {

	// MOCK DATA

	photo1b64 := "iVBORw0KGgoAAAANSUhEUgAAAAsAAAALCAYAAACprHcmAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsQAAA7EAZUrDhsAAABdSURBVChTY/wPBAxI4J2sKpTFwCD0+DaUBQFMUJoowPhWRgVuMsgkdJOR+SgmgyRgVqMrBAEmdHfhAiB1JHkQrBibO9GdAeIzobsLFwCpQ/Egsmkw05EBrSKFgQEAciA8LnIAvlsAAAAASUVORK5CYII="
	photo1Dec, erroremock := b64.StdEncoding.DecodeString(photo1b64)
	if erroremock != nil {
		logger.Error("error: ", erroremock)
	}
	photo2b64 := "iVBORw0KGgoAAAANSUhEUgAAAAsAAAALCAYAAACprHcmAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsQAAA7EAZUrDhsAAABVSURBVChTtZHBCgAgCEOn///P1STDRkGX3iETdU2yNkDBzOYNkBJ8xicos8appMo135RZyKe1kbj6uhHD43heMJpPPjMmzF193WDftmBVS/XKr08BOpMPOwHM/fW8AAAAAElFTkSuQmCC"
	photo2Dec, erroremock2 := b64.StdEncoding.DecodeString(photo2b64)
	if erroremock2 != nil {
		logger.Error("error: ", erroremock2)
	}

	veniero := new(model.User)
	veniero.Username = "veniero2"

	pippo := new(model.User)
	pippo.Username = "pippo"

	topolino := new(model.User)
	topolino.Username = "topolino"

	paperino := new(model.User)
	paperino.Username = "paperino"

	memdb.SaveUser(veniero.Username)
	memdb.SaveUser(pippo.Username)
	memdb.SaveUser(topolino.Username)
	memdb.SaveUser(paperino.Username)

	photo1 := memdb.SavePhoto(veniero.Username, photo1Dec)
	time.Sleep(1 * time.Second)
	photo2 := memdb.SavePhoto(pippo.Username, photo2Dec)
	time.Sleep(1 * time.Second)
	memdb.SavePhoto(topolino.Username, photo2Dec)
	time.Sleep(1 * time.Second)
	photo3 := memdb.SavePhoto(pippo.Username, photo1Dec)
	time.Sleep(1 * time.Second)
	memdb.SavePhoto(paperino.Username, photo2Dec)

	memdb.SaveFollow(model.FollowRequest{
		User:     veniero,
		Followee: pippo,
	})

	memdb.SaveFollow(model.FollowRequest{
		User:     veniero,
		Followee: paperino,
	})

	memdb.SaveFollow(model.FollowRequest{
		User:     veniero,
		Followee: topolino,
	})

	memdb.SaveFollow(model.FollowRequest{
		User:     pippo,
		Followee: veniero,
	})

	memdb.SaveFollow(model.FollowRequest{
		User:     topolino,
		Followee: veniero,
	})

	memdb.SaveBan(model.BanRequest{
		User:   veniero,
		Banned: paperino,
	})

	memdb.SaveComment(model.CommentRequest{
		User:    veniero,
		PhotoId: photo1.Id,
		Text:    "Che bel fiore!",
	})

	memdb.SaveComment(model.CommentRequest{
		User:    pippo,
		PhotoId: photo1.Id,
		Text:    "Congratulazioni",
	})

	memdb.SaveLike(model.LikeRequest{
		User:    veniero,
		PhotoId: photo1.Id,
	})

	memdb.SaveLike(model.LikeRequest{
		User:    veniero,
		PhotoId: photo2.Id,
	})
	memdb.SaveLike(model.LikeRequest{
		User:    pippo,
		PhotoId: photo3.Id,
	})

	memdb.SaveLike(model.LikeRequest{
		User:    topolino,
		PhotoId: photo3.Id,
	})

	_, errUpdateUsername := memdb.UpdateUsername("_", veniero.Username, "veniero")

	if errUpdateUsername != nil {
		logger.Error("error: updating username in mock data")
	}
}
