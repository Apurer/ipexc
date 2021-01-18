package ipexc

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"testing"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hi there, I love you!")
}

func server() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var timeout = time.Duration(2 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

func TestFunctions(t *testing.T) {

	transport := http.Transport{
		Dial: dialTimeout,
	}

	client := http.Client{
		Transport: &transport,
	}

	// start http server
	go server()
	//make request to http server
	_, err := client.Get("http://localhost:8080")
	if err != nil {
		t.Errorf("error during making http request to server with message: %q", err)
	}

	// change iptables rules
	args := []string{"-I", "INPUT", "-p", "tcp", "--dport", "8080", "-j", "DROP"}

	cmd := exec.Command("iptables", args...)

	err = cmd.Run()
	if err != nil {
		t.Errorf("error during preparing iptables rules for test: %q", err)
	}

	// run Insert function
	var port uint64 = 8080
	var ip uint32 = 2130706433
	err = Insert(port, ip)
	if err != nil {
		t.Errorf("error during running Insert function : %q", err)
	}
	// make request to http server
	_, err = client.Get("http://localhost:8080")
	if err != nil {
		t.Errorf("error during making http request to server after running Insert function with message: %q", err)
	}

	// run Delete function
	err = Delete(port, ip)
	if err != nil {
		t.Errorf("error during running Delete function : %q", err)
	}

	// clean up of iptables rules
	args[0] = "-D"

	cmd = exec.Command("iptables", args...)

	err = cmd.Run()
	if err != nil {
		t.Errorf("error during removing rules from iptables: %q", err)
	}
}
