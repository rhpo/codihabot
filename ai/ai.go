package ai

// ai/generate.go
import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"instabot/postgen"

	"google.golang.org/genai"

	"math/rand"
)

const ISDEV = false

const testStr = `{
  "color": "#FF2D20",
  "name": "laravel-php-framework",
  "caption": "Build robust web apps effortlessly with Laravel, the elegant PHP framework. Speed up your development today!",
  "slides": {
    "hero": {
      "title": "Laravel: Elegant PHP Web Framework",
      "description": "Rapidly develop powerful web applications with Laravel's expressive syntax and rich ecosystem.",
      "image": "Laravel framework logo TRANSPARENT"
    },
    "info": [
      {
        "title": "What is Laravel?",
        "content": "Laravel is a powerful and elegant PHP web application framework known for its expressive, clean syntax. It adheres to the Model-View-Controller (MVC) architectural pattern, providing a structured approach to web development. One of its main goals is to make the development process enjoyable for developers by easing common tasks required in most web projects, such as authentication, routing, sessions, and caching. Its comprehensive features and robust ecosystem make it a go-to choice for building anything from small personal projects to large enterprise-level applications.<br><br>Laravel simplifies complex operations, allowing developers to focus more on creating unique features rather than reinventing the wheel. For instance, creating a new user authentication system takes just a few commands using Laravel Breeze or Jetstream, vastly reducing development time compared to building it from scratch. You can define routes like"
      },
      {
        "title": "Key Features & Benefits",
        "content": "Laravel boasts a rich set of features that significantly enhance development efficiency and application performance. Its built-in Artisan command-line interface provides numerous helpful commands, like generating boilerplate code (php artisan make:controller) or running database migrations (php artisan migrate). The Blade templating engine offers powerful, yet simple, ways to build views with features like template inheritance and component rendering.<br><br>Security is paramount, and Laravel provides robust protections against common web vulnerabilities, including CSRF and XSS attacks, right out of the box. Middleware offers a convenient mechanism for filtering HTTP requests entering your application, perfect for authentication or CORS handling. Additionally, Laravel's robust queue system and caching drivers help optimize performance for high-traffic applications. This comprehensive toolkit ensures that developers can build secure, scalable, and maintainable applications with ease, fostering a productive and enjoyable coding experience."
      }
    ],
    "finish": {
      "description": "Unlock your full potential as a web developer. Laravel simplifies complex tasks, empowering you to create amazing applications efficiently.",
      "image": "Laravel framework logo TRANSPARENT"
    }
  }
}`

var topics = [...]string{"aarch64", "adonisjs", "aerospike", "aframe", "aftereffects", "akka", "algolia", "almalinux", "alpinejs", "amazonwebservices", "anaconda", "android", "androidstudio", "angular", "angularjs", "angularmaterial", "ansible", "ansys", "antdesign", "apache", "apacheairflow", "apachekafka", "apachespark", "apex", "apl", "apollographql", "appcelerator", "apple", "appwrite", "archlinux", "arduino", "argocd", "artixlinux", "astro", "atom", "awk", "axios", "azure", "azuredevops", "azuresqldatabase", "babel", "babylonjs", "backbonejs", "ballerina", "bamboo", "bash", "bazel", "beats", "behance", "bevyengine", "biome", "bitbucket", "blazor", "blender", "bootstrap", "bower", "browserstack", "bulma", "bun", "c", "cairo", "cakephp", "canva", "capacitor", "carbon", "cassandra", "centos", "ceylon", "chakraui", "chartjs", "chrome", "circleci", "clarity", "clickhouse", "clion", "clojure", "clojurescript", "cloudflare", "cloudflareworkers", "cloudrun", "cmake", "cobol", "codeac", "codecov", "codeigniter", "codepen", "coffeescript", "composer", "confluence", "consul", "contao", "corejs", "cosmosdb", "couchbase", "couchdb", "cpanel", "cplusplus", "crystal", "csharp", "css3", "cucumber", "cypressio", "d3js", "dart", "datadog", "datagrip", "dataspell", "datatables", "dbeaver", "debian", "delphi", "denojs", "detaspace", "devicon", "digitalocean", "discloud", "discordjs", "django", "djangorest", "docker", "doctrine", "dot-net", "dotnetcore", "dovecot", "dreamweaver", "dropwizard", "drupal", "duckdb", "dyalog", "dynamodb", "dynatrace", "eclipse", "ecto", "elasticsearch", "electron", "eleventy", "elixir", "elm", "emacs", "embeddedc", "ember", "entityframeworkcore", "envoy", "erlang", "eslint", "expo", "express", "facebook", "fastapi", "fastify", "faunadb", "feathersjs", "fedora", "fiber", "figma", "filamentphp", "filezilla", "firebase", "firebird", "firefox", "flask", "flutter", "forgejo", "fortran", "foundation", "framermotion", "framework7", "fsharp", "fusion", "gardener", "gatling", "gatsby", "gazebo", "gcc", "gentoo", "ghost", "gimp", "git", "gitbook", "github", "githubactions", "githubcodespaces", "gitkraken", "gitlab", "gitpod", "gitter", "gleam", "glitch", "go", "godot", "goland", "google", "googlecloud", "googlecolab", "gradle", "grafana", "grails", "graphql", "groovy", "grpc", "grunt", "gulp", "hadoop", "handlebars", "harbor", "hardhat", "harvester", "haskell", "haxe", "helm", "heroku", "hibernate", "homebrew", "hoppscotch", "html5", "htmx", "hugo", "hyperv", "ie10", "ifttt", "illustrator", "inertiajs", "influxdb", "inkscape", "insomnia", "intellij", "ionic", "jaegertracing", "jamstack", "jasmine", "java", "javascript", "jeet", "jekyll", "jenkins", "jest", "jetbrains", "jetpackcompose", "jhipster", "jira", "jiraalign", "jquery", "json", "jule", "julia", "junit", "jupyter", "k3os", "k3s", "k6", "kaggle", "kaldi", "kalilinux", "karatelabs", "karma", "kdeneon", "keras", "kibana", "knexjs", "knockout", "kotlin", "krakenjs", "ktor", "kubeflow", "kubernetes", "labview", "laminas", "laravel", "laraveljetstream", "latex", "leetcode", "less", "libgdx", "linkedin", "linux", "linuxmint", "liquibase", "livewire", "llvm", "lodash", "logstash", "love2d", "lua", "lumen", "magento", "mapbox", "mariadb", "markdown", "materializecss", "materialui", "matlab", "matplotlib", "mattermost", "maven", "maya", "memcached", "mercurial", "meteor", "microsoftsqlserver", "minitab", "mithril", "mobx", "mocha", "modx", "moleculer", "mongodb", "mongoose", "monogame", "moodle", "msdos", "mysql", "nano", "nats", "neo4j", "neovim", "nestjs", "netbeans", "netbox", "netlify", "networkx", "newrelic", "nextjs", "nginx", "ngrok", "ngrx", "nhibernate", "nim", "nimble", "nixos", "nodejs", "nodemon", "nodered", "nodewebkit", "nomad", "norg", "notion", "npm", "npss", "nuget", "numpy", "nuxt", "nuxtjs", "oauth", "objectivec", "ocaml", "ohmyzsh", "okta", "openal", "openapi", "opencl", "opencv", "opengl", "openstack", "opensuse", "opentelemetry", "opera", "oracle", "ory", "p5js", "packer", "pandas", "passport", "perl", "pfsense", "phalcon", "phoenix", "photonengine", "photoshop", "php", "phpstorm", "pixijs", "playwright", "plotly", "pm2", "pnpm", "podman", "poetry", "polygon", "portainer", "postcss", "postgresql", "postman", "powershell", "premierepro", "primeng", "prisma", "processing", "processwire", "prolog", "prometheus", "protractor", "proxmox", "pug", "pulsar", "pulumi", "puppeteer", "purescript", "putty", "pycharm", "pypi", "pyscript", "pytest", "python", "pytorch", "qodana", "qt", "qtest", "quarkus", "quasar", "qwik", "r", "rabbitmq", "racket", "radstudio", "rails", "railway", "rancher", "raspberrypi", "reach", "react", "reactbootstrap", "reactnative", "reactnavigation", "reactrouter", "readthedocs", "realm", "rect", "redhat", "redis", "redux", "reflex", "remix", "renpy", "replit", "rexx", "rider", "rocksdb", "rockylinux", "rollup", "ros", "rspec", "rstudio", "ruby", "rubymine", "rust", "rxjs", "safari", "salesforce", "sanity", "sass", "scala", "scalingo", "scikitlearn", "sdl", "selenium", "sema", "sentry", "sequelize", "shopware", "shotgrid", "sketch", "slack", "socketio", "solidity", "solidjs", "sonarqube", "sourceengine", "sourcetree", "spack", "spicedb", "splunk", "spring", "spss", "spyder", "sqlalchemy", "sqldeveloper", "sqlite", "ssh", "stackblitz", "stackoverflow", "stata", "stenciljs", "storybook", "streamlit", "styledcomponents", "stylus", "subversion", "sulu", "supabase", "surrealdb", "svelte", "svgo", "swagger", "swift", "swiper", "symfony", "tailwindcss", "talos", "tauri", "teleport", "tensorflow", "terraform", "terramate", "tex", "thealgorithms", "threedsmax", "threejs", "thymeleaf", "titaniumsdk", "tmux", "tomcat", "tortoisegit", "towergit", "traefikmesh", "traefikproxy", "travis", "trello", "trpc", "turbo", "twilio", "twitter", "typescript", "typo3", "ubuntu", "unifiedmodelinglanguage", "unity", "unix", "unrealengine", "uwsgi", "v8", "vaadin", "vagrant", "vala", "vault", "veevalidate", "vercel", "vertx", "vim", "visualbasic", "visualstudio", "vite", "vitejs", "vitess", "vitest", "vscode", "vscodium", "vsphere", "vuejs", "vuestorefront", "vuetify", "vulkan", "vyper", "waku", "wasm", "web3js", "webflow", "webgpu", "weblate", "webpack", "webstorm", "windows11", "windows8", "wolfram", "woocommerce", "wordpress", "xamarin", "xcode", "xd", "xml", "yaml", "yarn", "yii", "yugabytedb", "yunohost", "zend", "zig", "zsh", "zustand"}

const AIprompt = `⚠️ VITAL RULE — READ FIRST AND NEVER IGNORE:  
The topic ID is **predefined** and given as: [%s].  
- It must be treated as exactly **one word, lowercase**.  
- Do not alter, replace, or invent any other topic.  
- The topic ALWAYS comes from the **tech lexicon**.  
  Example: if the topic is "insomnia", it refers to the **API platform Insomnia**, not the medical condition.  
This restriction is vitally important: if the topic differs from [%s] or is interpreted outside the tech lexicon, the output is invalid.  

⚠️ IMPORTANT (CONTENT RULES):  
- The slides are for Instagram posts.  
- Content must be minimized for engagement.  
- BUT it must not be too short: enough detail to teach and feel serious.  
- All content must be coherent, professional, aesthetic, and beginner-friendly.  
- Avoid strange characters or filler.  
- Each slide must deliver a clear, useful takeaway.  

Color logic:  
- "color" = dimmed version of a color that represents the topic.  
- Must guarantee readability with white text.  
- Never use random or unrelated colors.  

Markdown rules:  
- Use #### for subtitles.  
- Use * for flat bullet points (no nesting).  
- Use \n for new lines, \n\n for paragraph breaks.  
- Use <br><br> between text and code blocks.  

Code rules:  
- At least one info slide must include a short code snippet (max 7 lines).  
- Code must be correct, simple, and relevant.  
- Code must inspire learning (show ease or usefulness).  

Slide rules:  
- Keep explanations clear, general, and beginner-friendly.  
- Avoid long text blocks: compact, readable, engaging.  
- Use minimal words, but keep depth to feel serious and valuable.  
- Each info slide must teach one idea.  
- Slides must focus on *learning*, not just promotion.  

JSON rules:  
- Output must be valid JSON only.  
- Do NOT minify JSON — keep it readable.  
- Values must use the least words possible while staying clear, coherent, and useful.  
- Follow this schema exactly:  

{
  "color": string (hex, dimmed, representative),
  "name": string (topic id, lowercase, must match [%s]),
  "caption": string (Instagram caption, multi-line with \n, must include hashtags),
  "slides": {
    "hero": {
      "title": string (≤7 words, markdown, no colon),
      "description": string (≤20 words),
      "image": ""
    },
    "info": [
      {
        "title": string,
        "content": string (markdown, minimized but not too short, coherent, engaging, educational)
      },
      {
        "title": string,
        "content": string
      },
      {
        "title": string,
        "content": string
      }
    ],
    "finish": {
      "description": string,
      "image": ""
    }
  }
}

`

/*func filterTopics(topics []string, alreadyTalked []string) []string {
	// put alreadyTalked into a set for quick lookup
	exclude := make(map[string]struct{}, len(alreadyTalked))
	for _, t := range alreadyTalked {
		exclude[t] = struct{}{}
	}

	// keep only the ones not in exclude
	result := make([]string, 0, len(topics))
	for _, t := range topics {
		if _, found := exclude[t]; !found {
			result = append(result, t)
		}
	}
	return result
}*/

// FilterTopics returns a slice of strings from A that are not present in B.
func FilterTopics(A, B []string) []string {
	// put B into a set (map for fast lookup)
	setB := make(map[string]struct{})
	for _, b := range B {
		setB[b] = struct{}{}
	}

	// filter A
	result := []string{}
	for _, a := range A {
		if _, found := setB[a]; !found {
			result = append(result, a)
		}
	}
	return result
}


func GeneratePost(ctx context.Context) (postgen.Post, error) {

	var out string
	var topic string

	if ISDEV {
		out = testStr
	} else {
		apiKey := os.Getenv("GEMINI_API_KEY")

		if apiKey == "" {
			return postgen.Post{}, fmt.Errorf("GEMINI_API_KEY not set")
		}

		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey: apiKey,
		})
		if err != nil {
			return postgen.Post{}, err
		}
		alreadyTalkedAbout := make([]string, 0)

		entries, err := os.ReadDir("./posts")
		if err != nil {
			panic(err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				alreadyTalkedAbout = append(alreadyTalkedAbout, entry.Name())
			}
		}

		// Filter topics
		filtered := FilterTopics(topics[:], alreadyTalkedAbout)

		// prompt := fmt.Sprintf(AIprompt, strings.Join(filtered, ", "), strings.Join(alreadyTalkedAbout, ", "))
		topic = filtered[rand.Intn(len(filtered))]

		prompt := fmt.Sprintf(AIprompt, topic, topic, topic)

		// fmt.Println("Already talked about: ", strings.Join(alreadyTalkedAbout, ", "))
		// fmt.Println("Available Topics: ", strings.Join(filtered[0:], ", "))

		model := "gemini-2.5-flash-lite"

		resp, err := client.Models.GenerateContent(ctx, model, genai.Text(prompt), &genai.GenerateContentConfig{})
		if err != nil {
			return postgen.Post{}, fmt.Errorf("failed to generate content: %w", err)
		} else {
			out = resp.Candidates[0].Content.Parts[0].Text
		}
	}

	out = strings.TrimPrefix(out, "```json")
	out = strings.TrimSuffix(out, "```")

	var p postgen.Post
	if err := json.Unmarshal([]byte(out), &p); err != nil {
		return postgen.Post{}, fmt.Errorf("failed to parse model JSON: %w\nraw output: %s", err, out)
	}

	// assigning images

	// var err error
	// imageClient := images.NewClient(os.Getenv("SEARCH_ENGINE_API_KEY"), os.Getenv("SEARCH_ENGINE_ID"))
	// p.Slides.Hero.Image, err = imageClient.SearchPNGImage(p.Slides.Hero.Image)
	// p.Slides.Finish.Image = p.Slides.Hero.Image
	// 	if err != nil {
	// 		return postgen.Post{}, fmt.Errorf("failed to find image: %w", err)
	// 	}

	// finish.image is hero.image is https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/{post.name}/{post.name}-original.svg
	p.Slides.Hero.Image = fmt.Sprintf("https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/%s/%s-original.svg", p.Name, p.Name)
	p.Slides.Finish.Image = p.Slides.Hero.Image

	println("Image: ", p.Slides.Hero.Image)

	// Convert markdown to HTML with your existing postgen logic
	p = postgen.NewPost(p)

	p.Name = topic

	return p, nil
}
