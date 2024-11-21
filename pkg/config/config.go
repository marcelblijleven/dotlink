package config

type MapPattern struct {
	Source string `mapstructure:"source"`
	Target string `mapstructure:"target"`
}

type OsSettings struct {
	IgnorePatterns []string     `mapstructure:"ignore_patterns"`
	MapPatterns    []MapPattern `mapstructure:"map_patterns"`
}

type OsSpecific map[string]OsSettings

type Config struct {
	Target                   string       `mapstructure:"target"`
	OsSpecific               OsSpecific   `mapstructure:"os_specific"`
	IgnorePatternsFromConfig []string     `mapstructure:"ignore_patterns"`
	mapPatternsFromConfig    []MapPattern `mapstructure:"map_patterns"`
}

// extendIgnorePatterns adds the provided strings to the IgnorePatterns slice.
// It will also add the defaults.
func (c *Config) extendIgnorePatterns(flagPatterns ...string) {
	ignorePatterns := []string{".dotlink.yaml"}
	ignorePatterns = append(ignorePatterns, flagPatterns...)
	c.IgnorePatternsFromConfig = append(c.IgnorePatternsFromConfig, ignorePatterns...)
}

// IgnorePatterns returns the ignore patterns from the configuration file,
// including OS specific patterns.
func (c Config) IgnorePatterns() []string {
	patterns := c.IgnorePatternsFromConfig
	os_specific, exists := c.OsSpecific[getMappedGoos()]

	if exists {
		patterns = append(patterns, os_specific.IgnorePatterns...)
	}
	return patterns
}

// MapPatterns returns the map patterns from the configuration file, including
// OS specific patterns.
func (c Config) MapPatterns() []MapPattern {
	patterns := c.mapPatternsFromConfig
	os_specific, exists := c.OsSpecific[getMappedGoos()]

	if exists {
		patterns = append(patterns, os_specific.MapPatterns...)
	}
	return patterns
}
