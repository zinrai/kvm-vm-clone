package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	sourceVM         string
	destVM           string
	noHostnameChange bool
)

func init() {
	flag.StringVar(&sourceVM, "source", "", "Name of the source VM to clone")
	flag.StringVar(&destVM, "dest", "", "Name for the new cloned VM")
	flag.BoolVar(&noHostnameChange, "no-hostname-change", false, "Skip hostname change (for FreeBSD or other unsupported OSes)")
}

func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	fmt.Printf("Executing command: %s %s\n", name, strings.Join(args, " "))

	output, err := cmd.CombinedOutput()
	fmt.Printf("Command output:\n%s\n", string(output))

	return string(output), err
}

func checkCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func checkRequiredCommands() error {
	requiredCommands := []string{"sudo", "virt-clone", "virsh"}
	if !noHostnameChange {
		requiredCommands = append(requiredCommands, "virt-customize")
	}
	missingCommands := []string{}

	for _, cmd := range requiredCommands {
		if !checkCommand(cmd) {
			missingCommands = append(missingCommands, cmd)
		}
	}

	if len(missingCommands) > 0 {
		return fmt.Errorf("Missing required commands: %v", missingCommands)
	}

	return nil
}

func isVMRunning(vmName string) (bool, error) {
	output, err := runCommand("sudo", "virsh", "list", "--name", "--state-running")
	if err != nil {
		return false, fmt.Errorf("Failed to get list of running VMs: %v", err)
	}

	runningVMs := strings.Split(strings.TrimSpace(output), "\n")
	for _, vm := range runningVMs {
		if vm == vmName {
			return true, nil
		}
	}

	return false, nil
}

func cloneVM(sourceVM, newVM string) error {
	_, err := runCommand("sudo", "virt-clone", "--original", sourceVM, "--name", newVM, "--auto-clone")
	return err
}

func setVMHostname(vmName, hostname string) error {
	_, err := runCommand("sudo", "virt-customize", "-d", vmName, "--hostname", hostname)
	return err
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s --source source_vm --dest new_vm [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nClones a KVM virtual machine using virt-clone and sets the new hostname.\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if sourceVM == "" || destVM == "" {
		fmt.Println("Error: Both --source and --dest flags are required")
		flag.Usage()
		os.Exit(1)
	}

	if err := checkRequiredCommands(); err != nil {
		log.Fatalf("Error: %v\nPlease install the missing commands and try again.", err)
	}

	running, err := isVMRunning(sourceVM)
	if err != nil {
		log.Fatalf("Error checking source VM state: %v", err)
	}
	if running {
		log.Fatalf("Error: Source VM '%s' is currently running. Please stop the VM before cloning.", sourceVM)
	}

	fmt.Printf("Cloning VM '%s' to '%s'...\n", sourceVM, destVM)
	if err := cloneVM(sourceVM, destVM); err != nil {
		log.Fatalf("Failed to clone VM: %v", err)
	}

	if !noHostnameChange {
		fmt.Printf("Setting hostname of new VM to '%s'...\n", destVM)
		if err := setVMHostname(destVM, destVM); err != nil {
			log.Fatalf("Failed to set hostname: %v", err)
		}
		fmt.Printf("VM '%s' cloned successfully to '%s' and hostname updated\n", sourceVM, destVM)
	} else {
		fmt.Printf("VM '%s' cloned successfully to '%s' (hostname change skipped)\n", sourceVM, destVM)
	}
}
