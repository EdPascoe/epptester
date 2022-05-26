# epptester

A test tool for EPP connections

Most EPP systems don't give decent error messages and it's often extremely hard to figure out if the issue is at the
registry or at the client side.
This tool is intended to help fix that. There are versions here for Linux, Mac, and Windows. If you need something
else (Arm/Raspberry PI???) let me know and I'll start including it.

Example Usage:
./bin/epptester -cert dnservices.pem -key dnservices.pem --tls 1.2 -port 3121 -host localhost -username dnservices
-password e4a192c3

Options:
* -host Epp server  (default "127.0.0.1")
* -port Server port  (default 3121)
* -cert certificate PEM file (default "cert.pem")
* -key string key PEM file (default "key.pem")
* -username EPP Username
* -password EPP Password
* -tls string TLS version to test with: 1.1, 1.2, 1.3, any (default "any")
* -version Show version information

If the certificate and key can both be in the same file.



