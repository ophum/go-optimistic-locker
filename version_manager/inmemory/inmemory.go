package inmemory

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sync"

	"github.com/google/uuid"
	versionmanager "github.com/ophum/go-optimistic-locker/version_manager"
)

type inmemoryStore struct {
	versions map[string]string
	mu       *sync.RWMutex
}

func NewInmemoryStore() versionmanager.VersionManager {
	return &inmemoryStore{
		versions: make(map[string]string),
		mu:       &sync.RWMutex{},
	}
}

func (s *inmemoryStore) generateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func (s *inmemoryStore) generateHash(seed any) (string, error) {
	j, err := json.Marshal(&seed)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(j)
	return hex.EncodeToString(hash[:]), nil
}

func (s *inmemoryStore) Create(ctx context.Context, path string, opts ...versionmanager.Option) (string, error) {
	return s.put(ctx, path, opts...)
}

func (s *inmemoryStore) Update(ctx context.Context, path string, opts ...versionmanager.Option) (string, error) {
	return s.put(ctx, path, opts...)
}

func (s *inmemoryStore) put(ctx context.Context, path string, opts ...versionmanager.Option) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var option versionmanager.Options
	option.Apply(opts...)

	version := option.GetString("version")
	if version == "" {
		var err error
		version, err = s.generateUUID()
		if err != nil {
			return "", err
		}
	}

	s.versions[path] = version
	return version, nil
}
func (s *inmemoryStore) Get(ctx context.Context, path string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	version, ok := s.versions[path]
	if !ok {
		return "", errors.New("not found")
	}
	return version, nil
}

func (s *inmemoryStore) Delete(ctx context.Context, path string) error {
	s.mu.Unlock()
	defer s.mu.Unlock()
	delete(s.versions, path)
	return nil
}
