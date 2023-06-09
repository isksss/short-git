package main

import (
	"log"
	"os/exec"
	"strings"
)

func main() {
	userName := getGitUserName()
	branchName := createBranchName(userName)

	if !branchExists(branchName) {
		createAndPushNewBranch(branchName)
	} else {
		checkoutBranch(branchName)
	}

	commitAndPushChanges()

	pullAllBranches()
}

func getGitUserName() string {
	userName, err := executeCommand("git", "config", "--get", "user.name")
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(userName)
}

func createBranchName(userName string) string {
	return userName + "_branch"
}

func branchExists(branchName string) bool {
	_, err := executeCommand("git", "rev-parse", "--verify", branchName)
	return err == nil
}

func createAndPushNewBranch(branchName string) {
	_, err := executeCommand("git", "checkout", "-b", branchName)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Created new branch: %s\n", branchName)

	_, err = executeCommand("git", "push", "-u", "origin", branchName)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Pushed new branch: %s to remote\n", branchName)
}

func checkoutBranch(branchName string) {
	currentBranch, err := executeCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		log.Fatal(err)
	}

	if strings.TrimSpace(currentBranch) != branchName {
		_, err = executeCommand("git", "checkout", branchName)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Printf("Checked out to branch: %s\n", branchName)
		}
	}
}

func commitAndPushChanges() {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	if len(strings.TrimSpace(string(out))) > 0 {
		executeCommand("git", "add", "-A")
		executeCommand("git", "commit", "-m", "auto commit")
		executeCommand("git", "push", "--all")
	} else {
		log.Println("No changes to commit")
	}
}

// ... The rest of the code

func pullAllBranches() {
	currentBranch := getCurrentBranch()

	branches := getAllBranches()

	fetchUpdatesForBranches(branches)

	returnToOriginalBranch(currentBranch)
}

func getCurrentBranch() string {
	currentBranch, err := executeCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(currentBranch)
}

func getAllBranches() []string {
	branchList, err := executeCommand("git", "branch", "--all")
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(strings.TrimSpace(branchList), "\n")
}

func fetchUpdatesForBranches(branches []string) {
	for _, branch := range branches {
		branch = strings.TrimPrefix(branch, "* ")
		branch = strings.TrimSpace(branch)

		log.Printf("Checking out branch: %s\n", branch)

		_, err := executeCommand("git", "checkout", branch)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("Fetching updates for branch: %s\n", branch)

		_, err = executeCommand("git", "fetch")
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("Successfully fetched updates for branch: %s\n", branch)
	}
}

func returnToOriginalBranch(currentBranch string) {
	log.Printf("Returning to original branch: %s\n", currentBranch)

	_, err := executeCommand("git", "checkout", currentBranch)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Successfully returned to original branch: %s\n", currentBranch)
	}
}

func executeCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	out, err := cmd.CombinedOutput() // Standard output and standard error are combined
	output := string(out)

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			switch exitCode {
			case 128:
				log.Printf("Git command failed with exit code 128, not a git repository. Output was:\n%s", output)
			default:
				log.Printf("Git command failed with exit code %d. Output was:\n%s", exitCode, output)
			}
		} else {
			log.Printf("Failed to execute command: %v", err)
		}
	}

	return output, err
}
