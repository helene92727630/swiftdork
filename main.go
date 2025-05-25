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

var bip39EnglishWords = []string{
	"abandon", "ability", "able", "about", "above", "absent", "absorb", "abstract", "absurd", "abuse",
	"access", "acid", "acoustic", "across", "act", "action", "active", "actor", "actual", "adapt",
	"add", "addict", "address", "adjust", "admit", "adult", "advance", "advice", "aerobic", "affair",
	"affect", "affordable", "afraid", "again", "age", "agent", "agree", "ahead", "aim", "air",
	"airport", "aisle", "alarm", "album", "alcohol", "alert", "alien", "all", "alley", "allow",
	"almost", "alone", "alpha", "already", "also", "alter", "always", "amateur", "amazing", "ambition",
	"ambush", "amend", "amino", "among", "amount", "amuse", "analyst", "anchor", "ancient", "anger",
	"angle", "angry", "animal", "ankle", "announce", "annual", "another", "answer", "antenna", "antique",
	"open", "opera", "opinion", "oppose", "option", "orange", "orbit", "orchard", "order", "ordinary",
	"organize", "orient", "original", "orphan", "ostrich", "other", "outdoor", "outer", "output", "outside",
	"oval", "oven", "over", "own", "owner", "oyster", "ozone", "pair", "palm", "panel",
	"panic", "panther", "paper", "parade", "parent", "park", "parrot", "party", "pass", "patch",
	"path", "patient", "patrol", "pattern", "pause", "pave", "payment", "peace", "peanut", "pear",
	"peasant", "pelican", "pen", "penalty", "pencil", "people", "pepper", "perfect", "perfume", "permit",
	"person", "pet", "phantom", "phase", "photo", "phrase", "physical", "piano", "picnic", "picture",
	"piece", "pig", "pigeon", "pill", "pilot", "pink", "pioneer", "pipe", "pistol", "pitch",
	"pizza", "place", "plain", "planet", "plastic", "plate", "play", "please", "pledge", "plug",
	"plus", "poem", "poet", "point", "polar", "pole", "police", "power", "practice", "praise",
	"predict", "prefer", "prepare", "present", "pretty", "prevent", "price", "pride", "primary", "print",
	"priority", "prison", "private", "prize", "problem", "process", "produce", "profit", "program", "project",
	"promote", "proof", "property", "prosper", "protect", "proud", "provide", "public", "pudding", "pull",
	"pulp", "pulse", "punch", "pupil", "puppy", "purchase", "purity", "purpose", "purse", "push",
	"put", "puzzle", "pyramid", "quality", "quantum", "quarter", "question", "quick", "quit", "quiz",
	"quote", "rabbit", "race", "rack", "radar", "radio", "rail", "rain", "raise", "rally",
	"ramp", "ranch", "random", "range", "rapid", "rare", "rate", "rather", "raven", "raw",
	"reach", "read", "ready", "real", "reason", "rebel", "rebuild", "recall", "receive", "recipe",
	"record", "recycle", "reduce", "reflect", "reform", "refuse", "region", "regret", "regular", "reject",
	"relax", "release", "relief", "rely", "remain", "remember", "remind", "remove", "render", "renew",
	"rent", "repair", "repeat", "replace", "report", "represent", "republic", "require", "rescue", "research",
	"resist", "resource", "response", "result", "retire", "retreat", "return", "reveal", "review", "reward",
	"rhythm", "ribbon", "rice", "ridge", "rifle", "right", "rigid", "ring", "riot", "ripple",
	"risk", "ritual", "rival", "river", "road", "roast", "robot", "robust", "rocket", "romance",
	"roof", "rookie", "room", "rose", "rotate", "rough", "round", "route", "royal", "rubber",
	"rude", "rug", "rule", "run", "runway", "rural", "sad", "saddle", "sadness", "safe",
	"sail", "salad", "salmon", "salon", "salt", "salute", "same", "sample", "sand", "satisfy",
	"satoshi", "sauce", "save", "say", "scale", "scan", "scare", "scatter", "scene", "scheme",
	"school", "science", "scissors", "score", "scorpion", "scout", "scrap", "screen", "script", "scrub",
	"sea", "search", "season", "seat", "second", "secret", "section", "security", "see", "seed",
	"seek", "segment", "select", "sell", "send", "sense", "sentence", "series", "service", "session",
	"set", "settle", "setup", "seven", "shadow", "shaft", "shallow", "share", "shed", "shell",
	"sheriff", "shield", "shift", "shine", "ship", "shiver", "shock", "shoe", "shoot", "shop",
	"short", "shoulder", "shove", "shrimp", "shrink", "sibling", "side", "sidewalk", "siege", "sight",
	"sign", "silent", "silk", "silly", "silver", "similar", "simple", "since", "sing", "siren",
	"sister", "situate", "six", "size", "skate", "sketch", "ski", "skill", "skin", "skirt",
	"slack", "slave", "sleep", "slim", "slip", "slope", "slot", "slow", "slush", "small",
	"smart", "smile", "smoke", "smooth", "snack", "snail", "snake", "sneak", "sneeze", "sniff",
	"spirit", "split", "spoil", "sponsor", "spoon", "sport", "spot", "spray", "spread", "spring",
	"spy", "square", "squeeze", "squirrel", "stable", "stadium", "staff", "stage", "stairs", "stamp",
	"stand", "start", "state", "stay", "steak", "steel", "stem", "step", "stereo", "stick",
	"still", "sting", "stomach", "stone", "stool", "story", "stove", "strategy", "street", "strike",
	"strong", "struggle", "student", "stuff", "stumble", "style", "subject", "submerge", "submit", "subsidy",
	"suburb", "success", "such", "sudden", "suffer", "sugar", "suggest", "suit", "summer", "sun",
	"sunny", "sunset", "super", "supply", "supreme", "sure", "surface", "surge", "surprise", "surround",
	"survey", "suspect", "sustain", "swallow", "swamp", "swap", "swarm", "swear", "sweet", "swift",
	"swim", "swing", "switch", "sword", "symbol", "symptom", "syrup", "system", "table", "tackle",
	"tag", "talent", "talk", "tank", "tape", "target", "task", "taste", "taxi", "teach",
	"team", "tell", "temple", "tenant", "tennis", "tent", "term", "test", "text", "thank",
	"that", "the", "then", "theory", "there", "they", "thing", "this", "thought", "three",
	"thrive", "throw", "thumb", "thunder", "ticket", "tide", "tiger", "tilt", "timber", "time",
	"tiny", "tip", "tired", "tissue", "title", "toast", "tobacco", "today", "toddler", "toe",
	"together", "toilet", "token", "tomato", "tomorrow", "tone", "tongue", "tooth", "top", "topic",
	"topple", "torch", "tornado", "tortoise", "toss", "total", "tourist", "toward", "tower", "town",
	"toy", "track", "trade", "traffic", "tragic", "train", "transfer", "trap", "travel", "tray",
	"treat", "tree", "trend", "trial", "tribe", "trick", "trigger", "trim", "trip", "trophy",
	"trouble", "true", "truly", "trumpet", "trust", "truth", "try", "tsunami", "tube", "tuition",
	"tumble", "tuna", "tunnel", "turkey", "turn", "turtle", "twelve", "twenty", "twice", "twin",
	"type", "ugly", "umbrella", "uncover", "under", "undo", "unfold", "unhappy", "uniform", "unique",
	"unit", "universe", "unknown", "unlock", "until", "unusual", "unveil", "update", "upgrade", "uphold",
	"upon", "upper", "upset", "urban", "urge", "usage", "use", "used", "useful", "useless",
	"usual", "utility", "vacant", "vacuum", "vague", "valid", "valley", "valve", "van", "vanish",
	"vapor", "various", "vast", "vault", "vehicle", "velvet", "vendor", "venture", "venue", "verb",
	"verify", "version", "very", "vessel", "veteran", "viable", "vibrant", "vicious", "victory", "video",
	"view", "village", "vintage", "violin", "virtual", "virus", "visa", "visit", "visual", "vital",
	"vivid", "vocal", "voice", "void", "volcano", "volume", "vote", "vowel", "voyage", "wage",
	"wagon", "wait", "walk", "wall", "walrus", "want", "war", "wardrobe", "warm", "warning",
	"wash", "waste", "water", "wave", "way", "we", "wealth", "wear", "web", "wedding",
	"weekend", "weird", "welcome", "west", "wet", "whale", "what", "wheat", "wheel", "when",
	"where", "while", "whisper", "wide", "widget", "wisdom", "wise", "wish", "witness", "wolf",
	"woman", "wonder", "wood", "wool", "word", "work", "world", "worry", "worth", "wrap",
	"wreck", "wrestle", "wrist", "write", "wrong", "yard", "year", "yellow", "you", "young",
	"yourself", "zero", "zigzag", "zone", "zoo", "zombie",
}

func GenerateVariations(baseWords []string, count int) []string {
	if len(baseWords) == 0 {
		return []string{}
	}

	var variations []string
	variations = append(variations, baseWords...)

	suffixes := []string{"_backup", "_conf", "_log", "_old", "_dev", "_test", "-db", "s", "es", "data"}
	prefixes := []string{"dev_", "test_", "old_", "new_"}
	numbers := []string{"1", "123", "2023", "2024", "2025", "00", "01"}

	for len(variations) < count {
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
		if len(variations) >= count {
			break
		}
	}
	return variations
}

func LoadWordsFromFile(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("⚠️ Внимание: Не удалось загрузить слова из %s. Этот словарь будет пропущен.\n", filename)
		return []string{}
	}
	defer f.Close()

	var words []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			words = append(words, strings.ToLower(line))
		}
	}
	return words
}

func NewDorkGenerator() *DorkGenerator {
	g := &DorkGenerator{
		Dictionaries: map[string][]string{
			"admin": {
				"admin", "login", "panel", "control", "dashboard", "manager", "backend",
				"console", "secure", "signin", "cpanel", "administrator",
			},
			"types": {
				"php", "asp", "aspx", "jsp", "cfm", "html", "htm", "xml", "json", "sql",
				"txt", "doc", "pdf", "xls", "csv", "bak", "zip",
			},
			"vulnerabilities": {
				"vulnerable", "exploit", "inurl", "intitle", "intext", "index of",
				"password", "config", "backup", "database", "leak", "exposed",
				"credentials", "shell", "error", "dump", "log", "debug",
			},
			"years": GenerateYears(2000, time.Now().Year()),
			"servers": {
				"server", "hosting", "cloud", "vps", "dedicated", "shared",
				"apache", "nginx", "iis", "tomcat", "lighttpd", "caddy",
			},
            "locations": {
                "london", "paris", "tokyo", "newyork", "berlin", "moscow", "beijing",
                "sydney", "dubai", "rio", "cairo", "rome", "madrid", "seoul", "mumbai",
                "toronto", "mexicocity", "buenosaires", "johannesburg", "vancouver",
            },
		},
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
            Timeout: 10 * time.Second,
        },
	}
	g.Templates = g.LoadTemplates("base_templates.txt")

    g.Dictionaries["english_words"] = bip39EnglishWords

    g.Dictionaries["russian_words"] = LoadWordsFromFile("common_russian_words.txt")

    g.Dictionaries["common_words"] = []string{}
    g.Dictionaries["common_words"] = append(g.Dictionaries["common_words"], bip39EnglishWords...)
    g.Dictionaries["common_words"] = append(g.Dictionaries["common_words"], g.Dictionaries["russian_words"]...)
    
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
		fmt.Printf("⚠️ Внимание: Не удалось загрузить шаблоны из %s. Используются шаблоны по умолчанию.\n", filename)
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
	return templates
}

func (g *DorkGenerator) LoadCountries() {
	g.Countries = map[string][]string{"codes": {}, "names": {}}
	resp, err := g.httpClient.Get("https://restcountries.com/v3.1/all")
	if err != nil {
		fmt.Printf("⚠️ Внимание: Ошибка загрузки стран (%v), используется резервный список.\n", err)
		g.Countries["codes"] = []string{"us", "ru", "de", "fr", "it"}
		g.Countries["names"] = []string{"united states", "russia", "germany"}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("⚠️ Внимание: Ошибка чтения данных стран (%v), используется резервный список.\n", err)
		g.Countries["codes"] = []string{"us", "ru", "de", "fr", "it"}
		g.Countries["names"] = []string{"united states", "russia", "germany"}
		return
	}

	var data []map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("⚠️ Внимание: Ошибка парсинга данных стран (%v), используется резервный список.\n", err)
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
		fmt.Printf("⚠️ Внимание: Ошибка загрузки доменов (%v), используется резервный список.\n", err)
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

func (g *DorkGenerator) GenerateDorks(target string, count int, countries, domains []string) []string {
	rand.Seed(time.Now().UnixNano())
	dorks := make(map[string]struct{})
	var wg sync.WaitGroup
	var mu sync.Mutex

	targetTerms := g.GetTargetTerms(target)
	filteredCountries := g.FilterList(countries, g.Countries["codes"])
	filteredDomains := g.FilterList(domains, g.Domains)

	if target == "" {
        targetTerms = append(targetTerms, g.Dictionaries["common_words"]...)
        targetTerms = append(targetTerms, g.Dictionaries["locations"]...)
        targetTerms = append(targetTerms, g.Dictionaries["english_words"]...)
        targetTerms = append(targetTerms, g.Dictionaries["russian_words"]...)
	}

	workerCount := 10
	jobs := make(chan struct{}, workerCount)

    if len(g.Templates) == 0 {
        fmt.Println("❌ Ошибка: Список шаблонов пуст. Невозможно сгенерировать Dork-запросы.")
        return []string{}
    }
    if len(targetTerms) == 0 && len(g.Dictionaries["admin"]) == 0 &&
       len(g.Dictionaries["types"]) == 0 && len(g.Dictionaries["vulnerabilities"]) == 0 &&
       len(filteredCountries) == 0 && len(g.Dictionaries["years"]) == 0 &&
       len(filteredDomains) == 0 && len(g.Dictionaries["servers"]) == 0 &&
       len(g.Dictionaries["common_words"]) == 0 && len(g.Dictionaries["locations"]) == 0 &&
       len(g.Dictionaries["english_words"]) == 0 && len(g.Dictionaries["russian_words"]) == 0 {
        fmt.Println("❌ Ошибка: Все словари/списки слов пусты. Невозможно сгенерировать Dork-запросы.")
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
                    "common_word":   getRand(g.Dictionaries["common_words"]),
                    "location":      getRand(g.Dictionaries["locations"]),
                    "english_word":  getRand(g.Dictionaries["english_words"]),
                    "russian_word":  getRand(g.Dictionaries["russian_words"]),
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

    for {
        mu.Lock()
        currentDorkCount := len(dorks)
        mu.Unlock()

        if currentDorkCount >= count {
            break
        }

        select {
        case jobs <- struct{}{}:
        default:
            time.Sleep(10 * time.Millisecond)
        }
    }
    close(jobs)

	wg.Wait()

	var result []string
	for d := range dorks {
		result = append(result, d)
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
		fmt.Printf("⚠️ Внимание: Не найдено совпадений для фильтра %v. Используется полный список.\n", filter)
		return full
	}
	return filteredResult
}

func SaveToFile(lines []string, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("❌ Ошибка сохранения в файл: %w", err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	for _, line := range lines {
		_, err := fmt.Fprintln(writer, line)
		if err != nil {
			return fmt.Errorf("❌ Ошибка записи строки в файл: %w", err)
		}
	}
	return writer.Flush()
}

func main() {
	targetFlag := flag.String("цель", "", "Цель (например, wordpress, nginx и т.д.)")
	countryFlag := flag.String("страна", "", "Фильтр стран (через запятую или 'ww')")
	domainFlag := flag.String("домен", "", "Фильтр доменов (через запятую или 'ww')")
	quantityFlag := flag.Int("количество", 1000, "Количество Dork-запросов для генерации")
	outputFileFlag := flag.String("файл", "дорки.txt", "Файл для сохранения результатов")
	listTargetsFlag := flag.Bool("список", false, "Показать список доступных целей")
	flag.Parse()

	if *quantityFlag <= 0 {
		fmt.Println("❌ Ошибка: Количество Dork-запросов должно быть положительным числом.")
		os.Exit(1)
	}

	генератор := NewDorkGenerator()

	if *listTargetsFlag {
		fmt.Println("\n🔍 Доступные цели:")
		for категория, цели := range генератор.SpecialTargets {
			fmt.Printf("Категория: %s\n", категория)
			for цель := range цели {
				fmt.Printf(" - %s\n", цель)
			}
		}
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

	начало := time.Now()
	дорки := генератор.GenerateDorks(*targetFlag, *quantityFlag, inputCountries, inputDomains)
	длительность := time.Since(начало)

	fmt.Printf("✅ Сгенерировано %d Dork-запросов за %.2f секунд\n", len(дорки), длительность.Seconds())
	if err := SaveToFile(дорки, *outputFileFlag); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("💾 Результаты сохранены в:", *outputFileFlag)
}
