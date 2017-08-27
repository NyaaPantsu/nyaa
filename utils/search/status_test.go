package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus_Parse(t *testing.T) {
	assert := assert.New(t)
	status := Status(0)
	status.Parse("")
	assert.Zero(status, "Should be equal to 0")
	status.Parse("lol10")
	assert.Zero(status, "Should be equal to 0")
	status.Parse("2")
	assert.Equal(Status(2), status)
	status.Parse("10")
	assert.Zero(status, "Should be equal 0")
}

func TestStatus_String(t *testing.T) {
	assert := assert.New(t)
	status := Status(0)
	status.Parse("")
	assert.Empty(status.String(), "Should be empty")
	status.Parse("lol10")
	assert.Empty(status.String(), "Should be empty")
	status.Parse("2")
	assert.Equal("2", status.String())
	status.Parse("10")
	assert.Empty(status.String(), "Should be empty")
}
func TestStatus_ToESQuery(t *testing.T) {
	assert := assert.New(t)
	status := Status(0)
	assert.Empty(status.ToESQuery(), "Should be empty")
	status = Status(3)
	assert.Equal("status:>3", status.ToESQuery(), "Should be equal")
	status = Status(1)
	assert.Equal("status:>1", status.ToESQuery(), "Should be equal")
	status = Status(2)
	assert.Equal("!status:2", status.ToESQuery(), "Should be equal")
}
func TestStatus_ToDBQuery(t *testing.T) {
	assert := assert.New(t)
	status := Status(0)
	assert.Empty(status.ToDBQuery(), "Should be empty")
	status = Status(3)
	assert.Equal("status >= 3", status.ToDBQuery(), "Should be equal")
	status = Status(1)
	assert.Equal("status >= 1", status.ToDBQuery(), "Should be equal")
	status = Status(2)
	assert.Equal("status <> 2", status.ToDBQuery(), "Should be equal")
}
