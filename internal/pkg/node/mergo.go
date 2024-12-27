package node

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"strings"

	"dario.cat/mergo"

	"github.com/warewulf/warewulf/internal/pkg/util"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

// copyProfile creates a deep copy of the given Profile object.
// It uses encoding/gob to serialize and deserialize the Profile, ensuring
// that all nested fields are copied.
//
// Parameters:
// - original: The Profile object to be copied.
//
// Returns:
// - A new Profile object that is a deep copy of the input.
// - An error if serialization or deserialization fails.
func copyProfile(original Profile) (Profile, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	profile := Profile{}
	if err := enc.Encode(original); err != nil {
		return profile, err
	} else {
		if err := dec.Decode(&profile); err != nil {
			return profile, err
		} else {
			return profile, nil
		}
	}
}

// getNodeProfiles retrieves a list of profile IDs associated with a specific node ID.
// It retrives nested profiles and ensures the list is cleaned of duplicates
// and negations (denoted with a '~' prefix).
//
// Parameters:
// - id: The identifier of the node whose profiles are to be retrieved.
//
// Returns:
// - A slice of profile IDs associated with the given node ID.
func (config *NodesYaml) getNodeProfiles(id string) (profiles []string) {
	if node, ok := config.Nodes[id]; ok {
		for _, profileID := range node.Profiles {
			profiles = cleanList(append(profiles, profileID))
			if !strings.HasPrefix(profileID, "~") {
				profiles = config.appendProfileProfiles(profiles, profileID)
			}
		}
	}
	return cleanList(profiles)
}

// appendProfileProfiles recursively appends profile IDs associated with a given profile ID
// to the provided list of profile IDs. It recursively processes nested profiles and ensures
// the list is cleaned of duplicates and negations (denoted with a '~' prefix).
//
// Profiles are only added if they do not already exist in the list.
//
// Parameters:
// - profiles: A slice of strings representing the current list of profiles by ID.
// - id: The identifier of the profile whose associated profiles are to be appended.
//
// Returns:
//   - A slice of strings containing the updated list of profile IDs.
func (config *NodesYaml) appendProfileProfiles(profiles []string, id string) []string {
	if profile, ok := config.NodeProfiles[id]; ok {
		for _, subID := range profile.Profiles {
			if !util.InSlice(profiles, subID) {
				profiles = cleanList(append(profiles, subID))
				if !strings.HasPrefix(subID, "~") {
					profiles = config.appendProfileProfiles(profiles, subID)
				}
			}
		}
	}
	return profiles
}

// MergeNode merges the configuration of a node identified by `id` with all the profiles
// associated with it, producing a fully composed `Node` and a `fieldMap` detailing the
// sources of various configuration fields.
//
// It works by:
//   - Retrieving the base node configuration using `GetNodeOnly`.
//   - Gathering all profile IDs associated with the node via `getNodeProfiles`.
//   - For each profile:
//   - Merging fields from a deep copy of each profile into the node,
//     recording the origin of each configuration field (i.e., which profile provided it)
//     in a `fieldMap` so that traceability is maintained.
//   - Finally, merging the original node configuration back into the processed node, ensuring
//     that any fields not set by the profiles are preserved, and updating the `fieldMap`
//     accordingly.
//
// Parameters:
// - id: The identifier of the node to be merged with its profiles.
//
// Returns:
// - node: The resulting merged `Node` configuration.
// - fields: A `fieldMap` detailing the source(s) of each configuration field.
// - err: An error if any node or profile retrieval or merging operations fail.
func (config *NodesYaml) MergeNode(id string) (node Node, fields fieldMap, err error) {
	node, err = config.GetNodeOnly(id)
	if err != nil {
		return node, fields, err
	}
	originalNode := node
	node = Node{}

	fields = make(fieldMap)

	for _, profileID := range config.getNodeProfiles(id) {
		if profile, err := config.GetProfile(profileID); err != nil {
			wwlog.Warn("profile not found: %s", profileID)
			continue
		} else if profile, err := copyProfile(profile); err != nil {
			wwlog.Warn("error processing profile %s: %v", profileID, err)
			continue
		} else {
			if err = mergo.Merge(&node.Profile, profile, mergo.WithAppendSlice, mergo.WithOverride); err != nil {
				return node, fields, err
			}
			for _, fieldName := range listFields(profile) {
				if value, err := getNestedFieldValue(profile, fieldName); err == nil && valueStr(value) != "" {
					source := profileID
					prevSource := fields.Source(fieldName)
					if value.Kind() == reflect.Slice && prevSource != "" {
						source = strings.Join([]string{prevSource, source}, ",")
					}
					if value, err := getNestedFieldString(node, fieldName); err == nil {
						fields.Set(fieldName, source, value)
					}
				}
			}
		}
	}

	if err = mergo.Merge(&node, originalNode, mergo.WithAppendSlice, mergo.WithOverride); err != nil {
		return node, fields, err
	}
	for _, fieldName := range listFields(originalNode) {
		if value, err := getNestedFieldValue(originalNode, fieldName); err == nil && valueStr(value) != "" {
			source := ""
			prevSource := fields.Source(fieldName)
			if value.Kind() == reflect.Slice && prevSource != "" {
				source = strings.Join([]string{prevSource, id}, ",")
			}
			if value, err := getNestedFieldString(node, fieldName); err == nil {
				fields.Set(fieldName, source, value)
			}
		}
	}

	node.Profiles = originalNode.Profiles
	if len(node.Profiles) > 0 {
		fields.Set("Profiles", "", strings.Join(originalNode.Profiles, ","))
		fields["Profiles"].Source = ""
	} else {
		delete(fields, "Profiles")
	}

	node.id = id
	node.valid = true
	node.updatePrimaryNetDev()
	return node, fields, nil
}
