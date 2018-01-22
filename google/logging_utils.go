package google

import (
	"fmt"
	"regexp"
)

// loggingSinkResourceTypes contains all the possible Stackdriver Logging resource types. Used to parse ids safely.
var loggingSinkResourceTypes = []string{
	"billingAccount",
	"folders",
	"organizations",
	"projects",
}

// LoggingSinkId represents the parts that make up the canonical id used within terraform for a logging resource.
type LoggingSinkId struct {
	resourceType string
	resourceId   string
	name         string
}

// loggingSinkIdRegex matches valid logging sink canonical ids
var loggingSinkIdRegex = regexp.MustCompile("(.+)/(.+)/sinks/(.+)")

// canonicalId returns the LoggingSinkId as the canonical id used within terraform.
func (l LoggingSinkId) canonicalId() string {
	return fmt.Sprintf("%s/%s/sinks/%s", l.resourceType, l.resourceId, l.name)
}

// parent returns the "parent-level" resource that the sink is in (e.g. `folders/foo` for id `folders/foo/sinks/bar`)
func (l LoggingSinkId) parent() string {
	return fmt.Sprintf("%s/%s", l.resourceType, l.resourceId)
}

// parseLoggingSinkId parses a canonical id into a LoggingSinkId, or returns an error on failure.
func parseLoggingSinkId(id string) (*LoggingSinkId, error) {
	parts := loggingSinkIdRegex.FindStringSubmatch(id)
	if parts == nil {
		return nil, fmt.Errorf("unable to parse logging sink id %#v", id)
	}
	// If our resourceType is not a valid logging sink resource type, complain loudly
	validLoggingSinkResourceType := false
	for _, v := range loggingSinkResourceTypes {
		if v == parts[1] {
			validLoggingSinkResourceType = true
			break
		}
	}

	if !validLoggingSinkResourceType {
		return nil, fmt.Errorf("Logging resource type %s is not valid. Valid resource types: %#v", parts[1],
			loggingSinkResourceTypes)
	}
	return &LoggingSinkId{
		resourceType: parts[1],
		resourceId:   parts[2],
		name:         parts[3],
	}, nil
}

// loggingExclusionResourceTypes contains all the possible Stackdriver Logging resource types. Used to parse ids safely.
var loggingExclusionResourceTypes = []string{
	"billingAccount",
	"folders",
	"organizations",
	"projects",
}

// LoggingExclusionId represents the parts that make up the canonical id used within terraform for a logging resource.
type LoggingExclusionId struct {
	resourceType string
	resourceId   string
	name         string
}

// loggingExclusionIdRegex matches valid logging exclusion canonical ids
var loggingExclusionIdRegex = regexp.MustCompile("(.+)/(.+)/exclusions/(.+)")

// canonicalId returns the LoggingExclusionId as the canonical id used within terraform.
func (l LoggingExclusionId) canonicalId() string {
	return fmt.Sprintf("%s/%s/exclusions/%s", l.resourceType, l.resourceId, l.name)
}

// parent returns the "parent-level" resource that the exclusion is in (e.g. `folders/foo` for id `folders/foo/exclusions/bar`)
func (l LoggingExclusionId) parent() string {
	return fmt.Sprintf("%s/%s", l.resourceType, l.resourceId)
}

// parseLoggingExclusionId parses a canonical id into a LoggingExclusionId, or returns an error on failure.
func parseLoggingExclusionId(id string) (*LoggingExclusionId, error) {
	parts := loggingExclusionIdRegex.FindStringSubmatch(id)
	if parts == nil {
		return nil, fmt.Errorf("unable to parse logging exclusion id %#v", id)
	}
	// If our resourceType is not a valid logging exclusion resource type, complain loudly
	validLoggingExclusionResourceType := false
	for _, v := range loggingExclusionResourceTypes {
		if v == parts[1] {
			validLoggingExclusionResourceType = true
			break
		}
	}

	if !validLoggingExclusionResourceType {
		return nil, fmt.Errorf("Logging resource type %s is not valid. Valid resource types: %#v", parts[1],
			loggingExclusionResourceTypes)
	}
	return &LoggingExclusionId{
		resourceType: parts[1],
		resourceId:   parts[2],
		name:         parts[3],
	}, nil
}
