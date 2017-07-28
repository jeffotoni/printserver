# printserver

A simple api rest to receive an encrypted POST and responsible for printing a label on the Zebra printer in the linux environment.

We will use the lpr native linux command to do the printing.

The way lpr works, in a nutshell, is: it reads in the file and hands the printable data over to the linux printing daemon, lpd. Lpd is a legacy piece of software for Unix and Linux, but it is supported under the modern system used by most Linux distributions, CUPS (the Common Unix Printing System).

You may need to manually install CUPS, and lpr itself, to print this way. If you are operating Debian, or a Debian-derived Linux system like Ubuntu that uses the APT package managements system, you can install them by running the following command:

Sudo apt-get update & & sudo apt-get install cups-client lpr

This command will install the Common Unix Printing System on your system. You should now be able to set up CUPS by directing any web browser to the address: http: // localhost: 631

The good thing is that we will send everything encrypted, but we can choose to encrypt the content before sending or not.


# Packages

go get "github.com/didip/tollbooth"

# Install

$ go build printserver.go

$ sudo cp printserver /usr/bin
