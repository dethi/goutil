package envflag

import "strings"

// EnvironmentPrefix defines a string that will be implicitely prefixed to a
// flag name before looking it up in the environment variables.
var EnvironmentPrefix = ""

// parseEnv parses flags from environment variables.
// Flags already set will be ignored.
func (f *FlagSet) parseEnv(environ []string) error {
	m := f.formal

	env := make(map[string]string)
	for _, s := range environ {
		i := strings.Index(s, "=")
		if i < 1 {
			continue
		}
		env[s[0:i]] = s[i+1:]
	}

	for _, flag := range m {
		name := flag.Name
		_, set := f.actual[name]
		if set {
			continue
		}

		flag, alreadythere := m[name]
		if !alreadythere {
			if name == "help" || name == "h" { // special case for nice help message.
				f.usage()
				return ErrHelp
			}
			return f.failf("environment variable provided but not defined: %s", name)
		}

		envKey := normalizeNameForEnv(f.envPrefix, flag.Name)
		value, isSet := env[envKey]
		if !isSet {
			continue
		}

		if err := flag.Value.Set(value); err != nil {
			return f.failf("invalid value %q for environment variable %s: %v", value, name, err)
		}

		// update f.actual
		if f.actual == nil {
			f.actual = make(map[string]*Flag)
		}
		f.actual[name] = flag
	}
	return nil
}

func normalizeNameForEnv(prefix, name string) string {
	envKey := strings.ToUpper(name)
	if prefix != "" {
		envKey = prefix + "_" + envKey
	}
	envKey = strings.Replace(envKey, "-", "_", -1)
	envKey = strings.Replace(envKey, ".", "_", -1)
	return envKey
}

// NewFlagSetWithEnvPrefix returns a new empty flag set with the specified name,
// environment variable prefix, and error handling property.
func NewFlagSetWithEnvPrefix(name string, prefix string, errorHandling ErrorHandling) *FlagSet {
	f := NewFlagSet(name, errorHandling)
	f.envPrefix = prefix
	return f
}
