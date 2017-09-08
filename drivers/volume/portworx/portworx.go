package portworx

import (
	"fmt"
	"log"
	"strings"
	"time"

	dockerclient "github.com/fsouza/go-dockerclient"

	"github.com/libopenstorage/openstorage/api"
	clusterclient "github.com/libopenstorage/openstorage/api/client/cluster"
	volumeclient "github.com/libopenstorage/openstorage/api/client/volume"
	"github.com/libopenstorage/openstorage/cluster"
	"github.com/libopenstorage/openstorage/volume"
	"github.com/portworx/torpedo/drivers/scheduler"
	torpedovolume "github.com/portworx/torpedo/drivers/volume"
)

var (
	docker *dockerclient.Client
)

// DriverName is the name of the portworx driver implementation
const DriverName = "pxd"

type portworx struct {
	hostConfig     *dockerclient.HostConfig
	clusterManager cluster.Cluster
	volDriver      volume.VolumeDriver
	schedDriver    scheduler.Driver
}

func (d *portworx) String() string {
	return DriverName
}

func (d *portworx) Init(sched string) error {
	log.Printf("Using the Portworx volume driver under scheduler: %v\n", sched)
	var err error
	d.schedDriver, err = scheduler.Get(sched)
	if err != nil {
		return err
	}

	nodes := d.schedDriver.GetNodes()

	var endpoint string
	for _, n := range nodes {
		if n.Type == scheduler.NodeTypeWorker {
			endpoint = n.Addresses[0]
			break
		}
	}

	if len(endpoint) == 0 {
		return fmt.Errorf("failed to get endpoint for portworx volume driver")
	}

	log.Printf("Using %v as endpoint for portworx volume driver\n", endpoint)
	clnt, err := clusterclient.NewClusterClient("http://"+endpoint+":9001", "v1")
	if err != nil {
		return err
	}
	d.clusterManager = clusterclient.ClusterManager(clnt)

	clnt, err = volumeclient.NewDriverClient("http://"+endpoint+":9001", "pxd-sched", "")
	if err != nil {
		return err
	}
	d.volDriver = volumeclient.VolumeDriver(clnt)

	cluster, err := d.clusterManager.Enumerate()
	if err != nil {
		return err
	}

	log.Printf("The following Portworx nodes are in the cluster:\n")
	for _, n := range cluster.Nodes {
		log.Printf(
			"\tNode UID: %v\tNode IP: %v\tNode Status: %v\n",
			n.Id,
			n.DataIp,
			n.Status,
		)
	}

	return err
}

func (d *portworx) CleanupVolume(name string) error {
	locator := &api.VolumeLocator{}

	volumes, err := d.volDriver.Enumerate(locator, nil)
	if err != nil {
		return err
	}

	for _, v := range volumes {
		if v.Locator.Name == name {
			// First unmount this volume at all mount paths...
			for _, path := range v.AttachPath {
				if err = d.volDriver.Unmount(v.Id, path); err != nil {
					err = fmt.Errorf(
						"Error while unmounting %v at %v because of: %v",
						v.Id,
						path,
						err,
					)
					log.Printf("%v", err)
					return err
				}
			}

			if err = d.volDriver.Detach(v.Id); err != nil {
				err = fmt.Errorf(
					"Error while detaching %v because of: %v",
					v.Id,
					err,
				)
				log.Printf("%v", err)
				return err
			}

			if err = d.volDriver.Delete(v.Id); err != nil {
				err = fmt.Errorf(
					"Error while deleting %v because of: %v",
					v.Id,
					err,
				)
				log.Printf("%v", err)
				return err
			}

			log.Printf("Succesfully removed Portworx volume %v\n", name)

			return nil
		}
	}

	return nil
}

func (d *portworx) InspectVolume(name string) error {
	return nil
}

// Portworx runs as a container - so all we need to do is ask docker to
// stop the running portworx container.
func (d *portworx) StopDriver(ip string) error {
	endpoint := "tcp://" + ip + ":2375"
	docker, err := dockerclient.NewClient(endpoint)
	if err != nil {
		return err
	}

	if err = docker.Ping(); err != nil {
		return err
	}

	// Find and stop the Portworx container
	lo := dockerclient.ListContainersOptions{
		All:  true,
		Size: false,
	}

	allContainers, err := docker.ListContainers(lo)
	if err != nil {
		return err
	}

	for _, c := range allContainers {
		info, err := docker.InspectContainer(c.ID)
		if err != nil {
			return err
		}

		if strings.Contains(info.Config.Image, "px") {
			if !info.State.Running {
				return fmt.Errorf(
					"portworx container with UID %v is not running",
					c.ID,
				)
			}

			d.hostConfig = info.HostConfig
			log.Printf("Stopping Portworx container with UID: %v\n", c.ID)
			if err = docker.StopContainer(c.ID, 0); err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("Could not find the Portworx container on %v", ip)
}

func (d *portworx) WaitStart(ip string) error {
	// Wait for Portworx to become usable.
	status, _ := d.clusterManager.NodeStatus()
	for i := 0; status != api.Status_STATUS_OK; i++ {
		if i > 60 {
			return fmt.Errorf(
				"Portworx did not start up in time: Status is %v",
				status,
			)
		}

		time.Sleep(1 * time.Second)
		status, _ = d.clusterManager.NodeStatus()
	}

	return nil
}

func (d *portworx) StartDriver(ip string) error {
	endpoint := "tcp://" + ip + ":2375"
	docker, err := dockerclient.NewClient(endpoint)
	if err != nil {
		return err
	}

	if err = docker.Ping(); err != nil {
		return err
	}

	// Find and stop the Portworx container
	lo := dockerclient.ListContainersOptions{
		All:  true,
		Size: false,
	}

	allContainers, err := docker.ListContainers(lo)
	if err != nil {
		return err
	}

	for _, c := range allContainers {
		info, err := docker.InspectContainer(c.ID)
		if err != nil {
			return err
		}

		if strings.Contains(info.Config.Image, "px") {
			if info.State.Running {
				return fmt.Errorf(
					"portworx container with UID %v is not stopped",
					c.ID,
				)
			}

			log.Printf("Starting Portworx container with UID: %v\n", c.ID)
			if err = docker.StartContainer(c.ID, d.hostConfig); err != nil {
				return err
			}

			return d.WaitStart(ip)
		}
	}

	log.Printf("Could not fine the Portworx container.\n")
	return fmt.Errorf("could not find the Portworx container on %v", ip)
}

func init() {
	torpedovolume.Register(DriverName, &portworx{})
}
