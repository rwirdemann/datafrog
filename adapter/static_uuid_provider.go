package adapter

type StaticUUIDProvider struct{}

// NewString always returns the same UUID and should be used only for testing.
func (StaticUUIDProvider) NewString() string {
	return "023a6a95-6c8a-4483-bcfb-17b1c58c317f"
}
