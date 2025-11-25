package pg

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

// GetDSN возвращает dsn шаблона бд для тестов
// К этой бд не нужно подключатся. Ее нужно копировать для
// создания новой бд, с которой будет работать тест
func GetDSN() string {
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		dsn = "postgresql://localhost/chat_test?sslmode=disable"
	}
	return dsn
}

// CreateDB Создаст новую бд для тестов
// И возвращает строку подклюления к ней
func CreateDB(t *testing.T, name string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	templateDsn, err := url.Parse(GetDSN())
	require.NoError(t, err)
	templateDB := strings.Trim(templateDsn.Path, "/")

	conn, err := pgx.Connect(ctx, generalDsn(*templateDsn).String())
	require.NoError(t, err)
	defer func() { _ = conn.Close(ctx) }()

	_, err = conn.Exec(ctx, dropConnections(templateDB))
	require.NoError(t, err)

	_, err = conn.Exec(ctx, dropDatabase(name))
	if err != nil {
		fmt.Printf("unable to drop database %s: %s", name, err)
	}
	require.NoError(t, err)

	_, err = conn.Exec(ctx, createDatabaseFromTemplate(templateDB, name))
	require.NoError(t, err)

	newDSN := *templateDsn
	newDSN.Path = name

	fmt.Printf("CREATE TEST DB: %s\n", newDSN.String())
	return newDSN.String()
}

// DbNameFor Создаст имя БД по названию теста
// Если имя окажется слишком длинным - вместо имени будет
// использован хеш от имени
func DbNameFor(t *testing.T) string {
	dbName := t.Name()
	dbName = strings.ToLower(dbName)
	dbName = strings.ReplaceAll(dbName, "/", "_")

	// Если имя теста оказалось слишком длинным
	// то будет использован хэш
	if len(dbName) > 63 {
		// nolint:gosec // Это только для тестов
		hash := sha1.Sum([]byte(dbName))
		dbName = hex.EncodeToString(hash[:])
	}

	return dbName
}

func generalDsn(dsn url.URL) *url.URL {
	// Эта БД есть всегда. Поэтому она выбрана для подключения
	// А вот к templateDB подключаться нельзя, т.к. она будет
	// использована для клонирвоания
	// Нельзя клонировать БД к которой кто-то подключился
	dsn.Path = "template1"
	return &dsn
}

func dropConnections(dbName string) string {
	return fmt.Sprintf(`
		SELECT pg_terminate_backend(pg_stat_activity.pid) 
		FROM pg_stat_activity 
		WHERE pg_stat_activity.datname = '%s' AND pid <> pg_backend_pid();`,
		dbName)
}

func dropDatabase(dbName string) string {
	return fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)
}

func createDatabaseFromTemplate(template, dbName string) string {
	return fmt.Sprintf(`CREATE DATABASE %s WITH TEMPLATE %q`, dbName, template)
}
