package main

import (
	"fmt"
	"log"
	"os"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
	"github.com/debarshibasak/kubestrike/providers"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "provider",
				Usage:    "set a provider",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "master-count",
				Usage:    "master count",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "worker-count",
				Usage:    "worker count",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "cluster-name",
				Usage:    "name of the cluster",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "config",
				Usage: "configuration file",
			},
			&cli.StringFlag{
				Name:  "cni",
				Usage: "choose the networking layer",
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "enable verbose mode",
			},
		},
		Action: func(c *cli.Context) error {
			provider := c.String("provider")
			fmt.Println(provider)

			log.Println("creating vm...")

			provider = c.String("provider")

			if provider == "" {
				log.Fatal("provider is not set")
			}

			c.String("config")

			masterNodes, workerNodes, haproxy, err := providers.Get(
				provider,
				c.Int("master-count"),
				c.Int("worker-count"),
			)
			if err != nil {
				log.Fatal(err)
			}

			var networking *kubeadmclient.Networking

			cni := c.String("cni")
			if cni == "" {
				networking = kubeadmclient.Flannel
			} else {
				networking := kubeadmclient.LookupNetworking(cni)
				if networking == nil {
					log.Fatal("network plugin in empty")
				}
			}

			log.Println("creating cluster...")

			kubeadmClient := kubeadmclient.Kubeadm{
				ClusterName: c.String("cluster-name"),
				HaProxyNode: haproxy,
				MasterNodes: masterNodes,
				WorkerNodes: workerNodes,
				VerboseMode: c.Bool("verbose"),
				Netorking:   networking,
			}

			err = kubeadmClient.CreateCluster()
			if err != nil {
				log.Fatal(err)
			}

			printSummary(kubeadmClient)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func printSummary(kubeadm kubeadmclient.Kubeadm) {
	fmt.Println("master machines")
	fmt.Println("-----------")
	for _, master := range kubeadm.MasterNodes {
		fmt.Println(master)
	}

	if kubeadm.HaProxyNode != nil {
		fmt.Println("-----------")
		fmt.Println("haproxy machines")
		fmt.Println("-----------")
		fmt.Println(kubeadm.HaProxyNode)
	}

	fmt.Println("-----------")
	fmt.Println("workers machines")
	fmt.Println("-----------")

	for _, worker := range kubeadm.WorkerNodes {
		fmt.Println(worker)
	}

}
