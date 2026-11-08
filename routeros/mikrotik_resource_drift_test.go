package routeros

import (
	"reflect"
	"testing"
)

func TestDriftMapIpServiceAvailableFrom(t *testing.T) {
	tests := []struct {
		ros     string
		reverse bool
		key     string
		want    string
	}{
		{"7.22.3", false, "address", ""},
		{"7.22.3", true, "available-from", ""},
		{"7.23.8", false, "address", ""},
		{"7.24", false, "address", "available-from"},
		{"7.24.6", false, "address", "available-from"},
		{"7.24.6", true, "available-from", "address"},
		{"7.25.1", false, "address", "available-from"},
	}

	for _, tt := range tests {
		got := driftAttributeSlice.GetDriftMap(tt.ros, "/ip/service", tt.reverse)[tt.key]
		if got != tt.want {
			t.Errorf("GetDriftMap(%q, reverse=%v)[%q] = %q, want %q", tt.ros, tt.reverse, tt.key, got, tt.want)
		}
	}
}

func TestBuildDatasourceReadFilter(t *testing.T) {
	setTestRouterOSVersion(t, "7.24.6")

	tests := []struct {
		path   string
		filter map[string]any
		want   []string
	}{
		{"/ip/service", map[string]any{"address": "10.0.0.1/32"}, []string{"available-from=10.0.0.1/32"}},
		{"/ip/service", map[string]any{"name": "ssh"}, []string{"name=ssh"}},
		{"/ip/ssh", map[string]any{"always_allow_password_login": "yes"}, []string{"password-authentication=yes"}},
	}

	for _, tt := range tests {
		if got := buildDatasourceReadFilter(tt.filter, tt.path); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("buildDatasourceReadFilter(%v, %q) = %v, want %v", tt.filter, tt.path, got, tt.want)
		}
	}
}

func TestDriftRemovedBgpAddPathOut(t *testing.T) {
	for _, path := range []string{"/routing/bgp/template", "/routing/bgp/connection"} {
		if _, ok := driftAttributeSlice.GetRemoved("7.21", path)["add_path_out"]; ok {
			t.Errorf("%v: add_path_out removed on 7.21", path)
		}
		if _, ok := driftAttributeSlice.GetRemoved("7.22", path)["add_path_out"]; !ok {
			t.Errorf("%v: add_path_out not removed on 7.22", path)
		}
		if _, ok := driftAttributeSlice.GetDriftMap("7.22", path, false)["add_path_out"]; ok {
			t.Errorf("%v: add_path_out is in the forward drift map", path)
		}
		if _, ok := driftAttributeSlice.GetDriftMap("7.22", path, true)[driftRemoved]; ok {
			t.Errorf("%v: the removed marker is in the reverse drift map", path)
		}
	}
}
