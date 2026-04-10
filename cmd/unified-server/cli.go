package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type usageSection struct {
	Title string
	Flags []string
}

var cliUsageSections = []usageSection{
	{Title: "General", Flags: []string{"version", "domain", "port", "support-email"}},
	{Title: "Database", Flags: []string{"db-type", "db-path", "db-conn"}},
	{Title: "Map defaults", Flags: []string{"default-lat", "default-lon", "default-zoom", "default-layer", "auto-locate-default"}},
	{Title: "Realtime & archives", Flags: []string{"safecast-realtime", "json-archive-path", "json-archive-frequency", "import-tgz-url", "import-tgz-file"}},
	{Title: "Safecast API fetcher", Flags: []string{"safecast-fetcher", "safecast-fetcher-interval", "safecast-fetcher-batch-size", "safecast-fetcher-start-date"}},
	{Title: "Authentication", Flags: []string{"require-auth", "allow-registration", "session-secret", "session-cookie-name", "session-duration", "base-url"}},
	{Title: "Email (SMTP)", Flags: []string{"smtp-host", "smtp-port", "smtp-username", "smtp-password", "smtp-from", "smtp-from-name"}},
	{Title: "Self-upgrade", Flags: []string{"selfupgrade", "selfupgrade-url"}},
}

type cliColorTheme struct {
	Enabled bool
	Section string
	Flag    string
	Usage   string
	Default string
	Reset   string
}

func resolveCLIColorTheme(out io.Writer) cliColorTheme {
	theme := cliColorTheme{}
	file, ok := out.(*os.File)
	if !ok {
		return theme
	}
	if os.Getenv("NO_COLOR") != "" {
		return theme
	}
	info, err := file.Stat()
	if err != nil {
		return theme
	}
	if (info.Mode() & os.ModeCharDevice) == 0 {
		return theme
	}

	theme.Enabled = true
	theme.Section = "\033[38;5;25m"
	theme.Flag = "\033[38;5;208m"
	theme.Usage = "\033[38;5;240m"
	theme.Default = "\033[38;5;34m"
	theme.Reset = "\033[0m"
	return theme
}

func configureCLIUsage() {
	flag.CommandLine.SetOutput(os.Stdout)
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		theme := resolveCLIColorTheme(out)

		fmt.Fprintf(out, "Usage: %s [flags]\n\n", os.Args[0])
		if theme.Enabled {
			fmt.Fprintf(out, "%sFlags:%s\n", theme.Section, theme.Reset)
		} else {
			fmt.Fprintln(out, "Flags:")
		}

		printed := map[string]bool{}
		for _, section := range cliUsageSections {
			var sectionFlags []*flag.Flag
			for _, name := range section.Flags {
				if f := flag.Lookup(name); f != nil {
					sectionFlags = append(sectionFlags, f)
					printed[f.Name] = true
				}
			}
			if len(sectionFlags) == 0 {
				continue
			}

			if theme.Enabled {
				fmt.Fprintf(out, "%s%s:%s\n", theme.Section, section.Title, theme.Reset)
			} else {
				fmt.Fprintf(out, "%s:\n", section.Title)
			}
			for _, f := range sectionFlags {
				writeFlagUsage(out, f, theme)
			}
			fmt.Fprintln(out)
		}

		var leftovers []string
		flag.VisitAll(func(f *flag.Flag) {
			if !printed[f.Name] {
				leftovers = append(leftovers, f.Name)
			}
		})
		if len(leftovers) > 0 {
			sort.Strings(leftovers)
			if theme.Enabled {
				fmt.Fprintf(out, "%sAdditional flags:%s\n", theme.Section, theme.Reset)
			} else {
				fmt.Fprintln(out, "Additional flags:")
			}
			for _, name := range leftovers {
				if f := flag.Lookup(name); f != nil {
					writeFlagUsage(out, f, theme)
				}
			}
		}

		printCLILicenseNote(out, theme)
	}
}

func writeFlagUsage(out io.Writer, f *flag.Flag, theme cliColorTheme) {
	if f == nil {
		return
	}
	name, usage := flag.UnquoteUsage(f)
	if theme.Enabled {
		fmt.Fprintf(out, "  %s-%s%s", theme.Flag, f.Name, theme.Reset)
	} else {
		fmt.Fprintf(out, "  -%s", f.Name)
	}
	if name != "" {
		fmt.Fprintf(out, " %s", name)
	}
	fmt.Fprintln(out)

	if usage != "" {
		for _, part := range strings.Split(usage, "\n") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if theme.Enabled {
				fmt.Fprintf(out, "      %s%s%s\n", theme.Usage, part, theme.Reset)
			} else {
				fmt.Fprintf(out, "      %s\n", part)
			}
		}
	}

	if def := strings.TrimSpace(f.DefValue); def != "" {
		if theme.Enabled {
			fmt.Fprintf(out, "      %sDefault:%s %s%s%s\n", theme.Flag, theme.Reset, theme.Default, def, theme.Reset)
		} else {
			fmt.Fprintf(out, "      Default: %s\n", def)
		}
	}
}

func printCLILicenseNote(out io.Writer, theme cliColorTheme) {
	if out == nil {
		return
	}

	fmt.Fprintln(out)
	if theme.Enabled {
		fmt.Fprintf(out, "%sLicense & community:%s\n", theme.Section, theme.Reset)
	} else {
		fmt.Fprintln(out, "License & community:")
	}

	lines := []string{
		"Code: MIT License.",
		"Research datasets: CC0 1.0 Universal (Public Domain).",
		"Thank you for using this program and sharing your tracks. This work is fragile — care for it, and it will grow.",
		"Support the sources, share honest knowledge, and run your own nodes so the maps stay free and safe.",
	}

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if theme.Enabled {
			fmt.Fprintf(out, "  %s%s%s\n", theme.Usage, line, theme.Reset)
		} else {
			fmt.Fprintf(out, "  %s\n", line)
		}
	}
}
