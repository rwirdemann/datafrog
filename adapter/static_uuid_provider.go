package adapter

type StaticUUIDProvider struct{}

func (StaticUUIDProvider) NewString() string {
	return "023a6a95-6c8a-4483-bcfb-17b1c58c317f"
}
