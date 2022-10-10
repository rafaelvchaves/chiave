package crdt

var graphHandlers = map[string]Handler[Graph]{
	"ADDV": AddVertexHandler{},
	"RMV":  RemoveVertexHandler{},
	"ADDE": AddEdgeHandler{},
	"RME":  RemoveEdgeHandler{},
}

var graphQueries = map[string]Query[Graph]{
	"NEIGH":   NeighborQuery{},
	"EXISTSV": ExistsVertexQuery{},
	"EXISTSE": ExistsEdgeQuery{},
}

var counterHandlers = map[string]Handler[Counter]{
	"INC": IncrementHandler{},
	"DEC": DecrementHandler{},
}

var counterQueries = map[string]Query[Counter]{
	"VALUE": ValueQuery{},
}

var setHandlers = map[string]Handler[ORSet]{
	"ADD": AddHandler{},
	"REM": RemoveHandler{},
}

var setQueries = map[string]Query[ORSet]{
	"EXISTS": ExistsQuery{},
}
