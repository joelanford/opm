package declcfg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/joelanford/opm/internal/property"
)

func FormatDir(configDir string) error {
	configDir, err := filepath.Abs(configDir)
	if err != nil {
		return fmt.Errorf("get absolute path: %v", err)
	}
	cfg, err := LoadDir(configDir)
	if err != nil {
		return fmt.Errorf("load configs: %v", err)
	}

	w := dirWalker{}
	otherFiles := map[string]struct{}{}
	bundleFiles := map[string]struct{}{}
	if err := w.WalkFiles(configDir, func(path string, r io.Reader) error {
		fileCfg, err := readJSON(r)
		if err != nil {
			return err
		}
		if fileCfg == nil {
			otherFiles[path] = struct{}{}
			return nil
		}
		if len(fileCfg.Others) > 0 {
			schemas := []string{}
			for _, b := range fileCfg.Others {
				schemas = append(schemas, fmt.Sprintf("%q", b.Schema))
			}
			set := sets.NewString(schemas...)
			return fmt.Errorf("unknown meta objects exist (schemas: %s), can't safely fmt due to possibility of breaking unknown file references", strings.Join(set.List(), ","))
		}
		for _, pkg := range fileCfg.Packages {
			if len(pkg.Properties) > 0 {
				propertyTypes := []string{}
				for _, p := range pkg.Properties {
					propertyTypes = append(propertyTypes, fmt.Sprintf("%q", p.Type))
				}
				set := sets.NewString(propertyTypes...)
				return fmt.Errorf("unknown properties on package %q (types: %s): can't safely fmt due to possibility of breaking unknown file references", pkg.Name, strings.Join(set.List(), ","))
			}
		}
		for _, b := range fileCfg.Bundles {
			props, err := property.Parse(b.Properties)
			if err != nil {
				return fmt.Errorf("parse properties for bundle %q: %v", b.Name, err)
			}
			if len(props.Others) > 0 {
				propertyTypes := []string{}
				for _, p := range props.Others {
					propertyTypes = append(propertyTypes, fmt.Sprintf("%q", p.Type))
				}
				set := sets.NewString(propertyTypes...)
				return fmt.Errorf("unknown properties on bundle %q (types: %s): can't safely fmt due to possibility of breaking unknown file references", b.Name, strings.Join(set.List(), ","))
			}
			for _, obj := range props.BundleObjects {
				if obj.IsRef() {
					relPath := filepath.Join(filepath.Dir(path), obj.GetRef())
					absObjPath, err := filepath.Abs(relPath)
					if err != nil {
						return fmt.Errorf("get absolute path for %q: %v", obj.GetRef(), err)
					}
					bundleFiles[absObjPath] = struct{}{}
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	for f := range bundleFiles {
		delete(otherFiles, f)
	}

	bakDir := fmt.Sprintf("%s.bak-%s", configDir, rand.String(5))
	if err := os.Rename(configDir, bakDir); err != nil {
		return fmt.Errorf("backup input config dir: %v", err)
	}

	if err := writeDir(*cfg, configDir); err != nil {
		_ = os.Rename(bakDir, configDir)
		return fmt.Errorf("write fmt-ed configs: %v", err)
	}
	for f := range otherFiles {
		f = strings.Replace(f, configDir, bakDir, 1)
		if err := copyTo(configDir, bakDir, f); err != nil {
			_ = os.RemoveAll(configDir)
			_ = os.Rename(bakDir, configDir)
			return fmt.Errorf("copy file %q: %v", f, err)
		}
	}
	if err := os.RemoveAll(bakDir); err != nil {
		return fmt.Errorf("remove input config dir: %v", err)
	}
	return nil
}

func copyTo(dest, base, fromPath string) error {
	relPath, err := filepath.Rel(base, fromPath)
	if err != nil {
		return fmt.Errorf("get relative path for %q: %v", fromPath, err)
	}
	destPath := filepath.Join(dest, relPath)
	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(destPath), 0777); err != nil {
		return err
	}
	dst, err := os.Create(destPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("copy %q to %q: %v", fromPath, destPath, err)
	}
	return nil
}
