package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	var cfg Config
	cfg.typeFilter = []byte{'h', 'i', '3'} // default type filter

	var ignoredTypesString, aliasesString string

	flag.BoolVar(&cfg.dirsOnly, "d", false, "List directories only")
	flag.BoolVar(&cfg.fullPath, "f", false, "Prints the full selector for each item")
	flag.BoolVar(&cfg.printType, "t", false, "Prints item type for each item")
	flag.BoolVar(&cfg.url, "u", false, "Prints a Gopher URL for each item")
	flag.BoolVar(&cfg.html, "h", false, "Outputs the tree as an HTML page with links to items")
	flag.BoolVar(&cfg.gopher, "g", false, "Outputs the tree as a Gopher menu with links to items")
	flag.BoolVar(&cfg.realTime, "r", false, "Prints individual lines of the tree as they are generated in real time")
	flag.IntVar(&cfg.maxDepth, "L", 0, "Max depth level of the tree")
	flag.BoolVar(&cfg.disableNotices, "N", false, "Hide notices for already indexed and foreign items")
	flag.StringVar(&ignoredTypesString, "T", "", "Comma-separated list of item types you're willing to ignore")
	flag.StringVar(&aliasesString, "a", "", "Comma-separated list of target server's alias hostnames")
	flag.BoolVar(&cfg.foreign, "F", false, "Show links to files and directories on foreign hosts")

	flag.Parse()

	// disgusting code
	ignoredTypeList := strings.Split(ignoredTypesString, ",")
	for _, ignType := range ignoredTypeList {
		if len(ignType) > 1 {
			fmt.Fprintf(os.Stderr, "%s: Warning: item types are represented by a single character. \"%s\" contains %d characters and will be ignored.\n",
				os.Args[0], ignType, len(ignType))
		} else if len(ignType) == 1 {
			cfg.typeFilter = append(cfg.typeFilter, ignType[0])
		}
	}

	cfg.aliases = strings.Split(aliasesString, ",")

	args := flag.Args()
	if len(args) == 0 {
		usage(os.Args[0])
		os.Exit(1)
	}

	argument := args[0]

	rootInfo, err := URLToMenuItem(argument)
	cfg.aliases = append(cfg.aliases, rootInfo.Host)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	title := rootInfo.URL()
	if cfg.html {
		title = rootInfo.HTML() + "<br />"
	} else if cfg.gopher {
		rootInfo.DisplayName = title
		title = rootInfo.String()
	} else {
		title += "\n"
	}
	fmt.Print(title)

	tree, err := gopherTree(cfg, rootInfo, "", &[]string{""}, 1)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	} else {
		fmt.Print(tree)
	}

}

func usage(progname string) {
	fmt.Fprintf(os.Stderr, "Usage: %s [ option ... ] URL\n", progname)
}
