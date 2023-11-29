A simple commandline hacker news client.

Features:
* Browse the front page anonymously (i.e. no login) and sort by new, hot, best
* Search for stories via the Algolia API and sort by date, popularity
* Format output for terminal markdown viewing (via e.g. `kitty`) or csv

Examples:
* Get top 30 stories on the front page:

  ```sh
  hn
  ```
  
* Get newest 50 stories on the front page and output as markdown:

  ```sh
  hn --ranking new --limit 50 --style markdown | mdcat
  ```
  
* Search for stories containing "foobar" ranked by date and format output as csv:

  ```sh
  hn --query "foobar" --ranking date --style csv
  ```

Full CLI:
```
Options:
    -h, --help      show this help message and exit
    -v, --version   show program version information and exit
    -l, --limit     max number of results to fetch (default: 30)
    -s, --style     output style, one of plain|markdown|md|csv (default: plain)
    -r, --ranking   ranking method
                    top|new|best for front page items (default: top)
                    date|popularity for search result items (default: popularity)
    -q, --query     search query

Notes:
	The output for --style=csv is: id,type,by,timestamp,title,url,score,comments
```
