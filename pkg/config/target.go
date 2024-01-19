package config

type Targets []Target

func (tgts Targets) Aliases() []string {
	aliases := make([]string, 0)
	for _, tgt := range tgts {
		aliases = append(aliases, tgt.Aliases...)
	}
	return aliases
}

type Target struct {
	Title       string   `yaml:"title" json:"title"`
	Description string   `yaml:"description" json:"description"`
	Aliases     []string `yaml:"aliases" json:"aliases"`
	Urls        []ArgUrl `yaml:"urls" json:"urls"`
}

func (tgt Target) ArgUrls(fn func(alias string, url ArgUrl)) {
	for _, alias := range tgt.Aliases {
		for _, argUrl := range tgt.Urls {
			fn(alias, argUrl)
		}
	}
}
