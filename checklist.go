package main

type checklist struct {
	status  bool
	comment string
}

type checklists struct {
	packageJsonExists                 checklist
	mainJsExists                      checklist
	mainJsHaveStudentIdComment        checklist
	rootShowingHtml                   checklist
	serveInPort5000                   checklist
	htmlContainH1ElementWithStudentId checklist
}

func newChecklists() checklists {
	defaultChecklist := checklist{
		status:  false,
		comment: "",
	}
	return checklists{
		packageJsonExists:                 defaultChecklist,
		mainJsExists:                      defaultChecklist,
		mainJsHaveStudentIdComment:        defaultChecklist,
		rootShowingHtml:                   defaultChecklist,
		serveInPort5000:                   defaultChecklist,
		htmlContainH1ElementWithStudentId: defaultChecklist,
	}
}
