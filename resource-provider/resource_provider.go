package resourceprovider

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/bmatcuk/go-vagrant"
	"github.com/viniciusbds/arrebol-pb-resource-manager/constants"
)

var (
	numNodes = 0
)

func AddNode(vcpu, memory float32) error {
	nodeName := fmt.Sprintf("node%d", numNodes+1)

	vagrantfilePath := path.Join(constants.VAGRANT_PATH, nodeName)

	err := os.Mkdir(vagrantfilePath, os.ModePerm)
	if err != nil {
		return err
	}

	input, err := ioutil.ReadFile(constants.VAGRANTFILE_TEMPLATE_PATH)
	if err != nil {
		return err
	}

	output := bytes.Replace(input, []byte("VBOX_NAME"), []byte(nodeName), -1)
	output = bytes.Replace(output, []byte("MEMORY"), []byte(fmt.Sprintf("%v", memory)), -1)
	output = bytes.Replace(output, []byte("CPUS"), []byte(fmt.Sprintf("%v", vcpu)), -1)

	if err = ioutil.WriteFile(path.Join(vagrantfilePath, "Vagrantfile"), output, os.ModePerm); err != nil {
		return err
	}

	client, err := vagrant.NewVagrantClient(vagrantfilePath)
	if err != nil {
		return err
	}

	upcmd := client.Up()
	upcmd.Verbose = true
	if err := upcmd.Run(); err != nil {
		return err
	}
	if upcmd.Error != nil {
		return err
	}
	numNodes++
	return nil
}

func RemoveNode(nodeName string) error {
	vagrantfilePath := path.Join(constants.VAGRANT_PATH, nodeName)

	client, err := vagrant.NewVagrantClient(vagrantfilePath)
	if err != nil {
		return err
	}
	destroycmd := client.Destroy()
	destroycmd.Verbose = true
	if err := destroycmd.Run(); err != nil {
		return err
	}
	if destroycmd.Error != nil {
		return err
	}
	if err := os.RemoveAll(path.Join(constants.VAGRANT_PATH, nodeName)); err != nil {
		return err
	}
	return nil
}