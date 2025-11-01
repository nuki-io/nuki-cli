# Nuki CLI

A versatile command line interface to manage Nuki devices through bluetooth and online services.

## ⚠️ Disclaimer ⚠️

This tool is still in a very early development stage. Use at your own risk.
It is maintained on a best-effort basis by Nuki employees on a volunteer basis. No claims as to its completeness or correctness are made.

This project is not an officially supported product of Nuki.

## Requirements

- Go 1.23
- Stringer (`go install golang.org/x/tools/cmd/stringer@latest`)

## How to use

### Building

```
make build
```

Will produce a binary `nukictl` in the repo root.

## Package structure

### nukible

This contains the low level handling of BLE with the library tinygo.org/x/bluetooth. It tries to abstract away some peculiarities of different OS.

### blecommands

Contains the command structs that are used to communicate with smartlocks. Every command has its own struct with adequate Go datatypes, like mapping datetime structures or adding enums.

If you just want to communicate with a device and want to implement all flow and error handling by yourself, this is the package you need.

### bleflows

High-level abstraction that aims to provide use case specific functions, while dealing with all protocol stuff and flows under the hood for you.
