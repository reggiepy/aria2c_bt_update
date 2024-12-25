# golang auto update aria2c bt-tracker task

[![go version](https://img.shields.io/github/go-mod/go-version/reggiepy/aria2c_bt_update?color=success&filename=go.mod&style=flat)](https://github.com/reggiepy/aria2c_bt_update)
[![release](https://img.shields.io/github/v/tag/reggiepy/aria2c_bt_update?color=success&label=release)](https://github.com/reggiepy/aria2c_bt_update)
[![build status](https://img.shields.io/badge/build-pass-success.svg?style=flat)](https://github.com/reggiepy/aria2c_bt_update)
[![License](https://img.shields.io/badge/license-GNU%203.0-success.svg?style=flat)](https://github.com/reggiepy/aria2c_bt_update)
[![Go Report Card](https://goreportcard.com/badge/github.com/reggiepy/aria2c_bt_update)](https://goreportcard.com/report/github.com/reggiepy/aria2c_bt_update)

## Installation

```bash
git clone https://github.com/reggiepy/aria2c_bt_update.git
cd aria2c_bt_update
go mod tidy
```

## Usage

```bash
go run cmd/aria2c_bt_update/main.go
```

build
```bash
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
# go env -w CGO_ENABLED=0 GOOS=windows  GOARCH=amd64
go build github.com/reggiepy/aria2c_bt_updater/cmd/aria2c_bt_update
go build -ldflags="-s -w" github.com/reggiepy/aria2c_bt_updater/cmd/aria2c_bt_update
```

## Architecture
