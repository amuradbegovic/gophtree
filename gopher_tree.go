package main

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

type Config struct {
	dirsOnly bool
	fullPath bool
	url      bool
	html     bool
	gopher   bool
	realTime bool
	maxDepth int
}

func gopherTree(cfg Config, rootInfo MenuItem, indentation string, filter *[]string, depth int) (string, error) {
	response, err := getGopher(rootInfo)
	if err != nil {
		return "", err
	}

	menu := responseToMenu(response)
	menu = cleanupMenu(menu, rootInfo)

	*filter = append(*filter, rootInfo.Selector)

	tree := ""

	for i, item := range menu {
		if cfg.dirsOnly && item.Type != '1' {
			continue
		}

		branch := ""

		pipe := "├── "
		if i == len(menu)-1 {
			pipe = "└── "
		}

		//tree += indentation + pipe + path.Base(item.Selector)

		branch += indentation + pipe
		if cfg.html {
			branch += fmt.Sprintf("<a href=\"%s\">", item.URL())
		}
		if cfg.url {
			branch += item.URL()
		} else {
			if cfg.fullPath {
				branch += item.Selector
			} else {
				branch += MakePath(rootInfo.Selector, item.Selector)
			}
		}
		if cfg.html {
			branch += "</a>"
		}

		if contains(*filter, item.Selector) {
			branch += " (already indexed)"
			if cfg.html {
				branch += "<br />"
			}
			branch += "\n"

			continue
		}

		if cfg.html {
			branch += "<br />"
		}

		if cfg.gopher {
			item.DisplayName = branch
			branch = item.String()
		} else {
			branch += "\n"
		}

		if cfg.realTime {
			fmt.Print(branch)
		}

		tree += branch

		*filter = append(*filter, item.Selector)

		if item.Type == '1' {
			if cfg.maxDepth != 0 {
				if depth >= cfg.maxDepth {
					continue
				}
			}
			addedInd := "│   "
			if i == len(menu)-1 {
				addedInd = "    "
			}
			if cfg.html {
				addedInd = strings.Replace(addedInd, " ", "&nbsp;&nbsp;", -1)
			}
			toAdd, err := gopherTree(cfg, item, indentation+addedInd, filter, depth+1)
			if err != nil {
				return "", err
			}
			tree += toAdd
		}
	}

	return tree, nil
}

/*func getDepth(indentation string) int { // ugly hack tbh
	return len(indentation) / 4
}*/

func contains(slice []string, element string) bool {
	for _, el := range slice {
		if el == element {
			return true
		}
	}
	return false
}

func cleanupMenu(original_menu []MenuItem, rootInfo MenuItem) (new_menu []MenuItem) {
	for _, item := range original_menu {
		if item.Selector != "" && item.Selector != "Err" && item.Type != 'h' && strings.HasSuffix(item.Host, rootInfo.Host) {
			//fmt.Println(item)
			new_menu = append(new_menu, item)
		}
	}
	return new_menu
}

func MakePath(menuPath string, linkPath string) string {
	result := linkPath
	if strings.HasSuffix(menuPath, ".gph") || strings.HasSuffix(menuPath, ".gophermap") {
		menuPath = path.Dir(menuPath)
	}

	if strings.HasPrefix(linkPath, path.Dir(menuPath)) {
		result, _ = filepath.Rel(menuPath, linkPath)
	}
	return result
}
