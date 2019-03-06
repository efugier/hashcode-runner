package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

const header = "\033[95m"
const blue = "\033[94m"
const green = "\033[92m"
const orange = "\033[93m"
const red = "\033[91m"
const endc = "\033[0m"
const bold = "\033[1m"
const underline = "\033[4m"

type result struct {
	oldScore int
	newScore int
	status   string
}

func makebold(s string) string {
	return bold + s + endc
}
func okgreen(s string) string {
	return green + s + endc
}
func okblue(s string) string {
	return blue + s + endc
}
func warnorange(s string) string {
	return orange + s + endc
}
func nokred(s string) string {
	return red + s + endc
}

func main() {
	datasets := flag.String("datasets", "A", "datasets to evaluate")
	model := flag.String("model", "./model.sh", "model model executable")
	scorer := flag.String("scorer", "./scorer.sh", "scorer executable")
	datafolder := flag.String("datafolder", "data", "folder containing datasets")
	submissionsfolder := flag.String("submissionsfolder", "submissions", "folder containing submissions")
	realtimeoutput := flag.Bool("realtimeoutput", false, "print the output of the models in stdout in real time (always the case if there is only one dataset)")

	flag.Parse()

	results := make([]result, len(*datasets))

	if len(*datasets) < 2 {
		*realtimeoutput = true
	}

	// Run the computations in parallel
	{
		wg := sync.WaitGroup{}
		for i, dataset := range *datasets {
			results[i] = result{-1, -1, nokred("worse")}
			wg.Add(1)
			go func(res *result, dataset string) {
				defer wg.Done()
				testDataset(dataset, *model, *scorer, res, *datafolder, *submissionsfolder, *realtimeoutput)
			}(&results[i], string(dataset))
		}
		wg.Wait()
	}

	resTable := tablewriter.NewWriter(os.Stdout)
	resTable.SetHeader([]string{"Test case", "Old score", "New score", "Status"})
	resTable.SetRowLine(true)
	for i, dataset := range *datasets {
		r := results[i]
		resTable.Append([]string{makebold(string(dataset)), strconv.Itoa(r.oldScore), strconv.Itoa(r.newScore), r.status})
	}
	resTable.Render()
}

func testDataset(dataset string, model string, scorer string, res *result, datafolder string, submissionsfolder string, realtimeoutput bool) {
	// Compute file names
	inputFileName := fmt.Sprintf("%s/%s.in", datafolder, dataset)
	scoreFileName := fmt.Sprintf("%s/%s.score", submissionsfolder, dataset)
	outputFileName := fmt.Sprintf("%s/%s.out", submissionsfolder, dataset)
	tmpOutputFileName := fmt.Sprintf("%s-tmp/%s.out.tmp", submissionsfolder, dataset)

	// Read oldScore
	var oldScore int
	{
		oldscore, err := ioutil.ReadFile(scoreFileName)
		if err != nil {
			log.Println("Couldn't load ", scoreFileName, ":", err)
		} else {
			oldScore, err = strconv.Atoi(strings.TrimSpace(string(oldscore)))
			if err != nil {
				log.Println("Couldn't parse", scoreFileName, "as int: ", oldscore)
				return
			}
			res.oldScore = oldScore
		}
	}

	// Run the model
	var modelOutput []byte
	{
		modelCmd := exec.Command(model, inputFileName, tmpOutputFileName)
		if realtimeoutput {
			modelCmd.Stdout, modelCmd.Stderr = os.Stdout, os.Stderr
			err := modelCmd.Run()
			if err != nil {
				log.Println("Error executing model", dataset, ":", err)
				return
			}
			modelOutput = []byte("See above\n")
		} else {
			modelOutput, err := modelCmd.CombinedOutput()
			if err != nil {
				log.Println("Error executing model", dataset, ":", err)
				log.Println("model output for dataset", dataset, ":\n", string(modelOutput), "---")
				return
			}
		}
	}

	// Compute new score
	var newScore int
	var scorerStderr bytes.Buffer
	{
		scorerCmd := exec.Command(scorer, inputFileName, tmpOutputFileName)
		scorerCmd.Stderr = &scorerStderr
		scoreout, err := scorerCmd.Output()
		if err != nil {
			log.Println("Error computing new score:", err)
			return
		}

		newScore, err = strconv.Atoi(strings.TrimSpace(string(scoreout)))
		if err != nil {
			log.Println("Couldn't parse new score as int: ", string(scoreout))
			return
		}

		if newScore > oldScore {
			err := swapFiles(tmpOutputFileName, outputFileName)
			if err != nil {
				log.Println("Failed to swap output with the submission:", err)
				return
			}

			err = ioutil.WriteFile(scoreFileName, scoreout, 0777)
			if err != nil {
				log.Println("Error writing new score:", err)
				return
			}
			res.status = okgreen("better")
		} else if newScore == oldScore {
			res.status = warnorange("same")
		}
	}
	res.newScore = newScore

	// Output infos
	fmt.Printf("%s\n%s\n%s%s\n%s",
		okblue(makebold(fmt.Sprintf("Dataset %s finished with a score of %d.", dataset, newScore))),
		makebold("- model stdout & stderr:"),
		modelOutput,
		makebold("- scorer stderr:"),
		scorerStderr.Bytes())
}
