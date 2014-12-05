multiredir is a small program to redirect a bunch of TCP ports to another
hosts/ports.

Usage:

	multiredir:
	  -r=: redirection to serve: bind_host:port/to_host:port (multiple -r params can be used)
	  -t=5s: connect timeout
	  -v=false: be verbose on errors

Example:

	multiredir -r :80/192.0.2.1:8080 -r :443/192.0.2.1:443 -r :11211/192.0.2.10:11211
