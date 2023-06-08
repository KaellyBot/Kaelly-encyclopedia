package constants

type Source struct {
	Name string
	Icon string
	URL  string
}

func GetEncyclopediasSource() Source {
	return Source{
		Name: "Dofusdude",
		Icon: "https://docs.dofusdu.de/favicon-32x32.jpg",
		URL:  "http://dofusdu.de",
	}
}
