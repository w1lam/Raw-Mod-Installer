package lists

type GithubContentResponse struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Sha    string `json:"sha"`
	Size   int    `json:"size"`
	URL    string `json:"url"`
	RawURL string `json:"download_url"`
	Type   string `json:"type"`
}

type AvailableLists struct {
	ModPacks        map[string]ResolvedModPack
	ResourceBundles map[string]ResolvedResourceBundle
}

// GetAllAvailableLists gets the url for a modpack from a list
func GetAllAvailableLists() (AvailableLists, error) {
	modPacks, err := GetAvailableModPacks()
	if err != nil {
		return AvailableLists{}, err
	}

	resourceBundles, err := GetAvailableResourceBundles()
	if err != nil {
		return AvailableLists{}, err
	}

	return AvailableLists{
		ModPacks:        modPacks,
		ResourceBundles: resourceBundles,
	}, nil
}
