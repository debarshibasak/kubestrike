## Kubestrike
[![CodeFactor](https://www.codefactor.io/repository/github/debarshibasak/kubestrike/badge?s=c522561b3a0c2ea3b686df1947f7114e466cca22)](https://www.codefactor.io/repository/github/debarshibasak/kubestrike) 
[![CircleCI](https://circleci.com/gh/debarshibasak/kubestrike.svg?style=svg)](https://circleci.com/gh/debarshibasak/kubestrike)

Kubestrike is a tooling for creating kubernetes clusters in an automated fashion on ubuntu (in future for centos and redhat).
This is an alternative to kubespray. 
It does not provide as many features as kubespray, however mission of the project is to be as fast as possible to provision clusters.
You can create HA and Non-HA clusters with kubestrike.
In an HA scenario, it also provisions an HA Proxy.

Also, it support cluster creation across various Cloud based kubernetes engines.

#### Installation

```.env
go install github.com/debarshi/kubestrike
```

#### Support providers

- Multipass ([example](https://github.com/debarshibasak/kubestrike/tree/master/example/multipass))
- Baremetal/VM ([example](https://github.com/debarshibasak/kubestrike/tree/master/example/baremetal))

#### Using CLI to create clusters 

Run the automation as follows.

A manifest file looks like this.
```
apiVersion: v1
kind: CreateCluster
provider: Multipass
clusterName: testcluster
multipass:
  masterCount: 1
  workerCount: 3
networking:
  podCidr: 10.233.0.0/18
  serviceCidr: 10.233.64.0/18
  plugin: flannel
```

To execute the automation, you have to run as follows.

```
kubestrike --config examples/create_cluster.yml --run
```

#### Roadmap
- Following Providers will be supported soon :-
- AWS
- GCP
- Digital Ocean

#### Supporting this project
- Donate on patreon for testing this project
- If you want to join this project, please feel free to create pull requests.
- You can support my effort with donation at [patreon](https://www.patreon.com/bePatron?u=31747625)


<a href="https://www.patreon.com/bePatron?u=31747625" data-patreon-widget-type="become-patron-button">Become a Patron!</a><script async src="https://c6.patreon.com/becomePatronButton.bundle.js"></script>
