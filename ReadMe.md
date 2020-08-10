# Echelon - hierarchical progress in terminal on steroids

[![Build Status](https://api.cirrus-ci.com/github/cirruslabs/echelon.svg)](https://cirrus-ci.com/github/cirruslabs/echelon)

Library to show tree-like hierarchical progress in VT100 compatible terminals.

Here is an example how it looks for running Dockerized tasks via [Cirrus CLI](https://github/cirruslabs/cirrus-cli):

![Cirrus CLI Demo](images/cirrus-cli-demo.gif)

## Features

* Customizable and works with any VT100 compatible terminal
* Implements incremental drawing algorithm to optimize drawing performance
* Can be used from multiple go routines

## Example

Please check `demo` folder for a simple example or how *echelon* is used in [Cirrus CLI](https://github/cirruslabs/cirrus-cli):
