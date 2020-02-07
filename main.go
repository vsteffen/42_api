package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vsteffen/42_api/reqApi42"
	_ "github.com/vsteffen/42_api/tools"
	cst "github.com/vsteffen/42_api/tools/constants"
	"os"
	"time"
	"fmt"
	"bufio"
	"regexp"
	"github.com/manifoldco/promptui"
)

type projectParent struct {
	this	*reqApi42.API42ProjectParent
	sons	[]*reqApi42.API42Project
}

type projectsPerType struct {
	parents		map[uint]*projectParent
	directs		[]*reqApi42.API42Project
}

func askStringClean (askStr string) (string) {
	fmt.Print(askStr)
	scannerStdin := bufio.NewScanner(os.Stdin)
	scannerStdin.Scan()

	if err := scannerStdin.Err(); err != nil {
		log.Fatal().Err(err).Msg("askString: Failed to read user input")
	}
	regexWhitespace := regexp.MustCompile(`\s+`)
	str := regexWhitespace.ReplaceAllString(scannerStdin.Text(), " ")
	return str
}

func findExaminer(api42 *reqApi42.API42) {
	prompt := promptui.Select{
		Label:	"Does your project have a parent?",
		Items:	[]string{"Yes", "No"},
		HideHelp: true,
	}
	indexAction, input, err := prompt.Run()

	if err != nil {
		log.Fatal().Err(err).Msg("PromptUI: failed")
	}

	var parentProjectName string
	if indexAction == 0 {
		parentProjectName = askStringClean("Please, enter the parent project name: ")
	} else {
		fmt.Println("Not a parent")
	}
	_ = parentProjectName
	fmt.Println(input)
}

func sortProjectsPerType(api42Projects *[]reqApi42.API42Project) (*projectsPerType) {
	var allProjects projectsPerType
	allProjects.parents = make(map[uint]*projectParent)
	allProjects.directs = make([]*reqApi42.API42Project, 0)
	for index, project := range *api42Projects {
		if project.Parent == nil {
			allProjects.directs = append(allProjects.directs, &(*api42Projects)[index])
		} else {
			projectDeref := (*api42Projects)[index]
			if parentMapValue, ok := allProjects.parents[projectDeref.Parent.ID]; !ok {
				allProjects.parents[projectDeref.Parent.ID] = &projectParent{projectDeref.Parent, []*reqApi42.API42Project{&(*api42Projects)[index]}}
			} else {
				parentMapValue.sons = append(parentMapValue.sons, &(*api42Projects)[index])
			}
		}
	}
	return &allProjects
}

func debugPrintProjectsPerType(allProjects *projectsPerType) {
	fmt.Println("###################################")
	for _, parent := range allProjects.parents {
		fmt.Println(parent.this.Name)
		for _, son := range parent.sons {
			fmt.Println("-> " + son.Name)
		}
		fmt.Println("----------------")
	}
	fmt.Println("+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+")
	for _, direct := range allProjects.directs {
		fmt.Println(direct.Name)
	}
	fmt.Println("###################################")
}

func main() {
	flags := []interface{}{}
	flags = append(flags, flag.Bool("refresh", false, "force to refresh token"))
	flags = append(flags, flag.Bool("check-default-values", false, "send a request to verify the default values"))
	flag.Parse()
	nonFlags := flag.Args()
	if len(nonFlags) > 0 {
		flag.Usage()
		os.Exit(1)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Stamp})

	fmt.Print(cst.MenuHello)

	api42 := reqApi42.New(flags)
	allProjects := sortProjectsPerType(api42.GetProjects())
	debugPrintProjectsPerType(allProjects)
	// os.Exit(0)
	// _ = allProjects
	var indexAction int
	var err error
	menuActions := []string{cst.MenuActionFind, cst.MenuActionUpdateLocations, cst.MenuActionUpdateProjects, cst.MenuActionUpdateCursus, cst.MenuActionUpdateCampus, cst.MenuActionQuit}
	for {
		prompt := promptui.Select{
			Label:	"Choose an action",
			Items:	menuActions,
			HideHelp: true,
		}

		indexAction, _, err = prompt.Run()

		if err != nil {
			log.Fatal().Err(err).Msg("PromptUI: failed")
		}

		switch menuActions[indexAction] {
		case cst.MenuActionFind:
			findExaminer(api42)
		case cst.MenuActionUpdateLocations:
			api42.UpdateLocations()
		case cst.MenuActionUpdateProjects:
			api42.UpdateProjects()
			allProjects = sortProjectsPerType(api42.GetProjects())
		case cst.MenuActionUpdateCursus:
			cursusName := askStringClean("Please, enter the cursus name: ")
			api42.UpdateCursus(cursusName)
		case cst.MenuActionUpdateCampus:
			campusName := askStringClean("Please, enter the campus name: ")
			api42.UpdateCampus(campusName)
		case cst.MenuActionQuit:
			fmt.Println("Goodbye!")
			os.Exit(0)
		default:
			log.Fatal().Msg("PromptUI: indexAction out of bound")
		}

	}

	// fmt.Printf("You choose %s\n", input)





	// fmt.Println(api42.GetProjects())
}


/*
--> Update locations
--> Update campus
--> Update cursus
--> Find project examiner
	---> Project parent ?
		---> Search your project
			---> Show n results or matching project
--> Exit
*/
