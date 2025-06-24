package vm

import (
	"encoding/json"
	"github.com/clarechu/infra-pulse/src/models"
	nets "github.com/shirou/gopsutil/v3/net"
	"reflect"
	"testing"
)

func TestGetNetWorks(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "info",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := GetNetWorks(); (err != nil) != tt.wantErr {
				t.Errorf("GetNetWorks() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log()
		})
	}
}

func TestGetHostname(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "hostname",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHostname()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostname() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Logf("GetHostname() got = %s", got)
			}
		})
	}
}

func TestGetInfo(t *testing.T) {
	tests := []struct {
		name    string
		wantOs  *models.OSInfo
		wantErr bool
	}{
		{
			name:    "xxx",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOs, err := GetInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOs, tt.wantOs) {
				marshal, _ := json.Marshal(gotOs)
				t.Logf("GetInfo() gotOs = %+v", string(marshal))
			}
		})
	}
}

func TestReleaseInfo(t *testing.T) {
	tests := []struct {
		name    string
		want    *models.Release
		wantErr bool
	}{
		{
			name: "x",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReleaseInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReleaseInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Logf("ReleaseInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMainIp(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "get main public",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMainIp()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMainIp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}

func TestGetLastBootTime(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name: "last time",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLastBootTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLastBootTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("GetLastBootTime() got = %v, want %v", got, tt.want)
		})
	}
}

func TestGetProcesses(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "demo1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetProcesses()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProcesses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("got: %+v", got)
		})
	}
}

func TestPorts(t *testing.T) {
	tests := []struct {
		name    string
		want    []nets.ProtoCountersStat
		wantErr bool
	}{
		{
			name:    "xx",
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Ports()
			if (err != nil) != tt.wantErr {
				t.Errorf("Ports() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Ports() got = %v, want %v", got, tt.want)
		})
	}
}

func TestGetProcessesByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    []*models.Processes
		wantErr bool
	}{
		{
			name: "xx",
			args: args{
				name: "CloudHub",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetProcessesByName(tt.args.name, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProcessesByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("got error:%+v", got)
		})
	}
}

func TestDeleteSliceString(t *testing.T) {
	type args struct {
		s    []string
		elem string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "aa",
			args: args{
				s:    []string{"1", "", "", "2", "", "3", "", "4"},
				elem: "",
			},
			want: []string{
				"1", "2", "3", "4",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeleteSliceString(tt.args.s, tt.args.elem); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteSliceString() = %v, want %v", got, tt.want)
			}
		})
	}
}
