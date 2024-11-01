package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitConfig(cfgFile string, target *string, ignoreExtra []string, configuration *Config) {
	homedir, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Configure all defaults here
	viper.SetDefault("target", homedir)
	viper.SetDefault("ignore_patterns", []string{".gitignore", ".dotlink.yaml", "dist/"})

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in current directory with name ".dotlink" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".dotlink.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		log.Fatal(err)
	}

	if err := viper.Unmarshal(configuration); err != nil {
		log.Fatalf("could not unmarshal configuration into struct: %e", err)
	}
	//
	// Check if target is set, if not assign it to value from config
	// wrap in cobra.CheckErr to exit if any error is returned.
	cobra.CheckErr(GetTarget(target))
	configuration.Target = *target
	configuration.extendIgnorePatterns(ignoreExtra...)
}
