package pbsubstreams

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type ModuleKind int

const (
	ModuleKindStore = ModuleKind(iota)
	ModuleKindMap
	ModuleKindBlockIndex
)

func (m *Module) BlockFilterQueryString() (string, error) {
	if m.BlockFilter == nil {
		return "", nil
	}
	switch q := m.BlockFilter.Query.(type) {
	case *Module_BlockFilter_QueryString:
		return q.QueryString, nil
	case *Module_BlockFilter_QueryFromParams:
		for _, input := range m.Inputs {
			if p := input.GetParams(); p != nil {
				return p.Value, nil
			}
		}
		return "", fmt.Errorf("getting blockFilterQueryString: no params input")
	default:
		return "", fmt.Errorf("getting blockFilterQueryString: unsupported query type")
	}
}

func (x *Module) ModuleKind() ModuleKind {
	if x.Kind == nil {
		// Log the module to a file when kind is nil
		logFile := fmt.Sprintf("manifest_%d.log", time.Now().UnixNano())
		absPath, _ := os.Getwd()
		fullPath := fmt.Sprintf("%s/%s", absPath, logFile)
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			defer f.Close()
			f.WriteString(fmt.Sprintf("Module with nil kind: %s\n", x.String()))
		}
		panic(fmt.Sprintf("nil kind for module %q (logged to %s)", x.Name, fullPath))
	}

	switch x.Kind.(type) {
	case *Module_KindMap_:
		return ModuleKindMap
	case *Module_KindStore_:
		return ModuleKindStore
	case *Module_KindBlockIndex_:
		return ModuleKindBlockIndex
	}

	// Log the module to a file when kind is unknown
	logFile := fmt.Sprintf("manifest_%d.log", time.Now().UnixNano())
	absPath, _ := os.Getwd()
	fullPath := fmt.Sprintf("%s/%s", absPath, logFile)
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		defer f.Close()
		f.WriteString(fmt.Sprintf("Module with unknown kind: %s\n", x.String()))
	}
	panic(fmt.Sprintf("unsupported kind: %T (logged to %s)", x.Kind, fullPath))
}

func (x *Module_Input) Pretty() string {
	var result string
	switch x.Input.(type) {
	case *Module_Input_Map_:
		result = x.GetMap().GetModuleName()
	case *Module_Input_Store_:
		result = x.GetStore().GetModuleName()
	case *Module_Input_Source_:
		result = x.GetSource().GetType()
	case *Module_Input_Params_:
		result = x.GetParams().GetValue()
	default:
		result = "unknown"
	}

	return strings.TrimSpace(result)
}

func (x Module_KindStore_UpdatePolicy) Pretty() string {
	return strings.TrimPrefix(strings.ToLower(x.String()), "update_policy_")
}

func (x Module_Input_Store_Mode) Pretty() string {
	return strings.ToLower(Module_Input_Store_Mode_name[int32(x)])
}
