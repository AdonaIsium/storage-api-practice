package store

type VolumeFilter struct {
	NameEquals string
	TagKey     string
	TagValue   string
}

type HostFilter struct {
	NameEquals string
	Identity   string
}

type MappingFilter struct {
	VolumeID string
	HostID   string
}
