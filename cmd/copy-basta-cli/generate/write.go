package generate

import "text/template"

func write(root string, files []file) error {
	for _, file := range files {
		if file.template {
			t, err := newTemplate(file.path, string(file.content))
			if err != nil {
				return err
			}
			panic("deal with t")
		}
		panic("deal with creating new file")
	}
	return nil
}

func newTemplate(name string, tmpl string) (*template.Template, error) {
	return template.New(name).Parse(tmpl)
}
