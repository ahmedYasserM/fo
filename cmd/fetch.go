package cmd

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ahmedYasserM/fo/internal/colors"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch [URL]",
	Short: "Fetches sample test cases from a Codeforces problem URL",
	Long: `Fetch downloads sample input and output from a given Codeforces problem URL.
The samples are saved to 'testcases.txt'.

Example:
  fo fetch https://codeforces.com/contest/1234/problem/A
  fo fetch https://codeforces.com/group/abcdef/contest/12345/problem/B`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rawurl := args[0]
		return fetchSamples(rawurl)
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}

// parseURL parses Codeforces URLs including group contests.
// Returns urlType, groupID (empty if none), contestID, problemIndex.
func parseURL(rawurl string) (string, string, string, string, error) {
	parsed, err := url.Parse(rawurl)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to parse URL: %w", err)
	}

	pathParts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(pathParts) < 3 {
		return "", "", "", "", fmt.Errorf("invalid URL path structure")
	}

	switch pathParts[0] {
	case "group":
		// /group/<groupId>/contest/<contestId>/problem/<problemIndex>
		if len(pathParts) >= 6 && pathParts[2] == "contest" && pathParts[4] == "problem" {
			groupId := pathParts[1]
			contestId := pathParts[3]
			problemIndex := strings.ToUpper(pathParts[5])
			return "group_contest", groupId, contestId, problemIndex, nil
		}
		return "", "", "", "", fmt.Errorf("unsupported group URL structure: %q", parsed.Path)

	case "contest":
		if len(pathParts) < 4 || pathParts[2] != "problem" {
			return "", "", "", "", fmt.Errorf("expected /contest/{id}/problem/{index}")
		}
		return "contest", "", pathParts[1], strings.ToUpper(pathParts[3]), nil

	case "problemset":
		if len(pathParts) < 5 || pathParts[1] != "problem" {
			return "", "", "", "", fmt.Errorf("expected /problemset/problem/{id}/{index}")
		}
		return "problemset", "", pathParts[2], strings.ToUpper(pathParts[3]), nil

	case "gym":
		if len(pathParts) < 4 || pathParts[2] != "problem" {
			return "", "", "", "", fmt.Errorf("expected /gym/{id}/problem/{index}")
		}
		return "gym", "", pathParts[1], strings.ToUpper(pathParts[3]), nil

	default:
		return "", "", "", "", fmt.Errorf("unsupported URL type: %q", pathParts[0])
	}
}

func fetchSamples(rawurl string) error {
	urlType, groupID, contestID, problemIndex, err := parseURL(rawurl)
	if err != nil {
		return fmt.Errorf("error parsing URL: %w", err)
	}

	var fetchURL string
	switch urlType {
	case "group_contest":
		fetchURL = fmt.Sprintf("https://codeforces.com/group/%s/contest/%s/problem/%s", groupID, contestID, problemIndex)
	case "contest":
		fetchURL = fmt.Sprintf("https://codeforces.com/contest/%s/problem/%s", contestID, problemIndex)
	case "problemset":
		fetchURL = fmt.Sprintf("https://codeforces.com/problemset/problem/%s/%s", contestID, problemIndex)
	case "gym":
		fetchURL = fmt.Sprintf("https://codeforces.com/gym/%s/problem/%s", contestID, problemIndex)
	default:
		return fmt.Errorf("unsupported URL type: %s", urlType)
	}

	fmt.Printf("Fetching samples from: %s%s%s\n", colors.CYAN, fetchURL, colors.RESET)

	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("failed to create cookie jar: %w", err)
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: 20 * time.Second,
	}

	homeReq, err := http.NewRequest("GET", "https://codeforces.com/", nil)
	if err != nil {
		return fmt.Errorf("failed to create homepage request: %w", err)
	}
	setCommonHeaders(homeReq)
	_, err = client.Do(homeReq)
	if err != nil {
		return fmt.Errorf("failed to fetch homepage: %w", err)
	}

	req, err := http.NewRequest("GET", fetchURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	setCommonHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return fmt.Errorf("received HTTP 403 Forbidden. Server may block automated access. Try using VPN or increasing delays.")
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to fetch page, status code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse HTML: %w", err)
	}

	sampleDiv := findDivByClass(doc, "sample-test")
	if sampleDiv == nil {
		return fmt.Errorf("could not find sample-test block on the page")
	}

	inputs := findAllDivsByClass(sampleDiv, "input")
	outputs := findAllDivsByClass(sampleDiv, "output")

	if len(inputs) == 0 || len(outputs) == 0 || len(inputs) != len(outputs) {
		return fmt.Errorf("sample input/output counts do not match or missing")
	}

	f, err := os.Create("testcases.txt")
	if err != nil {
		return fmt.Errorf("could not create testcases.txt: %w", err)
	}
	defer f.Close()

	for i := range inputs {
		inputText := extractPreText(getPreChild(inputs[i]))
		outputText := extractPreText(getPreChild(outputs[i]))
		_, err := fmt.Fprintf(f, "--- Sample #%d Input ---\n%s\n\n", i+1, inputText)
		if err != nil {
			return fmt.Errorf("failed to write input sample: %w", err)
		}
		_, err = fmt.Fprintf(f, "--- Sample #%d Output ---\n%s\n\n", i+1, outputText)
		if err != nil {
			return fmt.Errorf("failed to write output sample: %w", err)
		}
	}

	fmt.Printf("%s✅ Saved %d sample(s) to testcases.txt%s\n", colors.GREEN, len(inputs), colors.RESET)
	return nil
}

func setCommonHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", "https://codeforces.com/")
	req.Header.Set("Connection", "keep-alive")
}

func findDivByClass(n *html.Node, className string) *html.Node {
	var found *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if found != nil {
			return
		}
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == className {
					found = n
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return found
}

func findAllDivsByClass(n *html.Node, className string) []*html.Node {
	var result []*html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == className {
					result = append(result, n)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return result
}

func extractPreText(n *html.Node) string {
	if n == nil {
		return ""
	}
	var b strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
			b.WriteString(n.Data)
		case html.ElementNode:
			if n.Data == "br" {
				b.WriteByte('\n')
			} else {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c)
				}
			}
		}
	}
	f(n)
	return b.String()
}

func getPreChild(n *html.Node) *html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "pre" {
			return c
		}
	}
	return nil
}
