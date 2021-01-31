package ipexc

import (
	"os/exec"
)

// Insert function for inserting new records into tables INPUT and OUTPUT
func Insert(port string, ip string) error {

	args := []string{"-I", "INPUT", "-p", "tcp", "-s", ip, "--dport", port, "-m", "state", "--state", "NEW,ESTABLISHED", "-j", "ACCEPT"}

	cmd := exec.Command("iptables", args...)

	err := cmd.Run()
	if err != nil {
		return err
	}

	args[1] = "OUTPUT"
	args[11] = "ESTABLISHED"

	cmd = exec.Command("iptables", args...)

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Delete function for deleting new records from tables INPUT and OUTPUT
func Delete(port string, ip string) error {

	args := []string{"-D", "INPUT", "-p", "tcp", "-s", ip, "--dport", port, "-m", "state", "--state", "NEW,ESTABLISHED", "-j", "ACCEPT"}

	cmd := exec.Command("iptables", args...)

	err := cmd.Run()
	if err != nil {
		return err
	}

	args[1] = "OUTPUT"
	args[11] = "ESTABLISHED"

	cmd = exec.Command("iptables", args...)

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
