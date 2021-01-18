package ipexc

import (
	"encoding/binary"
	"net"
	"os/exec"
	"strconv"
)

func ipToint(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func intToip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}

// Insert function for inserting new records into tables INPUT and OUTPUT
func Insert(port uint64, ip uint32) error {

	args := []string{"-I", "INPUT", "-p", "tcp", "-s", intToip(ip).String(), "--dport", strconv.FormatUint(port, 10), "-m", "state", "--state", "NEW,ESTABLISHED", "-j", "ACCEPT"}

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
func Delete(port uint64, ip uint32) error {

	args := []string{"-D", "INPUT", "-p", "tcp", "-s", intToip(ip).String(), "--dport", strconv.FormatUint(port, 10), "-m", "state", "--state", "NEW,ESTABLISHED", "-j", "ACCEPT"}

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
