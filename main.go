package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/debarshibasak/kubestrike/v1alpha1/config"
)

func main() {

	configuration := flag.String("config", "", "location of configuration")
	install := flag.Bool("install", false, "install operation")
	uninstall := flag.Bool("uninstall", false, "uninstall operation")
	strictInstalltion := flag.Bool("use-strict", false, "uninstall operation")

	flag.Parse()

	log.Println("[kubestrike] started")

	if *install {
		configRaw, err := ioutil.ReadFile(*configuration)
		if err != nil {
			log.Fatal(err)
		}

		clusterOrchestration, err := config.NewParser(*strictInstalltion).Parse(configRaw)
		if err != nil {
			log.Fatal(err)
		}

		if err := clusterOrchestration.Validate(); err != nil {
			log.Fatal(err)
		}

		if err := clusterOrchestration.Install(); err != nil {
			log.Fatal(err)
		}
		return
	}

	if *uninstall {
		log.Println("[kubestrike] not implemented yet")
		return
	}

	log.Fatal("no execution options provided")
}
