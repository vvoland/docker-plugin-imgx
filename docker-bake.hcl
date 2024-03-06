target "bin-image" {
    platforms = [
        "linux/amd64",
        "linux/arm64",
        "darwin/arm64",
        "darwin/amd64",
        "windows/amd64",
        "windows/arm64"
    ]
    target = "binary"
    tags = ["pawelgronowski465/docker-plugin-imgx:latest"]
}

target "binary" {
    target = "binary"
    platforms = ["local"]
    output = ["./bin"]
}

group "default" {
    targets = ["binary"]
}

