package app

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

var (
	// Used for flags.
	cfgFile    string
	configPath string
	debug      bool
	quiet      bool
	trace      bool
	useNew     bool

	rootCmd = &cobra.Command{
		Use:     "mdm ",
		Short:   "音乐下载工具",
		Long:    `音乐下载工具`,
		Version: "0.1.0",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) <= 0{
				return errors.New("参数数量不对")
			}

			viper.Set("urls", args)

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if quiet {
				logrus.SetOutput(ioutil.Discard)
			} else {
				if debug {
					logrus.SetLevel(logrus.DebugLevel)
					logrus.Debugf("Loglevel set to [%v]", logrus.DebugLevel)
				}
				if trace {
					logrus.SetLevel(logrus.TraceLevel)
					logrus.Tracef("Loglevel set to [%v]", logrus.TraceLevel)
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return Download()
		},
	}
)

func MainErr() error {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "db2db.toml",
		"Specify an alternate cluster Toml file")
	rootCmd.PersistentFlags().StringVarP(&configPath, "configPath", "p", ".",
		"Specify an alternate config location")
	ViperIntBindAndSetP(rootCmd,"threadNumber", "n", 1,
		"线程数量")
	ViperStringBindAndSetP(rootCmd,"url", "u", "",
		"访问url")
	ViperStringBindAndSetP(rootCmd,"downloadDir", "d", "./songs",
		"下载")

	return rootCmd.Execute()

}

func initConfig() {
	if configPath != "" {
		viper.AddConfigPath(configPath)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
	}
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile(".db2db.toml.bak")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func ViperBoolBindAndSetP(command *cobra.Command, name string, shorthand string, defaultValue bool, usage string) {
	command.PersistentFlags().BoolP(name, shorthand, defaultValue, usage)
	err := viper.BindPFlag(name, command.PersistentFlags().Lookup(name))
	if err != nil {
		logrus.Fatalf("设置参数%s err:%s", name, err.Error())
	}
	viper.SetDefault(name, defaultValue)
}

func ViperBoolVarBindAndSetP(command *cobra.Command, p *bool, name string, shorthand string, defaultValue bool, usage string) {
	command.PersistentFlags().BoolVarP(p, name, shorthand, defaultValue, usage)
	err := viper.BindPFlag(name, command.PersistentFlags().Lookup(name))
	if err != nil {
		logrus.Fatalf("设置参数%s err:%s", name, err.Error())
	}
	viper.SetDefault(name, defaultValue)
}
func ViperStringBindAndSetP(command *cobra.Command, name string, shorthand string, defaultValue string, usage string) {
	command.PersistentFlags().StringP(name, shorthand, defaultValue, usage)
	err := viper.BindPFlag(name, command.PersistentFlags().Lookup(name))
	if err != nil {
		logrus.Fatalf("设置参数%s err:%s", name, err.Error())
	}
	viper.SetDefault(name, defaultValue)
}

func ViperIntBindAndSetP(command *cobra.Command, name string, shorthand string, defaultValue int, usage string) {
	command.PersistentFlags().IntP(name, shorthand, defaultValue, usage)
	err := viper.BindPFlag(name, command.PersistentFlags().Lookup(name))
	if err != nil {
		logrus.Fatalf("设置参数%s err:%s", name, err.Error())
	}
	viper.SetDefault(name, defaultValue)
}

func ViperBoolBindAndSet(command *cobra.Command, name string, shorthand string, defaultValue bool, usage string) {
	command.Flags().BoolP(name, shorthand, defaultValue, usage)
	err := viper.BindPFlag(name, command.Flags().Lookup(name))
	if err != nil {
		logrus.Fatalf("设置参数%s err:%s", name, err.Error())
	}
	viper.SetDefault(name, defaultValue)
}

func ViperBoolVarBindAndSet(command *cobra.Command, p *bool, name string, shorthand string, defaultValue bool, usage string) {
	command.Flags().BoolVarP(p, name, shorthand, defaultValue, usage)
	err := viper.BindPFlag(name, command.Flags().Lookup(name))
	if err != nil {
		logrus.Fatalf("设置参数%s err:%s", name, err.Error())
	}
	viper.SetDefault(name, defaultValue)
}
func ViperStringBindAndSet(command *cobra.Command, name string, shorthand string, defaultValue string, usage string) {
	command.Flags().StringP(name, shorthand, defaultValue, usage)
	err := viper.BindPFlag(name, command.Flags().Lookup(name))
	if err != nil {
		logrus.Fatalf("设置参数%s err:%s", name, err.Error())
	}
	viper.SetDefault(name, defaultValue)
}

func ViperIntBindAndSet(command *cobra.Command, name string, shorthand string, defaultValue int, usage string) {
	command.Flags().IntP(name, shorthand, defaultValue, usage)
	err := viper.BindPFlag(name, command.Flags().Lookup(name))
	if err != nil {
		logrus.Fatalf("设置参数%s err:%s", name, err.Error())
	}
	viper.SetDefault(name, defaultValue)
}