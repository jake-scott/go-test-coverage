package coverage_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/vladopajic/go-test-coverage/v2/pkg/testcoverage/coverage"
	"github.com/vladopajic/go-test-coverage/v2/pkg/testcoverage/testdata"
)

const (
	profileOK  = "../testdata/" + testdata.ProfileOK
	profileNOK = "../testdata/" + testdata.ProfileNOK
	prefix     = "github.com/vladopajic/go-test-coverage/v2"
)

func Test_GenerateCoverageStats(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		return
	}

	// should not be able to read directory
	stats, err := GenerateCoverageStats(Config{Profile: t.TempDir()})
	assert.Error(t, err)
	assert.Empty(t, stats)

	// should get error parsing invalid profile file
	stats, err = GenerateCoverageStats(Config{Profile: profileNOK})
	assert.Error(t, err)
	assert.Empty(t, stats)

	// should be okay to read valid profile
	stats1, err := GenerateCoverageStats(Config{Profile: profileOK})
	assert.NoError(t, err)
	assert.NotEmpty(t, stats1)

	// should be okay to read valid profile
	stats2, err := GenerateCoverageStats(Config{
		Profile:      profileOK,
		ExcludePaths: []string{`cover\.go$`},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, stats2)
	// stats2 should have less total statements because cover.go should have been excluded
	assert.Greater(t, CalcTotalStats(stats1).Total, CalcTotalStats(stats2).Total)

	// should remove prefix from stats
	stats3, err := GenerateCoverageStats(Config{
		Profile:     profileOK,
		LocalPrefix: prefix,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, stats3)
	assert.Equal(t, CalcTotalStats(stats1), CalcTotalStats(stats3))
	assert.NotContains(t, stats3[0].Name, prefix)
}

func Test_findFile(t *testing.T) {
	t.Parallel()

	const filename = "pkg/testcoverage/coverage/cover.go"

	file, noPrefixName, err := FindFile(prefix+"/"+filename, "")
	assert.NoError(t, err)
	assert.Equal(t, filename, noPrefixName)
	assert.True(t, strings.HasSuffix(file, filename))

	file, noPrefixName, err = FindFile(prefix+"/"+filename, prefix)
	assert.NoError(t, err)
	assert.Equal(t, filename, noPrefixName)
	assert.True(t, strings.HasSuffix(file, filename))

	_, _, err = FindFile(prefix+"/main1.go", "")
	assert.Error(t, err)

	_, _, err = FindFile("", "")
	assert.Error(t, err)

	_, _, err = FindFile(prefix, "")
	assert.Error(t, err)
}
