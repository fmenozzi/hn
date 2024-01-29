A simple commandline hacker news client.

Features:
* Browse the front page anonymously (i.e. no login) and sort by new, hot, best
* Search for stories via the Algolia API and sort by date, popularity
* Format output for terminal markdown viewing (via e.g. [`mdcat`](https://github.com/swsnr/mdcat))
    * Markdown via `mdcat` et al only possible on supported terminals (e.g. [`kitty`](https://sw.kovidgoyal.net/kitty/), [`iTerm2`](https://iterm2.com/))
* Format output in json for scripting
    * See [here](https://github.com/HackerNews/API) for API details

Examples:
* Get top 30 stories on the front page:

  ```sh
  hn
  ```
  
* Get newest 50 stories on the front page and output as markdown using `mdcat`:

  ```sh
  hn --ranking new --limit 50 --style markdown | mdcat
  ```
  
* Search for stories containing "foobar" ranked by date and output as json:

  ```sh
  hn --query "foobar" --ranking date --style json
  ```

Full CLI:
```
Options:
    -h, --help      show this help message and exit
    -v, --version   show program version information and exit
    -l, --limit     max number of results to fetch (default: 30)
    -s, --style     output style, one of plain, markdown, md, json (default: plain)
    -r, --ranking   ranking method
                        one of top, new, best for front page items (default: top)
                        one of date, popularity for search result items (default: popularity)
    -q, --query     search query
    -t, --tags      filter search results on specific tags (default: story)

Notes:
    Search tags are ANDed by default but can be ORed if between parentheses. For
    example, "author_pg,(story,poll)" filters on "author_pg AND (type=story OR type=poll)".
    See https://hn.algolia.com/api for more.
```

This code is licensed under the [GNU General Public License version 3](https://www.gnu.org/licenses/gpl-3.0.en.html).
