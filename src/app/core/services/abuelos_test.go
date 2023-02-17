package services_test

import (
	"testing"

	"github.com/blackshark537/abuelos/src/app/core/services"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	service := services.ProjectAbuelos("")
	assert.Greater(t, len(service), 0, "Expect to be greater than 0")
}
