package postgres_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	repository "github.com/LieAlbertTriAdrian/clean-arch-golang/internal/postgres"
)

func TestEncodeCursor(t *testing.T) {
	utcLoc, err := time.LoadLocation("UTC")
	require.NoError(t, err)
	date := time.Date(2019, 10, 23, 17, 48, 0, 0, utcLoc)
	res := repository.EncodeCursor(date.Unix())
	require.Equal(t, "MTU3MTg1Mjg4MA==", res)
}

func TestDencodeCursor(t *testing.T) {
	expected := int64(1571852880)
	cursor := "MTU3MTg1Mjg4MA=="
	res, err := repository.DecodeCursor(cursor)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestDencodeCursorError(t *testing.T) {
	expected := int64(1571852880)
	cursor := "MTU3MTg1Mjg4MA"
	res, err := repository.DecodeCursor(cursor)
	require.Error(t, err)
	require.NotEqual(t, expected, res)
	require.Zero(t, res)
}
