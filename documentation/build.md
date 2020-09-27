# Project build
In this documentation you can find the information about project build

## Build for current system
In these instructions we assume, that you need to build this project for your current system. For build of the project you need to follow the next steps:

``Watning! The following steps will work for MacOs and Linux systems``
1. Clone the latest version of the project to your machine
``` 
git clone git@github.com:sharovik/wt.git ~/go/src/github.com/sharovik/wt
```
2. Go to the project dir and run next command:
``` 
make build
```
3. If there is no errors, you will see the next binary files
-- `./bin/wt` - the application which is ready for run.

## Cross platform build
If you want to run cross-platform build, please use the following instructions

### Before cross-platform build
For cross-platform build I use `karalabe/xgo-latest`. So please before project build do the following steps
1. Install `docker` and `go` to your system
2. Run this command `docker pull karalabe/xgo-latest`
3. Your project should be in `GOPATH` folder or `GOPATH` should point to the directory where you clone this project

### Build
For build please run this command
``` 
make build-project-cross-platform
```

### Supported OS

This command will build the `wt` for the following types of the system:
#### MacOS
- darwin-386
- darwin-amd64
#### Linux
- linux-386
- linux-amd64
#### Windows
- windows-386
- windows-amd64
