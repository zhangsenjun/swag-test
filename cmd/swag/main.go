package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/swaggo/swag"
	"github.com/swaggo/swag/format"
	"github.com/zhangsenjun/swag-test/gen"
)

const (
	searchDirFlag         = "dir"
	excludeFlag           = "exclude"
	generalInfoFlag       = "generalInfo"
	propertyStrategyFlag  = "propertyStrategy"
	outputFlag            = "output"
	outputTypesFlag       = "outputTypes"
	parseVendorFlag       = "parseVendor"
	parseDependencyFlag   = "parseDependency"
	markdownFilesFlag     = "markdownFiles"
	codeExampleFilesFlag  = "codeExampleFiles"
	parseInternalFlag     = "parseInternal"
	generatedTimeFlag     = "generatedTime"
	requiredByDefaultFlag = "requiredByDefault"
	parseDepthFlag        = "parseDepth"
	instanceNameFlag      = "instanceName"
	overridesFileFlag     = "overridesFile"
	parseGoListFlag       = "parseGoList"
	quietFlag             = "quiet"
	extandFilesFlag       = "extandFiles"
)

var initFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:    quietFlag,
		Aliases: []string{"q"},
		Usage:   "Make the logger quiet.",
	},
	&cli.StringFlag{
		Name:    generalInfoFlag,
		Aliases: []string{"g"},
		Value:   "main.go",
		Usage:   "Go file path in which 'swagger general API Info' is written",
	},
	&cli.StringFlag{
		Name:    searchDirFlag,
		Aliases: []string{"d"},
		Value:   "./",
		Usage:   "Directories you want to parse,comma separated and general-info file must be in the first one",
	},
	&cli.StringFlag{
		Name:  excludeFlag,
		Usage: "Exclude directories and files when searching, comma separated",
	},
	&cli.StringFlag{
		Name:    propertyStrategyFlag,
		Aliases: []string{"p"},
		Value:   swag.CamelCase,
		Usage:   "Property Naming Strategy like " + swag.SnakeCase + "," + swag.CamelCase + "," + swag.PascalCase,
	},
	&cli.StringFlag{
		Name:    outputFlag,
		Aliases: []string{"o"},
		Value:   "./docs",
		Usage:   "Output directory for all the generated files(swagger.json, swagger.yaml and docs.go)",
	},
	&cli.StringFlag{
		Name:    outputTypesFlag,
		Aliases: []string{"ot"},
		Value:   "go,json,yaml",
		Usage:   "Output types of generated files (docs.go, swagger.json, swagger.yaml) like go,json,yaml",
	},
	&cli.BoolFlag{
		Name:  parseVendorFlag,
		Usage: "Parse go files in 'vendor' folder, disabled by default",
	},
	&cli.BoolFlag{
		Name:    parseDependencyFlag,
		Aliases: []string{"pd"},
		Usage:   "Parse go files inside dependency folder, disabled by default",
	},
	&cli.StringFlag{
		Name:    markdownFilesFlag,
		Aliases: []string{"md"},
		Value:   "",
		Usage:   "Parse folder containing markdown files to use as description, disabled by default",
	},
	&cli.StringFlag{
		Name:    codeExampleFilesFlag,
		Aliases: []string{"cef"},
		Value:   "",
		Usage:   "Parse folder containing code example files to use for the x-codeSamples extension, disabled by default",
	},
	&cli.BoolFlag{
		Name:  parseInternalFlag,
		Usage: "Parse go files in internal packages, disabled by default",
	},
	&cli.BoolFlag{
		Name:  generatedTimeFlag,
		Usage: "Generate timestamp at the top of docs.go, disabled by default",
	},
	&cli.IntFlag{
		Name:  parseDepthFlag,
		Value: 100,
		Usage: "Dependency parse depth",
	},
	&cli.BoolFlag{
		Name:  requiredByDefaultFlag,
		Usage: "Set validation required for all fields by default",
	},
	&cli.StringFlag{
		Name:  instanceNameFlag,
		Value: "",
		Usage: "This parameter can be used to name different swagger document instances. It is optional.",
	},
	&cli.StringFlag{
		Name:  overridesFileFlag,
		Value: gen.DefaultOverridesFile,
		Usage: "File to read global type overrides from.",
	},
	&cli.BoolFlag{
		Name:  parseGoListFlag,
		Value: true,
		Usage: "Parse dependency via 'go list'",
	},
}

var updateFlags = append(
	[]cli.Flag{
		&cli.StringFlag{
			Name:    extandFilesFlag,
			Value:   "./docs/common/extands.json",
			Aliases: []string{"efs"},
			Usage:   "Use of multiple files `|` Split. Defaults path is: ./docs/common/extands.json ",
		},
	},
	initFlags...,
)

func initAction(ctx *cli.Context) error {
	strategy := ctx.String(propertyStrategyFlag)

	switch strategy {
	case swag.CamelCase, swag.SnakeCase, swag.PascalCase:
	default:
		return fmt.Errorf("not supported %s propertyStrategy", strategy)
	}

	outputTypes := strings.Split(ctx.String(outputTypesFlag), ",")
	if len(outputTypes) == 0 {
		return fmt.Errorf("no output types specified")
	}
	var logger swag.Debugger
	if ctx.Bool(quietFlag) {
		logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}

	return gen.New().Build(&gen.Config{
		SearchDir:           ctx.String(searchDirFlag),
		Excludes:            ctx.String(excludeFlag),
		MainAPIFile:         ctx.String(generalInfoFlag),
		PropNamingStrategy:  strategy,
		OutputDir:           ctx.String(outputFlag),
		OutputTypes:         outputTypes,
		ParseVendor:         ctx.Bool(parseVendorFlag),
		ParseDependency:     ctx.Bool(parseDependencyFlag),
		MarkdownFilesDir:    ctx.String(markdownFilesFlag),
		ParseInternal:       ctx.Bool(parseInternalFlag),
		GeneratedTime:       ctx.Bool(generatedTimeFlag),
		RequiredByDefault:   ctx.Bool(requiredByDefaultFlag),
		CodeExampleFilesDir: ctx.String(codeExampleFilesFlag),
		ParseDepth:          ctx.Int(parseDepthFlag),
		InstanceName:        ctx.String(instanceNameFlag),
		OverridesFile:       ctx.String(overridesFileFlag),
		ParseGoList:         ctx.Bool(parseGoListFlag),
		Debugger:            logger,
	})
}

func updateAction(ctx *cli.Context) error {
	err := initAction(ctx)
	if err != nil {
		log.Println("Failed to execute `swag init`, update failed")
	}
	log.Println("swag init executed successfully")
	log.Println(ctx.String(extandFilesFlag))
	err = updateData(ctx.String(outputFlag), ctx.String(extandFilesFlag))
	if err != nil {
		return err
	}

	log.Println("update success")
	return nil
}

func main() {
	app := cli.NewApp()
	app.Version = swag.Version
	app.Usage = "Automatically generate RESTful API documentation with Swagger 2.0 for Go."
	app.Commands = []*cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Create docs.go",
			Action:  initAction,
			Flags:   initFlags,
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "update docs",
			Action:  updateAction,
			Flags:   updateFlags,
		},
		{
			Name:    "fmt",
			Aliases: []string{"f"},
			Usage:   "format swag comments",
			Action: func(c *cli.Context) error {
				searchDir := c.String(searchDirFlag)
				excludeDir := c.String(excludeFlag)
				mainFile := c.String(generalInfoFlag)

				return format.New().Build(&format.Config{
					SearchDir: searchDir,
					Excludes:  excludeDir,
					MainFile:  mainFile,
				})
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    searchDirFlag,
					Aliases: []string{"d"},
					Value:   "./",
					Usage:   "Directories you want to parse,comma separated and general-info file must be in the first one",
				},
				&cli.StringFlag{
					Name:  excludeFlag,
					Usage: "Exclude directories and files when searching, comma separated",
				},
				&cli.StringFlag{
					Name:    generalInfoFlag,
					Aliases: []string{"g"},
					Value:   "main.go",
					Usage:   "Go file path in which 'swagger general API Info' is written",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func updateData(docsDirPath string, extandFilesPath string) error {
	const tempPlaceholder string = `"schemes": "Placeholder",`
	const illegalStr string = `"schemes": {{ marshal .Schemes }},`
	var filePath string = docsDirPath + "/docs.go"

	docsFileBytes, err := ioutil.ReadFile(filePath)
	docsFileStr := string(docsFileBytes)
	if err != nil {
		log.Println("open file fail")
		return err
	}

	//The index is increased by one to only take the content in {}, which conforms to the json syntax
	var constBegin int = strings.Index(docsFileStr, "`{") + 1
	var constEnd int = strings.Index(docsFileStr, "}`") + 1

	docTemplate := docsFileStr[constBegin:constEnd]
	jsonTempStr := strings.Replace(docTemplate, illegalStr, tempPlaceholder, 1)

	templateMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(jsonTempStr), &templateMap)
	if err != nil {
		log.Println("json unmarshal fail, check constants in docs.go")
		return err
	}
	definitionsMap := templateMap["definitions"].(map[string]interface{})
	for _, extandsFilePath := range strings.Split(extandFilesPath, "|") {
		err = appendDefinitions(definitionsMap, extandsFilePath)
		if err != nil {
			log.Println("append definitions fail, the file name: " + extandsFilePath)
			return err
		}
	}
	err = replaceType(templateMap["paths"].(map[string]interface{}), definitionsMap)
	if err != nil {
		return err
	}

	tempBytes, err := json.MarshalIndent(templateMap, "", "\t")
	if err != nil {
		log.Println("json marshalIndent fail")
		return err
	}
	jsonTempStr = string(tempBytes)
	jsonTempStr = strings.Replace(jsonTempStr, tempPlaceholder, illegalStr, 1)
	newDocsFile := strings.Replace(docsFileStr, docTemplate, jsonTempStr, 1)

	err = ioutil.WriteFile(filePath, []byte(newDocsFile), 0666)
	if err != nil {
		log.Println("write file fail")
		return err
	}
	return nil
}

func replaceType(pathsMap map[string]interface{}, definitionsMap map[string]interface{}) error {
	const typeRefPrefix = "#/definitions/"
	for path, pathMap := range pathsMap {
		postMap := pathMap.(map[string]interface{})["post"]
		responsesMap := postMap.(map[string]interface{})["responses"].(map[string]interface{})
		for _, content := range responsesMap {
			description := strings.Trim(content.(map[string]interface{})["description"].(string), " ")
			if strings.Index(description, "third-lib-") != -1 {
				schema := content.(map[string]interface{})["schema"].(map[string]interface{})
				descStrs := strings.Split(description, "-")
				if len(descStrs) < 3 {
					log.Println("The interface third-lib data structure description format is incorrect.\nThe error interface is " + path)
					return errors.New("update fail")
				}
				if _, exit := definitionsMap[descStrs[2]]; !exit {
					log.Println(descStrs[2] + " type not exist. Please add type to the extands.json")
					return errors.New("update fail")
				}
				replaceType := typeRefPrefix + descStrs[2]
				schema["$ref"] = replaceType
				delete(schema, "type")
			}
		}
	}
	return nil
}

func appendDefinitions(definitionsMap map[string]interface{}, extandsFilePath string) error {
	extantsFileBytes, err := ioutil.ReadFile(extandsFilePath)
	if err != nil {
		log.Println("open extands.json file fail")
		return err
	}
	var extandsMap map[string]interface{}
	err = json.Unmarshal(extantsFileBytes, &extandsMap)
	if err != nil {
		log.Println("extants.json unmarshal fail")
		return err
	}
	for modelType, content := range extandsMap {
		definitionsMap[modelType] = content
	}
	return nil
}
