#!/usr/bin/env nu

let build_time = date now | format date "%Y-%m-%d %H:%M"
let version = (git describe --tag --long --dirty=-dirty)
print {"version": $version, "build_time": $build_time}

def build [] {
    let ldflags = $'-X "main.Version=($version)" -X "main.VersionTime=($build_time)"'
    timeit { go build -gcflags="all=-N -l" -ldflags $ldflags } | print ($in | format duration sec)
}

def build_linux [] {
    print 'build linux'
    GOOS=linux GOARCH=amd64 build
}

def build_windows [] {
    print 'build windows'
    GOOS=windows GOARCH=amd64 build
}

def main [] {
    build_windows; build_linux
}

def "main windows" [] {
    build_windows
}

def "main win" [] {
    build_windows
}

def "main linux" [] {
    build_windows
}
