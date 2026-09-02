package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"strings"

	"github.com/ofgrenudo/415/pkg/aws"
	"github.com/ofgrenudo/415/pkg/net"
)

type domainList []string

func (d *domainList) String() string {
	return strings.Join(*d, ",")
}

func (d *domainList) Set(v string) error {
	*d = append(*d, v)
	return nil
}

func main() {
	var domains domainList
	flag.Var(&domains, "domain", "Domain to update in Route 53 (repeatable)")
	flag.Var(&domains, "d", "Shorthand for --domain (repeatable)")
	flag.Parse()

	if len(domains) == 0 {
		slog.Error("At least one --domain/-d flag is required.")
		os.Exit(1)
	}

	slog.Info("Starting Dynamic DNS....")

	ip := net.GetIPv4()
	if ip == nil {
		slog.Error("Failed to determine public IP. Aborting.")
		return
	}
	if ip.IsPrivate() {
		slog.Error("IP Returned is a private IP. You cannot assign a Private IP to a Public DNS record.", "IP", ip)
		return
	}
	slog.Info("Successfully received public IP.", "IP", ip.String())

	ctx := context.Background()
	client, err := aws.NewClient(ctx)
	if err != nil {
		slog.Error("Failed to create AWS Route 53 client.", "Error Message", err)
		os.Exit(1)
	}

	failed := false
	for _, domain := range domains {
		zoneID, err := aws.FindHostedZoneID(ctx, client, domain)
		if err != nil {
			slog.Error("Failed to find hosted zone.", "Domain", domain, "Error Message", err)
			failed = true
			continue
		}

		currentIP, err := aws.GetCurrentARecordIP(ctx, client, zoneID, domain)
		if err != nil {
			slog.Error("Failed to read current A record.", "Domain", domain, "Error Message", err)
			failed = true
			continue
		}

		if currentIP == ip.String() {
			slog.Info("A record already up to date.", "Domain", domain, "IP", currentIP)
			continue
		}

		if err := aws.UpsertARecord(ctx, client, zoneID, domain, ip.String()); err != nil {
			slog.Error("Failed to update A record.", "Domain", domain, "Error Message", err)
			failed = true
			continue
		}

		slog.Info("Successfully updated A record.", "Domain", domain, "IP", ip.String())
	}

	if failed {
		os.Exit(1)
	}
}
