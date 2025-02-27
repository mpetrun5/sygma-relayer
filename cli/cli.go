// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ChainSafe/chainbridge-core/flags"
)

var (
	rootCMD = &cobra.Command{
		Use: "",
	}
)

func init() {
	flags.BindFlags(rootCMD)
	rootCMD.PersistentFlags().String("name", "", "relayer name")
	_ = viper.BindPFlag("name", rootCMD.PersistentFlags().Lookup("name"))
}

func Execute() {
	rootCMD.AddCommand(runCMD, peerInfoCMD)
	if err := rootCMD.Execute(); err != nil {
		log.Fatal().Err(err).Msg("failed to execute root cmd")
	}
}
