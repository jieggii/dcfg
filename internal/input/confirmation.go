package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ConfirmationPrompt(prompt string) (bool, error) {
	fmt.Printf("%v [Y/n] ", prompt)
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimRight(choice, "\n")
	if choice == "" || choice == "y" || choice == "Y" {
		return true, nil
	} else {
		return false, nil
	}
}
