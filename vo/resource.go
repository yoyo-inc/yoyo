package vo

type QueryResourceVO struct {
	ResourceType string `form:"resourceType"`
	Filename     string `form:"filename"`
}
