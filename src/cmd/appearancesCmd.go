package cmd

import (
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/simivar/tibia-sprites-exporter/src/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	AppearancesOutPath string
)

func init() {
	rootCmd.AddCommand(appearancesCmd)

	appearancesCmd.Flags().StringVar(&AppearancesOutPath, "out", defaultAppearancesOutPath(), "output JSON path for appearances export")
	_ = viper.BindPFlag("appearancesOut", appearancesCmd.Flags().Lookup("out"))
}

var appearancesCmd = &cobra.Command{
	Use:   "appearances",
	Short: "Exports full appearances data (protobuf) to a single JSON file",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Tibia appearances export running")

		catalogDir := app.ExpandPath(viper.GetString("catalog"))
		catalogFile := filepath.Join(catalogDir, "catalog-content.json")
		outPath := app.ExpandPath(viper.GetString("appearancesOut"))

		appearancesFileName := app.GetAppearancesFileNameFromCatalogContent(catalogFile)
		log.Info().Msgf("Appearances file name: %s", appearancesFileName)

		if err := app.ExportAppearancesJSON(catalogDir, appearancesFileName, outPath); err != nil {
			log.Fatal().Err(err).Msg("Appearances export failed")
		}

		log.Info().Str("out", outPath).Msg("Tibia appearances export finished")
	},
}

func defaultAppearancesOutPath() string {
	return app.ExpandPath("./output/appearances.json")
}