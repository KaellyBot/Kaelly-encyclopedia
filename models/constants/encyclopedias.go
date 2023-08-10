package constants

type Source struct {
	Name string
	Icon string
	URL  string
}

func GetEncyclopediasSource() Source {
	return Source{
		Name: "dofusdude",
		Icon: "https://docs.dofusdu.de/favicon-96x96.png",
		URL:  "https://github.com/dofusdude",
	}
}
