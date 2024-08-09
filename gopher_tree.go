package main

import (
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

type Config struct {
	dirsOnly       bool
	fullPath       bool
	url            bool
	html           bool
	gopher         bool
	realTime       bool
	maxDepth       int
	disableNotices bool
	typeFilter     []byte
	aliases        []string
}

// this should be rewritten so it can be called with a link to a file as a parameter, it would also eliminate some code in main.go
func gopherTree(cfg Config, rootInfo MenuItem, indentation string, filter *[]string, depth int) (string, error) {
	response, err := getGopher(rootInfo)
	if err != nil {
		return "", err
	}

	menu := responseToMenu(response)
	menu = cleanupMenu(menu, rootInfo, cfg)

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

		branch += indentation + pipe
		if cfg.html {
			branch += fmt.Sprintf("<a href=\"%s\">", item.URL())
		} else if cfg.url || (item.Host != rootInfo.Host) {
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

		if slices.Contains(*filter, item.Selector) {
			if !cfg.disableNotices {
				branch += " (already indexed)"
			}
			if cfg.html {
				branch += "<br />"
			}
			branch += "\n"

			continue
		}

		if !cfg.disableNotices && item.Host != rootInfo.Host {
			branch += " (foreign host)"
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
			if cfg.maxDepth > 0 && depth >= cfg.maxDepth {
				continue
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

func cleanupMenu(originalMenu []MenuItem, rootInfo MenuItem, cfg Config) (newMenu []MenuItem) {
	for _, item := range originalMenu {
		if item.Selector != "" && !slices.Contains(cfg.typeFilter, item.Type) && strings.HasSuffix(item.Host, rootInfo.Host) {
			newMenu = append(newMenu, item)
		}
	}

	return newMenu
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
