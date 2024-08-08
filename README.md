# micli

Micli is stand for Market Information based on cli, this is a simple cli to interact with the market information.
Currently, this cli compatible with darwin / mac terminal and only support the crypto market information 

> NOTE: This project is still in development, if you want to contribute, please feel free to fork and create a pull request

## How to install

### Prerequisite
- Go 1.20 or higher
- Make sure setup all go environment variable

### Install
- Clone this repository
- Run `./build.sh` to build the binary
- create directory `~/micli` in your home directory
- move the binary of `micli` to your `~/micli` directory
- move the binary to `~/micli` directory
- add `~/micli` to your PATH
- run `micli` in your terminal

## Feature
[WIP]

## Demo
[WIP]

## architecture code [if you want to contribute]

### model
model contains the data model and the logic to manipulate it.

### repo
repo contains the logic to interact with the external api or data

### service
service contains the logic to interact with the model and the repo

### view
view contains the logic to interact with the user
