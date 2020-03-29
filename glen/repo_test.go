package glen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRepo(t *testing.T) {
	t.Parallel()

	expected := &Repo{
		LocalPath:  ".",
		RemoteName: "origin",
	}
	repo := NewRepo()

	assert.Equal(t, expected, repo)
}
