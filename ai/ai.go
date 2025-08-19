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

const AIprompt = `
Choose any random topic ID (ex. language/framework etc...) topic FROM THIS LIST ONLY: %s

INSIDE CODEBLOCKS USE \n INSTEAD OF <br>
AND USE \n\n TO SEPARATE BETWEEN PARAGRAPHS (DOUBLE \n)
USE THE #### (h4) IN MARKDOWN SOMETIMES TO MARK SUBTITLES OR POINTS
YOU MUST ALSO NOT WRITE A BIG CHUNK OF TEXT, BUT RATHER FOCUS ON KEY POINTS AND ESSENTIAL INFORMATION. MAKE PARAGRAPHS SEPARATED BY <br> TO ENHANCE VISIBILITY
AFTER WRITING A "*" (to indicate a point on a list) ALWAYS MAKE SURE YOU ADD A SPACE BETWEEN IT AND THE WORD

IF YOU TALK ABOUT A PROGRAMMING LANGUAGE, TRY TO INCLUDE A REALLY SHORT AND SIMPLE CODE SNIPPET TO
MOTIVATE THE USER (MAXIMUM: 7 LINES OF CODE) TO USE THAT PROGRAMMING LANGUAGE, LIKE:

#include <stdio.h>

int main(char** argv[]) {
	printf("hello world!");
}

OR

package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello world!");
}

AND REMEMBER: ONLY REPLY WITH THE JSON AND WITHOUT ` + "```" + `

ONE INFO SLIDE MUST CONTAIN CODE EXAMPLES.

+ USE <br><br> BETWEEN TEXT AND CODEBLOCKS OR VICE VERSA

DONT MAKE NESTES LISTS (YOU CAN ONLY MAKE SIMPLE LISTS) LIKE
* list
* * sublist

MAKE THE INFO CONTENT WANTED TO BE READ, DONT MAKE IT LONG OR SATURATED.

DONT USE LONG CODES, MAKE SURE TO USE SIMPLE CODES NOT WITH LONG NUMBER OF LINES OR LONG LINES (MAX 6 LINES PER CODEBLOCK)
MAX AMOUNT OF INFO SLIDES: 3 (MINIMIZE IT AS MAX)
REPLACE EACH \n by <br>
AND FOR MARKDOWN POINTS, USE * as a point

FOR INLINE SNIPPETS DONT USE the ` + "``" + ` to specify a codeblock, just write it in italic **

DONT YAP A LOT, WRITE SHORT CODE SNIPPETS, VIEWERS LIKE TO SEE CODE MORE THAN TALKING.

WARNING!!!
IN THE INFO SLIDE, YOU MUST GENERATE 300 WORDS, NOT MORE, AND NOT LESS. FOR EACH SLIDE.


MAKE THE TOPIC EASY, MAKE IT so that beginners can understand the content, not only experienced programmers or familliar with programming, DONT GO INTO DEEP DETAILS, JUST TALK GENERATLLY ABOUT A PROGRAMMING LANGUAGE OR A FRAMEWORK, DONT. GO. INTO. DETAILS.
{
  "color": string (USE CSS color hex THAT MATCHES THE TOPIC but a little more dimmed, IT MUST MATCH THE COLOR OF THE PROGRAMMING LANGUAGE/FRAMEWORK, make sure its visible as bakcground for white text),
  "name": string (the topic id that you chose from the list (id. ex: javascript) IT MUST EXACTLY MATCH THE ID THAT YOU PICKED FROM THE PREVIOUS LIST),
  "caption": string (short caption for the post for instagram, Should be long and contain the following hashtags ex. 🤔 What the hell is even Svelte? You keep hearing the name, but no idea what it actually is? Let’s fix that. Svelte is a modern web framework — like React, but cleaner, faster, and way easier to learn ⚡ Perfect for students, beginners, and devs who just want to build without the nonsense. We’re finally BACK and this post kicks off our next workshop season 🔥 Got questions? Confused by all the nerdy talk? We got you. 👉 Join the Discord — link in bio! We’ll explain it all in plain human language 💬 #WTFisSvelte #Svelte #WebFrameworks #CODIHA #codihaclub #codiha #FrontendDev #LearnToCode #PFEReady #ModernWeb #WebDevSimplified) WITH \n don't forget,
  "slides": {
    "hero": {
      "title": string (markdown allowed) MAX OF 7 WORDS (example: WTF is Svelte?, use other words to your own style, like: have you heard about ... ?) [don't use the colon : in the title, make it attractive, use words that attract the user],
      "description": string (markdown allowed) MAX OF 20 WORDS,
      "image": LEAVE EMPTY STRING
    },
    "info": [

      { "title": string (markdown allowed), "content": string (markdown allowed) // for the content don't write a huge chunk of text, use a bit of markdown, lists, new lines to separate paragraphs... etc, YOU MUST GENERATE EXACTLY WHAT YOUVE BEEN TOLD IN THE WARNING PREVIOUSLY) }
	  ...
    ],
    "finish": {
      "description": string (markdown allowed),
      "image": LEAVE EMPTY STRING
    }
  }
}
`

func filterTopics(topics []string, alreadyTalked []string) []string {
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
}

func GeneratePost(ctx context.Context) (postgen.Post, error) {

	var out string

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
		filtered := filterTopics(topics[:], alreadyTalkedAbout)

		// prompt := fmt.Sprintf(AIprompt, strings.Join(filtered, ", "), strings.Join(alreadyTalkedAbout, ", "))
		prompt := fmt.Sprintf(AIprompt, strings.Join(filtered, ", "))

		fmt.Println("Already talked about: ", strings.Join(alreadyTalkedAbout, ", "))
		fmt.Println("Available Topics: ", strings.Join(filtered, ", "))

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
	return p, nil
}
