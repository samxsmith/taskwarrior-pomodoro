package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	intervalInMins, err := strconv.Atoi(textLinePrompt("For how many minutes?"))
	if err != nil {
		fmt.Println("error with interval: ", err)
		return
	}

	taskWarriorID := textLinePrompt("Which taskwarrior task ID?")
	if taskWarriorID == "" {
		fmt.Println("No task ID passed")
		return
	}

	b, err := exec.Command("task", "_get", fmt.Sprintf("%s.description", taskWarriorID)).Output()
	if err != nil {
		fmt.Println("task warrior error: ", err)
		return
	}

	ok, err := confirmPrompt(fmt.Sprintf("You want to run task: %s", string(b)))
	if !ok || err != nil {
		fmt.Println("Won't run: ", err)
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go handleInterrupt(sigs, taskWarriorID)

	fmt.Println("Will start task: ", string(b))

	err = exec.Command("task", "start", taskWarriorID).Run()
	if err != nil {
		fmt.Println("error starting task: ", err)
		done(taskWarriorID)
		return
	}

	fmt.Printf("Will run for %d minutes... \n", intervalInMins)

	waitDuration := time.Duration(intervalInMins) * time.Minute
	fmt.Println("Will finish at: ", time.Now().Add(waitDuration).Format(time.Kitchen))
	time.Sleep(waitDuration)

	fmt.Println("time's up!!")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go playSound(&wg)
	done(taskWarriorID)
	wg.Wait()
}

func playSound(wg *sync.WaitGroup) {
	for i := 0; i < 5; i++ {
		exec.Command("afplay", "/System/Library/Sounds/Blow.aiff").Run()
	}
	wg.Done()
}

func handleInterrupt(sigs chan os.Signal, taskID string) {
	sig := <-sigs
	fmt.Println()
	fmt.Println("got interrupt: ", sig)
	fmt.Println("beginning graceful shut down")
	done(taskID)
	os.Exit(1)
}

func done(taskID string) {
	fmt.Println("Stopping task")
	exec.Command("task", taskID, "stop").Run()
}

func textLinePrompt(msg string) string {
	print(fmt.Sprintf("%v: ", msg))
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	return input
}

func confirmPrompt(msg string) (bool, error) {
	fmt.Println(msg)
	var input string
	print("y/N: ")
	_, e := fmt.Scanf("%s", &input)
	if e != nil {
		return false, e
	}
	normalisedInput := strings.ToLower(strings.TrimSpace(input))
	if normalisedInput == "y" || normalisedInput == "yes" {
		return true, e
	}
	return false, e
}
