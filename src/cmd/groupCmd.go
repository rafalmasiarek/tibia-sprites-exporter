package cmd

import (
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/simivar/tibia-sprites-exporter/src/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	GroupedOutputPath  string
	GroupsIndexOutPath string
)

func init() {
	rootCmd.AddCommand(groupCmd)

	groupCmd.Flags().StringVar(&SplitOutputPath, "splitOutput", defaultSplitOutputPath(), "split sprites output path")
	groupCmd.Flags().StringVar(&GroupedOutputPath, "groupedOutput", defaultGroupedOutputPath(), "grouped sprites by appearances.json output path")

	// Optional JSON index for downstream tooling (e.g., PHP outfit pipelines).
	groupCmd.Flags().StringVar(&GroupsIndexOutPath, "groupsIndexOut", "", "write groups index JSON (optional)")

	_ = viper.BindPFlag("splitOutput", groupCmd.Flags().Lookup("splitOutput"))
	_ = viper.BindPFlag("groupedOutput", groupCmd.Flags().Lookup("groupedOutput"))
	_ = viper.BindPFlag("groupsIndexOut", groupCmd.Flags().Lookup("groupsIndexOut"))
}

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Groups sprites from the Tibia client based on the appearances file",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Tibia Sprites group running")

		catalogDir := app.ExpandPath(viper.GetString("catalog"))
		catalogFile := filepath.Join(catalogDir, "catalog-content.json")
		splitOutput := app.ExpandPath(viper.GetString("splitOutput"))
		groupedOutput := app.ExpandPath(viper.GetString("groupedOutput"))
		groupsIndexOut := app.ExpandPath(viper.GetString("groupsIndexOut"))

		appearancesFileName := app.GetAppearancesFileNameFromCatalogContent(catalogFile)
		log.Info().Msgf("Appearances file name: %s", appearancesFileName)

		app.GroupSplitSprites(catalogDir, appearancesFileName, splitOutput, groupedOutput, groupsIndexOut)

		log.Info().Msg("Tibia Sprites group finished")
	},
}

func defaultGroupedOutputPath() string {
	return app.ExpandPath(
		"./output/grouped",
	)
}
