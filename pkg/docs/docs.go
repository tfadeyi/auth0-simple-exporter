// Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Oluwole Fadeyi (@tfadeyi)"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "https://github.com/tfadeyi/auth0-simple-exporter/blob/main/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/metrics": {
            "get": {
                "description": "Exposes the Auth0 metrics collected by the exporter in a prometheus format.",
                "produces": [
                    "application/json",
                    "text/plain; charset=utf-8"
                ],
                "summary": "Auth0 metrics in Prometheus format.",
                "responses": {}
            }
        },
        "/probe": {
            "get": {
                "description": "Exposes the exporter's own metrics, i.e: target_scrape_request_total.",
                "produces": [
                    "application/json",
                    "text/plain; charset=utf-8"
                ],
                "summary": "Exporter's own metrics for its operations.",
                "responses": {}
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.7",
	Host:             "localhost:9301",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Auth0 simple exporter",
	Description:      "A simple Prometheus exporter for Auth0 log [events](https://auth0.com/docs/api/management/v2#!/Logs/get_logs),\nwhich allows you to collect metrics from Auth0 and expose them in a format that can be consumed by Prometheus.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}