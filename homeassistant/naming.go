package homeassistant

import "strings"

// entityNameSeparator joins the device and entity name parts. Matches the
// frontend DEFAULT_SEPARATOR used by formatEntityName.
const entityNameSeparator = " "

// entityNameStripSuffixes mirrors the frontend SUFFIXES used when removing a
// device-name prefix from an entity name.
var entityNameStripSuffixes = []string{" ", ": ", " - "}

// EntityNamer resolves structured "device + entity" display names from the
// entity and device registries. It mirrors the Home Assistant frontend
// computeEntityName / formatEntityName model documented in the entity naming
// migration (epic home-assistant/epics#44).
type EntityNamer struct {
	entities map[string]EntityRegistryDisplayEntry
	devices  map[string]DeviceRegistryEntry
}

// NewEntityNamer builds a namer from a registry display snapshot and the device
// registry list.
func NewEntityNamer(display EntityRegistryDisplayResponse, devices []DeviceRegistryEntry) *EntityNamer {
	entityMap := make(map[string]EntityRegistryDisplayEntry, len(display.Entities))
	for _, entry := range display.Entities {
		entityMap[entry.EntityID] = entry
	}

	deviceMap := make(map[string]DeviceRegistryEntry, len(devices))
	for _, device := range devices {
		deviceMap[device.ID] = device
	}

	return &EntityNamer{entities: entityMap, devices: deviceMap}
}

// DisplayName returns the structured display name for an entity, composing the
// device name and the entity-specific name. It falls back to friendlyName when
// the entity is missing from the registry or has no resolvable structured name.
func (namer *EntityNamer) DisplayName(entityID string, friendlyName string) string {
	if namer == nil {
		return friendlyName
	}

	entry, ok := namer.entities[entityID]
	if !ok {
		return friendlyName
	}

	entityName := entry.Name

	var device DeviceRegistryEntry
	hasDevice := false
	if entry.DeviceID != "" {
		device, hasDevice = namer.devices[entry.DeviceID]
	}

	if !hasDevice {
		if entityName != "" {
			return entityName
		}
		return friendlyName
	}

	deviceName := computeDeviceName(device)

	// If the entity name matches the device name, treat the entity name as
	// empty and use the device name alone.
	if deviceName == entityName {
		entityName = ""
	} else if deviceName != "" && entityName != "" {
		if stripped := stripPrefixFromEntityName(entityName, deviceName); stripped != "" {
			entityName = stripped
		}
	}

	parts := make([]string, 0, 2)
	if deviceName != "" {
		parts = append(parts, deviceName)
	}
	if entityName != "" {
		parts = append(parts, entityName)
	}

	if len(parts) == 0 {
		return friendlyName
	}

	return strings.Join(parts, entityNameSeparator)
}

// computeDeviceName mirrors the frontend computeDeviceName: the user-set name
// takes precedence over the integration-provided name.
func computeDeviceName(device DeviceRegistryEntry) string {
	if name := strings.TrimSpace(device.NameByUser); name != "" {
		return name
	}

	return strings.TrimSpace(device.Name)
}

// stripPrefixFromEntityName removes a device-name prefix from an entity name,
// mirroring the frontend strip_prefix_from_entity_name. It returns an empty
// string when the prefix does not apply.
func stripPrefixFromEntityName(entityName string, prefix string) string {
	lowerName := strings.ToLower(entityName)
	lowerPrefix := strings.ToLower(prefix)

	for _, suffix := range entityNameStripSuffixes {
		prefixWithSuffix := lowerPrefix + suffix
		if !strings.HasPrefix(lowerName, prefixWithSuffix) {
			continue
		}

		newName := entityName[len(prefixWithSuffix):]
		if newName == "" {
			continue
		}

		// Match the frontend: capitalize the first word unless it already
		// contains an upper-case letter (e.g. a brand name).
		firstWord := ""
		if idx := strings.Index(newName, " "); idx >= 0 {
			firstWord = newName[:idx]
		}

		if hasUpperCase(firstWord) {
			return newName
		}

		return capitalizeFirst(newName)
	}

	return ""
}

func hasUpperCase(value string) bool {
	return strings.ToLower(value) != value
}

func capitalizeFirst(value string) string {
	if value == "" {
		return value
	}

	runes := []rune(value)
	runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
	return string(runes)
}
