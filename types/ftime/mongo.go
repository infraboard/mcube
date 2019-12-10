package ftime

// MarshalBSON 实现JSON 序列化接口
func (t Time) MarshalBSON() ([]byte, error) {
	return t.MarshalJSON()
}

// UnmarshalBSON 实现JSON 反序列化接口
func (t *Time) UnmarshalBSON(b []byte) error {
	return t.UnmarshalJSON(b)
}
