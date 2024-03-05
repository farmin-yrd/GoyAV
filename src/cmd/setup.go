package main

import (
	"errors"
	"fmt"
	"goyav/internal/adapter/antivirus"
	"goyav/internal/adapter/storage/byterepo"
	"goyav/internal/adapter/storage/docrepo"
	"goyav/internal/core/port"
	"goyav/pkg/helper"
	"log/slog"
	"os"
	"strconv"
)

const (
	// Default upload size limit in bytes : 1 Mib
	DefaultMaxUploadSize uint64 = 1 << 20
	defaultUploadTimeout uint64 = 10
)

// setup initializes the GoyAV application with necessary configurations.
// It configures the host, port, max upload size, version, and information for the application,
// along with initializing byte, document repositories and antivirus analyzer
func setup(host *string, port *int64, maxUploadSize *uint64, uploadTimeout *uint64, ver *string, info *string, b *port.ByteRepository, d *port.DocumentRepository, a *port.AntivirusAnalyzer) error {
	var err error

	setLogger()

	// Configure host and port
	*host = helper.GetEnvWithDefault("GOYAV_HOST", "localhost")
	*port, err = strconv.ParseInt(helper.GetEnvWithDefault("GOYAV_PORT", "80"), 10, 64)
	if err != nil {
		return errors.New("GOYAV_PORT must be a valid port number")
	}
	slog.Info("server configuration", "host", *host, "port", *port)

	// Configure version
	*ver, err = helper.GetEnvWithError("GOYAV_VERSION")
	if err != nil {
		return errors.New("GOYAV_VERSION must be set")
	}
	slog.Info("application version set", "version", *ver)

	// Configure information
	*info = helper.GetEnvWithDefault("GOYAV_INFORMATION", "GoyAV")
	slog.Info("application information set", "information", *info)

	// Configure maximum upload size
	*maxUploadSize, err = strconv.ParseUint(helper.GetEnvWithDefault("GOYAV_MAX_UPLOAD_SIZE", ""), 10, 64)
	if err != nil || *maxUploadSize == 0 {
		*maxUploadSize = DefaultMaxUploadSize
		slog.Warn("setting maximum upload size set to default", "default (MiB)", *maxUploadSize)
	}
	slog.Info("maximum upload size set", "size (MiB)", *maxUploadSize)

	// Configure upload timeout
	*uploadTimeout, err = strconv.ParseUint(helper.GetEnvWithDefault("GOYAV_UPLOAD_TIMEOUT", ""), 10, 64)
	if err != nil || *uploadTimeout <= 0 {
		*uploadTimeout = defaultUploadTimeout
		slog.Warn("setting upload timeout to default", "default (seconds)", defaultUploadTimeout)
	}
	slog.Info("upload timeout set", "timeout (seconds)", uploadTimeout)

	// Initialize byte repository
	if err = setupMinioByteRepository(b); err != nil {
		return fmt.Errorf("error while creating byte repository: %w", err)
	}

	// Initialize document repository
	if err = setupPostgresDocumentRepository(d); err != nil {
		return fmt.Errorf("error while creating document repository: %w", err)
	}

	// Initialize antivirus analyzer
	if err = setupClamAVAnalyzer(a); err != nil {
		return fmt.Errorf("error while creating antivirus analyzer: %w", err)
	}

	return nil
}

// setupMinioByteRepository configures a Minio byte repository for storing binary data of files.
func setupMinioByteRepository(b *port.ByteRepository) error {
	var err error

	// Retrieve Minio host configuration
	minioHost := helper.GetEnvWithDefault("GOYAV_MINIO_HOST", "127.0.0.1")
	slog.Info("configuring minio", "host", minioHost)

	// Parse and validate Minio port
	minioPort, err := strconv.ParseUint(helper.GetEnvWithDefault("GOYAV_MINIO_PORT", "9000"), 10, 64)
	if err != nil {
		return errors.New("GOYAV_MINIO_PORT must be a valid port number")
	}
	slog.Info("configuring minio", "port", minioPort)

	// Retrieve Minio access key ID with error check
	minioAccessKeyID, err := helper.GetEnvWithError("GOYAV_MINIO_ACCESS_KEY")
	if err != nil {
		return err
	}
	slog.Info("configuring minio", "access key ID", minioAccessKeyID)

	// Retrieve Minio secret key with error check
	minioSecretKey, err := helper.GetEnvWithError("GOYAV_MINIO_SECRET_KEY")
	if err != nil {
		return err
	}
	slog.Debug("configuring minio", "minio secret key", minioSecretKey)

	// Retrieve Minio bucket name configuration
	minioBucketName := helper.GetEnvWithDefault("GOYAV_MINIO_BUCKET_NAME", "goyav")
	slog.Info("configuring minio", "minio bucket name", minioBucketName)

	// Parse and validate Minio SSL usage
	minioUseSSL, err := strconv.ParseBool(helper.GetEnvWithDefault("GOYAV_MINIO_USE_SSL", "false"))
	if err != nil {
		return errors.New("GOYAV_MINIO_USE_SSL must be true or false")
	}
	slog.Info("configuring minio", "use ssl", minioUseSSL)

	*b, err = byterepo.NewMinio(
		minioHost,
		minioPort,
		minioAccessKeyID,
		minioSecretKey,
		minioBucketName,
		minioUseSSL,
	)
	if err != nil {
		return err
	}

	slog.Info("minio repository setup complete")
	return nil
}

// setupPostgresDocumentRepository configures a Postgres document repository.
func setupPostgresDocumentRepository(d *port.DocumentRepository) error {
	var err error

	// Retrieve PostgreSQL host configuration
	pgHost := helper.GetEnvWithDefault("GOYAV_POSTGRES_HOST", "127.0.0.1")
	slog.Info("configuring postgres", "host", pgHost)

	// Parse and validate PostgreSQL port
	pgPort, err := strconv.ParseUint(helper.GetEnvWithDefault("GOYAV_POSTGRES_PORT", "5432"), 10, 64)
	if err != nil {
		return errors.New("GOYAV_POSTGRES_PORT must be a valid port number")
	}
	slog.Info("configuring postgres", "port", pgPort)

	// Retrieve PostgreSQL user
	pgUser, err := helper.GetEnvWithError("GOYAV_POSTGRES_USER")
	if err != nil {
		return err
	}
	slog.Info("configuring postgres", "user", pgUser)

	// Retrieve PostgreSQL user password
	pgUserPassword, err := helper.GetEnvWithError("GOYAV_POSTGRES_USER_PASSWORD")
	if err != nil {
		return err
	}
	slog.Debug("configuring postgres", "password", pgUserPassword)

	// Retrieve PostgreSQL database name
	pgDB, err := helper.GetEnvWithError("GOYAV_POSTGRES_DB")
	if err != nil {
		return err
	}
	slog.Info("configuring postgres", "database name", pgDB)

	// Retrieve PostgreSQL schema
	pgSchema, err := helper.GetEnvWithError("GOYAV_POSTGRES_SCHEMA")
	if err != nil {
		return err
	}
	slog.Info("configuring postgres", "postgres schema name", pgSchema)

	// Initialize the PostgreSQL document repository
	*d, err = docrepo.NewPotgresDocumentRepository(pgHost, pgPort, pgDB, pgSchema, pgUser, pgUserPassword, false)
	if err != nil {
		return err
	}

	slog.Info("postgres repository setup complete")
	return nil
}

// setupClamAVAnalyzer configures a ClamAV antivirus analyzer.
func setupClamAVAnalyzer(a *port.AntivirusAnalyzer) error {
	var err error

	// Retrieve ClamAV host configuration
	clamdHost := helper.GetEnvWithDefault("GOYAV_CLAMAV_HOST", "127.0.0.1")
	slog.Info("configuring clamav", "host", clamdHost)

	// Parse and validate ClamAV port
	clamdPort, err := strconv.ParseUint(helper.GetEnvWithDefault("GOYAV_CLAMAV_PORT", "3310"), 10, 64)
	if err != nil {
		return errors.New("GOYAV_CLAMAV_PORT must be a valid port number")
	}
	slog.Info("configuring clamav", "port", clamdPort)

	// Parse and validate ClamAV timeout
	clamdTimeout, err := strconv.ParseUint(helper.GetEnvWithDefault("GOYAV_CLAMAV_TIMEOUT", "30"), 10, 64)
	if err != nil {
		return errors.New("GOYAV_CLAMAV_TIMEOUT must be a strictly positive number")
	}
	slog.Info("configuring clamav", "timeout", clamdTimeout)

	// Initialize the ClamAV analyzer
	*a, err = antivirus.NewClamavAnalyser(clamdHost, clamdPort, clamdTimeout)
	if err != nil {
		return err
	}

	slog.Info("clamav analyzer setup complete")
	return nil
}

func setLogger() {
	var level slog.Level = slog.LevelInfo

	isDubugMode, _ := strconv.ParseBool(helper.GetEnvWithDefault("GOYAV_DEBUG_MODE", "false"))
	if isDubugMode {
		level = slog.LevelDebug
	}

	slog.SetDefault(
		slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: level,
			}),
		),
	)
}
