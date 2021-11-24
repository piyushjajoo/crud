/*
Copyright © 2021 Piyush Jajoo piyush.jajoo1991@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/piyushjajoo/crud/pkg"

	"github.com/spf13/cobra"
)

var api, helm bool
var name string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:     "init <module name>",
	Short:   "init creates the scaffolding for the go based micro-service",
	Aliases: []string{"initialize", "initialise", "create"},
	Long: `
Init (cobra init) command initializes the go module along with a bare-bone http-server.
Please make sure you have go installed and GOPATH set. Also make sure you have helm v3 installed as well.

By default if no flags provided it initializes following -
1. go.mod and go.sum files
2. Dockerfile to build your micro-service along with build.sh script
3. main.go with bare http-server written in gorilla mux
4. README.md with basic Summary

If you want api documentation provide --swagger flag. If you want helm chart provide --chart flag.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cobra.CheckErr(fmt.Errorf("init needs the module name"))
		}

		cobra.CheckErr(validateModuleName(args[0])) // validates module name

		projectPath, err := createProject(args) // create project
		cobra.CheckErr(err)
		fmt.Printf("Your micro-service scaffolding is created at\n%s\n", projectPath)
	},
}

// validateModuleName validates the syntax of the module name
func validateModuleName(moduleName string) error {
	// FIXME add implementation
	return nil
}

// createProject initializes the Project object
func createProject(args []string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// parse the project directory name from module name
	projectDirName := getProjectDirName(args[0])
	wd = fmt.Sprintf("%s/%s", wd, projectDirName)

	project := &pkg.Project{
		AbsolutePath:    wd,
		ModuleName:      args[0],
		ProjectDirName:  projectDirName,
		CreateApiDoc:    api,
		CreateHelmChart: helm,
	}

	// create the project
	err = project.Create()
	if err != nil {
		return "", err
	}

	return project.AbsolutePath, nil
}

// getProjectDirName returns the project directory name from the module name
func getProjectDirName(moduleName string) string {
	moduleNameSplit := strings.Split(moduleName, "/")
	return moduleNameSplit[len(moduleNameSplit)-1]
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&name, "name", "n", "", "module name for the go module, last part of the name will be used for directory name (e.g 'github.com/piyushjajoo/crud' is the module name and crud is the directory name)")
	initCmd.Flags().BoolVarP(&api, "swagger", "s", false, "to generate swagger api documentation file")
	initCmd.Flags().BoolVarP(&helm, "chart", "c", false, "to generate helm chart")
}
