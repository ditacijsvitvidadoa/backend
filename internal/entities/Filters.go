package entities

type CategoryInfo struct {
	LabelUA string `json:"labelUA"`
	Value   string `json:"value"`
}

type Filter struct {
	Title string         `json:"title"`
	Items []CategoryInfo `json:"items"`
}

type FilterCategory struct {
	Categories Filter `json:"categories"`
	Age        Filter `json:"age"`
	Brand      Filter `json:"brand"`
	Material   Filter `json:"material"`
	Type       Filter `json:"type"`
	Discount   Filter `json:"discount"`
}
