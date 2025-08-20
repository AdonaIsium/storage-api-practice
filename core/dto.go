package core

type CreateVolumeRequest struct {
	Name      string
	SizeBytes int64
	Thin      bool
	QosPolicy string
	Tags      map[string]string
}

type CreateHostRequest struct {
	Name       string
	Identities []HostIdentity
}

type MapOptions struct {
	LUN *int
}
