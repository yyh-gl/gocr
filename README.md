# GoCR

GoCR is a notifier of code review requests.

You can easily set the repository and notification destination.

## Supported tools

- Code Management services
  - GitHub
- Message Services
  - Slack 

## Introduction

1. `$ go get -u github.com/yyh-gl/gocr`
1. Set repositories and notification destinations into `.gocr.yml`.  
Example `.gocr.yml` is [here](https://github.com/yyh-gl/gocr/blob/master/.gocr.example.yml).
1. `$ gocr -c /path/to/.goct.yml`
