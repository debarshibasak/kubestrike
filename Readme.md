## Kubestrike

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

- Multipass
- Baremetal/VM

#### Using CLI to create clusters 

Run the automation as follows.
```
kubestrike --config examples/multipass_config.yaml --install
```
This command will actually acquire 2+2+1 instances in multipass. 1 extra instance to provision HAProxy.
Currently the cli only supports multipass.

#### Roadmap
- Following Providers will be supported soon :-
- AWS
- GCP
- Digital Ocean

#### Supporting this project
- I need funding for testing this project
- If you want to join this project, please feel free to create pull requests.
- You can support my effort with donation at [patreon](https://www.patreon.com/bePatron?u=31747625)

<a href="https://www.patreon.com/bePatron?u=31747625" data-patreon-widget-type="become-patron-button">Become a Patron!</a><script async src="https://c6.patreon.com/becomePatronButton.bundle.js"></script>