package saunter

import tmpl "html/template"

var indexTemplate *tmpl.Template

func generateIndexTemplate() {
	const source = `
  <!DOCTYPE html>
  <html lang="en">
    <head>
      <meta charset="UTF-8" />
      <title>Swagger UI</title>
      <link rel="stylesheet" type="text/css" href="/swagger-static/swagger-ui.css" />
      <link
        rel="icon"
        type="image/png"
        href="/swagger-static/favicon-32x32.png"
        sizes="32x32"
      />
      <link
        rel="icon"
        type="image/png"
        href="/swagger-static/favicon-16x16.png"
        sizes="16x16"
      />
      <style>
        html {
          box-sizing: border-box;
          overflow: -moz-scrollbars-vertical;
          overflow-y: scroll;
        }
        *,
        *:before,
        *:after {
          box-sizing: inherit;
        }

        body {
          margin: 0;
          background: #fafafa;
        }
      </style>
    </head>
    <body>
      <div id="swagger-ui"></div>
      <script src="/swagger-static/swagger-ui-bundle.js" charset="UTF-8"></script>
      <script>
        window.onload = function () {
          window.ui = SwaggerUIBundle({
            spec: {{.Spec}},
            dom_id: '#swagger-ui',
            deepLinking: true,
            presets: [SwaggerUIBundle.presets.apis],
            defaultModelExpandDepth: 10,
            defaultModelsExpandDepth: -1,
          });
        };
      </script>
      <style>
        .swagger-ui .information-container .info {
          margin: 20px 0;
        }
        .swagger-ui .scheme-container {
          padding: unset;
          background: unset;
          box-shadow: unset;
          margin: -60px 0 0 0;
          padding-bottom: 30px;
        }
        .swagger-ui .download-contents {
          display: none;
        }
        .swagger-ui .copy-to-clipboard {
          bottom: 5px;
          right: 10px;
          width: 20px;
          height: 20px;
        }
        .swagger-ui .copy-to-clipboard button {
          padding-left: 18px;
          height: 18px;
        }
      </style>
    </body>
  </html>
  `
	template, err := tmpl.New("index").Parse(source)
	if err != nil {
		panic(err)
	}

	indexTemplate = template
}
