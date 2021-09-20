// +build ignore

package main

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-provider-aws/internal/namevaluesfilters"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

const filename = `service_filters_gen.go`

// Representing types such as []*ec2.Filter, []*rds.Filter, ...
var sliceServiceNames = []string{
	"autoscaling",
	"databasemigrationservice",
	"docdb",
	"ec2",
	"elasticinference",
	"elasticsearchservice",
	"fsx",
	"imagebuilder",
	"licensemanager",
	"neptune",
	"rds",
	"resourcegroupstaggingapi",
	"route53resolver",
}

type TemplateData struct {
	SliceServiceNames []string
}

func main() {
	// Always sort to reduce any potential generation churn
	sort.Strings(sliceServiceNames)

	templateData := TemplateData{
		SliceServiceNames: sliceServiceNames,
	}
	templateFuncMap := template.FuncMap{
		"FilterPackage":         namevaluesfilters.ServiceFilterPackage,
		"FilterType":            namevaluesfilters.ServiceFilterType,
		"FilterTypeNameField":   namevaluesfilters.ServiceFilterTypeNameField,
		"FilterTypeValuesField": namevaluesfilters.ServiceFilterTypeValuesField,
		"Title":                 strings.Title,
	}

	tmpl, err := template.New("servicefilters").Funcs(templateFuncMap).Parse(templateBody)

	if err != nil {
		log.Fatalf("error parsing template: %s", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, templateData)

	if err != nil {
		log.Fatalf("error executing template: %s", err)
	}

	generatedFileContents, err := format.Source(buffer.Bytes())

	if err != nil {
		log.Fatalf("error formatting generated file: %s", err)
	}

	f, err := os.Create(filename)

	if err != nil {
		log.Fatalf("error creating file (%s): %s", filename, err)
	}

	defer f.Close()

	_, err = f.Write(generatedFileContents)

	if err != nil {
		log.Fatalf("error writing to file (%s): %s", filename, err)
	}
}

var templateBody = `
// Code generated by generators/servicefilters/main.go; DO NOT EDIT.

package namevaluesfilters

import (
	"github.com/aws/aws-sdk-go/aws"
{{- range .SliceServiceNames }}
{{- if eq . (. | FilterPackage) }}
	"github.com/aws/aws-sdk-go/service/{{ . }}"
{{- end }}
{{- end }}
)

// []*SERVICE.Filter handling
{{- range .SliceServiceNames }}

// {{ . | Title }}Filters returns {{ . }} service filters.
func (filters NameValuesFilters) {{ . | Title }}Filters() []*{{ . | FilterPackage }}.{{ . | FilterType }} {
	m := filters.Map()

	if len(m) == 0 {
		return nil
	}

	result := make([]*{{ . | FilterPackage }}.{{ . | FilterType }}, 0, len(m))

	for k, v := range m {
		filter := &{{ . | FilterPackage }}.{{ . | FilterType }}{
			{{ . | FilterTypeNameField }}:   aws.String(k),
			{{ . | FilterTypeValuesField }}: aws.StringSlice(v),
		}

		result = append(result, filter)
	}

	return result
}
{{- end }}
`
