package mapper

import (
	"testing"
)

func TestValidateMappingString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid uid only", "1000", false},
		{"valid uid and gid", "1000:1000", false},
		{"valid uid mapping", "1000=1000", false},
		{"valid uid and gid mapping", "1000:1000=1000:1000", false},
		{"invalid format", "abc", true},
		{"invalid format with colon", "1000:abc", true},
		{"invalid format with equals", "1000=abc", true},
		{"invalid format with both", "1000:abc=1000:def", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMappingString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMappingString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateInput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Mapping
		wantErr bool
	}{
		{
			name:  "valid uid only",
			input: "1000",
			want: &Mapping{
				ContainerUID: 1000,
				ContainerGID: 1000,
				HostUID:      1000,
				HostGID:      1000,
			},
			wantErr: false,
		},
		{
			name:  "valid uid and gid",
			input: "1000:2000",
			want: &Mapping{
				ContainerUID: 1000,
				ContainerGID: 2000,
				HostUID:      1000,
				HostGID:      2000,
			},
			wantErr: false,
		},
		{
			name:  "valid uid mapping",
			input: "1000=2000",
			want: &Mapping{
				ContainerUID: 1000,
				ContainerGID: 1000,
				HostUID:      2000,
				HostGID:      2000,
			},
			wantErr: false,
		},
		{
			name:  "valid uid and gid mapping",
			input: "1000:2000=3000:4000",
			want: &Mapping{
				ContainerUID: 1000,
				ContainerGID: 2000,
				HostUID:      3000,
				HostGID:      4000,
			},
			wantErr: false,
		},
		{
			name:    "invalid uid range",
			input:   "70000",
			wantErr: true,
		},
		{
			name:    "invalid gid range",
			input:   "1000:70000",
			wantErr: true,
		},
		{
			name:    "invalid host uid range",
			input:   "1000=70000",
			wantErr: true,
		},
		{
			name:    "invalid host gid range",
			input:   "1000:1000=1000:70000",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				if got.ContainerUID != tt.want.ContainerUID {
					t.Errorf("ContainerUID = %v, want %v", got.ContainerUID, tt.want.ContainerUID)
				}
				if got.ContainerGID != tt.want.ContainerGID {
					t.Errorf("ContainerGID = %v, want %v", got.ContainerGID, tt.want.ContainerGID)
				}
				if got.HostUID != tt.want.HostUID {
					t.Errorf("HostUID = %v, want %v", got.HostUID, tt.want.HostUID)
				}
				if got.HostGID != tt.want.HostGID {
					t.Errorf("HostGID = %v, want %v", got.HostGID, tt.want.HostGID)
				}
			}
		})
	}
}

func TestCreateMap(t *testing.T) {
	tests := []struct {
		name   string
		idType string
		idList [][2]int
		want   []string
	}{
		{
			name:   "single mapping",
			idType: "u",
			idList: [][2]int{{1000, 1000}},
			want: []string{
				"lxc.idmap: u 0 100000 1000",
				"lxc.idmap: u 1000 1000 1",
				"lxc.idmap: u 1001 101001 64535",
			},
		},
		{
			name:   "multiple mappings",
			idType: "g",
			idList: [][2]int{{1000, 1000}, {1002, 1002}},
			want: []string{
				"lxc.idmap: g 0 100000 1000",
				"lxc.idmap: g 1000 1000 1",
				"lxc.idmap: g 1001 101001 1",
				"lxc.idmap: g 1002 1002 1",
				"lxc.idmap: g 1003 101003 64533",
			},
		},
		{
			name:   "sequential mappings",
			idType: "u",
			idList: [][2]int{{1000, 1000}, {1001, 1001}},
			want: []string{
				"lxc.idmap: u 0 100000 1000",
				"lxc.idmap: u 1000 1000 1",
				"lxc.idmap: u 1001 1001 1",
				"lxc.idmap: u 1002 101002 64534",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateMap(tt.idType, tt.idList)
			if len(got) != len(tt.want) {
				t.Errorf("CreateMap() returned %d lines, want %d", len(got), len(tt.want))
				return
			}
			for i, line := range got {
				if line != tt.want[i] {
					t.Errorf("CreateMap()[%d] = %v, want %v", i, line, tt.want[i])
				}
			}
		})
	}
}
