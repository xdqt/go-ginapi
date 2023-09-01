package esexample

type Ellis struct {
	Name string `json:"name"`
}

type Dynamic struct {
	Operation string `json:"operation"`
	Value     string `json:"value"`
	Field     string `json:"field"`
}
