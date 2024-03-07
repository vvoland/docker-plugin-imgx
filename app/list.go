package app

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/containerd/platforms"
	"github.com/docker/cli/cli/command"
	imagetypes "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/go-units"
	"github.com/fatih/color"
)

func List(ctx context.Context, dockerCLI command.Cli, args []string) error {
	images, err := dockerCLI.Client().ImageList(ctx, imagetypes.ListOptions{
		All:            true,
		ContainerCount: true,
	})
	if err != nil {
		return err
	}

	var view []topImage
	for _, img := range images {
		details := imageDetails{
			ID:   img.ID,
			Size: units.HumanSizeWithPrecision(float64(img.Size), 3),
			Used: img.Containers > 0,
		}

		var children []subImage
		for _, platform := range img.PlatformImages {
			sub := subImage{
				Platform:  platforms.Format(platform.Platform),
				Available: platform.Available,
				Details: imageDetails{
					ID:   platform.ID,
					Size: units.HumanSizeWithPrecision(float64(platform.ContentSize+platform.UnpackedSize), 3),
					Used: platform.Containers > 0,
				},
			}

			children = append(children, sub)
		}

		for _, tag := range img.RepoTags {
			view = append(view, topImage{
				Name:     tag,
				Details:  details,
				Children: children,
			})
		}
	}

	return printImageTree(dockerCLI, view)
}

type imageDetails struct {
	ID   string
	Size string
	Used bool
}

type topImage struct {
	Name     string
	Details  imageDetails
	Children []subImage
}

type subImage struct {
	Platform  string
	Available bool
	Details   imageDetails
}

func printImageTree(dockerCLI command.Cli, images []topImage) error {
	out := dockerCLI.Out()
	_, width := out.GetTtySize()

	headers := []header{
		{Title: "Image", Width: 0, Left: true},
		{Title: "ID", Width: 12},
		{Title: "Size", Width: 8},
		{Title: "Used", Width: 4},
	}

	const spacing = 3
	nameWidth := int(width)
	for _, h := range headers {
		if h.Width == 0 {
			continue
		}
		nameWidth -= h.Width
		nameWidth -= spacing
	}

	maxImageName := len(headers[0].Title)
	for _, img := range images {
		if len(img.Name) > maxImageName {
			maxImageName = len(img.Name)
		}
		for _, sub := range img.Children {
			if len(sub.Platform) > maxImageName {
				maxImageName = len(sub.Platform)
			}
		}
	}

	if nameWidth > maxImageName+spacing {
		nameWidth = maxImageName + spacing
	}

	if nameWidth < 0 {
		headers = headers[:1]
		nameWidth = int(width)
	}
	headers[0].Width = nameWidth

	headerColor := color.New(color.FgHiWhite).Add(color.Bold)

	// Print headers
	for i, h := range headers {
		if i > 0 {
			_, _ = fmt.Fprint(out, strings.Repeat(" ", spacing))
		}

		headerColor.Fprint(out, h.PrintC(headerColor, h.Title))
	}

	_, _ = fmt.Fprintln(out)

	topNameColor := color.New(color.FgBlue).Add(color.Underline).Add(color.Bold)
	normalColor := color.New(color.FgWhite)
	normalFaintedColor := color.New(color.FgWhite).Add(color.Faint)
	greenColor := color.New(color.FgGreen)

	printDetails := func(clr *color.Color, details imageDetails) {
		truncID := stringid.TruncateID(details.ID)
		fmt.Fprint(out, headers[1].Print(clr, truncID))
		fmt.Fprint(out, strings.Repeat(" ", spacing))

		fmt.Fprint(out, headers[2].Print(clr, details.Size))
		fmt.Fprint(out, strings.Repeat(" ", spacing))

		if details.Used {
			fmt.Fprint(out, headers[3].Print(greenColor, " ✔ ️"))
		} else {
			fmt.Fprint(out, headers[3].Print(clr, " "))
		}
	}

	// Print images
	for _, img := range images {
		fmt.Fprint(out, headers[0].Print(topNameColor, img.Name))
		fmt.Fprint(out, strings.Repeat(" ", spacing))

		printDetails(normalColor, img.Details)

		_, _ = fmt.Fprintln(out, "")
		for idx, sub := range img.Children {
			clr := normalColor
			if !sub.Available {
				clr = normalFaintedColor
			}

			if idx != len(img.Children)-1 {
				fmt.Fprint(out, headers[0].Print(clr, "├─ "+sub.Platform))
			} else {
				fmt.Fprint(out, headers[0].Print(clr, "└─ "+sub.Platform))
			}

			fmt.Fprint(out, strings.Repeat(" ", spacing))
			printDetails(clr, sub.Details)

			fmt.Fprintln(out, "")
		}
	}

	return nil
}

func maybeUint(v int64) *uint {
	u := uint(v)
	return &u
}

type header struct {
	Title string
	Width int
	Left  bool
}

func truncateRunes(s string, length int) string {
	runes := []rune(s)
	if len(runes) > length {
		return string(runes[:length])
	}
	return s
}

func (h header) Print(color *color.Color, s string) (out string) {
	if h.Left {
		return h.PrintL(color, s)
	}
	return h.PrintC(color, s)
}

func (h header) PrintC(color *color.Color, s string) (out string) {
	ln := utf8.RuneCountInString(s)
	if h.Left {
		return h.PrintL(color, s)
	}

	if ln > int(h.Width) {
		return color.Sprint(truncateRunes(s, h.Width))
	}

	fill := int(h.Width) - ln

	l := fill / 2
	r := fill - l

	return strings.Repeat(" ", l) + color.Sprint(s) + strings.Repeat(" ", r)
}

func (h header) PrintL(color *color.Color, s string) string {
	ln := utf8.RuneCountInString(s)
	if ln > int(h.Width) {
		return color.Sprint(truncateRunes(s, h.Width))
	}

	return color.Sprint(s) + strings.Repeat(" ", int(h.Width)-ln)
}
