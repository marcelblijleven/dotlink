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
	Target         string       `mapstructure:"target"`
	OsSpecific     OsSpecific   `mapstructure:"os_specific"`
	ignorePatterns []string     `mapstructure:"ignore_patterns"`
	mapPatterns    []MapPattern `mapstructure:"map_patterns"`
}

func NewConfig(target string, ignorePatterns []string, mapPatterns []MapPattern) Config {
	return Config{
		Target:         target,
		ignorePatterns: ignorePatterns,
		mapPatterns:    mapPatterns,
		OsSpecific:     make(OsSpecific),
	}
}

// extendIgnorePatterns adds the provided strings to the IgnorePatterns slice.
// It will also add the defaults.
func (c *Config) extendIgnorePatterns(flagPatterns ...string) {
	ignorePatterns := []string{".git", ".dotlink.yaml"}
	ignorePatterns = append(ignorePatterns, flagPatterns...)
	c.ignorePatterns = append(c.ignorePatterns, ignorePatterns...)
}

// IgnorePatterns returns the ignore patterns from the configuration file,
// including OS specific patterns.
func (c Config) IgnorePatterns() []string {
	patterns := c.ignorePatterns
	os_specific, exists := c.OsSpecific[getMappedGoos()]

	if exists {
		patterns = append(patterns, os_specific.IgnorePatterns...)
	}
	return patterns
}

// MapPatterns returns the map patterns from the configuration file, including
// OS specific patterns.
func (c Config) MapPatterns() []MapPattern {
	patterns := c.mapPatterns
	os_specific, exists := c.OsSpecific[getMappedGoos()]

	if exists {
		patterns = append(patterns, os_specific.MapPatterns...)
	}
	return patterns
}
