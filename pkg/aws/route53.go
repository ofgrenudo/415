package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

const defaultTTL = int64(300)

func NewClient(ctx context.Context) (*route53.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	return route53.NewFromConfig(cfg), nil
}

// FindHostedZoneID finds the most specific hosted zone that owns domain by
// walking the domain's labels from most-specific to least (e.g.
// "home.example.com" -> "example.com" -> "com").
func FindHostedZoneID(ctx context.Context, client *route53.Client, domain string) (string, error) {
	fqdn := strings.TrimSuffix(domain, ".") + "."
	labels := strings.Split(strings.TrimSuffix(fqdn, "."), ".")

	for i := range labels {
		candidate := strings.Join(labels[i:], ".") + "."
		out, err := client.ListHostedZonesByName(ctx, &route53.ListHostedZonesByNameInput{
			DNSName: aws.String(candidate),
		})
		if err != nil {
			return "", fmt.Errorf("failed to list hosted zones for %q: %w", candidate, err)
		}
		for _, zone := range out.HostedZones {
			if aws.ToString(zone.Name) == candidate {
				return strings.TrimPrefix(aws.ToString(zone.Id), "/hostedzone/"), nil
			}
		}
	}

	return "", fmt.Errorf("no hosted zone found for domain %q", domain)
}

// GetCurrentARecordIP returns the IP currently published for domain's A
// record in the given hosted zone, or "" if no such record exists.
func GetCurrentARecordIP(ctx context.Context, client *route53.Client, zoneID, domain string) (string, error) {
	fqdn := strings.TrimSuffix(domain, ".") + "."

	out, err := client.ListResourceRecordSets(ctx, &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zoneID),
		StartRecordName: aws.String(fqdn),
		StartRecordType: types.RRTypeA,
		MaxItems:        aws.Int32(1),
	})
	if err != nil {
		return "", fmt.Errorf("failed to list resource record sets for %q: %w", domain, err)
	}

	for _, rrset := range out.ResourceRecordSets {
		if rrset.Type != types.RRTypeA || aws.ToString(rrset.Name) != fqdn {
			continue
		}
		if len(rrset.ResourceRecords) == 0 {
			return "", nil
		}
		return aws.ToString(rrset.ResourceRecords[0].Value), nil
	}

	return "", nil
}

// UpsertARecord creates or updates domain's A record in the given hosted
// zone to point at ip.
func UpsertARecord(ctx context.Context, client *route53.Client, zoneID, domain, ip string) error {
	fqdn := strings.TrimSuffix(domain, ".") + "."

	_, err := client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionUpsert,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(fqdn),
						Type: types.RRTypeA,
						TTL:  aws.Int64(defaultTTL),
						ResourceRecords: []types.ResourceRecord{
							{Value: aws.String(ip)},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to upsert A record for %q: %w", domain, err)
	}

	return nil
}
