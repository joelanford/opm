package declcfg

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joelanford/declcfg/internal/property"
)

func TestReadJSON(t *testing.T) {
	type spec struct {
		name              string
		file              string
		assertion         require.ErrorAssertionFunc
		expectNumPackages int
		expectNumBundles  int
		expectNumOthers   int
	}
	specs := []spec{
		{
			name:              "Ignored/NotJSON",
			file:              "testdata/invalid/not-json.txt",
			assertion:         require.NoError,
			expectNumPackages: 0,
			expectNumBundles:  0,
			expectNumOthers:   0,
		},
		{
			name:              "Ignored/NotJSONObject",
			file:              "testdata/invalid/not-json-object.json",
			assertion:         require.NoError,
			expectNumPackages: 0,
			expectNumBundles:  0,
			expectNumOthers:   0,
		},
		{
			name:      "Error/InvalidPackageJSON",
			file:      "testdata/invalid/invalid-package-json.json",
			assertion: require.Error,
		},
		{
			name:      "Error/InvalidBundleJSON",
			file:      "testdata/invalid/invalid-bundle-json.json",
			assertion: require.Error,
		},
		{
			name:              "Success/UnrecognizedSchema",
			file:              "testdata/valid/unrecognized-schema.json",
			assertion:         require.NoError,
			expectNumPackages: 1,
			expectNumBundles:  1,
			expectNumOthers:   1,
		},
		{
			name:              "Success/ValidFile",
			file:              "testdata/valid/etcd.json",
			assertion:         require.NoError,
			expectNumPackages: 1,
			expectNumBundles:  6,
			expectNumOthers:   0,
		},
	}

	for _, s := range specs {
		t.Run(s.name, func(t *testing.T) {
			f, err := os.Open(s.file)
			require.NoError(t, err)

			cfg, err := readJSON(f)
			s.assertion(t, err)
			if err == nil {
				require.NotNil(t, cfg)
				assert.Equal(t, len(cfg.Packages), s.expectNumPackages, "unexpected package count")
				assert.Equal(t, len(cfg.Bundles), s.expectNumBundles, "unexpected bundle count")
				assert.Equal(t, len(cfg.Others), s.expectNumOthers, "unexpected others count")
			}
		})
	}
}

func TestLoadDir(t *testing.T) {
	type spec struct {
		name      string
		dir       string
		assertion require.ErrorAssertionFunc
		expected  *DeclarativeConfig
	}
	specs := []spec{
		{
			name:      "Error/NonExistentDir",
			dir:       "testdata/nonexistent",
			assertion: require.Error,
		},
		{
			name:      "Error/Invalid",
			dir:       "testdata/invalid",
			assertion: require.Error,
		},
		{
			name:      "Success/ValidDir",
			dir:       "testdata/valid",
			assertion: require.NoError,
			expected: &DeclarativeConfig{
				Packages: []Package{
					{Schema: "olm.package", Name: "cockroachdb", DefaultChannel: "stable-5.x", Icon: &Icon{Data: []uint8{0x3c, 0x73, 0x76, 0x67, 0x20, 0x78, 0x6d, 0x6c, 0x6e, 0x73, 0x3d, 0x22, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x77, 0x77, 0x77, 0x2e, 0x77, 0x33, 0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x32, 0x30, 0x30, 0x30, 0x2f, 0x73, 0x76, 0x67, 0x22, 0x20, 0x76, 0x69, 0x65, 0x77, 0x42, 0x6f, 0x78, 0x3d, 0x22, 0x30, 0x20, 0x30, 0x20, 0x33, 0x31, 0x2e, 0x38, 0x32, 0x20, 0x33, 0x32, 0x22, 0x20, 0x77, 0x69, 0x64, 0x74, 0x68, 0x3d, 0x22, 0x32, 0x34, 0x38, 0x36, 0x22, 0x20, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x3d, 0x22, 0x32, 0x35, 0x30, 0x30, 0x22, 0x3e, 0x3c, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x3e, 0x43, 0x4c, 0x3c, 0x2f, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x3e, 0x3c, 0x70, 0x61, 0x74, 0x68, 0x20, 0x64, 0x3d, 0x22, 0x4d, 0x31, 0x39, 0x2e, 0x34, 0x32, 0x20, 0x39, 0x2e, 0x31, 0x37, 0x61, 0x31, 0x35, 0x2e, 0x33, 0x39, 0x20, 0x31, 0x35, 0x2e, 0x33, 0x39, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x2d, 0x33, 0x2e, 0x35, 0x31, 0x2e, 0x34, 0x20, 0x31, 0x35, 0x2e, 0x34, 0x36, 0x20, 0x31, 0x35, 0x2e, 0x34, 0x36, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x2d, 0x33, 0x2e, 0x35, 0x31, 0x2d, 0x2e, 0x34, 0x20, 0x31, 0x35, 0x2e, 0x36, 0x33, 0x20, 0x31, 0x35, 0x2e, 0x36, 0x33, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x20, 0x33, 0x2e, 0x35, 0x31, 0x2d, 0x33, 0x2e, 0x39, 0x31, 0x20, 0x31, 0x35, 0x2e, 0x37, 0x31, 0x20, 0x31, 0x35, 0x2e, 0x37, 0x31, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x20, 0x33, 0x2e, 0x35, 0x31, 0x20, 0x33, 0x2e, 0x39, 0x31, 0x7a, 0x4d, 0x33, 0x30, 0x20, 0x2e, 0x35, 0x37, 0x41, 0x31, 0x37, 0x2e, 0x32, 0x32, 0x20, 0x31, 0x37, 0x2e, 0x32, 0x32, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x20, 0x32, 0x35, 0x2e, 0x35, 0x39, 0x20, 0x30, 0x61, 0x31, 0x37, 0x2e, 0x34, 0x20, 0x31, 0x37, 0x2e, 0x34, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x2d, 0x39, 0x2e, 0x36, 0x38, 0x20, 0x32, 0x2e, 0x39, 0x33, 0x41, 0x31, 0x37, 0x2e, 0x33, 0x38, 0x20, 0x31, 0x37, 0x2e, 0x33, 0x38, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x20, 0x36, 0x2e, 0x32, 0x33, 0x20, 0x30, 0x61, 0x31, 0x37, 0x2e, 0x32, 0x32, 0x20, 0x31, 0x37, 0x2e, 0x32, 0x32, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x2d, 0x34, 0x2e, 0x34, 0x34, 0x2e, 0x35, 0x37, 0x41, 0x31, 0x36, 0x2e, 0x32, 0x32, 0x20, 0x31, 0x36, 0x2e, 0x32, 0x32, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x2e, 0x31, 0x33, 0x61, 0x2e, 0x30, 0x37, 0x2e, 0x30, 0x37, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x20, 0x2e, 0x30, 0x39, 0x20, 0x31, 0x37, 0x2e, 0x33, 0x32, 0x20, 0x31, 0x37, 0x2e, 0x33, 0x32, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x20, 0x2e, 0x38, 0x33, 0x20, 0x31, 0x2e, 0x35, 0x37, 0x2e, 0x30, 0x37, 0x2e, 0x30, 0x37, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x20, 0x2e, 0x30, 0x38, 0x20, 0x30, 0x20, 0x31, 0x36, 0x2e, 0x33, 0x39, 0x20, 0x31, 0x36, 0x2e, 0x33, 0x39, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x20, 0x31, 0x2e, 0x38, 0x31, 0x2d, 0x2e, 0x35, 0x34, 0x20, 0x31, 0x35, 0x2e, 0x36, 0x35, 0x20, 0x31, 0x35, 0x2e, 0x36, 0x35, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x20, 0x31, 0x31, 0x2e, 0x35, 0x39, 0x20, 0x31, 0x2e, 0x38, 0x38, 0x20, 0x31, 0x37, 0x2e, 0x35, 0x32, 0x20, 0x31, 0x37, 0x2e, 0x35, 0x32, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x2d, 0x33, 0x2e, 0x37, 0x38, 0x20, 0x34, 0x2e, 0x34, 0x38, 0x63, 0x2d, 0x2e, 0x32, 0x2e, 0x33, 0x32, 0x2d, 0x2e, 0x33, 0x37, 0x2e, 0x36, 0x35, 0x2d, 0x2e, 0x35, 0x35, 0x20, 0x31, 0x73, 0x2d, 0x2e, 0x32, 0x32, 0x2e, 0x34, 0x35, 0x2d, 0x2e, 0x33, 0x33, 0x2e, 0x36, 0x39, 0x2d, 0x2e, 0x33, 0x31, 0x2e, 0x37, 0x32, 0x2d, 0x2e, 0x34, 0x34, 0x20, 0x31, 0x2e, 0x30, 0x38, 0x61, 0x31, 0x37, 0x2e, 0x34, 0x36, 0x20, 0x31, 0x37, 0x2e, 0x34, 0x36, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x20, 0x34, 0x2e, 0x32, 0x39, 0x20, 0x31, 0x38, 0x2e, 0x37, 0x63, 0x2e, 0x32, 0x36, 0x2e, 0x32, 0x35, 0x2e, 0x35, 0x33, 0x2e, 0x34, 0x39, 0x2e, 0x38, 0x31, 0x2e, 0x37, 0x33, 0x73, 0x2e, 0x34, 0x34, 0x2e, 0x33, 0x37, 0x2e, 0x36, 0x37, 0x2e, 0x35, 0x34, 0x2e, 0x35, 0x39, 0x2e, 0x34, 0x34, 0x2e, 0x38, 0x39, 0x2e, 0x36, 0x34, 0x61, 0x2e, 0x30, 0x37, 0x2e, 0x30, 0x37, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x20, 0x2e, 0x30, 0x38, 0x20, 0x30, 0x63, 0x2e, 0x33, 0x2d, 0x2e, 0x32, 0x31, 0x2e, 0x36, 0x2d, 0x2e, 0x34, 0x32, 0x2e, 0x38, 0x39, 0x2d, 0x2e, 0x36, 0x34, 0x73, 0x2e, 0x34, 0x35, 0x2d, 0x2e, 0x33, 0x35, 0x2e, 0x36, 0x37, 0x2d, 0x2e, 0x35, 0x34, 0x2e, 0x35, 0x35, 0x2d, 0x2e, 0x34, 0x38, 0x2e, 0x38, 0x31, 0x2d, 0x2e, 0x37, 0x33, 0x61, 0x31, 0x37, 0x2e, 0x34, 0x35, 0x20, 0x31, 0x37, 0x2e, 0x34, 0x35, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x20, 0x35, 0x2e, 0x33, 0x38, 0x2d, 0x31, 0x32, 0x2e, 0x36, 0x31, 0x20, 0x31, 0x37, 0x2e, 0x33, 0x39, 0x20, 0x31, 0x37, 0x2e, 0x33, 0x39, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x2d, 0x31, 0x2e, 0x30, 0x39, 0x2d, 0x36, 0x2e, 0x30, 0x39, 0x63, 0x2d, 0x2e, 0x31, 0x34, 0x2d, 0x2e, 0x33, 0x37, 0x2d, 0x2e, 0x32, 0x39, 0x2d, 0x2e, 0x37, 0x33, 0x2d, 0x2e, 0x34, 0x35, 0x2d, 0x31, 0x2e, 0x30, 0x39, 0x73, 0x2d, 0x2e, 0x32, 0x32, 0x2d, 0x2e, 0x34, 0x37, 0x2d, 0x2e, 0x33, 0x33, 0x2d, 0x2e, 0x36, 0x39, 0x2d, 0x2e, 0x33, 0x35, 0x2d, 0x2e, 0x36, 0x36, 0x2d, 0x2e, 0x35, 0x35, 0x2d, 0x31, 0x61, 0x31, 0x37, 0x2e, 0x36, 0x31, 0x20, 0x31, 0x37, 0x2e, 0x36, 0x31, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x2d, 0x33, 0x2e, 0x37, 0x38, 0x2d, 0x34, 0x2e, 0x34, 0x38, 0x20, 0x31, 0x35, 0x2e, 0x36, 0x35, 0x20, 0x31, 0x35, 0x2e, 0x36, 0x35, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x20, 0x31, 0x31, 0x2e, 0x36, 0x2d, 0x31, 0x2e, 0x38, 0x34, 0x20, 0x31, 0x36, 0x2e, 0x31, 0x33, 0x20, 0x31, 0x36, 0x2e, 0x31, 0x33, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x20, 0x31, 0x2e, 0x38, 0x31, 0x2e, 0x35, 0x34, 0x2e, 0x30, 0x37, 0x2e, 0x30, 0x37, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x20, 0x2e, 0x30, 0x38, 0x20, 0x30, 0x71, 0x2e, 0x34, 0x34, 0x2d, 0x2e, 0x37, 0x36, 0x2e, 0x38, 0x32, 0x2d, 0x31, 0x2e, 0x35, 0x36, 0x61, 0x2e, 0x30, 0x37, 0x2e, 0x30, 0x37, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x2d, 0x2e, 0x30, 0x39, 0x41, 0x31, 0x36, 0x2e, 0x38, 0x39, 0x20, 0x31, 0x36, 0x2e, 0x38, 0x39, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x20, 0x33, 0x30, 0x20, 0x2e, 0x35, 0x37, 0x7a, 0x22, 0x20, 0x66, 0x69, 0x6c, 0x6c, 0x3d, 0x22, 0x23, 0x31, 0x35, 0x31, 0x66, 0x33, 0x34, 0x22, 0x2f, 0x3e, 0x3c, 0x70, 0x61, 0x74, 0x68, 0x20, 0x64, 0x3d, 0x22, 0x4d, 0x32, 0x31, 0x2e, 0x38, 0x32, 0x20, 0x31, 0x37, 0x2e, 0x34, 0x37, 0x61, 0x31, 0x35, 0x2e, 0x35, 0x31, 0x20, 0x31, 0x35, 0x2e, 0x35, 0x31, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x2d, 0x34, 0x2e, 0x32, 0x35, 0x20, 0x31, 0x30, 0x2e, 0x36, 0x39, 0x20, 0x31, 0x35, 0x2e, 0x36, 0x36, 0x20, 0x31, 0x35, 0x2e, 0x36, 0x36, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x2d, 0x2e, 0x37, 0x32, 0x2d, 0x34, 0x2e, 0x36, 0x38, 0x20, 0x31, 0x35, 0x2e, 0x35, 0x20, 0x31, 0x35, 0x2e, 0x35, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x20, 0x34, 0x2e, 0x32, 0x35, 0x2d, 0x31, 0x30, 0x2e, 0x36, 0x39, 0x20, 0x31, 0x35, 0x2e, 0x36, 0x32, 0x20, 0x31, 0x35, 0x2e, 0x36, 0x32, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x20, 0x2e, 0x37, 0x32, 0x20, 0x34, 0x2e, 0x36, 0x38, 0x22, 0x20, 0x66, 0x69, 0x6c, 0x6c, 0x3d, 0x22, 0x23, 0x33, 0x34, 0x38, 0x35, 0x34, 0x30, 0x22, 0x2f, 0x3e, 0x3c, 0x70, 0x61, 0x74, 0x68, 0x20, 0x64, 0x3d, 0x22, 0x4d, 0x31, 0x35, 0x20, 0x32, 0x33, 0x2e, 0x34, 0x38, 0x61, 0x31, 0x35, 0x2e, 0x35, 0x35, 0x20, 0x31, 0x35, 0x2e, 0x35, 0x35, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x2d, 0x2e, 0x37, 0x32, 0x20, 0x34, 0x2e, 0x36, 0x38, 0x20, 0x31, 0x35, 0x2e, 0x35, 0x34, 0x20, 0x31, 0x35, 0x2e, 0x35, 0x34, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x2d, 0x33, 0x2e, 0x35, 0x33, 0x2d, 0x31, 0x35, 0x2e, 0x33, 0x37, 0x41, 0x31, 0x35, 0x2e, 0x35, 0x20, 0x31, 0x35, 0x2e, 0x35, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x20, 0x31, 0x35, 0x20, 0x32, 0x33, 0x2e, 0x34, 0x38, 0x22, 0x20, 0x66, 0x69, 0x6c, 0x6c, 0x3d, 0x22, 0x23, 0x37, 0x64, 0x62, 0x63, 0x34, 0x32, 0x22, 0x2f, 0x3e, 0x3c, 0x2f, 0x73, 0x76, 0x67, 0x3e}, MediaType: "image/svg+xml"}, Description: ""},
					{Schema: "olm.package", Name: "etcd", DefaultChannel: "singlenamespace-alpha", Icon: &Icon{Data: []uint8{0x3c, 0x73, 0x76, 0x67, 0x20, 0x77, 0x69, 0x64, 0x74, 0x68, 0x3d, 0x22, 0x32, 0x35, 0x30, 0x30, 0x22, 0x20, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x3d, 0x22, 0x32, 0x34, 0x32, 0x32, 0x22, 0x20, 0x76, 0x69, 0x65, 0x77, 0x42, 0x6f, 0x78, 0x3d, 0x22, 0x30, 0x20, 0x30, 0x20, 0x32, 0x35, 0x36, 0x20, 0x32, 0x34, 0x38, 0x22, 0x20, 0x78, 0x6d, 0x6c, 0x6e, 0x73, 0x3d, 0x22, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x77, 0x77, 0x77, 0x2e, 0x77, 0x33, 0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x32, 0x30, 0x30, 0x30, 0x2f, 0x73, 0x76, 0x67, 0x22, 0x20, 0x70, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x41, 0x73, 0x70, 0x65, 0x63, 0x74, 0x52, 0x61, 0x74, 0x69, 0x6f, 0x3d, 0x22, 0x78, 0x4d, 0x69, 0x64, 0x59, 0x4d, 0x69, 0x64, 0x22, 0x3e, 0x3c, 0x70, 0x61, 0x74, 0x68, 0x20, 0x64, 0x3d, 0x22, 0x4d, 0x32, 0x35, 0x32, 0x2e, 0x33, 0x38, 0x36, 0x20, 0x31, 0x32, 0x38, 0x2e, 0x30, 0x36, 0x34, 0x63, 0x2d, 0x31, 0x2e, 0x32, 0x30, 0x32, 0x2e, 0x31, 0x2d, 0x32, 0x2e, 0x34, 0x31, 0x2e, 0x31, 0x34, 0x37, 0x2d, 0x33, 0x2e, 0x36, 0x39, 0x33, 0x2e, 0x31, 0x34, 0x37, 0x2d, 0x37, 0x2e, 0x34, 0x34, 0x36, 0x20, 0x30, 0x2d, 0x31, 0x34, 0x2e, 0x36, 0x37, 0x2d, 0x31, 0x2e, 0x37, 0x34, 0x36, 0x2d, 0x32, 0x31, 0x2e, 0x31, 0x38, 0x37, 0x2d, 0x34, 0x2e, 0x39, 0x34, 0x34, 0x20, 0x32, 0x2e, 0x31, 0x37, 0x2d, 0x31, 0x32, 0x2e, 0x34, 0x34, 0x37, 0x20, 0x33, 0x2e, 0x30, 0x39, 0x32, 0x2d, 0x32, 0x34, 0x2e, 0x39, 0x38, 0x37, 0x20, 0x32, 0x2e, 0x38, 0x35, 0x2d, 0x33, 0x37, 0x2e, 0x34, 0x38, 0x31, 0x2d, 0x37, 0x2e, 0x30, 0x36, 0x35, 0x2d, 0x31, 0x30, 0x2e, 0x32, 0x32, 0x2d, 0x31, 0x35, 0x2e, 0x31, 0x34, 0x2d, 0x31, 0x39, 0x2e, 0x38, 0x36, 0x33, 0x2d, 0x32, 0x34, 0x2e, 0x32, 0x35, 0x36, 0x2d, 0x32, 0x38, 0x2e, 0x37, 0x34, 0x37, 0x20, 0x33, 0x2e, 0x39, 0x35, 0x35, 0x2d, 0x37, 0x2e, 0x34, 0x31, 0x35, 0x20, 0x39, 0x2e, 0x38, 0x30, 0x31, 0x2d, 0x31, 0x33, 0x2e, 0x37, 0x39, 0x35, 0x20, 0x31, 0x37, 0x2e, 0x31, 0x2d, 0x31, 0x38, 0x2e, 0x33, 0x31, 0x39, 0x6c, 0x33, 0x2e, 0x31, 0x33, 0x33, 0x2d, 0x31, 0x2e, 0x39, 0x33, 0x37, 0x2d, 0x32, 0x2e, 0x34, 0x34, 0x32, 0x2d, 0x32, 0x2e, 0x37, 0x35, 0x34, 0x63, 0x2d, 0x31, 0x32, 0x2e, 0x35, 0x38, 0x31, 0x2d, 0x31, 0x34, 0x2e, 0x31, 0x36, 0x37, 0x2d, 0x32, 0x37, 0x2e, 0x35, 0x39, 0x36, 0x2d, 0x32, 0x35, 0x2e, 0x31, 0x32, 0x2d, 0x34, 0x34, 0x2e, 0x36, 0x32, 0x2d, 0x33, 0x32, 0x2e, 0x35, 0x35, 0x32, 0x4c, 0x31, 0x37, 0x35, 0x2e, 0x38, 0x37, 0x36, 0x20, 0x30, 0x6c, 0x2d, 0x2e, 0x38, 0x36, 0x32, 0x20, 0x33, 0x2e, 0x35, 0x38, 0x38, 0x63, 0x2d, 0x32, 0x2e, 0x30, 0x33, 0x20, 0x38, 0x2e, 0x33, 0x36, 0x33, 0x2d, 0x36, 0x2e, 0x32, 0x37, 0x34, 0x20, 0x31, 0x35, 0x2e, 0x39, 0x30, 0x38, 0x2d, 0x31, 0x32, 0x2e, 0x31, 0x20, 0x32, 0x31, 0x2e, 0x39, 0x36, 0x32, 0x61, 0x31, 0x39, 0x33, 0x2e, 0x38, 0x34, 0x32, 0x20, 0x31, 0x39, 0x33, 0x2e, 0x38, 0x34, 0x32, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x2d, 0x33, 0x34, 0x2e, 0x39, 0x35, 0x36, 0x2d, 0x31, 0x34, 0x2e, 0x34, 0x30, 0x35, 0x41, 0x31, 0x39, 0x34, 0x2e, 0x30, 0x31, 0x32, 0x20, 0x31, 0x39, 0x34, 0x2e, 0x30, 0x31, 0x32, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x20, 0x39, 0x33, 0x2e, 0x30, 0x35, 0x36, 0x20, 0x32, 0x35, 0x2e, 0x35, 0x32, 0x43, 0x38, 0x37, 0x2e, 0x32, 0x35, 0x34, 0x20, 0x31, 0x39, 0x2e, 0x34, 0x37, 0x33, 0x20, 0x38, 0x33, 0x2e, 0x30, 0x32, 0x20, 0x31, 0x31, 0x2e, 0x39, 0x34, 0x37, 0x20, 0x38, 0x30, 0x2e, 0x39, 0x39, 0x39, 0x20, 0x33, 0x2e, 0x36, 0x30, 0x38, 0x4c, 0x38, 0x30, 0x2e, 0x31, 0x33, 0x2e, 0x30, 0x32, 0x6c, 0x2d, 0x33, 0x2e, 0x33, 0x38, 0x32, 0x20, 0x31, 0x2e, 0x34, 0x37, 0x43, 0x35, 0x39, 0x2e, 0x39, 0x33, 0x39, 0x20, 0x38, 0x2e, 0x38, 0x31, 0x35, 0x20, 0x34, 0x34, 0x2e, 0x35, 0x31, 0x20, 0x32, 0x30, 0x2e, 0x30, 0x36, 0x35, 0x20, 0x33, 0x32, 0x2e, 0x31, 0x33, 0x35, 0x20, 0x33, 0x34, 0x2e, 0x30, 0x32, 0x6c, 0x2d, 0x32, 0x2e, 0x34, 0x34, 0x39, 0x20, 0x32, 0x2e, 0x37, 0x36, 0x20, 0x33, 0x2e, 0x31, 0x33, 0x20, 0x31, 0x2e, 0x39, 0x33, 0x37, 0x63, 0x37, 0x2e, 0x32, 0x37, 0x36, 0x20, 0x34, 0x2e, 0x35, 0x30, 0x36, 0x20, 0x31, 0x33, 0x2e, 0x31, 0x30, 0x36, 0x20, 0x31, 0x30, 0x2e, 0x38, 0x34, 0x39, 0x20, 0x31, 0x37, 0x2e, 0x30, 0x35, 0x34, 0x20, 0x31, 0x38, 0x2e, 0x32, 0x32, 0x33, 0x2d, 0x39, 0x2e, 0x30, 0x38, 0x38, 0x20, 0x38, 0x2e, 0x38, 0x35, 0x2d, 0x31, 0x37, 0x2e, 0x31, 0x35, 0x34, 0x20, 0x31, 0x38, 0x2e, 0x34, 0x36, 0x32, 0x2d, 0x32, 0x34, 0x2e, 0x32, 0x31, 0x34, 0x20, 0x32, 0x38, 0x2e, 0x36, 0x33, 0x35, 0x2d, 0x2e, 0x32, 0x37, 0x35, 0x20, 0x31, 0x32, 0x2e, 0x34, 0x38, 0x39, 0x2e, 0x36, 0x20, 0x32, 0x35, 0x2e, 0x31, 0x32, 0x20, 0x32, 0x2e, 0x37, 0x38, 0x20, 0x33, 0x37, 0x2e, 0x37, 0x34, 0x2d, 0x36, 0x2e, 0x34, 0x38, 0x34, 0x20, 0x33, 0x2e, 0x31, 0x36, 0x37, 0x2d, 0x31, 0x33, 0x2e, 0x36, 0x36, 0x38, 0x20, 0x34, 0x2e, 0x38, 0x39, 0x34, 0x2d, 0x32, 0x31, 0x2e, 0x30, 0x36, 0x35, 0x20, 0x34, 0x2e, 0x38, 0x39, 0x34, 0x2d, 0x31, 0x2e, 0x32, 0x39, 0x38, 0x20, 0x30, 0x2d, 0x32, 0x2e, 0x35, 0x31, 0x33, 0x2d, 0x2e, 0x30, 0x34, 0x37, 0x2d, 0x33, 0x2e, 0x36, 0x39, 0x33, 0x2d, 0x2e, 0x31, 0x34, 0x35, 0x4c, 0x30, 0x20, 0x31, 0x32, 0x37, 0x2e, 0x37, 0x38, 0x35, 0x6c, 0x2e, 0x33, 0x34, 0x35, 0x20, 0x33, 0x2e, 0x36, 0x37, 0x31, 0x63, 0x31, 0x2e, 0x38, 0x30, 0x32, 0x20, 0x31, 0x38, 0x2e, 0x35, 0x37, 0x38, 0x20, 0x37, 0x2e, 0x35, 0x37, 0x20, 0x33, 0x36, 0x2e, 0x32, 0x34, 0x37, 0x20, 0x31, 0x37, 0x2e, 0x31, 0x35, 0x34, 0x20, 0x35, 0x32, 0x2e, 0x35, 0x32, 0x33, 0x6c, 0x31, 0x2e, 0x38, 0x37, 0x20, 0x33, 0x2e, 0x31, 0x37, 0x36, 0x20, 0x32, 0x2e, 0x38, 0x31, 0x2d, 0x32, 0x2e, 0x33, 0x38, 0x34, 0x61, 0x34, 0x38, 0x2e, 0x30, 0x34, 0x20, 0x34, 0x38, 0x2e, 0x30, 0x34, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x20, 0x32, 0x32, 0x2e, 0x37, 0x33, 0x37, 0x2d, 0x31, 0x30, 0x2e, 0x36, 0x35, 0x20, 0x31, 0x39, 0x34, 0x2e, 0x38, 0x36, 0x20, 0x31, 0x39, 0x34, 0x2e, 0x38, 0x36, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x39, 0x2e, 0x34, 0x36, 0x20, 0x33, 0x31, 0x2e, 0x36, 0x39, 0x36, 0x63, 0x31, 0x31, 0x2e, 0x38, 0x32, 0x38, 0x20, 0x34, 0x2e, 0x31, 0x33, 0x37, 0x20, 0x32, 0x34, 0x2e, 0x31, 0x35, 0x31, 0x20, 0x37, 0x2e, 0x32, 0x32, 0x35, 0x20, 0x33, 0x36, 0x2e, 0x38, 0x37, 0x38, 0x20, 0x39, 0x2e, 0x30, 0x36, 0x33, 0x20, 0x31, 0x2e, 0x32, 0x32, 0x20, 0x38, 0x2e, 0x34, 0x31, 0x37, 0x2e, 0x32, 0x34, 0x38, 0x20, 0x31, 0x37, 0x2e, 0x31, 0x32, 0x32, 0x2d, 0x33, 0x2e, 0x30, 0x37, 0x32, 0x20, 0x32, 0x35, 0x2e, 0x31, 0x37, 0x31, 0x6c, 0x2d, 0x31, 0x2e, 0x34, 0x20, 0x33, 0x2e, 0x34, 0x31, 0x31, 0x20, 0x33, 0x2e, 0x36, 0x2e, 0x37, 0x39, 0x33, 0x63, 0x39, 0x2e, 0x32, 0x32, 0x20, 0x32, 0x2e, 0x30, 0x32, 0x37, 0x20, 0x31, 0x38, 0x2e, 0x35, 0x32, 0x33, 0x20, 0x33, 0x2e, 0x30, 0x36, 0x20, 0x32, 0x37, 0x2e, 0x36, 0x33, 0x31, 0x20, 0x33, 0x2e, 0x30, 0x36, 0x6c, 0x32, 0x37, 0x2e, 0x36, 0x32, 0x33, 0x2d, 0x33, 0x2e, 0x30, 0x36, 0x20, 0x33, 0x2e, 0x36, 0x30, 0x34, 0x2d, 0x2e, 0x37, 0x39, 0x33, 0x2d, 0x31, 0x2e, 0x34, 0x30, 0x33, 0x2d, 0x33, 0x2e, 0x34, 0x31, 0x37, 0x63, 0x2d, 0x33, 0x2e, 0x33, 0x31, 0x32, 0x2d, 0x38, 0x2e, 0x30, 0x35, 0x2d, 0x34, 0x2e, 0x32, 0x38, 0x34, 0x2d, 0x31, 0x36, 0x2e, 0x37, 0x36, 0x35, 0x2d, 0x33, 0x2e, 0x30, 0x36, 0x33, 0x2d, 0x32, 0x35, 0x2e, 0x31, 0x38, 0x33, 0x20, 0x31, 0x32, 0x2e, 0x36, 0x37, 0x36, 0x2d, 0x31, 0x2e, 0x38, 0x34, 0x20, 0x32, 0x34, 0x2e, 0x39, 0x35, 0x34, 0x2d, 0x34, 0x2e, 0x39, 0x32, 0x20, 0x33, 0x36, 0x2e, 0x37, 0x33, 0x38, 0x2d, 0x39, 0x2e, 0x30, 0x34, 0x35, 0x61, 0x31, 0x39, 0x35, 0x2e, 0x31, 0x30, 0x38, 0x20, 0x31, 0x39, 0x35, 0x2e, 0x31, 0x30, 0x38, 0x20, 0x30, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x39, 0x2e, 0x34, 0x38, 0x32, 0x2d, 0x33, 0x31, 0x2e, 0x37, 0x32, 0x36, 0x20, 0x34, 0x38, 0x2e, 0x32, 0x35, 0x34, 0x20, 0x34, 0x38, 0x2e, 0x32, 0x35, 0x34, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x20, 0x32, 0x32, 0x2e, 0x38, 0x34, 0x38, 0x20, 0x31, 0x30, 0x2e, 0x36, 0x36, 0x6c, 0x32, 0x2e, 0x38, 0x30, 0x39, 0x20, 0x32, 0x2e, 0x33, 0x38, 0x20, 0x31, 0x2e, 0x38, 0x36, 0x32, 0x2d, 0x33, 0x2e, 0x31, 0x36, 0x38, 0x63, 0x39, 0x2e, 0x36, 0x2d, 0x31, 0x36, 0x2e, 0x32, 0x39, 0x37, 0x20, 0x31, 0x35, 0x2e, 0x33, 0x36, 0x38, 0x2d, 0x33, 0x33, 0x2e, 0x39, 0x36, 0x35, 0x20, 0x31, 0x37, 0x2e, 0x31, 0x34, 0x32, 0x2d, 0x35, 0x32, 0x2e, 0x35, 0x31, 0x33, 0x6c, 0x2e, 0x33, 0x34, 0x35, 0x2d, 0x33, 0x2e, 0x36, 0x36, 0x35, 0x2d, 0x33, 0x2e, 0x36, 0x31, 0x34, 0x2e, 0x32, 0x37, 0x39, 0x7a, 0x4d, 0x31, 0x36, 0x37, 0x2e, 0x34, 0x39, 0x20, 0x31, 0x37, 0x32, 0x2e, 0x39, 0x36, 0x63, 0x2d, 0x31, 0x33, 0x2e, 0x30, 0x36, 0x38, 0x20, 0x33, 0x2e, 0x35, 0x35, 0x34, 0x2d, 0x32, 0x36, 0x2e, 0x33, 0x34, 0x20, 0x35, 0x2e, 0x33, 0x34, 0x38, 0x2d, 0x33, 0x39, 0x2e, 0x35, 0x33, 0x32, 0x20, 0x35, 0x2e, 0x33, 0x34, 0x38, 0x2d, 0x31, 0x33, 0x2e, 0x32, 0x32, 0x38, 0x20, 0x30, 0x2d, 0x32, 0x36, 0x2e, 0x34, 0x38, 0x33, 0x2d, 0x31, 0x2e, 0x37, 0x39, 0x33, 0x2d, 0x33, 0x39, 0x2e, 0x35, 0x36, 0x33, 0x2d, 0x35, 0x2e, 0x33, 0x34, 0x38, 0x61, 0x31, 0x35, 0x33, 0x2e, 0x32, 0x35, 0x35, 0x20, 0x31, 0x35, 0x33, 0x2e, 0x32, 0x35, 0x35, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x2d, 0x31, 0x36, 0x2e, 0x39, 0x33, 0x32, 0x2d, 0x33, 0x35, 0x2e, 0x36, 0x37, 0x63, 0x2d, 0x34, 0x2e, 0x30, 0x36, 0x36, 0x2d, 0x31, 0x32, 0x2e, 0x35, 0x31, 0x37, 0x2d, 0x36, 0x2e, 0x34, 0x34, 0x35, 0x2d, 0x32, 0x35, 0x2e, 0x36, 0x33, 0x2d, 0x37, 0x2e, 0x31, 0x33, 0x35, 0x2d, 0x33, 0x39, 0x2e, 0x31, 0x33, 0x34, 0x20, 0x38, 0x2e, 0x34, 0x34, 0x36, 0x2d, 0x31, 0x30, 0x2e, 0x34, 0x34, 0x33, 0x20, 0x31, 0x38, 0x2e, 0x30, 0x35, 0x32, 0x2d, 0x31, 0x39, 0x2e, 0x35, 0x39, 0x31, 0x20, 0x32, 0x38, 0x2e, 0x36, 0x36, 0x35, 0x2d, 0x32, 0x37, 0x2e, 0x32, 0x39, 0x33, 0x61, 0x31, 0x35, 0x32, 0x2e, 0x36, 0x32, 0x20, 0x31, 0x35, 0x32, 0x2e, 0x36, 0x32, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x20, 0x33, 0x34, 0x2e, 0x39, 0x36, 0x35, 0x2d, 0x31, 0x39, 0x2e, 0x30, 0x31, 0x31, 0x20, 0x31, 0x35, 0x33, 0x2e, 0x32, 0x34, 0x32, 0x20, 0x31, 0x35, 0x33, 0x2e, 0x32, 0x34, 0x32, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x20, 0x33, 0x34, 0x2e, 0x38, 0x39, 0x38, 0x20, 0x31, 0x38, 0x2e, 0x39, 0x37, 0x63, 0x31, 0x30, 0x2e, 0x36, 0x35, 0x34, 0x20, 0x37, 0x2e, 0x37, 0x34, 0x33, 0x20, 0x32, 0x30, 0x2e, 0x33, 0x30, 0x32, 0x20, 0x31, 0x36, 0x2e, 0x39, 0x36, 0x32, 0x20, 0x32, 0x38, 0x2e, 0x37, 0x39, 0x20, 0x32, 0x37, 0x2e, 0x34, 0x37, 0x2d, 0x2e, 0x37, 0x32, 0x34, 0x20, 0x31, 0x33, 0x2e, 0x34, 0x32, 0x37, 0x2d, 0x33, 0x2e, 0x31, 0x33, 0x32, 0x20, 0x32, 0x36, 0x2e, 0x34, 0x36, 0x35, 0x2d, 0x37, 0x2e, 0x32, 0x30, 0x34, 0x20, 0x33, 0x38, 0x2e, 0x39, 0x36, 0x31, 0x61, 0x31, 0x35, 0x32, 0x2e, 0x37, 0x36, 0x37, 0x20, 0x31, 0x35, 0x32, 0x2e, 0x37, 0x36, 0x37, 0x20, 0x30, 0x20, 0x30, 0x20, 0x31, 0x2d, 0x31, 0x36, 0x2e, 0x39, 0x35, 0x32, 0x20, 0x33, 0x35, 0x2e, 0x37, 0x30, 0x37, 0x7a, 0x6d, 0x2d, 0x32, 0x38, 0x2e, 0x37, 0x34, 0x2d, 0x36, 0x32, 0x2e, 0x39, 0x39, 0x38, 0x63, 0x30, 0x20, 0x39, 0x2e, 0x32, 0x33, 0x32, 0x20, 0x37, 0x2e, 0x34, 0x38, 0x32, 0x20, 0x31, 0x36, 0x2e, 0x37, 0x20, 0x31, 0x36, 0x2e, 0x37, 0x30, 0x32, 0x20, 0x31, 0x36, 0x2e, 0x37, 0x20, 0x39, 0x2e, 0x32, 0x31, 0x37, 0x20, 0x30, 0x20, 0x31, 0x36, 0x2e, 0x36, 0x39, 0x2d, 0x37, 0x2e, 0x34, 0x36, 0x36, 0x20, 0x31, 0x36, 0x2e, 0x36, 0x39, 0x2d, 0x31, 0x36, 0x2e, 0x37, 0x20, 0x30, 0x2d, 0x39, 0x2e, 0x31, 0x39, 0x36, 0x2d, 0x37, 0x2e, 0x34, 0x37, 0x33, 0x2d, 0x31, 0x36, 0x2e, 0x36, 0x39, 0x32, 0x2d, 0x31, 0x36, 0x2e, 0x36, 0x39, 0x2d, 0x31, 0x36, 0x2e, 0x36, 0x39, 0x32, 0x2d, 0x39, 0x2e, 0x32, 0x32, 0x20, 0x30, 0x2d, 0x31, 0x36, 0x2e, 0x37, 0x30, 0x31, 0x20, 0x37, 0x2e, 0x34, 0x39, 0x36, 0x2d, 0x31, 0x36, 0x2e, 0x37, 0x30, 0x31, 0x20, 0x31, 0x36, 0x2e, 0x36, 0x39, 0x32, 0x7a, 0x6d, 0x2d, 0x32, 0x31, 0x2e, 0x35, 0x37, 0x38, 0x20, 0x30, 0x63, 0x30, 0x20, 0x39, 0x2e, 0x32, 0x33, 0x32, 0x2d, 0x37, 0x2e, 0x34, 0x38, 0x20, 0x31, 0x36, 0x2e, 0x37, 0x2d, 0x31, 0x36, 0x2e, 0x37, 0x20, 0x31, 0x36, 0x2e, 0x37, 0x2d, 0x39, 0x2e, 0x32, 0x32, 0x36, 0x20, 0x30, 0x2d, 0x31, 0x36, 0x2e, 0x36, 0x38, 0x35, 0x2d, 0x37, 0x2e, 0x34, 0x36, 0x36, 0x2d, 0x31, 0x36, 0x2e, 0x36, 0x38, 0x35, 0x2d, 0x31, 0x36, 0x2e, 0x37, 0x20, 0x30, 0x2d, 0x39, 0x2e, 0x31, 0x39, 0x33, 0x20, 0x37, 0x2e, 0x34, 0x36, 0x2d, 0x31, 0x36, 0x2e, 0x36, 0x38, 0x39, 0x20, 0x31, 0x36, 0x2e, 0x36, 0x38, 0x36, 0x2d, 0x31, 0x36, 0x2e, 0x36, 0x38, 0x39, 0x20, 0x39, 0x2e, 0x32, 0x32, 0x20, 0x30, 0x20, 0x31, 0x36, 0x2e, 0x37, 0x20, 0x37, 0x2e, 0x34, 0x39, 0x36, 0x20, 0x31, 0x36, 0x2e, 0x37, 0x20, 0x31, 0x36, 0x2e, 0x36, 0x39, 0x7a, 0x22, 0x20, 0x66, 0x69, 0x6c, 0x6c, 0x3d, 0x22, 0x23, 0x34, 0x31, 0x39, 0x45, 0x44, 0x41, 0x22, 0x2f, 0x3e, 0x3c, 0x2f, 0x73, 0x76, 0x67, 0x3e, 0xa}, MediaType: "image/svg+xml"}, Description: "A message about etcd operator, a description of channels"},
					{Schema: "olm.package", Name: "", DefaultChannel: "", Icon: nil, Description: ""},
				},
				Bundles: []Bundle{
					{
						Schema:  "olm.bundle",
						Name:    "cockroachdb.v2.0.9",
						Package: "cockroachdb",
						Image:   "quay.io/openshift-community-operators/cockroachdb:v2.0.9",
						Properties: []property.Property{
							{Type: "olm.channel", Value: json.RawMessage(`{"name":"stable"}`)},
							{Type: "olm.package", Value: json.RawMessage(`{"packageName":"cockroachdb","version":"2.0.9"}`)},
						},
					},
					{
						Schema:  "olm.bundle",
						Name:    "cockroachdb.v2.1.11",
						Package: "cockroachdb",
						Image:   "quay.io/openshift-community-operators/cockroachdb:v2.1.11",
						Properties: []property.Property{
							{Type: "olm.channel", Value: json.RawMessage(`{"name":"stable","replaces":"cockroachdb.v2.1.1"}`)},
							{Type: "olm.package", Value: json.RawMessage(`{"packageName":"cockroachdb","version":"2.1.11"}`)},
						},
					},
					{
						Schema:  "olm.bundle",
						Name:    "cockroachdb.v2.1.1",
						Package: "cockroachdb",
						Image:   "quay.io/openshift-community-operators/cockroachdb:v2.1.1",
						Properties: []property.Property{
							{Type: "olm.channel", Value: json.RawMessage(`{"name":"stable","replaces":"cockroachdb.v2.0.9"}`)},
							{Type: "olm.package", Value: json.RawMessage(`{"packageName":"cockroachdb","version":"2.1.1"}`)},
						},
					},
					{
						Schema:  "olm.bundle",
						Name:    "cockroachdb.v3.0.7",
						Package: "cockroachdb",
						Image:   "quay.io/openshift-community-operators/cockroachdb:v3.0.7",
						Properties: []property.Property{
							{Type: "olm.channel", Value: json.RawMessage(`{"name":"stable-3.x"}`)},
							{Type: "olm.package", Value: json.RawMessage(`{"packageName":"cockroachdb","version":"3.0.7"}`)},
						},
					},
					{
						Schema:  "olm.bundle",
						Name:    "cockroachdb.v5.0.3",
						Package: "cockroachdb",
						Image:   "quay.io/openshift-community-operators/cockroachdb:v5.0.3",
						Properties: []property.Property{
							{Type: "olm.channel", Value: json.RawMessage(`{"name":"stable-5.x"}`)},
							{Type: "olm.package", Value: json.RawMessage(`{"packageName":"cockroachdb","version":"5.0.3"}`)},
						},
					},
					{
						Schema:  "olm.bundle",
						Name:    "etcdoperator-community.v0.6.1",
						Package: "etcd",
						Image:   "quay.io/operatorhubio/etcd:v0.6.1",
						Properties: []property.Property{
							{Type: "olm.package", Value: json.RawMessage(`{"packageName":"etcd","version":"0.6.1"}`)},
							{Type: "olm.gvk", Value: json.RawMessage(`{"group":"etcd.database.coreos.com","kind":"EtcdCluster","version":"v1beta2"}`)},
							{Type: "olm.channel", Value: json.RawMessage(`{"name":"alpha"}`)},
						},
						RelatedImages: []RelatedImage{{Name: "etcdv0.6.1", Image: "quay.io/coreos/etcd-operator@sha256:bd944a211eaf8f31da5e6d69e8541e7cada8f16a9f7a5a570b22478997819943"}},
					},
					{
						Schema:  "olm.bundle",
						Name:    "etcdoperator.v0.9.0",
						Package: "etcd",
						Image:   "quay.io/operatorhubio/etcd:v0.9.0",
						Properties: []property.Property{
							{Type: "olm.package", Value: json.RawMessage(`{"packageName":"etcd","version":"0.9.0"}`)},
							{Type: "olm.gvk", Value: json.RawMessage(`{"group":"etcd.database.coreos.com","kind":"EtcdBackup","version":"v1beta2"}`)},
							{Type: "olm.channel", Value: json.RawMessage(`{"name":"singlenamespace-alpha"}`)},
							{Type: "olm.channel", Value: json.RawMessage(`{"name":"clusterwide-alpha"}`)},
						},
						RelatedImages: []RelatedImage{{Name: "etcdv0.9.0", Image: "quay.io/coreos/etcd-operator@sha256:db563baa8194fcfe39d1df744ed70024b0f1f9e9b55b5923c2f3a413c44dc6b8"}},
					},
					{
						Schema:  "olm.bundle",
						Name:    "etcdoperator.v0.9.2",
						Package: "etcd",
						Image:   "quay.io/operatorhubio/etcd:v0.9.2",
						Properties: []property.Property{
							{Type: "olm.package", Value: json.RawMessage(`{"packageName":"etcd","version":"0.9.2"}`)},
							{Type: "olm.gvk", Value: json.RawMessage(`{"group":"etcd.database.coreos.com","kind":"EtcdRestore","version":"v1beta2"}`)},
							{Type: "olm.channel", Value: json.RawMessage(`{"name":"singlenamespace-alpha","replaces":"etcdoperator.v0.9.0"}`)},
						},
						RelatedImages: []RelatedImage{{Name: "etcdv0.9.2", Image: "quay.io/coreos/etcd-operator@sha256:c0301e4686c3ed4206e370b42de5a3bd2229b9fb4906cf85f3f30650424abec2"}},
					},
					{
						Schema:  "olm.bundle",
						Name:    "etcdoperator.v0.9.2-clusterwide",
						Package: "etcd",
						Image:   "quay.io/operatorhubio/etcd:v0.9.2-clusterwide",
						Properties: []property.Property{
							{Type: "olm.package", Value: json.RawMessage(`{"packageName":"etcd","version":"0.9.2-clusterwide"}`)},
							{Type: "olm.gvk", Value: json.RawMessage(`{"group":"etcd.database.coreos.com","kind":"EtcdBackup","version":"v1beta2"}`)},
							{Type: "olm.skipRange", Value: json.RawMessage(`">=0.9.0<=0.9.1"`)},
							{Type: "olm.skips", Value: json.RawMessage(`"etcdoperator.v0.6.1"`)},
							{Type: "olm.skips", Value: json.RawMessage(`"etcdoperator.v0.9.0"`)},
							{Type: "olm.channel", Value: json.RawMessage(`{"name":"clusterwide-alpha","replaces":"etcdoperator.v0.9.0"}`)},
						},
						RelatedImages: []RelatedImage{{Name: "etcdv0.9.2", Image: "quay.io/coreos/etcd-operator@sha256:c0301e4686c3ed4206e370b42de5a3bd2229b9fb4906cf85f3f30650424abec2"}},
					},
					{
						Schema:  "olm.bundle",
						Name:    "etcdoperator.v0.9.4",
						Package: "etcd",
						Image:   "quay.io/operatorhubio/etcd:v0.9.4",
						Properties: []property.Property{
							{Type: "olm.package", Value: json.RawMessage(`{"packageName":"etcd","version":"0.9.4"}`)},
							{Type: "olm.package.required", Value: json.RawMessage(`{"packageName":"test","versionRange":">=1.2.3<2.0.0-0"}`)},
							{Type: "olm.gvk", Value: json.RawMessage(`{"group":"etcd.database.coreos.com","kind":"EtcdBackup","version":"v1beta2"}`)},
							{Type: "olm.gvk.required", Value: json.RawMessage(`{"group":"testapi.coreos.com","kind":"Testapi","version":"v1"}`)},
							{Type: "olm.channel", Value: json.RawMessage(`{"name":"singlenamespace-alpha","replaces":"etcdoperator.v0.9.2"}`)},
						},
						RelatedImages: []RelatedImage{{Name: "etcdv0.9.2", Image: "quay.io/coreos/etcd-operator@sha256:66a37fd61a06a43969854ee6d3e21087a98b93838e284a6086b13917f96b0d9b"}},
					},
					{
						Schema:  "olm.bundle",
						Name:    "etcdoperator.v0.9.4-clusterwide",
						Package: "etcd",
						Image:   "quay.io/operatorhubio/etcd:v0.9.4-clusterwide",
						Properties: []property.Property{
							{Type: "olm.package", Value: json.RawMessage(`{"packageName":"etcd","version":"0.9.4-clusterwide"}`)},
							{Type: "olm.gvk", Value: json.RawMessage(`{"group":"etcd.database.coreos.com","kind":"EtcdBackup","version":"v1beta2"}`)},
							{Type: "olm.channel", Value: json.RawMessage(`{"name":"clusterwide-alpha","replaces":"etcdoperator.v0.9.2-clusterwide"}`)},
						},
						RelatedImages: []RelatedImage{{Name: "etcdv0.9.2", Image: "quay.io/coreos/etcd-operator@sha256:66a37fd61a06a43969854ee6d3e21087a98b93838e284a6086b13917f96b0d9b"}},
					},
					{
						Schema: "olm.bundle",
					},
				},
				Others: []Meta{
					{Schema: "unexpected", Package: "", Blob: json.RawMessage(`"{ "schema":  "unexpected" }`)},
				},
			},
		},
	}

	for _, s := range specs {
		t.Run(s.name, func(t *testing.T) {
			cfg, err := LoadDir(s.dir)
			s.assertion(t, err)
			if err == nil {
				require.NotNil(t, cfg)
				equalsDeclarativeConfig(t, *s.expected, *cfg)
			}
		})
	}
}