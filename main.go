package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	driverName    = "vstorage.csi.virtuozzo.com"
	driverVersion = "v1.0.0"
	endpoint      = "/csi/csi.sock"
	storagePath   = "/pstorage"
)

type virtuozzoStorageDriver struct {
	csi.UnimplementedIdentityServer
	csi.UnimplementedControllerServer
	csi.UnimplementedNodeServer
}

func NewVirtuozzoStorageDriver() *virtuozzoStorageDriver {
	fmt.Println("Initializing Virtuozzo Storage Driver")
	return &virtuozzoStorageDriver{}
}

func (d *virtuozzoStorageDriver) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	fmt.Println("GetPluginInfo called")
	return &csi.GetPluginInfoResponse{
		Name:          driverName,
		VendorVersion: driverVersion,
	}, nil
}

func (d *virtuozzoStorageDriver) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	volumeID := req.GetName()
	volumePath := filepath.Join(storagePath, volumeID)

	fmt.Printf("Creating volume with ID %s at path %s\n", volumeID, volumePath)

	if err := os.MkdirAll(volumePath, 0755); err != nil {
		fmt.Printf("Error creating volume directory: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to create volume directory: %v", err)
	}
    #cmd, err := exec.Command("/usr/sbin/ploop", "init" "-s 1G", string(volumePath))
	#    if (err != nil) {
	#	    fmt.Println(err)
	#    return
				  }
	cmd.Close()



	fmt.Printf("Volume %s has been successfully created at %s\n", volumeID, volumePath)

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      volumeID,
			CapacityBytes: req.GetCapacityRange().GetRequiredBytes(),
		},
	}, nil
}

func (d *virtuozzoStorageDriver) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	volumeID := req.GetVolumeId()
	volumePath := filepath.Join(storagePath, volumeID)

	fmt.Printf("Deleting volume with ID %s at path %s\n", volumeID, volumePath)

	if err := os.RemoveAll(volumePath); err != nil {
		fmt.Printf("Error deleting volume directory: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to delete volume directory: %v", err)
	}

	fmt.Printf("Volume %s has been successfully deleted from %s\n", volumeID, volumePath)

	return &csi.DeleteVolumeResponse{}, nil
}

func (d *virtuozzoStorageDriver) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
    return &csi.NodeGetInfoResponse{
        NodeId: "rke-w3",
    }, nil
}

func (d *virtuozzoStorageDriver) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
    return &csi.ProbeResponse{}, nil
	}


func (d *virtuozzoStorageDriver) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
    fmt.Println("GetPluginCapabilities called")
    return &csi.GetPluginCapabilitiesResponse{
        Capabilities: []*csi.PluginCapability{
            {
                Type: &csi.PluginCapability_Service_{
                    Service: &csi.PluginCapability_Service{
                        Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
                    },
                },
            },
            // Добавьте другие capabilities в зависимости от поддерживаемых вашим драйвером функций
        },
    }, nil
}

func (d *virtuozzoStorageDriver) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
    fmt.Println("ControllerGetCapabilities called")
    return &csi.ControllerGetCapabilitiesResponse{
        Capabilities: []*csi.ControllerServiceCapability{
            {
                Type: &csi.ControllerServiceCapability_Rpc{
                    Rpc: &csi.ControllerServiceCapability_RPC{
                        Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
                    },
                },
            },
            // Добавьте другие capabilities в зависимости от поддерживаемых вашим драйвером функций
        },
    }, nil
}

func (d *virtuozzoStorageDriver) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
    fmt.Println("NodeGetCapabilities called")
    return &csi.NodeGetCapabilitiesResponse{
        Capabilities: []*csi.NodeServiceCapability{
            {
                Type: &csi.NodeServiceCapability_Rpc{
                    Rpc: &csi.NodeServiceCapability_RPC{
                        Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
                    },
                },
            },
            // Добавьте другие capabilities в зависимости от поддерживаемых вашим драйвером функций
        },
    }, nil
}



func main() {
    fmt.Printf("Starting Virtuozzo Storage CSI driver on endpoint %s\n", endpoint)
    listener, err := net.Listen("unix", endpoint)
    if err != nil {
        fmt.Printf("Failed to listen: %v\n", err)
        os.Exit(1)
    }

    driver := NewVirtuozzoStorageDriver() // Создаём экземпляр драйвера один раз

    server := grpc.NewServer()
    csi.RegisterIdentityServer(server, driver) // Передаём тот же экземпляр в каждую функцию регистрации
    csi.RegisterControllerServer(server, driver)
    csi.RegisterNodeServer(server, driver)

    if err := server.Serve(listener); err != nil {
        fmt.Printf("Failed to serve: %v\n", err)
        os.Exit(1)
    }
}


