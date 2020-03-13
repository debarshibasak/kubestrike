## Kubestrike

Kubestrike is a tooling for creating kubernetes clusters in an automated fashion on ubuntu, centos and redhat.
This is an alternative to kubespray. 
It does not provide as many features as kubespray, however mission of the project is to be as faster as possible to provision clusters.
You can create HA and Non-HA clusters with kubestrike.
In an HA scenario, it also provisions an HA Proxy.

Also, it support cluster creation across various Cloud based kubernetes engines.

#### Installation

```.env
go install github.com/debarshi/kubestrike
```

#### Support providers

- Multipass
- Baremetal/VM

#### Using CLI to create clusters 

Build the project as follows

```
go install github.com/debarshibasak/kubestrike
kubestrike --provider multipass --master-count 2 --worker-count 2 --cluster-name test
```
This command will actually acquire 2+2+1 instances in multipass. 1 extra instance to provision HAProxy.
Currently the cli only supports multipass.

#### Roadmap
- Following Providers will be supported soon :-
- AWS
- GCP
- Digital Ocean