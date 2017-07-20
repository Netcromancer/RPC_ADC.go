//A replica of the RPC_ADC.sh script for OpenStack Deployments.
//Original script written by Nathan Pawelek (nathan.pawelek@rackspace.com)
//Go version written by Brandon Bruce (brandon.bruce@rackspace.com)

package main

import (
    "fmt"
    "log"
    "os"
//    "io/ioutil"
//    m "math"
//    "net/http"
    "os/exec"
//    "strconv"
)

//Main calls all the functions.

func main(){
    fmt.Println("INTRO MESSAGE HERE")
    installUpdates()
    getVersion()
    getManufacturer()
    validateDNS()
    removeDeleteme()
    expandNova()
    validateHosts()
}

//Find a place to put global variable for a backup directory (/home/rack/).


//This will update apt and install all the needed packages.
func installUpdates(){
    updates, err := exec.Command("/bin/bash", "-c", "apt-get -qq update && apt-get install dmidecode dnsutils").Output()
        if err != nil {
            log.Fatal(err)
            fmt.Println("Failure to update packages.")
            fmt.Printf("%s", err)
            os.Exit(1)
        }
    fmt.Printf("%s", updates)
}

//This will verify the correct version of Ubuntu (14.04 or 16.04).
func getVersion(){
    version, err := exec.Command("/bin/bash", "-c", "lsb_release -r | awk '{print $NF}'").Output()
    if err != nil {
        log.Fatal(err)
        fmt.Println("Cannot verify Ubuntu release version.")
        fmt.Printf("%s", err)
        os.Exit(1)
    }
    if string(version) != "14.04" && string(version) != "16.04" {
        fmt.Println("Incorrect version. Must be Ubuntu 14.04 or 16.04.")
        fmt.Printf("%s", version)
        os.Exit(1)
    } 
    fmt.Printf("%s", version)
}

//This gathers the manufacturer.
func getManufacturer(){
    manufacturer, err := exec.Command("/bin/bash", "-c", "dmidecode --type 3 | awk '/Manufacturer/ {print $2}'").Output()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s", manufacturer)
}

//This verifies that DNS is working properly.
func validateDNS(){
    dns, err := exec.Command("/bin/bash", "-c", "dig +short rackspace.com").Output()
        if err != nil {
            log.Fatal(err)
            fmt.Println("Unable to validate DNS.")
            fmt.Printf("%s", err)
            os.Exit(1)
        }
    fmt.Printf("%s", dns)
}

//This removes the deleteme lvm.
func removeDeleteme(){
    lvm_count := exec.Command("/bin/bash", "-c", "lvs | grep deleteme | wc -l")
    if lvm_count > 1 {
        fmt.Println("Multiple deleteme logical volumes detected. Please remove the deleteme volumes manually.")
        os.Exit(1)
    } else if lvm_count == 0 {
        fmt.Println("No deleteme logical volumes detected.")
    } else {
        exec.Command("/bin/bash", "-c", "umount /delete")
        exec.Command("/bin/bash", "-c", "rmdir /deleteme")
        exec.Command("/bin/bash", "-c", "lvchange -an /dev/mapper/$(lvs | grep deleteme | awk ' { print $2 }')-deleteme00")
        exec.Command("/bin/bash", "-c", "lvremove -f /dev/mapper/$(lvs | grep deleteme | awk ' { print $2 }')-deleteme00")
        exec.Command("/bin/bash", "-c", "cp /etc/fstab /home/rack/fstab.bak")
        exec.Command("/bin/bash", "-c", "sed -i '/\/deleteme/d' /etc/fstab")
        fmt.Println("The deleteme volume has been successfully removed.")
    }
}

//This expands the "nova" volume.
func expandNova(){
    vg_count := exec.Command("/bin/bash", "-c", "vgs | wc -l")
    space := exec.Command("/bin/bash", "-c", "vgs | awk '{print $7}' | grep -v 'VFree' | sed "s/\..*//"")
    if vg_count > 2 {
        fmt.Println("Mulitple volume groups detected. Please expand the nova volume manually.")
        os.Exit(1)
    } else {
        } if { space != 0
            exec.Command("/bin/bash", "-c", "umount /var/lib/nova")
            exec.Command("/bin/bash", "-c", "lvresize -f -l+100%FREE /dev/$(lvs | grep nova | awk ' { print $2 }')/$(lvs | grep nova | awk '{print $1}')"
            exec.Command("/bin/bash", "-c", "resize2fs -f /dev/mapper/$(lvs | grep nova | awk ' { print $2 }')-$(lvs | grep nova | awk '{print $1}')"
            exec.Command("/bin/bash", "-c", "mount -a"
//Need to add failure catching and exiting/warning here
        } else {
            fmt.Println("The nova logical volume was not detected or there was no free space for expansion.")
        }
}

//This standardizes the /etc/hosts and /etc/hostname files
func validateHosts(){
    iface := exec.Command("/bin/bash", "-c", "route -n | awk '/^0.0.0.0/ {print $NF}'")
    addy := exec.Command("/bin/bash", "-c", "ip a sh ${INTERFACE} | awk '/inet / {sub(/\/[0-9]+$/, "", $2); print $2; exit}'")
    fqdn := exec.Command("/bin/bash", "-c", "echo -e $RS_SERVER_NAME")
    hostname := exec.Command("/bin/bash", "-c", "echo ${FQDN} | awk -F. '{print $1}'")
    if addy != nil && fqdn != nil && hostname != nil{
        exec.Command("/bin/bash", "-c", "cp -f /etc/hosts /home/rack
    }
}

