# 2025gopherconuk
Code from my 2025 GopherCon UK introductory talk on Unix Programming in Go

## Running Examples

In general the example programs are run from the commandline in the form

	go run "file name" helpers.go [parameters, ...]

However there are some exceptions.

Where signal handling is used as the "go" command effectively interposes a subprocess between a parent go program calling it and the child process it launches. In this case it's necessary to build a binary for the child executable, leading to

	go build 04_child.go helpers.go
	go run 04_parent.go helpers.go

For examples 10_create and 10_use which use CGO it's necessary to build them separately with the commands

	cd 10_create
	go build
	./10_create

	cd 10_use
	go build
	./10_use

For examples 11_create and 11_use which use Syscalls it's necessary to add darwin.go (on macOS)

	go run 11_create.go helpers.go darwin.go
	go run 11_use.go helpers.go darwin.go


darwin.go should be replaced with an alternative OS-specific file on other platforms.


## Talk Abstract

Are you writing Go programs that run on Linux or MacOS? Then you're writing Go programs for a UNIX environment! But how much does that figure in your daily experience?

This talk provides an introduction to programming for the UNIX environment using Go.

We'll start with the premise of an operating system kernel, the program which manages a UNIX system’s computing resources. This will lead to discussion of memory management, file handles and pipes, before building into more advanced topics around process coordination and privileges, interprocess communication, networking, and concurrency.

On a UNIX system many facilities we think of as part of the operating system are provided by third-party libraries written in languages such as C or Rust. Whilst this is a deep topic with many quirks, we’ll finish up with a gentle introduction to CGO, the Go tool for incorporating C code into our programs. With CGO we’re able to access a wealth of pre-existing functionality, whether that be direct access to LIBC’s malloc() function or creating a wrapper for a database engine like SQLite.

By the end of this talk you should have a broad overview of the UNIX environment from a programmer’s perspective, which of its facilities are directly accessible from Go’s standard library, and how to access those which aren’t from CGO.

The focus throughout will be on code examples, building from simple beginnings to more advanced use cases, to provide a springboard for further exploration.


## What is UNIX

* kernel
** micro
** monolith
* device drivers
* boot loader
* init process
* Managing computing resources and programs
* Controlling user access to resources

## Memory management

* executable vs data
* allocating pages
* virtual memory
* dynamic loading

## block devices, directory trees, and file handles

* character vs binary
* file descriptors
* file i/o

## process coordination

* locks
* semaphores
* shared memory

## pipes, sockets, and IPC

* stdin, stdout, stderr
* command line parameters
* domain sockets
* tcp/ip sockets

## privileges and resource control

* users
* groups
* file permissions

## threads

## signals & timers

## CGO

* code example of using malloc
* build tags for different platforms