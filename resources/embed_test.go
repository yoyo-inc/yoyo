package resources_test

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yoyo-inc/yoyo/resources"
)

func TestEmbed(t *testing.T) {
	t.Run("report", func(t *testing.T) {
		fsys, err := fs.Sub(resources.InternalReportTplDir, "report")
		assert.Nil(t, err)
		_, err = fsys.Open("default/template.html")
		assert.Nil(t, err)
	})
}
