# jira2trello
[![CI](https://github.com/Brialius/jira2trello/actions/workflows/ci.yml/badge.svg)](https://github.com/Brialius/jira2trello/actions/workflows/ci.yml)[![Go Report Card](https://goreportcard.com/badge/github.com/Brialius/jira2trello)](https://goreportcard.com/report/github.com/Brialius/jira2trello)

## Usage
```
Usage:
     jira2trello [command]
   
   Available Commands:
     configure     Ask configuration settings and save them to file
     help          Help about any command
     report        Report based on trello cards or jira query
     sync          Jira to Trello sync
     update        Update jira2trello
     weekly-report Weekly report based on jira query
   
   Flags:
         --config string   config file (default is $HOME/.jira2trello.yaml)
         --debug           write debug info to logs and files
     -h, --help            help for jira2trello
   
   Use "jira2trello [command] --help" for more information about a command.   
```
## Screenshots
#### Sync command
![image](https://user-images.githubusercontent.com/6441812/143793782-159757dc-12fe-46c9-a502-1229f346f4d3.png)

## Build
### make goals
|Goal|Description|
|----|-----------|
|build (default)|build binaries|
|setup|download and install required dependencies|
|test|run tests|
|install|install binary to `$GOPATH/bin`|
|lint|run linters|
|clean|run `go clean`|
|generate|run `go generate ./...`|
|mod-refresh|run `go mod tidy` and `go mod vendor`|
|ci|run all steps needed for CI|
|version|show current git tag if any matched to `v*` exists|
|release|set git tag and push to repo `make release ver=v1.2.3`|
