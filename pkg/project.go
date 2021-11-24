package pkg

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/piyushjajoo/crud/tpl"
)

type Project struct {
	ModuleName      string
	ProjectDirName  string
	AbsolutePath    string
	CreateApiDoc    bool
	CreateHelmChart bool
}

const (
	GorillaMuxModuleName = "github.com/gorilla/mux"
	EnvConfigModuleName  = "github.com/kelseyhightower/envconfig"
	ValidatorModuleName  = "github.com/go-playground/validator"
)

func (p *Project) Create() error {

	cwd, err := os.Getwd()
	if err != nil {
		log.Println("error getting current working directory:", err)
		return err
	}

	// create the project directory
	if err = createDir(p.AbsolutePath); err != nil {
		log.Println("error creating project directory at path", p.AbsolutePath, ":", err)
		return err
	}

	// initialize go module
	if err := os.Chdir(p.AbsolutePath); err != nil {
		log.Println("error changing directory to path", p.AbsolutePath, ":", err)
		return err
	} else {
		// initialize go module
		if err = goMod(p.ModuleName); err != nil {
			log.Println("error initializing go module", p.ModuleName, "at path", p.AbsolutePath, ":", err)
			return err
		}

		// go get gorilla mux
		if err := goGet(GorillaMuxModuleName); err != nil {
			log.Println("error getting module", GorillaMuxModuleName, ":", err)
			return err
		}

		// go get github.com/kelseyhightower/envconfig
		if err := goGet(EnvConfigModuleName); err != nil {
			log.Println("error getting module", EnvConfigModuleName, ":", err)
			return err
		}

		// go get "github.com/go-playground/validator"
		if err := goGet(ValidatorModuleName); err != nil {
			log.Println("error getting module", ValidatorModuleName, ":", err)
			return err
		}

		// change the directory to cwd
		err = os.Chdir(cwd)
		if err != nil {
			log.Println("error changing current working directory to", cwd, "after creating project directory:", err)
			return err
		}
	}

	// create main.go
	mainFile, err := os.Create(fmt.Sprintf("%s/main.go", p.AbsolutePath))
	if err != nil {
		log.Println("error creating main.go at", p.AbsolutePath, ":", err)
		return err
	}
	defer mainFile.Close()

	mainTemplate := template.Must(template.New("main").Parse(string(tpl.MainTemplate())))
	err = mainTemplate.Execute(mainFile, p)
	if err != nil {
		log.Println("error executing template for main.go", err)
		return err
	}

	// create routes, consts, conf, models and utils directories under pkg directory
	pkgDir := p.AbsolutePath + "/pkg"
	if err = createDir(pkgDir); err != nil {
		log.Println("error creating pkg directory at", p.AbsolutePath, ":", err)
		return err
	}

	routesDir := pkgDir + "/routes"
	if err = createDir(routesDir); err != nil {
		log.Println("error creating routes directory at", pkgDir, ":", err)
		return err
	} else {
		routesFile, err := os.Create(fmt.Sprintf("%s/routes.go", routesDir))
		if err != nil {
			log.Println("error creating routes.go at", routesDir, ":", err)
			return err
		}
		defer routesFile.Close()

		routesTemplate := template.Must(template.New("routes").Parse(string(tpl.RoutesTemplate())))
		err = routesTemplate.Execute(routesFile, p)
		if err != nil {
			log.Println("error executing template for routes.go", err)
			return err
		}
	}

	confDir := pkgDir + "/conf"
	if err = createDir(confDir); err != nil {
		log.Println("error creating conf directory at", pkgDir, ":", err)
		return err
	} else {
		confFile, err := os.Create(fmt.Sprintf("%s/conf.go", confDir))
		if err != nil {
			log.Println("error creating conf.go at", confDir, ":", err)
			return err
		}
		defer confFile.Close()

		confTemplate := template.Must(template.New("conf").Parse(string(tpl.ConfTemplate())))
		err = confTemplate.Execute(confFile, p)
		if err != nil {
			log.Println("error executing template for conf.go", err)
			return err
		}
	}

	modelsDir := pkgDir + "/models"
	if err = createDir(modelsDir); err != nil {
		log.Println("error creating models directory at", modelsDir, ":", err)
		return err
	} else {
		modelsFile, err := os.Create(fmt.Sprintf("%s/models.go", modelsDir))
		if err != nil {
			log.Println("error creating models.go at", modelsDir, ":", err)
			return err
		}
		defer modelsFile.Close()

		modelsTemplate := template.Must(template.New("models").Parse(string(tpl.ModelsTemplate())))
		err = modelsTemplate.Execute(modelsFile, p)
		if err != nil {
			log.Println("error executing template for models.go", err)
			return err
		}
	}

	utilsDir := pkgDir + "/utils"
	if err = createDir(utilsDir); err != nil {
		log.Println("error creating utils directory at", utilsDir, ":", err)
		return err
	} else {
		utilsFile, err := os.Create(fmt.Sprintf("%s/utils.go", utilsDir))
		if err != nil {
			log.Println("error creating utils.go at", utilsDir, ":", err)
			return err
		}
		defer utilsFile.Close()

		utilsTemplate := template.Must(template.New("utils").Parse(string(tpl.UtilsTemplate())))
		err = utilsTemplate.Execute(utilsFile, p)
		if err != nil {
			log.Println("error executing template for utils.go", err)
			return err
		}
	}

	// if api flag is set, create api documentation
	if p.CreateApiDoc {
		apiDir := p.AbsolutePath + "/api"
		if err = createDir(apiDir); err != nil {
			log.Println("error creating api directory at", p.AbsolutePath, ":", err)
			return err
		}
		swaggerFile, err := os.Create(fmt.Sprintf("%s/swagger.json", apiDir))
		if err != nil {
			log.Println("error creating swagger.json at", apiDir, ":", err)
			return err
		}
		defer swaggerFile.Close()
	}

	// if helm flag is set, create helm chart
	if p.CreateHelmChart {
		if err := os.Chdir(p.AbsolutePath); err != nil {
			log.Println("error changing directory to path", p.AbsolutePath, ":", err)
			return err
		} else {
			chartsDir := p.AbsolutePath+"/charts"
			if err = createDir(chartsDir); err != nil {
				log.Println("error creating charts directory at", p.AbsolutePath, ":", err)
				return err
			}
			err = exec.Command("helm", "create", chartsDir+"/"+p.ProjectDirName).Run()
			if err != nil {
				log.Println("error creating helm chart at", chartsDir, ":", err)
				return err
			}
			// change the directory to cwd
			err = os.Chdir(cwd)
			if err != nil {
				log.Println("error changing current working directory to", cwd, "after creating helm chart:", err)
				return err
			}
		}
	}

	// create README.md
	readmeFile, err := os.Create(fmt.Sprintf("%s/README.md", p.AbsolutePath))
	if err != nil {
		log.Println("error creating README.md file at", p.AbsolutePath, ":", err)
		return err
	}
	defer readmeFile.Close()

	// create Dockerfile
	dockerfile, err := os.Create(fmt.Sprintf("%s/Dockerfile", p.AbsolutePath))
	if err != nil {
		log.Println("error creating Dockerfile file at", p.AbsolutePath, ":", err)
		return err
	}
	defer dockerfile.Close()

	dockerfileTemplate := template.Must(template.New("dockerfile").Parse(string(tpl.DockerfileTemplate())))
	err = dockerfileTemplate.Execute(dockerfile, p)
	if err != nil {
		log.Println("error executing Dockerfile template:", err)
		return err
	}

	// create build.sh to build the docker image
	buildFile, err := os.Create(fmt.Sprintf("%s/build.sh", p.AbsolutePath))
	if err != nil {
		log.Println("error creating build.sh file at", p.AbsolutePath, ":", err)
		return err
	}
	defer buildFile.Close()

	buildFileTemplate := template.Must(template.New("build").Parse(string(tpl.BuildFileTemplate())))
	err = buildFileTemplate.Execute(buildFile, p)
	if err != nil {
		log.Println("error executing build.sh template:", err)
		return err
	}

	// provide executable permissions to build.sh
	err = os.Chmod(p.AbsolutePath+"/build.sh", 0755)
	if err != nil {
		log.Println("error changing permissions for build.sh:", err)
		return err
	}

	return nil
}

// goMod runs the go mod init <module name> command
func goMod(moduleName string) error {
	return exec.Command("go", "mod", "init", moduleName).Run()
}

// goGet runs the go get <module>
func goGet(moduleName string) error {
	return exec.Command("go", "get", moduleName).Run()
}

// createDir creates a directory if it doesn't exist at the provided absolute path
func createDir(absolutePath string) error {
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		// create project directory with permissions for User (all), Group (read and execute) and Others (Read)
		if err = os.Mkdir(absolutePath, 0754); err != nil {
			return err
		}
	}
	return nil
}
