# blog-aggregator
Welcome to blog-aggregator! (Or *gator* for short). Below you'll find general requirements and instructions to get this CLI up and running.

## Requirements
This CLI has a few requirements:
1. Latest version of PostgreSQL (I created this on v16.8)
2. Latest version of Go (I created this on v1.24.1)

## How to install
To get everything up and running, start by running these commands once you have the requirements above installed to ensure you can have the same toolset I have:
1. "sudo -u postgres psql" (for Mac "psql postgres")
2. "CREATE DATABASE gator;"
3. "\c gator"
4. "ALTER USER postgres PASSWORD *password of your choice*;"
5. "exit" (should be out of the database now)
5. "go install github.com/pressly/goose/v3/cmd/goose@latest"
6. **From the sql/schema directory** "goose postgres postgres://postgres:password@localhost:5432/gator up" (password is whatever you made it)
7. Final step is to run "go install" in the root directory so you can simply run "blog-aggregator" for the program name.

**OPTIONAL** "go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest" (if you want to create your own SQL queries and generate them to Go code)

## Commands
* "blog-aggregator register {username}" (this will register a username of your choice to the database and create a JSON config file called ".gatorconfig.json" at the ~\ location.)
* "blog-aggregator login {username}" (this will update the config file and successfully log you in as long as the username is in the database.)
* "blog-aggregator addfeed {name of feed} {url of feed}" (this is how you can add feeds to aggregate and browse, it automatically has you follow the feeds you add.)
* "blog-aggregator agg {time between requests}" (time is optional) ex: "gator agg 1m30s" (this command runs indefinitely until you hit ctrl+C to stop it)
* "blog-aggregator browse {limit of posts}" (limit is option, defaults to two posts)
* "blog-aggregator follow {url of feed}"
* "blog-aggregator unfollow {url of feed}"
* "blog-aggregator users" (lists out all users and shows currently logged in one)
* "blog-aggregatorr feeds"
* "blog-aggregator following" (lists out all feeds you are following)

Above are all the commands you'll likely use for this CLI, there is also a "blog-aggregator reset" command if you want to wipe the whole database to a clean slate.
