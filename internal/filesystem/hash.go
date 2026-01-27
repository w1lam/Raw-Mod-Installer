package filesystem

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileHashEntry struct {
	Path   string
	Size   int64
	Sha512 string
}

func ComputeDirHash(dir string) (string, error) {
	var entries []FileHashEntry

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(strings.ToLower(d.Name()), ".json") {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		sha, err := computeFileSHA512(path)
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		entries = append(entries, FileHashEntry{
			Path:   filepath.ToSlash(rel),
			Size:   info.Size(),
			Sha512: sha,
		})

		return nil
	})
	if err != nil {
		return "", nil
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	h := sha512.New()
	enc := json.NewEncoder(h)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(entries); err != nil {
		return "", nil
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func computeFileSHA512(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha512.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
