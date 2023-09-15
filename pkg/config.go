package rsync

const WRAPPER_CONFIG = ".rsync-wrapper.yaml"

type RsyncWrapper struct {
	ExcludeDirs            []string `yaml:"excludeDirs,omitempty"`
	DestinationPath        string   `yaml:"destinationPath"`
	DestinationAddress     string   `yaml:"destinationAddress"`
	RemoveDestinationFiles bool     `yaml:"removeDestinationFiles,omitempty"`
	DryRun                 bool     `yaml:"dryRun,omitempty"`
}
