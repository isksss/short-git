package main

import (
	"log"
	"os/exec"
	"strings"
)

func main() {
	// git configからユーザ名を取得します。
	userName, err := executeCommand("git", "config", "--get", "user.name")
	if err != nil {
		log.Fatal(err)
	}

	// 改行を削除し、ブランチ名を作成します。
	branchName := strings.TrimSpace(userName) + "_branch"

	// ブランチが存在するか確認します。
	_, err = executeCommand("git", "rev-parse", "--verify", branchName)
	if err != nil {
		// ブランチが存在しない場合、新しいブランチを作成します。
		_, err = executeCommand("git", "checkout", "-b", branchName)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Created new branch: %s\n", branchName)

		// 新しいブランチをリモートにプッシュします。
		_, err = executeCommand("git", "push", "-u", "origin", branchName)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Pushed new branch: %s to remote\n", branchName)
	}

	// git statusコマンドを実行します。
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	// 変更がある場合、それをコミットしてプッシュします。
	if len(strings.TrimSpace(string(out))) > 0 {
		executeCommand("git", "add", "-A")
		executeCommand("git", "commit", "-m", "auto commit")
		executeCommand("git", "push", "--all")
	} else {
		log.Println("No changes to commit")
	}

	// 全てのブランチでgit pullを実行します。
	pullAllBranches()
}

func pullAllBranches() {
	// 現在のブランチを取得します。
	currentBranch, err := executeCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		log.Fatal(err)
	}

	// すべてのブランチ名を取得します。
	branchList, err := executeCommand("git", "branch", "--all")
	if err != nil {
		log.Fatal(err)
	}

	branches := strings.Split(strings.TrimSpace(branchList), "\n")

	// 各ブランチでgit fetchを実行します。
	for _, branch := range branches {
		branch = strings.TrimPrefix(branch, "* ")
		branch = strings.TrimSpace(branch)

		_, err = executeCommand("git", "checkout", branch)
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = executeCommand("git", "fetch")
		if err != nil {
			log.Println(err)
		}
	}

	// 元のブランチに戻します。
	_, err = executeCommand("git", "checkout", strings.TrimSpace(currentBranch))
	if err != nil {
		log.Fatal(err)
	}
}

func executeCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	out, err := cmd.Output()
	return string(out), err
}
