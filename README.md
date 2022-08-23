# epptester

A test tool for EPP connections

Most EPP systems don't give decent error messages and it's often extremely hard to figure out if the issue is at the
registry or at the client side.
This tool is intended to help fix that. There are versions here for Linux, Mac, and Windows. If you need something
else (Arm/Raspberry PI???) let me know and I'll start including it.

Example Usage:
./bin/epptester -cert mycert.pem -key mycert.pem --tls 1.2 -port 3121 -host localhost -username someuser
-password 1234567890

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

# Epppromtester

An EPP tester intended to be used with [Node Exporter](https://github.com/prometheus/node_exporter) 

## Usage
``` epppromtester -C configfile.yaml ```

See eppprom.yaml.sample for an example of  fields and zones. Multiple epp servers can be tested from one config file as long as 
they all share the same client certificates
The tester will evaluate both IPv4 and IPv6 connections if both A and AAAA records have been set for the hostname.

## Fields in config file
* __promfile__  Text collector format file to write the results to.
* __refresh__ How often to run a check if being run as a daemon under systemd. Defaults to 0 a single run and exit.
* __cert__  Client Certificate
* __key__  Client Private key

### Zone list.
Under the __zones__ key
* __name__ Name of the epp zone or server for the prometheus file.
* __host__ Hostname to connect to.
* __port__ Port to connect to. [700]
* __ipversion__ One of ipv4, ipv6, or both. Which ip stack to test  [both]
* __tlsversion__ TLS version to test with. Currently one of 1.0, 1.1, 1.2, 1.3, any [any]


