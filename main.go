package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vsteffen/42_api/reqApi42"
	"github.com/vsteffen/42_api/tools"
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

func findProjectName(searchStr string, projects *[]*reqApi42.API42Project) ([]*reqApi42.API42Project, []string, bool) {
	matchProjects := make([]*reqApi42.API42Project, cst.FindNameMaxResults)
	matchCosts := make([]int, cst.FindNameMaxResults)
	highestCost := cst.MaxInt

	for indexInit, _ := range matchCosts {
		matchCosts[indexInit] = cst.MaxInt
	}

	for indexProject, project := range *projects {
		currentCost := tools.EditDistance(searchStr, project.Name)
		if currentCost == 0 {
			matchCosts[0] = currentCost
			matchProjects[0] = (*projects)[indexProject]
			return matchProjects, []string{project.Name}, true
		}
		if currentCost < highestCost {
			for indexMatchCost, cost := range matchCosts {
				if currentCost < cost {
					copy(matchCosts[indexMatchCost + 1:], matchCosts[indexMatchCost:])
					copy(matchProjects[indexMatchCost + 1:], matchProjects[indexMatchCost:])
					matchCosts[indexMatchCost] = currentCost
					matchProjects[indexMatchCost] = (*projects)[indexProject]
					if indexMatchCost + 1 == cst.FindNameMaxResults {
						highestCost = currentCost
					}
					break
				}
			}
		}
	}
	matchStrings := make([]string, 0)
	for _, project := range matchProjects {
		if project == nil {
			break
		}
		matchStrings = append(matchStrings, project.Name)
	}
	return matchProjects, matchStrings, false
}

func findProjectParentName(searchStr string, parents *map[uint]*projectParent) ([]*projectParent, []string, bool) {
	matchParent := make([]*projectParent, cst.FindNameMaxResults)
	matchCosts := make([]int, cst.FindNameMaxResults)
	highestCost := cst.MaxInt

	for indexInit, _ := range matchCosts {
		matchCosts[indexInit] = cst.MaxInt
	}

	for indexProject, project := range *parents {
		currentCost := tools.EditDistance(searchStr, project.this.Name)
		if currentCost == 0 {
			matchCosts[0] = currentCost
			matchParent[0] = (*parents)[indexProject]
			return matchParent, []string{project.this.Name}, true
		}
		if currentCost < highestCost {
			for indexMatchCost, cost := range matchCosts {
				if currentCost < cost {
					copy(matchCosts[indexMatchCost + 1:], matchCosts[indexMatchCost:])
					copy(matchParent[indexMatchCost + 1:], matchParent[indexMatchCost:])
					matchCosts[indexMatchCost] = currentCost
					matchParent[indexMatchCost] = (*parents)[indexProject]
					if indexMatchCost + 1 == cst.FindNameMaxResults {
						highestCost = currentCost
					}
					break
				}
			}
		}
	}
	matchStrings := make([]string, 0)
	for _, parent := range matchParent {
		matchStrings = append(matchStrings, parent.this.Name)
	}
	return matchParent, matchStrings, false
}

func getIndexNameChoice(items []string) (int) {
	items = append(items, "Cancel")
	prompt := promptui.Select{
		Label:	"Found these projects name. Choose or cancel",
		Items:	items,
		HideHelp: true,
	}
	indexProjectFind, _, err := prompt.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("PromptUI: failed")
	}
	if indexProjectFind == cst.FindNameMaxResults {
		return -1
	}
	return indexProjectFind
}

func findExaminer(allProjects *projectsPerType) {
	prompt := promptui.Select{
		Label:	"Does your project have a parent",
		Items:	[]string{"Yes", "No"},
		HideHelp: true,
	}
	indexAction, _, err := prompt.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("PromptUI: failed")
	}

	var realProjectsToSearch *[]*reqApi42.API42Project
	if indexAction == 0 {
		parentProjectName := askStringClean("Please, enter the parent project name: ")
		parentFind, parentsFindNames, fullMatch := findProjectParentName(parentProjectName, &allProjects.parents)
		if fullMatch {
			realProjectsToSearch = &(parentFind[0].sons)
		} else {
			indexChoose := getIndexNameChoice(parentsFindNames)
			if indexChoose == -1 {
				return
			}
			realProjectsToSearch = &(parentFind[indexChoose].sons)
		}
	} else {
		realProjectsToSearch = &allProjects.directs
	}
	projectName := askStringClean("Please, enter the project name: ")
	projectsFind, projectsFindNames, fullMatch := findProjectName(projectName, realProjectsToSearch)
	var projectSelected *reqApi42.API42Project
	if fullMatch {
		projectSelected = projectsFind[0]
	} else {
		indexChoose := getIndexNameChoice(projectsFindNames)
		if indexChoose == -1 {
			return
		}
		projectSelected = projectsFind[indexChoose]
	}
	fmt.Println(projectSelected)
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

	var indexAction int
	var err error
	menuActions := []string{
		cst.MenuActionFind,
		cst.MenuActionUpdateLocations,
		cst.MenuActionUpdateProjects,
		cst.MenuActionUpdateCursus,
		cst.MenuActionUpdateCampus,
		cst.MenuActionRefreshTokens,
		cst.MenuActionQuit,
	}
	for {
		prompt := promptui.Select{
			Label:	"Choose an action",
			Items:	menuActions,
			HideHelp: true,
		}

		indexAction, _, err = prompt.Run()

		if err != nil {
			log.Fatal().Err(err).Msg("Prompt: failed")
		}

		switch menuActions[indexAction] {
		case cst.MenuActionFind:
			if len(allProjects.directs) == 0 || len(allProjects.parents) == 0 {
				log.Error().Msg("Prompt: list of projects empty")
			} else {
				findExaminer(allProjects)
			}
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
		case cst.MenuActionRefreshTokens:
			api42.RefreshToken()
		case cst.MenuActionQuit:
			fmt.Println("Goodbye!")
			os.Exit(0)
		default:
			log.Fatal().Msg("Prompt: indexAction out of bound")
		}

	}
}


/*
--> Update locations
--> Update campus
--> Update cursus
--> Find project examiner
	---> Project parent ?
		---> Search your project
			---> Show n results or matching project
--> Refresh tokens
--> Exit
*/
