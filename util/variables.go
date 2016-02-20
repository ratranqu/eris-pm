package util

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/eris-ltd/eris-pm/definitions"

	log "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

func PreProcess(toProcess string, do *definitions.Do) (string, error) {
	// $block.... $account.... etc. should be caught. hell$$o should not
	// :$libAddr needs to be caught
	catchEr := regexp.MustCompile("(^|\\s|:)\\$([a-zA-Z0-9_]+)")

	// If there's a match then run through the replacement process
	if catchEr.MatchString(toProcess) {
		// find what we need to catch.

		processedString := toProcess
		for _, jobMatch := range catchEr.FindAllStringSubmatch(toProcess, -1) {
			jobName := jobMatch[2]
			varName := "$" + jobName

			// first parse the reserved words.
			if strings.Contains(jobName, "block") {

				block, err := replaceBlockVariable(jobName, do)
				if err != nil {
					return "", err
				}
				strings.Replace(processedString, varName, block, 1)
			}

			// second we loop through the jobNames to do a result replace
			for _, job := range do.Package.Jobs {
				if string(jobName) == job.JobName {
					log.WithFields(log.Fields{
						"jobname": string(jobName),
						"result":  job.JobResult,
					}).Debug("Fixing Variables =>")
					processedString = strings.Replace(processedString, varName, job.JobResult, 1)
				}
			}
		}
		return processedString, nil
	}

	// if no matches, return original
	return toProcess, nil
}

func replaceBlockVariable(toReplace string, do *definitions.Do) (string, error) {
	block, err := ChainStatus("latest_block_height", do)
	if err != nil {
		return "", err
	}

	if toReplace == "block" {
		return block, nil
	}

	catchEr := regexp.MustCompile("block\\+(\\d*)")
	if catchEr.MatchString(toReplace) {
		height := catchEr.FindStringSubmatch(toReplace)[1]
		h1, err := strconv.Atoi(height)
		if err != nil {
			return "", err
		}
		h2, err := strconv.Atoi(block)
		if err != nil {
			return "", err
		}
		height = strconv.Itoa(h1 + h2)
		return height, nil
	}

	catchEr = regexp.MustCompile("block\\-(\\d*)")
	if catchEr.MatchString(toReplace) {
		height := catchEr.FindStringSubmatch(toReplace)[1]
		h1, err := strconv.Atoi(height)
		if err != nil {
			return "", err
		}
		h2, err := strconv.Atoi(block)
		if err != nil {
			return "", err
		}
		height = strconv.Itoa(h1 - h2)
		return height, nil
	}

	return toReplace, nil
}
