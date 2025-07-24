package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"

	"github.com/ahmedYasserM/fo/internal/colors"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [URL]",
	Short: "Fetch sample test cases from a Codeforces problem URL",
	Long: `Fetch downloads sample input and output from a given Codeforces problem URL.
The samples are saved to 'testcases.txt'.

Example:
  fo fetch https://codeforces.com/contest/1234/problem/A`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fetchSamples(args[0])
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}

// extractPreText extracts multiline text inside <pre>, handling:
// - multiple <div> lines (e.g. class="test-example-line")
// - or <br> tags as line separators
// It preserves the exact formatting of sample inputs/outputs.
func extractPreText(s *goquery.Selection) string {
	// Check if <pre> contains <div> children; if yes, extract each div as line
	divs := s.ChildrenFiltered("div")
	if divs.Length() > 0 {
		var lines []string
		divs.Each(func(i int, sel *goquery.Selection) {
			lines = append(lines, sel.Text())
		})
		return strings.Join(lines, "\n")
	}

	// Otherwise, extract text handling <br> tags as newlines
	var builder strings.Builder
	for _, node := range s.Nodes {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			extractNodeTextWithBr(c, &builder)
		}
	}

	return builder.String()
}

// Helper: recursively extract text traversing HTML nodes
// Converts <br> tags to newline characters
func extractNodeTextWithBr(n *html.Node, builder *strings.Builder) {
	switch n.Type {
	case html.TextNode:
		builder.WriteString(n.Data)
	case html.ElementNode:
		if n.Data == "br" {
			builder.WriteByte('\n')
		} else {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				extractNodeTextWithBr(c, builder)
			}
		}
	}
}

func fetchSamples(rawurl string) error {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
			"(KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"),
	)

	var inputs []string
	var outputs []string

	c.OnHTML("div.sample-test", func(e *colly.HTMLElement) {
		fmt.Println("Found sample-test block on the page") // Debug print

		e.ForEach("div.input", func(idx int, el *colly.HTMLElement) {
			el.DOM.Find("pre").Each(func(_ int, pre *goquery.Selection) {
				text := strings.TrimSpace(extractPreText(pre))
				fmt.Printf("Sample input #%d:\n%s\n---\n", idx+1, text) // Debug print
				inputs = append(inputs, text)
			})
		})

		e.ForEach("div.output", func(idx int, el *colly.HTMLElement) {
			el.DOM.Find("pre").Each(func(_ int, pre *goquery.Selection) {
				text := strings.TrimSpace(extractPreText(pre))
				fmt.Printf("Sample output #%d:\n%s\n---\n", idx+1, text) // Debug print
				outputs = append(outputs, text)
			})
		})
	})

	err := c.Visit(rawurl)
	if err != nil {
		return fmt.Errorf("failed to visit URL %w", err)
	}

	if len(inputs) == 0 || len(inputs) != len(outputs) {
		return fmt.Errorf("could not find matching sample inputs and outputs")
	}

	f, err := os.Create("testcases.txt")
	if err != nil {
		return fmt.Errorf("failed to create testcases.txt: %w", err)
	}
	defer f.Close()

	for i := range inputs {
		fmt.Fprintf(f, "--- Sample #%d Input ---\n%s\n\n", i+1, inputs[i])
		fmt.Fprintf(f, "--- Sample #%d Output ---\n%s\n\n", i+1, outputs[i])
	}

	fmt.Printf("%s✅ Saved %d sample(s) to testcases.txt%s\n", colors.GREEN, len(inputs), colors.RESET)
	return nil
}
