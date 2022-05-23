package epptester

import (
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {
	name := "Gladys"
	if name != "Gladys" {
		t.Fatalf("Wrong name")
	}
}

// Checking TestFindserverip returns something the correct result
func TestFindserverip(t *testing.T) {
	ip := Findserverip("dns.google")
	// fmt.Println("IP: %s", ip)
	if ip != "8.8.8.8" && ip != "8.8.4.4" {
		t.Fatalf("Findserverip(dns.google) returned '%s'", ip)
	}
}

// Check that Findmyip doesn't crash
func TestFindmyip(t *testing.T) {
	Findmyip()
}
