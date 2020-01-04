package main

import(
	"time"
)

type TemplateExecutorEntry struct {
	Tmpl Template
	Prefix string
	Data interface {}
	Content chan string
}

type TemplateExecutor struct {
	Entries []TemplateExecutorEntry

	Init bool
}

var instance *TemplateExecutor

func NewTemplateInstance() *TemplateExecutor {
	if instance == nil {
		instance = new(TemplateExecutor)
		instance.Init()
	}

	return instance
}

func (t *TemplateExecutor) Init() {
	if t.Init {
		return
	}

	go func() {
		for {
			t.Exec()
			time.Sleep(1 * time.Second)
		}
	}()

	t.Init = true
}

func (t *TemplateExecutor) Register(tmpl Template, data interface{}) chan string {
	channel := make(chan string)
	entry := TemplateExecutorEntry {
		Tmpl: tmpl,
		Prefix: tmpl.Prefix,
		Data: data,
		Content: channel,
	}

	append(t.Entries, entry)

	return channel
}

func (t *TemplateExecutor) Exec() {

	for _, e := range t.Entries {
		if e.Tmpl.DoExec() {
			content, err := e.Tmpl.Exec(e.Prefix, e.Data)

			if(err == nil) {
				e.Content <- content
			}
		}
	}
}
