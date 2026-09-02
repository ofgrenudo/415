package net

import (
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
)

func GetIPv4() net.IP {
	slog.Debug("Getting WAN IP....")
	resp, err := http.Get("https://ipv4.icanhazip.com")
	if err != nil {
		slog.Error("There was an unexpected error getting the IP from https://ipv4.icanhazip.com", "Error Message", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("There was an unexpected error reading the response from https://ipv4.icanhazip.com", "Error Message", err)
		return nil
	}
	if resp.StatusCode != 200 {
		slog.Error("Unexpected status code from https://ipv4.icanhazip.com", "StatusCode", resp.StatusCode, "Body", string(body))
		return nil
	}

	ip := net.ParseIP(strings.TrimSpace(string(body)))
	if ip == nil {
		slog.Error("Could not parse IP from response", "Body", string(body))
		return nil
	}

	slog.Debug("successfully received IP", "IP", ip.String())
	return ip
}
