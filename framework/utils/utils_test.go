package utils_test

import (
	"testing"

	"github.com/robertvitoriano/go-encoder-microservice/framework/utils"
	"github.com/stretchr/testify/require"
)

func TestIsJson(t *testing.T) {

	testJSON := `{"hello_world":"name"}`

	err := utils.IsJson(testJSON)

	require.Nil(t, err)
}
