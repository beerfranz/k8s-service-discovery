package main

import (
	"strconv"
    "io"
    "bytes"
)

func convertToTraefikYaml(backends []Backend) *bytes.Buffer {
	buffer := bytes.NewBufferString("")

	io.WriteString(buffer, "http:\n  services:\n")
	for _, backend := range backends {
		io.WriteString(buffer, "    " + backend.Name + ":\n        loadBalancer:\n          servers:\n")

		for _, target := range backend.Targets {
			port := strconv.Itoa(int(target.Port))
			io.WriteString(buffer, "            - url: \"" + target.Protocol + "://" + target.Ip + ":" + port + "\"\n")
		}
	}

	return buffer
}
