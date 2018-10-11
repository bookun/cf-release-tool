[![Build Status](https://travis-ci.org/bookun/cf-release-tool.svg?branch=master)](https://travis-ci.org/bookun/cf-release-tool) [![Maintainability](https://api.codeclimate.com/v1/badges/33d4eb3e51099945979d/maintainability)](https://codeclimate.com/github/bookun/cf-release-tool/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/33d4eb3e51099945979d/test_coverage)](https://codeclimate.com/github/bookun/cf-release-tool/test_coverage)  
# Blue-Green deployment tool for PHP application in cloud foundry

## Overview
cf-release-tool is a plugin for the CF command line tool that executes Blue-Green deployment.
Blue-Green deployment is zero-downtime deploys.

## Features
* provides Blue-Green deployment
* Pushs app based on git branch that you want to release
* <WIP> map test route to *green app*. if user approves, map production route to it.

## How to Build

1. `go get github.com/bookun/cf-release-tool`
2. `cd $GOPATH/src/github.com/bookun/cf-release-tool`
3. `go build -o ReleaseTool`
4. `cf install-plugin ReleaseTool`
5. `cf release -h`

## How to Use

*Caution*
Please append env variables to your manifest file, `ORG`, `SPACE`, `HOST`, `DOMAIN`. [sample](https://github.com/bookun/cf-release-tool/blob/v1.0/testdata/manifest1.yml)

* pushs app based on master branch and manifest.yml
    `cf release`

* pushs app based on "branch" and manifest.yml
    `cf release -b <branch>`

* pushs app based on "branch" and "your/manifest-file/path"
    `cf release -b <branch> -f <your/manifest-file/path>`

## Uninstall
`cf uninstall-plugin ReleaseTool`
