package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type DorkGenerator struct {
	Templates      []string
	Dictionaries   map[string][]string
	SpecialTargets map[string]map[string][]string
	Countries      map[string][]string
	Domains        []string
	httpClient     *http.Client
}

var AutoWordlistURLs = map[string]string{
	"english_auto": "https://raw.githubusercontent.com/dwyl/english-words/master/words.txt",
	"spanish_auto": "https://raw.githubusercontent.com/javierarce/palabras/master/listado-general.txt",
	"german_auto":  "https://gist.githubusercontent.com/MarvinJWendt/2f4f4154b8ae218600eb091a5706b5f4/raw/36b70dd6be330aa61cd4d4cdfda6234dcb0b8784/wordlist-german.txt",
	"russian_auto": "https://raw.githubusercontent.com/hingston/russian/refs/heads/master/100000-russian-words.txt",
}

func GenerateVariations(baseWords []string, count int) []string {
	if len(baseWords) == 0 {
		return []string{}
	}

	var variations []string
	if len(baseWords) < count {
		variations = append(variations, baseWords...)
	} else {
		tempWords := make([]string, len(baseWords))
		copy(tempWords, baseWords)
		rand.Shuffle(len(tempWords), func(i, j int) {
			tempWords[i], tempWords[j] = tempWords[j], tempWords[i]
		})
		variations = append(variations, tempWords[:count/2]...)
	}

	suffixes := []string{"_backup", "_conf", "_log", "_old", "_dev", "_test", "-db", "s", "es", "data"}
	prefixes := []string{"dev_", "test_", "old_", "new_"}
	numbers := []string{"1", "123", "2023", "2024", "2025", "00", "01"}

	for len(variations) < count {
		if len(baseWords) == 0 {
			break
		}
		word := baseWords[rand.Intn(len(baseWords))]

		switch rand.Intn(4) {
		case 0:
			if len(suffixes) > 0 {
				variations = append(variations, word+suffixes[rand.Intn(len(suffixes))])
			}
		case 1:
			if len(prefixes) > 0 {
				variations = append(variations, prefixes[rand.Intn(len(prefixes))]+word)
			}
		case 2:
			if len(numbers) > 0 {
				variations = append(variations, word+numbers[rand.Intn(len(numbers))])
			}
		case 3:
			if len(baseWords) > 1 {
				word2 := baseWords[rand.Intn(len(baseWords))]
				if word != word2 && len(word)+len(word2) > 5 {
					variations = append(variations, word+word2)
					variations = append(variations, word+"-"+word2)
					variations = append(variations, word+"_"+word2)
				}
			}
		}
	}
	return variations
}

func LoadWordsFromFileOrURL(filename string, githubURL string, httpClient *http.Client) []string {
	f, err := os.Open(filename)
	if err == nil {
		defer f.Close()
		var words []string
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				words = append(words, strings.ToLower(line))
			}
		}
		fmt.Printf("✅ Successfully loaded words from local file: %s\n", filename)
		return words
	}

	if githubURL != "" {
		fmt.Printf("⚠️ Warning: Local file %s not found. Attempting to download from GitHub: %s\n", filename, githubURL)
		req, err := http.NewRequest("GET", githubURL, nil)
		if err != nil {
			fmt.Printf("❌ Error creating request for %s: %v. Skipping this dictionary.\n", githubURL, err)
			return []string{}
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; GoogleDorkGenerator/1.0; +http://example.com/bot)")

		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Printf("❌ Error downloading words from %s: %v. Skipping this dictionary.\n", githubURL, err)
			return []string{}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("❌ Error downloading words from %s: Received status code %d. Skipping this dictionary.\n", githubURL, resp.StatusCode)
			return []string{}
		}

		localFile, err := os.Create(filename)
		if err != nil {
			fmt.Printf("❌ Error creating local file %s to save downloaded words: %v. Continuing without saving locally.\n", filename, err)
			var words []string
			scanner := bufio.NewScanner(resp.Body)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" && !strings.HasPrefix(line, "#") {
					words = append(words, strings.ToLower(line))
				}
			}
			return words
		} else {
			defer localFile.Close()
			teeReader := io.TeeReader(resp.Body, localFile)
			var words []string
			scanner := bufio.NewScanner(teeReader)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" && !strings.HasPrefix(line, "#") {
					words = append(words, strings.ToLower(line))
				}
			}
			fmt.Printf("✅ Successfully downloaded and loaded words from GitHub, saved to: %s\n", filename)
			return words
		}
	}

	fmt.Printf("⚠️ Warning: Could not load words for %s. Neither local file nor GitHub URL provided or accessible. Skipping this dictionary.\n", filename)
	return []string{}
}

func (g *DorkGenerator) loadAutoDictionaries() {
	fmt.Println("⚙️ AUTO mode enabled: Loading multiple language dictionaries...")
	for lang, url := range AutoWordlistURLs {
		filename := fmt.Sprintf("common_%s_words.txt", strings.Split(lang, "_")[0])
		words := LoadWordsFromFileOrURL(filename, url, g.httpClient)
		if len(words) > 0 {
			g.Dictionaries[lang] = words
			fmt.Printf("Loaded %d words for '%s'.\n", len(words), lang)
		}
	}
}

func NewDorkGenerator(customWordsSource string) *DorkGenerator {
	g := &DorkGenerator{
		Dictionaries: make(map[string][]string),
		SpecialTargets: map[string]map[string][]string{
			"cms": {
				"wordpress": {
					"wp-admin", "wp-content", "wp-includes", "wp-login", "wp-config",
					"wp-json", "xmlrpc.php", "wp-cron.php", "wp-signup.php",
					"wp-activate.php", "wp-links-opml.php",
				},
				"joomla": {
					"administrator", "joomla", "index.php?option=com",
					"index.php?option=com_users", "index.php?option=com_content",
					"index.php?option=com_contact", "index.php?option=com_weblinks",
				},
				"drupal": {
					"user/login", "drupal", "?q=user/login", "?q=node", "?q=admin",
					"?q=filter/tips", "user/register", "user/password", "?q=search",
				},
				"magento": {
					"adminhtml", "magento", "/admin/dashboard", "downloader",
					"rss/catalog", "rss/order", "customer/account/login", "checkout/cart",
					"catalogsearch/result",
				},
				"prestashop": {
					"authentication", "prestashop", "admin123",
					"order-history", "my-account", "addresses", "identity",
					"guest-tracking", "order-follow",
				},
			},
			"frameworks": {
				"laravel": {"laravel", "storage/logs", "/login", "/admin", ".env", "artisan"},
				"django":  {"admin/login", "django", "/accounts/login", "settings.py", "manage.py"},
				"rails":   {"rails", "admin", "/users/sign_in", "routes.rb", "secrets.yml"},
			},
			"servers": {
				"apache": {"server-status", "apache", "htaccess", "htpasswd", "cgi-bin", "error.log"},
				"nginx":  {"nginx-status", "nginx", "nginx.conf", "default.conf", "error.log", "access.log"},
				"iis":    {"iisstart.htm", "microsoft-iis", "web.config", "global.asax", "bin/", "App_Code/"},
			},
		},
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
	g.Templates = g.LoadTemplates("base_templates.txt")

	g.Dictionaries["admin"] = []string{
		"admin", "login", "panel", "control", "dashboard", "manager", "backend",
		"console", "secure", "signin", "cpanel", "administrator",
	}
	g.Dictionaries["types"] = []string{
		"php", "asp", "aspx", "jsp", "cfm", "html", "htm", "xml", "json", "sql",
		"txt", "doc", "pdf", "xls", "csv", "bak", "zip",
	}
	g.Dictionaries["vulnerabilities"] = []string{
		"vulnerable", "exploit", "inurl", "intitle", "intext", "index of",
		"password", "config", "backup", "database", "leak", "exposed",
		"credentials", "shell", "error", "dump", "log", "debug",
	}
	g.Dictionaries["years"] = GenerateYears(2000, time.Now().Year()+1)
	g.Dictionaries["servers"] = []string{
		"server", "hosting", "cloud", "vps", "dedicated", "shared",
		"apache", "nginx", "iis", "tomcat", "lighttpd", "caddy",
	}
	g.Dictionaries["locations"] = []string{
		"london", "paris", "tokyo", "newyork", "berlin", "moscow", "beijing",
		"sydney", "dubai", "rio", "cairo", "rome", "madrid", "seoul", "mumbai",
		"toronto", "mexicocity", "buenosaires", "johannesburg", "vancouver",
	}

	if strings.ToLower(customWordsSource) == "auto" {
		g.loadAutoDictionaries()
	} else if customWordsSource != "" {
		fmt.Printf("⚙️ Loading custom wordlist from: %s\n", customWordsSource)
		words := LoadWordsFromFileOrURL(customWordsSource, "", g.httpClient)
		if len(words) > 0 {
			g.Dictionaries["custom_wordlist"] = words
			fmt.Printf("Loaded %d words into 'custom_wordlist'.\n", len(words))
		} else {
			fmt.Printf("❌ No words loaded from custom wordlist: %s\n", customWordsSource)
		}
	} else {
		fmt.Println("No custom wordlist or 'auto' mode selected. Loading default common dictionaries (English, Russian, Spanish, etc.).")
		g.Dictionaries["english_words"] = LoadWordsFromFileOrURL("common_english_words.txt", AutoWordlistURLs["english_auto"], g.httpClient)
		g.Dictionaries["russian_words"] = LoadWordsFromFileOrURL("common_russian_words.txt", AutoWordlistURLs["russian_auto"], g.httpClient)
		g.Dictionaries["spanish_words"] = LoadWordsFromFileOrURL("common_spanish_words.txt", AutoWordlistURLs["spanish_auto"], g.httpClient)
		g.Dictionaries["german_words"] = LoadWordsFromFileOrURL("common_german_words.txt", AutoWordlistURLs["german_auto"], g.httpClient)
	}

	g.Dictionaries["common_words"] = []string{}
	for dictName, words := range g.Dictionaries {
		if strings.HasSuffix(dictName, "_auto") || strings.HasSuffix(dictName, "_words") || dictName == "custom_wordlist" {
			g.Dictionaries["common_words"] = append(g.Dictionaries["common_words"], words...)
		}
	}

	g.Dictionaries["admin"] = append(g.Dictionaries["admin"], GenerateVariations(g.Dictionaries["admin"], 50)...)
	if len(g.Dictionaries["common_words"]) > 0 {
		g.Dictionaries["common_words"] = append(g.Dictionaries["common_words"], GenerateVariations(g.Dictionaries["common_words"], 100)...)
	}
	if len(g.Dictionaries["vulnerabilities"]) > 0 {
		g.Dictionaries["vulnerabilities"] = append(g.Dictionaries["vulnerabilities"], GenerateVariations(g.Dictionaries["vulnerabilities"], 50)...)
	}

	g.LoadCountries()
	g.LoadDomains()

	if len(g.Countries["names"]) > 0 {
		g.Dictionaries["common_words"] = append(g.Dictionaries["common_words"], g.Countries["names"]...)
	}

	return g
}

func GenerateYears(start, end int) []string {
	var years []string
	for i := start; i <= end; i++ {
		years = append(years, fmt.Sprintf("%d", i))
	}
	return years
}

func (g *DorkGenerator) LoadTemplates(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("⚠️ Warning: Failed to load templates from %s. Using default templates.\n", filename)
		return []string{"inurl:{target} {type}", "intitle:{target} {admin}"}
	}
	defer f.Close()

	var templates []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			templates = append(templates, line)
		}
	}
	if len(templates) == 0 {
		fmt.Println("⚠️ Warning: No templates found in file. Using default templates.")
		return []string{"inurl:{target} {type}", "intitle:{target} {admin}"}
	}
	return templates
}

func (g *DorkGenerator) LoadCountries() {
	g.Countries = map[string][]string{"codes": {}, "names": {}}
	resp, err := g.httpClient.Get("https://restcountries.com/v3.1/all")
	if err != nil {
		fmt.Printf("⚠️ Warning: Error loading countries (%v), using fallback list.\n", err)
		g.Countries["codes"] = []string{"us", "ru", "de", "fr", "it"}
		g.Countries["names"] = []string{"united states", "russia", "germany"}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("⚠️ Warning: Error reading country data (%v), using fallback list.\n", err)
		g.Countries["codes"] = []string{"us", "ru", "de", "fr", "it"}
		g.Countries["names"] = []string{"united states", "russia", "germany"}
		return
	}

	var data []map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("⚠️ Warning: Error parsing country data (%v), using fallback list.\n", err)
		g.Countries["codes"] = []string{"us", "ru", "de", "fr", "it"}
		g.Countries["names"] = []string{"united states", "russia", "germany"}
		return
	}

	for _, c := range data {
		if code, ok := c["cca2"].(string); ok {
			g.Countries["codes"] = append(g.Countries["codes"], strings.ToLower(code))
		}
		if name, ok := c["name"].(map[string]interface{}); ok {
			if common, ok := name["common"].(string); ok {
				g.Countries["names"] = append(g.Countries["names"], strings.ToLower(common))
			}
		}
	}
}

func (g *DorkGenerator) LoadDomains() {
	resp, err := g.httpClient.Get("https://data.iana.org/TLD/tlds-alpha-by-domain.txt")
	if err != nil {
		fmt.Printf("⚠️ Warning: Error loading domains (%v), using fallback list.\n", err)
		g.Domains = []string{"com", "net", "org", "io", "ru", "de"}
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			g.Domains = append(g.Domains, line)
		}
	}
}

func (g *DorkGenerator) GenerateDorks(target string, count int, countries, domains []string, selectedDictionaries []string) []string {
	rand.Seed(time.Now().UnixNano())
	dorks := make(map[string]struct{})
	var wg sync.WaitGroup
	var mu sync.Mutex

	targetTerms := g.GetTargetTerms(target)
	filteredCountries := g.FilterList(countries, g.Countries["codes"])
	filteredDomains := g.FilterList(domains, g.Domains)

	var dynamicWords []string
	if len(selectedDictionaries) == 0 {
		dynamicWords = append(dynamicWords, g.Dictionaries["common_words"]...)
		dynamicWords = append(dynamicWords, g.Dictionaries["locations"]...)
	} else {
		for _, dictName := range selectedDictionaries {
			if dict, ok := g.Dictionaries[dictName]; ok {
				dynamicWords = append(dynamicWords, dict...)
			} else {
				fmt.Printf("⚠️ Warning: Dictionary '%s' not found. Skipping.\n", dictName)
			}
		}
	}

	if target == "" && len(dynamicWords) > 0 {
		targetTerms = append(targetTerms, dynamicWords...)
	} else if target != "" && len(targetTerms) == 0 {
		targetTerms = append(targetTerms, dynamicWords...)
	}

	workerCount := 10
	jobs := make(chan struct{}, workerCount)

	if len(g.Templates) == 0 {
		fmt.Println("❌ Error: Template list is empty. Cannot generate Dork queries.")
		return []string{}
	}

	if len(targetTerms) == 0 &&
		len(g.Dictionaries["admin"]) == 0 &&
		len(g.Dictionaries["types"]) == 0 &&
		len(g.Dictionaries["vulnerabilities"]) == 0 &&
		len(filteredCountries) == 0 &&
		len(g.Dictionaries["years"]) == 0 &&
		len(filteredDomains) == 0 &&
		len(g.Dictionaries["servers"]) == 0 &&
		len(dynamicWords) == 0 {
		fmt.Println("❌ Error: All relevant dictionaries/word lists are empty. Cannot generate Dork queries.")
		return []string{}
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				mu.Lock()
				if len(dorks) >= count {
					mu.Unlock()
					return
				}
				mu.Unlock()

				getRand := func(list []string) string {
					if len(list) == 0 {
						return ""
					}
					return list[rand.Intn(len(list))]
				}

				template := getRand(g.Templates)
				replacements := map[string]string{
					"target":        getRand(targetTerms),
					"admin":         getRand(g.Dictionaries["admin"]),
					"type":          getRand(g.Dictionaries["types"]),
					"vulnerability": getRand(g.Dictionaries["vulnerabilities"]),
					"country":       getRand(filteredCountries),
					"year":          getRand(g.Dictionaries["years"]),
					"domain":        getRand(filteredDomains),
					"server":        getRand(g.Dictionaries["servers"]),
					"common_word":   getRand(dynamicWords),
					"location":      getRand(g.Dictionaries["locations"]),
					"english_word":  getRand(g.Dictionaries["english_auto"]),
					"russian_word":  getRand(g.Dictionaries["russian_auto"]),
					"spanish_word":  getRand(g.Dictionaries["spanish_auto"]),
					"german_word":   getRand(g.Dictionaries["german_auto"]),
					"custom_word":   getRand(g.Dictionaries["custom_wordlist"]),
				}

				dork := template
				for key, value := range replacements {
					dork = strings.ReplaceAll(dork, "{"+key+"}", value)
				}
				dork = strings.Join(strings.Fields(dork), " ")

				mu.Lock()
				dorks[dork] = struct{}{}
				mu.Unlock()
			}
		}()
	}

	maxAttempts := count * 5
	attempts := 0
	for {
		mu.Lock()
		currentDorkCount := len(dorks)
		mu.Unlock()

		if currentDorkCount >= count || attempts >= maxAttempts {
			break
		}

		select {
		case jobs <- struct{}{}:
			attempts++
		default:
			time.Sleep(5 * time.Millisecond)
		}
	}
	close(jobs)

	wg.Wait()

	var result []string
	for d := range dorks {
		result = append(result, d)
	}
	if len(result) > count {
		result = result[:count]
	}
	return result
}

func (g *DorkGenerator) GetTargetTerms(target string) []string {
	if target == "" {
		return g.Dictionaries["vulnerabilities"]
	}
	target = strings.ToLower(target)
	for _, category := range g.SpecialTargets {
		if terms, ok := category[target]; ok {
			return terms
		}
	}
	return []string{target}
}

func (g *DorkGenerator) FilterList(filter, full []string) []string {
	if len(filter) == 0 || (len(filter) == 1 && strings.ToLower(filter[0]) == "ww") {
		return full
	}
	var filteredResult []string
	for _, element := range full {
		for _, f := range filter {
			if element == f {
				filteredResult = append(filteredResult, element)
				break
			}
		}
	}
	if len(filteredResult) == 0 {
		fmt.Printf("⚠️ Warning: No matches found for filter %v. Using full list.\n", filter)
		return full
	}
	return filteredResult
}

func SaveToFile(lines []string, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("❌ Error saving to file %s: %w", filename, err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	for _, line := range lines {
		_, err := fmt.Fprintln(writer, line)
		if err != nil {
			return fmt.Errorf("❌ Error writing line to file %s: %w", filename, err)
		}
	}
	return writer.Flush()
}

func main() {
	targetFlag := flag.String("target", "", "Target (e.g., wordpress, nginx, etc.). Leaves dynamic for random search if not set.")
	countryFlag := flag.String("country", "", "Country filter (comma-separated ISO 3166-1 alpha-2 codes, e.g., 'us,de' or 'ww' for worldwide).")
	domainFlag := flag.String("domain", "", "Domain filter (comma-separated TLDs, e.g., 'com,org' or 'ww' for all).")
	quantityFlag := flag.Int("quantity", 1000, "Number of Dork queries to generate.")
	outputFileFlag := flag.String("output", "dorks.txt", "File to save generated Dork queries.")
	listTargetsFlag := flag.Bool("list-targets", false, "Show list of available targets.")
	listDictionariesFlag := flag.Bool("list-dictionaries", false, "Show list of available dictionaries.")
	dictionariesFlag := flag.String("dictionaries", "", "Comma-separated list of dictionaries to use (e.g., 'admin,types'). Use 'all' for all built-in dictionaries or 'auto_langs' for auto-downloaded languages.")
	customWordsSourceFlag := flag.String("custom_words", "", "Path to a local wordlist file, or 'auto' to download common language wordlists from GitHub.")

	flag.Parse()

	if *quantityFlag <= 0 {
		fmt.Println("❌ Error: Number of Dork queries must be a positive number.")
		os.Exit(1)
	}

	generator := NewDorkGenerator(*customWordsSourceFlag)

	if *listTargetsFlag {
		fmt.Println("\n🔍 Available targets:")
		for category, targets := range generator.SpecialTargets {
			fmt.Printf("Category: %s\n", category)
			for targetName := range targets {
				fmt.Printf(" - %s\n", targetName)
			}
		}
		return
	}

	if *listDictionariesFlag {
		fmt.Println("\n📚 Available dictionaries:")
		for dictName := range generator.Dictionaries {
			if !strings.HasSuffix(dictName, "_auto") {
				fmt.Printf(" - %s (%d words)\n", dictName, len(generator.Dictionaries[dictName]))
			}
		}
		fmt.Println("\nNote: Dictionaries ending with '_auto' are loaded in 'auto' mode but not directly selectable by name unless you want to use them specifically.")
		fmt.Println("You can also use 'auto_langs' with --dictionaries to include all auto-downloaded language wordlists.")
		return
	}

	inputCountries := []string{}
	if *countryFlag != "" {
		inputCountries = strings.Split(*countryFlag, ",")
		for i, c := range inputCountries {
			inputCountries[i] = strings.ToLower(strings.TrimSpace(c))
		}
	}

	inputDomains := []string{}
	if *domainFlag != "" {
		inputDomains = strings.Split(*domainFlag, ",")
		for i, d := range inputDomains {
			inputDomains[i] = strings.ToLower(strings.TrimSpace(d))
		}
	}

	selectedDictionaries := []string{}
	if *dictionariesFlag != "" {
		if strings.ToLower(*dictionariesFlag) == "all" {
			for dictName := range generator.Dictionaries {
				if !strings.HasSuffix(dictName, "_auto") {
					selectedDictionaries = append(selectedDictionaries, dictName)
				}
			}
			for lang := range AutoWordlistURLs {
				selectedDictionaries = append(selectedDictionaries, lang)
			}

		} else if strings.ToLower(*dictionariesFlag) == "auto_langs" {
			for lang := range AutoWordlistURLs {
				selectedDictionaries = append(selectedDictionaries, lang)
			}
		} else {
			selectedDictionaries = strings.Split(*dictionariesFlag, ",")
			for i, d := range selectedDictionaries {
				selectedDictionaries[i] = strings.ToLower(strings.TrimSpace(d))
			}
		}
	}

	start := time.Now()
	dorks := generator.GenerateDorks(*targetFlag, *quantityFlag, inputCountries, inputDomains, selectedDictionaries)
	duration := time.Since(start)

	fmt.Printf("✅ Generated %d Dork queries in %.2f seconds\n", len(dorks), duration.Seconds())
	if err := SaveToFile(dorks, *outputFileFlag); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("💾 Results saved to:", *outputFileFlag)
}
