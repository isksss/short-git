package main

import (
	"log"
	"os/exec"
	"strings"
)

func main() {
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
}

func executeCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	_, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
}
