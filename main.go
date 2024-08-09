/*
Options I should give to the user:
- delay a certain amount of ms between requests as to not overflow the target server
- show/hide warnings (already indexed, etc...)
- index only target server, not linked servers - give an option to mark and/or index "foreign" hosts
- aliases for target server
- exclude menu item types
- show size

write a man page
add license
git
*/

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	var cfg Config

	flag.BoolVar(&cfg.dirsOnly, "d", false, "List directories only")
	flag.BoolVar(&cfg.fullPath, "f", false, "Prints the full selector for each item")
	flag.BoolVar(&cfg.url, "u", false, "Prints a Gopher URL for each item")
	flag.BoolVar(&cfg.html, "h", false, "Outputs the tree as an HTML page with links to items")
	flag.BoolVar(&cfg.gopher, "g", false, "Outputs the tree as a Gopher menu with links to items")
	flag.BoolVar(&cfg.realTime, "r", false, "Prints individual lines of the tree as they are generated in real time")
	flag.IntVar(&cfg.maxDepth, "L", 0, "level\n\tMax display depth of the directory tree")

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		usage(os.Args[0])
		os.Exit(1)
	}

	for _, argument := range args {
		rootInfo, err := URLToMenuItem(argument)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		tree, err := gopherTree(cfg, rootInfo, "", &[]string{""}, 1)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		} else {
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
			fmt.Print(tree)
		}
	}
}

func usage(progname string) {
	fmt.Fprintf(os.Stderr, "Usage: %s [ option ... ] url1 [ url2 ... ]\n", progname)
}
