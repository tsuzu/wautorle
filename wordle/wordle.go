package wordle

import (
	"fmt"
	"strings"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type Automator struct {
	svc *selenium.Service
	wd  selenium.WebDriver
}

// NewAutomator creates a new Automator for Wordle.
func NewAutomator() (*Automator, error) {

	port := 5050

	svc, err := selenium.NewChromeDriverService("chromedriver", port)

	if err != nil {
		return nil, err
	}

	chrCaps := chrome.Capabilities{
		Args: []string{
			"--headless",
			"--user-agent=Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Mobile/15E148 Safari/604.1",
		},
		Prefs: map[string]interface{}{
			"profile.content_settings.exceptions.clipboard": map[string]interface{}{
				"[*.]nytimes.com,*": map[string]interface{}{
					"last_modified": time.Now().UnixMilli(),
					"setting":       1,
				},
			},
		},
	}

	caps := selenium.Capabilities{"browserName": "chrome"}
	caps.AddChrome(chrCaps)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		svc.Stop()

		return nil, err
	}
	if err := wd.Get("https://www.nytimes.com/games/wordle/index.html"); err != nil {
		svc.Stop()
		wd.Close()

		return nil, err
	}

	return &Automator{
		svc: svc,
		wd:  wd,
	}, nil
}

// Close closes SeleniumUtil.
func (a *Automator) Close() error {
	a.svc.Stop()
	a.wd.Close()

	return nil
}

func (a *Automator) SetStateStats(stateStats string) error {
	_, err := a.wd.ExecuteScript(`
	var arg = JSON.parse(arguments[0]);
	localStorage.setItem('nyt-wordle-state', JSON.stringify(arg.state));
	localStorage.setItem('nyt-wordle-statistics', JSON.stringify(arg.stats));
	`, []interface{}{
		stateStats,
	})

	if err := a.wd.Refresh(); err != nil {
		return err
	}

	return err
}

func (a *Automator) GetStateStats() (stateStats string, err error) {
	r, err := a.wd.ExecuteScript(`return JSON.stringify({
		state: JSON.parse(localStorage.getItem('nyt-wordle-state')),
		stats: JSON.parse(localStorage.getItem('nyt-wordle-statistics')),
	})`, nil)

	if err != nil {
		return "", err
	}

	return r.(string), nil
}

func (a *Automator) Line(i int) (string, error) {
	const script = `return ((line) => {
		const row = document.querySelector('game-app').shadowRoot.querySelectorAll('game-row')[line];
		const eval = row._evaluation;
	
		return Array.from(row._letters).reduce((sum, current, idx) => {
			return sum + (eval[idx] == 'correct' ? 'G' :
				eval[idx] == 'present' ? 'O' :
				'') + current
		}, '');
	})(arguments[0])
	`

	result, err := a.wd.ExecuteScript(script, []interface{}{
		i,
	})

	if err != nil {
		return "", err
	}

	return result.(string), nil
}

func (a *Automator) Screenshot() ([]byte, error) {
	return a.wd.Screenshot()
}

func (a *Automator) Enter(word string) error {
	const script = `
		document.querySelector('game-app').shadowRoot.querySelector('game-keyboard').shadowRoot.querySelector('button[data-key=%c]').click();
	`

	const enterScript = `
		document.querySelector('game-app').shadowRoot.querySelector('game-keyboard').shadowRoot.querySelector('button.one-and-a-half').click()
	`

	scripts := make([]string, len(word))

	for i, c := range []byte(word) {
		scripts[i] = fmt.Sprintf(script, c)
	}

	_, err := a.wd.ExecuteScript(strings.Join(scripts, "\n")+enterScript, nil)

	return err
}

func (a *Automator) Finished() (bool, error) {
	const script = `
		return !!document.querySelector('game-app').shadowRoot.querySelector('game-stats')
	`

	result, err := a.wd.ExecuteScript(script, nil)

	if err != nil {
		return false, err
	}

	return result.(bool), nil
}

func (a *Automator) CopyResult() (string, error) {
	const script = `
		navigator.share = (a) => { window.resultText = a.text; }
		navigator.canShare = () => true;

		document.querySelector('game-app').shadowRoot.querySelector('game-stats').shadowRoot.querySelector('#share-button').click();
		return window.resultText;
	`

	result, err := a.wd.ExecuteScript(script, nil)

	if err != nil {
		return "", err
	}

	return result.(string), nil
}
