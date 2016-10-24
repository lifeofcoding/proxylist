package main

import (
	"regexp"
	"strings"
)

type Proxy struct {
	Unique     int    `json:"-"`
	Filter     int    `json:"-"`
	LastUpdate string `json:"lastupdate"`
	Address    string `json:"address"`
	Port       string `json:"port"`
	Country    string `json:"country"`
	Speed      string `json:"speed"`
	Connection string `json:"connection"`
	Protocol   string `json:"protocol"`
	Anonymity  string `json:"anonymity"`
}

func (p *Proxy) ParseAddress(line string) {
	var end int = strings.Index(line, "</style>")
	var content string = line[(end + 8):len(line)]

	// Remove unnecessary empty HTML tags.
	content = strings.Replace(content, "<span></span>", "", -1)

	// Remove explicit invisible HTML tags with inline CSS style.
	re := regexp.MustCompile(`<(span|div) style="display:none">[^<]+<\/(span|div)>`)
	content = re.ReplaceAllString(content, "")

	// Remove implicit invisible HTML tags with CSS classes.
	var invisible []string = p.InvisibleTags(line)
	for _, cssClass := range invisible {
		re = regexp.MustCompile(`<(span|div) class="` + cssClass + `">[^<]+<\/(span|div)>`)
		content = re.ReplaceAllString(content, "")
	}

	// Remove the display inline style from the remaining HTML code.
	re = regexp.MustCompile(`<(span|div) style="display: inline">([^<]+)<\/(span|div)>`)
	content = re.ReplaceAllString(content, "$2")

	// Remove the unnecessary HTML code from visible CSS classes.
	re = regexp.MustCompile(`<(span|div) class="[^"]+">([^<]+)<\/(span|div)>`)
	content = re.ReplaceAllString(content, "$2")

	// Clean final string.
	re = regexp.MustCompile(`([0-9\.]{7,15}).*`)
	content = re.ReplaceAllString(content, "$1")

	if content != "" {
		p.Address = content
	}
}

func (p *Proxy) InvisibleTags(line string) []string {
	var invisible []string
	var styleStart int = strings.Index(line, "<style>")
	var styleEnd int = strings.Index(line, "</style>")
	var section string = line[(styleStart + 7):styleEnd]

	re := regexp.MustCompile(`\.([0-9a-zA-Z_\-]{4})\{display:(inline|none)\}`)
	parts := re.FindAllStringSubmatch(section, -1)

	for _, data := range parts {
		if data[2] == "none" {
			invisible = append(invisible, data[1])
		}
	}

	return invisible
}

func (p *Proxy) ParseLastUpdate(line string) {
	re := regexp.MustCompile(`<span class="updatets[^>]+>([^<]+)<\/span>`)
	parts := re.FindStringSubmatch(line)

	if len(parts) == 2 {
		p.LastUpdate = strings.TrimSpace(parts[1])
	}
}

func (p *Proxy) ParsePort(line string) {
	re := regexp.MustCompile(`<td>(\S+)\s+<\/td>`)
	parts := re.FindStringSubmatch(line)

	if len(parts) == 2 {
		p.Port = parts[1]
	}
}

func (p *Proxy) ParseCountry(line string) {
	re := regexp.MustCompile(`\/>([^<]+)<\/span><\/td>`)
	parts := re.FindStringSubmatch(line)

	if len(parts) == 2 {
		p.Country = strings.TrimSpace(parts[1])
	}
}

func (p *Proxy) ParseSpeed(line string) {
	re := regexp.MustCompile(`style="width: ([0-9]+%);`)
	parts := re.FindStringSubmatch(line)

	if len(parts) == 2 {
		p.Speed = parts[1]
	}
}

func (p *Proxy) ParseConnection(line string) {
	re := regexp.MustCompile(`style="width: ([0-9]+%);`)
	parts := re.FindStringSubmatch(line)

	if len(parts) == 2 {
		p.Connection = parts[1]
	}
}

func (p *Proxy) ParseProtocol(line string) {
	re := regexp.MustCompile(`<td>(\S+)\s+<\/td>`)
	parts := re.FindStringSubmatch(line)

	if len(parts) == 2 {
		p.Protocol = parts[1]
	}
}

func (p *Proxy) ParseAnonymity(line string) {
	re := regexp.MustCompile(`<td nowrap>([^>]+)<\/td>`)
	parts := re.FindStringSubmatch(line)

	if len(parts) == 2 {
		p.Anonymity = strings.TrimSpace(parts[1])
	}
}
