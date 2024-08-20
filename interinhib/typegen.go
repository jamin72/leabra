// Code generated by "core generate -add-types"; DO NOT EDIT.

package interinhib

import (
	"cogentcore.org/core/types"
)

var _ = types.AddType(&types.Type{Name: "github.com/emer/leabra/v2/interinhib.InterInhib", IDName: "inter-inhib", Doc: "InterInhib specifies inhibition between layers, where\nthe receiving layer either does a Max or Add of portion of\ninhibition from other layer(s).", Fields: []types.Field{{Name: "Lays", Doc: "layers to receive inhibition from"}, {Name: "Gi", Doc: "multiplier on Gi from other layers"}, {Name: "Add", Doc: "add inhibition -- otherwise Max"}}})