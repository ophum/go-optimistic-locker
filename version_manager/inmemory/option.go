package inmemory

import versionmanager "github.com/ophum/go-optimistic-locker/version_manager"

func SpecificVersion(version string) versionmanager.Option {
	return versionmanager.WithParam("version", version)
}
