// Package cmd /*
package cmd

import (
	"fmt"
	"github.com/quantosnetwork/Quantos/address"

	"github.com/spf13/cobra"
)

// addressCmd represents the address command
var addressCmd = &cobra.Command{
	Use:   "address",
	Short: "QBit (Quantos) address manager",
	Long:  `QBit (Quantos) address manager`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("address called")
	},
}

//var bip39 bool

var newAddressCmd = &cobra.Command{
	Use:   "new-address",
	Short: "create a new wallet address",
	Long:  "create a new wallet address",
	Run: func(cmd *cobra.Command, args []string) {

		//	bip39 := true

		compress, _ := cmd.Flags().GetBool("compress")
		pass, _ := cmd.Flags().GetString("pass")
		mnemonic, _ := cmd.Flags().GetString("mnemonic")
		address.NewQBITAddress(compress, pass, mnemonic)
		//address.NewAddress(compress, pass, mnemonic)

	},
}

func init() {
	rootCmd.AddCommand(addressCmd)

	newAddressCmd.Flags().BoolP("bip39", "b", false, "mnemonic code for generating deterministic keys")
	newAddressCmd.Flags().BoolP("compress", "c", true, "generate a compressed public key")
	newAddressCmd.Flags().String("pass", "", "protect bip39 address with passphrase")
	newAddressCmd.Flags().Int("number", 10, "set number of keys to generate")
	newAddressCmd.Flags().String("mnemonic", "", "optional list of words to re-generate a root key")

	addressCmd.AddCommand(newAddressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addressCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addressCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
