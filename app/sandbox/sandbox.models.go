package sandbox

import "../db"

type sandbox struct {
	compilerName     string
	sourceType       string
	sandboxStoreRoot string
	codeDir          string
	contestID        string
	questionID       string
	userID           string
	outputFilename   string
	inputFilename    string
	compileStmt      string
	inputCode        string
	result           string
	timeoutValue     int
	executionTime    string
	containerID      string
	store            db.Store
}
