package utils_test

import (
	"testing"

	"github.com/Codeflix-FullCycle/encoder/framework/utils"
	"github.com/stretchr/testify/require"
)

func TestIsJson(t *testing.T) {
	js := `{
		"id":"f8004c83-5617-4f0c-9826-dc20e0fea240",
		"file_path":"test.mp4",
		"status":"pending"
	}`

	err := utils.IsJson(js)

	require.Empty(t, err)
}
