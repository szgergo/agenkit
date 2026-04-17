package manifest

import (
	"fmt"
	"os"

	"path/filepath"

	"go.yaml.in/yaml/v4"
)

func LoadGlobalManifest() (*Manifest, error) {

	userHomeDir, err := os.UserHomeDir()

	if err != nil {
		return nil, fmt.Errorf("error getting user home directory: %w", err)
	}

	return gatherConfigFromPath(createAgenkitConfigPath(userHomeDir))

}

func LoadProjectManifest() (*Manifest, error) {

	cwd, err := os.Getwd()

	if err != nil {
		return nil, fmt.Errorf("error getting current working directory: %w", err)
	}

	return gatherConfigFromPath(createAgenkitConfigPath(cwd))
}

func createAgenkitConfigPath(base string) string {
	return filepath.Join(base, ".agenkit", "agenkit.yml")
}

func gatherConfigFromPath(path string) (*Manifest, error) {

	rawConfig, err := os.ReadFile(path)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}

	var loadedConfig *Manifest

	loadedConfig, parseErr := parseConfig(rawConfig)

	if parseErr != nil {
		return nil, fmt.Errorf("invalid config file: %w", parseErr)
	}

	return loadedConfig, nil

}

func parseConfig(rawConfig []byte) (*Manifest, error) {

	var config Manifest

	err := yaml.Load(rawConfig, &config, yaml.WithKnownFields())

	if err != nil {
		return nil, err
	}

	return &config, nil
}
