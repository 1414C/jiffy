---
title: "Installing Jiffy"
date: 2018-02-05T20:39:34-07:00
weight: 10
draft: true
---

The easiest way to install Jiffy is to use *go get* on the command-line to pull the latest version from github, as shown below:

```bash

    $ go get -u github.com/1414C/jiffy

```

This will pull the jiffy github repository into your $GOPATH/src/github.com folder, as well as any dependencies referenced by the jiffy source code.  The *-u* flag is included to instruct *go get* to check for and pull updates to jiffy packages and their dependencies.  This is the least sophisticated way of managing dependencies in go; check the status of [dep](https://github.com/golang/dep) if you are interested in more advanced dependency management.

Once the jiffy source code and dependencies have been installed into your $GOPATH, you can use *go build* to compile a binary from the jiffy sources.  The easiest way to do this is to open a terminal window, switch to $GOPATH/src/github.com/1414C/jiffy and run *go build* as shown below.

```bash

    $ go build -v

```

This will result in the creation of a binary file called 'jiffy'.  You may move the 'jiffy' binary anywhere in your $PATH, but convention would have you install it in /usr/local/bin.  Once you have moved the 'jiffy' binary to its new home, open a new terminal window and use the *which* command to ensure that you can access the binary.

```bash

    $ which jiffy -a

```

If *which* cannot find the jiffy binary (or finds the wrong one!), you need to make sure that your $PATH is set correctly.  
