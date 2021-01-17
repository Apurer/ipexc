package ipexc

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"testing"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hi there, I love you!")
}

func server() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func TestFunctions(t *testing.T) {
	// start http server
	go server()
	// make request to http server
	_, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Errorf("error during making http request to server with message: %q", err)
	}

	// change iptables rules
	args := []string{"-P", "INPUT", "-p", "tcp", "--dport", "8080", "-j", "DROP"}

	cmd := exec.Command("iptables", args...)

	err = cmd.Run()
	if err != nil {
		t.Errorf("error during preparing iptables rules for test: %q", err)
	}

	// make request to http server
	_, err = http.Get("http://localhost:8080")
	if err == nil {
		t.Errorf("iptables rules did not work: %v", args)
	}

	// run Insert function
	var port uint64 = 8080
	var ip uint32 = 2130706433
	err = Insert(port, ip)
	if err != nil {
		t.Errorf("error during running Insert function : %q", err)
	}
	// make request to http server
	_, err = http.Get("http://localhost:8080")
	if err != nil {
		t.Errorf("error during making http request to server after running Insert function with message: %q", err)
	}

	// run Delete function
	err = Delete(port, ip)
	if err != nil {
		t.Errorf("error during running Delete function : %q", err)
	}

	// make request to http server
	_, err = http.Get("http://localhost:8080")
	if err == nil {
		t.Error("running function Delete did not work")
	}

	// clean up of ip tables rules
	args[0] = "-D"

	cmd = exec.Command("iptables", args...)

	err = cmd.Run()
	if err != nil {
		t.Errorf("error during removing rules from iptables: %q", err)
	}
}
