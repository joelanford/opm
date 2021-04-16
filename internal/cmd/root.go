package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/operator-framework/operator-registry/pkg/action"
	"github.com/operator-framework/operator-registry/pkg/image/containerdregistry"
	"github.com/operator-framework/operator-registry/pkg/registry"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/joelanford/opm/internal/declcfg"
	"github.com/joelanford/opm/internal/model"
	"github.com/joelanford/opm/internal/property"
)

func Execute() {
	if err := newCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
func newCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "opm",
		Short: "Manage your OLM indexes",
	}
	alpha := &cobra.Command{
		Use:   "alpha",
		Short: "Small focused tools to manage your OLM declarative config-based indexes.",
	}

	root.AddCommand(alpha)
	alpha.AddCommand(
		newBlobCmd(),
		newValidateCmd(),
		newFmtCmd(),
	)
	return root
}

func newBlobCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "blob <schema> <image>",
		Short: "Generate a declarative config blob for the specified schema",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			schemas := strings.Split(args[0], ",")
			imageRef := args[1]

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			var encodes []func(cfg declcfg.DeclarativeConfig) error

			for _, schema := range schemas {
				switch schema {
				case "olm.package":
					encodes = append(encodes, func(cfg declcfg.DeclarativeConfig) error { return enc.Encode(cfg.Packages[0]) })
				case "olm.bundle":
					encodes = append(encodes, func(cfg declcfg.DeclarativeConfig) error { return enc.Encode(cfg.Bundles[0]) })
				default:
					log.Fatalf("cannot generate blob for schema %q", schema)
				}
			}

			logger := logrus.New()
			logger.SetOutput(ioutil.Discard)
			nullLogger := logrus.NewEntry(logger)
			logrus.SetOutput(ioutil.Discard)

			reg, err := containerdregistry.NewRegistry(containerdregistry.WithLog(nullLogger))
			if err != nil {
				log.Fatal(err)
			}
			defer reg.Destroy()
			extractor := action.NewImageBundleExtractor(imageRef, reg, nullLogger)
			bundle, err := extractor.ExtractBundle(cmd.Context())
			if err != nil {
				log.Fatal(err)
			}

			bundles, err := registry.ConvertRegistryBundleToModelBundles(bundle)
			if err != nil {
				log.Fatal(err)
			}
			if len(bundles) == 0 {
				log.Fatal("no bundles extracted from image")
			}

			mPackage := &model.Package{
				Name:        bundles[0].Package.Name,
				Description: bundles[0].Package.Description,
				Channels:    make(map[string]*model.Channel),
			}
			if bundles[0].Package.Icon != nil && len(bundles[0].Package.Icon.Data) > 0 {
				mPackage.Icon = &model.Icon{
					Data:      bundles[0].Package.Icon.Data,
					MediaType: bundles[0].Package.Icon.MediaType,
				}

			}
			for _, b := range bundles {
				ch := b.Channel
				mCh := &model.Channel{
					Package: mPackage,
					Name:    ch.Name,
					Bundles: make(map[string]*model.Bundle),
				}
				mPackage.Channels[ch.Name] = mCh
				if bundle.Annotations != nil && bundle.Annotations.DefaultChannelName == ch.Name {
					mPackage.DefaultChannel = mCh
				}
				mB := &model.Bundle{
					Package:  mPackage,
					Channel:  mCh,
					Name:     b.Name,
					Image:    b.Image,
					Replaces: b.Replaces,
					Skips:    b.Skips,
					Objects:  b.Objects,
					CsvJSON:  b.CsvJSON,
				}
				for _, p := range b.Properties {
					mB.Properties = append(mB.Properties, property.Property(p))
				}
				for _, ri := range b.RelatedImages {
					mB.RelatedImages = append(mB.RelatedImages, model.RelatedImage(ri))
				}
				mCh.Bundles[b.Name] = mB
			}

			cfg := declcfg.ConvertFromModel(model.Model{mPackage.Name: mPackage})
			for _, encode := range encodes {
				if err := encode(cfg); err != nil {
					log.Fatal(err)
				}
			}
		},
	}
}

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <configsDir>",
		Short: "Validate a declarative config directory",
		Run: func(_ *cobra.Command, args []string) {
			configsDir := args[0]
			cfg, err := declcfg.LoadDir(configsDir)
			if err != nil {
				log.Fatal(err)
			}

			// TODO(joelanford): Use default schema, allow user-provided schema
			if _, err := declcfg.ConvertToModel(*cfg); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func newFmtCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fmt <configsDir>",
		Short: "Re-format a declarative config directory to a standard layout",
		Run: func(_ *cobra.Command, args []string) {
			configsDir := args[0]
			if err := declcfg.FormatDir(configsDir); err != nil {
				log.Fatal(err)
			}
		},
	}
}
