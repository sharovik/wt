package configuration

type Config struct {
	WorkingBranch      string
	DestinationBranch  string
	VcsType            string
	Path               string
	Ext                string
	PathToIgnoreFile   string
	DisplayTemplate    string
	IgnoreFromAnalysis string
	MaxAnalysisDepth   int
	WithToBeChecked    bool
	Version            bool
	CpuProfile         string
	MemProfile         string
}

var C Config