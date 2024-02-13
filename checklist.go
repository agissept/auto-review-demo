package main

type checklist struct {
	status  bool
	comment string
}

//goland:noinspection ALL
type checklists struct {
	packageJsonExists                 checklist `json:"package_json_exists"`
	mainJsExists                      checklist `json:"main_js_exists"`
	mainJsHaveStudentIdComment        checklist `json:"main_js_have_student_id_comment"`
	rootShowingHtml                   checklist `json:"root_showing_html"`
	serveInPort5000                   checklist `json:"serve_in_port_5000"`
	htmlContainH1ElementWithStudentId checklist `json:"html_contain_h_1_element_with_student_id"`
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
