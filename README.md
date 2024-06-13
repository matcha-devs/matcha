# Matcha

Matcha is the first truly comprehensive, web-based personal finance tracker and budgeting tool.

[![Go](https://github.com/matcha-devs/matcha/actions/workflows/go.yml/badge.svg)](
https://github.com/matcha-devs/matcha/actions/workflows/go.yml)

## Goal features:

* Reliable connections and performance
* Net worth trajectory
* Optimal paycheck allocator
* AI transaction categorizing
* Rewards currency value integration
* Investment balancing
* True value of different tax bucket holdings
* Subscriptions summaries
* Global accounts integration

## Usage

Access our service and create your account at https://www.[domain].com

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.
Please make sure to update tests as appropriate.

After cloning this repository, add ```MYSQL_PASSWORD='your sql password'```
for environment variable in the configuration. Then, run ```go mod tidy```
to import dependencies. For the directory choose main file path ( ```'file path'\matcha\cmd\main```).

Run this to prevent tracking and modifying of the common config file:
```git update-index --skip-worktree .idea/runConfigurations/go_build_github_com_matcha_devs_matcha_cmd.xml```

To test modules, run ```go test ./...``` in the terminal.

## Dependencies

* Web service developed in **Go**

## Authors

* [Seoyoung Cho](https://github.com/seoyoungcho213) **co-owner** (backend)
* [Carlos Cotera](https://github.com/carlosacj55) **co-owner** (frontend)
* [Ali A Shah](https://github.com/alishah634) (backend)
* [Andrea Goh](https://github.com/andreag0101) (frontend)

## License

[LICENSE.md](LICENSE.md)
