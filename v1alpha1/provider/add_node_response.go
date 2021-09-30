package provider

import "github.com/debarshibasak/machina"

type AddNodeResponse struct {
	Master *machina.Node
	Worker []*machina.Node
}
