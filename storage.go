package casc

import (
	"encoding/hex"
	"errors"
	"io"

	"github.com/jybp/casc/downloader"
	"github.com/jybp/casc/root/diablo3"
)

const (
	RegionUS = "us"
	RegionEU = "eu"
	RegionKR = "kr"
	RegionTW = "tw"
	RegionCN = "cn"
)

const (
	Diablo3 = "d3"
	// WorldOfWarcraft  = "wow"
	// HeroesOfTheStorm = "hero"
	// Hearthstone      = "hsb"
	// Overwatch        = "pro"
	// Starcraft1       = "s1"
	// Starcraft2       = "s2"
	// Warcraft3        = "w3"
)

// Downloader is the interface that wraps the Get method.
// Get should retrieve the file associated with rawurl.
type Downloader interface {
	Get(rawurl string) (io.ReadSeeker, error)
}

// each app has its own way of relating file names to content hash
type root interface {
	Files() ([]string, error)
	ContentHash(filename string) ([]byte, error)
}

type Storage struct {
	App        string
	Region     string
	Downloader Downloader

	extractor *extractor
	root      root
}

func (s *Storage) app() string {
	if s.App == "" {
		return Diablo3
	}
	return s.App
}

func (s *Storage) region() string {
	if s.Region == "" {
		return RegionUS
	}
	return s.Region
}

func (s *Storage) downloader() Downloader {
	if s.Downloader == nil {
		return &downloader.HTTP{}
	}
	return s.Downloader
}

// Version returns the version of s.App on s.Region.
func (s *Storage) Version() (string, error) {
	if err := s.setupExtractor(); err != nil {
		return "", err
	}
	return s.extractor.version.Name, nil
}

// Files enumerates all files
func (s *Storage) Files() ([]string, error) {
	if err := s.setupRoot(); err != nil {
		return nil, err
	}
	return s.root.Files()
}

// Extract extracts the file with the filename name
func (s *Storage) Extract(filename string) ([]byte, error) {
	if err := s.setupRoot(); err != nil {
		return nil, err
	}
	contentHash, err := s.root.ContentHash(filename)
	if err != nil {
		return nil, err
	}
	return s.extractor.extract(contentHash)
}

// initialize s.extractor
func (s *Storage) setupExtractor() error {
	if s.extractor != nil {
		return nil
	}
	extractor, err := newExtractor(s.downloader(), s.app(), s.region())
	if err != nil {
		return err
	}
	s.extractor = extractor
	return nil
}

// initialize s.root
func (s *Storage) setupRoot() error {
	if s.root != nil {
		return nil
	}
	if err := s.setupExtractor(); err != nil {
		return err
	}
	rootHash := make([]byte, hex.DecodedLen(len(s.extractor.build.RootHash)))
	if _, err := hex.Decode(rootHash, []byte(s.extractor.build.RootHash)); err != nil {
		return err
	}
	switch s.App {
	case Diablo3:
		s.root = &diablo3.Root{RootHash: rootHash, Extract: s.extractor.extract}
	default:
		return errors.New("unsupported app")
	}
	return nil
}
