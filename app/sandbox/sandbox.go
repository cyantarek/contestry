package sandbox

import (
	"os"
	"time"
	"os/exec"
	"bytes"
	"fmt"
	"../db"
	"../models"
	"log"
)

type ISandbox interface {
	Prepare()
	Run()
	GetResult() string
}

/*
Types of bash commands used

01. Creating Docker container:
"sh", "-c", " " + "docker run -i -d -v \"$(pwd)/"+ s.codeDir + ":/code\" sandbox:v2"
>> sh -c docker run -i -d -v "./code:/code" sandbox:v2

02. Compile source code:
"sh", "-c", " " + "docker exec -i " + s.containerID + " " + "bash -c 'cd /code; go build -o out " + s.inputFilename + "'"
>> sh -c docker exec -i b245da5563 bash -c 'cd /code; go build -o out solution.go'

03. Run output program:
"sh", "-c", "docker exec -i " + s.containerID + " bash -c 'yes " + question.Input + "  | ./code/out'"
>> sh -c docker exec -i b245da5563 bash -c 'yes "Hello" | ./code/out'

04. Killing and Removing Docker Container:
"sh", "-c", " " + "docker kill " + s.containerID + " && docker rm " + s.containerID
>> sh -c docker kill b245da5563 && docker rm b245da5563
 */

func NewSandbox(sandboxRoot, contestID, questionID, userID, inputCode, sourceType string) *sandbox {
	return &sandbox{
		sandboxStoreRoot: sandboxRoot,
		contestID:        contestID,
		questionID:       questionID,
		userID:           userID,
		outputFilename:   "out",
		inputCode:        inputCode,
		sourceType:       sourceType,
		store:            db.GetGormDb("test.db", "sqlite3"),
	}
}

func (s *sandbox) Prepare() {
	s.codeDir = s.sandboxStoreRoot + "/" + s.contestID + "/" + s.questionID + "/" + s.userID
	os.RemoveAll(s.codeDir)
	os.MkdirAll(s.codeDir, os.ModePerm)

	switch s.sourceType {
	case "go":
		s.inputFilename = "solution.go"
		s.compilerName = "go"
		s.compileStmt = "bash -c 'cd /code; go build -o out " + s.inputFilename + "'"
	case "c":
		s.inputFilename = "solution.c"
		s.compilerName = "gcc"
		s.compileStmt = "bash -c 'cd /code; " + s.compilerName + " -o out " + s.inputFilename + "'"
	case "cpp":
		s.inputFilename = "solution.cpp"
		s.compilerName = "g++"
		s.compileStmt = "bash -c 'cd /code; " + s.compilerName + " -o out " + s.inputFilename + "'"
	case "java":
		s.inputFilename = "solution.jar"
		s.compilerName = "javac"
	case "python":
		s.inputFilename = "solution.py"
		s.compilerName = "python3"
		s.compileStmt = "bash -c 'cd /code; " + s.compilerName + " " + s.inputFilename + "'"
	}

	savedCodeFile, err := os.OpenFile(s.codeDir+"/"+s.inputFilename, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Println(err.Error())
	}
	bs, err := savedCodeFile.WriteAt([]byte(s.inputCode), 0)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println(bs, "bytes written")
	}
}

func (s *sandbox) Run() {
	start := time.Now()
	dockerLaunchStmt := "docker run -i -d -v \"$(pwd)/"+ s.codeDir + ":/code\" sandbox:v2"
	var dockerCmd = exec.Command("sh", "-c", " " + dockerLaunchStmt)
	var stdErr bytes.Buffer
	dockerCmd.Stderr = &stdErr
	out, err := dockerCmd.Output()
	if err != nil {
		s.result = stdErr.String()
	} else {
		s.containerID = string(out)[:len(string(out))-1]
	}

	go s.timeoutContainerCleaner(s.containerID, s.contestID, s.questionID, s.userID)
	go s.timeoutContainerCleaner(s.containerID, s.contestID, s.questionID, s.userID)

	compileStmt := "docker exec -i " + s.containerID + " " + s.compileStmt
	var compileCmd = exec.Command("sh", "-c", compileStmt)
	compileCmd.Stderr = &stdErr

	//input := "hello world"
	var question models.Question
	s.store.GetSingle(&question, "id = ?", s.questionID)

	//var runCmd = exec.Command("sh", "-c", "docker exec -i " + s.containerID+ " yes " + input + " | ./code/out") ./code/out
	var runCmd = exec.Command("sh", "-c", "docker exec -i " + s.containerID + " bash -c 'yes " + question.Input + "  | ./code/out'")
	runCmd.Stderr = &stdErr


	if s.sourceType == "python" {
		out, err = compileCmd.Output()
		if err != nil {
			s.result = stdErr.String()
		} else {
			s.result = string(out)

			var solution models.Solution
			s.store.GetRaw(&solution, "select * from solutions where solutions.contest_id = ? and solutions.question_id = ? and solutions.user_id = ?", s.contestID, s.questionID, s.userID)
			if solution.ID > 0 {
				solution.Result = s.result[:len(s.result)-1]
				if solution.Result == question.CorrectAns {
					solution.Point = question.Point
				} else {
					solution.Point = 0.0
				}
				solution.ExecTime = time.Since(start).String()
				s.store.UpdateData(&solution)
			} else {
				var solution models.Solution
				solution.ContestID = s.contestID
				solution.QuestionID = s.questionID
				solution.UserID = s.userID
				solution.Result = s.result[:len(s.result)-1]
				if solution.Result == question.CorrectAns {
					solution.Point = question.Point
				} else {
					solution.Point = 0.0
				}
				solution.ExecTime = time.Since(start).String()
				s.store.InsertData(&solution)
			}
		}
	} else {
		_, err = compileCmd.Output()
		if err != nil {
			s.result = stdErr.String()

			var solution models.Solution
			s.store.GetRaw(&solution, "select * from solutions where solutions.contest_id = ? and solutions.question_id = ? and solutions.user_id = ?", s.contestID, s.questionID, s.userID)
			if solution.ID > 0 {
				solution.Result = s.result
				solution.Point = 0.0
				solution.ExecTime = time.Since(start).String()
				s.store.UpdateData(&solution)
			} else {
				var solution models.Solution
				solution.ContestID = s.contestID
				solution.QuestionID = s.questionID
				solution.UserID = s.userID
				solution.Result = s.result[:len(s.result)-1]
				solution.Point = 0.0
				solution.ExecTime = time.Since(start).String()
				s.store.InsertData(&solution)
			}
		} else {
			out, err := runCmd.Output()
			if err != nil {
				s.result = stdErr.String()

				var solution models.Solution
				s.store.GetRaw(&solution, "select * from solutions where solutions.contest_id = ? and solutions.question_id = ? and solutions.user_id = ?", s.contestID, s.questionID, s.userID)
				if solution.ID > 0 {
					solution.Result = s.result
					solution.Point = 0.0
					solution.ExecTime = time.Since(start).String()
					s.store.UpdateData(&solution)
				} else {
					var solution models.Solution
					solution.ContestID = s.contestID
					solution.QuestionID = s.questionID
					solution.UserID = s.userID
					solution.Result = s.result[:len(s.result)-1]
					solution.Point = 0.0
					solution.ExecTime = time.Since(start).String()
					s.store.InsertData(&solution)
				}
			} else {
				s.result = string(out)

				var solution models.Solution
				s.store.GetRaw(&solution, "select * from solutions where solutions.contest_id = ? and solutions.question_id = ? and solutions.user_id = ?", s.contestID, s.questionID, s.userID)
				if solution.ID > 0 {
					solution.Result = s.result[:len(s.result)-1]
					if solution.Result == question.CorrectAns {
						fmt.Println(solution.Result, question.CorrectAns)
						solution.Point = question.Point
					} else {
						fmt.Println("Wrong")
						fmt.Println(len(question.CorrectAns), len(solution.Result))
						fmt.Println(solution.Result)
						solution.Point = 0.0
					}
					solution.ExecTime = time.Since(start).String()
					s.store.UpdateData(&solution)
				} else {
					var solution models.Solution
					solution.ContestID = s.contestID
					solution.QuestionID = s.questionID
					solution.UserID = s.userID
					solution.Result = s.result[:len(s.result)-1]
					if solution.Result == question.CorrectAns {
						solution.Point = question.Point
					} else {
						solution.Point = 0.0
					}
					solution.ExecTime = time.Since(start).String()
					s.store.InsertData(&solution)
				}
			}
		}
	}

	dockerKillStmt := "docker kill " + s.containerID + " && docker rm " + s.containerID
	var dockerKillCmd = exec.Command("sh", "-c", " " + dockerKillStmt)

	dockerKillCmd.Stderr = &stdErr
	_, err = dockerKillCmd.Output()
	if err != nil {
		log.Println(stdErr.String())
	}
}


func (s *sandbox) timeoutContainerCleaner(containerId, contestId, questionId, userId string) {
	fmt.Println("Docker Cleaner Launched")
	time.Sleep(time.Second*5)
	//after 5 seconds, if container still is alive, clean it
	dockerKillStmt := "docker kill " + containerId + " && docker rm " + containerId
	var dockerKillCmd = exec.Command("sh", "-c", " " + dockerKillStmt)

	_, err := dockerKillCmd.Output()
	if err != nil {
		log.Println("CCCC", err.Error())
	} else {
		var solution models.Solution
		s.store.GetRaw(&solution, "select * from solutions where solutions.contest_id = ? and solutions.question_id = ? and solutions.user_id = ?", contestId, questionId, userId)
		if solution.ID > 0 {
			solution.Result = "Time Out"
			solution.ExecTime = "5.05s"
			s.store.UpdateData(&solution)
		} else {
			var solution models.Solution
			solution.ContestID = contestId
			solution.QuestionID = questionId
			solution.UserID = userId
			solution.Result = "Time Out"
			solution.ExecTime = "5.05s"
			s.store.InsertData(&solution)
		}
	}
	fmt.Println("Docker Cleaner Job Finished")
}

func (s *sandbox) GetResult() string {
	return s.result
}