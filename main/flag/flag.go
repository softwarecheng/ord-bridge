package flag

import (
	"fmt"
	"os"

	"github.com/softwarecheng/ord-bridge/common/db/maintain"
	"github.com/softwarecheng/ord-bridge/common/util"
	"github.com/softwarecheng/ord-bridge/main/g"
	"github.com/spf13/cobra"
)

func checkChain(chain string) (err error) {
	switch chain {
	case "testnet":
		fallthrough
	case "mainnet":
	default:
		return fmt.Errorf("unsupported chain: %s", chain)
	}
	return nil
}

func ParseCmdParams() error {
	var chain string
	var dbDir string
	var confPath string
	rootCmd := &cobra.Command{
		Use:   "ord-bridge",
		Short: "ord bridge, a ord inscription data server",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := initConf(confPath)
			if err != nil {
				return fmt.Errorf("init config error: %v", err)
			}
			err = initLog(&g.Cfg.Log)
			if err != nil {
				return fmt.Errorf("init log error: %v", err)
			}
			return nil
		},
	}
	rootCmd.Flags().StringVarP(&confPath, "conf", "", "config.yaml", "run with config file")
	rootCmd.Long = `ord bridge, a ord inscription data server. 
	rpc service for ord data , convert or forward data inerface from (ordinals.com) database`

	rootCmd.Example = `	'ordx-bridge init --chain=testnet' or 'ordx-bridge init --chain=mainnet'
	'ordx-bridge --config=config.yaml', default config.yaml in curent directory
	'ordx-bridge dbgc --db ./db/testnet'
	'ordx-bridge import --chain=testnet --db ./db/testnet --data=/data/ordx-data-backup/ord-latest/export.ordx'`

	rootCmd.PersistentFlags().StringVarP(&dbDir, "db", "", "", "specify the db directory")
	rootCmd.PersistentFlags().StringVarP(&chain, "chain", "", "", "Specify the chain as testnet or mainnet")

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "generate config config.yaml in curent directory",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := checkChain(chain)
			if err != nil {
				cmd.Usage()
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := genDefaultCfg(chain)
			return err
		},
	}

	rootCmd.AddCommand(initCmd)

	var importDataPath string
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "import data from ord(ordinals.com) database, use ord0.16.0-ordx version complete",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := checkChain(chain)
			if err != nil {
				cmd.Usage()
				return err
			}
			if _, err := os.Stat(importDataPath); os.IsNotExist(err) {
				cmd.Usage()
				return fmt.Errorf("import data path isn't exist")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dbDir, err := util.GetAbsolutePath(dbDir)
			if err != nil {
				return fmt.Errorf("get db dir failed, error: %s", err)
			}
			err = importOrdData(chain, dbDir, importDataPath)
			if err != nil {
				return fmt.Errorf("ord import file %s failed, error: %v", importDataPath, err)
			}
			return nil
		},
	}
	importCmd.Flags().StringVarP(&importDataPath, "data", "", "", "specify the import data path")
	rootCmd.AddCommand(importCmd)

	dbgcCmd := &cobra.Command{
		Use:   "dgbc",
		Short: "gc database log",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(dbDir); os.IsNotExist(err) {
				cmd.Usage()
				return fmt.Errorf("db path isn't exist")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dbDir, err := util.GetAbsolutePath(dbDir)
			if err != nil {
				return fmt.Errorf("get db dir failed, error: %s", err)
			}
			err = maintain.GcFromPath(dbDir, 0.5)
			if err != nil {
				return err
			}
			return nil
		},
	}
	rootCmd.AddCommand(dbgcCmd)

	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
