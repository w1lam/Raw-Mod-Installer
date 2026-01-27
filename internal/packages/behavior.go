package packages

import (
	"path/filepath"

	"github.com/w1lam/Packages/modrinth"
	"github.com/w1lam/Raw-Mod-Installer/internal/paths"
)

// PackageType is the type of a package modpack/resourcebundle/shaderbundle?
type PackageType string

// AllPTypes are all pkg types
var AllPTypes = []PackageType{PackageModPack, PackageResourceBundle, PackageShaderBundle}

// PKG TYPES
const (
	PackageModPack        PackageType = "modpack"
	PackageResourceBundle PackageType = "resourcebundle"
	PackageShaderBundle   PackageType = "shaderbundle"
)

// PackageBehavior defines the behavior for a package type
type PackageBehavior struct {
	StorageDir func(*paths.Paths) string
	ActiveDir  func(*paths.Paths) string

	EnsureLoader bool
	EnableAfter  bool

	ResolveFilter modrinth.EntryFilter
}

// PackageBehaviors defines the behavior for a package type
var PackageBehaviors = map[PackageType]PackageBehavior{
	PackageModPack: {
		StorageDir: func(p *paths.Paths) string {
			return filepath.Join(p.PackagesDir, "modpacks")
		},
		ActiveDir: func(p *paths.Paths) string {
			return p.ModsDir
		},
		EnsureLoader: true,
		EnableAfter:  true,
	},

	PackageResourceBundle: {
		StorageDir: func(p *paths.Paths) string {
			return filepath.Join(p.PackagesDir, "resourcebundles")
		},
		ActiveDir: func(p *paths.Paths) string {
			return p.ResourcePacksDir
		},
		EnsureLoader: false,
		EnableAfter:  true,
	},

	PackageShaderBundle: {
		StorageDir: func(p *paths.Paths) string {
			return filepath.Join(p.PackagesDir, "shaderbundles")
		},
		ActiveDir: func(p *paths.Paths) string {
			return p.ShaderPacksDir
		},
		EnsureLoader: false,
		EnableAfter:  true,
	},
}
