package transport

// Append and Range describe the stable gRPC boundary. The service uses the
// same domain structs so a generated protobuf adapter can be added without
// coupling storage to transport concerns.
type AppendRequest struct {
	Tenant  string
	Samples []PointRequest
}
type PointRequest struct {
	Metric    string
	Labels    map[string]string
	Timestamp int64
	Value     float64
	Quality   uint32
}
type AppendResponse struct {
	Accepted  int
	Duplicate int
	Error     string
}
type RangeRequest struct {
	Tenant, Query      string
	StartUnix, EndUnix int64
}
type RangeResponse struct{ Series map[string][]PointResponse }
type PointResponse struct {
	Timestamp int64
	Value     float64
	Quality   uint32
}
