package react

import (
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/bendigiorgio/ikou/internal/app/utils"
	esbuild "github.com/evanw/esbuild/pkg/api"
	"go.uber.org/zap"
	v8 "rogchap.com/v8go"
)

// [Yaffle/TextEncoderTextDecoder.js](https://gist.github.com/Yaffle/5458286)
var textEncoderPolyfill = `function TextEncoder(){} TextEncoder.prototype.encode=function(string){var octets=[],length=string.length,i=0;while(i<length){var codePoint=string.codePointAt(i),c=0,bits=0;codePoint<=0x7F?(c=0,bits=0x00):codePoint<=0x7FF?(c=6,bits=0xC0):codePoint<=0xFFFF?(c=12,bits=0xE0):codePoint<=0x1FFFFF&&(c=18,bits=0xF0),octets.push(bits|(codePoint>>c)),c-=6;while(c>=0){octets.push(0x80|((codePoint>>c)&0x3F)),c-=6}i+=codePoint>=0x10000?2:1}return octets};function TextDecoder(){} TextDecoder.prototype.decode=function(octets){var string="",i=0;while(i<octets.length){var octet=octets[i],bytesNeeded=0,codePoint=0;octet<=0x7F?(bytesNeeded=0,codePoint=octet&0xFF):octet<=0xDF?(bytesNeeded=1,codePoint=octet&0x1F):octet<=0xEF?(bytesNeeded=2,codePoint=octet&0x0F):octet<=0xF4&&(bytesNeeded=3,codePoint=octet&0x07),octets.length-i-bytesNeeded>0?function(){for(var k=0;k<bytesNeeded;){octet=octets[i+k+1],codePoint=(codePoint<<6)|(octet&0x3F),k+=1}}():codePoint=0xFFFD,bytesNeeded=octets.length-i,string+=String.fromCodePoint(codePoint),i+=bytesNeeded+1}return string};`
var processPolyfill = `var process = {env: {NODE_ENV: "production"}};`
var consolePolyfill = `var console = {log: function(){}};`

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>React App</title>
	<link href="/public/style.css" rel="stylesheet">
</head>
<body>
    <div id="app">{{.RenderedContent}}</div>
	<script type="module">
	  {{ .JS }}
	</script>
	<script id="IKOU_PROPS">window.APP_PROPS = {{.InitialProps}};</script>
</body>
</html>
`

const clientHtmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>React App</title>
	<link href="/public/style.css" rel="stylesheet">
</head>
<body>
    <div id="app">{{.RenderedContent}}</div>
	<script type="module">
	  {{ .JS }}
	renderClientSide();
	</script>
	<script id="IKOU_PROPS">window.APP_PROPS = {{.InitialProps}};</script>
</body>
</html>
`

type PageData struct {
	RenderedContent template.HTML
	InitialProps    template.JS
	JS              template.JS
	Tmpl            *template.Template
}

type PageProps struct {
	PageRoute string
}

// buildBackend compiles the specified TypeScript or TSX file into a single JavaScript bundle using esbuild.
// The resulting bundle is formatted as an Immediately Invoked Function Expression (IIFE) for use in v8.
//
// Parameters:
//   - pagePath: The file path of the TypeScript or TSX entry point to be bundled.
//
// Returns:
//   - A string containing the contents of the bundled JavaScript file.
//   - An error if the build process fails or if no output files are generated.
func buildBackend(pagePath string) (string, error) {
	defer utils.Logger.Sync()
	result := esbuild.Build(esbuild.BuildOptions{
		EntryPoints:       []string{pagePath},
		Bundle:            true,
		Write:             false,
		Outdir:            "out/",
		Format:            esbuild.FormatIIFE, // IIFE format for use in v8
		Platform:          esbuild.PlatformBrowser,
		Target:            esbuild.ESNext,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Metafile:          false,
		LogLevel:          esbuild.LogLevelError,

		Banner: map[string]string{
			"js": textEncoderPolyfill + processPolyfill + consolePolyfill,
		},
		Loader: map[string]esbuild.Loader{
			".tsx": esbuild.LoaderTSX,
			".ts":  esbuild.LoaderTS,
		},
	})

	if len(result.OutputFiles) == 0 {
		utils.Logger.Sugar().Fatal("Server build error:", result.Errors)
		return "", fmt.Errorf("no output files from backend build")
	}

	return string(result.OutputFiles[0].Contents), nil
}

// buildClient takes a client entry point file path, uses esbuild to bundle it,
// and returns the bundled client-side JavaScript as a string. If the build
// process does not produce any output files, it returns an error.
//
// Parameters:
//   - clientEntry: A string representing the file path of the client entry point.
//
// Returns:
//   - A string containing the bundled client-side JavaScript.
//   - An error if the build process fails or produces no output files.
func buildClient(clientEntry string) (string, error) {
	defer utils.Logger.Sync()
	clientResult := esbuild.Build(esbuild.BuildOptions{
		EntryPoints: []string{clientEntry},
		Bundle:      true,
		Write:       true,
		LogLevel:    esbuild.LogLevelError,
	})

	if len(clientResult.OutputFiles) == 0 {
		utils.Logger.Sugar().Fatal("Client build error:", clientResult.Errors)
		return "", fmt.Errorf("no output files from client build")
	}

	return string(clientResult.OutputFiles[0].Contents), nil
}

// RenderPage renders a React page either as a static site generation (SSG) or server-side rendering (SSR).
//
// Parameters:
// - isSSG: A boolean indicating if the page should be rendered as SSG.
// - clientEntry: The entry point for the client-side bundle.
// - props: The properties to be passed to the React component.
// - pagePath: The path of the page to be rendered.
//
// Returns:
// - PageData: A struct containing the rendered HTML content, initial props, JavaScript bundle, and the HTML template.
// - error: An error if any occurred during the rendering process.
func RenderPage(isSSG bool, clientEntry string, props PageProps, pagePath string) (PageData, error) {
	defer utils.Logger.Sync()

	var renderedHTML string
	var err error

	propsWithPage := struct {
		PageProps
		PagePath string `json:"pagePath"`
	}{
		PageProps: props,
		PagePath:  pagePath,
	}

	clientBundle := ""

	backendBundle, err := buildBackend(pagePath)
	if err != nil {
		utils.Logger.Error("Error building backend bundle", zap.Error(err))
		return PageData{}, err
	}

	ctx := v8.NewContext(nil)

	_, err = ctx.RunScript(backendBundle, "bundle.js")
	if err != nil {
		utils.Logger.Error("Error running backend bundle", zap.Error(err))
		return PageData{}, err
	}

	val, err := ctx.RunScript("renderApp()", "render.js")
	if err != nil {
		utils.Logger.Sugar().Fatalf("Failed to render React component: %v", err)
	}
	renderedHTML = val.String()

	if !isSSG {
		clientBundle, err = buildClient(clientEntry)
		if err != nil {
			utils.Logger.Error("Error building client bundle", zap.Error(err))
			return PageData{}, err
		}
	}

	jsonProps, err := json.Marshal(propsWithPage)
	if err != nil {
		utils.Logger.Sugar().Fatalf("Failed to marshal props: %v", err)
		return PageData{}, err
	}

	tmpl, err := template.New("webpage").Parse(htmlTemplate)
	if err != nil {
		utils.Logger.Sugar().Fatal("Error parsing template:", err)
		return PageData{}, err
	}

	if !isSSG {
		tmpl, err = template.New("webpage").Parse(clientHtmlTemplate)
		if err != nil {
			utils.Logger.Sugar().Fatal("Error parsing client template:", err)
			return PageData{}, err
		}
	}

	return PageData{
		RenderedContent: template.HTML(renderedHTML),
		InitialProps:    template.JS(jsonProps),
		JS:              template.JS(clientBundle),
		Tmpl:            tmpl,
	}, nil
}
