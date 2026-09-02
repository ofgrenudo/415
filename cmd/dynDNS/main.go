package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/ofgrenudo/415/pkg/aws"
	"github.com/ofgrenudo/415/pkg/net"
)

const usageText = `NAME
       dynDNS -- dynamic DNS updater for AWS Route 53

SYNOPSIS
       dynDNS -d domain [-d domain ...]
       dynDNS --domain domain [--domain domain ...]
       dynDNS -h | --help

DESCRIPTION
       dynDNS determines the caller's current public IPv4 address by
       querying https://ipv4.icanhazip.com, then updates (UPSERT) the
       A record of each domain given via -d/--domain in AWS Route 53
       to point at that address.

       For each domain, the hosted zone that owns it is discovered
       automatically by walking up the domain's labels (for example
       "host.example.com" is checked against the zones "example.com"
       and "com" until a matching hosted zone is found). The existing
       A record, if any, is read first; the Route 53 API is only
       called to change the record when the published address differs
       from the address just retrieved, so unnecessary writes are
       avoided.

       This program is intended to be run periodically, for example
       from cron or a systemd timer, on a host whose public IP address
       may change over time (such as a residential internet
       connection), keeping one or more DNS records pointed at that
       host.

OPTIONS
       -d domain, --domain domain
              Domain name whose A record should be kept in sync with
              the caller's public IP. May be given more than once to
              update several domains in a single run. At least one is
              required.

       -h, --help
              Print this help message and exit.

EXIT STATUS
       0      All requested domains were checked successfully; each
              was either already up to date or was updated.

       1      No -d/--domain flag was given, the public IP could not
              be determined, an AWS client could not be created, or
              at least one domain failed to update (its hosted zone
              could not be found, its current record could not be
              read, or the Route 53 update itself failed). Other
              domains in the same invocation are still attempted.

ENVIRONMENT
       dynDNS authenticates to AWS using the standard AWS SDK for Go
       v2 credential chain. In order of precedence this includes:

              AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY,
              AWS_SESSION_TOKEN
                     Static credentials supplied directly.

              AWS_PROFILE
                     Named profile to use from the shared credentials
                     and config files.

              AWS_REGION, AWS_DEFAULT_REGION
                     AWS region to use for the Route 53 client.

       ~/.aws/credentials, ~/.aws/config
              Shared credentials and config files, as used by the AWS
              CLI (see "aws configure").

       In the absence of any of the above, an IAM role attached to the
       host (for example an EC2 instance profile) is used if present.

       The IAM identity used must be permitted to call
       route53:ListHostedZonesByName, route53:ListResourceRecordSets,
       and route53:ChangeResourceRecordSets on the relevant hosted
       zone(s).

EXAMPLES
       Update a single domain:

              dynDNS -d home.example.com

       Update several domains in one run:

              dynDNS -d home.example.com -d vpn.example.com \
                     -d requests.example.com

SEE ALSO
       dyndns(8), aws(1)

AUTHOR
       Joshua W.B.
`

func usage() {
	fmt.Fprint(os.Stderr, usageText)
}

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
	flag.Usage = usage
	flag.Var(&domains, "domain", "Domain to update in Route 53 (repeatable)")
	flag.Var(&domains, "d", "Shorthand for --domain (repeatable)")
	flag.Parse()

	if len(domains) == 0 {
		slog.Error("At least one --domain/-d flag is required.")
		usage()
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
